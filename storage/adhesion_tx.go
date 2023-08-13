package storage

import (
	"context"
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
			return err
		}

		err = q.DesactivateInvitation(ctx, adhesion.ID)
		if err != nil {
			return err
		}

		_, err = q.CreateInvitation(ctx, CreateInvitationParams{
			Link:       arg.JoinId,
			AdhesionID: adhesion.ID,
		})
		return err
	})

	return &org, err
}

func (store *SQLStorage) ApprovedInvitationTx(ctx context.Context, link string) error {
	err := store.execTx(ctx, func(q *Queries) error {

		invitation, err := q.DesactivateInvitationFromLink(ctx, link)
		if err != nil {
			return err
		}

		_, err = q.ApprovedAdhesion(ctx, invitation.AdhesionID)
		return err
	})

	return err
}
