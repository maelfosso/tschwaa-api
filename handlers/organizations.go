package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"tschwaa.com/api/models"
	"tschwaa.com/api/storage"
)

type createOrg interface {
	CreateOrganizationWithMembershipTx(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error)
}

type listOrg interface {
	ListOrganizationOfMember(ctx context.Context, memberID uint64) ([]*models.Organization, error)
}

type getOrg interface {
	GetOrganization(ctx context.Context, id uint64) (*models.Organization, error)
	GetCurrentSession(ctx context.Context, organizationID uint64) (*models.Session, error)
}

type CreateOrganizationRequest struct {
	Name        string `json:"name,omitempty" validate:"nonzero,nonnil"`
	Description string `json:"description,omitempty"`
}

func CreateOrganization(mux chi.Router, o createOrg) {
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		decoder := json.NewDecoder(r.Body)

		var inputs CreateOrganizationRequest
		if err := decoder.Decode(&inputs); err != nil {
			http.Error(w, "error when decoding the organization json data", http.StatusBadRequest)
			return
		}

		currentMember := GetCurrentMember(r)

		org, err := o.CreateOrganizationWithMembershipTx(ctx, storage.CreateOrganizationParams{
			Name:        inputs.Name,
			Description: &inputs.Description,
			CreatedBy:   &currentMember.ID,
		})
		if err != nil {
			http.Error(w, "error when creating the organization", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(org.ID); err != nil {
			http.Error(w, "error when encoding the request response", http.StatusBadRequest)
			return
		}
	})
}

func ListOrganizations(mux chi.Router, o listOrg) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		fmt.Println("JWT Claims - ", GetCurrentMember(r))
		currentMember := GetCurrentMember(r)
		orgs, err := o.ListOrganizationOfMember(ctx, currentMember.ID)
		if err != nil {
			http.Error(w, "error occured when fetching the organizations", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(orgs); err != nil {
			http.Error(w, "error when encoding all the organizations", http.StatusBadRequest)
			return
		}
	})
}

type GetOrganizationResponse struct {
	Organization   *models.Organization `json:"organization"`
	CurrentSession *models.Session      `json:"current_session"`
}

func GetOrganization(mux chi.Router, o getOrg) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orgIdParam := chi.URLParamFromCtx(ctx, "orgID")
		orgId, _ := strconv.ParseUint(orgIdParam, 10, 64)
		log.Println("Get Org ID: ", orgId)

		org, err := o.GetOrganization(ctx, orgId)
		log.Println("Get Organization ", org, err)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "that organization does not exist", http.StatusBadRequest)
				return
			} else {
				http.Error(w, "error occured when get information about the organization", http.StatusBadRequest)
				return
			}
		}

		currentSession, err := o.GetCurrentSession(ctx, org.ID)
		response := GetOrganizationResponse{
			Organization:   org,
			CurrentSession: currentSession,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "error when encoding all the organization", http.StatusBadRequest)
			return
		}
	})
}
