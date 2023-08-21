CREATE TABLE invitations (
  id INTEGER GENERATED ALWAYS AS IDENTITY,
  link TEXT NOT NULL,
  active BOOLEAN DEFAULT TRUE,
  membership_id INTEGER,

  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  CONSTRAINT fk_invitations_memberships_membership_id
    FOREIGN KEY (membership_id) REFERENCES memberships(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
