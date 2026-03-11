-- Revert to original role constraint (requires no 'author' users exist)
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('admin', 'editor', 'viewer'));
