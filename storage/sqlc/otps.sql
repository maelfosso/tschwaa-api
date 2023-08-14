-- name: GetActiveOTPFromPhone :one
SELECT *
FROM otps
WHERE phone = $1 AND active = TRUE;

-- name: DeactivateOTP :exec
UPDATE otps
SET active = FALSE
WHERE id = $1;

-- name: CreateOTP :one
INSERT INTO otps (
  wa_message_id, phone, pin_code
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: CheckOTP :one
SELECT *
FROM otps
WHERE phone = $1 AND pin_code = $2 AND active = TRUE;
