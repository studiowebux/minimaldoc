package api

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
	"github.com/studiowebux/minimaldoc/internal/server/config"
	"github.com/studiowebux/minimaldoc/internal/server/store"
)

// bootstrapMu serializes bootstrap attempts to prevent the TOCTOU race
// where two concurrent requests both pass the "no sites exist" check.
var bootstrapMu sync.Mutex

// BootstrapResult contains the created site and user info.
type BootstrapResult struct {
	SiteID   string
	SiteName string
	APIKey   string
	UserID   string
	Email    string
	Password string
}

// Bootstrap creates the initial site and admin user.
// Returns the generated API key and credentials.
// Serialized with bootstrapMu to prevent duplicate site creation from concurrent requests.
func Bootstrap(ctx context.Context, db store.Store, cfg *config.Config, siteName, domain, adminEmail, adminPassword string) (*BootstrapResult, error) {
	bootstrapMu.Lock()
	defer bootstrapMu.Unlock()

	// Check if any sites exist (under lock, so no race)
	sites, err := db.ListSites(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing sites: %w", err)
	}
	if len(sites) > 0 {
		return nil, fmt.Errorf("bootstrap already done: %d site(s) exist", len(sites))
	}

	// Generate IDs and keys
	siteID, err := auth.GenerateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate site ID: %w", err)
	}

	apiKey, err := auth.GenerateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}
	apiKeyHash := auth.HashAPIKey(apiKey)

	userID, err := auth.GenerateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate user ID: %w", err)
	}

	// Generate password if not provided
	if adminPassword == "" {
		adminPassword, err = auth.GenerateAPIKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate password: %w", err)
		}
		adminPassword = adminPassword[:16] // Shorter for usability
	}

	passwordHash, err := auth.HashPassword(adminPassword, cfg.Auth.BCryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create site
	_, err = db.CreateSite(ctx, siteID, siteName, domain, apiKeyHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create site: %w", err)
	}
	slog.Info("created site", "name", siteName, "id", siteID)

	// Create admin user
	_, err = db.CreateUser(ctx, userID, siteID, adminEmail, passwordHash, "admin", "Admin")
	if err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}
	slog.Info("created admin user", "email", adminEmail)

	return &BootstrapResult{
		SiteID:   siteID,
		SiteName: siteName,
		APIKey:   apiKey,
		UserID:   userID,
		Email:    adminEmail,
		Password: adminPassword,
	}, nil
}
