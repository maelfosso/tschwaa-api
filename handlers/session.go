package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"tschwaa.com/api/models"
	"tschwaa.com/api/storage"
)

type createSession interface {
	CreateSessionTx(ctx context.Context, arg storage.CreateSessionParams) (*models.Session, error)
}

type getCurrentSession interface {
	GetCurrentSession(ctx context.Context, organizationID uint64) (*models.Session, error)
}

type CreateSessionRequest struct {
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	OrganizationId uint64 `json:"organization_id"`
}

func CreateSession(mux chi.Router, s createSession) {
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		orgIdParam := chi.URLParamFromCtx(r.Context(), "orgID")
		orgId, _ := strconv.ParseUint(orgIdParam, 10, 64)
		log.Println("Get Org ID x2: ", orgId)

		decoder := json.NewDecoder(r.Body)

		var inputs CreateSessionRequest
		if err := decoder.Decode(&inputs); err != nil {
			http.Error(w, "error when decoding the session json data", http.StatusBadRequest)
			return
		}

		startDate, _ := time.Parse("2006-01-02", inputs.StartDate)
		endDate, _ := time.Parse("2006-01-02", inputs.EndDate)
		session, err := s.CreateSessionTx(r.Context(), storage.CreateSessionParams{
			StartDate:      startDate,
			EndDate:        endDate,
			InProgress:     true,
			OrganizationID: orgId,
		})
		if err != nil {
			log.Println("error when creating a session")
			http.Error(w, "ERR_CREATE_SESSION_101", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(session); err != nil {
			log.Println("error when encoding all the organization")
			http.Error(w, "ERR_CREATE_SESSION_102", http.StatusBadRequest)
			return
		}
	})
}

func GetCurrentSession(mux chi.Router, s getCurrentSession) {
	mux.Get("/current", func(w http.ResponseWriter, r *http.Request) {
		orgIdParam := chi.URLParamFromCtx(r.Context(), "orgID")
		orgId, _ := strconv.ParseUint(orgIdParam, 10, 64)
		log.Println("Get Org ID: ", orgId)

		session, err := s.GetCurrentSession(r.Context(), orgId)
		if err != nil {
			log.Println("error when getting the current session : ", err)
			http.Error(w, "ERR_GET_CURR_SESSION_101", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(session); err != nil {
			log.Println("error when encoding the session information")
			http.Error(w, "ERR_GET_CURR_SESSION_102", http.StatusBadRequest)
			return
		}
	})
}
