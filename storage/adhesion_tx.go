package storage

import (
	"context"
	"fmt"
	"time"

	"tschwaa.com/api/models"
)

type CreateAdhesionInvitationParams struct {
	MemberID       uint64 `db:"member_id" json:"member_id"`
	OrganizationID uint64 `db:"organization_id" json:"organization_id"`
	Joined         bool   `db:"joined" json:"joined"`
	JoinId         string
}

func (store *SQLStorage) CreateInvitationTx(ctx context.Context, arg CreateAdhesionInvitationParams) (*models.Organization, error) {
	var org models.Organization

	err := store.execTx(ctx, func(q *Queries) error {
		adhesion, err := q.CreateAdhesion(ctx, CreateAdhesionParams{
			MemberID:       arg.MemberID,
			OrganizationID: arg.OrganizationID,
			Joined:         false,
			JoinedAt:       time.Now(),
		})
		if err != nil {
			return fmt.Errorf("error when creating adhesion of member[%d] into organization[%d]: %w", arg.MemberID, arg.OrganizationID, err)
		}

		err = q.DesactivateInvitation(ctx, adhesion.ID)
		if err != nil {
			return fmt.Errorf("error when desactivating invitation from adhesion[%d]: %w", adhesion.ID, err)
		}

		_, err = q.CreateInvitation(ctx, CreateInvitationParams{
			Link:       arg.JoinId,
			AdhesionID: adhesion.ID,
		})
		return fmt.Errorf("error when creating invitation %s of %d: %w", arg.JoinId, adhesion.ID, err)
	})

	return &org, err
}

func (store *SQLStorage) ApprovedInvitationTx(ctx context.Context, link string) error {
	err := store.execTx(ctx, func(q *Queries) error {

		invitation, err := q.DesactivateInvitationFromLink(ctx, link)
		if err != nil {
			return fmt.Errorf("error when desactivating invitation from link %s: %w", link, err)
		}

		_, err = q.ApprovedAdhesion(ctx, invitation.AdhesionID)
		return fmt.Errorf("error when approving adhesion %d: %w", invitation.AdhesionID, err)
	})

	return err
}
