CREATE TABLE IF NOT EXISTS members_of_session (
  id INTEGER GENERATED ALWAYS AS IDENTITY,
  membership_id INTEGER NOT NULL,
  session_id INTEGER NOT NULL,

  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  CONSTRAINT fk_mos_sessions_session_id
    FOREIGN KEY (session_id) REFERENCES sessions(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT fk_mos_memberships_membership_id
    FOREIGN KEY (membership_id) REFERENCES memberships(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
