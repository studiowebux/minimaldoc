package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
)

// Site represents a registered site.
type Site struct {
	ID        string
	Name      string
	Domain    sql.NullString
	APIKey    string
	Config    string
	CreatedAt string // SQLite returns string
	UpdatedAt string
}

// User represents a user account.
type User struct {
	ID            string
	SiteID        string
	Email         string
	PasswordHash  sql.NullString
	Role          string
	OAuthProvider sql.NullString
	OAuthID       sql.NullString
	Name          sql.NullString
	AvatarURL     sql.NullString
	EmailVerified bool
	CreatedAt     string
	UpdatedAt     string
	LastLoginAt   sql.NullString
}

// Session represents a user session.
type Session struct {
	ID        string
	UserID    string
	TokenHash string
	IPAddress sql.NullString
	UserAgent sql.NullString
	ExpiresAt string
	CreatedAt string
}

// PageView represents an analytics page view.
type PageView struct {
	ID          int64
	SiteID      string
	Path        string
	Referrer    sql.NullString
	Country     sql.NullString
	DeviceType  sql.NullString
	Browser     sql.NullString
	OS          sql.NullString
	SessionHash sql.NullString
	CreatedAt   string
}

// Rating represents page feedback.
type Rating struct {
	ID          int64
	SiteID      string
	Path        string
	Rating      int
	Feedback    sql.NullString
	SessionHash sql.NullString
	CreatedAt   string
}

// Subscriber represents a newsletter subscriber.
type Subscriber struct {
	ID             string
	SiteID         string
	Email          string
	Verified       bool
	VerifyToken    sql.NullString
	VerifySentAt   sql.NullString
	SubscribedAt   string
	UnsubscribedAt sql.NullString
}

// Site queries

func (db *DB) CreateSite(ctx context.Context, id, name, domain, apiKeyHash string) (*Site, error) {
	query := `
		INSERT INTO sites (id, name, domain, api_key, config)
		VALUES ($1, $2, $3, $4, '{}')
	`
	_, err := db.ExecContext(ctx, query, id, name, nullString(domain), apiKeyHash)
	if err != nil {
		return nil, err
	}

	// Fetch the created site
	return db.GetSiteByID(ctx, id)
}

func (db *DB) GetSiteByID(ctx context.Context, id string) (*Site, error) {
	query := `SELECT id, name, domain, api_key, config, created_at, updated_at FROM sites WHERE id = $1`
	var s Site
	err := db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.Name, &s.Domain, &s.APIKey, &s.Config, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (db *DB) GetSiteByAPIKey(ctx context.Context, apiKeyHash string) (*Site, error) {
	query := `SELECT id, name, domain, api_key, config, created_at, updated_at FROM sites WHERE api_key = $1`
	var s Site
	err := db.QueryRowContext(ctx, query, apiKeyHash).Scan(
		&s.ID, &s.Name, &s.Domain, &s.APIKey, &s.Config, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (db *DB) ListSites(ctx context.Context) ([]Site, error) {
	query := `SELECT id, name, domain, api_key, config, created_at, updated_at FROM sites ORDER BY created_at DESC`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []Site
	for rows.Next() {
		var s Site
		if err := rows.Scan(&s.ID, &s.Name, &s.Domain, &s.APIKey, &s.Config, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sites = append(sites, s)
	}
	return sites, rows.Err()
}

func (db *DB) UpdateSiteAPIKey(ctx context.Context, siteID, newAPIKeyHash string) error {
	query := `UPDATE sites SET api_key = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := db.ExecContext(ctx, query, newAPIKeyHash, siteID)
	return err
}

func (db *DB) DeleteSite(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM sites WHERE id = $1`, id)
	return err
}

// User queries

func (db *DB) CreateUser(ctx context.Context, id, siteID, email, passwordHash, role, name string) (*User, error) {
	query := `
		INSERT INTO users (id, site_id, email, password_hash, role, name, email_verified)
		VALUES ($1, $2, $3, $4, $5, $6, false)
		RETURNING id, site_id, email, password_hash, role, name, email_verified, created_at, updated_at
	`
	var u User
	err := db.QueryRowContext(ctx, query, id, siteID, email, nullString(passwordHash), role, nullString(name)).Scan(
		&u.ID, &u.SiteID, &u.Email, &u.PasswordHash, &u.Role, &u.Name, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt,
	)
	return &u, err
}

func (db *DB) GetUserByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, site_id, email, password_hash, role, oauth_provider, oauth_id, name, avatar_url, email_verified, created_at, updated_at, last_login_at
		FROM users WHERE id = $1
	`
	var u User
	err := db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.SiteID, &u.Email, &u.PasswordHash, &u.Role, &u.OAuthProvider, &u.OAuthID,
		&u.Name, &u.AvatarURL, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (db *DB) GetUserByEmail(ctx context.Context, siteID, email string) (*User, error) {
	query := `
		SELECT id, site_id, email, password_hash, role, oauth_provider, oauth_id, name, avatar_url, email_verified, created_at, updated_at, last_login_at
		FROM users WHERE site_id = $1 AND email = $2
	`
	var u User
	err := db.QueryRowContext(ctx, query, siteID, email).Scan(
		&u.ID, &u.SiteID, &u.Email, &u.PasswordHash, &u.Role, &u.OAuthProvider, &u.OAuthID,
		&u.Name, &u.AvatarURL, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (db *DB) GetUserByOAuth(ctx context.Context, provider, providerID string) (*User, error) {
	query := `
		SELECT id, site_id, email, password_hash, role, oauth_provider, oauth_id, name, avatar_url, email_verified, created_at, updated_at, last_login_at
		FROM users WHERE oauth_provider = $1 AND oauth_id = $2
	`
	var u User
	err := db.QueryRowContext(ctx, query, provider, providerID).Scan(
		&u.ID, &u.SiteID, &u.Email, &u.PasswordHash, &u.Role, &u.OAuthProvider, &u.OAuthID,
		&u.Name, &u.AvatarURL, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (db *DB) ListUsers(ctx context.Context, siteID string) ([]User, error) {
	query := `
		SELECT id, site_id, email, password_hash, role, oauth_provider, oauth_id, name, avatar_url, email_verified, created_at, updated_at, last_login_at
		FROM users WHERE site_id = $1 ORDER BY created_at DESC
	`
	rows, err := db.QueryContext(ctx, query, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.SiteID, &u.Email, &u.PasswordHash, &u.Role, &u.OAuthProvider, &u.OAuthID,
			&u.Name, &u.AvatarURL, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (db *DB) UpdateUserLastLogin(ctx context.Context, userID string) error {
	query := `UPDATE users SET last_login_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := db.ExecContext(ctx, query, userID)
	return err
}

func (db *DB) DeleteUser(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

// Session queries

func (db *DB) CreateSession(ctx context.Context, id, userID, tokenHash, ipAddress, userAgent string, expiresAt time.Time) error {
	query := `INSERT INTO sessions (id, user_id, token_hash, ip_address, user_agent, expires_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.ExecContext(ctx, query, id, userID, tokenHash, nullString(ipAddress), nullString(userAgent), expiresAt.Format(time.RFC3339))
	return err
}

func (db *DB) GetSessionByToken(ctx context.Context, tokenHash string) (*Session, error) {
	query := `SELECT id, user_id, token_hash, ip_address, user_agent, expires_at, created_at FROM sessions WHERE token_hash = $1 AND expires_at > CURRENT_TIMESTAMP`
	var s Session
	err := db.QueryRowContext(ctx, query, tokenHash).Scan(
		&s.ID, &s.UserID, &s.TokenHash, &s.IPAddress, &s.UserAgent, &s.ExpiresAt, &s.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (db *DB) DeleteSession(ctx context.Context, tokenHash string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM sessions WHERE token_hash = $1`, tokenHash)
	return err
}

func (db *DB) DeleteUserSessions(ctx context.Context, userID string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}

func (db *DB) CleanExpiredSessions(ctx context.Context) (int64, error) {
	result, err := db.ExecContext(ctx, `DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP`)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

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

// Helper to create user info from DB user
func (u *User) ToAuthInfo() *auth.UserInfo {
	return &auth.UserInfo{
		ID:            u.ID,
		Email:         u.Email,
		Name:          u.Name.String,
		AvatarURL:     u.AvatarURL.String,
		EmailVerified: u.EmailVerified,
		Provider:      u.OAuthProvider.String,
		ProviderID:    u.OAuthID.String,
		Role:          u.Role,
	}
}

// Helper for nullable strings
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
