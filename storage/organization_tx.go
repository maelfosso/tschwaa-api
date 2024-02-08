package storage

import (
	"context"
	"fmt"
	"time"

	"tschwaa.com/api/models"
	"tschwaa.com/api/utils"
)

func (store *SQLStorage) CreateOrganizationWithMembershipTx(ctx context.Context, arg CreateOrganizationParams) (*models.Organization, error) {
	var result *models.Organization

	err := store.execTx(ctx, func(q *Queries) error {
		org, err := q.CreateOrganization(ctx, arg)
		if err != nil {
			return utils.Fail(
				fmt.Sprintf("error when creating organizatioin %s", arg.Name),
				"ERR_CRT_ORG_MBRSHP_01", err)
		}
		result = org

		_, err = q.CreateMembership(ctx, CreateMembershipParams{
			MemberID:       *arg.CreatedBy,
			OrganizationID: org.ID,
			Joined:         true,
			JoinedAt:       time.Now(),
		})
		return utils.Fail(
			fmt.Sprintf("error when creating membership of member[%d] into organization[%d]", *arg.CreatedBy, org.ID),
			"ERR_CRT_ORG_MBRSHP_02", err)
	})

	return result, err
}
