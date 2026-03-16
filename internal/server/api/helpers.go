package api

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
	"github.com/studiowebux/minimaldoc/internal/server/store"
)

// parsePagination extracts and clamps limit/offset from query params.
// limit: 1–100 (default 20), offset: 0–100000 (default 0).
func parsePagination(c *gin.Context) (limit, offset int) {
	limit = 20
	offset = 0
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
			if limit > 100 {
				limit = 100
			}
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
			if offset > 100_000 {
				offset = 100_000
			}
		}
	}
	return
}

// cspNonce returns the per-request CSP nonce from gin context.
func cspNonce(c *gin.Context) string {
	if nonce, ok := c.Get("csp_nonce"); ok {
		return nonce.(string)
	}
	return ""
}

// Response helpers

// respondHTML sends an HTML response with proper content type.
func respondHTML(c *gin.Context, html string) {
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// respondHTMLError sends an HTML error message.
func respondHTMLError(c *gin.Context, msg string) {
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, fmt.Sprintf(`<p class="error">%s</p>`, template.HTMLEscapeString(msg)))
}

// respondHTMLEmpty sends an HTML empty state message.
func respondHTMLEmpty(c *gin.Context, msg string) {
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, fmt.Sprintf(`<p class="empty">%s</p>`, template.HTMLEscapeString(msg)))
}

// Author role authorization

// canAuthorEditPost checks if an author can edit a specific post.
// Returns true if user is admin/editor, or if author owns the post.
// Returns false with error response if author doesn't own the post.
func canAuthorEditPost(c *gin.Context, post *store.BlogPost) bool {
	role, err := getUserRole(c)
	if err != nil {
		return false
	}

	// Admin and editor can edit any post
	if role == "admin" || role == "editor" {
		return true
	}

	// Author can only edit own posts
	if role == "author" {
		userID, err := getUserID(c)
		if err != nil {
			return false
		}
		return post.AuthorID.Valid && post.AuthorID.String == userID
	}

	return false
}

// denyAuthorEdit sends a forbidden response for author trying to edit others' posts.
func denyAuthorEdit(c *gin.Context, action string) {
	respondError(c, http.StatusForbidden, ErrOwnPostsOnly, fmt.Sprintf("can only %s own posts", action))
}

// HTML generation helpers

// StatCard represents a stat card for display.
type StatCard struct {
	Value     interface{}
	Label     string
	Secondary bool
}

// buildStatCards generates HTML for a list of stat cards.
func buildStatCards(cards []StatCard) string {
	var html strings.Builder
	for _, card := range cards {
		class := "stat-card"
		if card.Secondary {
			class += " secondary"
		}
		html.WriteString(fmt.Sprintf(`
		<div class="%s">
			<div class="stat-value">%s</div>
			<div class="stat-label">%s</div>
		</div>`, class, template.HTMLEscapeString(fmt.Sprintf("%v", card.Value)), template.HTMLEscapeString(card.Label)))
	}
	return html.String()
}

// TableColumn defines a column in a data table.
type TableColumn struct {
	Header string
	Class  string
}

// TableRow represents a row of cells.
type TableRow struct {
	Cells   []string
	Actions string // Optional action buttons HTML
}

// buildDataTable generates a standard data table HTML.
func buildDataTable(columns []TableColumn, rows []TableRow) string {
	var html strings.Builder

	html.WriteString(`<table class="data-table">
	<thead>
		<tr>`)

	for _, col := range columns {
		html.WriteString(fmt.Sprintf(`
			<th>%s</th>`, template.HTMLEscapeString(col.Header)))
	}

	if len(rows) > 0 && rows[0].Actions != "" {
		html.WriteString(`
			<th>Actions</th>`)
	}

	html.WriteString(`
		</tr>
	</thead>
	<tbody>`)

	for _, row := range rows {
		html.WriteString(`
		<tr>`)
		for i, cell := range row.Cells {
			class := ""
			if i < len(columns) && columns[i].Class != "" {
				class = fmt.Sprintf(` class="%s"`, columns[i].Class)
			}
			html.WriteString(fmt.Sprintf(`
			<td%s>%s</td>`, class, cell))
		}
		if row.Actions != "" {
			html.WriteString(fmt.Sprintf(`
			<td class="actions">%s</td>`, row.Actions))
		}
		html.WriteString(`
		</tr>`)
	}

	html.WriteString(`
	</tbody>
</table>`)

	return html.String()
}

// escapeHTML is a shorthand for template.HTMLEscapeString.
func escapeHTML(s string) string {
	return template.HTMLEscapeString(s)
}

// getDocsURL returns the base URL from docs config, or "/" as fallback.
func (r *Router) getDocsURL() string {
	if r.config.Docs != nil && r.config.Docs.BaseURL != "" {
		return strings.TrimSuffix(r.config.Docs.BaseURL, "/")
	}
	return "/"
}

// getDocsTitle returns the title from docs config, or the fallback value.
func (r *Router) getDocsTitle(fallback string) string {
	if r.config.Docs != nil && r.config.Docs.Title != "" {
		return r.config.Docs.Title
	}
	return fallback
}

// getPublicPageData returns common template data for public pages (blog, forum).
func (r *Router) getPublicPageData(c *gin.Context, currentPage string) gin.H {
	// Try context first, then query param, then config default
	siteID, _ := getSiteID(c)
	if siteID == "" {
		siteID = c.Query("site_id")
	}
	if siteID == "" && r.config.Docs != nil {
		siteID = r.config.Docs.SiteID
	}

	siteName := "MinimalDoc"
	if siteID != "" {
		if site, err := r.db.GetSiteByID(c.Request.Context(), siteID); err == nil && site != nil {
			siteName = site.Name
		}
	}

	// Check if user is authenticated via cookie
	var user map[string]interface{}
	if token, err := c.Cookie(r.config.Auth.SessionCookieKey); err == nil && token != "" {
		if claims, err := auth.ValidateToken(token, r.config.Auth.JWTSecret); err == nil {
			user = map[string]interface{}{
				"id":    claims.UserID,
				"email": claims.Email,
				"name":  claims.Email, // Use email as name fallback
				"role":  claims.Role,
			}
		}
	}

	return gin.H{
		"site_id":       siteID,
		"site_name":     r.getDocsTitle(siteName),
		"current_page":  currentPage,
		"docs_url":      r.getDocsURL(),
		"user":          user,
		"authenticated": user != nil,
		"Nonce":         cspNonce(c),
	}
}

// getSiteIDWithFallback gets site ID from context, query param, or config.
func (r *Router) getSiteIDWithFallback(c *gin.Context) string {
	// Try context first (set by auth middleware)
	if siteID, err := getSiteID(c); err == nil && siteID != "" {
		return siteID
	}
	// Try query param
	if siteID := c.Query("site_id"); siteID != "" {
		return siteID
	}
	// Fallback to config
	if r.config.Docs != nil && r.config.Docs.SiteID != "" {
		return r.config.Docs.SiteID
	}
	return ""
}
