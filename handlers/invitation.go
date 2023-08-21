package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"tschwaa.com/api/models"
	"tschwaa.com/api/storage"
)

const INVITATIION_TIME_OUT_AFTER_DAYS = 1

type getInvitation interface {
	GetInvitation(ctx context.Context, link string) (*models.Invitation, error)
}

type joinOrganization interface {
	GetInvitation(ctx context.Context, link string) (*models.Invitation, error)
	CreateUserWithMemberTx(ctx context.Context, arg storage.CreateUserWithMemberParams) (uint64, error)
	ApprovedInvitationTx(ctx context.Context, link string) error
}

type JoinOrganizationInputs struct {
	ID        uint64 `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Sex       string `json:"sex,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	Code      string `json:"code,omitempty"`
}

type JoinOrganizationResults struct {
	Code         string              `json:"code,omitempty"`
	Active       bool                `json:"active,omitempty"`
	Link         string              `json:"link,omitempty"`
	CreatedAt    time.Time           `json:"created_at,omitempty"`
	Membership   *models.Membership  `json:"membership,omitempty"`
	Member       models.Member       `json:"member,omitempty"`
	Organization models.Organization `json:"organization,omitempty"`
}

func GetInvitation(mux chi.Router, d getInvitation) {
	mux.Get("/{joinId}", func(w http.ResponseWriter, r *http.Request) {
		joinId := chi.URLParamFromCtx(r.Context(), "joinId")
		log.Println("Get Invitation ID: ", joinId)

		currentUser := GetCurrentMember(r)
		log.Println("CURRENT USER : ", currentUser)

		invitation, err := d.GetInvitation(r.Context(), joinId)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("that invitation does not exist; ", err)
				http.Error(w, "ERR_GINV_601", http.StatusBadRequest)
				return
			} else {
				log.Println("error occured when get information about the organization; ", err)
				http.Error(w, "ERR_GINV_602", http.StatusBadRequest)
				return
			}
		}

		log.Println("Invitation: ", invitation)
		log.Println("Member ", invitation.Member)
		log.Println("Current User", currentUser)

		if currentUser != nil {
			if !(invitation.Member.Phone == currentUser.Phone ||
				invitation.Member.Email == currentUser.Email) {
				log.Println("the invited member is not the signed member")
				http.Error(w, "ERR_GINV_606", http.StatusBadRequest)
				return
			}
		}

		// Check if the invitation is active
		if !invitation.Active {
			log.Println("the invitation is no longer active")
			http.Error(w, "ERR_GINV_603", http.StatusBadRequest)
			return
		}

		// Check if the invitation is outdated
		if invitation.CreatedAt.After(time.Now().AddDate(0, 0, INVITATIION_TIME_OUT_AFTER_DAYS)) {
			log.Println("the invitation is outdated")
			http.Error(w, "ERR_GINV_604", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		var result JoinOrganizationResults
		result.Link = invitation.Link
		result.Member = invitation.Member
		result.Organization = invitation.Organization
		if invitation.Member.UserID > 0 {
			result.Membership = nil

			if currentUser != nil {
				// Return member info: S601
				result.Code = "S601"
				// if err := json.NewEncoder(w).Encode(result); err != nil {
				// 	log.Println("error when encoding the successful response; ", err)
				// 	http.Error(w, "ERR_GINV_607", http.StatusBadRequest)
				// 	return
				// }
			} else {
				// Return member info: S602
				result.Code = "S602"
				// if err := json.NewEncoder(w).Encode(result); err != nil {
				// 	log.Println("error when encoding the successfult response; ", err)
				// 	http.Error(w, "ERR_GINV_608", http.StatusBadRequest)
				// 	return
				// }
			}
			data, _ := json.Marshal(result)
			log.Printf("Result - Current User: %s\n", data)
			if err := json.NewEncoder(w).Encode(result); err != nil {
				log.Println("error when encoding the successful response; ", err)
				http.Error(w, "ERR_GINV_607", http.StatusBadRequest)
				return
			}
		} else {
			// result.Membership = &invitation.Membership
			result.CreatedAt = invitation.CreatedAt
			result.Active = invitation.Active
			result.Code = ""
			if err := json.NewEncoder(w).Encode(result); err != nil {
				log.Println("error when encoding the successful response; ", err)
				http.Error(w, "ERR_GINV_605", http.StatusBadRequest)
				return
			}
		}

	})
}

func JoinOrganization(mux chi.Router, j joinOrganization) {
	mux.Post("/{joinId}", func(w http.ResponseWriter, r *http.Request) {
		joinId := chi.URLParamFromCtx(r.Context(), "joinId")
		log.Println("Get Invitation ID: ", joinId)

		_, err := j.GetInvitation(r.Context(), joinId)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("that invitation does not exist; ", err)
				http.Error(w, "ERR_JOIN_601", http.StatusBadRequest)
				return
			} else {
				log.Println("error occured when get information about the organization; ", err)
				http.Error(w, "ERR_JOIN_602", http.StatusBadRequest)
				return
			}
		}
		decoder := json.NewDecoder(r.Body)

		var data JoinOrganizationInputs
		if err := decoder.Decode(&data); err != nil {
			log.Println("error when decoding the organization json data", err)
			http.Error(w, "ERR_JOIN_611", http.StatusBadRequest)
			return
		}

		// Is the user already a member?
		if len(data.Code) > 0 {

		} else {
			var member models.Member
			member.ID = data.ID
			member.FirstName = data.FirstName
			member.LastName = data.LastName
			member.Sex = data.Sex
			member.Phone = data.Phone
			member.Email = data.Email

			// The member already exists so, let's create the user (authentication)
			_, err := j.CreateUserWithMemberTx(r.Context(), storage.CreateUserWithMemberParams{
				Phone:    data.Phone,
				Email:    data.Email,
				Password: data.Password,
				Token:    "",
				MemberID: data.ID,
			})
			if err != nil {
				log.Println("error when creating the user: %w", err)
				http.Error(w, "ERR_JOIN_613", http.StatusBadRequest)
				return
			}
		}

		err = j.ApprovedInvitationTx(r.Context(), joinId)
		if err != nil {
			log.Println("error when closing invitation", err)
			http.Error(w, "ERR_JOIN_613", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(true); err != nil {
			log.Println("error when encoding all the organization; ", err)
			http.Error(w, "ERR_JOIN_619", http.StatusBadRequest)
			return
		}
	})
}
