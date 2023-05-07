CREATE TABLE users (
  id Integer Primary Key Generated Always as Identity,
  phone VARCHAR(15) UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  token TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- INSERT INTO users
-- SELECT phone, email, password, token, created_at, updated_at
-- FROM members;
