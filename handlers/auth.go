package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"tschwaa.com/api/models"
)

type auth interface {
	Signup(ctx context.Context, user models.User) (string, error)
	Signin(ctx context.Context, credentials models.SignInCredentials) (*models.SignInResult, error)
}

func Signup(mux chi.Router, s auth, log *zap.Logger) {
	mux.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var data models.SignUpCredentials
		if err := decoder.Decode(&data); err != nil {
			http.Error(w, "error decoding the user model", http.StatusBadRequest)
			return
		}

		var user models.User
		user.Firstname = data.Firstname
		user.Lastname = data.Lastname
		user.Password = data.Password
		user.Phone = data.Phone
		user.Email = data.Email

		if !user.IsValid() {
			log.Info("Error SignUp", zap.Error(fmt.Errorf("user is invalid")))
			http.Error(w, "user is invalid", http.StatusBadRequest)
			return
		}

		if _, err := s.Signup(r.Context(), user); err != nil {
			log.Info("Error SignUp", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(true); err != nil {
			log.Info("Error SignUp", zap.Error(err))
			http.Error(w, "error encoding the result", http.StatusBadRequest)
			return
		}
	})
}

func Signin(mux chi.Router, s auth) {
	mux.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
		log.Println("into /signin")

		decoder := json.NewDecoder(r.Body)

		var credenials models.SignInCredentials
		if err := decoder.Decode(&credenials); err != nil {
			http.Error(w, "error decoding credentials", http.StatusBadRequest)
			return
		}

		user, err := s.Signin(r.Context(), credenials)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "error enconding the result", http.StatusBadRequest)
			return
		}
	})
}
