ALTER TABLE memberships
  ADD COLUMN position TEXT NOT NULL DEFAULT 'Member',
  ADD COLUMN status TEXT NOT NULL DEFAULT 'Resident',
  ADD COLUMN role TEXT NOT NULL DEFAULT 'member'
;
