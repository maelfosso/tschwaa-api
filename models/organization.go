package models

import (
	"time"

	"gopkg.in/validator.v2"
)

type Organization struct {
	ID          uint64 `json:"id,omitempty"`
	Name        string `json:"name,omitempty" validate:"nonzero,nonnil"`
	Description string `json:"description,omitempty"`
	CreatedBy   int    `json:"createdBy,omitempty" validate:"min=1"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (o Organization) IsValid() bool {
	if err := validator.Validate(o); err != nil {
		return false
	}

	return true
}

type Member struct {
	ID        uint64 `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Sex       string `json:"sex,omitempty"`
	Phone     string `json:"phone_number,omitempty"`
	Joined    bool   `json:"joined,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Adhesion struct {
	ID       uint64 `json:"id,omitempty"`
	MemberID uint64 `json:"member_id,omitempty"`
	OrgID    uint64 `json:"org_id,omitempty"`

	Joined   bool      `json:"joined,omitempty"`
	JoinedAt time.Time `json:"joined_at,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Invitation struct {
	Link         string       `json:"link,omitempty"`
	Active       bool         `json:"active,omitempty"`
	CreatedAt    time.Time    `json:"created_at,omitempty"`
	Adhesion     Adhesion     `json:"adhesion,omitempty"`
	Member       Member       `json:"member,omitempty"`
	Organization Organization `json:"organization,omitempty"`
}