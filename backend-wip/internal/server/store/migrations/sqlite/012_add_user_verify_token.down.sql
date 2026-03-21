-- Remove verify_token column from users table
DROP INDEX IF EXISTS idx_users_verify_token;
-- SQLite doesn't support DROP COLUMN directly, would need table recreation
-- For simplicity, leave column in place on rollback
