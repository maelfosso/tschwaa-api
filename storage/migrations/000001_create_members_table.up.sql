CREATE TABLE IF NOT EXISTS members (
  id Integer Primary Key Generated Always as Identity,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  sex TEXT NOT NULL DEFAULT 'male',
  phone VARCHAR(15) UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  token TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
