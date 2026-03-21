package store

import (
	"context"
	"database/sql"
)

// subscriberColumns is the canonical column list for subscriber SELECT queries.
const subscriberColumns = `id, site_id, email, verified, verify_token, verify_sent_at, subscribed_at, unsubscribed_at, oauth_provider, oauth_display_name, verified_via`

// scanSubscriber scans a row into a Subscriber struct using the canonical column order.
func scanSubscriber(scanner interface{ Scan(...any) error }) (Subscriber, error) {
	var s Subscriber
	err := scanner.Scan(
		&s.ID, &s.SiteID, &s.Email, &s.Verified, &s.VerifyToken, &s.VerifySentAt,
		&s.SubscribedAt, &s.UnsubscribedAt,
		&s.OAuthProvider, &s.OAuthDisplayName, &s.VerifiedVia,
	)
	return s, err
}

// Newsletter queries

func (db *DB) CreateSubscriber(ctx context.Context, id, siteID, email, verifyTokenHash string) error {
	query := `INSERT INTO subscribers (id, site_id, email, verify_token, verify_sent_at, verify_token_expires_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, datetime('now', '+24 hours'))`
	if db.driver == "postgres" {
		query = `INSERT INTO subscribers (id, site_id, email, verify_token, verify_sent_at, verify_token_expires_at) VALUES ($1, $2, $3, $4, NOW(), NOW() + INTERVAL '24 hours')`
	}
	_, err := db.ExecContext(ctx, query, id, siteID, email, verifyTokenHash)
	return err
}

// CreateVerifiedSubscriber inserts or upgrades a subscriber as OAuth-verified.
// ON CONFLICT upgrades an existing email-only subscriber to OAuth-verified.
func (db *DB) CreateVerifiedSubscriber(ctx context.Context, id, siteID, email, provider, displayName string) error {
	query := `
		INSERT INTO subscribers (id, site_id, email, verified, verify_token, oauth_provider, oauth_display_name, verified_via)
		VALUES ($1, $2, $3, true, NULL, $4, $5, $4)
		ON CONFLICT (site_id, email) DO UPDATE SET
			verified = true,
			verify_token = NULL,
			oauth_provider = EXCLUDED.oauth_provider,
			oauth_display_name = EXCLUDED.oauth_display_name,
			verified_via = EXCLUDED.verified_via
	`
	_, err := db.ExecContext(ctx, query, id, siteID, email, provider, displayName)
	return err
}

func (db *DB) GetSubscriberByEmail(ctx context.Context, siteID, email string) (*Subscriber, error) {
	query := `SELECT ` + subscriberColumns + ` FROM subscribers WHERE site_id = $1 AND email = $2`
	s, err := scanSubscriber(db.QueryRowContext(ctx, query, siteID, email))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (db *DB) GetSubscriberByToken(ctx context.Context, siteID, tokenHash string) (*Subscriber, error) {
	query := `SELECT ` + subscriberColumns + ` FROM subscribers WHERE site_id = $1 AND verify_token = $2 AND verify_token_expires_at > CURRENT_TIMESTAMP`
	s, err := scanSubscriber(db.QueryRowContext(ctx, query, siteID, tokenHash))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (db *DB) VerifySubscriber(ctx context.Context, siteID, tokenHash string) error {
	query := `UPDATE subscribers SET verified = true, verify_token = NULL, verify_token_expires_at = NULL WHERE site_id = $1 AND verify_token = $2 AND verify_token_expires_at > CURRENT_TIMESTAMP`
	result, err := db.ExecContext(ctx, query, siteID, tokenHash)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (db *DB) UpdateSubscriberToken(ctx context.Context, siteID, email, newTokenHash string) error {
	query := `UPDATE subscribers SET verify_token = $1, verify_sent_at = CURRENT_TIMESTAMP, verify_token_expires_at = datetime('now', '+24 hours') WHERE site_id = $2 AND email = $3 AND verified = false`
	if db.driver == "postgres" {
		query = `UPDATE subscribers SET verify_token = $1, verify_sent_at = NOW(), verify_token_expires_at = NOW() + INTERVAL '24 hours' WHERE site_id = $2 AND email = $3 AND verified = false`
	}
	_, err := db.ExecContext(ctx, query, newTokenHash, siteID, email)
	return err
}

func (db *DB) UnsubscribeByEmail(ctx context.Context, siteID, email string) error {
	query := `UPDATE subscribers SET unsubscribed_at = CURRENT_TIMESTAMP WHERE site_id = $1 AND email = $2`
	_, err := db.ExecContext(ctx, query, siteID, email)
	return err
}

func (db *DB) ListSubscribers(ctx context.Context, siteID string, verifiedOnly bool) ([]Subscriber, error) {
	query := `SELECT ` + subscriberColumns + ` FROM subscribers WHERE site_id = $1 AND unsubscribed_at IS NULL`
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
		s, err := scanSubscriber(rows)
		if err != nil {
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
