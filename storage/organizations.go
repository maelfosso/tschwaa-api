package storage

import (
	"context"
	"fmt"
	"strconv"

	"tschwaa.com/api/model"
)

func (d *Database) CreateOrganization(ctx context.Context, org model.Organization) (int64, error) {
	query := `
		INSERT INTO organizations(name, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var lastInsertId int64 = 0
	err := d.DB.QueryRowContext(ctx, query, org.Name, org.Description, org.CreatedBy).Scan(&lastInsertId)
	return lastInsertId, err
}

func (d *Database) ListAllOrganizationFromUser(ctx context.Context, id uint64) ([]model.Organization, error) {
	query := `
		SELECT id, name, description
		FROM organizations
		WHERE created_by = $1
	`
	rows, err := d.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	var orgs = []model.Organization{}
	for rows.Next() {
		var id, name, description string
		if err := rows.Scan(&id, &name, &description); err != nil {
			return nil, fmt.Errorf("error when parsing the organizations result")
		}

		i, _ := strconv.ParseUint(id, 10, 64)
		orgs = append(orgs, model.Organization{
			ID:          i,
			Name:        name,
			Description: description,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error when parsing the organizations result")
	}
	return orgs, nil
}

func (d *Database) GetOrganization(ctx context.Context, orgId uint64) (*model.Organization, error) {
	var org model.Organization

	query := `
		SELECT id, name, description, created_at
		FROM organizations
		WHERE id = $1
	`
	if err := d.DB.QueryRowContext(ctx, query, orgId).Scan(&org.ID, &org.Name, &org.Description, &org.CreatedAt); err == nil {
		return &org, nil
		// } else if err == sql.ErrNoRows {
		// 	return nil, fmt.Errorf("organization does not exist")
	} else {
		return nil, err
	}
}
