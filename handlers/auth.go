package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fatih/structs"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"tschwaa.com/api/helpers"
	"tschwaa.com/api/models"
	"tschwaa.com/api/requests"
	"tschwaa.com/api/services"
	"tschwaa.com/api/storage"
)

type SignUpInputs struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Sex       string `json:"sex,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

type SignInInputs struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResult struct {
	ID    uint64  `json:"id,omitempty"`
	Name  string  `json:"name",omitempty`
	Email string  `json:"email",omitempty`
	Phone string  `json:"phone,omitempty"`
	Token *string `json:"access_token",omitempty`
}

type JWTClaims struct {
	User models.Member
	jwt.StandardClaims
}

type authWeb interface {
	GetMemberByUsername(ctx context.Context, arg storage.GetMemberByUsernameParams) (*models.Member, error)
	CreateMember(ctx context.Context, arg storage.CreateMemberParams) (*models.Member, error)
	CreateUserWithMemberTx(ctx context.Context, arg storage.CreateUserWithMemberParams) (uint64, error)
	CreateMemberWithAssociatedUserTx(ctx context.Context, arg storage.CreateMemberWithAssociatedUserParams) error
	GetUserByUsername(ctx context.Context, arg storage.GetUserByUsernameParams) (*models.User, error)
	GetMemberByID(ctx context.Context, id uint64) (*models.Member, error)
}

func SignUp(mux chi.Router, s authWeb) {
	mux.Post("/sign-up", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		decoder := json.NewDecoder(r.Body)

		var inputs SignUpInputs
		if err := decoder.Decode(&inputs); err != nil {
			log.Println("error decoding the user model", err)
			http.Error(w, "error decoding the user model", http.StatusBadRequest)
			return
		}

		// Check if a user with the same email exist
		existingMember, err := s.GetMemberByUsername(ctx, storage.GetMemberByUsernameParams{
			Phone: inputs.Phone,
			Email: inputs.Email,
		})
		if err != nil || existingMember != nil {
			err := fmt.Errorf("member with the email/phone already exists: %w", err)
			log.Println("Error GetMemberByUsername", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.CreateMemberWithAssociatedUserTx(ctx, storage.CreateMemberWithAssociatedUserParams{
			FirstName: inputs.FirstName,
			LastName:  inputs.LastName,
			Sex:       inputs.Sex,
			Email:     inputs.Email,
			Phone:     inputs.Phone,
			Password:  inputs.Password,
		})
		if err != nil {
			log.Println("Error CreateMemberWithAssociatedUserTx: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(true); err != nil {
			http.Error(w, "error encoding the result", http.StatusBadRequest)
			return
		}
	})
}

func SignIn(mux chi.Router, s authWeb) {
	mux.Post("/sign-in", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log.Println("into /sign-in")

		decoder := json.NewDecoder(r.Body)

		var credentials SignInInputs
		if err := decoder.Decode(&credentials); err != nil {
			http.Error(w, "error decoding credentials", http.StatusBadRequest)
			return
		}

		existingUser, err := s.GetUserByUsername(ctx, storage.GetUserByUsernameParams{
			Phone: credentials.Username,
			Email: credentials.Username,
		})
		if err != nil || existingUser == nil {
			err = fmt.Errorf("user with that username does not exist: %w", err)
			log.Println("Error CreateUser", zap.Error(err))
			http.Error(w, "user with that username does not exist", http.StatusBadRequest)
			return
		}

		if !helpers.IsPasswordMatched(existingUser.Password, credentials.Password) {
			err = fmt.Errorf("the password is not correct: %w", err)
			log.Println("Error CreateUser", zap.Error(err))
			http.Error(w, "the password is not correct", http.StatusBadRequest)
			return
		}

		existingMember, err := s.GetMemberByID(ctx, existingUser.MemberID)
		if err != nil || existingMember == nil {
			err = fmt.Errorf("member related to the user does not exist: %w", err)
			log.Println("Error CreateUser", zap.Error(err))
			http.Error(w, "member related to the user does not exist", http.StatusBadRequest)
			return
		}

		var signInResult SignInResult
		signInResult.Name = fmt.Sprintf("%s %s", existingMember.FirstName, existingMember.LastName)
		signInResult.Email = existingMember.Email
		signInResult.ID = existingMember.ID

		tokenString, err := services.GenerateJWTToken(structs.Map(&signInResult))
		if err != nil {
			log.Println("Error CreateUser", zap.Error(err))
			http.Error(w, "error when creating token", http.StatusBadRequest)
			return
		}

		signInResult.Token = nil
		http.SetCookie(
			w,
			&http.Cookie{
				Name:     "jwt",
				Value:    tokenString,
				Path:     "/",
				Expires:  time.Now().Add(1 * time.Hour),
				HttpOnly: true,
				Secure:   true,
				MaxAge:   3600,
				SameSite: http.SameSiteLaxMode,
			},
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(signInResult); err != nil {
			http.Error(w, "error encoding the result", http.StatusBadRequest)
			return
		}
	})
}

type authMobile interface {
	DoesUserExist(ctx context.Context, phoneNumber string) (bool, error)
	CreateOTPTx(ctx context.Context, arg storage.CreateOTPParams) (*models.Otp, error)
	DeactivateOTP(ctx context.Context, id uint64) error
	CheckOTP(ctx context.Context, arg storage.CheckOTPParams) (*models.Otp, error)
	GetMemberByPhone(ctx context.Context, phone string) (*models.Member, error)
}

type GetOtpRequest struct {
	Phone    string `json:"phone,omitempty"`
	Language string `json:"language,omitempty"`
}

type CheckOtpRequest struct {
	Phone    string `json:"phone,omitempty"`
	Language string `json:"language,omitempty"`
	PinCode  string `json:"pin_code,omitempty"`
}

func GetOtp(mux chi.Router, a authMobile) {
	mux.Post("/otp", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var input GetOtpRequest

		decoder := json.NewDecoder(r.Body)

		// extract the phone number
		err := decoder.Decode(&input)
		log.Println("extract the phone number: ", err, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check if the user exists as a member (user can't exists its member)
		// if the user doesn't exist then ask him to create an account using the web app
		// or to be invited by the administrator of its organization
		exists, err := a.DoesUserExist(ctx, input.Phone)
		if err != nil {
			log.Println("error when sending the Otp via WhatsApp: ", err)
			http.Error(w, "ERR_COTP_150", http.StatusBadRequest)
			return
		}
		if exists == false {
			log.Println("no user with the phone number: ", input.Phone)
			http.Error(w, "ERR_COTP_151", http.StatusBadRequest)
			return
		}

		// generate the pin code of 4 digits
		now := time.Now()
		pinCode := helpers.GeneratePinCode(now)

		// send the pin code to a the phone number using Whatsapp API
		res, err := requests.SendTschwaaOtp(
			input.Phone,
			input.Language,
			pinCode,
		)
		if err != nil {
			log.Println("error when sending the Otp via WhatsApp: ", err)
			http.Error(w, "ERR_COTP_152", http.StatusBadRequest)
			return
		}

		_, err = a.CreateOTPTx(r.Context(), storage.CreateOTPParams{
			WaMessageID: res.Messages[0].ID,
			Phone:       input.Phone,
			PinCode:     pinCode,
		})
		if err != nil {
			log.Println("error when saving the Otp: ", err)
			http.Error(w, "ERR_COTP_153", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

func CheckOtp(mux chi.Router, a authMobile) {
	mux.Post("/otp/check", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// read the request body
		var input CheckOtpRequest

		// read the request body
		decoder := json.NewDecoder(r.Body)

		// extract the phone number and the pin code
		err := decoder.Decode(&input)

		log.Println("extract the phone number: ", err, input)
		if err != nil {
			log.Println("error when extracting the request body: ", err)
			http.Error(w, "ERR_CTOP_101", http.StatusBadRequest)
			return
		}

		// check that the phone number is correct
		otp, err := a.CheckOTP(ctx, storage.CheckOTPParams{
			Phone:   input.Phone,
			PinCode: input.PinCode,
		})
		if err != nil {
			log.Println("error when checking the otp: ", err)
			http.Error(w, fmt.Sprintf("ERR_COTP_102_%s", err), http.StatusBadRequest)
			return
		}

		// Set otp's active to FALSE
		err = a.DeactivateOTP(ctx, otp.ID)
		if err != nil {
			log.Println("error when changing the active state of the current Otp line: ", err)
			http.Error(w, "ERR_COTP_103", http.StatusBadRequest)
			return
		}

		// Generating the JWT Token
		member, err := a.GetMemberByPhone(ctx, input.Phone)
		if err != nil {
			log.Println("error when looking for user: ", err)
			http.Error(w, "ERR_COTP_104", http.StatusBadRequest)
			return
		}

		var signInResult SignInResult
		signInResult.Name = fmt.Sprintf("%s %s", member.FirstName, member.LastName)
		signInResult.Phone = member.Phone

		tokenString, err := services.GenerateJWTToken(structs.Map(signInResult))
		if err != nil {
			log.Println("error when generating jwt token ", err)
			http.Error(w, "ERR_COTP_105", http.StatusBadRequest)
			return
		}

		signInResult.Token = nil // tokenString
		http.SetCookie(
			w,
			&http.Cookie{
				Name:     "jwt",
				Value:    tokenString,
				Path:     "/",
				Expires:  time.Now().Add(1 * time.Hour),
				HttpOnly: true,
				Secure:   true,
				MaxAge:   3600,
				SameSite: http.SameSiteLaxMode,
			},
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(signInResult); err != nil {
			log.Println("error when encoding auth result: ", err)
			http.Error(w, "ERR_COTP_106", http.StatusBadRequest)
			return
		}
	})
}

func ResendOtp(mux chi.Router, m authMobile) {
	mux.Post("/otp/resend", func(w http.ResponseWriter, r *http.Request) {
		// // read the request body
		// var input CheckOtpRequest

		// // read the request body
		// decoder := json.NewDecoder(r.Body)

		// // extract the phone number and the pin code
		// err := decoder.Decode(&input)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		// // check that the pin code is 6 digit
		// var m *models.Otp

		// // check that the phone number is correct
		// m, err = a.CheckOtp(r.Context(), input.Phone, input.PinCode)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		// m.Active = false
		// a.SaveOtp(r.Context(), *m)
	})
}
