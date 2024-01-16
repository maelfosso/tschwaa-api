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

-- name: GetSession :one
SELECT *
FROM sessions
WHERE organization_id = $1 AND id = $2;

-- name: CreateSessionPlace :one
INSERT INTO session_places(type, session_id)
VALUES ($1, $2)
RETURNING *;

-- name: CreateSessionPlaceOnline :one
INSERT INTO session_places_online(type, url, session_place_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateSessionPlaceGivenVenue :one
INSERT INTO session_places_given_venue(name, location, session_place_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateSessionPlaceMemberHome :one
INSERT INTO session_places_member_home(session_place_id)
VALUES ($1)
RETURNING *;

-- name: DeleteSessionPlaceOnline :exec
DELETE
FROM session_places_online
WHERE id = $1;

-- name: DeleteSessionPlaceGivenVenue :exec
DELETE
FROM session_places_given_venue
WHERE id = $1;

-- name: DeleteSessionPlaceMemberHome :exec
DELETE
FROM session_places_member_home
WHERE id = $1;

-- name: UpdateSessionPlace :one
UPDATE session_places
SET type = $2
WHERE id = $1
RETURNING *;

-- name: GetSessionPlace :one
SELECT *
FROM session_places
WHERE session_id = $1;

-- name: GetSessionPlaceOnline :one
SELECT *
FROM session_places_online
WHERE session_place_id = $1;

-- name: GetSessionPlaceGiveVenue :one
SELECT *
FROM session_places_given_venue
WHERE session_place_id = $1;

-- name: GetSessionPlaceMemberHome :one
SELECT *
FROM session_places_member_home
WHERE session_place_id = $1;
