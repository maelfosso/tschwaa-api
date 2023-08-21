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
}

type CreateOrganizationRequest struct {
	Name        string `json:"name,omitempty" validate:"nonzero,nonnil"`
	Description string `json:"description,omitempty"`
}

func CreateOrganization(mux chi.Router, o createOrg) {
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var inputs CreateOrganizationRequest
		if err := decoder.Decode(&inputs); err != nil {
			http.Error(w, "error when decoding the organization json data", http.StatusBadRequest)
			return
		}

		currentMember := GetCurrentMember(r)

		org, err := o.CreateOrganizationWithMembershipTx(r.Context(), storage.CreateOrganizationParams{
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

		fmt.Println("JWT Claims - ", GetCurrentMember(r))
		currentMember := GetCurrentMember(r)
		orgs, err := o.ListOrganizationOfMember(r.Context(), currentMember.ID)
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

func GetOrganization(mux chi.Router, o getOrg) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		orgIdParam := chi.URLParamFromCtx(r.Context(), "orgID")
		orgId, _ := strconv.ParseUint(orgIdParam, 10, 64)
		log.Println("Get Org ID: ", orgId)

		org, err := o.GetOrganization(r.Context(), orgId)
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(org); err != nil {
			http.Error(w, "error when encoding all the organization", http.StatusBadRequest)
			return
		}
	})
}

// func ArticleCtx(next http.Handler) http.Handler {
//   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//     articleID := chi.URLParam(r, "articleID")
//     article, err := dbGetArticle(articleID)
//     if err != nil {
//       http.Error(w, http.StatusText(404), 404)
//       return
//     }
//     ctx := context.WithValue(r.Context(), "article", article)
//     next.ServeHTTP(w, r.WithContext(ctx))
//   })
// }

// func getArticle(w http.ResponseWriter, r *http.Request) {
//   ctx := r.Context()
//   article, ok := ctx.Value("article").(*Article)
//   if !ok {
//     http.Error(w, http.StatusText(422), 422)
//     return
//   }
//   w.Write([]byte(fmt.Sprintf("title:%s", article.Title)))
// }
