-- Add duration tracking to page_views
ALTER TABLE page_views ADD COLUMN duration_seconds INTEGER DEFAULT NULL;
