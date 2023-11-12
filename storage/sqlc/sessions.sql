-- name: GetCurrentSession :one
SELECT *
FROM sessions
WHERE current = TRUE;
