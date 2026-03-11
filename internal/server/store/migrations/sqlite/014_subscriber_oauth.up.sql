-- Add OAuth identity columns to subscribers
ALTER TABLE subscribers ADD COLUMN oauth_provider TEXT DEFAULT '';
ALTER TABLE subscribers ADD COLUMN oauth_display_name TEXT DEFAULT '';
ALTER TABLE subscribers ADD COLUMN verified_via TEXT DEFAULT '';

-- Index for looking up subscribers by OAuth provider
CREATE INDEX idx_subscribers_oauth ON subscribers(site_id, oauth_provider) WHERE oauth_provider != '';
