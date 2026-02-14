-- Add bounce tracking to page_views
ALTER TABLE page_views ADD COLUMN is_bounce BOOLEAN DEFAULT NULL;
