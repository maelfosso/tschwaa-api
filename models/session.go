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

// type MembersOfSession struct {
// 	ID           uint64 `json:"id"`
// 	MembershipID uint64 `json:"membership_id"`
// 	SessionID    uint64 `json:"session_id"`

// 	CreatedAt time.Time `json:"created_at,omitempty"`
// 	UpdatedAt time.Time `json:"updated_at,omitempty"`
// }

type MembersOfSession struct {
	ID           *uint64    `json:"id"`
	SessionID    *uint64    `json:"session_id"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	MemberID     uint64     `json:"member_id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Sex          string     `json:"sex"`
	Phone        string     `json:"phone"`
	MembershipID uint64     `json:"membership_id"`
	Position     string     `json:"position"`
	Role         string     `json:"role"`
	Status       string     `json:"status"`
	Joined       bool       `json:"joined"`
	JoinedAt     *time.Time `json:"joined_at"`
}

type SessionPlace struct {
	ID        uint64     `json:"id"`
	Type      string     `json:"type"`
	SessionID uint64     `json:"session_id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type SessionPlaceOnline struct {
	ID             uint64     `json:"id"`
	Type           string     `json:"type"`
	URL            string     `json:"url"`
	SessionPlaceID uint64     `json:"session_place_id"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type SessionPlaceGivenVenue struct {
	ID             uint64     `json:"id"`
	Name           string     `json:"name"`
	Location       string     `json:"location"`
	SessionPlaceID uint64     `json:"session_place_id"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type SessionPlaceMemberHome struct {
	ID             uint64     `json:"id"`
	SessionPlaceID uint64     `json:"session_place_id"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}
