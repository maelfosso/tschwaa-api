CREATE TABLE IF NOT EXISTS session_places_given_venue (
  id INTEGER GENERATED ALWAYS AS IDENTITY,
  name TEXT NOT NULL,
  location POINT NOT NULL,
  session_place_id INTEGER NOT NULL,

  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  CONSTRAINT fk_session_place_given_venue_session_places_session_place_id
    FOREIGN KEY (session_place_id) REFERENCES sessions(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
