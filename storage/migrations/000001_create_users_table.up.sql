CREATE TABLE IF NOT EXISTS users(
  id Integer Primary Key Generated Always as Identity,
  firstname VARCHAR(300),
  lastname VARCHAR(300),
  phone VARCHAR(15) UNIQUE NOT NULL,
  email VARCHAR (50) UNIQUE NOT NULL,
  password VARCHAR (50) NOT NULL,
  token VARCHAR(100) NOT NULL,
  createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
  updatedAt TIMESTAMP NOT NULL DEFAULT NOW()
);
