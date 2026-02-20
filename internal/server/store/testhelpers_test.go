package store

import (
	"context"
	"testing"

	"github.com/studiowebux/minimaldoc/internal/server/config"
)

func testContext() context.Context {
	return context.Background()
}

// setupTestDB creates an in-memory SQLite database with all migrations applied.
// Registers t.Cleanup to close the database when the test completes.
func setupTestDB(t *testing.T) *DB {
	t.Helper()

	cfg := config.DatabaseConfig{
		Driver:          "sqlite",
		URL:             ":memory:",
		MaxOpenConns:    1,
		MaxIdleConns:    1,
	}

	db, err := New(cfg)
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

// createTestSite creates a site for use in tests and returns its ID.
func createTestSite(t *testing.T, db *DB) string {
	t.Helper()

	ctx := testContext()
	site, err := db.CreateSite(ctx, "test-site-id", "Test Site", "test.example.com", "hashed-api-key")
	if err != nil {
		t.Fatalf("failed to create test site: %v", err)
	}
	return site.ID
}

// createTestUser creates a user for use in tests and returns its ID.
func createTestUser(t *testing.T, db *DB, siteID string) string {
	t.Helper()

	ctx := testContext()
	user, err := db.CreateUser(ctx, "test-user-id", siteID, "test@example.com", "$2a$12$fakehashfakehashfakehashfakehashfakehashfakehashfake", "admin", "Test User")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return user.ID
}
