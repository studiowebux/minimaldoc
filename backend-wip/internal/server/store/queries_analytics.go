package store

import (
	"context"
	"time"
)

// Analytics queries

func (db *DB) RecordPageView(ctx context.Context, siteID, path, referrer, country, deviceType, browser, os, sessionHash string) error {
	query := `
		INSERT INTO page_views (site_id, path, referrer, country, device_type, browser, os, session_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := db.ExecContext(ctx, query, siteID, path,
		nullString(referrer), nullString(country), nullString(deviceType),
		nullString(browser), nullString(os), nullString(sessionHash))
	return err
}

func (db *DB) GetPageViewStats(ctx context.Context, siteID string, since time.Time) (totalViews int64, uniqueSessions int64, err error) {
	query := `
		SELECT COUNT(*), COUNT(DISTINCT session_hash)
		FROM page_views
		WHERE site_id = $1 AND created_at >= $2
	`
	err = db.QueryRowContext(ctx, query, siteID, since.Format(time.RFC3339)).Scan(&totalViews, &uniqueSessions)
	return
}

func (db *DB) GetPageViewStatsExtended(ctx context.Context, siteID string, since time.Time) (totalViews, uniqueSessions int64, avgDuration float64, err error) {
	query := `
		SELECT
			COUNT(*),
			COUNT(DISTINCT session_hash),
			COALESCE(AVG(CASE WHEN duration_seconds > 0 THEN duration_seconds END), 0)
		FROM page_views
		WHERE site_id = $1 AND created_at >= $2
	`
	err = db.QueryRowContext(ctx, query, siteID, since.Format(time.RFC3339)).Scan(&totalViews, &uniqueSessions, &avgDuration)
	return
}

func (db *DB) UpdatePageViewDuration(ctx context.Context, siteID, path, sessionHash string, duration int) error {
	// Update the most recent page view for this session/path combination
	query := `
		UPDATE page_views
		SET duration_seconds = $1
		WHERE id = (
			SELECT id FROM page_views
			WHERE site_id = $2 AND path = $3 AND session_hash = $4
			ORDER BY created_at DESC
			LIMIT 1
		)
	`
	_, err := db.ExecContext(ctx, query, duration, siteID, path, sessionHash)
	return err
}

func (db *DB) UpdatePageViewDurationAndBounce(ctx context.Context, siteID, path, sessionHash string, duration int, isBounce *bool) error {
	// Update the most recent page view for this session/path combination
	if isBounce != nil && *isBounce {
		query := `
			UPDATE page_views
			SET duration_seconds = $1, is_bounce = 1
			WHERE id = (
				SELECT id FROM page_views
				WHERE site_id = $2 AND path = $3 AND session_hash = $4
				ORDER BY created_at DESC
				LIMIT 1
			)
		`
		_, err := db.ExecContext(ctx, query, duration, siteID, path, sessionHash)
		return err
	}
	return db.UpdatePageViewDuration(ctx, siteID, path, sessionHash, duration)
}

func (db *DB) GetBounceRate(ctx context.Context, siteID string, since time.Time) (float64, error) {
	query := `
		SELECT
			COALESCE(
				CAST(SUM(CASE WHEN is_bounce = 1 THEN 1 ELSE 0 END) AS FLOAT) /
				NULLIF(COUNT(DISTINCT session_hash), 0) * 100,
				0
			)
		FROM page_views
		WHERE site_id = $1 AND created_at >= $2 AND is_bounce IS NOT NULL
	`
	var rate float64
	err := db.QueryRowContext(ctx, query, siteID, since.Format(time.RFC3339)).Scan(&rate)
	return rate, err
}

func (db *DB) GetTopPages(ctx context.Context, siteID string, since time.Time, limit int) ([]struct {
	Path  string
	Views int64
}, error) {
	query := `
		SELECT path, COUNT(*) as views
		FROM page_views
		WHERE site_id = $1 AND created_at >= $2
		GROUP BY path
		ORDER BY views DESC
		LIMIT $3
	`
	rows, err := db.QueryContext(ctx, query, siteID, since.Format(time.RFC3339), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []struct {
		Path  string
		Views int64
	}
	for rows.Next() {
		var p struct {
			Path  string
			Views int64
		}
		if err := rows.Scan(&p.Path, &p.Views); err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	return pages, rows.Err()
}

// TrafficSource represents aggregated referrer data
type TrafficSource struct {
	Source string
	Visits int64
}

func (db *DB) GetTrafficSources(ctx context.Context, siteID string, since time.Time, limit int) ([]TrafficSource, error) {
	query := `
		SELECT
			CASE
				WHEN referrer IS NULL OR referrer = '' THEN 'Direct'
				WHEN referrer LIKE '%google.%' THEN 'Google'
				WHEN referrer LIKE '%bing.%' THEN 'Bing'
				WHEN referrer LIKE '%duckduckgo.%' THEN 'DuckDuckGo'
				WHEN referrer LIKE '%github.%' THEN 'GitHub'
				WHEN referrer LIKE '%twitter.%' OR referrer LIKE '%x.com%' THEN 'Twitter/X'
				WHEN referrer LIKE '%linkedin.%' THEN 'LinkedIn'
				WHEN referrer LIKE '%reddit.%' THEN 'Reddit'
				WHEN referrer LIKE '%facebook.%' THEN 'Facebook'
				ELSE 'Other'
			END as source,
			COUNT(*) as visits
		FROM page_views
		WHERE site_id = $1 AND created_at >= $2
		GROUP BY source
		ORDER BY visits DESC
		LIMIT $3
	`
	rows, err := db.QueryContext(ctx, query, siteID, since.Format(time.RFC3339), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []TrafficSource
	for rows.Next() {
		var s TrafficSource
		if err := rows.Scan(&s.Source, &s.Visits); err != nil {
			return nil, err
		}
		sources = append(sources, s)
	}
	return sources, rows.Err()
}

// DailyViewCount represents page views for a single day
type DailyViewCount struct {
	Date  string
	Views int64
}

func (db *DB) GetDailyViews(ctx context.Context, siteID string, days int) ([]DailyViewCount, error) {
	query := `
		SELECT
			date(created_at) as day,
			COUNT(*) as views
		FROM page_views
		WHERE site_id = $1 AND created_at >= datetime('now', '-' || $2 || ' days')
		GROUP BY day
		ORDER BY day ASC
	`
	rows, err := db.QueryContext(ctx, query, siteID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counts []DailyViewCount
	for rows.Next() {
		var c DailyViewCount
		if err := rows.Scan(&c.Date, &c.Views); err != nil {
			return nil, err
		}
		counts = append(counts, c)
	}
	return counts, rows.Err()
}
