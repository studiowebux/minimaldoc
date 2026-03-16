-- Hash newsletter verification tokens and add expiry.
-- Existing plaintext tokens are cleared; unverified subscribers must re-subscribe.
ALTER TABLE subscribers ADD COLUMN verify_token_expires_at TIMESTAMPTZ;
UPDATE subscribers SET verify_token = NULL WHERE verify_token IS NOT NULL;
