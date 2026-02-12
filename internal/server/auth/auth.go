// Package auth provides authentication for minimaldoc-server.
// Supports local email/password and OAuth 2.0/OIDC providers.
package auth

import (
	"context"
	"fmt"
)

// Provider defines the interface for authentication providers.
type Provider interface {
	// Name returns the provider identifier.
	Name() string

	// Authenticate validates credentials and returns user info.
	Authenticate(ctx context.Context, credentials map[string]string) (*UserInfo, error)

	// GetAuthURL returns the URL to redirect for OAuth login (OAuth providers only).
	GetAuthURL(state string) string

	// HandleCallback processes the OAuth callback and returns user info.
	HandleCallback(ctx context.Context, code string) (*UserInfo, error)
}

// UserInfo contains authenticated user information.
type UserInfo struct {
	ID            string
	Email         string
	Name          string
	AvatarURL     string
	EmailVerified bool
	Provider      string // "local", "google", "github", "cognito", etc.
	ProviderID    string // ID from OAuth provider
	Role          string // admin, editor, viewer
}

// ProviderRegistry manages authentication providers.
type ProviderRegistry struct {
	providers map[string]Provider
}

// NewRegistry creates a new provider registry.
func NewRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry.
func (r *ProviderRegistry) Register(p Provider) {
	r.providers[p.Name()] = p
}

// Get retrieves a provider by name.
func (r *ProviderRegistry) Get(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
	return p, nil
}

// List returns all registered provider names.
func (r *ProviderRegistry) List() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}
