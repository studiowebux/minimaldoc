package store

import (
	"context"
	"database/sql"
)

// Doc access queries

func (db *DB) CreateDocAccess(ctx context.Context, id, siteID, pathPattern, requiredRole, description string) (*DocAccess, error) {
	query := `
		INSERT INTO doc_access (id, site_id, path_pattern, required_role, description)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := db.ExecContext(ctx, query, id, siteID, pathPattern, requiredRole, nullString(description))
	if err != nil {
		return nil, err
	}
	return db.GetDocAccessByID(ctx, id)
}

func (db *DB) GetDocAccessByID(ctx context.Context, id string) (*DocAccess, error) {
	query := `
		SELECT id, site_id, path_pattern, required_role, description, created_at, updated_at
		FROM doc_access WHERE id = $1
	`
	var d DocAccess
	err := db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.SiteID, &d.PathPattern, &d.RequiredRole, &d.Description, &d.CreatedAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &d, err
}

func (db *DB) ListDocAccess(ctx context.Context, siteID string) ([]DocAccess, error) {
	query := `
		SELECT id, site_id, path_pattern, required_role, description, created_at, updated_at
		FROM doc_access WHERE site_id = $1
		ORDER BY path_pattern ASC
	`
	rows, err := db.QueryContext(ctx, query, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []DocAccess
	for rows.Next() {
		var d DocAccess
		if err := rows.Scan(&d.ID, &d.SiteID, &d.PathPattern, &d.RequiredRole, &d.Description, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, d)
	}
	return rules, rows.Err()
}

func (db *DB) UpdateDocAccess(ctx context.Context, id, pathPattern, requiredRole, description string) error {
	query := `
		UPDATE doc_access
		SET path_pattern = $1, required_role = $2, description = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`
	_, err := db.ExecContext(ctx, query, pathPattern, requiredRole, nullString(description), id)
	return err
}

func (db *DB) DeleteDocAccess(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM doc_access WHERE id = $1`, id)
	return err
}

// CheckDocAccess checks if a path matches any access rule and returns the required role.
// Returns nil if the path is not protected.
func (db *DB) CheckDocAccess(ctx context.Context, siteID, path string) (*DocAccess, error) {
	// Get all rules for the site and check pattern matching in Go
	// This allows for glob-style patterns
	rules, err := db.ListDocAccess(ctx, siteID)
	if err != nil {
		return nil, err
	}

	for _, rule := range rules {
		if matchPath(rule.PathPattern, path) {
			return &rule, nil
		}
	}
	return nil, nil
}

// matchPath checks if a path matches a pattern (supports * and ** wildcards).
func matchPath(pattern, path string) bool {
	// Simple glob matching:
	// * matches any sequence of non-/ characters
	// ** matches any sequence of characters including /
	// Exact match if no wildcards

	if pattern == path {
		return true
	}

	// Convert pattern to regex-like matching
	i, j := 0, 0
	for i < len(pattern) && j < len(path) {
		if pattern[i] == '*' {
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				// ** matches everything including /
				i += 2
				if i >= len(pattern) {
					return true // ** at end matches everything
				}
				// Find next literal char in pattern
				for j < len(path) {
					if matchPath(pattern[i:], path[j:]) {
						return true
					}
					j++
				}
				return false
			}
			// * matches non-/ characters
			i++
			if i >= len(pattern) {
				// * at end - check no more slashes
				for ; j < len(path); j++ {
					if path[j] == '/' {
						return false
					}
				}
				return true
			}
			// Find next literal char
			for j < len(path) && path[j] != '/' {
				if matchPath(pattern[i:], path[j:]) {
					return true
				}
				j++
			}
			continue
		}
		if pattern[i] != path[j] {
			return false
		}
		i++
		j++
	}

	// Handle trailing wildcards
	for i < len(pattern) {
		if pattern[i] == '*' {
			i++
			if i < len(pattern) && pattern[i] == '*' {
				i++
			}
		} else {
			break
		}
	}

	return i >= len(pattern) && j >= len(path)
}
