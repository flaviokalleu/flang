package servidor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/flavio/flang/compiler/ast"
	"github.com/flavio/flang/runtime/banco"
	wa "github.com/flavio/flang/runtime/whatsapp"
)

// Servidor is the embedded Flang web server.
type Servidor struct {
	Program *ast.Program
	DB      *banco.Banco
	Porta   string
	WS      *WSHub
	WA      *wa.Client
}

// Novo creates a new server.
func Novo(program *ast.Program, db *banco.Banco, porta string) *Servidor {
	return &Servidor{Program: program, DB: db, Porta: porta, WS: NewWSHub()}
}

// Iniciar starts the HTTP server.
func (s *Servidor) Iniciar() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.handlePagina)
	mux.HandleFunc("/api/", s.handleAPI)
	mux.HandleFunc("/ws", s.WS.HandleWS)

	handler := s.middleware(mux)
	return http.ListenAndServe(":"+s.Porta, handler)
}

func (s *Servidor) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Servidor) handlePagina(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(s.renderHTML()))
}

func (s *Servidor) handleAPI(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/")
	path = strings.TrimSuffix(path, "/")
	parts := strings.SplitN(path, "/", 2)

	modelo := parts[0]
	if modelo == "" {
		s.jsonError(w, "modelo não especificado", http.StatusBadRequest)
		return
	}

	if _, ok := s.DB.Models[modelo]; !ok {
		s.jsonError(w, fmt.Sprintf("modelo '%s' não existe", modelo), http.StatusNotFound)
		return
	}

	if len(parts) == 2 && parts[1] != "" {
		id, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			s.jsonError(w, "ID inválido", http.StatusBadRequest)
			return
		}
		s.handleAPIComID(w, r, modelo, id)
		return
	}

	switch r.Method {
	case http.MethodGet:
		items, err := s.DB.Listar(modelo)
		if err != nil {
			s.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.jsonOK(w, items)

	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.jsonError(w, "erro ao ler dados", http.StatusBadRequest)
			return
		}
		item, err := s.DB.Criar(modelo, body)
		if err != nil {
			s.jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Broadcast via WebSocket
		if id, ok := item["id"]; ok {
			var idInt int64
			switch v := id.(type) {
			case int64:
				idInt = v
			case float64:
				idInt = int64(v)
			}
			s.WS.Broadcast(WSMessage{Type: "criar", Model: modelo, ID: idInt, Data: item})
		}
		// Trigger WhatsApp notifications
		s.triggerNotifiers("criar", modelo, item)
		w.WriteHeader(http.StatusCreated)
		s.jsonOK(w, item)

	default:
		s.jsonError(w, "método não permitido", http.StatusMethodNotAllowed)
	}
}

func (s *Servidor) handleAPIComID(w http.ResponseWriter, r *http.Request, modelo string, id int64) {
	switch r.Method {
	case http.MethodGet:
		item, err := s.DB.Buscar(modelo, id)
		if err != nil {
			s.jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		s.jsonOK(w, item)

	case http.MethodPut:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.jsonError(w, "erro ao ler dados", http.StatusBadRequest)
			return
		}
		item, err := s.DB.Atualizar(modelo, id, body)
		if err != nil {
			s.jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.WS.Broadcast(WSMessage{Type: "atualizar", Model: modelo, ID: id, Data: item})
		s.triggerNotifiers("atualizar", modelo, item)
		s.jsonOK(w, item)

	case http.MethodDelete:
		if err := s.DB.Deletar(modelo, id); err != nil {
			s.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.WS.Broadcast(WSMessage{Type: "deletar", Model: modelo, ID: id})
		w.WriteHeader(http.StatusNoContent)

	default:
		s.jsonError(w, "método não permitido", http.StatusMethodNotAllowed)
	}
}

func (s *Servidor) jsonOK(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *Servidor) jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"erro": msg})
}

// triggerNotifiers checks and fires WhatsApp/other notifications.
func (s *Servidor) triggerNotifiers(triggerType string, modelo string, data map[string]any) {
	if s.WA == nil || !s.WA.IsConnected() {
		return
	}

	for _, notif := range s.Program.Notifiers {
		// Match trigger type and model
		match := false
		switch {
		case notif.Trigger == triggerType && notif.Model == modelo:
			match = true
		case notif.Trigger == triggerType && notif.Model == "":
			match = true
		case notif.Field != "" && notif.Value != "":
			// Conditional: when field equals value
			if val, ok := data[notif.Field]; ok {
				if fmt.Sprintf("%v", val) == notif.Value {
					match = true
				}
			}
		}

		if !match {
			continue
		}

		// Resolve destination phone
		phone := resolveField(notif.SendTo, data)
		if phone == "" {
			continue
		}

		// Resolve message template (replace {field} with values)
		msg := resolveTemplate(notif.Message, data)

		// Send via WhatsApp
		if notif.Channel == "whatsapp" {
			go func(p, m string) {
				if err := s.WA.EnviarMensagem(p, m); err != nil {
					fmt.Printf("[whatsapp] Erro ao enviar: %s\n", err)
				}
			}(phone, msg)
		}
	}
}

// resolveField gets a value from data, supports dotted paths.
func resolveField(field string, data map[string]any) string {
	if field == "" {
		return ""
	}
	// Direct phone number
	if field[0] >= '0' && field[0] <= '9' || field[0] == '+' {
		return field
	}
	// Try direct field
	parts := strings.SplitN(field, ".", 2)
	key := parts[len(parts)-1] // use last part (e.g. "telefone" from "cliente.telefone")
	if val, ok := data[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	if val, ok := data[parts[0]]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

// resolveTemplate replaces {field} placeholders with actual values.
func resolveTemplate(tmpl string, data map[string]any) string {
	result := tmpl
	for key, val := range data {
		result = strings.ReplaceAll(result, "{"+key+"}", fmt.Sprintf("%v", val))
	}
	return result
}
