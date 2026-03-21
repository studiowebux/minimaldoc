-- Add visibility column to blog_posts
ALTER TABLE blog_posts ADD COLUMN visibility VARCHAR(20) DEFAULT 'public';

-- Valid values: public, authenticated, role_viewer, role_author, role_editor, role_admin
CREATE INDEX idx_blog_posts_visibility ON blog_posts(site_id, visibility);
