-- name: GetCurrentSession :one
SELECT *
FROM sessions
WHERE current = TRUE;

-- name: CreateSession :one
INSERT INTO sessions(start_date, end_date, current, organization_id)
VALUES ($1, $2, $3, $4)
RETURNING *;
