package model

import (
	"time"

	"gopkg.in/validator.v2"
)

type Organization struct {
	ID          int    `json:"id,omitempty"`
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
