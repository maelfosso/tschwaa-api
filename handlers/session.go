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

type UpdateSessionMembersRequest struct {
	Members []models.OrganizationMember `json:"members"`
}

func UpdateSessionMembers(mux chi.Router) {
	mux.Patch("/members", func(w http.ResponseWriter, r *http.Request) {
		orgIdParam := chi.URLParamFromCtx(r.Context(), "orgID")
		orgId, _ := strconv.ParseUint(orgIdParam, 10, 64)
		log.Println("Get Org ID", orgId)

		sessionIdParam := chi.URLParamFromCtx(r.Context(), "sessionID")
		sessionId, _ := strconv.ParseUint(sessionIdParam, 10, 64)
		log.Println("Get Session ID x2: ", sessionId)

		decoder := json.NewDecoder(r.Body)

		var inputs UpdateSessionMembersRequest
		if err := decoder.Decode(&inputs); err != nil {
			http.Error(w, "error when decoding the session json data", http.StatusBadRequest)
			return
		}
		log.Println("Request - ", inputs)

		// 0. Delete organization'smembers of that session id
		// 1. Check if members is part of the organization: Returning its membership ID
		// 2. Insert the (membership, session) into MembersOfSession

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(true); err != nil {
			log.Println("error when encoding all the organization")
			http.Error(w, "ERR_CREATE_SESSION_102", http.StatusBadRequest)
			return
		}
	})
}

type AddMemberToSessionRequest struct {
	MembershipID uint64 `json:"membership_id"`
}

func AddMemberToSession(mux chi.Router, s addMemberToSession) {
	mux.Post("/members", func(w http.ResponseWriter, r *http.Request) {
		orgIdParam := chi.URLParamFromCtx(r.Context(), "orgID")
		orgId, _ := strconv.ParseUint(orgIdParam, 10, 64)
		log.Println("Get Org ID", orgId)

		sessionIdParam := chi.URLParamFromCtx(r.Context(), "sessionID")
		sessionId, _ := strconv.ParseUint(sessionIdParam, 10, 64)
		log.Println("Get Session ID x2: ", sessionId)

		decoder := json.NewDecoder(r.Body)

		var inputs AddMemberToSessionRequest
		if err := decoder.Decode(&inputs); err != nil {
			http.Error(w, "ERR_ADD_MBSHIP_SESS_101", http.StatusBadRequest)
			return
		}

		// Check if member is part of the organization
		membership, err := s.DoesMembershipConcernOrganization(r.Context(), storage.DoesMembershipConcernOrganizationParams{
			ID:             inputs.MembershipID,
			OrganizationID: orgId,
		})
		if err != nil {
			log.Printf("error when checking if membership[%d] concerns organization[%d]: %s", inputs.MembershipID, orgId, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_102", http.StatusBadRequest)
			return
		}
		if membership == nil {
			log.Printf("error membership[%d] not concern by organization[%d]: %s", inputs.MembershipID, orgId, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_103", http.StatusBadRequest)
			return
		}

		// Add member to session
		mos, err := s.AddMemberToSession(r.Context(), storage.AddMemberToSessionParams{
			MembershipID: inputs.MembershipID,
			SessionID:    sessionId,
		})
		if err != nil {
			log.Printf("error when adding membership[%d] to session[%d]: %s", inputs.MembershipID, sessionId, err)
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
	mux.Delete("/members", func(w http.ResponseWriter, r *http.Request) {
		orgIdParam := chi.URLParamFromCtx(r.Context(), "orgID")
		orgId, _ := strconv.ParseUint(orgIdParam, 10, 64)

		sessionIdParam := chi.URLParamFromCtx(r.Context(), "sessionID")
		sessionId, _ := strconv.ParseUint(sessionIdParam, 10, 64)

		decoder := json.NewDecoder(r.Body)

		var inputs RemoveMemberFromSessionRequest
		if err := decoder.Decode(&inputs); err != nil {
			http.Error(w, "ERR_RMV_MBSHIP_SESS_101", http.StatusBadRequest)
			return
		}

		// Check if member is part of the organization
		membership, err := s.DoesMembershipConcernOrganization(r.Context(), storage.DoesMembershipConcernOrganizationParams{
			ID:             inputs.MembershipID,
			OrganizationID: orgId,
		})
		if err != nil {
			log.Printf("error when checking if membership[%d] concerns organization[%d]: %s", inputs.MembershipID, orgId, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_102", http.StatusBadRequest)
			return
		}
		if membership == nil {
			log.Printf("error membership[%d] not concern by organization[%d]: %s", inputs.MembershipID, orgId, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_103", http.StatusBadRequest)
			return
		}

		// Add member to session
		mos, err := s.RemoveMemberFromSession(r.Context(), storage.RemoveMemberFromSessionParams{
			SessionID     : sessionId,
			OrganizationID: orgId,
			MemberID      : 
		})
		if err != nil {
			log.Printf("error when adding membership[%d] to session[%d]: %s", inputs.MembershipID, sessionId, err)
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
