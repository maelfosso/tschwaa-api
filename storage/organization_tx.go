package storage

import (
	"context"
	"time"

	"tschwaa.com/api/models"
)

func (store *SQLStorage) CreateOrganizationWithAdhesionTx(ctx context.Context, arg CreateOrganizationParams) (*models.Organization, error) {
	var org models.Organization

	err := store.execTx(ctx, func(q *Queries) error {
		org, err := q.CreateOrganization(ctx, arg)
		if err != nil {
			return err
		}

		_, err = q.CreateAdhesion(ctx, CreateAdhesionParams{
			MemberID:       *arg.CreatedBy,
			OrganizationID: org.ID,
			Joined:         true,
			JoinedAt:       time.Now(),
		})
		return err
	})

	return &org, err
}
