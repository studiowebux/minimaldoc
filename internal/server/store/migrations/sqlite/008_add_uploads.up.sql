-- Uploads table for storing file metadata
CREATE TABLE uploads (
    id TEXT PRIMARY KEY,
    site_id TEXT NOT NULL REFERENCES sites(id),
    user_id TEXT NOT NULL REFERENCES users(id),
    filename TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    size_bytes INTEGER NOT NULL,
    storage_path TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_uploads_site ON uploads(site_id);
CREATE INDEX idx_uploads_user ON uploads(user_id);
CREATE INDEX idx_uploads_created ON uploads(created_at DESC);
