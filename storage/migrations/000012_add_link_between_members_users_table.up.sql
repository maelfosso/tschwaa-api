ALTER TABLE members
  ADD COLUMN user_id INTEGER DEFAULT NULL,
  ADD CONSTRAINT fk_members_users_user_id
      FOREIGN KEY (user_id) REFERENCES users(id)
      ON DELETE SET NULL
      ON UPDATE CASCADE
;

ALTER TABLE users
  ADD COLUMN member_id INTEGER NOT NULL,
  ADD CONSTRAINT fk_users_members_member_id
      FOREIGN KEY (member_id) REFERENCES members(id)
      ON DELETE CASCADE
      ON UPDATE CASCADE
;
