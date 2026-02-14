DROP INDEX IF EXISTS idx_blog_posts_visibility;
ALTER TABLE blog_posts DROP COLUMN IF EXISTS visibility;
