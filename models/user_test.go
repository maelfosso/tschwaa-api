package models_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/matryer/is"
	"tschwaa.com/api/models"
)

func TestMember_IsValid(t *testing.T) {
	tests := []struct {
		firstname string
		lastname  string
		sex       string
		phone     string
		email     string
		password  string
		valid     bool
	}{
		{"", "doe", "male", "6932", "john.doe@mail.com", "awef", false},
		{"john", "", "male", "693234", "john.doe@mail.com", "aewfw", false},
		{"johb", "doe", "male", "", "hohw.doe@mail.com", "awe", false},
		{"john", "doe", "male", "9023", "john.doe@mail.com", "awe", true},
		{"john", "doe", "male", "69032432", "j", "aw", false},
		{"john", "doe", "male", "69032", "john.doe@gmail.com", "awe", true},
	}

	t.Run("reports valid users", func(t *testing.T) {
		for i, test := range tests {
			t.Run(fmt.Sprint(i), func(t *testing.T) {
				is := is.New(t)
				u := models.Member{0, test.firstname, test.lastname, test.sex, test.phone, test.email, test.password, "", time.Now(), time.Now()}
				is.Equal(test.valid, u.IsValid())
			})
		}
	})
}
