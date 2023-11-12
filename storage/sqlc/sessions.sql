-- name: GetCurrentSession :one
SELECT *
FROM sessions
WHERE in_progress = TRUE;

-- name: CreateSession :one
INSERT INTO sessions(start_date, end_date, in_progress, organization_id)
VALUES ($1, $2, $3, $4)
RETURNING *;
