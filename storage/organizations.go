package storage

import (
	"context"

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
