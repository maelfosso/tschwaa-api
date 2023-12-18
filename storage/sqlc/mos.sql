-- name: RemoveAllMembersFromSession :exec
DELETE
FROM members_of_session mos
USING memberships m
WHERE mos.membership_id = m.id
AND mos.session_id = $1 AND m.organization_id = $2;

-- name: RemoveMemberFromSession :exec
DELETE
FROM members_of_session mos
USING memberships m
WHERE mos.membership_id = m.id
AND mos.session_id = $1 AND m.organization_id = $2 AND m.member_id = $3;

-- name: AddMemberToSession :one
INSERT INTO members_of_session(membership_id, session_id)
VALUES ($1, $2)
RETURNING *;
