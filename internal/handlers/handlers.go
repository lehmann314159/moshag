package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/lehmann314159/moshag/internal/auth"
	"github.com/lehmann314159/moshag/internal/db"
	"github.com/lehmann314159/moshag/internal/ollama"
)

// Handlers holds dependencies for HTTP handlers.
type Handlers struct {
	templates *template.Template
	db        *db.DB
	ollama    *ollama.Client
}

// New creates a new Handlers instance.
func New(templates *template.Template, database *db.DB, ollamaClient *ollama.Client) *Handlers {
	return &Handlers{
		templates: templates,
		db:        database,
		ollama:    ollamaClient,
	}
}

// render executes a full page template.
func (h *Handlers) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// renderPartial executes a partial template (for HTMX responses).
func (h *Handlers) renderPartial(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("partial template error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// PageData is passed to full page templates.
type PageData struct {
	Title      string
	ActivePage string
	User       *auth.Session
}

// pageData creates a PageData with the session from the request context.
func pageData(r *http.Request, title, activePage string) PageData {
	return PageData{
		Title:      title,
		ActivePage: activePage,
		User:       auth.GetSession(r.Context()),
	}
}
