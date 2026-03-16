package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// adminDocAccess renders the doc access management page.
func (r *Router) adminDocAccess(c *gin.Context) {
	claims, _ := getUserClaims(c)

	c.HTML(http.StatusOK, "doc-access.html", gin.H{
		"Title":       "Doc Access Rules",
		"User":        claims,
		"AdminPath":   r.config.Server.AdminPath,
		"CurrentPage": "doc-access",
		"Nonce":       cspNonce(c),
	})
}

// fragmentDocAccessList renders the list of doc access rules.
func (r *Router) fragmentDocAccessList(c *gin.Context) {
	siteID, _ := getSiteID(c)

	rules, err := r.db.ListDocAccess(c.Request.Context(), siteID)
	if err != nil {
		respondHTMLError(c, "Error loading rules")
		return
	}

	if len(rules) == 0 {
		respondHTMLEmpty(c, "No access rules defined. All documentation is public.")
		return
	}

	html := `<table class="data-table">
		<thead>
			<tr>
				<th>Path Pattern</th>
				<th>Required Role</th>
				<th>Description</th>
				<th>Actions</th>
			</tr>
		</thead>
		<tbody>`

	for _, rule := range rules {
		desc := rule.Description.String
		if !rule.Description.Valid || desc == "" {
			desc = "-"
		}
		html += fmt.Sprintf(`
			<tr>
				<td class="path"><code>%s</code></td>
				<td class="role"><span class="role-badge role-%s">%s</span></td>
				<td class="description">%s</td>
				<td class="actions">
					<button class="btn btn-sm" hx-get="/admin/fragments/doc-access-form/%s" hx-target="#rule-form-container" hx-swap="innerHTML">Edit</button>
					<button class="btn btn-sm btn-danger" hx-delete="/api/docs/rules/%s" hx-confirm="Delete this rule?" hx-target="#rule-list" hx-swap="outerHTML" hx-trigger="click">Delete</button>
				</td>
			</tr>`,
			escapeHTML(rule.PathPattern),
			escapeHTML(rule.RequiredRole),
			escapeHTML(rule.RequiredRole),
			escapeHTML(desc),
			rule.ID,
			rule.ID,
		)
	}

	html += `</tbody></table>`

	respondHTML(c, html)
}

// fragmentDocAccessForm renders the form for creating/editing a doc access rule.
func (r *Router) fragmentDocAccessForm(c *gin.Context) {
	id := c.Param("id")

	var pathPattern, requiredRole, description string
	isEdit := false

	if id != "" {
		existing, err := r.db.GetDocAccessByID(c.Request.Context(), id)
		if err != nil {
			respondHTMLError(c, "Error loading rule")
			return
		}
		if existing == nil {
			respondHTMLError(c, "Rule not found")
			return
		}
		pathPattern = existing.PathPattern
		requiredRole = existing.RequiredRole
		description = existing.Description.String
		isEdit = true
	}

	// Build role options
	roles := []string{"viewer", "author", "editor", "admin"}
	roleOptions := ""
	for _, role := range roles {
		selected := ""
		if role == requiredRole {
			selected = " selected"
		}
		roleOptions += fmt.Sprintf(`<option value="%s"%s>%s</option>`, role, selected, role)
	}

	// Determine form action
	formMethod := "POST"
	formAction := "/api/docs/rules"
	submitText := "Create Rule"
	if isEdit {
		formMethod = "PUT"
		formAction = fmt.Sprintf("/api/docs/rules/%s", id)
		submitText = "Update Rule"
	}

	html := fmt.Sprintf(`
		<form class="rule-form" hx-%s="%s" hx-target="#rule-list" hx-swap="innerHTML" hx-on::after-request="if(event.detail.successful) { document.getElementById('rule-form-container').innerHTML = ''; htmx.ajax('GET', '/admin/fragments/doc-access-list', '#rule-list'); }">
			<div class="form-row">
				<div class="form-group">
					<label for="path_pattern">Path Pattern</label>
					<input type="text" id="path_pattern" name="path_pattern" value="%s" placeholder="/docs/internal/**" required>
					<small>Use * for single segment, ** for recursive match</small>
				</div>
				<div class="form-group">
					<label for="required_role">Required Role</label>
					<select id="required_role" name="required_role" required>
						%s
					</select>
				</div>
			</div>
			<div class="form-group">
				<label for="description">Description (optional)</label>
				<input type="text" id="description" name="description" value="%s" placeholder="Internal API documentation">
			</div>
			<div class="form-actions">
				<button type="submit" class="btn btn-primary">%s</button>
				<button type="button" class="btn" onclick="document.getElementById('rule-form-container').innerHTML = ''">Cancel</button>
			</div>
		</form>
	`,
		formMethod,
		formAction,
		escapeHTML(pathPattern),
		roleOptions,
		escapeHTML(description),
		submitText,
	)

	respondHTML(c, html)
}
