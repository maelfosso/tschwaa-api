package storage

import (
	"context"
	"database/sql"

	"tschwaa.com/api/models"
)

const addMemberInSession = `-- name: AddMemberInSession :one
INSERT INTO members_of_session(membership_id, session_id)
VALUES ($1, $2)
RETURNING id, membership_id, session_id, created_at, updated_at
`

type AddMemberInSessionParams struct {
	MembershipID int32 `db:"membership_id" json:"membership_id"`
	SessionID    int32 `db:"session_id" json:"session_id"`
}

func (q *Queries) AddMemberInSession(ctx context.Context, arg AddMemberInSessionParams) (*models.MembersOfSession, error) {
	row := q.db.QueryRowContext(ctx, addMemberInSession, arg.MembershipID, arg.SessionID)
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
	SessionID      int32         `db:"session_id" json:"session_id"`
	OrganizationID sql.NullInt32 `db:"organization_id" json:"organization_id"`
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
AND mos.session_id = $1 AND m.organization_id = $2 AND m.member_id = $3
`

type RemoveMemberFromSessionParams struct {
	SessionID      int32         `db:"session_id" json:"session_id"`
	OrganizationID sql.NullInt32 `db:"organization_id" json:"organization_id"`
	MemberID       sql.NullInt32 `db:"member_id" json:"member_id"`
}

func (q *Queries) RemoveMemberFromSession(ctx context.Context, arg RemoveMemberFromSessionParams) error {
	_, err := q.db.ExecContext(ctx, removeMemberFromSession, arg.SessionID, arg.OrganizationID, arg.MemberID)
	return err
}
