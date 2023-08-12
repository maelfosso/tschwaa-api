package storage

import (
	"context"
	"time"

	"tschwaa.com/api/models"
)

const createAdhesion = `-- name: CreateAdhesion :one
INSERT INTO adhesions(member_id, organization_id, joined, joined_at)
VALUES ($1, $2, $3, $4)
RETURNING id, member_id, organization_id, created_at, updated_at, joined, joined_at, position, status, role
`

type CreateAdhesionParams struct {
	MemberID       uint64    `db:"member_id" json:"member_id"`
	OrganizationID uint64    `db:"organization_id" json:"organization_id"`
	Joined         bool      `db:"joined" json:"joined"`
	JoinedAt       time.Time `db:"joined_at" json:"joined_at"`
}

func (q *Queries) CreateAdhesion(ctx context.Context, arg CreateAdhesionParams) (*models.Adhesion, error) {
	row := q.db.QueryRowContext(ctx, createAdhesion,
		arg.MemberID,
		arg.OrganizationID,
		arg.Joined,
		arg.JoinedAt,
	)
	var i models.Adhesion
	err := row.Scan(
		&i.ID,
		&i.MemberID,
		&i.OrganizationID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Joined,
		&i.JoinedAt,
		&i.Position,
		&i.Status,
		&i.Role,
	)
	return &i, err
}

const getMembersFromOrganization = `-- name: GetMembersFromOrganization :many
SELECT m.id, m.first_name, m.last_name, m.sex, m.phone, a.position, a.role, a.status, a.joined, a.joined_at
FROM adhesions a INNER JOIN members m on a.member_id = m.id
WHERE a.organization_id = $1
`

func (q *Queries) GetMembersFromOrganization(ctx context.Context, organizationID uint64) ([]*models.OrganizationMember, error) {
	rows, err := q.db.QueryContext(ctx, getMembersFromOrganization, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.OrganizationMember{}
	for rows.Next() {
		var i models.OrganizationMember
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Sex,
			&i.Phone,
			&i.Position,
			&i.Role,
			&i.Status,
			&i.Joined,
			&i.JoinedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
