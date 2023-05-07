package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"tschwaa.com/api/models"
)

type auth interface {
	// zap.Logger
	Signup(ctx context.Context, member models.Member, user models.User) (string, error)
	Signin(ctx context.Context, credentials models.SignInInputs) (*models.SignInResult, error)
}

func Signup(mux chi.Router, s auth) {
	mux.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var data models.SignUpInputs
		if err := decoder.Decode(&data); err != nil {
			log.Println("error decoding the user model", err)
			http.Error(w, "error decoding the user model", http.StatusBadRequest)
			return
		}

		var member models.Member
		member.FirstName = data.FirstName
		member.LastName = data.LastName
		member.Sex = data.Sex
		member.Phone = data.Phone
		member.Email = data.Email

		var user models.User
		user.Password = data.Password
		user.Phone = data.Phone
		user.Email = data.Email

		if !user.IsValid() {
			// log.Info("Error SignUp", zap.Error(fmt.Errorf("user is invalid")))
			http.Error(w, "user is invalid", http.StatusBadRequest)
			return
		}

		if _, err := s.Signup(r.Context(), member, user); err != nil {
			// log.Info("Error SignUp", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(true); err != nil {
			// log.Info("Error SignUp", zap.Error(err))
			http.Error(w, "error encoding the result", http.StatusBadRequest)
			return
		}
	})
}

func Signin(mux chi.Router, s auth) {
	mux.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
		log.Println("into /signin")

		decoder := json.NewDecoder(r.Body)

		var credenials models.SignInInputs
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
