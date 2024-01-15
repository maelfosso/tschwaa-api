CREATE TABLE IF NOT EXISTS session_places_member_home (
  id INTEGER GENERATED ALWAYS AS IDENTITY,
  session_place_id INTEGER NOT NULL,

  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  CONSTRAINT fk_session_place_member_home_session_places_session_place_id
    FOREIGN KEY (session_place_id) REFERENCES sessions(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
