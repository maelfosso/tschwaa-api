package storage

import (
	"context"

	"tschwaa.com/api/common"
	"tschwaa.com/api/models"
	"tschwaa.com/api/utils"
)

func (store *SQLStorage) GetSessionPlaceTx(ctx context.Context, sessionID uint64) (*models.ISessionPlace, error) {
	var subSessionPlace models.ISessionPlace

	err := store.execTx(ctx, func(q *Queries) error {
		sessionPlace, err := q.GetSessionPlace(ctx, sessionID)
		if err != nil {
			return utils.Fail(
				"error when getting session place",
				"ERR_CRT_SES_01",
				err,
			)
		}

		if sessionPlace.Type == common.SESSION_PLACE_ONLINE {
			subSessionPlace, err = q.GetSessionPlaceOnline(ctx, sessionPlace.ID)

			if err != nil {
				return utils.Fail(
					"error when getting session place online",
					"ERR_CRT_SES_01",
					err,
				)
			}
		} else if sessionPlace.Type == common.SESSION_PLACE_GIVEN_VENUE {
			subSessionPlace, err = q.GetSessionPlaceGiveVenue(ctx, sessionPlace.ID)

			if err != nil {
				return utils.Fail(
					"error when getting session place given venu",
					"ERR_CRT_SES_01",
					err,
				)
			}
		} else if sessionPlace.Type == common.SESSION_PLACE_MEMBER_HOME {
			subSessionPlace, err = q.GetSessionPlaceMemberHome(ctx, sessionPlace.ID)

			if err != nil {
				return utils.Fail(
					"error when getting session place member home",
					"ERR_CRT_SES_01",
					err,
				)
			}
		}

		return nil
	})

	return &subSessionPlace, err
}
