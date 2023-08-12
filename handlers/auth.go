package handlers

import (
	"context"
	"crypto/rand"
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
	ID    uint64 `json:"id,omitempty"`
	Name  string `json:"name",omitempty`
	Email string `json:"email",omitempty`
	Phone string `json:"phone,omitempty"`
	Token string `json:"access_token",omitempty`
}

type JwtClaims struct {
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

func createSecret() (string, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", secret), nil
}

func Signup(mux chi.Router, s authWeb) {
	mux.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
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
		}) // user.Phone, user.Email)
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

func Signin(mux chi.Router, s authWeb) {
	mux.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log.Println("into /signin")

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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if services.IsPasswordMatched(credentials.Password, existingUser.Password) {
			err = fmt.Errorf("the password is not correct: %w", err)
			log.Println("Error CreateUser", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		existingMember, err := s.GetMemberByID(ctx, existingUser.MemberID)
		if err != nil || existingMember == nil {
			err = fmt.Errorf("member related to the user does not exist: %w", err)
			log.Println("Error CreateUser", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var signInResult SignInResult
		signInResult.Name = fmt.Sprintf("%s %s", existingMember.FirstName, existingMember.LastName)
		signInResult.Email = existingMember.Email
		signInResult.ID = existingMember.ID

		tokenString, err := services.GenerateJWTToken(structs.Map(&signInResult))
		if err != nil {
			log.Println("Error CreateUser", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		signInResult.Token = tokenString

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(signInResult); err != nil {
			http.Error(w, "error enconding the result", http.StatusBadRequest)
			return
		}
	})
}

type authMobile interface {
	CreateUserIfNotExists(ctx context.Context, phoneNumber, name string) error
	CreateOTP(ctx context.Context, pinCode models.OTP) error
	SaveOTP(ctx context.Context, pinCode models.OTP) error
	CheckOTP(ctx context.Context, phoneNumber, pinCode string) (*models.OTP, error)
	FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
}

type GetOTPRequest struct {
	PhoneNumber string `json:"phone_number,omitempty"`
	Language    string `json:"language,omitempty"`
}

type CheckOTPRequest struct {
	PhoneNumber string `json:"phone_number,omitempty"`
	Language    string `json:"language,omitempty"`
	PinCode     string `json:"pin_code,omitempty"`
}

func GetOTP(mux chi.Router, a authMobile) {
	mux.Post("/otp", func(w http.ResponseWriter, r *http.Request) {
		var input GetOTPRequest

		decoder := json.NewDecoder(r.Body)

		// extract the phone number
		err := decoder.Decode(&input)
		log.Println("extract the phone number: ", err, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// generate the pin code of 4 digits
		now := time.Now()
		pinCode := helpers.GenerateOTP(now)

		// send the pin code to a the phone number using Whatsapp API
		res, err := SendWoZOTP(
			input.PhoneNumber,
			input.Language,
			pinCode,
		)
		if err != nil {
			log.Println("error when sending the OTP via WhatsApp: ", err)
			http.Error(w, "ERR_COTP_150", http.StatusBadRequest)
			return
		}

		// check if there is an user with this account
		err = a.CreateUserIfNotExists(r.Context(), input.PhoneNumber, "")
		if err != nil {
			log.Println("error when creating the user if he does not exist: ", err)
			http.Error(w, "ERR_COTP_151", http.StatusBadRequest)
			return
		}

		// if not, save the association phone number/pin code in the db
		var m models.OTP
		m.WaMessageId = res.Messages[0].ID
		m.PhoneNumber = input.PhoneNumber
		m.PinCode = pinCode
		m.Active = true

		err = a.CreateOTP(r.Context(), m)
		if err != nil {
			log.Println("error when saving the OTP: ", err)
			http.Error(w, "ERR_COTP_152", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

func CheckOTP(mux chi.Router, a authMobile) {
	mux.Post("/otp/check", func(w http.ResponseWriter, r *http.Request) {
		// read the request body
		var input CheckOTPRequest

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

		// check that the pin code is 6 digit
		var m *models.OTP

		// check that the phone number is correct
		m, err = a.CheckOTP(r.Context(), input.PhoneNumber, input.PinCode)
		if err != nil {
			log.Println("error when checking the otp: ", err)
			http.Error(w, fmt.Sprintf("ERR_COTP_102_%s", err), http.StatusBadRequest)
			return
		}

		m.Active = false
		err = a.SaveOTP(r.Context(), *m)
		if err != nil {
			log.Println("error when changing the active state of the current OTP line: ", err)
			http.Error(w, "ERR_COTP_103", http.StatusBadRequest)
			return
		}

		// Generating the JWT Token
		u, err := a.FindUserByPhoneNumber(r.Context(), input.PhoneNumber)
		if err != nil {
			log.Println("error when looking for user: ", err)
			http.Error(w, "ERR_COTP_104", http.StatusBadRequest)
			return
		}

		var signInResult SignInResult
		signInResult.Name = u.Name
		signInResult.PhoneNumber = u.PhoneNumber

		tokenString, err := services.GenerateJWTToken(structs.Map(signInResult))
		if err != nil {
			log.Println("error when generating jwt token ", err)
			http.Error(w, "ERR_COTP_105", http.StatusBadRequest)
			return
		}

		signInResult.Token = tokenString

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(signInResult); err != nil {
			log.Println("error when encoding auth result: ", err)
			http.Error(w, "ERR_COTP_106", http.StatusBadRequest)
			return
		}
	})
}

func ResendOTP(mux chi.Router, m authMobile) {
	mux.Post("/otp/resend", func(w http.ResponseWriter, r *http.Request) {
		// // read the request body
		// var input CheckOTPRequest

		// // read the request body
		// decoder := json.NewDecoder(r.Body)

		// // extract the phone number and the pin code
		// err := decoder.Decode(&input)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		// // check that the pin code is 6 digit
		// var m *models.OTP

		// // check that the phone number is correct
		// m, err = a.CheckOTP(r.Context(), input.PhoneNumber, input.PinCode)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		// m.Active = false
		// a.SaveOTP(r.Context(), *m)
	})
}
