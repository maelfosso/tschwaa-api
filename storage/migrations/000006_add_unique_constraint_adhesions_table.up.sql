ALTER TABLE adhesions ADD CONSTRAINT ak_adhesions_member_id_organization_id
  UNIQUE (member_id, organization_id);
