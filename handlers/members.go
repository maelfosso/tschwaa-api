package handlers

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"tschwaa.com/api/models"
	"tschwaa.com/api/requests"
	"tschwaa.com/api/storage"
)

type getOrgMembers interface {
	GetMembersFromOrganization(ctx context.Context, organizationID uint64) ([]*models.OrganizationMember, error)
}

type inviteMembersIntoOrganization interface {
	GetOrganization(ctx context.Context, id uint64) (*models.Organization, error)
	GetMemberByPhone(ctx context.Context, phone string) (*models.Member, error)
	CreateMember(ctx context.Context, arg storage.CreateMemberParams) (*models.Member, error)
	CreateAdhesion(ctx context.Context, arg storage.CreateAdhesionParams) (*models.Adhesion, error)
	CreateInvitationTx(ctx context.Context, arg storage.CreateAdhesionInvitationParams) (*models.Organization, error)
}

func GetOrganizationMembers(mux chi.Router, o getOrgMembers) {
	mux.Get("/members", func(w http.ResponseWriter, r *http.Request) {
		orgIdParam := chi.URLParamFromCtx(r.Context(), "orgID")
		orgId, _ := strconv.ParseUint(orgIdParam, 10, 64)
		log.Println("Get Org ID: ", orgId)

		members, err := o.GetMembersFromOrganization(r.Context(), orgId)
		if err != nil {
			log.Println("error occured when get information about the organization", err)
			http.Error(w, "error occured when get information about the organization", http.StatusBadRequest)
			return
		}
		log.Println("GEt Organization Members ", members)

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
		var members []models.Member
		if err := decoder.Decode(&members); err != nil {
			log.Println("error when decoding the members json data", err)
			http.Error(w, "ERR_IMIO_501", http.StatusBadRequest)
			return
		}

		log.Println("member invited", r.Body)
		log.Println(members)
		currentMember := GetCurrentMember(r)

		org, err := o.GetOrganization(r.Context(), orgId)
		if err != nil {
			log.Println("error occured when get information about the organization", err)
			http.Error(w, "ERR_IMIO_502", http.StatusBadRequest)
			return
		}
		wg := new(sync.WaitGroup)
		wg.Add(len(members))

		responseChannel := make(chan invitationSentResponse) //, 1)

		log.Println("***** START PROCESSING MEMBERS : ", len(members))
		for _, member := range members {
			go func(member models.Member, channel chan invitationSentResponse, wg *sync.WaitGroup) {
				defer wg.Done()

				log.Println("Start Processing - ", member)
				// Check if a member with the same phone number exist
				existingMember, err := o.GetMemberByPhone(r.Context(), member.Phone)
				if err != nil {
					log.Println("error when checking if member with phone number already exists", err)
					channel <- invitationSentResponse{
						PhoneNumber: member.Phone,
						Invited:     false,
						Error:       "ERR_IMIO_510",
					}
					return
				}

				// If member doesn't exist, create it
				if existingMember == nil {
					log.Println("error when creating a member", err)
					createdMember, err := o.CreateMember(r.Context(), storage.CreateMemberParams{
						FirstName: member.FirstName,
						LastName:  member.LastName,
						Email:     member.Email,
						Phone:     member.Phone,
						Sex:       member.Sex,
					})
					if err != nil {
						log.Println("error when creating a member", err)
						channel <- invitationSentResponse{
							PhoneNumber: member.Phone,
							Invited:     false,
							Error:       "ERR_IMIO_511",
						}
						return
					}
					member.ID = createdMember.ID
				} else {
					member.ID = existingMember.ID
					member.FirstName = existingMember.FirstName
					member.LastName = existingMember.LastName
					member.Sex = existingMember.Sex
				}

				log.Println("PRocessing Members - ", member)

				sha := sha512.New()
				sha.Write([]byte(
					fmt.Sprintf(
						"%d-%d-%d",
						member.ID, orgId, time.Now().UnixNano(),
					),
				))
				joinId := base64.URLEncoding.EncodeToString(sha.Sum(nil)[:])
				result, err := requests.SendInvitationToJoinOrganization(member, org.Name, joinId, currentMember.FirstName)
				if err != nil {
					log.Println("error when sending a whatsapp invitation to a member", err)
					channel <- invitationSentResponse{
						PhoneNumber: member.Phone,
						Invited:     false,
						Error:       "ERR_IMIO_512",
					}
					return
				}

				_, err = o.CreateInvitationTx(r.Context(), storage.CreateAdhesionInvitationParams{
					MemberID:       member.ID,
					OrganizationID: org.ID,
					Joined:         false,
					JoinId:         joinId,
				})
				if err != nil {
					log.Println("error when creating an invitation", err)
					channel <- invitationSentResponse{
						PhoneNumber: member.Phone,
						Invited:     true,
						Error:       "ERR_IMIO_514",
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

		log.Println("------ WAIT FOR MEMBERS TO BE ALL PROCESSED \n\n")
		go func() {
			wg.Wait()
			log.Println("------ CLOSING RESPONSE CHANNEL.. -----")
			close(responseChannel)
		}()

		// wg.Wait()
		// close(responseChannel)
		responses := []invitationSentResponse{}
		for val := range responseChannel {
			log.Println("Channel : ", val)
			responses = append(responses, val)
		}
		// wg.Wait()
		// close(responseChannel)

		// log.Println("\n\n**** ALL MEMBERS PROCESSED")

		// log.Println("invitation send response : ", responses)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(responses); err != nil {
			log.Println("error when encoding all the organization")
			http.Error(w, "ERR_IMIO_515", http.StatusBadRequest)
			return
		}
	})
}
