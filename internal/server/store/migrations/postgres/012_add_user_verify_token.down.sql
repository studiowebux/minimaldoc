-- Remove verify_token column from users table
DROP INDEX IF EXISTS idx_users_verify_token;
ALTER TABLE users DROP COLUMN IF EXISTS verify_token;
