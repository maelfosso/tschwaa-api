package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Refresh(mux chi.Router) {
	mux.Post("/refresh", func(w http.ResponseWriter, r *http.Request) {})
}

func IsTokenValid(mux chi.Router) {
	mux.Post("/valid", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(false); err != nil {
			http.Error(w, "error encoding the result", http.StatusBadRequest)
			return
		}
	})
}
