// Package config provides configuration management for minimaldoc-server.
// Configuration is loaded from environment variables and optional config file.
package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// DocsConfig holds configuration loaded from the static site's config.yaml.
type DocsConfig struct {
	Title   string `yaml:"title"`
	BaseURL string `yaml:"base_url"`
	SiteID  string `yaml:"site_id"`
}

// LoadDocsConfig loads the docs config.yaml from the specified path.
func LoadDocsConfig(path string) (*DocsConfig, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- path from trusted user configuration
	if err != nil {
		return nil, err
	}
	var cfg DocsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

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

	// Storage settings
	Storage StorageConfig

	// Rate limiting settings
	RateLimit RateLimitConfig

	// OpenTelemetry settings
	Telemetry TelemetryConfig

	// Forum settings
	Forum ForumConfig

	// Docs config (loaded from docs/config.yaml)
	Docs *DocsConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host           string
	Port           int // Public API port (tracking, feedback, newsletter)
	AdminPort      int // Admin UI and management API port
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	AdminPath      string // Path prefix for admin UI (default: /admin)
	APIPath        string // Path prefix for API (default: /api)
	CORSOrigins    []string
	TrustedProxies []string // IP addresses of trusted reverse proxies (for X-Forwarded-For)
	DocsDir        string   // Directory containing static docs (default: public)
	DocsConfigPath string // Path to docs config.yaml (default: docs/config.yaml)
	Environment    string // "production" or "development" (default: development)
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
	JWTSecret            string
	JWTExpiry            time.Duration
	RefreshExpiry        time.Duration
	BCryptCost           int
	SessionCookieKey     string
	EnableLocal          bool   // Enable email/password auth
	EnableOAuth          bool   // Enable OAuth providers
	AllowNewsletterOAuth bool   // Allow OAuth subscribe for newsletter
	BootstrapToken       string // Optional token required to call /api/bootstrap
	SecureCookies        bool   // Set Secure flag on cookies (requires HTTPS)
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
	Issuer      string // OIDC issuer URL
	AuthURL     string // Authorization endpoint (if not using discovery)
	TokenURL    string // Token endpoint (if not using discovery)
	UserInfoURL string // UserInfo endpoint (if not using discovery)
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
	BaseURL     string // Base URL for verification links (e.g., https://api.example.com)
}

// AIConfig holds Claude API settings.
type AIConfig struct {
	Enabled   bool
	APIKey    string
	Model     string // claude-3-sonnet, claude-3-opus, etc.
	MaxTokens int
}

// StorageConfig holds file storage settings.
type StorageConfig struct {
	Provider     string   // local or s3
	LocalPath    string   // Local filesystem path for uploads
	S3Bucket     string   // S3 bucket name
	S3Region     string   // S3 region
	S3AccessKey  string   // S3 access key
	S3SecretKey  string   // S3 secret key
	S3Endpoint   string   // Custom S3 endpoint (for MinIO, R2)
	S3PublicURL  string   // Public URL prefix for S3 objects
	MaxFileSize  int64    // Max upload size in bytes
	AllowedTypes []string // Allowed MIME types
}

// RateLimitConfig holds rate limiting settings.
type RateLimitConfig struct {
	Enabled      bool          // Enable rate limiting
	LoginLimit   int           // Max login attempts per window
	LoginWindow  time.Duration // Login rate limit window
	APILimit     int           // Max API requests per window
	APIWindow    time.Duration // API rate limit window
	SubmitLimit  int           // Max submissions per window (comments, feedback, newsletter)
	SubmitWindow time.Duration // Submit rate limit window
}

// TelemetryConfig holds OpenTelemetry settings.
type TelemetryConfig struct {
	Enabled     bool   // Enable OpenTelemetry tracing
	Endpoint    string // OTLP HTTP endpoint (e.g. "localhost:4318")
	ServiceName string // Service name for traces (default: minimaldoc-server)
}

// ForumConfig holds forum feature settings.
type ForumConfig struct {
	Enabled           bool          // Enable forum feature
	AllowAnonymous    bool          // Allow viewing without auth
	RequireAuth       bool          // Require auth to post
	MaxTopicsPerDay   int           // Rate limit for topic creation
	MaxPostsPerDay    int           // Rate limit for post creation
	EditWindow        time.Duration // Time window for editing posts
	ModerationMode    string        // none, first_post, all
	EmailEnabled      bool          // Enable email notifications
	EmailDigest       string        // daily, weekly, none
	EmailOnReply      bool          // Send email on reply
	EmailOnMention    bool          // Send email on @mention
	ReputationEnabled bool          // Enable reputation system
	RepTopicCreate    int           // Points for creating topic
	RepPostCreate     int           // Points for posting reply
	RepLikeReceived   int           // Points when liked
	RepSolutionMarked int           // Points for accepted solution
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host:           getEnv("SERVER_HOST", "0.0.0.0"),
			Port:           getEnvInt("SERVER_PORT", 8080),
			AdminPort:      getEnvInt("SERVER_ADMIN_PORT", 8090),
			ReadTimeout:    getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout:   getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			AdminPath:      getEnv("SERVER_ADMIN_PATH", "/admin"),
			APIPath:        getEnv("SERVER_API_PATH", "/api/v1"),
			CORSOrigins:    getEnvSlice("SERVER_CORS_ORIGINS", nil),
			TrustedProxies: getEnvSlice("SERVER_TRUSTED_PROXIES", nil),
			DocsDir:        getEnv("SERVER_DOCS_DIR", "public"),
			DocsConfigPath: getEnv("DOCS_CONFIG_PATH", "docs/config.yaml"),
			Environment:    getEnv("SERVER_ENV", "development"),
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
			JWTSecret:            getEnv("AUTH_JWT_SECRET", ""),
			JWTExpiry:            getEnvDuration("AUTH_JWT_EXPIRY", 15*time.Minute),
			RefreshExpiry:        getEnvDuration("AUTH_REFRESH_EXPIRY", 7*24*time.Hour),
			BCryptCost:           getEnvInt("AUTH_BCRYPT_COST", 12),
			SessionCookieKey:     getEnv("AUTH_SESSION_COOKIE", "minimaldoc_session"),
			EnableLocal:          getEnvBool("AUTH_ENABLE_LOCAL", true),
			EnableOAuth:          getEnvBool("AUTH_ENABLE_OAUTH", false),
			AllowNewsletterOAuth: getEnvBool("NEWSLETTER_ALLOW_OAUTH_SUBSCRIBE", false),
			BootstrapToken:       getEnv("BOOTSTRAP_TOKEN", ""),
			SecureCookies:        getEnvBool("AUTH_SECURE_COOKIES", false),
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
			BaseURL:     getEnv("EMAIL_BASE_URL", "http://localhost:8080"),
		},
		AI: AIConfig{
			Enabled:   getEnvBool("AI_ENABLED", false),
			APIKey:    getEnv("ANTHROPIC_API_KEY", ""),
			Model:     getEnv("AI_MODEL", "claude-sonnet-4-20250514"),
			MaxTokens: getEnvInt("AI_MAX_TOKENS", 4096),
		},
		Storage: StorageConfig{
			Provider:     getEnv("STORAGE_PROVIDER", "local"),
			LocalPath:    getEnv("STORAGE_LOCAL_PATH", "./uploads"),
			S3Bucket:     getEnv("STORAGE_S3_BUCKET", ""),
			S3Region:     getEnv("STORAGE_S3_REGION", "us-east-1"),
			S3AccessKey:  getEnv("STORAGE_S3_ACCESS_KEY", ""),
			S3SecretKey:  getEnv("STORAGE_S3_SECRET_KEY", ""),
			S3Endpoint:   getEnv("STORAGE_S3_ENDPOINT", ""),
			S3PublicURL:  getEnv("STORAGE_S3_PUBLIC_URL", ""),
			MaxFileSize:  getEnvInt64("STORAGE_MAX_FILE_SIZE", 5*1024*1024), // 5MB
			AllowedTypes: getEnvSlice("STORAGE_ALLOWED_TYPES", []string{"image/jpeg", "image/png", "image/gif", "image/webp"}),
		},
		RateLimit: RateLimitConfig{
			Enabled:      getEnvBool("RATE_LIMIT_ENABLED", true),
			LoginLimit:   getEnvInt("RATE_LIMIT_LOGIN_LIMIT", 5),
			LoginWindow:  getEnvDuration("RATE_LIMIT_LOGIN_WINDOW", 15*time.Minute),
			APILimit:     getEnvInt("RATE_LIMIT_API_LIMIT", 100),
			APIWindow:    getEnvDuration("RATE_LIMIT_API_WINDOW", time.Minute),
			SubmitLimit:  getEnvInt("RATE_LIMIT_SUBMIT_LIMIT", 10),
			SubmitWindow: getEnvDuration("RATE_LIMIT_SUBMIT_WINDOW", time.Minute),
		},
		Telemetry: TelemetryConfig{
			Enabled:     getEnvBool("OTEL_ENABLED", false),
			Endpoint:    getEnv("OTEL_ENDPOINT", ""),
			ServiceName: getEnv("OTEL_SERVICE_NAME", "minimaldoc-server"),
		},
		Forum: ForumConfig{
			Enabled:           getEnvBool("FORUM_ENABLED", true),
			AllowAnonymous:    getEnvBool("FORUM_ALLOW_ANONYMOUS", true),
			RequireAuth:       getEnvBool("FORUM_REQUIRE_AUTH", true),
			MaxTopicsPerDay:   getEnvInt("FORUM_MAX_TOPICS_PER_DAY", 10),
			MaxPostsPerDay:    getEnvInt("FORUM_MAX_POSTS_PER_DAY", 50),
			EditWindow:        getEnvDuration("FORUM_EDIT_WINDOW", time.Hour),
			ModerationMode:    getEnv("FORUM_MODERATION_MODE", "none"),
			EmailEnabled:      getEnvBool("FORUM_EMAIL_ENABLED", true),
			EmailDigest:       getEnv("FORUM_EMAIL_DIGEST", "daily"),
			EmailOnReply:      getEnvBool("FORUM_EMAIL_ON_REPLY", true),
			EmailOnMention:    getEnvBool("FORUM_EMAIL_ON_MENTION", true),
			ReputationEnabled: getEnvBool("FORUM_REPUTATION_ENABLED", true),
			RepTopicCreate:    getEnvInt("FORUM_REP_TOPIC_CREATE", 5),
			RepPostCreate:     getEnvInt("FORUM_REP_POST_CREATE", 2),
			RepLikeReceived:   getEnvInt("FORUM_REP_LIKE_RECEIVED", 1),
			RepSolutionMarked: getEnvInt("FORUM_REP_SOLUTION_MARKED", 10),
		},
	}

	// Load docs config (optional - don't fail if not found)
	if docsConfig, err := LoadDocsConfig(cfg.Server.DocsConfigPath); err == nil {
		cfg.Docs = docsConfig
	}

	// Validate required settings
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that required configuration is present and logs
// warnings for production misconfigurations.
func (c *Config) Validate() error {
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("AUTH_JWT_SECRET is required")
	}
	if len(c.Auth.JWTSecret) < 32 {
		return fmt.Errorf("AUTH_JWT_SECRET must be at least 32 characters")
	}
	if isWeakSecret(c.Auth.JWTSecret) {
		if c.Server.Environment == "production" {
			return fmt.Errorf("AUTH_JWT_SECRET has suspiciously low entropy — use a random secret")
		}
		slog.Warn("AUTH_JWT_SECRET has suspiciously low entropy — use a random secret in production")
	}
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	// Validate email base URL scheme
	if c.Email.BaseURL != "" {
		u, err := url.Parse(c.Email.BaseURL)
		if err != nil {
			return fmt.Errorf("EMAIL_BASE_URL is not a valid URL: %w", err)
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return fmt.Errorf("EMAIL_BASE_URL must use http or https scheme (got %q)", u.Scheme)
		}
	}

	c.warnProductionConfig()
	return nil
}

// isWeakSecret checks for placeholder or low-entropy JWT secrets.
func isWeakSecret(s string) bool {
	lower := strings.ToLower(s)
	for _, pattern := range []string{"change", "secret", "example", "placeholder", "default", "password"} {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	// Check all-same-character
	if len(s) > 0 {
		allSame := true
		for i := 1; i < len(s); i++ {
			if s[i] != s[0] {
				allSame = false
				break
			}
		}
		if allSame {
			return true
		}
	}
	return false
}

// warnProductionConfig logs warnings for settings that are unsafe in production.
func (c *Config) warnProductionConfig() {
	isProd := c.Server.Environment == "production"

	// CORS wildcard is unsafe in production
	for _, origin := range c.Server.CORSOrigins {
		if origin == "*" {
			if isProd {
				slog.Warn("CORS wildcard origin is unsafe in production, set SERVER_CORS_ORIGINS to specific domains")
			}
			break
		}
	}

	// Insecure cookies in production
	if isProd && !c.Auth.SecureCookies {
		slog.Warn("AUTH_SECURE_COOKIES is false in production, cookies will be sent over HTTP")
	}

	// Email base URL pointing to localhost or HTTP in production
	if isProd {
		if strings.Contains(c.Email.BaseURL, "localhost") || strings.Contains(c.Email.BaseURL, "127.0.0.1") {
			slog.Warn("EMAIL_BASE_URL contains localhost in production", "url", c.Email.BaseURL)
		}
		if strings.HasPrefix(c.Email.BaseURL, "http://") {
			slog.Warn("EMAIL_BASE_URL uses HTTP in production, use HTTPS", "url", c.Email.BaseURL)
		}
	}

	// SMTP credentials missing when provider is smtp
	if c.Email.Provider == "smtp" && (c.Email.SMTPUser == "" || c.Email.SMTPPass == "") {
		slog.Warn("EMAIL_PROVIDER is smtp but SMTP_USER or SMTP_PASS is empty")
	}

	// S3 credentials missing when storage provider is s3
	if c.Storage.Provider == "s3" && (c.Storage.S3AccessKey == "" || c.Storage.S3SecretKey == "") {
		slog.Warn("STORAGE_PROVIDER is s3 but S3 access key or secret key is empty")
	}

	// OAuth enabled but no providers configured
	if c.Auth.EnableOAuth && len(c.OAuth.Providers) == 0 {
		slog.Warn("AUTH_ENABLE_OAUTH is true but no OAuth providers are configured")
	}

	// AI enabled without API key
	if c.AI.Enabled && c.AI.APIKey == "" {
		slog.Warn("AI_ENABLED is true but ANTHROPIC_API_KEY is empty")
	}

	// Mock email provider in production
	if isProd && c.Email.Provider == "mock" {
		slog.Warn("EMAIL_PROVIDER is mock in production, emails will not be sent")
	}

	// SQLite in production
	if isProd && c.Database.Driver == "sqlite" {
		slog.Warn("DB_DRIVER is sqlite in production, consider using PostgreSQL")
	}
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

func getEnvInt64(key string, defaultVal int64) int64 {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
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
		parts := strings.Split(val, ",")
		result := make([]string, 0, len(parts))
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	return defaultVal
}
