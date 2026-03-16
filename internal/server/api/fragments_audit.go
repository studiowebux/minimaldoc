package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Audit log HTMX fragments for admin UI

// fragmentAuditStats returns stats cards HTML.
func (r *Router) fragmentAuditStats(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}

	stats, err := r.db.GetAuditLogStats(c.Request.Context(), siteID)
	if err != nil {
		respondHTMLError(c, "Failed to load stats")
		return
	}

	cards := buildStatCards([]StatCard{
		{Value: stats.Total, Label: "Total Logs"},
		{Value: stats.Today, Label: "Today"},
		{Value: stats.ThisWeek, Label: "This Week"},
	})

	respondHTML(c, cards)
}

// fragmentAuditLogList returns the audit log table HTML with filtering and pagination.
func (r *Router) fragmentAuditLogList(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}

	// Parse filters
	action := c.Query("action")
	entityType := c.Query("entity_type")

	limit, offset := parsePagination(c)

	logs, total, err := r.db.ListAuditLogs(c.Request.Context(), siteID, action, entityType, limit, offset)
	if err != nil {
		respondHTMLError(c, "Failed to load audit logs")
		return
	}

	if len(logs) == 0 {
		respondHTMLEmpty(c, "No audit logs found.")
		return
	}

	// Build table rows
	var rows []TableRow
	for _, log := range logs {
		// Format timestamp
		timestamp := formatAuditTimestamp(log.CreatedAt)

		// Format user info
		userInfo := log.UserEmail
		if userInfo == "" {
			userInfo = "System"
		}

		// Format action badge
		actionBadge := fmt.Sprintf(`<span class="badge badge-%s">%s</span>`, getActionBadgeClass(log.Action), escapeHTML(log.Action))

		// Format entity
		entityInfo := escapeHTML(log.EntityType)
		if log.EntityName.Valid && log.EntityName.String != "" {
			entityInfo += fmt.Sprintf(": %s", escapeHTML(log.EntityName.String))
		}

		// Format details
		details := ""
		if log.Details.Valid && log.Details.String != "" {
			details = escapeHTML(log.Details.String)
			if len(details) > 50 {
				details = details[:50] + "..."
			}
		}

		rows = append(rows, TableRow{
			Cells: []string{
				timestamp,
				escapeHTML(userInfo),
				actionBadge,
				entityInfo,
				details,
			},
		})
	}

	columns := []TableColumn{
		{Header: "Time"},
		{Header: "User"},
		{Header: "Action"},
		{Header: "Entity"},
		{Header: "Details"},
	}

	html := buildDataTable(columns, rows)

	// Add pagination if needed
	if total > limit {
		html += buildAuditPagination(action, entityType, limit, offset, total)
	}

	respondHTML(c, html)
}

// formatAuditTimestamp formats a timestamp for display.
func formatAuditTimestamp(ts string) string {
	// Try parsing as various formats
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, ts); err == nil {
			// Return relative time for recent, absolute for older
			since := time.Since(t)
			if since < time.Minute {
				return "just now"
			}
			if since < time.Hour {
				mins := int(since.Minutes())
				if mins == 1 {
					return "1 minute ago"
				}
				return fmt.Sprintf("%d minutes ago", mins)
			}
			if since < 24*time.Hour {
				hours := int(since.Hours())
				if hours == 1 {
					return "1 hour ago"
				}
				return fmt.Sprintf("%d hours ago", hours)
			}
			// Older entries show date/time
			return t.Format("Jan 2 15:04")
		}
	}
	return ts
}

// getActionBadgeClass returns the CSS class for an action badge.
func getActionBadgeClass(action string) string {
	switch action {
	case "create":
		return "success"
	case "update":
		return "info"
	case "delete":
		return "danger"
	case "publish":
		return "success"
	case "unpublish":
		return "warning"
	case "approve":
		return "success"
	case "reject":
		return "warning"
	case "login":
		return "info"
	case "logout":
		return "secondary"
	default:
		return "secondary"
	}
}

// buildAuditPagination builds pagination controls for the audit log.
func buildAuditPagination(action, entityType string, limit, offset, total int) string {
	var html strings.Builder

	currentPage := offset/limit + 1
	totalPages := (total + limit - 1) / limit

	html.WriteString(`<div class="pagination">`)

	// Build base query string for pagination links
	baseQuery := fmt.Sprintf("limit=%d", limit)
	if action != "" {
		baseQuery += "&action=" + action
	}
	if entityType != "" {
		baseQuery += "&entity_type=" + entityType
	}

	// Previous button
	if currentPage > 1 {
		prevOffset := offset - limit
		html.WriteString(fmt.Sprintf(`<button class="btn btn-secondary" hx-get="/admin/fragments/audit-log-list?%s&offset=%d" hx-target="#audit-log-list" hx-swap="innerHTML">Previous</button>`, baseQuery, prevOffset))
	} else {
		html.WriteString(`<button class="btn btn-secondary" disabled>Previous</button>`)
	}

	// Page info
	html.WriteString(fmt.Sprintf(`<span class="pagination-info">Page %d of %d (%d total)</span>`, currentPage, totalPages, total))

	// Next button
	if currentPage < totalPages {
		nextOffset := offset + limit
		html.WriteString(fmt.Sprintf(`<button class="btn btn-secondary" hx-get="/admin/fragments/audit-log-list?%s&offset=%d" hx-target="#audit-log-list" hx-swap="innerHTML">Next</button>`, baseQuery, nextOffset))
	} else {
		html.WriteString(`<button class="btn btn-secondary" disabled>Next</button>`)
	}

	html.WriteString(`</div>`)

	return html.String()
}
