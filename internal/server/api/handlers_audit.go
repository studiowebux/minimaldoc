package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Audit log handlers - Admin-only access for viewing audit trail

// adminAuditLog renders the audit log admin page.
func (r *Router) adminAuditLog(c *gin.Context) {
	claims, _ := getUserClaims(c)
	c.HTML(http.StatusOK, "audit-log.html", gin.H{
		"Title":       "Audit Log",
		"CurrentPage": "audit-log",
		"User":        claims,
	})
}

// listAuditLogsAPI returns audit logs as JSON for API consumers.
func (r *Router) listAuditLogsAPI(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED"})
		return
	}

	// Parse filters
	action := c.Query("action")
	entityType := c.Query("entity_type")

	// Parse pagination
	limit := 50
	offset := 0
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	logs, total, err := r.db.ListAuditLogs(c.Request.Context(), siteID, action, entityType, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DATABASE_ERROR"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// logAuditAction creates an audit log entry for an admin action.
// This is a helper method called by other handlers after performing actions.
func (r *Router) logAuditAction(c *gin.Context, action, entityType, entityID, entityName, details string) {
	siteID, err := getSiteID(c)
	if err != nil {
		return // Silently fail if no site context
	}

	userID, _ := getUserID(c)
	userEmail := ""
	if claims, err := getUserClaims(c); err == nil && claims != nil {
		userEmail = claims.Email
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	id := uuid.New().String()

	// Log asynchronously to avoid slowing down the main request
	go func() {
		_ = r.db.CreateAuditLog(c.Request.Context(), id, siteID, userID, userEmail, action, entityType, entityID, entityName, details, ipAddress, userAgent)
	}()
}
