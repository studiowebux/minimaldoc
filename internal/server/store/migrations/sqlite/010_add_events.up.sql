CREATE TABLE events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    site_id TEXT NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    category TEXT,
    path TEXT,
    value TEXT,
    session_hash TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX idx_events_site_created ON events(site_id, created_at DESC);
CREATE INDEX idx_events_name ON events(site_id, name, created_at);
