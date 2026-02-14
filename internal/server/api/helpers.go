package api

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/store"
)

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
	c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("can only %s own posts", action)})
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
			<div class="stat-value">%v</div>
			<div class="stat-label">%s</div>
		</div>`, class, card.Value, template.HTMLEscapeString(card.Label)))
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

// formatInt formats an integer for display.
func formatInt(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(n)/1000000)
}
