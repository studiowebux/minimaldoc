package store

import (
	"context"
	"time"
)

// Event queries

// RecordEvent stores a custom analytics event.
func (db *DB) RecordEvent(ctx context.Context, siteID, name, category, path, value, sessionHash string) error {
	query := `
		INSERT INTO events (site_id, name, category, path, value, session_hash)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := db.ExecContext(ctx, query, siteID, name,
		nullString(category), nullString(path), nullString(value), nullString(sessionHash))
	return err
}

// EventStat represents aggregated event counts.
type EventStat struct {
	Name  string
	Count int64
}

// GetEventStats returns event counts grouped by name for a site since a given time.
func (db *DB) GetEventStats(ctx context.Context, siteID string, since time.Time) ([]EventStat, error) {
	query := `
		SELECT name, COUNT(*) as count
		FROM events
		WHERE site_id = $1 AND created_at >= $2
		GROUP BY name
		ORDER BY count DESC
		LIMIT 20
	`
	rows, err := db.QueryContext(ctx, query, siteID, since.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []EventStat
	for rows.Next() {
		var s EventStat
		if err := rows.Scan(&s.Name, &s.Count); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}

// GetEventsByName returns events with a specific name for a site.
func (db *DB) GetEventsByName(ctx context.Context, siteID, name string, limit, offset int) ([]Event, error) {
	query := `
		SELECT id, site_id, name, category, path, value, session_hash, created_at
		FROM events
		WHERE site_id = $1 AND name = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := db.QueryContext(ctx, query, siteID, name, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.SiteID, &e.Name, &e.Category, &e.Path, &e.Value, &e.SessionHash, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

// ListRecentEvents returns the most recent events for a site.
func (db *DB) ListRecentEvents(ctx context.Context, siteID string, limit int) ([]Event, error) {
	query := `
		SELECT id, site_id, name, category, path, value, session_hash, created_at
		FROM events
		WHERE site_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := db.QueryContext(ctx, query, siteID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.SiteID, &e.Name, &e.Category, &e.Path, &e.Value, &e.SessionHash, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

// GetTotalEventCount returns the total number of events for a site since a given time.
func (db *DB) GetTotalEventCount(ctx context.Context, siteID string, since time.Time) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM events
		WHERE site_id = $1 AND created_at >= $2
	`
	var count int64
	err := db.QueryRowContext(ctx, query, siteID, since.Format(time.RFC3339)).Scan(&count)
	return count, err
}

// GetUniqueEventNames returns the count of unique event names for a site since a given time.
func (db *DB) GetUniqueEventNames(ctx context.Context, siteID string, since time.Time) (int64, error) {
	query := `
		SELECT COUNT(DISTINCT name)
		FROM events
		WHERE site_id = $1 AND created_at >= $2
	`
	var count int64
	err := db.QueryRowContext(ctx, query, siteID, since.Format(time.RFC3339)).Scan(&count)
	return count, err
}
