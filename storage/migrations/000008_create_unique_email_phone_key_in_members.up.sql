ALTER TABLE members DROP CONSTRAINT members_email_key;
ALTER TABLE members DROP CONSTRAINT members_phone_key;
ALTER TABLE members ADD CONSTRAINT un_members_email_phone_key UNIQUE (email, phone)