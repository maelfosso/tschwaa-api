package storage

import (
	"context"
	"database/sql"

	"tschwaa.com/api/models"
)

const createOrganization = `-- name: CreateOrganization :one
INSERT INTO organizations(name, description, created_by)
VALUES ($1, $2, $3)
RETURNING id, name, created_by, created_at, updated_at, description
`

type CreateOrganizationParams struct {
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description"`
	CreatedBy   *uint64 `db:"created_by" json:"created_by"`
}

func (q *Queries) CreateOrganization(ctx context.Context, arg CreateOrganizationParams) (*models.Organization, error) {
	row := q.db.QueryRowContext(ctx, createOrganization, arg.Name, arg.Description, arg.CreatedBy)
	var i models.Organization
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
	)
	return &i, err
}

const getOrganization = `-- name: GetOrganization :one
SELECT id, name, created_by, created_at, updated_at, description
FROM organizations
WHERE id = $1
`

func (q *Queries) GetOrganization(ctx context.Context, id uint64) (*models.Organization, error) {
	row := q.db.QueryRowContext(ctx, getOrganization, id)
	var i models.Organization
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
	)
	return &i, err
}

const listOrganizationOfMember = `-- name: ListOrganizationOfMember :many
SELECT o.id, o.name, o.created_by, o.created_at, o.updated_at, o.description
FROM organizations O INNER JOIN memberships A ON O.id = A.organization_id
WHERE A.member_id = $1
`

func (q *Queries) ListOrganizationOfMember(ctx context.Context, memberID uint64) ([]*models.Organization, error) {
	rows, err := q.db.QueryContext(ctx, listOrganizationOfMember, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.Organization{}
	for rows.Next() {
		var i models.Organization
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedBy,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Description,
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

const listOrganizations = `-- name: ListOrganizations :many
SELECT id, name, created_by, created_at, updated_at, description
FROM organizations
`

func (q *Queries) ListOrganizations(ctx context.Context) ([]*models.Organization, error) {
	rows, err := q.db.QueryContext(ctx, listOrganizations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.Organization{}
	for rows.Next() {
		var i models.Organization
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedBy,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Description,
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

const listOrganizationsCreatedBy = `-- name: ListOrganizationsCreatedBy :many
SELECT id, name, created_by, created_at, updated_at, description
FROM organizations
WHERE created_by = $1
`

func (q *Queries) ListOrganizationsCreatedBy(ctx context.Context, createdBy sql.NullInt32) ([]*models.Organization, error) {
	rows, err := q.db.QueryContext(ctx, listOrganizationsCreatedBy, createdBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.Organization{}
	for rows.Next() {
		var i models.Organization
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedBy,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Description,
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
