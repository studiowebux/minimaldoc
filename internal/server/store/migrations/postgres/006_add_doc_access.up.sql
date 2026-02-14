-- Access rules for protected documentation paths
CREATE TABLE doc_access (
    id UUID PRIMARY KEY,
    site_id UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    path_pattern TEXT NOT NULL,
    required_role TEXT NOT NULL DEFAULT 'viewer' CHECK (required_role IN ('viewer', 'author', 'editor', 'admin')),
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(site_id, path_pattern)
);

CREATE INDEX idx_doc_access_site ON doc_access(site_id);
CREATE INDEX idx_doc_access_path ON doc_access(path_pattern);
