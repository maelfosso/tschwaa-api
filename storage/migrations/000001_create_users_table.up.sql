CREATE TABLE IF NOT EXISTS users(
  id Integer Primary Key Generated Always as Identity,
  firstname TEXT NOT NULL,
  lastname TEXT NOT NULL,
  phone VARCHAR(15) UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  token TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
