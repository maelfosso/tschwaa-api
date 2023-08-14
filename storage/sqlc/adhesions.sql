-- name: CreateAdhesion :one
INSERT INTO adhesions(member_id, organization_id, joined, joined_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetMembersFromOrganization :many
SELECT m.id, m.first_name, m.last_name, m.sex, m.phone, a.position, a.role, a.status, a.joined, a.joined_at
FROM adhesions a INNER JOIN members m on a.member_id = m.id
WHERE a.organization_id = $1;

-- name: GetAdhesion :one
SELECT *
FROM adhesions
WHERE id = $1;

-- name: ApprovedAdhesion :one
UPDATE adhesions
SET joined = TRUE, joined_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreateInvitation :one
INSERT INTO invitations(link, adhesion_id)
VALUES ($1, $2)
RETURNING *;

-- name: DesactivateInvitation :exec
UPDATE invitations
SET active = FALSE
WHERE adhesion_id = $1 AND active = TRUE;

-- name: DesactivateInvitationFromLink :one
UPDATE invitations
SET active = FALSE
WHERE link = $1
RETURNING *;

-- name: GetInvitation :one
SELECT link, active, i.created_at,
  a.joined, a.member_id as member_id, a.organization_id as organization_id,
  m.id, m.first_name, m.last_name, m.sex, m.phone, m.email, m.user_id,
  o.id, o.name, o.description
FROM invitations i
INNER JOIN adhesions a ON i.adhesion_id = a.id
INNER JOIN members m ON a.member_id = m.id
INNER JOIN organizations o ON a.organization_id = o.id
WHERE link = $1;