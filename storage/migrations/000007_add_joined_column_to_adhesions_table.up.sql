ALTER TABLE memberships
  ADD COLUMN joined BOOLEAN DEFAULT FALSE,
  ADD COLUMN joined_at TIMESTAMP
;