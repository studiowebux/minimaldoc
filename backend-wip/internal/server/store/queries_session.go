package store

import (
	"context"
	"database/sql"
	"time"
)

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
