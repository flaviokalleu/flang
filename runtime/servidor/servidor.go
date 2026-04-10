package servidor

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/flavio/flang/compiler/ast"
	"github.com/flavio/flang/compiler/lexer"
	"github.com/flavio/flang/compiler/parser"
	authpkg "github.com/flavio/flang/runtime/auth"
	"github.com/flavio/flang/runtime/banco"
	emailpkg "github.com/flavio/flang/runtime/email"
	"github.com/flavio/flang/runtime/httpclient"
	interp "github.com/flavio/flang/runtime/interpreter"
	wa "github.com/flavio/flang/runtime/whatsapp"
)

// Servidor is the embedded Flang web server.
type Servidor struct {
	Program     *ast.Program
	DB          *banco.Banco
	Porta       string
	WS          *WSHub
	WA          *wa.Client
	Auth        *authpkg.Auth
	Email       *emailpkg.Client
	HTTPClient  *httpclient.Client
	Interpreter *interp.Interpreter
	rateLimiter map[string][]time.Time
	rateMu      sync.Mutex
	htmlCache   string
	htmlCacheMu sync.RWMutex
}

// Novo creates a new server.
func Novo(program *ast.Program, db *banco.Banco, porta string) *Servidor {
	return &Servidor{
		Program: program, DB: db, Porta: porta, WS: NewWSHub(),
		rateLimiter: make(map[string][]time.Time),
	}
}

// Iniciar starts the HTTP server.
func (s *Servidor) Iniciar() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.handlePagina)
	mux.HandleFunc("/api/", s.handleAPI)
	mux.HandleFunc("/upload", s.handleUpload)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	mux.HandleFunc("/ws", s.WS.HandleWS)
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth routes
	if s.Auth != nil {
		mux.HandleFunc("/api/login", s.Auth.Login)
		mux.HandleFunc("/api/registro", s.Auth.Registrar)
		mux.HandleFunc("/api/register", s.Auth.Registrar)
		mux.HandleFunc("/api/me", s.Auth.Me)
	}

	// Proxy endpoint for frontend to call external APIs
	mux.HandleFunc("/api/_proxy", s.handleProxy)

	// Scripting endpoints
	mux.HandleFunc("/api/_eval", s.handleEval)
	mux.HandleFunc("/api/_log", s.handleLog)

	// Custom routes
	for _, route := range s.Program.Routes {
		r := route // capture for closure
		mux.HandleFunc(r.Path, func(w http.ResponseWriter, req *http.Request) {
			if r.Method != "" && req.Method != r.Method {
				s.jsonError(w, "método não permitido", http.StatusMethodNotAllowed)
				return
			}
			if s.Interpreter != nil {
				output := s.Interpreter.EvalStatements(r.Handler, nil)
				w.Header().Set("Content-Type", "application/json")
				if len(output) > 0 {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"resultado": output[len(output)-1],
						"output":    output,
					})
				} else {
					json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
				}
			} else {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
			}
		})
	}

	// Custom pages
	for _, page := range s.Program.Pages {
		pg := page
		mux.HandleFunc(pg.Path, func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			title := pg.Title
			if title == "" {
				title = s.Program.System.Name
			}
			pageHTML := fmt.Sprintf(`<!DOCTYPE html><html><head><meta charset="utf-8"><title>%s</title>
<style>body{font-family:system-ui;margin:0;padding:20px}</style></head><body>%s</body></html>`,
				title, pg.Content)
			w.Write([]byte(pageHTML))
		})
	}

	// Apply middleware chain
	var handler http.Handler = mux
	if s.Auth != nil {
		handler = s.Auth.Middleware(handler)
	}
	handler = s.middleware(handler)

	// Periodic rate limiter cleanup
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			s.rateMu.Lock()
			now := time.Now()
			for ip, times := range s.rateLimiter {
				var recent []time.Time
				for _, t := range times {
					if now.Sub(t) < time.Minute {
						recent = append(recent, t)
					}
				}
				if len(recent) == 0 {
					delete(s.rateLimiter, ip)
				} else {
					s.rateLimiter[ip] = recent
				}
			}
			s.rateMu.Unlock()
		}
	}()

	server := &http.Server{
		Addr:           ":" + s.Porta,
		Handler:        handler,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}
	return server.ListenAndServe()
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
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		// Rate limiting for POST requests to API
		if strings.HasPrefix(r.URL.Path, "/api/") && r.Method == http.MethodPost {
			ip := r.RemoteAddr
			s.rateMu.Lock()
			now := time.Now()
			// Clean old entries
			var recent []time.Time
			for _, t := range s.rateLimiter[ip] {
				if now.Sub(t) < time.Minute {
					recent = append(recent, t)
				}
			}
			s.rateLimiter[ip] = append(recent, now)
			count := len(s.rateLimiter[ip])
			s.rateMu.Unlock()
			if count > 100 { // 100 POST requests per minute
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"erro":"Muitas requisições. Tente novamente em breve."}`))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Servidor) getCachedHTML() string {
	s.htmlCacheMu.RLock()
	if s.htmlCache != "" {
		defer s.htmlCacheMu.RUnlock()
		return s.htmlCache
	}
	s.htmlCacheMu.RUnlock()

	html := s.renderHTML()
	s.htmlCacheMu.Lock()
	s.htmlCache = html
	s.htmlCacheMu.Unlock()
	return html
}

func (s *Servidor) handlePagina(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(s.getCachedHTML()))
}

func (s *Servidor) handleEval(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Require authentication for code execution
	if s.Auth != nil {
		role := r.Header.Get("X-User-Role")
		if role != "admin" {
			s.jsonError(w, "Apenas administradores podem executar código", http.StatusForbidden)
			return
		}
	}

	if s.Interpreter == nil {
		http.Error(w, `{"error":"interpreter not initialized"}`, http.StatusInternalServerError)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 64*1024) // 64KB max

	var req struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
		return
	}

	if req.Code == "" {
		http.Error(w, `{"error":"empty code"}`, http.StatusBadRequest)
		return
	}

	// Wrap code in a logica block for parsing
	wrappedCode := "sistema eval\nlogica\n"
	for _, line := range strings.Split(req.Code, "\n") {
		wrappedCode += "  " + line + "\n"
	}

	lex := lexer.New(wrappedCode)
	tokens, err := lex.Tokenize()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  fmt.Sprintf("erro lexico: %s", err),
			"output": []string{},
		})
		return
	}

	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  fmt.Sprintf("erro de parsing: %s", err),
			"output": []string{},
		})
		return
	}

	output := s.Interpreter.EvalStatements(program.Scripts, program.Functions)
	if output == nil {
		output = []string{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"output": output,
	})
}

func (s *Servidor) handleLog(w http.ResponseWriter, r *http.Request) {
	if s.Interpreter == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"logs":[]}`))
		return
	}

	clear := r.URL.Query().Get("clear") == "true"
	logs := s.Interpreter.GetLogs(clear)
	if logs == nil {
		logs = []string{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"logs": logs,
	})
}

func (s *Servidor) handleAPI(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/")
	path = strings.TrimSuffix(path, "/")
	parts := strings.Split(path, "/")

	// Handle /api/_stats
	if parts[0] == "_stats" {
		s.handleStats(w, r)
		return
	}

	// Handle /api/_proxy (already routed via mux, but skip model check)
	if parts[0] == "_proxy" {
		return
	}

	// Handle scripting endpoints routed via mux
	if parts[0] == "_eval" || parts[0] == "_log" {
		return
	}

	modelo := parts[0]
	if modelo == "" {
		s.jsonError(w, "modelo não especificado", http.StatusBadRequest)
		return
	}

	if _, ok := s.DB.Models[modelo]; !ok {
		s.jsonError(w, fmt.Sprintf("modelo '%s' não existe", modelo), http.StatusNotFound)
		return
	}

	// Role-based access control for write operations
	if r.Method != http.MethodGet && s.Auth != nil {
		for _, screen := range s.Program.Screens {
			if strings.ToLower(screen.Name) == modelo || screenMatchesModel(screen, modelo) {
				if screen.Requires != "" && !s.Auth.CheckRole(r, screen.Requires) {
					s.jsonError(w, "Permissão negada", http.StatusForbidden)
					return
				}
				break
			}
		}
	}

	// Handle /api/{model}/export/csv and /api/{model}/export/json
	if len(parts) >= 2 && parts[1] == "export" {
		format := "json"
		if len(parts) >= 3 {
			format = parts[2]
		}
		s.handleExport(w, r, modelo, format)
		return
	}

	// Handle /api/{model}/{id}/restaurar
	if len(parts) == 3 && parts[2] == "restaurar" {
		id, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			s.jsonError(w, "ID inválido", http.StatusBadRequest)
			return
		}
		s.handleRestaurar(w, r, modelo, id)
		return
	}

	// Handle /api/{model}/{id}/{relation} - relationship expansion
	if len(parts) == 3 && parts[2] != "restaurar" && parts[2] != "export" {
		id, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			s.jsonError(w, "ID inválido", http.StatusBadRequest)
			return
		}
		relacao := parts[2]
		items, err := s.DB.BuscarRelacionados(modelo, id, relacao)
		if err != nil {
			s.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.jsonOK(w, items)
		return
	}

	if len(parts) >= 2 && parts[1] != "" {
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
		// Parse query params for pagination, filters, search
		q := r.URL.Query()
		pagina, _ := strconv.Atoi(q.Get("pagina"))
		if pagina == 0 {
			pagina, _ = strconv.Atoi(q.Get("page"))
		}
		limite, _ := strconv.Atoi(q.Get("limite"))
		if limite == 0 {
			limite, _ = strconv.Atoi(q.Get("limit"))
		}
		ordenar := q.Get("ordenar")
		if ordenar == "" {
			ordenar = q.Get("sort")
		}
		ordem := q.Get("ordem")
		if ordem == "" {
			ordem = q.Get("order")
		}
		busca := q.Get("busca")
		if busca == "" {
			busca = q.Get("search")
		}

		// Collect field filters
		filtros := make(map[string]string)
		if model, ok := s.DB.Models[modelo]; ok {
			for _, f := range model.Fields {
				fname := strings.ToLower(f.Name)
				if val := q.Get(fname); val != "" {
					filtros[fname] = val
				}
			}
		}

		params := &banco.ListarParams{
			Pagina: pagina, Limite: limite,
			Ordenar: ordenar, Ordem: ordem,
			Busca: busca, Filtros: filtros,
		}

		items, total, err := s.DB.Listar(modelo, params)
		if err != nil {
			s.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Return with pagination metadata
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", fmt.Sprintf("%d", total))
		json.NewEncoder(w).Encode(items)

	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB max
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
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB max
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

func (s *Servidor) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		s.jsonError(w, "metodo nao permitido", http.StatusMethodNotAllowed)
		return
	}

	// Limit to 32MB
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		s.jsonError(w, "erro ao ler arquivo: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure uploads directory exists
	if err := os.MkdirAll("uploads", 0755); err != nil {
		s.jsonError(w, "erro ao criar diretorio de uploads", http.StatusInternalServerError)
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)

	// Whitelist allowed extensions
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
		".svg": true, ".pdf": true, ".doc": true, ".docx": true, ".xls": true,
		".xlsx": true, ".csv": true, ".txt": true, ".mp4": true, ".mp3": true,
	}
	if !allowedExts[strings.ToLower(ext)] {
		s.jsonError(w, "Tipo de arquivo não permitido", http.StatusBadRequest)
		return
	}

	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	destPath := filepath.Join("uploads", name)

	dst, err := os.Create(destPath)
	if err != nil {
		s.jsonError(w, "erro ao salvar arquivo", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		s.jsonError(w, "erro ao escrever arquivo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"path": "/uploads/" + name, "name": header.Filename})
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

// triggerNotifiers checks and fires WhatsApp/email/other notifications.
func (s *Servidor) triggerNotifiers(triggerType string, modelo string, data map[string]any) {
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

		// Resolve destination
		dest := resolveField(notif.SendTo, data)
		if dest == "" {
			continue
		}

		// Resolve message template (replace {field} with values)
		msg := resolveTemplate(notif.Message, data)

		// Send via WhatsApp
		if notif.Channel == "whatsapp" && s.WA != nil && s.WA.IsConnected() {
			go func(p, m string) {
				if err := s.WA.EnviarMensagem(p, m); err != nil {
					fmt.Printf("[whatsapp] Erro ao enviar: %s\n", err)
				}
			}(dest, msg)
		}

		// Send via Email
		if notif.Channel == "email" && s.Email != nil {
			subject := resolveTemplate(notif.Subject, data)
			if subject == "" {
				subject = "Notificação"
			}
			go func(to, subj, body string) {
				if err := s.Email.EnviarEmail(to, subj, body); err != nil {
					fmt.Printf("[email] Erro ao enviar: %s\n", err)
				}
			}(dest, subject, msg)
		}
	}
}

// handleProxy allows the frontend to call external APIs through the server.
// POST /api/_proxy
// Body: {"method": "GET", "url": "https://...", "body": "..."}
func (s *Servidor) handleProxy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		s.jsonError(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Method string `json:"method"`
		URL    string `json:"url"`
		Body   string `json:"body"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.jsonError(w, "erro ao ler requisição", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		s.jsonError(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		s.jsonError(w, "URL é obrigatória", http.StatusBadRequest)
		return
	}

	// SSRF protection: block private/internal URLs
	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		s.jsonError(w, "URL inválida", http.StatusBadRequest)
		return
	}
	host := parsedURL.Hostname()
	// Block private IP ranges and dangerous hosts
	blockedPrefixes := []string{"127.", "10.", "192.168.", "172.16.", "172.17.", "172.18.", "172.19.", "172.20.", "172.21.", "172.22.", "172.23.", "172.24.", "172.25.", "172.26.", "172.27.", "172.28.", "172.29.", "172.30.", "172.31.", "169.254.", "0."}
	for _, prefix := range blockedPrefixes {
		if strings.HasPrefix(host, prefix) {
			s.jsonError(w, "URL bloqueada por segurança", http.StatusForbidden)
			return
		}
	}
	if host == "localhost" || host == "" || parsedURL.Scheme == "file" {
		s.jsonError(w, "URL bloqueada por segurança", http.StatusForbidden)
		return
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		s.jsonError(w, "Apenas HTTP/HTTPS permitido", http.StatusBadRequest)
		return
	}

	if req.Method == "" {
		req.Method = "GET"
	}

	if s.HTTPClient == nil {
		s.HTTPClient = httpclient.Novo()
	}

	var reqBody []byte
	if req.Body != "" {
		reqBody = []byte(req.Body)
	}

	resp, err := s.HTTPClient.Chamar(req.Method, req.URL, reqBody)
	if err != nil {
		s.jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
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

// handleStats returns record counts and status breakdowns per model.
func (s *Servidor) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.jsonError(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}

	type modelStats struct {
		Count    int64            `json:"count"`
		Statuses map[string]int64 `json:"statuses,omitempty"`
	}

	result := make(map[string]modelStats)
	for name := range s.DB.Models {
		count, _ := s.DB.Contar(name)
		ms := modelStats{Count: count}
		statuses, err := s.DB.ContarPorStatus(name)
		if err == nil && len(statuses) > 0 {
			ms.Statuses = statuses
		}
		result[name] = ms
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleExport exports all records as CSV or JSON.
func (s *Servidor) handleExport(w http.ResponseWriter, r *http.Request, modelo string, format string) {
	if r.Method != http.MethodGet {
		s.jsonError(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}

	items, err := s.DB.ListarTodos(modelo)
	if err != nil {
		s.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	model := s.DB.Models[modelo]
	timestamp := time.Now().Format("2006-01-02")

	switch format {
	case "csv":
		filename := fmt.Sprintf("%s_%s.csv", modelo, timestamp)
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

		writer := csv.NewWriter(w)
		// Write BOM for Excel UTF-8 compatibility
		w.Write([]byte{0xEF, 0xBB, 0xBF})

		// Header row
		headers := []string{"id"}
		for _, f := range model.Fields {
			headers = append(headers, strings.ToLower(f.Name))
		}
		headers = append(headers, "criado_em", "atualizado_em")
		writer.Write(headers)

		// Data rows
		for _, item := range items {
			var row []string
			for _, h := range headers {
				val := ""
				if v, ok := item[h]; ok && v != nil {
					val = fmt.Sprintf("%v", v)
					// Prevent CSV formula injection
					if len(val) > 0 && (val[0] == '=' || val[0] == '+' || val[0] == '-' || val[0] == '@') {
						val = "'" + val
					}
				}
				row = append(row, val)
			}
			writer.Write(row)
		}
		writer.Flush()

	default: // json
		filename := fmt.Sprintf("%s_%s.json", modelo, timestamp)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		json.NewEncoder(w).Encode(items)
	}
}

// screenMatchesModel checks if a screen references the given model via a list component.
func screenMatchesModel(screen *ast.Screen, modelo string) bool {
	for _, comp := range screen.Components {
		if comp.Type == ast.CompList && strings.ToLower(comp.Target) == modelo {
			return true
		}
	}
	return false
}

// handleRestaurar restores a soft-deleted record.
func (s *Servidor) handleRestaurar(w http.ResponseWriter, r *http.Request, modelo string, id int64) {
	if r.Method != http.MethodPut {
		s.jsonError(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}

	item, err := s.DB.Restaurar(modelo, id)
	if err != nil {
		s.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.WS.Broadcast(WSMessage{Type: "restaurar", Model: modelo, ID: id, Data: item})
	s.jsonOK(w, item)
}
