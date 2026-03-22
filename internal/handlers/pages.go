package handlers

import (
	"log"
	"net/http"

	"github.com/lehmann314159/moshag/internal/auth"
	"github.com/lehmann314159/moshag/internal/db"
)

// homeData is passed to the home page template.
type homeData struct {
	PageData
	Adventures []*db.Adventure
}

// Home renders the home page with the adventure list.
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	adventures, err := h.db.ListAdventures(currentUserID(r))
	if err != nil {
		log.Printf("list adventures: %v", err)
		adventures = nil
	}

	data := homeData{
		PageData:   pageData(r, "MOSHAG — Your Adventures", "home"),
		Adventures: adventures,
	}
	h.render(w, "base", data)
}

// LoginPage renders the login page. If already logged in, redirect to home.
func (h *Handlers) LoginPage(w http.ResponseWriter, r *http.Request) {
	if auth.GetSession(r.Context()) != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	h.render(w, "base", pageData(r, "MOSHAG — Sign In", "login"))
}
