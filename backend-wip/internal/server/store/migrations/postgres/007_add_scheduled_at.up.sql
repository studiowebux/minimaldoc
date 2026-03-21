-- Add scheduled publishing support
ALTER TABLE blog_posts ADD COLUMN scheduled_at TIMESTAMPTZ;

CREATE INDEX idx_blog_posts_scheduled ON blog_posts(scheduled_at) WHERE scheduled_at IS NOT NULL AND status = 'draft';
