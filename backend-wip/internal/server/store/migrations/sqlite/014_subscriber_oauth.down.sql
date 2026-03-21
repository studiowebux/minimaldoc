-- SQLite does not support DROP COLUMN before 3.35.0
-- These columns are additive and safe to leave; drop the index only
DROP INDEX IF EXISTS idx_subscribers_oauth;
