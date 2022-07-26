package model

import (
	"regexp"
	"strings"
	"unicode"
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
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"-"`
	Token     string `json:"-"`

	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
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
	if strings.HasPrefix(u.Phone, "+") {
		phone = u.Phone[1:len(u.Phone)]
	}
	for _, r := range phone {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}
