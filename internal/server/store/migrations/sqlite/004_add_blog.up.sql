-- Blog posts table
CREATE TABLE blog_posts (
    id TEXT PRIMARY KEY,
    site_id TEXT REFERENCES sites(id) ON DELETE CASCADE,
    author_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    slug TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    content TEXT NOT NULL,
    featured_image TEXT,
    tags TEXT DEFAULT '[]',
    category TEXT,
    status TEXT DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    published_at TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    UNIQUE(site_id, slug)
);

-- Blog comments table
CREATE TABLE blog_comments (
    id TEXT PRIMARY KEY,
    site_id TEXT REFERENCES sites(id) ON DELETE CASCADE,
    post_id TEXT REFERENCES blog_posts(id) ON DELETE CASCADE,
    parent_id TEXT REFERENCES blog_comments(id) ON DELETE CASCADE,
    author_name TEXT NOT NULL,
    author_email TEXT NOT NULL,
    content TEXT NOT NULL,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'spam')),
    ip_address TEXT,
    user_agent TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    moderated_at TEXT,
    moderated_by TEXT REFERENCES users(id) ON DELETE SET NULL
);

-- Indexes for blog posts
CREATE INDEX idx_blog_posts_site_status ON blog_posts(site_id, status);
CREATE INDEX idx_blog_posts_site_slug ON blog_posts(site_id, slug);
CREATE INDEX idx_blog_posts_author ON blog_posts(author_id);
CREATE INDEX idx_blog_posts_published ON blog_posts(site_id, status, published_at DESC);
CREATE INDEX idx_blog_posts_category ON blog_posts(site_id, category);

-- Indexes for blog comments
CREATE INDEX idx_blog_comments_post ON blog_comments(post_id);
CREATE INDEX idx_blog_comments_site_status ON blog_comments(site_id, status);
CREATE INDEX idx_blog_comments_parent ON blog_comments(parent_id);
CREATE INDEX idx_blog_comments_created ON blog_comments(created_at DESC);
