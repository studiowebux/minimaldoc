package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
)

// Admin UI handlers

func (r *Router) adminDashboard(c *gin.Context) {
	claims, err := getUserClaims(c)
	if err != nil {
		c.Redirect(http.StatusFound, r.config.Server.AdminPath+"/login")
		return
	}
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Title":       "Dashboard",
		"CurrentPage": "dashboard",
		"User":        claims,
		"Nonce":       cspNonce(c),
	})
}

func (r *Router) adminLogin(c *gin.Context) {
	// Get OAuth providers for login page
	var providers []string
	for _, p := range r.config.OAuth.Providers {
		providers = append(providers, p.Name)
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":           "Login",
		"oauth_providers": providers,
	})
}

func (r *Router) adminLoginPost(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Get first site (for now, single-tenant)
	sites, err := r.db.ListSites(c.Request.Context())
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login",
			"error": "Database error. Please try again.",
		})
		return
	}
	if len(sites) == 0 {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login",
			"error": "No site configured. Run bootstrap first.",
		})
		return
	}

	siteID := sites[0].ID

	// Get user
	user, err := r.db.GetUserByEmail(c.Request.Context(), siteID, email)
	if err != nil || user == nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login",
			"error": "Invalid credentials",
		})
		return
	}

	// Verify password with constant-time comparison to prevent timing attacks.
	passwordHash := user.PasswordHash.String
	if !user.PasswordHash.Valid {
		// Use a dummy hash to maintain constant time
		passwordHash = "$2a$12$000000000000000000000.0000000000000000000000000000000"
	}
	if !auth.VerifyPassword(password, passwordHash) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login",
			"error": "Invalid credentials",
		})
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role, user.SiteID, r.config.Auth.JWTSecret, r.config.Auth.JWTExpiry)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login",
			"error": "Authentication failed",
		})
		return
	}

	// Set cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(r.config.Auth.SessionCookieKey, token, int(r.config.Auth.JWTExpiry.Seconds()), "/", "", r.config.Auth.SecureCookies, true)

	// Update last login (log error but don't fail login)
	if err := r.db.UpdateUserLastLogin(c.Request.Context(), user.ID); err != nil {
		slog.Error("failed to update last login", "error", err)
	}

	c.Redirect(http.StatusFound, r.config.Server.AdminPath)
}

func (r *Router) adminLogout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(r.config.Auth.SessionCookieKey, "", -1, "/", "", r.config.Auth.SecureCookies, true)
	c.Redirect(http.StatusFound, r.config.Server.AdminPath+"/login")
}

func (r *Router) adminAnalytics(c *gin.Context) {
	claims, err := getUserClaims(c)
	if err != nil {
		c.Redirect(http.StatusFound, r.config.Server.AdminPath+"/login")
		return
	}
	c.HTML(http.StatusOK, "analytics.html", gin.H{
		"Title":       "Analytics",
		"CurrentPage": "analytics",
		"User":        claims,
		"Nonce":       cspNonce(c),
	})
}

func (r *Router) adminFeedback(c *gin.Context) {
	claims, err := getUserClaims(c)
	if err != nil {
		c.Redirect(http.StatusFound, r.config.Server.AdminPath+"/login")
		return
	}
	c.HTML(http.StatusOK, "feedback.html", gin.H{
		"Title":       "Feedback",
		"CurrentPage": "feedback",
		"User":        claims,
		"Nonce":       cspNonce(c),
	})
}

func (r *Router) adminSubscribers(c *gin.Context) {
	claims, err := getUserClaims(c)
	if err != nil {
		c.Redirect(http.StatusFound, r.config.Server.AdminPath+"/login")
		return
	}
	c.HTML(http.StatusOK, "subscribers.html", gin.H{
		"Title":       "Subscribers",
		"CurrentPage": "subscribers",
		"User":        claims,
		"Nonce":       cspNonce(c),
	})
}

func (r *Router) adminSettings(c *gin.Context) {
	claims, err := getUserClaims(c)
	if err != nil {
		c.Redirect(http.StatusFound, r.config.Server.AdminPath+"/login")
		return
	}
	siteID, err := getSiteID(c)
	if err != nil {
		c.Redirect(http.StatusFound, r.config.Server.AdminPath+"/login")
		return
	}
	site, err := r.db.GetSiteByID(c.Request.Context(), siteID)
	if err != nil {
		slog.Error("failed to get site for settings", "error", err)
	}

	serverURL := r.config.Email.BaseURL
	if serverURL == "" {
		serverURL = fmt.Sprintf("http://localhost:%d", r.config.Server.Port)
	}

	apiKey := ""
	if site != nil {
		apiKey = site.APIKey
	}

	c.HTML(http.StatusOK, "settings.html", gin.H{
		"Title":       "Settings",
		"CurrentPage": "settings",
		"User":        claims,
		"Nonce":       cspNonce(c),
		"SiteID":      siteID,
		"APIKey":      apiKey,
		"ServerURL":   serverURL,
	})
}
