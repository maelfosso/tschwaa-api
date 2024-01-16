package storage

import (
	"context"
	"database/sql"
	"time"

	"tschwaa.com/api/models"
)

const getCurrentSession = `-- name: GetCurrentSession :one
SELECT id, start_date, end_date, organization_id, in_progress, created_at, updated_at
FROM sessions
WHERE organization_id = $1 AND in_progress = TRUE
`

func (q *Queries) GetCurrentSession(ctx context.Context, organizationID uint64) (*models.Session, error) {
	row := q.db.QueryRowContext(ctx, getCurrentSession, organizationID)
	var i models.Session
	err := row.Scan(
		&i.ID,
		&i.StartDate,
		&i.EndDate,
		&i.OrganizationID,
		&i.InProgress,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return &i, err
}

const getSession = `-- name: GetSession :one
SELECT id, start_date, end_date, organization_id, in_progress, created_at, updated_at
FROM sessions
WHERE organization_id = $1 AND id = $2
`

type GetSessionParams struct {
	OrganizationID uint64 `db:"organization_id" json:"organization_id"`
	SessionID      uint64 `db:"id" json:"session_id"`
}

func (q *Queries) GetSession(ctx context.Context, arg GetSessionParams) (*models.Session, error) {
	row := q.db.QueryRowContext(ctx, getSession, arg.OrganizationID, arg.SessionID)
	var i models.Session
	err := row.Scan(
		&i.ID,
		&i.StartDate,
		&i.EndDate,
		&i.OrganizationID,
		&i.InProgress,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const noSessionInProgress = `-- name: NoSessionInProgress :exec
UPDATE sessions
SET in_progress = FALSE
WHERE organization_id = $1 AND in_progress = TRUE
`

func (q *Queries) NoSessionInProgress(ctx context.Context, organizationID uint64) error {
	_, err := q.db.ExecContext(ctx, noSessionInProgress, organizationID)
	return err
}

const createSession = `-- name: CreateSession :one
INSERT INTO sessions(start_date, end_date, in_progress, organization_id)
VALUES ($1, $2, $3, $4)
RETURNING id, start_date, end_date, organization_id, in_progress, created_at, updated_at
`

type CreateSessionParams struct {
	StartDate      time.Time `db:"start_date" json:"start_date"`
	EndDate        time.Time `db:"end_date" json:"end_date"`
	InProgress     bool      `db:"in_progress" json:"in_progress"`
	OrganizationID uint64    `db:"organization_id" json:"organization_id"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (*models.Session, error) {
	row := q.db.QueryRowContext(ctx, createSession,
		arg.StartDate,
		arg.EndDate,
		arg.InProgress,
		arg.OrganizationID,
	)
	var i models.Session
	err := row.Scan(
		&i.ID,
		&i.StartDate,
		&i.EndDate,
		&i.OrganizationID,
		&i.InProgress,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}
