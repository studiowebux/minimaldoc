-- Remove OAuth identity columns from subscribers
DROP INDEX IF EXISTS idx_subscribers_oauth;
ALTER TABLE subscribers DROP COLUMN IF EXISTS verified_via;
ALTER TABLE subscribers DROP COLUMN IF EXISTS oauth_display_name;
ALTER TABLE subscribers DROP COLUMN IF EXISTS oauth_provider;
