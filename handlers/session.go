package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
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

type getMembersOfSession interface {
	ListAllMembersOfSession(ctx context.Context, arg storage.ListAllMembersOfSessionParams) ([]*models.MembersOfSession, error)
}

type updateSessionMembers interface {
	DoesMembershipExist(ctx context.Context, arg storage.DoesMembershipExistParams) (*models.Membership, error)
	UpdateSessionMembersTx(ctx context.Context, arg storage.UpdateSessionMembersParams) error
}

type addMemberToSession interface {
	DoesMembershipConcernOrganization(ctx context.Context, arg storage.DoesMembershipConcernOrganizationParams) (*models.Membership, error)
	AddMemberToSession(ctx context.Context, arg storage.AddMemberToSessionParams) (*models.MembersOfSession, error)
}

type removeMemberFromSession interface {
	DoesMembershipConcernOrganization(ctx context.Context, arg storage.DoesMembershipConcernOrganizationParams) (*models.Membership, error)
	RemoveMemberFromSession(ctx context.Context, arg storage.RemoveMemberFromSessionParams) error
}

type CreateSessionRequest struct {
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	OrganizationId uint64 `json:"organization_id"`
}

func CreateSession(mux chi.Router, s createSession) {
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orgIdParam := chi.URLParamFromCtx(ctx, "orgID")
		orgID, _ := strconv.ParseUint(orgIdParam, 10, 64)
		log.Println("Get Org ID x2: ", orgID)

		decoder := json.NewDecoder(r.Body)

		var inputs CreateSessionRequest
		if err := decoder.Decode(&inputs); err != nil {
			http.Error(w, "error when decoding the session json data", http.StatusBadRequest)
			return
		}

		startDate, _ := time.Parse("2006-01-02", inputs.StartDate)
		endDate, _ := time.Parse("2006-01-02", inputs.EndDate)
		session, err := s.CreateSessionTx(ctx, storage.CreateSessionParams{
			StartDate:      startDate,
			EndDate:        endDate,
			InProgress:     true,
			OrganizationID: orgID,
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
		ctx := r.Context()

		orgIdParam := chi.URLParamFromCtx(ctx, "orgID")
		orgID, _ := strconv.ParseUint(orgIdParam, 10, 64)
		log.Println("Get Org ID: ", orgID)

		session, err := s.GetCurrentSession(ctx, orgID)
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

func GetMembersOfSession(mux chi.Router, svc getMembersOfSession) {
	mux.Get("/members", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orgIdParam := chi.URLParamFromCtx(ctx, "orgID")
		orgID, _ := strconv.ParseUint(orgIdParam, 10, 64)
		log.Println("Get Org ID", orgID)

		sessionIdParam := chi.URLParamFromCtx(ctx, "sessionID")
		sessionID, _ := strconv.ParseUint(sessionIdParam, 10, 64)
		log.Println("Get Session ID x2: ", sessionID)

		mos, err := svc.ListAllMembersOfSession(ctx, storage.ListAllMembersOfSessionParams{
			OrganizationID: orgID,
			SessionID:      sessionID,
		})
		if err != nil {
			log.Printf("error when listing all members of session[%d] of the organization[%d]: %w", sessionID, orgID, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(mos); err != nil {
			log.Println("error when encoding all the organization")
			http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
			return
		}
	})
}

type UpdateSessionMembersRequest struct {
	Members []models.OrganizationMember `json:"members"`
}

type checkingMembershipResponse struct {
	MemberId   uint64
	Membership *models.Membership
}

func UpdateSessionMembers(mux chi.Router, svc updateSessionMembers) {
	mux.Patch("/members", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orgIdParam := chi.URLParamFromCtx(ctx, "orgID")
		orgID, _ := strconv.ParseUint(orgIdParam, 10, 64)
		log.Println("Get Org ID", orgID)

		sessionIdParam := chi.URLParamFromCtx(ctx, "sessionID")
		sessionID, _ := strconv.ParseUint(sessionIdParam, 10, 64)
		log.Println("Get Session ID x2: ", sessionID)

		decoder := json.NewDecoder(r.Body)

		var inputs UpdateSessionMembersRequest
		if err := decoder.Decode(&inputs); err != nil {
			http.Error(w, "error when decoding the session json data", http.StatusBadRequest)
			return
		}
		log.Println("Request - ", inputs)

		// 1. Check if members are part of the organization: Returning their membership IDs
		mapMemberToMembership := make(map[uint64]*models.Membership)
		countNoMembership := 0
		wg := new(sync.WaitGroup)
		wg.Add(len(inputs.Members))
		for _, member := range inputs.Members {
			go func(member models.OrganizationMember, wg *sync.WaitGroup) {
				defer wg.Done()

				membership, err := svc.DoesMembershipExist(ctx, storage.DoesMembershipExistParams{
					MemberID:       member.ID,
					OrganizationID: orgID,
				})
				if err != nil {
					mapMemberToMembership[member.ID] = nil
					countNoMembership = countNoMembership + 1
				} else {
					mapMemberToMembership[member.ID] = membership
				}
			}(member, wg)
		}
		wg.Wait()

		if countNoMembership > 0 {
			log.Printf("error because there are [%d] members who does not belong to the organization [%d]", countNoMembership, orgID)
			http.Error(w, "error because there are membership not members of that organization", http.StatusBadRequest)
			return
		}

		memberships := make([]models.Membership, 0, len(mapMemberToMembership))
		for _, v := range mapMemberToMembership {
			memberships = append(memberships, *v)
		}
		err := svc.UpdateSessionMembersTx(ctx, storage.UpdateSessionMembersParams{
			OrganizationID: orgID,
			SessionID:      sessionID,
			Memberships:    memberships,
		})
		if err != nil {
			log.Printf("error because there are [%d] members who does not belong to the organization [%d]", countNoMembership, orgID)
			http.Error(w, "error when updating members of session", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(true); err != nil {
			log.Println("error when encoding all the organization")
			http.Error(w, "ERR_IMIO_515", http.StatusBadRequest)
			return
		}
	})
}

type AddMemberToSessionRequest struct {
	MembershipID uint64 `json:"membership_id"`
}

func AddMemberToSession(mux chi.Router, s addMemberToSession) {
	mux.Post("/members", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orgIdParam := chi.URLParamFromCtx(ctx, "orgID")
		orgID, _ := strconv.ParseUint(orgIdParam, 10, 64)
		log.Println("Get Org ID", orgID)

		sessionIdParam := chi.URLParamFromCtx(ctx, "sessionID")
		sessionID, _ := strconv.ParseUint(sessionIdParam, 10, 64)
		log.Println("Get Session ID x2: ", sessionID)

		decoder := json.NewDecoder(r.Body)

		var inputs AddMemberToSessionRequest
		if err := decoder.Decode(&inputs); err != nil {
			http.Error(w, "ERR_ADD_MBSHIP_SESS_101", http.StatusBadRequest)
			return
		}

		// Check if member is part of the organization
		membership, err := s.DoesMembershipConcernOrganization(ctx, storage.DoesMembershipConcernOrganizationParams{
			ID:             inputs.MembershipID,
			OrganizationID: orgID,
		})
		if err != nil {
			log.Printf("error when checking if membership[%d] concerns organization[%d]: %s", inputs.MembershipID, orgID, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_102", http.StatusBadRequest)
			return
		}
		if membership == nil {
			log.Printf("error membership[%d] not concern by organization[%d]: %s", inputs.MembershipID, orgID, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_103", http.StatusBadRequest)
			return
		}

		// Add member to session
		mos, err := s.AddMemberToSession(ctx, storage.AddMemberToSessionParams{
			MembershipID: inputs.MembershipID,
			SessionID:    sessionID,
		})
		if err != nil {
			log.Printf("error when adding membership[%d] to session[%d]: %s", inputs.MembershipID, sessionID, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_104", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(mos); err != nil {
			log.Println("error when encoding all the organization")
			http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
			return
		}
	})
}

type RemoveMemberFromSessionRequest struct {
	MembershipID uint64 `json:"membership_id"`
}

func RemoveMemberFromSession(mux chi.Router, s removeMemberFromSession) {
	mux.Delete("/members/{mosID}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orgIdParam := chi.URLParamFromCtx(ctx, "orgID")
		orgID, _ := strconv.ParseUint(orgIdParam, 10, 64)

		sessionIdParam := chi.URLParamFromCtx(ctx, "sessionID")
		sessionID, _ := strconv.ParseUint(sessionIdParam, 10, 64)

		mosIdParam := chi.URLParamFromCtx(ctx, "mosID")
		mosID, _ := strconv.ParseUint(mosIdParam, 10, 64)

		// remove member from session
		err := s.RemoveMemberFromSession(ctx, storage.RemoveMemberFromSessionParams{
			ID:             mosID,
			SessionID:      sessionID,
			OrganizationID: orgID,
		})
		if err != nil {
			log.Printf("error when removing mos[%d] of organization[%d] from session[%d]: %s", mosID, orgID, sessionID, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_104", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(true); err != nil {
			log.Println("error when encoding all the organization")
			http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
			return
		}
	})
}
