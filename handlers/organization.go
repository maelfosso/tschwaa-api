package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"tschwaa.com/api/model"
)

type createorg interface {
	CreateOrganization(ctx context.Context, org model.Organization) (int64, error)
}

type listorg interface {
	ListAllOrganizationFromUser(ctx context.Context, id uint64) ([]model.Organization, error)
}

func CreateOrganization(mux chi.Router, o createorg) {
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var org model.Organization
		if err := decoder.Decode(&org); err != nil {
			http.Error(w, "error when decoding the organization json data", http.StatusBadRequest)
			return
		}

		currentUser := getCurrentUser(r)
		org.CreatedBy = currentUser.ID

		if !org.IsValid() {
			http.Error(w, "error - organization is not valid", http.StatusBadRequest)
			return
		}

		orgId, err := o.CreateOrganization(r.Context(), org)
		if err != nil {
			http.Error(w, "error when creating the organization", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(orgId); err != nil {
			http.Error(w, "error when encoding the request response", http.StatusBadRequest)
			return
		}
	})
}

func ListOrganizations(mux chi.Router, o listorg) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("JWT Claims - ", getCurrentUser(r))
		currentUser := getCurrentUser(r)
		orgs, err := o.ListAllOrganizationFromUser(r.Context(), uint64(currentUser.ID))
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
