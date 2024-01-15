CREATE TYPE SessionPlacesOnlineType AS ENUM('telegram', 'whatsapp', 'google_meet', 'zoom');

CREATE TABLE IF NOT EXISTS session_places_online (
  id INTEGER GENERATED ALWAYS AS IDENTITY,
  type SessionPlacesOnlineType,
  url TEXT NOT NULL,
  session_place_id INTEGER NOT NULL,

  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  CONSTRAINT fk_session_place_online_session_places_session_place_id
    FOREIGN KEY (session_place_id) REFERENCES sessions(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
