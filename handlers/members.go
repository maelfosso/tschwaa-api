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
	DoesMembershipExist(ctx context.Context, arg storage.DoesMembershipExistParams) (*models.Membership, error)
	CreateMembership(ctx context.Context, arg storage.CreateMembershipParams) (*models.Membership, error)
	GetInvitationLinkFromMembership(ctx context.Context, membershipID uint64) (string, error)
	CreateInvitationTx(ctx context.Context, arg storage.CreateMembershipInvitationParams) (*models.Organization, error)
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
	Phone   string `json:"phone,omitempty"`
	Invited bool   `json:"invited,omitempty"`
	Error   string `json:"error,omitempty"`
}

func InviteMembersIntoOrganization(mux chi.Router, o inviteMembersIntoOrganization) {
	mux.Post("/members/invite", func(w http.ResponseWriter, r *http.Request) {
		orgIdParam := chi.URLParamFromCtx(r.Context(), "orgID")
		reInvitation, err := strconv.ParseBool(r.URL.Query().Get("reInvitation"))
		if err != nil {
		}
		log.Println("reInvitation : ", reInvitation)

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
				var joinId string

				defer wg.Done()

				log.Println("Start Processing - ", member)
				// Check if a member with the same phone number exist
				existingMember, err := o.GetMemberByPhone(r.Context(), member.Phone)
				if err != nil {
					log.Println("error when checking if member with phone number already exists", err)
					channel <- invitationSentResponse{
						Phone:   member.Phone,
						Invited: false,
						Error:   "ERR_IMIO_510",
					}
					return
				}

				// If member doesn't exist, create it
				if existingMember == nil {
					if reInvitation {
						log.Println("the member should exists", err)
						channel <- invitationSentResponse{
							Phone:   member.Phone,
							Invited: false,
							Error:   "ERR_IMIO_517",
						}
						return
					}

					createdMember, err := o.CreateMember(r.Context(), storage.CreateMemberParams{
						FirstName: member.FirstName,
						LastName:  member.LastName,
						Email:     member.Email,
						Phone:     member.Phone,
						Sex:       member.Sex,
					})
					if err != nil {
						log.Println("error when creating member", err)
						channel <- invitationSentResponse{
							Phone:   member.Phone,
							Invited: false,
							Error:   "ERR_IMIO_511",
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

				membership, err := o.DoesMembershipExist(r.Context(), storage.DoesMembershipExistParams{
					MemberID:       member.ID,
					OrganizationID: org.ID,
				})
				if err != nil {
					log.Printf("error when checking if member[%d] has already joined organization[%d]: %s", member.ID, org.ID, err)
					channel <- invitationSentResponse{
						Phone:   member.Phone,
						Invited: false,
						Error:   "ERR_IMIO_515",
					}
					return
				}
				if reInvitation {
					if membership == nil {
						log.Printf("member[%d] must be a member of organization[%d]", member.ID, org.ID)
						channel <- invitationSentResponse{
							Phone:   member.Phone,
							Invited: false,
							Error:   "ERR_IMIO_518",
						}
						return
					}

					joinId, err = o.GetInvitationLinkFromMembership(r.Context(), membership.ID)
					if err != nil {
						log.Printf("error occurred when getting the invitation link from membership[%d]", membership.ID)
						channel <- invitationSentResponse{
							Phone:   member.Phone,
							Invited: false,
							Error:   "ERR_IMIO_519",
						}
						return
					}
				} else {
					if membership != nil {
						log.Printf("member[%d] already in organization[%d]", member.ID, org.ID)
						channel <- invitationSentResponse{
							Phone:   member.Phone,
							Invited: false,
							Error:   "ERR_IMIO_516",
						}
						return
					}

					joinId = base64.URLEncoding.EncodeToString(sha.Sum(nil)[:])
					_, err = o.CreateInvitationTx(r.Context(), storage.CreateMembershipInvitationParams{
						MemberID:       member.ID,
						OrganizationID: org.ID,
						Joined:         false,
						JoinId:         joinId,
					})
					if err != nil {
						log.Printf("error when invitating member[%d] to join organization[%d]: %s", member.ID, org.ID, err)
						channel <- invitationSentResponse{
							Phone:   member.Phone,
							Invited: true,
							Error:   "ERR_IMIO_514",
						}
						return
					}
				}

				result, err := requests.SendInvitationToJoinOrganization(member, org.Name, joinId, currentMember.FirstName)
				if err != nil {
					log.Println("error when sending a whatsapp invitation to a member", err)
					channel <- invitationSentResponse{
						Phone:   member.Phone,
						Invited: false,
						Error:   err.Error(),
					}
					return
				}
				if len(result.Messages) >= 1 {
					log.Println("invitation successfully sent to member", member.Phone)
					channel <- invitationSentResponse{
						Phone:   member.Phone,
						Invited: true,
						Error:   "",
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

		responses := []invitationSentResponse{}
		for val := range responseChannel {
			log.Println("Channel : ", val)
			responses = append(responses, val)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(responses); err != nil {
			log.Println("error when encoding all the organization")
			http.Error(w, "ERR_IMIO_515", http.StatusBadRequest)
			return
		}
	})
}
