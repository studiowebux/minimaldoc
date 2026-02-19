-- Audit logs table for tracking admin actions
CREATE TABLE audit_logs (
    id TEXT PRIMARY KEY,
    site_id TEXT REFERENCES sites(id) ON DELETE CASCADE,
    user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    user_email TEXT,
    action TEXT NOT NULL CHECK (action IN ('create', 'update', 'delete', 'publish', 'unpublish', 'approve', 'reject', 'login', 'logout')),
    entity_type TEXT NOT NULL,
    entity_id TEXT,
    entity_name TEXT,
    details TEXT,
    ip_address TEXT,
    user_agent TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

-- Index for listing logs by site and time (most common query)
CREATE INDEX idx_audit_logs_site_created ON audit_logs(site_id, created_at DESC);

-- Index for filtering by user
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);

-- Index for filtering by entity
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);

-- Index for filtering by action
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
