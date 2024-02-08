package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"tschwaa.com/api/models"
	"tschwaa.com/api/services"
)

type GetCurrentUserResult struct {
	ID    uint64 `json:"id,omitempty"`
	Name  string `json:"name",omitempty`
	Email string `json:"email",omitempty`
	Phone string `json:"phone,omitempty"`
}

func GetCurrentUser(mux chi.Router) {
	mux.Get("/user", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		currentMember := ctx.Value(services.JWTMemberKey).(*models.Member)
		if currentMember == nil {
			log.Println("Error No current user: ")
			http.Error(w, "ERR_NO_CURRENT_USER", http.StatusBadRequest)
			return
		}

		var currentUserResult GetCurrentUserResult
		currentUserResult.Name = fmt.Sprintf("%s %s", currentMember.FirstName, currentMember.LastName)
		currentUserResult.Email = currentMember.Email
		currentUserResult.Phone = currentMember.Phone
		currentUserResult.ID = currentMember.ID

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(currentUserResult); err != nil {
			http.Error(w, "error encoding the result", http.StatusBadRequest)
			return
		}
	})
}
