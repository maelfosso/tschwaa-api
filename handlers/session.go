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
	"tschwaa.com/api/common"
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
	DoesMembershipConcernOrganization(ctx context.Context, arg storage.DoesMembershipConcernOrganizationParams) (*models.Membership, error)
	UpdateSessionMembersTx(ctx context.Context, arg storage.UpdateSessionMembersParams) ([]*models.MembersOfSession, error)
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
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
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
	MembershipIDs []uint64 `json:"membership_ids"`
}

func UpdateSessionMembers(mux chi.Router, svc updateSessionMembers) {
	mux.Patch("/", func(w http.ResponseWriter, r *http.Request) {
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

		memberships := make([]models.Membership, 0, len(inputs.MembershipIDs))
		countNoMembership := 0
		wg := new(sync.WaitGroup)
		wg.Add(len(inputs.MembershipIDs))
		for _, membershipID := range inputs.MembershipIDs {
			go func(membershipID uint64, wg *sync.WaitGroup) {
				defer wg.Done()

				membership, err := svc.DoesMembershipConcernOrganization(ctx, storage.DoesMembershipConcernOrganizationParams{
					ID:             membershipID,
					OrganizationID: orgID,
				})
				if err != nil {
					countNoMembership = countNoMembership + 1
				} else {
					memberships = append(memberships, *membership)
				}
			}(membershipID, wg)
		}
		wg.Wait()

		if countNoMembership > 0 {
			log.Printf("error because there are [%d] members who does not belong to the organization [%d]", countNoMembership, orgID)
			http.Error(w, "error because there are membership not members of that organization", http.StatusBadRequest)
			return
		}

		mos, err := svc.UpdateSessionMembersTx(ctx, storage.UpdateSessionMembersParams{
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
		if err := json.NewEncoder(w).Encode(mos); err != nil {
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
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
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

		membershipID := inputs.MembershipID

		membership, err := s.DoesMembershipConcernOrganization(ctx, storage.DoesMembershipConcernOrganizationParams{
			ID:             membershipID,
			OrganizationID: orgID,
		})
		if err != nil {
			log.Printf("error when checking if membership[%d] concerns organization[%d]: %s", membershipID, orgID, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_102", http.StatusBadRequest)
			return
		}
		if membership == nil {
			log.Printf("error membership[%d] not concern by organization[%d]: %s", membershipID, orgID, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_103", http.StatusBadRequest)
			return
		}

		mos, err := s.AddMemberToSession(ctx, storage.AddMemberToSessionParams{
			MembershipID: membershipID,
			SessionID:    sessionID,
		})
		if err != nil {
			log.Printf("error when adding membership[%d] to session[%d]: %s", membershipID, sessionID, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_104", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(mos.ID); err != nil {
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
	mux.Delete("/{mosID}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orgIdParam := chi.URLParamFromCtx(ctx, "orgID")
		orgID, _ := strconv.ParseUint(orgIdParam, 10, 64)

		sessionIdParam := chi.URLParamFromCtx(ctx, "sessionID")
		sessionID, _ := strconv.ParseUint(sessionIdParam, 10, 64)

		mosIdParam := chi.URLParamFromCtx(ctx, "mosID")
		mosID, _ := strconv.ParseUint(mosIdParam, 10, 64)

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

type getPlaceOfSession interface {
	GetSession(ctx context.Context, arg storage.GetSessionParams) (*models.Session, error)
	GetSessionPlaceTx(ctx context.Context, sessionID uint64) (models.ISessionPlace, error)
}

func GetPlaceOfSession(mux chi.Router, svc getPlaceOfSession) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orgIdParam := chi.URLParamFromCtx(ctx, "orgID")
		orgID, _ := strconv.ParseUint(orgIdParam, 10, 64)

		sessionIdParam := chi.URLParamFromCtx(ctx, "sessionID")
		sessionID, _ := strconv.ParseUint(sessionIdParam, 10, 64)

		session, err := svc.GetSession(ctx, storage.GetSessionParams{
			OrganizationID: orgID,
			SessionID:      sessionID,
		})
		if err != nil {
			log.Printf("error when listing all members of session[%d] of the organization[%d]: %w", sessionID, orgID, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
			return
		}
		if session == nil {
			log.Printf("session does not exist of session[%d] of the organization[%d]: %w", sessionID, orgID, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
			return
		}

		iSessionPlace, err := svc.GetSessionPlaceTx(ctx, sessionID)
		if err != nil {
			log.Printf("error when getting real session place of session[%d] of the organization[%d]: %w", sessionID, orgID, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(iSessionPlace); err != nil {
			log.Println("error when encoding all the organization")
			http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
			return
		}
	})
}

type updatePlaceOfSession interface {
	GetSessionPlace(ctx context.Context, arg storage.GetSessionPlaceParams) (*models.SessionPlace, error)
	GetSessionPlaceOnline(ctx context.Context, arg storage.GetSessionPlaceOnlineParams) (*models.SessionPlacesOnline, error)
	GetSessionPlaceGivenVenue(ctx context.Context, arg storage.GetSessionPlaceGivenVenueParams) (*models.SessionPlacesGivenVenue, error)
	GetSessionPlaceMemberHome(ctx context.Context, arg storage.GetSessionPlaceMemberHomeParams) (*models.SessionPlacesMemberHome, error)
	UpdateSessionPlaceOnline(ctx context.Context, arg storage.UpdateSessionPlaceOnlineParams) (*models.SessionPlacesOnline, error)
	UpdateSessionPlaceGivenVenue(ctx context.Context, arg storage.UpdateSessionPlaceGivenVenueParams) (*models.SessionPlacesGivenVenue, error)
}

type UpdatePlaceOfSessionRequest struct {
	ID       uint64 `json:"id"`
	Type     string `json:"type,omitempty"`
	Link     string `json:"Link,omitempty"`
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`
}

func UpdatePlaceOfSession(mux chi.Router, svc updatePlaceOfSession) {
	mux.Patch("/{sessionPlaceID}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orgIdParam := chi.URLParamFromCtx(ctx, "orgID")
		orgID, _ := strconv.ParseUint(orgIdParam, 10, 64)

		sessionIdParam := chi.URLParamFromCtx(ctx, "sessionID")
		sessionID, _ := strconv.ParseUint(sessionIdParam, 10, 64)

		sessionPlaceIdParam := chi.URLParamFromCtx(ctx, "sessionPlaceID")
		sessionPlaceID, _ := strconv.ParseUint(sessionPlaceIdParam, 10, 64)

		sessionPlace, err := svc.GetSessionPlace(ctx, storage.GetSessionPlaceParams{
			ID:        sessionPlaceID,
			SessionID: sessionID,
		})
		if err != nil {
			log.Printf("error when listing all members of session[%d] of the organization[%d]: %w", sessionID, orgID, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
			return
		}
		if sessionPlace == nil {
			log.Printf("there is no session place with id[%d] from session[%d] of the organization[%d]", sessionPlaceID, sessionID, orgID)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)

		var inputs UpdatePlaceOfSessionRequest
		if err := decoder.Decode(&inputs); err != nil {
			http.Error(w, "error when decoding the session json data", http.StatusBadRequest)
			return
		}

		if sessionPlace.PlaceType == common.SESSION_PLACE_ONLINE {
			sessionPlaceOnline, err := svc.GetSessionPlaceOnline(ctx, storage.GetSessionPlaceOnlineParams{
				ID:             inputs.ID,
				SessionPlaceID: sessionPlaceID,
			})
			if err != nil {
				log.Printf(
					"error when getting online session place[%d] of session place [%d] from session[%d] of the organization[%d]: %w",
					inputs.ID, sessionPlaceID, sessionID, orgID, err,
				)
				http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
				return
			}
			if sessionPlaceOnline == nil {
				log.Printf(
					"there is no online session place with id[%d] of session place [%d] from session[%d] of the organization[%d]",
					inputs.ID, sessionPlaceID, sessionID, orgID,
				)
				http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
				return
			}

			sessionPlaceOnline, err = svc.UpdateSessionPlaceOnline(ctx, storage.UpdateSessionPlaceOnlineParams{
				ID:             sessionPlaceOnline.ID,
				SessionPlaceID: sessionPlaceID,
				Type:           inputs.Type,
				Link:           inputs.Link,
			})
			if err != nil {
				log.Printf(
					"error when updating online session place[%d] of session place [%d] from session[%d] of the organization[%d]: %w",
					inputs.ID, sessionPlaceID, sessionID, orgID, err,
				)
				http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
				return
			}
		} else if sessionPlace.PlaceType == common.SESSION_PLACE_GIVEN_VENUE {
			sessionPlaceGivenVenue, err := svc.GetSessionPlaceGivenVenue(ctx, storage.GetSessionPlaceGivenVenueParams{
				ID:             inputs.ID,
				SessionPlaceID: sessionPlaceID,
			})
			if err != nil {
				log.Printf(
					"error when getting online session place[%d] of session place [%d] from session[%d] of the organization[%d]: %w",
					inputs.ID, sessionPlaceID, sessionID, orgID, err,
				)
				http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
				return
			}
			if sessionPlaceGivenVenue == nil {
				log.Printf(
					"there is no online session place with id[%d] of session place [%d] from session[%d] of the organization[%d]",
					inputs.ID, sessionPlaceID, sessionID, orgID,
				)
				http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
				return
			}

			sessionPlaceGivenVenue, err = svc.UpdateSessionPlaceGivenVenue(ctx, storage.UpdateSessionPlaceGivenVenueParams{
				ID:             sessionPlaceGivenVenue.ID,
				SessionPlaceID: sessionPlaceID,
				Name:           inputs.Name,
				Location:       inputs.Location,
			})
			if err != nil {
				log.Printf(
					"error when updating given venue session place[%d] of session place [%d] from session[%d] of the organization[%d]: %w",
					inputs.ID, sessionPlaceID, sessionID, orgID, err,
				)
				http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
				return
			}
		} else if sessionPlace.PlaceType == common.SESSION_PLACE_MEMBER_HOME {
			sessionPlaceMemberHone, err := svc.GetSessionPlaceMemberHome(ctx, storage.GetSessionPlaceMemberHomeParams{
				ID:             inputs.ID,
				SessionPlaceID: sessionPlaceID,
			})
			if err != nil {
				log.Printf(
					"error when getting online session place[%d] of session place [%d] from session[%d] of the organization[%d]: %w",
					inputs.ID, sessionPlaceID, sessionID, orgID, err,
				)
				http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
				return
			}
			if sessionPlaceMemberHone == nil {
				log.Printf(
					"there is no online session place with id[%d] of session place [%d] from session[%d] of the organization[%d]",
					inputs.ID, sessionPlaceID, sessionID, orgID,
				)
				http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
				return
			}
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

type changePlaceOfSession interface {
	ChangeSessionPlaceTx(ctx context.Context, arg storage.ChangeSessionPlaceParams) (models.ISessionPlace, error)
}

type ChangePlaceOfSessionRequest struct {
	SessionPlaceType string `json:"place_type,omitempty"`

	// Online
	Type string `json:"type,omitempty"`
	Link string `json:"link,omitempty"`

	// Given Venue
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`

	// Member Home
	Choice string `json:"choice,omitempty"`
}

func ChangePlaceOfSession(mux chi.Router, svc changePlaceOfSession) {
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orgIdParam := chi.URLParamFromCtx(ctx, "orgID")
		orgID, _ := strconv.ParseUint(orgIdParam, 10, 64)

		sessionIdParam := chi.URLParamFromCtx(ctx, "sessionID")
		sessionID, _ := strconv.ParseUint(sessionIdParam, 10, 64)
		log.Println("OrgId, sessionId: ", orgID, sessionID)
		decoder := json.NewDecoder(r.Body)

		var inputs ChangePlaceOfSessionRequest
		if err := decoder.Decode(&inputs); err != nil {
			http.Error(w, "error when decoding the session json data", http.StatusBadRequest)
			return
		}
		log.Println("ChangePoSRequest: ", inputs.SessionPlaceType)
		log.Println(inputs.Type, inputs.Link)
		log.Println(inputs.Name, inputs.Location)
		log.Println(inputs.Choice)

		iSessionPlace, err := svc.ChangeSessionPlaceTx(ctx, storage.ChangeSessionPlaceParams{
			SessionID:        sessionID,
			SessionPlaceType: inputs.SessionPlaceType,

			Type: &inputs.Type,
			Link: &inputs.Link,

			Name:     &inputs.Name,
			Location: &inputs.Location,

			// Choice: &inputs.Choice,
		})
		if err != nil {
			log.Printf("error when completely changing a session place of session[%d] of the organization[%d]: %w", sessionID, orgID, err)
			http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(iSessionPlace); err != nil {
			log.Println("error when encoding all the organization")
			http.Error(w, "ERR_ADD_MBSHIP_SESS_105", http.StatusBadRequest)
			return
		}
	})
}
