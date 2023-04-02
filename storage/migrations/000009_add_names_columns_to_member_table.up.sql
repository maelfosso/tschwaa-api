ALTER TABLE members
  RENAME COLUMN name TO first_name;
ALTER TABLE members
  ADD COLUMN last_name TEXT DEFAULT '';
