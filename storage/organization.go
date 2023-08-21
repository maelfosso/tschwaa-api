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

// func (d *Database) CreateOrganization(ctx context.Context, org models.models.models.Organization) (uint64, error) {
// 	query := `
// 		INSERT INTO organizations(name, description, created_by)
// 		VALUES ($1, $2, $3)
// 		RETURNING id
// 	`

// 	var lastInsertId uint64 = 0
// 	err := d.DB.QueryRowContext(ctx, query, org.Name, org.Description, org.CreatedBy).Scan(&lastInsertId)
// 	return lastInsertId, err
// }

// func (d *Database) ListAllOrganizationFromMember(ctx context.Context, id uint64) ([]models.models.models.Organization, error) {
// 	query := `
// 		SELECT id, name, description
// 		FROM organizations
// 		WHERE created_by = $1
// 	`
// 	rows, err := d.DB.QueryContext(ctx, query, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var orgs = []models.models.models.Organization{}
// 	for rows.Next() {
// 		var id, name, description string
// 		if err := rows.Scan(&id, &name, &description); err != nil {
// 			return nil, fmt.Errorf("error when parsing the organizations result")
// 		}

// 		i, _ := strconv.ParseUint(id, 10, 64)
// 		orgs = append(orgs, models.models.models.Organization{
// 			ID:          i,
// 			Name:        name,
// 			Description: description,
// 		})
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error when parsing the organizations result")
// 	}
// 	return orgs, nil
// }

// func (d *Database) GetOrganization(ctx context.Context, orgId uint64) (*models.models.models.Organization, error) {
// 	var org models.models.models.Organization

// 	query := `
// 		SELECT id, name, description, created_at
// 		FROM organizations
// 		WHERE id = $1
// 	`
// 	if err := d.DB.QueryRowContext(ctx, query, orgId).Scan(&org.ID, &org.Name, &org.Description, &org.CreatedAt); err == nil {
// 		return &org, nil
// 		// } else if err == sql.ErrNoRows {
// 		// 	return nil, fmt.Errorf("organization does not exist")
// 	} else {
// 		return nil, err
// 	}
// }
