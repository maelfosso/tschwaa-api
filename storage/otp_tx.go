package storage

import (
	"context"

	"tschwaa.com/api/models"
)

func (store *SQLStorage) CreateOTPTx(ctx context.Context, arg CreateOTPParams) (*models.Otp, error) {
	var res *models.Otp

	err := store.execTx(ctx, func(q *Queries) error {
		otp, err := q.GetActiveOTPFromPhone(ctx, arg.Phone)
		if err != nil {
			return err
		}

		if otp != nil {
			err = q.DeactivateOTP(ctx, otp.ID)
			if err != nil {
				return err
			}
		}

		otp, err = q.CreateOTP(ctx, arg)
		if err != nil {
			return err
		}

		res = otp
		return err
	})

	return res, err
}
