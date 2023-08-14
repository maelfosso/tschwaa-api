-- name: CreateOrganization :one
INSERT INTO organizations(name, description, created_by)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListOrganizations :many
SELECT *
FROM organizations;

-- name: ListOrganizationsCreatedBy :many
SELECT *
FROM organizations
WHERE created_by = $1;

-- name: ListOrganizationOfMember :many
SELECT O.*
FROM organizations O INNER JOIN adhesions A ON O.id = A.organization_id
WHERE A.member_id = $1;

-- name: GetOrganization :one
SELECT *
FROM organizations
WHERE id = $1;
