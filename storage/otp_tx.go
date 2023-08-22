package storage

import (
	"context"
	"database/sql"
	"fmt"

	"tschwaa.com/api/models"
	"tschwaa.com/api/utils"
)

func (store *SQLStorage) CreateOTPTx(ctx context.Context, arg CreateOTPParams) (*models.Otp, error) {
	var res *models.Otp

	err := store.execTx(ctx, func(q *Queries) error {
		otp, err := q.GetActiveOTPFromPhone(ctx, arg.Phone)
		if err != nil && err != sql.ErrNoRows {
			return utils.Fail(
				fmt.Sprintf("error when getting active otp from %s", arg.Phone),
				"ERR_CRT_OTP_01", err)
		}

		if otp != nil {
			err = q.DeactivateOTP(ctx, otp.ID)
			if err != nil {
				return utils.Fail(
					fmt.Sprintf("error when desactivating otp %d", otp.ID),
					"ERR_CRT_OTP_02", err)
			}
		}

		otp, err = q.CreateOTP(ctx, arg)
		if err != nil {
			return utils.Fail(
				fmt.Sprintf("error when creating the otp %s for %s", arg.PinCode, arg.Phone),
				"ERR_CRT_OTP_03", err)
		}

		res = otp
		return err
	})

	return res, err
}
