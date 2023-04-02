ALTER TABLE members
  DROP COLUMN last_name;
ALTER TABLE members
  RENAME COLUMN first_name TO name;
