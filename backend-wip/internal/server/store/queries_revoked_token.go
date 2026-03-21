package store

import (
	"context"
	"database/sql"
	"time"
)

func (db *DB) RevokeToken(ctx context.Context, jti string, expiresAt time.Time) error {
	query := `INSERT INTO revoked_tokens (jti, expires_at) VALUES ($1, $2) ON CONFLICT (jti) DO NOTHING`
	_, err := db.ExecContext(ctx, query, jti, expiresAt.Format(time.RFC3339))
	return err
}

func (db *DB) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	var dummy string
	err := db.QueryRowContext(ctx, `SELECT jti FROM revoked_tokens WHERE jti = $1`, jti).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *DB) CleanRevokedTokens(ctx context.Context) (int64, error) {
	result, err := db.ExecContext(ctx, `DELETE FROM revoked_tokens WHERE expires_at < CURRENT_TIMESTAMP`)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
