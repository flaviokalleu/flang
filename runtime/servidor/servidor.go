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
)

// Servidor is the embedded Flang web server.
type Servidor struct {
	Program *ast.Program
	DB      *banco.Banco
	Porta   string
	WS      *WSHub
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
