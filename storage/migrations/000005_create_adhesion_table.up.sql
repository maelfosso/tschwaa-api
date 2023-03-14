CREATE TABLE adhesions (
  id INTEGER GENERATED ALWAYS AS IDENTITY,
  member_id INTEGER,
  organization_id INTEGER,

  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  CONSTRAINT fk_adhesions_members_member_id
    FOREIGN KEY (member_id) REFERENCES members(id),
  CONSTRAINT fk_adhesions_organizations_organization_id
    FOREIGN KEY (organization_id) REFERENCES organizations(id)
);
