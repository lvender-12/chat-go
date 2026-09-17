ALTER TABLE users
ADD COLUMN password_salt CHAR(32) NOT NULL
AFTER password_hash;
