package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/config"
	"github.com/studiowebux/minimaldoc/internal/server/email"
	"github.com/studiowebux/minimaldoc/internal/server/store"
)

func testContext() context.Context {
	return context.Background()
}

func init() {
	gin.SetMode(gin.TestMode)
}

// testConfig returns a minimal config for testing.
func testConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Host:        "localhost",
			Port:        8080,
			AdminPort:   8090,
			AdminPath:   "/admin",
			APIPath:     "/api",
			CORSOrigins: []string{"*"},
			Environment: "development",
		},
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret-key-min-32-characters-long",
			JWTExpiry:         15 * time.Minute,
			RefreshExpiry:     7 * 24 * time.Hour,
			BCryptCost:       4, // Low cost for fast tests
			SessionCookieKey: "test_session",
			EnableLocal:      true,
		},
		Email: config.EmailConfig{
			Provider:    "mock",
			BaseURL:     "http://localhost:8080",
			FromAddress: "test@example.com",
			FromName:    "Test",
		},
		Storage: config.StorageConfig{
			Provider:  "local",
			LocalPath: "/tmp/test-uploads",
		},
		RateLimit: config.RateLimitConfig{
			Enabled: false,
		},
		Forum: config.ForumConfig{
			Enabled:           true,
			AllowAnonymous:    true,
			RequireAuth:       true,
			ReputationEnabled: true,
		},
	}
}

// setupTestDB creates an in-memory SQLite database with all migrations applied.
func setupTestDB(t *testing.T) *store.DB {
	t.Helper()

	cfg := config.DatabaseConfig{
		Driver:       "sqlite",
		URL:          ":memory:",
		MaxOpenConns: 1,
		MaxIdleConns: 1,
	}

	db, err := store.New(cfg)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	if err := db.Migrate(); err != nil {
		db.Close()
		t.Fatalf("failed to run migrations: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

// performRequest executes an HTTP request against a Gin engine.
func performRequest(engine *gin.Engine, method, path string, body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)
	return w
}

// performRequestWithAuth executes an HTTP request with a Bearer token.
func performRequestWithAuth(engine *gin.Engine, method, path, token string, body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	engine.ServeHTTP(w, req)
	return w
}

// setupTestRouter creates a Router with an in-memory SQLite DB and mock email sender.
// Returns the Router, DB, and mock email sender.
func setupTestRouter(t *testing.T) (*Router, *store.DB, *email.MockSender) {
	t.Helper()

	db := setupTestDB(t)
	cfg := testConfig()
	mockEmail := email.NewMockSender()

	// Build a minimal public router without template loading
	engine := gin.New()
	r := &Router{
		Engine:  engine,
		config:  cfg,
		db:      db,
		email:   mockEmail,
	}

	// Register health endpoints for testing
	r.GET("/healthz", r.liveness)
	r.GET("/readyz", r.readiness)

	// Register API routes
	api := r.Group(cfg.Server.APIPath)
	{
		api.GET("/health", r.healthCheck)
		api.POST("/bootstrap", r.bootstrap)

		auth := api.Group("/auth")
		{
			auth.POST("/login", r.login)
			auth.POST("/logout", r.logout)
			auth.POST("/refresh", r.refreshToken)
			auth.GET("/me", AuthMiddleware(cfg), r.getCurrentUser)
		}

		newsletter := api.Group("/newsletter")
		{
			newsletter.POST("/subscribe", r.subscribe)
			newsletter.GET("/verify", r.verifySubscription)
			newsletter.POST("/unsubscribe", r.unsubscribe)
		}
	}

	return r, db, mockEmail
}
