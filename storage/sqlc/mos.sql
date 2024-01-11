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

-- name: ListAllMembersOfSession :many
SELECT mos.id, mos.session_id, mos.created_at, mos.updated_at,
  m.id as member_id, m.first_name, m.last_name, m.sex, m.phone,
  a.id as membership_id, a.position, a.role, a.status, a.joined, a.joined_at
FROM members m
INNER JOIN memberships a ON m.id = a.member_id
LEFT JOIN members_of_session mos ON a.id = mos.membership_id AND a.organization_id = $1 AND mos.session_id = $2;
