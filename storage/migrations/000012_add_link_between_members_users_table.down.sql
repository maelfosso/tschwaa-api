ALTER TABLE users DROP CONSTRAINT fk_members_users_user_id;
ALTER TABLE users
  DROP COLUMN member_id
;

ALTER TABLE members DROP CONSTRAINT fk_users_members_member_id;
ALTER TABLE members
  DROP COLUMN user_id
;
