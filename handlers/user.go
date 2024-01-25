package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetCurrentUser(mux chi.Router) {
	mux.Get("/user", func(w http.ResponseWriter, r *http.Request) {

	})
}
