DROP INDEX IF EXISTS idx_blog_posts_scheduled;
ALTER TABLE blog_posts DROP COLUMN IF EXISTS scheduled_at;
