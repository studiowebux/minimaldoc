-- Remove duration tracking from page_views
ALTER TABLE page_views DROP COLUMN duration_seconds;
