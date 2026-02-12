-- MinimalDoc Server - Initial Schema (SQLite)

-- Sites registered with the backend
CREATE TABLE sites (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    domain TEXT UNIQUE,
    api_key TEXT UNIQUE NOT NULL,
    config TEXT DEFAULT '{}',
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

-- Users (admin, editors, viewers)
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    site_id TEXT REFERENCES sites(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    password_hash TEXT,
    role TEXT DEFAULT 'viewer' CHECK (role IN ('admin', 'editor', 'viewer')),
    oauth_provider TEXT,
    oauth_id TEXT,
    name TEXT,
    avatar_url TEXT,
    email_verified INTEGER DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    last_login_at TEXT,
    UNIQUE(site_id, email)
);

-- Sessions for stateful auth
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT UNIQUE NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    expires_at TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now'))
);

-- API keys for programmatic access
CREATE TABLE api_keys (
    id TEXT PRIMARY KEY,
    site_id TEXT REFERENCES sites(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    key_hash TEXT UNIQUE NOT NULL,
    scopes TEXT DEFAULT '[]',
    expires_at TEXT,
    last_used_at TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

-- Cookie-free analytics (no PII)
CREATE TABLE page_views (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    site_id TEXT REFERENCES sites(id) ON DELETE CASCADE,
    path TEXT NOT NULL,
    referrer TEXT,
    country TEXT,
    device_type TEXT CHECK (device_type IN ('desktop', 'mobile', 'tablet')),
    browser TEXT,
    os TEXT,
    session_hash TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

-- Custom events tracking
CREATE TABLE events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    site_id TEXT REFERENCES sites(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    path TEXT,
    properties TEXT DEFAULT '{}',
    session_hash TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

-- Page ratings and feedback
CREATE TABLE ratings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    site_id TEXT REFERENCES sites(id) ON DELETE CASCADE,
    path TEXT NOT NULL,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    feedback TEXT,
    session_hash TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

-- Newsletter subscribers
CREATE TABLE subscribers (
    id TEXT PRIMARY KEY,
    site_id TEXT REFERENCES sites(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    verified INTEGER DEFAULT 0,
    verify_token TEXT,
    verify_sent_at TEXT,
    subscribed_at TEXT DEFAULT (datetime('now')),
    unsubscribed_at TEXT,
    UNIQUE(site_id, email)
);

-- Indexes for common queries
CREATE INDEX idx_page_views_site_created ON page_views(site_id, created_at DESC);
CREATE INDEX idx_page_views_path ON page_views(site_id, path, created_at);
CREATE INDEX idx_events_site_created ON events(site_id, created_at DESC);
CREATE INDEX idx_events_name ON events(site_id, name, created_at);
CREATE INDEX idx_ratings_site_path ON ratings(site_id, path);
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);
CREATE INDEX idx_subscribers_site ON subscribers(site_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_oauth ON users(oauth_provider, oauth_id);
