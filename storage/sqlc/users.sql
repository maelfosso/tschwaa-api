-- name: GetMemberByPhone :one
SELECT *
FROM members
WHERE phone = $1;

-- name: 