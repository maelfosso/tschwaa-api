ALTER TABLE members DROP CONSTRAINT un_members_email_phone_key;
ALTER TABLE members ADD CONSTRAINT members_email_key UNIQUE(email);
ALTER TABLE members ADD CONSTRAINT members_phone_key UNIQUE(phone);