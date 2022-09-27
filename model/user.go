package model

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var emailAddressMatcher = regexp.MustCompile(
	// Start of string
	`^` +
		// Local part of the address. Note that \x60 is a backtick (`) character.
		`(?P<local>[a-zA-Z0-9.!#$%&'*+/=?^_\x60{|}~-]+)` +
		`@` +
		// Domain of the address
		`(?P<domain>[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)` +
		// End of string
		`$`,
)

type User struct {
	ID        int    `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"-"`
	Token     string `json:"access_token"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (u *User) IsValid() bool {

	if strings.TrimSpace(u.Firstname) == "" ||
		strings.TrimSpace(u.Lastname) == "" ||
		strings.TrimSpace(u.Phone) == "" ||
		strings.TrimSpace(u.Password) == "" ||
		strings.TrimSpace(u.Email) == "" {
		return false
	}

	if !emailAddressMatcher.MatchString(u.Email) {
		return false
	}

	var phone string = u.Phone
	if strings.HasPrefix(phone, "+") {
		phone = phone[1:len(phone)]
	}
	for _, r := range phone {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

func (u *User) HashPassword() error {
	var passwordBytes = []byte(u.Password)

	// Hash password with Bcrypt MinCost
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		return nil
	}

	// Cast the hashedPassword to string
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) IsPasswordMatched(currentPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(currentPassword),
		[]byte(u.Password),
	)
	return err == nil
}

type SignUpCredentials struct {
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

type SignInCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResult struct {
	Name  string `json:"name",omitempty`
	Email string `json:"email",omitempty`
	Token string `json:"access_token",omitempty`
}

type JwtClaims struct {
	User User
	jwt.StandardClaims
}
