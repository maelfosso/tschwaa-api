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

const approvedAdhesion = `-- name: ApprovedAdhesion :one
UPDATE adhesions
SET joined = TRUE, joined_at = NOW()
WHERE id = $1
RETURNING id, member_id, organization_id, created_at, updated_at, joined, joined_at, position, status, role
`

func (q *Queries) ApprovedAdhesion(ctx context.Context, id uint64) (*models.Adhesion, error) {
	row := q.db.QueryRowContext(ctx, approvedAdhesion, id)
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

const getAdhesion = `-- name: GetAdhesion :one
SELECT id, member_id, organization_id, created_at, updated_at, joined, joined_at, position, status, role
FROM adhesions
WHERE id = $1
`

func (q *Queries) GetAdhesion(ctx context.Context, id uint64) (*models.Adhesion, error) {
	row := q.db.QueryRowContext(ctx, getAdhesion, id)
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

const createInvitation = `-- name: CreateInvitation :one
INSERT INTO invitations(link, adhesion_id)
VALUES ($1, $2)
RETURNING id, link, active, adhesion_id, created_at, updated_at
`

type CreateInvitationParams struct {
	Link       string `db:"link" json:"link"`
	AdhesionID uint64 `db:"adhesion_id" json:"adhesion_id"`
}

func (q *Queries) CreateInvitation(ctx context.Context, arg CreateInvitationParams) (*models.Invitation, error) {
	row := q.db.QueryRowContext(ctx, createInvitation, arg.Link, arg.AdhesionID)
	var i models.Invitation
	err := row.Scan(
		&i.ID,
		&i.Link,
		&i.Active,
		&i.AdhesionID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getInvitation = `-- name: GetInvitation :one
SELECT link, active, i.created_at, i.updated_at,
  a.joined, a.member_id as member_id, a.organization_id as organization_id,
  m.id, m.first_name, m.last_name, m.sex, m.phone, m.email, m.user_id,
  o.id, o.name, o.description
FROM invitations i
INNER JOIN adhesions a ON i.adhesion_id = a.id
INNER JOIN members m ON a.member_id = m.id
INNER JOIN organizations o ON a.organization_id = o.id
WHERE link = $1
`

func (q *Queries) GetInvitation(ctx context.Context, link string) (*models.Invitation, error) {
	row := q.db.QueryRowContext(ctx, getInvitation, link)
	var i models.Invitation
	err := row.Scan(
		&i.Link,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Adhesion.Joined,
		&i.Adhesion.MemberID,
		&i.Adhesion.OrganizationID,
		&i.Member.ID,
		&i.Member.FirstName,
		&i.Member.LastName,
		&i.Member.Sex,
		&i.Member.Phone,
		&i.Member.Email,
		&i.Member.UserID,
		&i.Organization.ID,
		&i.Organization.Name,
		&i.Organization.Description,
	)
	return &i, err
}

const desactivateInvitation = `-- name: DesactivateInvitation :exec
UPDATE invitations
SET active = FALSE
WHERE adhesion_id = $1 AND active = TRUE
`

func (q *Queries) DesactivateInvitation(ctx context.Context, adhesionID uint64) error {
	_, err := q.db.ExecContext(ctx, desactivateInvitation, adhesionID)
	return err
}

const desactivateInvitationFromLink = `-- name: DesactivateInvitationFromLink :one
UPDATE invitations
SET active = FALSE
WHERE link = $1
RETURNING id, link, active, adhesion_id, created_at, updated_at
`

func (q *Queries) DesactivateInvitationFromLink(ctx context.Context, link string) (*models.Invitation, error) {
	row := q.db.QueryRowContext(ctx, desactivateInvitationFromLink, link)
	var i models.Invitation
	err := row.Scan(
		&i.ID,
		&i.Link,
		&i.Active,
		&i.AdhesionID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}
