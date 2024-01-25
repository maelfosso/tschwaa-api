package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"tschwaa.com/api/services"
)

func GetCurrentUser(mux chi.Router) {
	mux.Get("/user", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		currentMember := ctx.Value(services.JWTMemberKey)
		if currentMember == nil {
			log.Println("Error No current user: ")
			http.Error(w, "ERR_NO_CURRENT_USER", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(currentMember); err != nil {
			http.Error(w, "error encoding the result", http.StatusBadRequest)
			return
		}
	})
}
