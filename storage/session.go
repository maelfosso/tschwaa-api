package storage

import (
	"context"

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
