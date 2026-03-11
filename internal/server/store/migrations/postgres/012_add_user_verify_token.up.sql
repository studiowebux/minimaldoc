-- Add verify_token column to users table for email verification
ALTER TABLE users ADD COLUMN verify_token TEXT;

-- Index for token lookup
CREATE INDEX idx_users_verify_token ON users(verify_token);
