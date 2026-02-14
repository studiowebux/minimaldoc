package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// adminUsers renders the user management page.
func (r *Router) adminUsers(c *gin.Context) {
	claims, _ := getUserClaims(c)

	c.HTML(http.StatusOK, "users.html", gin.H{
		"Title":       "User Management",
		"User":        claims,
		"AdminPath":   r.config.Server.AdminPath,
		"CurrentPage": "users",
	})
}

// fragmentUserStats renders user statistics.
func (r *Router) fragmentUserStats(c *gin.Context) {
	siteID, _ := getSiteID(c)

	total, err := r.db.CountUsers(c.Request.Context(), siteID)
	if err != nil {
		respondHTMLError(c, "Error loading stats")
		return
	}

	admins, _ := r.db.CountUsersByRole(c.Request.Context(), siteID, "admin")
	editors, _ := r.db.CountUsersByRole(c.Request.Context(), siteID, "editor")
	authors, _ := r.db.CountUsersByRole(c.Request.Context(), siteID, "author")
	viewers, _ := r.db.CountUsersByRole(c.Request.Context(), siteID, "viewer")

	html := buildStatCards([]StatCard{
		{Value: total, Label: "Total Users"},
		{Value: admins, Label: "Admins"},
		{Value: editors, Label: "Editors"},
		{Value: authors + viewers, Label: "Authors + Viewers"},
	})

	respondHTML(c, html)
}

// fragmentUserList renders the list of users.
func (r *Router) fragmentUserList(c *gin.Context) {
	siteID, _ := getSiteID(c)
	roleFilter := c.Query("role-filter")

	users, err := r.db.ListUsers(c.Request.Context(), siteID)
	if err != nil {
		respondHTMLError(c, "Error loading users")
		return
	}

	// Filter by role if specified
	if roleFilter != "" {
		var filtered []struct {
			ID          string
			Email       string
			Role        string
			Name        string
			LastLoginAt string
			CreatedAt   string
		}
		for _, u := range users {
			if u.Role == roleFilter {
				name := ""
				if u.Name.Valid {
					name = u.Name.String
				}
				lastLogin := "-"
				if u.LastLoginAt.Valid {
					lastLogin = u.LastLoginAt.String
				}
				filtered = append(filtered, struct {
					ID          string
					Email       string
					Role        string
					Name        string
					LastLoginAt string
					CreatedAt   string
				}{u.ID, u.Email, u.Role, name, lastLogin, u.CreatedAt})
			}
		}
		if len(filtered) == 0 {
			respondHTMLEmpty(c, "No users with this role")
			return
		}

		html := buildUserTable(filtered)
		respondHTML(c, html)
		return
	}

	if len(users) == 0 {
		respondHTMLEmpty(c, "No users found")
		return
	}

	// Convert to display format
	var displayUsers []struct {
		ID          string
		Email       string
		Role        string
		Name        string
		LastLoginAt string
		CreatedAt   string
	}
	for _, u := range users {
		name := ""
		if u.Name.Valid {
			name = u.Name.String
		}
		lastLogin := "-"
		if u.LastLoginAt.Valid {
			lastLogin = u.LastLoginAt.String
		}
		displayUsers = append(displayUsers, struct {
			ID          string
			Email       string
			Role        string
			Name        string
			LastLoginAt string
			CreatedAt   string
		}{u.ID, u.Email, u.Role, name, lastLogin, u.CreatedAt})
	}

	html := buildUserTable(displayUsers)
	respondHTML(c, html)
}

func buildUserTable(users []struct {
	ID          string
	Email       string
	Role        string
	Name        string
	LastLoginAt string
	CreatedAt   string
}) string {
	html := `<table class="data-table">
		<thead>
			<tr>
				<th>Email</th>
				<th>Name</th>
				<th>Role</th>
				<th>Last Login</th>
				<th>Actions</th>
			</tr>
		</thead>
		<tbody>`

	for _, u := range users {
		html += fmt.Sprintf(`
			<tr>
				<td class="email">%s</td>
				<td class="name">%s</td>
				<td class="role"><span class="role-badge role-%s">%s</span></td>
				<td class="date">%s</td>
				<td class="actions">
					<button class="btn btn-sm" hx-get="/admin/fragments/user-form/%s" hx-target="#user-form-container" hx-swap="innerHTML">Edit</button>
					<button class="btn btn-sm btn-danger" hx-delete="/api/users/%s" hx-confirm="Delete this user?" hx-target="#user-list" hx-swap="innerHTML" hx-trigger="click" hx-on::after-request="htmx.ajax('GET', '/admin/fragments/user-list', '#user-list');">Delete</button>
				</td>
			</tr>`,
			escapeHTML(u.Email),
			escapeHTML(u.Name),
			escapeHTML(u.Role),
			escapeHTML(u.Role),
			u.LastLoginAt,
			u.ID,
			u.ID,
		)
	}

	html += `</tbody></table>`
	return html
}

// fragmentUserForm renders the form for creating/editing a user.
func (r *Router) fragmentUserForm(c *gin.Context) {
	id := c.Param("id")

	var email, role, name string
	isEdit := false

	if id != "" {
		existing, err := r.db.GetUserByID(c.Request.Context(), id)
		if err != nil {
			respondHTMLError(c, "Error loading user")
			return
		}
		if existing == nil {
			respondHTMLError(c, "User not found")
			return
		}
		email = existing.Email
		role = existing.Role
		if existing.Name.Valid {
			name = existing.Name.String
		}
		isEdit = true
	}

	// Build role options
	roles := []string{"viewer", "author", "editor", "admin"}
	roleOptions := ""
	for _, r := range roles {
		selected := ""
		if r == role {
			selected = " selected"
		}
		roleOptions += fmt.Sprintf(`<option value="%s"%s>%s</option>`, r, selected, r)
	}

	// Determine form action
	formMethod := "post"
	formAction := "/api/users"
	submitText := "Create User"
	passwordLabel := "Password"
	passwordRequired := "required"
	if isEdit {
		formMethod = "put"
		formAction = fmt.Sprintf("/api/users/%s", id)
		submitText = "Update User"
		passwordLabel = "New Password (leave empty to keep current)"
		passwordRequired = ""
	}

	html := fmt.Sprintf(`
		<form class="user-form" hx-%s="%s" hx-ext="json-enc" hx-target="#user-list" hx-swap="innerHTML" hx-on::after-request="if(event.detail.successful) { document.getElementById('user-form-container').innerHTML = ''; htmx.ajax('GET', '/admin/fragments/user-list', '#user-list'); htmx.ajax('GET', '/admin/fragments/user-stats', '#user-stats'); }">
			<div class="form-row">
				<div class="form-group">
					<label for="email">Email</label>
					<input type="email" id="email" name="email" value="%s" required>
				</div>
				<div class="form-group">
					<label for="role">Role</label>
					<select id="role" name="role" required>
						%s
					</select>
				</div>
			</div>
			<div class="form-row">
				<div class="form-group">
					<label for="name">Name</label>
					<input type="text" id="name" name="name" value="%s" placeholder="Display name">
				</div>
				<div class="form-group">
					<label for="password">%s</label>
					<input type="password" id="password" name="password" minlength="8" %s>
				</div>
			</div>
			<div class="form-actions">
				<button type="submit" class="btn btn-primary">%s</button>
				<button type="button" class="btn" onclick="document.getElementById('user-form-container').innerHTML = ''">Cancel</button>
			</div>
		</form>
	`,
		formMethod,
		formAction,
		escapeHTML(email),
		roleOptions,
		escapeHTML(name),
		passwordLabel,
		passwordRequired,
		submitText,
	)

	respondHTML(c, html)
}
