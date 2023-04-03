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

const INVITATIION_TIME_OUT_AFTER_DAYS = 1

type getInvitation interface {
	GetInvitation(ctx context.Context, link string) (*models.Invitation, error)
}

type joinOrganization interface {
	UpdateMember(ctx context.Context, member models.Member) error
	DisableInvitation(ctx context.Context, link string) (uint64, error)
	ApprovedAdhesion(ctx context.Context, adhesionId uint64) error
}

func GetInvitation(mux chi.Router, d getInvitation) {
	mux.Get("/join/{joinId}", func(w http.ResponseWriter, r *http.Request) {
		joinId := chi.URLParamFromCtx(r.Context(), "joinId")
		log.Println("Get Invitation ID: ", joinId)

		invitation, err := d.GetInvitation(r.Context(), joinId)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("that invitation does not exist; ", err)
				http.Error(w, "ERR_JOIN_601", http.StatusBadRequest)
				return
			} else {
				log.Println("error occured when get information about the organization; ", err)
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
		if invitation.CreatedAt.After(time.Now().AddDate(0, 0, INVITATIION_TIME_OUT_AFTER_DAYS)) {
			log.Println("the invitation is outdated")
			http.Error(w, "ERR_JOIN_604", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(invitation); err != nil {
			log.Println("error when encoding all the organization; ", err)
			http.Error(w, "ERR_IMIO_515", http.StatusBadRequest)
			return
		}

	})
}

func JoinOrganization(mux chi.Router, j joinOrganization) {
	mux.Post("/join/{joinId}", func(w http.ResponseWriter, r *http.Request) {
		joinId := chi.URLParamFromCtx(r.Context(), "joinId")
		log.Println("Get Invitation ID: ", joinId)

		decoder := json.NewDecoder(r.Body)

		var member models.Member
		if err := decoder.Decode(&member); err != nil {
			log.Println("error when decoding the organization json data", err)
			http.Error(w, "ERR_JOIN_611", http.StatusBadRequest)
			return
		}

		if err := j.UpdateMember(r.Context(), member); err != nil {
			log.Println("error when updating the organization", err)
			http.Error(w, "ERR_JOIN_612", http.StatusBadRequest)
			return
		}

		adhesionId, err := j.DisableInvitation(r.Context(), joinId)
		if err != nil {
			log.Println("error when closing invitation", err)
			http.Error(w, "ERR_JOIN_613", http.StatusBadRequest)
			return
		}

		if err := j.ApprovedAdhesion(r.Context(), adhesionId); err != nil {
			log.Println("error when approving adhesion", err)
			http.Error(w, "ERR_JOIN_614", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(true); err != nil {
			log.Println("error when encoding all the organization; ", err)
			http.Error(w, "ERR_JOIN_619", http.StatusBadRequest)
			return
		}
	})
}
