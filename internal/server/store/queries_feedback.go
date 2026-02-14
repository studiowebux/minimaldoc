package store

import (
	"context"
)

// Feedback queries

func (db *DB) RecordRating(ctx context.Context, siteID, path string, rating int, feedback, sessionHash string) error {
	query := `INSERT INTO ratings (site_id, path, rating, feedback, session_hash) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.ExecContext(ctx, query, siteID, path, rating, nullString(feedback), nullString(sessionHash))
	return err
}

func (db *DB) GetRatingStats(ctx context.Context, siteID string) (avgRating float64, totalRatings int64, err error) {
	query := `SELECT COALESCE(AVG(rating), 0), COUNT(*) FROM ratings WHERE site_id = $1`
	err = db.QueryRowContext(ctx, query, siteID).Scan(&avgRating, &totalRatings)
	return
}

func (db *DB) GetRatingStatsExtended(ctx context.Context, siteID string) (avgRating float64, totalRatings, withComments, thisWeek int64, err error) {
	query := `
		SELECT
			COALESCE(AVG(rating), 0),
			COUNT(*),
			COUNT(CASE WHEN feedback IS NOT NULL AND feedback != '' THEN 1 END),
			COUNT(CASE WHEN created_at >= datetime('now', '-7 days') THEN 1 END)
		FROM ratings WHERE site_id = $1
	`
	err = db.QueryRowContext(ctx, query, siteID).Scan(&avgRating, &totalRatings, &withComments, &thisWeek)
	return
}

func (db *DB) GetRatingsByPath(ctx context.Context, siteID, path string) (avgRating float64, count int64, err error) {
	query := `SELECT COALESCE(AVG(rating), 0), COUNT(*) FROM ratings WHERE site_id = $1 AND path = $2`
	err = db.QueryRowContext(ctx, query, siteID, path).Scan(&avgRating, &count)
	return
}

func (db *DB) ListRatings(ctx context.Context, siteID string, limit, offset int) ([]Rating, error) {
	query := `
		SELECT id, site_id, path, rating, feedback, session_hash, created_at
		FROM ratings WHERE site_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := db.QueryContext(ctx, query, siteID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []Rating
	for rows.Next() {
		var r Rating
		if err := rows.Scan(&r.ID, &r.SiteID, &r.Path, &r.Rating, &r.Feedback, &r.SessionHash, &r.CreatedAt); err != nil {
			return nil, err
		}
		ratings = append(ratings, r)
	}
	return ratings, rows.Err()
}
