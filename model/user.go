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
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	// Country   string `json:"country"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
