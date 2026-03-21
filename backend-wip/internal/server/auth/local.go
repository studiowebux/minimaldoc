package auth

import (
	"context"
	"database/sql"
	"fmt"
)

// LocalProvider implements email/password authentication.
type LocalProvider struct {
	db         *sql.DB
	bcryptCost int
}

// NewLocalProvider creates a new local authentication provider.
func NewLocalProvider(db *sql.DB, bcryptCost int) *LocalProvider {
	return &LocalProvider{
		db:         db,
		bcryptCost: bcryptCost,
	}
}

// Name returns the provider identifier.
func (p *LocalProvider) Name() string {
	return "local"
}

// Authenticate validates email/password credentials.
func (p *LocalProvider) Authenticate(ctx context.Context, credentials map[string]string) (*UserInfo, error) {
	email := credentials["email"]
	password := credentials["password"]
	siteID := credentials["site_id"]

	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password required")
	}

	// Query user from database
	var user struct {
		ID           string
		Email        string
		PasswordHash string
		Name         sql.NullString
		AvatarURL    sql.NullString
		Role         string
	}

	query := `
		SELECT id, email, password_hash, name, avatar_url, role
		FROM users
		WHERE email = $1 AND site_id = $2 AND password_hash IS NOT NULL
	`

	err := p.db.QueryRowContext(ctx, query, email, siteID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.AvatarURL,
		&user.Role,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid credentials")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Verify password
	if !VerifyPassword(password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &UserInfo{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name.String,
		AvatarURL:     user.AvatarURL.String,
		EmailVerified: true,
		Provider:      "local",
		Role:          user.Role,
	}, nil
}

// GetAuthURL is not applicable for local auth.
func (p *LocalProvider) GetAuthURL(state string) string {
	return ""
}

// HandleCallback is not applicable for local auth.
func (p *LocalProvider) HandleCallback(ctx context.Context, code string) (*UserInfo, error) {
	return nil, fmt.Errorf("callback not supported for local auth")
}

// CreateUser creates a new user with email/password.
func (p *LocalProvider) CreateUser(ctx context.Context, siteID, email, password, name, role string) (*UserInfo, error) {
	// Hash password
	hash, err := HashPassword(password, p.bcryptCost)
	if err != nil {
		return nil, err
	}

	// Generate UUID (simplified, should use proper UUID library)
	id, err := GenerateVerificationToken()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users (id, site_id, email, password_hash, name, role, email_verified)
		VALUES ($1, $2, $3, $4, $5, $6, false)
		RETURNING id
	`

	err = p.db.QueryRowContext(ctx, query, id, siteID, email, hash, name, role).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &UserInfo{
		ID:            id,
		Email:         email,
		Name:          name,
		EmailVerified: true,
		Provider:      "local",
		Role:          role,
	}, nil
}

// UpdatePassword changes a user's password.
func (p *LocalProvider) UpdatePassword(ctx context.Context, userID, newPassword string) error {
	hash, err := HashPassword(newPassword, p.bcryptCost)
	if err != nil {
		return err
	}

	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	_, err = p.db.ExecContext(ctx, query, hash, userID)
	return err
}
