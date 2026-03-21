-- MinimalDoc Server - Rollback Initial Schema (PostgreSQL)

DROP INDEX IF EXISTS idx_users_oauth;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_subscribers_site;
DROP INDEX IF EXISTS idx_sessions_expires;
DROP INDEX IF EXISTS idx_sessions_user;
DROP INDEX IF EXISTS idx_ratings_site_path;
DROP INDEX IF EXISTS idx_events_name;
DROP INDEX IF EXISTS idx_events_site_created;
DROP INDEX IF EXISTS idx_page_views_path;
DROP INDEX IF EXISTS idx_page_views_site_created;

DROP TABLE IF EXISTS subscribers;
DROP TABLE IF EXISTS ratings;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS page_views;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS sites;
