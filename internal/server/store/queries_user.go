package store

import (
	"context"
	"database/sql"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
)

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

func (db *DB) GetUserByOAuth(ctx context.Context, siteID, provider, providerID string) (*User, error) {
	query := `
		SELECT id, site_id, email, password_hash, role, oauth_provider, oauth_id, name, avatar_url, email_verified, created_at, updated_at, last_login_at
		FROM users WHERE site_id = $1 AND oauth_provider = $2 AND oauth_id = $3
	`
	var u User
	err := db.QueryRowContext(ctx, query, siteID, provider, providerID).Scan(
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

func (db *DB) UpdateUser(ctx context.Context, id, email, role, name string) error {
	query := `UPDATE users SET email = $1, role = $2, name = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`
	_, err := db.ExecContext(ctx, query, email, role, nullString(name), id)
	return err
}

func (db *DB) UpdateUserPassword(ctx context.Context, id, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := db.ExecContext(ctx, query, passwordHash, id)
	return err
}

func (db *DB) CountUsers(ctx context.Context, siteID string) (int, error) {
	var count int
	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE site_id = $1`, siteID).Scan(&count)
	return count, err
}

func (db *DB) CountUsersByRole(ctx context.Context, siteID, role string) (int, error) {
	var count int
	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE site_id = $1 AND role = $2`, siteID, role).Scan(&count)
	return count, err
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

// CreateUserWithVerification creates a new user with a verification token.
func (db *DB) CreateUserWithVerification(ctx context.Context, id, siteID, email, passwordHash, role, name, verifyToken string) (*User, error) {
	query := `
		INSERT INTO users (id, site_id, email, password_hash, role, name, email_verified, verify_token)
		VALUES ($1, $2, $3, $4, $5, $6, false, $7)
		RETURNING id, site_id, email, password_hash, role, name, email_verified, verify_token, created_at, updated_at
	`
	var u User
	err := db.QueryRowContext(ctx, query, id, siteID, email, nullString(passwordHash), role, nullString(name), nullString(verifyToken)).Scan(
		&u.ID, &u.SiteID, &u.Email, &u.PasswordHash, &u.Role, &u.Name, &u.EmailVerified, &u.VerifyToken, &u.CreatedAt, &u.UpdatedAt,
	)
	return &u, err
}

// GetUserByVerifyToken finds a user by their verification token.
func (db *DB) GetUserByVerifyToken(ctx context.Context, token string) (*User, error) {
	query := `
		SELECT id, site_id, email, password_hash, role, oauth_provider, oauth_id, name, avatar_url, email_verified, verify_token, created_at, updated_at, last_login_at
		FROM users WHERE verify_token = $1
	`
	var u User
	err := db.QueryRowContext(ctx, query, token).Scan(
		&u.ID, &u.SiteID, &u.Email, &u.PasswordHash, &u.Role, &u.OAuthProvider, &u.OAuthID,
		&u.Name, &u.AvatarURL, &u.EmailVerified, &u.VerifyToken, &u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

// VerifyUserEmail marks a user's email as verified and clears the token.
func (db *DB) VerifyUserEmail(ctx context.Context, userID string) error {
	query := `UPDATE users SET email_verified = true, verify_token = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := db.ExecContext(ctx, query, userID)
	return err
}

// CreateUserWithOAuth creates a new user with OAuth credentials.
func (db *DB) CreateUserWithOAuth(ctx context.Context, id, siteID, email, oauthProvider, oauthID, name, avatarURL, role string) (*User, error) {
	query := `
		INSERT INTO users (id, site_id, email, oauth_provider, oauth_id, name, avatar_url, role, email_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true)
		RETURNING id, site_id, email, password_hash, role, oauth_provider, oauth_id, name, avatar_url, email_verified, created_at, updated_at
	`
	var u User
	err := db.QueryRowContext(ctx, query, id, siteID, email, nullString(oauthProvider), nullString(oauthID), nullString(name), nullString(avatarURL), role).Scan(
		&u.ID, &u.SiteID, &u.Email, &u.PasswordHash, &u.Role, &u.OAuthProvider, &u.OAuthID,
		&u.Name, &u.AvatarURL, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt,
	)
	return &u, err
}

// LinkOAuthToUser links OAuth credentials to an existing user.
func (db *DB) LinkOAuthToUser(ctx context.Context, userID, oauthProvider, oauthID string) error {
	query := `UPDATE users SET oauth_provider = $1, oauth_id = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	_, err := db.ExecContext(ctx, query, oauthProvider, oauthID, userID)
	return err
}

// SetPasswordResetToken stores a hashed reset token with 15-minute expiry.
func (db *DB) SetPasswordResetToken(ctx context.Context, userID, tokenHash string) error {
	query := `UPDATE users SET reset_token = $1, reset_token_expires_at = datetime('now', '+15 minutes'), updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	if db.driver == "postgres" {
		query = `UPDATE users SET reset_token = $1, reset_token_expires_at = NOW() + INTERVAL '15 minutes', updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	}
	_, err := db.ExecContext(ctx, query, tokenHash, userID)
	return err
}

// GetUserByResetToken finds a user by their hashed reset token if not expired.
func (db *DB) GetUserByResetToken(ctx context.Context, tokenHash string) (*User, error) {
	query := `
		SELECT id, site_id, email, password_hash, role, oauth_provider, oauth_id, name, avatar_url, email_verified, verify_token, created_at, updated_at, last_login_at
		FROM users WHERE reset_token = $1 AND reset_token_expires_at > CURRENT_TIMESTAMP
	`
	var u User
	err := db.QueryRowContext(ctx, query, tokenHash).Scan(
		&u.ID, &u.SiteID, &u.Email, &u.PasswordHash, &u.Role, &u.OAuthProvider, &u.OAuthID,
		&u.Name, &u.AvatarURL, &u.EmailVerified, &u.VerifyToken, &u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

// ClearPasswordResetToken clears the reset token after use.
func (db *DB) ClearPasswordResetToken(ctx context.Context, userID string) error {
	query := `UPDATE users SET reset_token = NULL, reset_token_expires_at = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := db.ExecContext(ctx, query, userID)
	return err
}
