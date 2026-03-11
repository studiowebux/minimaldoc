package store

import (
	"context"
	"database/sql"
)

// CreateUpload stores upload metadata.
func (db *DB) CreateUpload(ctx context.Context, id, siteID, userID, filename, mimeType string, sizeBytes int64, storagePath, url string) (*Upload, error) {
	query := `
		INSERT INTO uploads (id, site_id, user_id, filename, mime_type, size_bytes, storage_path, url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := db.ExecContext(ctx, query, id, siteID, userID, filename, mimeType, sizeBytes, storagePath, url)
	if err != nil {
		return nil, err
	}
	return db.GetUpload(ctx, id)
}

// GetUpload retrieves an upload by ID.
func (db *DB) GetUpload(ctx context.Context, id string) (*Upload, error) {
	query := `SELECT id, site_id, user_id, filename, mime_type, size_bytes, storage_path, url, created_at FROM uploads WHERE id = $1`
	var u Upload
	err := db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.SiteID, &u.UserID, &u.Filename, &u.MimeType, &u.SizeBytes, &u.StoragePath, &u.URL, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

// DeleteUpload removes an upload record.
func (db *DB) DeleteUpload(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM uploads WHERE id = $1`, id)
	return err
}

// ListUploads returns uploads for a site, optionally filtered by user.
func (db *DB) ListUploads(ctx context.Context, siteID, userID string, limit, offset int) ([]Upload, error) {
	var query string
	var args []interface{}

	if userID != "" {
		query = `SELECT id, site_id, user_id, filename, mime_type, size_bytes, storage_path, url, created_at
				 FROM uploads WHERE site_id = $1 AND user_id = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`
		args = []interface{}{siteID, userID, limit, offset}
	} else {
		query = `SELECT id, site_id, user_id, filename, mime_type, size_bytes, storage_path, url, created_at
				 FROM uploads WHERE site_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{siteID, limit, offset}
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uploads []Upload
	for rows.Next() {
		var u Upload
		if err := rows.Scan(&u.ID, &u.SiteID, &u.UserID, &u.Filename, &u.MimeType, &u.SizeBytes, &u.StoragePath, &u.URL, &u.CreatedAt); err != nil {
			return nil, err
		}
		uploads = append(uploads, u)
	}
	return uploads, rows.Err()
}
