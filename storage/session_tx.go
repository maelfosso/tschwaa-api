package storage

import (
	"context"

	"tschwaa.com/api/models"
	"tschwaa.com/api/utils"
)

func (store *SQLStorage) CreateSessionTx(ctx context.Context, arg CreateSessionParams) (*models.Session, error) {
	var session *models.Session

	err := store.execTx(ctx, func(q *Queries) error {
		err := q.NoSessionInProgress(ctx, arg.OrganizationID)
		if err != nil {
			return utils.Fail(
				"error when setting no progress",
				"ERR_CRT_SES_01",
				err,
			)
		}

		session, err = q.CreateSession(ctx, CreateSessionParams{
			StartDate:      arg.StartDate,
			EndDate:        arg.EndDate,
			InProgress:     true,
			OrganizationID: arg.OrganizationID,
		})
		return utils.Fail(
			"error when creating session",
			"ERR_CRT_SES_02",
			err,
		)
	})

	return session, err
}
