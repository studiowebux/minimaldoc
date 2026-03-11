-- Remove bounce tracking from page_views
ALTER TABLE page_views DROP COLUMN is_bounce;
