package models

import "time"

type Session struct {
	ID             uint64    `json:"id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	InProgress     bool      `json:"in_progress"`
	OrganizationID uint64    `json:"organization_id"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type MembersOfSession struct {
	ID           uint64 `json:"id"`
	MembershipID uint64 `json:"membership_id"`
	SessionID    uint64 `json:"session_id"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
