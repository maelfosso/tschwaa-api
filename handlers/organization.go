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
	"tschwaa.com/api/model"
)

type createOrg interface {
	CreateOrganization(ctx context.Context, org model.Organization) (int64, error)
}

type listOrg interface {
	ListAllOrganizationFromUser(ctx context.Context, id uint64) ([]model.Organization, error)
}

type getOrg interface {
	GetOrganization(ctx context.Context, orgId uint64) (*model.Organization, error)
}

type getOrgMembers interface {
	GetOrganizationMembers(ctx context.Context, orgId uint64) ([]model.Member, error)
}

func CreateOrganization(mux chi.Router, o createOrg) {
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

func ListOrganizations(mux chi.Router, o listOrg) {
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

func GetOrganizationMembers(mux chi.Router, o getOrgMembers) {
	mux.Get("/members", func(w http.ResponseWriter, r *http.Request) {
		orgIdParam := chi.URLParamFromCtx(r.Context(), "orgID")
		orgId, _ := strconv.ParseUint(orgIdParam, 10, 64)
		log.Println("Get Org ID: ", orgId)

		members, err := o.GetOrganizationMembers(r.Context(), orgId)
		if err != nil {
			if err == sql.ErrNoRows {
				members = []model.Member{}
			} else {
				http.Error(w, "error occured when get information about the organization", http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(members); err != nil {
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
