-- Access rules for protected documentation paths
CREATE TABLE doc_access (
    id TEXT PRIMARY KEY,
    site_id TEXT NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    path_pattern TEXT NOT NULL,
    required_role TEXT NOT NULL DEFAULT 'viewer',
    description TEXT,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(site_id, path_pattern)
);

CREATE INDEX idx_doc_access_site ON doc_access(site_id);
CREATE INDEX idx_doc_access_path ON doc_access(path_pattern);
