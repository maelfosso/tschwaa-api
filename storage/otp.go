package storage

import (
	"context"

	"tschwaa.com/api/models"
)

const createOTP = `-- name: CreateOTP :one
INSERT INTO otps (
  wa_message_id, phone, pin_code
) VALUES (
  $1, $2, $3
)
RETURNING id, wa_message_id, phone, pin_code, active
`

type CreateOTPParams struct {
	WaMessageID string `db:"wa_message_id" json:"wa_message_id"`
	Phone       string `db:"phone" json:"phone"`
	PinCode     string `db:"pin_code" json:"pin_code"`
}

func (q *Queries) CreateOTP(ctx context.Context, arg CreateOTPParams) (*models.Otp, error) {
	row := q.db.QueryRowContext(ctx, createOTP, arg.WaMessageID, arg.Phone, arg.PinCode)
	var i models.Otp
	err := row.Scan(
		&i.ID,
		&i.WaMessageID,
		&i.Phone,
		&i.PinCode,
		&i.Active,
	)
	return &i, err
}

const deactivateOTP = `-- name: DeactivateOTP :exec
UPDATE otps
SET active = FALSE
WHERE id = $1
`

func (q *Queries) DeactivateOTP(ctx context.Context, id uint64) error {
	_, err := q.db.ExecContext(ctx, deactivateOTP, id)
	return err
}

const getActiveOTPFromPhone = `-- name: GetActiveOTPFromPhone :one
SELECT id, wa_message_id, phone, pin_code, active
FROM otps
WHERE phone = $1 AND active = TRUE
`

func (q *Queries) GetActiveOTPFromPhone(ctx context.Context, phone string) (*models.Otp, error) {
	row := q.db.QueryRowContext(ctx, getActiveOTPFromPhone, phone)
	var i models.Otp
	err := row.Scan(
		&i.ID,
		&i.WaMessageID,
		&i.Phone,
		&i.PinCode,
		&i.Active,
	)
	return &i, err
}

const checkOTP = `-- name: CheckOTP :one
SELECT id, wa_message_id, phone, pin_code, active
FROM otps
WHERE phone = $1 AND pin_code = $2 AND active = TRUE
`

type CheckOTPParams struct {
	Phone   string `db:"phone" json:"phone"`
	PinCode string `db:"pin_code" json:"pin_code"`
}

func (q *Queries) CheckOTP(ctx context.Context, arg CheckOTPParams) (*models.Otp, error) {
	row := q.db.QueryRowContext(ctx, checkOTP, arg.Phone, arg.PinCode)
	var i models.Otp
	err := row.Scan(
		&i.ID,
		&i.WaMessageID,
		&i.Phone,
		&i.PinCode,
		&i.Active,
	)
	return &i, err
}
