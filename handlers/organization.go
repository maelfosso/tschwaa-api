package handlers

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"tschwaa.com/api/model"
	"tschwaa.com/api/requests"
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

type inviteMembersIntoOrganization interface {
	GetOrganization(ctx context.Context, orgId uint64) (*model.Organization, error)
	FindMemberByPhoneNumber(ctx context.Context, phone string) (*model.Member, error)
	CreateMember(ctx context.Context, member model.Member) (uint64, error)
	CreateAdhesion(ctx context.Context, memberId, orgId uint64) (uint64, error)
	CreateInvitation(ctx context.Context, joinId string, adhesionId uint64) (uint64, error)
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
				log.Println("error occured when get information about the organization", err)
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

type invitationSentResponse struct {
	PhoneNumber string `json:"phone_number,omitempty"`
	Invited     bool   `json:"invited,omitempty"`
	Error       string `json:"error,omitempty"`
}

func InviteMembersIntoOrganization(mux chi.Router, o inviteMembersIntoOrganization) {
	mux.Post("/members/invite", func(w http.ResponseWriter, r *http.Request) {
		orgIdParam := chi.URLParamFromCtx(r.Context(), "orgID")
		orgId, _ := strconv.ParseUint(orgIdParam, 10, 64)
		log.Println("Get Org ID: ", orgId)

		decoder := json.NewDecoder(r.Body)
		var members []model.Member
		if err := decoder.Decode(&members); err != nil {
			log.Println("error when decoding the members json data", err)
			http.Error(w, "err-imio-501", http.StatusBadRequest)
			return
		}

		currentUser := getCurrentUser(r)

		org, err := o.GetOrganization(r.Context(), orgId)
		if err != nil {
			log.Println("error occured when get information about the organization", err)
			http.Error(w, "err-imio-502", http.StatusBadRequest)
			return
		}
		wg := new(sync.WaitGroup)
		wg.Add(len(members))

		responseChannel := make(chan invitationSentResponse) //, 1)

		log.Println("***** START PROCESSING MEMBERS : ", len(members))
		for _, member := range members {
			go func(member model.Member, channel chan invitationSentResponse, wg *sync.WaitGroup) {
				defer wg.Done()

				// Check if a member with the same phone number exist
				existingMember, err := o.FindMemberByPhoneNumber(r.Context(), member.Phone)
				if err != nil {
					log.Println("error when checking if member with phone number already exists", err)
					channel <- invitationSentResponse{
						PhoneNumber: member.Phone,
						Invited:     false,
						Error:       "err-imio-510",
					}
					return
				}

				// If member doesn't exist, create it
				if existingMember == nil {
					log.Println("error when creating a member", err)
					memberId, err := o.CreateMember(r.Context(), member)
					if err != nil {
						log.Println("error when creating a member", err)
						channel <- invitationSentResponse{
							PhoneNumber: member.Phone,
							Invited:     false,
							Error:       "err-imio-511",
						}
						return
					}
					member.ID = memberId
				} else {
					member.ID = existingMember.ID
					member.Name = existingMember.Name
					member.Sex = existingMember.Sex
				}

				sha := sha512.New()
				sha.Write([]byte(
					fmt.Sprintf(
						"%d-%d-%d",
						member.ID, orgId, time.Now().UnixNano(),
					),
				))
				joinId := base64.URLEncoding.EncodeToString(sha.Sum(nil)[:])
				result, err := requests.SendInvitationToJoinOrganization(member, org.Name, joinId, currentUser.Firstname)
				if err != nil {
					log.Println("error when sending a whatsapp invitation to a member", err)
					channel <- invitationSentResponse{
						PhoneNumber: member.Phone,
						Invited:     false,
						Error:       "err-imio-512",
					}
					return
				}

				adhesionId, err := o.CreateAdhesion(r.Context(), member.ID, org.ID)
				if err != nil {
					log.Println("error when creating an adhesion to a member", err)
					channel <- invitationSentResponse{
						PhoneNumber: member.Phone,
						Invited:     true,
						Error:       "err-imio-513",
					}
					return
				}

				_, err = o.CreateInvitation(r.Context(), joinId, adhesionId)
				if err != nil {
					log.Println("error when creating an invitation", err)
					channel <- invitationSentResponse{
						PhoneNumber: member.Phone,
						Invited:     true,
						Error:       "err-imio-514",
					}
					return
				}

				if len(result.Messages) >= 1 {
					log.Println("invitation successfully sent to member", member.Phone)
					channel <- invitationSentResponse{
						PhoneNumber: member.Phone,
						Invited:     true,
						Error:       "",
					}
					return
				}
			}(member, responseChannel, wg)
		}

		go func() {
			wg.Wait()
			close(responseChannel)
		}()

		log.Println("------ WAIT FOR MEMBERS TO BE ALL PROCESSED \n\n")
		// wg.Wait()
		// close(responseChannel)
		responses := []invitationSentResponse{}
		for val := range responseChannel {
			log.Println("Channel : ", val)
			responses = append(responses, val)
		}
		// wg.Wait()
		// close(responseChannel)

		log.Println("\n\n**** ALL MEMBERS PROCESSED")

		log.Println("invitation send response : ", responses)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(responses); err != nil {
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
