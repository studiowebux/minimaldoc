package store

import (
	"context"
	"database/sql"
)

// Newsletter queries

func (db *DB) CreateSubscriber(ctx context.Context, id, siteID, email, verifyToken string) error {
	query := `INSERT INTO subscribers (id, site_id, email, verify_token, verify_sent_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)`
	_, err := db.ExecContext(ctx, query, id, siteID, email, verifyToken)
	return err
}

func (db *DB) GetSubscriberByEmail(ctx context.Context, siteID, email string) (*Subscriber, error) {
	query := `
		SELECT id, site_id, email, verified, verify_token, verify_sent_at, subscribed_at, unsubscribed_at
		FROM subscribers WHERE site_id = $1 AND email = $2
	`
	var s Subscriber
	err := db.QueryRowContext(ctx, query, siteID, email).Scan(
		&s.ID, &s.SiteID, &s.Email, &s.Verified, &s.VerifyToken, &s.VerifySentAt, &s.SubscribedAt, &s.UnsubscribedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (db *DB) GetSubscriberByToken(ctx context.Context, siteID, token string) (*Subscriber, error) {
	query := `
		SELECT id, site_id, email, verified, verify_token, verify_sent_at, subscribed_at, unsubscribed_at
		FROM subscribers WHERE site_id = $1 AND verify_token = $2
	`
	var s Subscriber
	err := db.QueryRowContext(ctx, query, siteID, token).Scan(
		&s.ID, &s.SiteID, &s.Email, &s.Verified, &s.VerifyToken, &s.VerifySentAt, &s.SubscribedAt, &s.UnsubscribedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (db *DB) VerifySubscriber(ctx context.Context, siteID, token string) error {
	query := `UPDATE subscribers SET verified = true, verify_token = NULL WHERE site_id = $1 AND verify_token = $2`
	result, err := db.ExecContext(ctx, query, siteID, token)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (db *DB) UnsubscribeByEmail(ctx context.Context, siteID, email string) error {
	query := `UPDATE subscribers SET unsubscribed_at = CURRENT_TIMESTAMP WHERE site_id = $1 AND email = $2`
	_, err := db.ExecContext(ctx, query, siteID, email)
	return err
}

func (db *DB) ListSubscribers(ctx context.Context, siteID string, verifiedOnly bool) ([]Subscriber, error) {
	query := `
		SELECT id, site_id, email, verified, verify_token, verify_sent_at, subscribed_at, unsubscribed_at
		FROM subscribers WHERE site_id = $1 AND unsubscribed_at IS NULL
	`
	if verifiedOnly {
		query += ` AND verified = true`
	}
	query += ` ORDER BY subscribed_at DESC`

	rows, err := db.QueryContext(ctx, query, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []Subscriber
	for rows.Next() {
		var s Subscriber
		if err := rows.Scan(&s.ID, &s.SiteID, &s.Email, &s.Verified, &s.VerifyToken, &s.VerifySentAt, &s.SubscribedAt, &s.UnsubscribedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, rows.Err()
}

func (db *DB) CountSubscribers(ctx context.Context, siteID string, verifiedOnly bool) (int64, error) {
	query := `SELECT COUNT(*) FROM subscribers WHERE site_id = $1 AND unsubscribed_at IS NULL`
	if verifiedOnly {
		query += ` AND verified = true`
	}
	var count int64
	err := db.QueryRowContext(ctx, query, siteID).Scan(&count)
	return count, err
}

func (db *DB) GetSubscriberStatsExtended(ctx context.Context, siteID string) (total, verified, pending, thisMonth int64, err error) {
	query := `
		SELECT
			COUNT(*),
			COUNT(CASE WHEN verified = true THEN 1 END),
			COUNT(CASE WHEN verified = false THEN 1 END),
			COUNT(CASE WHEN subscribed_at >= datetime('now', '-30 days') THEN 1 END)
		FROM subscribers WHERE site_id = $1 AND unsubscribed_at IS NULL
	`
	err = db.QueryRowContext(ctx, query, siteID).Scan(&total, &verified, &pending, &thisMonth)
	return
}
