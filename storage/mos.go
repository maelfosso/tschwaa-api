package storage

import (
	"context"

	"tschwaa.com/api/models"
)

const listAllMembersOfSession = `-- name: ListAllMembersOfSession :many
SELECT mos.id, mos.session_id, mos.created_at, mos.updated_at,
  m.id as member_id, m.first_name, m.last_name, m.sex, m.phone,
  a.id as membership_id, a.position, a.role, a.status, a.joined, a.joined_at
FROM members m
INNER JOIN memberships a ON m.id = a.member_id
LEFT JOIN members_of_session mos ON a.id = mos.membership_id AND a.organization_id = $1 AND mos.session_id = $2
`

type ListAllMembersOfSessionParams struct {
	OrganizationID uint64 `json:"organization_id"`
	SessionID      uint64 `json:"session_id"`
}

func (q *Queries) ListAllMembersOfSession(ctx context.Context, arg ListAllMembersOfSessionParams) ([]*models.MembersOfSession, error) {
	rows, err := q.db.QueryContext(ctx, listAllMembersOfSession, arg.OrganizationID, arg.SessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.MembersOfSession{}
	for rows.Next() {
		var i models.MembersOfSession
		if err := rows.Scan(
			&i.ID,
			&i.SessionID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.MemberID,
			&i.FirstName,
			&i.LastName,
			&i.Sex,
			&i.Phone,
			&i.MembershipID,
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

const addMemberToSession = `-- name: AddMemberToSession :one
INSERT INTO members_of_session(membership_id, session_id)
VALUES ($1, $2)
RETURNING id, membership_id, session_id, created_at, updated_at
`

type AddMemberToSessionParams struct {
	MembershipID uint64 `db:"membership_id" json:"membership_id"`
	SessionID    uint64 `db:"session_id" json:"session_id"`
}

func (q *Queries) AddMemberToSession(ctx context.Context, arg AddMemberToSessionParams) (*models.MembersOfSession, error) {
	row := q.db.QueryRowContext(ctx, addMemberToSession, arg.MembershipID, arg.SessionID)
	var i models.MembersOfSession
	err := row.Scan(
		&i.ID,
		&i.MembershipID,
		&i.SessionID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const removeAllMembersFromSession = `-- name: RemoveAllMembersFromSession :exec
DELETE
FROM members_of_session mos
USING memberships m
WHERE mos.membership_id = m.id
AND mos.session_id = $1 AND m.organization_id = $2
`

type RemoveAllMembersFromSessionParams struct {
	SessionID      uint64 `db:"session_id" json:"session_id"`
	OrganizationID uint64 `db:"organization_id" json:"organization_id"`
}

func (q *Queries) RemoveAllMembersFromSession(ctx context.Context, arg RemoveAllMembersFromSessionParams) error {
	_, err := q.db.ExecContext(ctx, removeAllMembersFromSession, arg.SessionID, arg.OrganizationID)
	return err
}

const removeMemberFromSession = `-- name: RemoveMemberFromSession :exec
DELETE
FROM members_of_session mos
USING memberships m
WHERE mos.membership_id = m.id
AND mos.id = $1 AND m.organization_id = $2 AND mos.session_id = $3
`

type RemoveMemberFromSessionParams struct {
	ID             uint64 `json:"id"`
	OrganizationID uint64 `json:"organization_id"`
	SessionID      uint64 `json:"session_id"`
}

func (q *Queries) RemoveMemberFromSession(ctx context.Context, arg RemoveMemberFromSessionParams) error {
	_, err := q.db.ExecContext(ctx, removeMemberFromSession, arg.ID, arg.OrganizationID, arg.SessionID)
	return err
}
