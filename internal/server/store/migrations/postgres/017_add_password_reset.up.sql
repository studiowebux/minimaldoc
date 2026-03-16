-- Add password reset token columns to users table.
ALTER TABLE users ADD COLUMN reset_token TEXT;
ALTER TABLE users ADD COLUMN reset_token_expires_at TIMESTAMPTZ;
CREATE INDEX idx_users_reset_token ON users(reset_token);
