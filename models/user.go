package models

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
	ID    uint64 `json:"id,omitempty"`
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`

	Password string `json:"-"`
	Token    string `json:"access_token"`

	MemberID uint64 `json:"id,omitempty"`
	Member   Member `json:"member,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (u *User) IsValid() bool {
	if strings.TrimSpace(u.Email) == "" ||
		strings.TrimSpace(u.Phone) == "" ||
		strings.TrimSpace(u.Password) == "" ||
		strings.TrimSpace(u.Email) == "" {
		return false
	}

	if !emailAddressMatcher.MatchString(u.Email) {
		return false
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

type Member struct {
	ID        uint64 `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Sex       string `json:"sex,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`

	UserID uint64 `json:"user_id,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (u *Member) IsValid() bool {
	if strings.TrimSpace(u.FirstName) == "" ||
		strings.TrimSpace(u.LastName) == "" ||
		strings.TrimSpace(u.Sex) == "" ||
		strings.TrimSpace(u.Phone) == "" ||
		// strings.TrimSpace(u.Password) == "" ||
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
	ID    uint64 `json:"name,omitempty"`
	Name  string `json:"name",omitempty`
	Email string `json:"email",omitempty`
	Token string `json:"access_token",omitempty`
}

type JwtClaims struct {
	User Member
	jwt.StandardClaims
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
	Organization Organization `json:"organization,omitempty"`
	Member       Member       `json:"member,omitempty"`
	Link         string       `json:"link,omitempty"`
	Code         string       `json:"code,omitempty"`
}
