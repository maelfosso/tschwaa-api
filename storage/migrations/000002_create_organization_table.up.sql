CREATE TABLE IF NOT EXISTS organizations (
  id INTEGER GENERATED ALWAYS AS IDENTITY,
  name TEXT NOT NULL,
  created_by INTEGER,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  CONSTRAINT fk_memberships_members_created_by
    FOREIGN KEY (created_by) REFERENCES members(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE
);
