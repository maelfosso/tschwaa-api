package storage

import (
	"context"
	"database/sql"
	"fmt"

	"tschwaa.com/api/models"
)

func (store *SQLStorage) CreateOTPTx(ctx context.Context, arg CreateOTPParams) (*models.Otp, error) {
	var res *models.Otp

	err := store.execTx(ctx, func(q *Queries) error {
		otp, err := q.GetActiveOTPFromPhone(ctx, arg.Phone)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error when getting active otp from %s: %w", arg.Phone, err)
		}

		if otp != nil {
			err = q.DeactivateOTP(ctx, otp.ID)
			if err != nil {
				return fmt.Errorf("error when desactivating otp %d: %w", otp.ID, err)
			}
		}

		otp, err = q.CreateOTP(ctx, arg)
		if err != nil {
			return fmt.Errorf("error when creating the otp %s for %s: %w", arg.PinCode, arg.Phone, err)
		}

		res = otp
		return err
	})

	return res, err
}
