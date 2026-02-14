-- Drop indexes
DROP INDEX IF EXISTS idx_blog_comments_created;
DROP INDEX IF EXISTS idx_blog_comments_parent;
DROP INDEX IF EXISTS idx_blog_comments_site_status;
DROP INDEX IF EXISTS idx_blog_comments_post;
DROP INDEX IF EXISTS idx_blog_posts_category;
DROP INDEX IF EXISTS idx_blog_posts_published;
DROP INDEX IF EXISTS idx_blog_posts_author;
DROP INDEX IF EXISTS idx_blog_posts_site_slug;
DROP INDEX IF EXISTS idx_blog_posts_site_status;

-- Drop tables
DROP TABLE IF EXISTS blog_comments;
DROP TABLE IF EXISTS blog_posts;
