package storage

import (
	"context"
	"fmt"
	"time"

	"tschwaa.com/api/models"
	"tschwaa.com/api/utils"
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
			return utils.Fail(
				fmt.Sprintf("error when creating membership of member[%d] into organization[%d]", arg.MemberID, arg.OrganizationID),
				"ERR_CRT_INV_01",
				err,
			)
		}

		err = q.DesactivateInvitation(ctx, membership.ID)
		if err != nil {
			return utils.Fail(
				fmt.Sprintf("error when desactivating invitation from membership[%d]", membership.ID),
				"ERR_CRT_INV_02",
				err,
			)
		}

		_, err = q.CreateInvitation(ctx, CreateInvitationParams{
			Link:         arg.JoinId,
			MembershipID: membership.ID,
		})
		return utils.Fail(
			fmt.Sprintf("error when creating invitation %s of %d", arg.JoinId, membership.ID),
			"ERR_CRT_INV_03",
			err,
		)
	})

	return &org, err
}

func (store *SQLStorage) ApprovedInvitationTx(ctx context.Context, link string) error {
	err := store.execTx(ctx, func(q *Queries) error {

		invitation, err := q.DesactivateInvitationFromLink(ctx, link)
		if err != nil {
			return utils.Fail(
				fmt.Sprintf("error when desactivating invitation from link %s", link),
				"ERR_APR_ORG_INV_01", err)
		}

		_, err = q.ApprovedMembership(ctx, invitation.MembershipID)
		return utils.Fail(
			fmt.Sprintf("error when approving membership %d", invitation.MembershipID),
			"ERR_APR_ORG_INV_02", err)
	})

	return err
}
