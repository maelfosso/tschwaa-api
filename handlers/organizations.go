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
)

type createOrg interface {
	CreateOrganization(ctx context.Context, org models.Organization) (uint64, error)
	CreateAdhesion(ctx context.Context, memberId, orgId uint64, joined bool) (uint64, error)
}

type listOrg interface {
	ListAllOrganizationFromMember(ctx context.Context, id uint64) ([]models.Organization, error)
}

type getOrg interface {
	GetOrganization(ctx context.Context, orgId uint64) (*models.Organization, error)
}

func CreateOrganization(mux chi.Router, o createOrg) {
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var org models.Organization
		if err := decoder.Decode(&org); err != nil {
			http.Error(w, "error when decoding the organization json data", http.StatusBadRequest)
			return
		}

		currentMember := GetCurrentMember(r)
		org.CreatedBy = currentMember.ID

		if !org.IsValid() {
			http.Error(w, "error - organization is not valid", http.StatusBadRequest)
			return
		}

		orgId, err := o.CreateOrganization(r.Context(), org)
		if err != nil {
			http.Error(w, "error when creating the organization", http.StatusBadRequest)
			return
		}

		_, err = o.CreateAdhesion(r.Context(), currentMember.ID, uint64(orgId), true)
		if err != nil {
			http.Error(w, "error when creating adhesion", http.StatusBadRequest)
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

func ListOrganizations(mux chi.Router, o listOrg) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("JWT Claims - ", GetCurrentMember(r))
		currentMember := GetCurrentMember(r)
		orgs, err := o.ListAllOrganizationFromMember(r.Context(), uint64(currentMember.ID))
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
