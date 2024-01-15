CREATE TYPE SessionPlacesType AS ENUM('online', 'given_venue', 'member_home');

CREATE TABLE IF NOT EXISTS session_places (
  id INTEGER GENERATED ALWAYS AS IDENTITY,
  type SessionPlacesType,
  session_id INTEGER NOT NULL,

  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  CONSTRAINT fk_session_places_sessions_session_id
    FOREIGN KEY (session_id) REFERENCES sessions(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
