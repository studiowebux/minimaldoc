package store

import (
	"context"
	"database/sql"
)

// Site queries

func (db *DB) CreateSite(ctx context.Context, id, name, domain, apiKeyHash string) (*Site, error) {
	query := `
		INSERT INTO sites (id, name, domain, api_key, config)
		VALUES ($1, $2, $3, $4, '{}')
	`
	_, err := db.ExecContext(ctx, query, id, name, nullString(domain), apiKeyHash)
	if err != nil {
		return nil, err
	}

	// Fetch the created site
	return db.GetSiteByID(ctx, id)
}

func (db *DB) GetSiteByID(ctx context.Context, id string) (*Site, error) {
	query := `SELECT id, name, domain, api_key, config, created_at, updated_at FROM sites WHERE id = $1`
	var s Site
	err := db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.Name, &s.Domain, &s.APIKey, &s.Config, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (db *DB) GetSiteByAPIKey(ctx context.Context, apiKeyHash string) (*Site, error) {
	query := `SELECT id, name, domain, api_key, config, created_at, updated_at FROM sites WHERE api_key = $1`
	var s Site
	err := db.QueryRowContext(ctx, query, apiKeyHash).Scan(
		&s.ID, &s.Name, &s.Domain, &s.APIKey, &s.Config, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (db *DB) ListSites(ctx context.Context) ([]Site, error) {
	query := `SELECT id, name, domain, api_key, config, created_at, updated_at FROM sites ORDER BY created_at DESC`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []Site
	for rows.Next() {
		var s Site
		if err := rows.Scan(&s.ID, &s.Name, &s.Domain, &s.APIKey, &s.Config, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sites = append(sites, s)
	}
	return sites, rows.Err()
}

func (db *DB) UpdateSite(ctx context.Context, id, name, domain string) error {
	query := `UPDATE sites SET name = $1, domain = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	_, err := db.ExecContext(ctx, query, name, nullString(domain), id)
	return err
}

func (db *DB) UpdateSiteAPIKey(ctx context.Context, siteID, newAPIKeyHash string) error {
	query := `UPDATE sites SET api_key = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := db.ExecContext(ctx, query, newAPIKeyHash, siteID)
	return err
}

func (db *DB) DeleteSite(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM sites WHERE id = $1`, id)
	return err
}
