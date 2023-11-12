CREATE TABLE IF NOT EXISTS sessions (
  id INTEGER GENERATED ALWAYS AS IDENTITY,
  start_date TIMESTAMP NOT NULL,
  end_date TIMESTAMP NOT NULL,
  organization_id INTEGER NOT NULL,
  in_progress BOOLEAN DEFAULT FALSE,
  
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  CONSTRAINT fk_sessions_organizations_organization_id
    FOREIGN KEY (organization_id) REFERENCES organizations(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
