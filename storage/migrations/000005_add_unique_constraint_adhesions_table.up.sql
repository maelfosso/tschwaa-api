ALTER TABLE memberships ADD CONSTRAINT ak_memberships_member_id_organization_id
  UNIQUE (member_id, organization_id);
