package storage

import (
	"context"

	"tschwaa.com/api/common"
	"tschwaa.com/api/models"
	"tschwaa.com/api/utils"
)

func (store *SQLStorage) GetSessionPlaceTx(ctx context.Context, sessionID uint64) (models.ISessionPlace, error) {
	var subSessionPlace models.ISessionPlace

	err := store.execTx(ctx, func(q *Queries) error {
		sessionPlace, err := q.GetSessionPlaceFromSession(ctx, sessionID)
		if err != nil {
			return utils.Fail(
				"error when getting session place",
				"ERR_CRT_SES_01",
				err,
			)
		}

		if sessionPlace.Type == common.SESSION_PLACE_ONLINE {
			subSessionPlace, err = q.GetSessionPlaceOnlineFromSessionPlace(ctx, sessionPlace.ID)

			if err != nil {
				return utils.Fail(
					"error when getting session place online",
					"ERR_CRT_SES_01",
					err,
				)
			}
			// TODO: check if subSessionPlace exists
		} else if sessionPlace.Type == common.SESSION_PLACE_GIVEN_VENUE {
			subSessionPlace, err = q.GetSessionPlaceGivenVenueFromSessionPlace(ctx, sessionPlace.ID)

			if err != nil {
				return utils.Fail(
					"error when getting session place given venu",
					"ERR_CRT_SES_01",
					err,
				)
			}
			// TODO: check if subSessionPlace exists
		} else if sessionPlace.Type == common.SESSION_PLACE_MEMBER_HOME {
			subSessionPlace, err = q.GetSessionPlaceMemberHomeFromSessionPlace(ctx, sessionPlace.ID)

			if err != nil {
				return utils.Fail(
					"error when getting session place member home",
					"ERR_CRT_SES_01",
					err,
				)
			}
			// TODO: check if subSessionPlace exists
		}

		return nil
	})

	return subSessionPlace, err
}

type DeleteSessionPlaceTxParams struct {
	ISessionPlace models.ISessionPlace
	SessionID     uint64
}

func (store *SQLStorage) DeleteSessionPlaceTx(ctx context.Context, arg DeleteSessionPlaceTxParams) error {
	err := store.execTx(ctx, func(q *Queries) error {
		if arg.ISessionPlace.GetType() == common.SESSION_PLACE_ONLINE {
			err := store.DeleteSessionPlaceOnline(ctx, arg.ISessionPlace.GetID())
			if err != nil {
				return utils.Fail(
					"error when deleting online session place",
					"ERR_CRT_SES_01",
					err,
				)
			}
		} else if arg.ISessionPlace.GetType() == common.SESSION_PLACE_GIVEN_VENUE {
			err := store.DeleteSessionPlaceGivenVenue(ctx, arg.ISessionPlace.GetID())
			if err != nil {
				return utils.Fail(
					"error when deleting given venue session place",
					"ERR_CRT_SES_01",
					err,
				)
			}
		} else if arg.ISessionPlace.GetType() == common.SESSION_PLACE_MEMBER_HOME {
			err := store.DeleteSessionPlaceMemberHome(ctx, arg.ISessionPlace.GetID())
			if err != nil {
				return utils.Fail(
					"error when deleting member home session place",
					"ERR_CRT_SES_01",
					err,
				)
			}
		}

		err := store.DeleteSessionPlace(ctx, DeleteSessionPlaceParams{
			ID:        arg.ISessionPlace.GetSessionPlaceID(),
			SessionID: arg.SessionID,
		})
		if err != nil {
			return utils.Fail(
				"error when deleting session place",
				"ERR_CRT_SES_01",
				err,
			)
		}

		return nil
	})

	return err
}

type CreateSessionPlaceTxParams struct {
	SessionID        uint64
	SessionPlaceType string

	OnlineType *string
	Url        *string

	Name     *string
	Location *string
}

func (store *SQLStorage) CreateSessionPlaceTx(ctx context.Context, arg CreateSessionPlaceTxParams) (models.ISessionPlace, error) {
	var iSessionPlace models.ISessionPlace

	err := store.execTx(ctx, func(q *Queries) error {
		sessionPlace, err := store.CreateSessionPlace(ctx, CreateSessionPlaceParams{
			Type:      arg.SessionPlaceType,
			SessionID: arg.SessionID,
		})
		if err != nil {
			return utils.Fail(
				"error when creating a session place",
				"ERR_CRT_SES_01",
				err,
			)
		}

		if sessionPlace.Type == common.SESSION_PLACE_ONLINE {
			iSessionPlace, err = store.CreateSessionPlaceOnline(ctx, CreateSessionPlaceOnlineParams{
				SessionPlaceID: sessionPlace.ID,
				Type:           *arg.OnlineType,
				Url:            *arg.Url,
			})
			if err != nil {
				return utils.Fail(
					"error when creating online session place",
					"ERR_CRT_SES_01",
					err,
				)
			}
		} else if sessionPlace.Type == common.SESSION_PLACE_GIVEN_VENUE {
			iSessionPlace, err = store.CreateSessionPlaceGivenVenue(ctx, CreateSessionPlaceGivenVenueParams{
				SessionPlaceID: sessionPlace.ID,
				Name:           *arg.Name,
				Location:       *arg.Location,
			})
			if err != nil {
				return utils.Fail(
					"error when creating online session place",
					"ERR_CRT_SES_01",
					err,
				)
			}
		} else if sessionPlace.Type == common.SESSION_PLACE_MEMBER_HOME {
			iSessionPlace, err = store.CreateSessionPlaceMemberHome(ctx, sessionPlace.ID)
			if err != nil {
				return utils.Fail(
					"error when creating online session place",
					"ERR_CRT_SES_01",
					err,
				)
			}
		}

		return nil
	})

	return iSessionPlace, err
}
