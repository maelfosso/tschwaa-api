CREATE TABLE invitations (
  id INTEGER GENERATED ALWAYS AS IDENTITY,
  link TEXT NOT NULL,
  active BOOLEAN DEFAULT TRUE,
  adhesion_id INTEGER,

  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  CONSTRAINT fk_invitations_adhesions_adhesion_id
    FOREIGN KEY (adhesion_id) REFERENCES adhesions(id)
);
