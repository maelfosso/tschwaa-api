package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"tschwaa.com/api/models"
)

type getInvitation interface {
	GetInvitation(ctx context.Context, link string) (*models.Invitation, error)
}

func GetInvitation(mux chi.Router, d getInvitation) {
	mux.Get("/join/{joinId}", func(w http.ResponseWriter, r *http.Request) {
		joinId := chi.URLParamFromCtx(r.Context(), "joinId")
		// joinId, _ := strconv.ParseUint(joinIdParam, 10, 64)
		log.Println("Get Invitation ID: ", joinId)

		invitation, err := d.GetInvitation(r.Context(), joinId)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("that organization does not exist")
				http.Error(w, "ERR_JOIN_601", http.StatusBadRequest)
				return
			} else {
				log.Println("error occured when get information about the organization")
				http.Error(w, "ERR_JOIN_602", http.StatusBadRequest)
				return
			}
		}

		// Check if the invitation is active
		if !invitation.Active {
			log.Println("the invitation is no longer active")
			http.Error(w, "ERR_JOIN_603", http.StatusBadRequest)
			return
		}

		// Check if the invitation is outdated
		if invitation.CreatedAt.After(time.Now()) {
			log.Println("the invitation is outdated")
			http.Error(w, "ERR_JOIN_604", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(invitation); err != nil {
			log.Println("error when encoding all the organization")
			http.Error(w, "ERR_IMIO_515", http.StatusBadRequest)
			return
		}

	})
}

func JoinOrganization(mux chi.Router) {
	mux.Post("/join/{joinId}", func(w http.ResponseWriter, r *http.Request) {

	})
}
