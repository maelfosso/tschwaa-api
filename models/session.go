package models

import (
	"fmt"
	"time"

	"tschwaa.com/api/common"
)

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
	PlaceType string     `json:"place_type"`
	SessionID uint64     `json:"session_id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type ISessionPlace interface {
	GetID() uint64
	GetSessionPlaceID() uint64
	GetType() string
	GetDetails() string
}

type SessionPlacesOnline struct {
	*SessionPlace

	ID             uint64     `json:"id"`
	Platform       string     `json:"platform"`
	Link           string     `json:"Link"`
	SessionPlaceID uint64     `json:"session_place_id"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

func (place SessionPlacesOnline) GetID() uint64 {
	return place.ID
}

func (place SessionPlacesOnline) GetSessionPlaceID() uint64 {
	return place.SessionPlaceID
}

func (place SessionPlacesOnline) GetType() string {
	return common.SESSION_PLACE_ONLINE
}

func (place SessionPlacesOnline) GetDetails() string {
	return fmt.Sprintf("%s (%s)", place.Platform, place.Link)
}

type SessionPlacesGivenVenue struct {
	*SessionPlace

	ID             uint64     `json:"id"`
	Name           string     `json:"name"`
	Location       string     `json:"location"`
	SessionPlaceID uint64     `json:"session_place_id"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

func (place SessionPlacesGivenVenue) GetID() uint64 {
	return place.ID
}

func (place SessionPlacesGivenVenue) GetSessionPlaceID() uint64 {
	return place.SessionPlaceID
}

func (place SessionPlacesGivenVenue) GetType() string {
	return common.SESSION_PLACE_GIVEN_VENUE
}

func (place SessionPlacesGivenVenue) GetDetails() string {
	return fmt.Sprintf("%s (%s)", place.Name, place.Location)
}

type SessionPlacesMemberHome struct {
	*SessionPlace

	ID             uint64     `json:"id"`
	SessionPlaceID uint64     `json:"session_place_id"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

func (place SessionPlacesMemberHome) GetID() uint64 {
	return place.ID
}

func (place SessionPlacesMemberHome) GetSessionPlaceID() uint64 {
	return place.SessionPlaceID
}

func (place SessionPlacesMemberHome) GetType() string {
	return common.SESSION_PLACE_MEMBER_HOME
}

func (place SessionPlacesMemberHome) GetDetails() string {
	return ""
}
