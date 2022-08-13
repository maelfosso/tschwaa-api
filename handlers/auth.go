package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"tschwaa.com/api/model"
)

type signupper interface {
	Signup(ctx context.Context, user model.User) (string, error)
}

func Signup(mux chi.Router, s signupper) {
	mux.Post("/auth/signup", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var user model.User
		if err := decoder.Decode(&user); err != nil {
			http.Error(w, "error decoding the user model", http.StatusBadRequest)
			return
		}

		if !user.IsValid() {
			http.Error(w, "user is invalid", http.StatusBadRequest)
			return
		}

		if _, err := s.Signup(r.Context(), user); err != nil {
			http.Error(w, "error signing user, try again", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(true); err != nil {
			http.Error(w, "error encoding the result", http.StatusBadRequest)
			return
		}
	})
}
