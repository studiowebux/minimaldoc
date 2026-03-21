package store

import (
	"context"
	"fmt"
	"strings"
)

// Audit log queries - Repository pattern for audit log operations

// CreateAuditLog creates a new audit log entry.
func (db *DB) CreateAuditLog(ctx context.Context, id, siteID, userID, userEmail, action, entityType, entityID, entityName, details, ipAddress, userAgent string) error {
	query := `
		INSERT INTO audit_logs (id, site_id, user_id, user_email, action, entity_type, entity_id, entity_name, details, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := db.ExecContext(ctx, query, id, siteID, nullString(userID), userEmail, action, entityType, nullString(entityID), nullString(entityName), nullString(details), nullString(ipAddress), nullString(userAgent))
	return err
}

// ListAuditLogs returns audit logs with optional filtering and pagination.
func (db *DB) ListAuditLogs(ctx context.Context, siteID, action, entityType string, limit, offset int) ([]AuditLog, int, error) {
	// Build WHERE clause dynamically
	var conditions []string
	var args []any
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("site_id = $%d", argIndex))
	args = append(args, siteID)
	argIndex++

	if action != "" {
		conditions = append(conditions, fmt.Sprintf("action = $%d", argIndex))
		args = append(args, action)
		argIndex++
	}

	if entityType != "" {
		conditions = append(conditions, fmt.Sprintf("entity_type = $%d", argIndex))
		args = append(args, entityType)
		argIndex++
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_logs %s", whereClause)
	var total int
	if err := db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch logs with pagination
	query := fmt.Sprintf(`
		SELECT id, site_id, user_id, user_email, action, entity_type, entity_id, entity_name, details, ip_address, user_agent, created_at
		FROM audit_logs
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var l AuditLog
		if err := rows.Scan(&l.ID, &l.SiteID, &l.UserID, &l.UserEmail, &l.Action, &l.EntityType, &l.EntityID, &l.EntityName, &l.Details, &l.IPAddress, &l.UserAgent, &l.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}

	return logs, total, rows.Err()
}

// AuditLogStats holds aggregate counts for audit logs.
type AuditLogStats struct {
	Total    int
	Today    int
	ThisWeek int
}

// GetAuditLogStats returns aggregate counts for audit logs.
func (db *DB) GetAuditLogStats(ctx context.Context, siteID string) (*AuditLogStats, error) {
	stats := &AuditLogStats{}

	// Total count
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_logs WHERE site_id = $1`, siteID).Scan(&stats.Total); err != nil {
		return nil, err
	}

	// Today's count - use date() for SQLite compatibility
	todayQuery := `SELECT COUNT(*) FROM audit_logs WHERE site_id = $1 AND date(created_at) = date('now')`
	if db.Driver() == "postgres" {
		todayQuery = `SELECT COUNT(*) FROM audit_logs WHERE site_id = $1 AND created_at >= CURRENT_DATE`
	}
	if err := db.QueryRowContext(ctx, todayQuery, siteID).Scan(&stats.Today); err != nil {
		return nil, err
	}

	// This week's count
	weekQuery := `SELECT COUNT(*) FROM audit_logs WHERE site_id = $1 AND date(created_at) >= date('now', '-7 days')`
	if db.Driver() == "postgres" {
		weekQuery = `SELECT COUNT(*) FROM audit_logs WHERE site_id = $1 AND created_at >= CURRENT_DATE - INTERVAL '7 days'`
	}
	if err := db.QueryRowContext(ctx, weekQuery, siteID).Scan(&stats.ThisWeek); err != nil {
		return nil, err
	}

	return stats, nil
}

// GetAuditLogEntityTypes returns distinct entity types in the audit log for filtering.
func (db *DB) GetAuditLogEntityTypes(ctx context.Context, siteID string) ([]string, error) {
	query := `SELECT DISTINCT entity_type FROM audit_logs WHERE site_id = $1 ORDER BY entity_type`
	rows, err := db.QueryContext(ctx, query, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, rows.Err()
}
