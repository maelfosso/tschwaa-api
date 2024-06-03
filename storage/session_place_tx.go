package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"tschwaa.com/api/common"
	"tschwaa.com/api/models"
	"tschwaa.com/api/utils"
)

func (store *SQLStorage) GetSessionPlaceTx(ctx context.Context, sessionID uint64) (models.ISessionPlace, error) {
	var subSessionPlace models.ISessionPlace

	err := store.execTx(ctx, func(q *Queries) error {
		sessionPlace, err := q.GetSessionPlaceFromSession(ctx, sessionID)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			} else {
				return utils.Fail(
					"error when getting session place",
					"ERR_CRT_SES_01",
					err,
				)
			}
		}

		if sessionPlace == nil {
			return utils.Fail("no session place", "ERR_CRT_SES_02", nil)
		}

		if sessionPlace.PlaceType == common.SESSION_PLACE_ONLINE {
			_subSessionPlace, err := q.GetSessionPlaceOnlineFromSessionPlace(ctx, sessionPlace.ID)

			if err != nil {
				return utils.Fail(
					"error when getting session place online",
					"ERR_CRT_SES_01",
					err,
				)
			}

			// TODO: check if subSessionPlace exists
			if _subSessionPlace == nil {
				subSessionPlace = models.SessionPlacesOnline{
					SessionPlace: sessionPlace,

					SessionPlaceID: sessionPlace.ID,
				}
			} else {
				subSessionPlace = models.SessionPlacesOnline{
					SessionPlaceID: sessionPlace.ID,

					ID:        _subSessionPlace.ID,
					Platform:  _subSessionPlace.Platform,
					Link:      _subSessionPlace.Link,
					CreatedAt: _subSessionPlace.CreatedAt,
					UpdatedAt: _subSessionPlace.UpdatedAt,
				}
			}
		} else if sessionPlace.PlaceType == common.SESSION_PLACE_GIVEN_VENUE {
			_subSessionPlace, err := q.GetSessionPlaceGivenVenueFromSessionPlace(ctx, sessionPlace.ID)

			if err != nil {
				return utils.Fail(
					"error when getting session place given venu",
					"ERR_CRT_SES_01",
					err,
				)
			}
			// TODO: check if subSessionPlace exists
			if _subSessionPlace == nil {
				subSessionPlace = models.SessionPlacesGivenVenue{
					SessionPlace: sessionPlace,

					SessionPlaceID: sessionPlace.ID,
				}
			} else {
				subSessionPlace = models.SessionPlacesGivenVenue{
					SessionPlace: sessionPlace,

					SessionPlaceID: sessionPlace.ID,

					ID:        _subSessionPlace.ID,
					Name:      _subSessionPlace.Name,
					Location:  _subSessionPlace.Location,
					CreatedAt: _subSessionPlace.CreatedAt,
					UpdatedAt: _subSessionPlace.UpdatedAt,
				}
			}
		} else if sessionPlace.PlaceType == common.SESSION_PLACE_MEMBER_HOME {
			_subSessionPlace, err := q.GetSessionPlaceMemberHomeFromSessionPlace(ctx, sessionPlace.ID)

			if err != nil {
				return utils.Fail(
					"error when getting session place member home",
					"ERR_CRT_SES_01",
					err,
				)
			}
			// TODO: check if subSessionPlace exists
			if _subSessionPlace == nil {
				subSessionPlace = models.SessionPlacesMemberHome{
					SessionPlace: sessionPlace,

					SessionPlaceID: sessionPlace.ID,
				}
			} else {
				subSessionPlace = models.SessionPlacesMemberHome{
					SessionPlace: sessionPlace,

					SessionPlaceID: sessionPlace.ID,

					ID:        _subSessionPlace.ID,
					CreatedAt: _subSessionPlace.CreatedAt,
					UpdatedAt: _subSessionPlace.UpdatedAt,
				}
			}
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

	Platform *string
	Link     *string

	Name     *string
	Location *string

	Choice *string
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

		if sessionPlace.PlaceType == common.SESSION_PLACE_ONLINE {
			iSessionPlace, err = store.CreateSessionPlaceOnline(ctx, CreateSessionPlaceOnlineParams{
				SessionPlaceID: sessionPlace.ID,
				Type:           *arg.Platform,
				Link:           *arg.Link,
			})
			if err != nil {
				return utils.Fail(
					"error when creating online session place",
					"ERR_CRT_SES_01",
					err,
				)
			}
		} else if sessionPlace.PlaceType == common.SESSION_PLACE_GIVEN_VENUE {
			iSessionPlace, err = store.CreateSessionPlaceGivenVenue(ctx, CreateSessionPlaceGivenVenueParams{
				SessionPlaceID: sessionPlace.ID,
				Name:           *arg.Name,
				Location:       *arg.Location,
			})
			if err != nil {
				return utils.Fail(
					"error when creating given venue session place",
					"ERR_CRT_SES_01",
					err,
				)
			}
		} else if sessionPlace.PlaceType == common.SESSION_PLACE_MEMBER_HOME {
			iSessionPlace, err = store.CreateSessionPlaceMemberHome(ctx, sessionPlace.ID)
			if err != nil {
				return utils.Fail(
					"error when creating member home session place",
					"ERR_CRT_SES_01",
					err,
				)
			}
		}

		return nil
	})

	return iSessionPlace, err
}

type ChangeSessionPlaceParams struct {
	SessionID        uint64
	SessionPlaceType string

	Platform *string
	Link     *string

	Name     *string
	Location *string
}

func (store *SQLStorage) ChangeSessionPlaceTx(ctx context.Context, arg ChangeSessionPlaceParams) (models.ISessionPlace, error) {
	var iSessionPlace models.ISessionPlace

	err := store.execTx(ctx, func(q *Queries) error {
		iSessionPlace, err := store.GetSessionPlaceTx(ctx, arg.SessionID)
		log.Println("GetSessionPlaceTx")
		log.Println("iSessionPlace : ", iSessionPlace)
		log.Println("err: ", err)
		log.Println(iSessionPlace != nil, utils.CheckNilInterface(iSessionPlace))
		if err != nil {
			return utils.Fail(
				fmt.Sprintf("error when getting full session place of session[%d]: %w", arg.SessionID, err),
				"ERR_CRT_SES_01",
				err,
			)
		}
		if !utils.CheckNilInterface(iSessionPlace) {
			log.Println("DeleteSessionPlaceTx")
			err := store.DeleteSessionPlaceTx(ctx, DeleteSessionPlaceTxParams{
				ISessionPlace: iSessionPlace,
				SessionID:     arg.SessionID,
			})
			if err != nil {
				return utils.Fail(
					fmt.Sprintf("error when completely deleting a session place of session[%d]: %w", arg.SessionID, err),
					"ERR_CRT_SES_01",
					err,
				)
			}
		}

		log.Println("CreateSessionPlaceTx")
		iSessionPlace, err = store.CreateSessionPlaceTx(ctx, CreateSessionPlaceTxParams{
			SessionID:        arg.SessionID,
			SessionPlaceType: arg.SessionPlaceType,

			Platform: arg.Platform,
			Link:     arg.Link,

			Name:     arg.Name,
			Location: arg.Location,
		})
		if err != nil {
			return utils.Fail(
				fmt.Sprintf("error when completely creating a session place of session[%d]: %w", arg.SessionID, err),
				"ERR_CRT_SES_01",
				err,
			)
		}

		return nil
	})

	return iSessionPlace, err
}
