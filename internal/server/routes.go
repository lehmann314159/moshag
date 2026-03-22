package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// setupRoutes configures all HTTP routes.
func (s *Server) setupRoutes() {
	// Static files.
	fs := http.FileServer(http.Dir("static"))
	s.router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Public pages.
	s.router.Get("/", s.handlers.Home)
	s.router.Get("/auth/login", s.authHandlers.LoginPage)
	s.router.Get("/auth/{provider}", s.authHandlers.BeginAuth)
	s.router.Get("/auth/{provider}/callback", s.authHandlers.Callback)
	s.router.Post("/auth/logout", s.authHandlers.Logout)

	// Adventure routes (auth not required for now).
	s.router.Route("/adventures", func(r chi.Router) {
		r.Post("/new", s.handlers.NewAdventure)
		r.Get("/{id}", s.handlers.ShowAdventure)
		r.Get("/{id}/messages", s.handlers.Messages)
		r.Post("/{id}/start", s.handlers.StartChat)
		r.Post("/{id}/chat", s.handlers.Chat)
		r.Post("/{id}/done", s.handlers.Done)
		r.Post("/{id}/roll", s.handlers.Roll)
		r.Post("/{id}/next", s.handlers.NextStep)
		r.Delete("/{id}", s.handlers.DeleteAdventure)
	})
}
