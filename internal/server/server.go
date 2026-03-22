package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/securecookie"
	"github.com/lehmann314159/moshag/internal/auth"
	"github.com/lehmann314159/moshag/internal/db"
	"github.com/lehmann314159/moshag/internal/handlers"
	"github.com/lehmann314159/moshag/internal/ollama"
	"github.com/lehmann314159/moshag/internal/tables"
)

// Server represents the HTTP server.
type Server struct {
	port         string
	router       *chi.Mux
	templates    *template.Template
	handlers     *handlers.Handlers
	authHandlers *auth.AuthHandlers
}

// New creates a new server instance.
func New(port, ollamaURL, ollamaModel, dbPath string, authCfg auth.AuthConfig, sessionSecret, sessionEncryptKey string) (*Server, error) {
	s := &Server{
		port:   port,
		router: chi.NewRouter(),
	}

	s.parseTemplates()

	// Open database.
	database, err := db.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Session cookie store.
	hashKey := deriveKey(sessionSecret)
	encKey := deriveKey(sessionEncryptKey)
	secureCookie := os.Getenv("SECURE_COOKIES") == "true"
	cookieStore := auth.NewCookieStore(hashKey, encKey, secureCookie)

	// State cookie signer (uses same hash key, no encryption needed).
	stateSC := securecookie.New(hashKey, nil)

	oc := ollama.NewClient(ollamaURL, ollamaModel)
	s.handlers = handlers.New(s.templates, database, oc)
	s.authHandlers = auth.NewAuthHandlers(authCfg, cookieStore, stateSC, s.templates, database)

	// Middleware.
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Compress(5))
	s.router.Use(auth.Middleware(cookieStore))

	s.setupRoutes()

	return s, nil
}

// parseTemplates loads all HTML templates.
func (s *Server) parseTemplates() {
	funcMap := template.FuncMap{
		"safe": func(str string) template.HTML {
			return template.HTML(str)
		},
		"jsonString": func(str string) template.JS {
			b, _ := json.Marshal(str)
			return template.JS(b)
		},
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"seq": func(n int) []int {
			s := make([]int, n)
			for i := range s {
				s[i] = i
			}
			return s
		},
		"stepLabel": tables.StepLabel,
		"stepIndex": func(steps []string, current string) int {
			for i, s := range steps {
				if s == current {
					return i
				}
			}
			return 0
		},
		"dict": func(args ...any) map[string]any {
			m := make(map[string]any, len(args)/2)
			for i := 0; i+1 < len(args); i += 2 {
				m[fmt.Sprintf("%v", args[i])] = args[i+1]
			}
			return m
		},
		"stepIsComplete": func(steps []string, current, check string) bool {
			currentIdx := -1
			checkIdx := -1
			for i, s := range steps {
				if s == current {
					currentIdx = i
				}
				if s == check {
					checkIdx = i
				}
			}
			return checkIdx < currentIdx
		},
	}

	tmpl := template.New("").Funcs(funcMap)
	tmpl = template.Must(tmpl.ParseGlob(filepath.Join("templates", "layouts", "*.html")))
	tmpl = template.Must(tmpl.ParseGlob(filepath.Join("templates", "pages", "*.html")))
	tmpl = template.Must(tmpl.ParseGlob(filepath.Join("templates", "partials", "*.html")))

	s.templates = tmpl
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	return http.ListenAndServe(":"+s.port, s.router)
}

// deriveKey decodes a hex-encoded key or generates a random 32-byte key.
func deriveKey(hexKey string) []byte {
	if hexKey != "" {
		b, err := hex.DecodeString(hexKey)
		if err == nil && len(b) == 32 {
			return b
		}
		log.Printf("Warning: invalid key (expected 64 hex chars), generating random key")
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate random key: " + err.Error())
	}
	return b
}
