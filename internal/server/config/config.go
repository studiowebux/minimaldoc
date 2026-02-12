// Package config provides configuration management for minimaldoc-server.
// Configuration is loaded from environment variables and optional config file.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all server configuration.
type Config struct {
	// Server settings
	Server ServerConfig

	// Database settings
	Database DatabaseConfig

	// Authentication settings
	Auth AuthConfig

	// OAuth providers
	OAuth OAuthConfig

	// Email settings
	Email EmailConfig

	// AI settings (Claude)
	AI AIConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	AdminPath    string // Path prefix for admin UI (default: /admin)
	APIPath      string // Path prefix for API (default: /api)
	CORSOrigins  []string
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	Driver          string // "postgres" or "sqlite"
	URL             string // Connection URL
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	MigrationsPath  string
}

// AuthConfig holds authentication settings.
type AuthConfig struct {
	JWTSecret        string
	JWTExpiry        time.Duration
	RefreshExpiry    time.Duration
	BCryptCost       int
	SessionCookieKey string
	EnableLocal      bool // Enable email/password auth
	EnableOAuth      bool // Enable OAuth providers
}

// OAuthConfig holds OAuth 2.0 / OIDC provider settings.
type OAuthConfig struct {
	Providers []OAuthProvider
}

// OAuthProvider represents a single OAuth provider configuration.
type OAuthProvider struct {
	Name         string // cognito, auth0, google, github, oidc
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string

	// OIDC-specific (for generic OIDC or Cognito)
	Issuer       string // OIDC issuer URL
	AuthURL      string // Authorization endpoint (if not using discovery)
	TokenURL     string // Token endpoint (if not using discovery)
	UserInfoURL  string // UserInfo endpoint (if not using discovery)
}

// EmailConfig holds email sending settings.
type EmailConfig struct {
	Provider    string // smtp, mock (ses/sendgrid: add when needed)
	SMTPHost    string
	SMTPPort    int
	SMTPUser    string
	SMTPPass    string
	FromAddress string
	FromName    string
}

// AIConfig holds Claude API settings.
type AIConfig struct {
	Enabled   bool
	APIKey    string
	Model     string // claude-3-sonnet, claude-3-opus, etc.
	MaxTokens int
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			AdminPath:    getEnv("SERVER_ADMIN_PATH", "/admin"),
			APIPath:      getEnv("SERVER_API_PATH", "/api"),
			CORSOrigins:  getEnvSlice("SERVER_CORS_ORIGINS", []string{"*"}),
		},
		Database: DatabaseConfig{
			Driver:          getEnv("DB_DRIVER", "sqlite"),
			URL:             getEnv("DATABASE_URL", "minimaldoc.db"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			MigrationsPath:  getEnv("DB_MIGRATIONS_PATH", "migrations"),
		},
		Auth: AuthConfig{
			JWTSecret:        getEnv("AUTH_JWT_SECRET", ""),
			JWTExpiry:        getEnvDuration("AUTH_JWT_EXPIRY", 15*time.Minute),
			RefreshExpiry:    getEnvDuration("AUTH_REFRESH_EXPIRY", 7*24*time.Hour),
			BCryptCost:       getEnvInt("AUTH_BCRYPT_COST", 12),
			SessionCookieKey: getEnv("AUTH_SESSION_COOKIE", "minimaldoc_session"),
			EnableLocal:      getEnvBool("AUTH_ENABLE_LOCAL", true),
			EnableOAuth:      getEnvBool("AUTH_ENABLE_OAUTH", false),
		},
		OAuth: loadOAuthConfig(),
		Email: EmailConfig{
			Provider:    getEnv("EMAIL_PROVIDER", "mock"),
			SMTPHost:    getEnv("SMTP_HOST", "localhost"),
			SMTPPort:    getEnvInt("SMTP_PORT", 587),
			SMTPUser:    getEnv("SMTP_USER", ""),
			SMTPPass:    getEnv("SMTP_PASS", ""),
			FromAddress: getEnv("EMAIL_FROM_ADDRESS", "noreply@example.com"),
			FromName:    getEnv("EMAIL_FROM_NAME", "MinimalDoc"),
		},
		AI: AIConfig{
			Enabled:   getEnvBool("AI_ENABLED", false),
			APIKey:    getEnv("ANTHROPIC_API_KEY", ""),
			Model:     getEnv("AI_MODEL", "claude-sonnet-4-20250514"),
			MaxTokens: getEnvInt("AI_MAX_TOKENS", 4096),
		},
	}

	// Validate required settings
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that required configuration is present.
func (c *Config) Validate() error {
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("AUTH_JWT_SECRET is required")
	}
	if len(c.Auth.JWTSecret) < 32 {
		return fmt.Errorf("AUTH_JWT_SECRET must be at least 32 characters")
	}
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	return nil
}

// loadOAuthConfig loads OAuth provider configurations from environment.
func loadOAuthConfig() OAuthConfig {
	cfg := OAuthConfig{}

	// Check for each known provider
	providers := []string{"cognito", "auth0", "google", "github"}

	for _, name := range providers {
		prefix := "OAUTH_" + strings.ToUpper(name) + "_"
		clientID := getEnv(prefix+"CLIENT_ID", "")
		if clientID == "" {
			continue
		}

		provider := OAuthProvider{
			Name:         name,
			ClientID:     clientID,
			ClientSecret: getEnv(prefix+"CLIENT_SECRET", ""),
			RedirectURL:  getEnv(prefix+"REDIRECT_URL", ""),
			Scopes:       getEnvSlice(prefix+"SCOPES", defaultScopes(name)),
			Issuer:       getEnv(prefix+"ISSUER", ""),
			AuthURL:      getEnv(prefix+"AUTH_URL", ""),
			TokenURL:     getEnv(prefix+"TOKEN_URL", ""),
			UserInfoURL:  getEnv(prefix+"USERINFO_URL", ""),
		}
		cfg.Providers = append(cfg.Providers, provider)
	}

	return cfg
}

// defaultScopes returns default OAuth scopes for known providers.
func defaultScopes(provider string) []string {
	switch provider {
	case "google":
		return []string{"openid", "email", "profile"}
	case "github":
		return []string{"read:user", "user:email"}
	default:
		return []string{"openid", "email", "profile"}
	}
}

// Helper functions for environment variable parsing

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}

func getEnvSlice(key string, defaultVal []string) []string {
	if val := os.Getenv(key); val != "" {
		return strings.Split(val, ",")
	}
	return defaultVal
}
