// Package store provides database access for minimaldoc-server.
// Supports PostgreSQL and SQLite backends.
package store

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"

	"github.com/studiowebux/minimaldoc/internal/server/config"
)

//go:embed migrations/postgres/*.sql
var postgresMigrations embed.FS

//go:embed migrations/sqlite/*.sql
var sqliteMigrations embed.FS

// DB wraps the database connection and provides query methods.
type DB struct {
	*sql.DB
	driver string
}

// New creates a new database connection.
func New(cfg config.DatabaseConfig) (*DB, error) {
	var (
		db  *sql.DB
		err error
	)

	driver := strings.ToLower(cfg.Driver)

	switch driver {
	case "postgres", "postgresql":
		db, err = sql.Open("postgres", cfg.URL)
		driver = "postgres"
	case "sqlite", "sqlite3":
		db, err = sql.Open("sqlite", cfg.URL)
		driver = "sqlite"
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{DB: db, driver: driver}, nil
}

// Migrate runs database migrations using embedded SQL files.
func (db *DB) Migrate() error {
	// Create migrations table if not exists
	if err := db.createMigrationsTable(); err != nil {
		return err
	}

	// Get current version
	currentVersion, err := db.getMigrationVersion()
	if err != nil {
		return err
	}

	// Run pending migrations
	migrations, err := db.loadMigrations()
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if m.version > currentVersion {
			log.Printf("Running migration %d: %s", m.version, m.name)
			if err := db.runMigration(m); err != nil {
				return fmt.Errorf("migration %d failed: %w", m.version, err)
			}
		}
	}

	return nil
}

type migration struct {
	version int
	name    string
	sql     string
}

func (db *DB) createMigrationsTable() error {
	query := `CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}

func (db *DB) getMigrationVersion() (int, error) {
	var version int
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func (db *DB) loadMigrations() ([]migration, error) {
	var fs embed.FS
	var dir string

	switch db.driver {
	case "postgres":
		fs = postgresMigrations
		dir = "migrations/postgres"
	case "sqlite":
		fs = sqliteMigrations
		dir = "migrations/sqlite"
	}

	entries, err := fs.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrations []migration
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".up.sql") {
			continue
		}

		// Parse version from filename (e.g., "001_initial.up.sql")
		var version int
		var name string
		_, err := fmt.Sscanf(e.Name(), "%d_%s", &version, &name)
		if err != nil {
			continue
		}

		content, err := fs.ReadFile(dir + "/" + e.Name())
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migration{
			version: version,
			name:    strings.TrimSuffix(name, ".up.sql"),
			sql:     string(content),
		})
	}

	return migrations, nil
}

func (db *DB) runMigration(m migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Run migration SQL
	if _, err := tx.Exec(m.sql); err != nil {
		return err
	}

	// Record migration
	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", m.version); err != nil {
		return err
	}

	return tx.Commit()
}

// Driver returns the database driver name.
func (db *DB) Driver() string {
	return db.driver
}

// IsPostgres returns true if using PostgreSQL.
func (db *DB) IsPostgres() bool {
	return db.driver == "postgres"
}

// IsSQLite returns true if using SQLite.
func (db *DB) IsSQLite() bool {
	return db.driver == "sqlite"
}
