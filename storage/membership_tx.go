package storage

import (
	"context"
	"fmt"
	"time"

	"tschwaa.com/api/models"
)

type CreateMembershipInvitationParams struct {
	MemberID       uint64 `db:"member_id" json:"member_id"`
	OrganizationID uint64 `db:"organization_id" json:"organization_id"`
	Joined         bool   `db:"joined" json:"joined"`
	JoinId         string
}

func (store *SQLStorage) CreateInvitationTx(ctx context.Context, arg CreateMembershipInvitationParams) (*models.Organization, error) {
	var org models.Organization

	err := store.execTx(ctx, func(q *Queries) error {
		membership, err := q.CreateMembership(ctx, CreateMembershipParams{
			MemberID:       arg.MemberID,
			OrganizationID: arg.OrganizationID,
			Joined:         false,
			JoinedAt:       time.Now(),
		})
		if err != nil {
			return fmt.Errorf("error when creating membership of member[%d] into organization[%d]: %w", arg.MemberID, arg.OrganizationID, err)
		}

		err = q.DesactivateInvitation(ctx, membership.ID)
		if err != nil {
			return fmt.Errorf("error when desactivating invitation from membership[%d]: %w", membership.ID, err)
		}

		_, err = q.CreateInvitation(ctx, CreateInvitationParams{
			Link:         arg.JoinId,
			MembershipID: membership.ID,
		})
		return fmt.Errorf("error when creating invitation %s of %d: %w", arg.JoinId, membership.ID, err)
	})

	return &org, err
}

func (store *SQLStorage) ApprovedInvitationTx(ctx context.Context, link string) error {
	err := store.execTx(ctx, func(q *Queries) error {

		invitation, err := q.DesactivateInvitationFromLink(ctx, link)
		if err != nil {
			return fmt.Errorf("error when desactivating invitation from link %s: %w", link, err)
		}

		_, err = q.ApprovedMembership(ctx, invitation.MembershipID)
		return fmt.Errorf("error when approving membership %d: %w", invitation.MembershipID, err)
	})

	return err
}
