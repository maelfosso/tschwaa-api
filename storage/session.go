package storage

import (
	"context"
	"time"

	"tschwaa.com/api/models"
)

const getCurrentSession = `-- name: GetCurrentSession :one
SELECT id, start_date, end_date, organization_id, current, created_at, updated_at
FROM sessions
WHERE current = TRUE
`

func (q *Queries) GetCurrentSession(ctx context.Context) (*models.Session, error) {
	row := q.db.QueryRowContext(ctx, getCurrentSession)
	var i models.Session
	err := row.Scan(
		&i.ID,
		&i.StartDate,
		&i.EndDate,
		&i.Current,
		&i.OrganizationID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const createSession = `-- name: CreateSession :one
INSERT INTO sessions(start_date, end_date, current, organization_id)
VALUES ($1, $2, $3, $4)
RETURNING id, start_date, end_date, organization_id, current, created_at, updated_at
`

type CreateSessionParams struct {
	StartDate      time.Time `db:"start_date" json:"start_date"`
	EndDate        time.Time `db:"end_date" json:"end_date"`
	Current        bool      `db:"current" json:"current"`
	OrganizationID uint64    `db:"organization_id" json:"organization_id"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (*models.Session, error) {
	row := q.db.QueryRowContext(ctx, createSession,
		arg.StartDate,
		arg.EndDate,
		arg.Current,
		arg.OrganizationID,
	)
	var i models.Session
	err := row.Scan(
		&i.ID,
		&i.StartDate,
		&i.EndDate,
		&i.OrganizationID,
		&i.Current,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}
