package storage

import (
	"context"
	"fmt"
	"time"

	"tschwaa.com/api/models"
)

func (store *SQLStorage) CreateOrganizationWithAdhesionTx(ctx context.Context, arg CreateOrganizationParams) (*models.Organization, error) {
	var org models.Organization

	err := store.execTx(ctx, func(q *Queries) error {
		org, err := q.CreateOrganization(ctx, arg)
		if err != nil {
			return fmt.Errorf("error when creating organizatioin %s: %s", arg.Name, err)
		}

		_, err = q.CreateAdhesion(ctx, CreateAdhesionParams{
			MemberID:       *arg.CreatedBy,
			OrganizationID: org.ID,
			Joined:         true,
			JoinedAt:       time.Now(),
		})
		return fmt.Errorf("error when creating adhesion of member[%d] into organization[%d]: %w", *arg.CreatedBy, org.ID, err)
	})

	return &org, err
}
