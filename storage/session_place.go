package storage

import (
	"context"

	"tschwaa.com/api/models"
)

const createSessionPlace = `-- name: CreateSessionPlace :one
INSERT INTO session_places(type, session_id)
VALUES ($1, $2)
RETURNING id, type, session_id, created_at, updated_at
`

type CreateSessionPlaceParams struct {
	Type      string `db:"type" json:"type"`
	SessionID uint64 `db:"session_id" json:"session_id"`
}

func (q *Queries) CreateSessionPlace(ctx context.Context, arg CreateSessionPlaceParams) (*models.SessionPlace, error) {
	row := q.db.QueryRowContext(ctx, createSessionPlace, arg.Type, arg.SessionID)
	var i models.SessionPlace
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.SessionID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const createSessionPlaceGivenVenue = `-- name: CreateSessionPlaceGivenVenue :one
INSERT INTO session_places_given_venue(name, location, session_place_id)
VALUES ($1, $2, $3)
RETURNING id, name, location, session_place_id, created_at, updated_at
`

type CreateSessionPlaceGivenVenueParams struct {
	Name           string      `db:"name" json:"name"`
	Location       interface{} `db:"location" json:"location"`
	SessionPlaceID uint64      `db:"session_place_id" json:"session_place_id"`
}

func (q *Queries) CreateSessionPlaceGivenVenue(ctx context.Context, arg CreateSessionPlaceGivenVenueParams) (*models.SessionPlacesGivenVenue, error) {
	row := q.db.QueryRowContext(ctx, createSessionPlaceGivenVenue, arg.Name, arg.Location, arg.SessionPlaceID)
	var i models.SessionPlacesGivenVenue
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Location,
		&i.SessionPlaceID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const createSessionPlaceMemberHome = `-- name: CreateSessionPlaceMemberHome :one
INSERT INTO session_places_member_home(session_place_id)
VALUES ($1)
RETURNING id, session_place_id, created_at, updated_at
`

func (q *Queries) CreateSessionPlaceMemberHome(ctx context.Context, sessionPlaceID uint64) (*models.SessionPlacesMemberHome, error) {
	row := q.db.QueryRowContext(ctx, createSessionPlaceMemberHome, sessionPlaceID)
	var i models.SessionPlacesMemberHome
	err := row.Scan(
		&i.ID,
		&i.SessionPlaceID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const createSessionPlaceOnline = `-- name: CreateSessionPlaceOnline :one
INSERT INTO session_places_online(type, url, session_place_id)
VALUES ($1, $2, $3)
RETURNING id, type, url, session_place_id, created_at, updated_at
`

type CreateSessionPlaceOnlineParams struct {
	Type           string `db:"type" json:"type"`
	Url            string `db:"url" json:"url"`
	SessionPlaceID uint64 `db:"session_place_id" json:"session_place_id"`
}

func (q *Queries) CreateSessionPlaceOnline(ctx context.Context, arg CreateSessionPlaceOnlineParams) (*models.SessionPlacesOnline, error) {
	row := q.db.QueryRowContext(ctx, createSessionPlaceOnline, arg.Type, arg.Url, arg.SessionPlaceID)
	var i models.SessionPlacesOnline
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Url,
		&i.SessionPlaceID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const deleteSessionPlaceGivenVenue = `-- name: DeleteSessionPlaceGivenVenue :exec
DELETE
FROM session_places_given_venue
WHERE id = $1
`

func (q *Queries) DeleteSessionPlaceGivenVenue(ctx context.Context, id uint64) error {
	_, err := q.db.ExecContext(ctx, deleteSessionPlaceGivenVenue, id)
	return err
}

const deleteSessionPlaceMemberHome = `-- name: DeleteSessionPlaceMemberHome :exec
DELETE
FROM session_places_member_home
WHERE id = $1
`

func (q *Queries) DeleteSessionPlaceMemberHome(ctx context.Context, id uint64) error {
	_, err := q.db.ExecContext(ctx, deleteSessionPlaceMemberHome, id)
	return err
}

const deleteSessionPlaceOnline = `-- name: DeleteSessionPlaceOnline :exec
DELETE
FROM session_places_online
WHERE id = $1
`

func (q *Queries) DeleteSessionPlaceOnline(ctx context.Context, id uint64) error {
	_, err := q.db.ExecContext(ctx, deleteSessionPlaceOnline, id)
	return err
}

const updateSessionPlace = `-- name: UpdateSessionPlace :one
UPDATE session_places
SET type = $2
WHERE id = $1
RETURNING id, type, session_id, created_at, updated_at
`

type UpdateSessionPlaceParams struct {
	ID   uint64 `db:"id" json:"id"`
	Type string `db:"type" json:"type"`
}

func (q *Queries) UpdateSessionPlace(ctx context.Context, arg UpdateSessionPlaceParams) (*models.SessionPlace, error) {
	row := q.db.QueryRowContext(ctx, updateSessionPlace, arg.ID, arg.Type)
	var i models.SessionPlace
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.SessionID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getSessionPlaceFromSession = `-- name: GetSessionPlaceFromSession :one
SELECT id, type, session_id, created_at, updated_at
FROM session_places
WHERE session_id = $1
`

func (q *Queries) GetSessionPlaceFromSession(ctx context.Context, sessionID uint64) (*models.SessionPlace, error) {
	row := q.db.QueryRowContext(ctx, getSessionPlaceFromSession, sessionID)
	var i models.SessionPlace
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.SessionID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getSessionPlace = `-- name: GetSessionPlace :one
SELECT id, type, session_id, created_at, updated_at
FROM session_places
WHERE id = $1 and session_id = $2
`

type GetSessionPlaceParams struct {
	ID        uint64 `json:"id"`
	SessionID uint64 `json:"session_id"`
}

func (q *Queries) GetSessionPlace(ctx context.Context, arg GetSessionPlaceParams) (*models.SessionPlace, error) {
	row := q.db.QueryRowContext(ctx, getSessionPlace, arg.ID, arg.SessionID)
	var i models.SessionPlace
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.SessionID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getSessionPlaceGiveVenue = `-- name: GetSessionPlaceGiveVenue :one
SELECT id, name, location, session_place_id, created_at, updated_at
FROM session_places_given_venue
WHERE session_place_id = $1
`

func (q *Queries) GetSessionPlaceGiveVenue(ctx context.Context, sessionPlaceID uint64) (*models.SessionPlacesGivenVenue, error) {
	row := q.db.QueryRowContext(ctx, getSessionPlaceGiveVenue, sessionPlaceID)
	var i models.SessionPlacesGivenVenue
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Location,
		&i.SessionPlaceID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getSessionPlaceMemberHome = `-- name: GetSessionPlaceMemberHome :one
SELECT id, session_place_id, created_at, updated_at
FROM session_places_member_home
WHERE session_place_id = $1
`

func (q *Queries) GetSessionPlaceMemberHome(ctx context.Context, sessionPlaceID uint64) (*models.SessionPlacesMemberHome, error) {
	row := q.db.QueryRowContext(ctx, getSessionPlaceMemberHome, sessionPlaceID)
	var i models.SessionPlacesMemberHome
	err := row.Scan(
		&i.ID,
		&i.SessionPlaceID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getSessionPlaceOnline = `-- name: GetSessionPlaceOnline :one
SELECT id, type, url, session_place_id, created_at, updated_at
FROM session_places_online
WHERE session_place_id = $1
`

func (q *Queries) GetSessionPlaceOnline(ctx context.Context, sessionPlaceID uint64) (*models.SessionPlacesOnline, error) {
	row := q.db.QueryRowContext(ctx, getSessionPlaceOnline, sessionPlaceID)
	var i models.SessionPlacesOnline
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Url,
		&i.SessionPlaceID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}
