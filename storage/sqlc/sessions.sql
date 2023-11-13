-- name: GetCurrentSession :one
SELECT *
FROM sessions
WHERE organization_id = $1 AND in_progress = TRUE;

-- name: NoSessionInProgress :exec
UPDATE sessions
SET in_progress = FALSE
WHERE organization_id = $1 AND in_progress = TRUE;

-- name: CreateSession :one
INSERT INTO sessions(start_date, end_date, in_progress, organization_id)
VALUES ($1, $2, $3, $4)
RETURNING *;
