-- Revoked JWT access tokens (for logout invalidation).
-- Entries are cleaned up after expires_at passes.
CREATE TABLE revoked_tokens (
    jti TEXT PRIMARY KEY,
    expires_at TEXT NOT NULL
);
CREATE INDEX idx_revoked_tokens_expires ON revoked_tokens(expires_at);
