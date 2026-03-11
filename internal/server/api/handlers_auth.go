package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
)

// BootstrapRequest represents bootstrap parameters.
type BootstrapRequest struct {
	SiteName       string `json:"site_name" binding:"required"`
	Domain         string `json:"domain"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password"`
	BootstrapToken string `json:"bootstrap_token"`
}

// LoginRequest represents login credentials.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	SiteID   string `json:"site_id" binding:"required"`
}

// Bootstrap handler

func (r *Router) bootstrap(c *gin.Context) {
	var req BootstrapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check bootstrap token if configured
	if r.config.Auth.BootstrapToken != "" {
		if req.BootstrapToken != r.config.Auth.BootstrapToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid bootstrap token"})
			return
		}
	}

	// Require password to be provided (never auto-generate and return)
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters"})
		return
	}

	result, err := Bootstrap(c.Request.Context(), r.db, r.config, req.SiteName, req.Domain, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Never return password in response
	c.JSON(http.StatusOK, gin.H{
		"site_id": result.SiteID,
		"api_key": result.APIKey,
		"user_id": result.UserID,
		"email":   result.Email,
		"message": "Bootstrap complete. Save your API key securely.",
	})
}

// Auth handlers

func (r *Router) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from database
	user, err := r.db.GetUserByEmail(c.Request.Context(), req.SiteID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Verify password with constant-time comparison to prevent timing attacks.
	// Always call VerifyPassword even if hash is invalid to prevent timing side-channel.
	passwordHash := user.PasswordHash.String
	if !user.PasswordHash.Valid {
		// Use a dummy hash to maintain constant time - this will always fail
		passwordHash = "$2a$12$000000000000000000000.0000000000000000000000000000000"
	}
	if !auth.VerifyPassword(req.Password, passwordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate tokens
	accessToken, err := auth.GenerateToken(user.ID, user.Email, user.Role, user.SiteID, r.config.Auth.JWTSecret, r.config.Auth.JWTExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, r.config.Auth.JWTSecret, r.config.Auth.RefreshExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	// Update last login
	_ = r.db.UpdateUserLastLogin(c.Request.Context(), user.ID)

	// Set cookie for admin UI with SameSite protection
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(r.config.Auth.SessionCookieKey, accessToken, int(r.config.Auth.JWTExpiry.Seconds()), "/", "", r.config.Auth.SecureCookies, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    int(r.config.Auth.JWTExpiry.Seconds()),
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
			"name":  user.Name.String,
		},
	})
}

func (r *Router) logout(c *gin.Context) {
	// Clear cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(r.config.Auth.SessionCookieKey, "", -1, "/", "", r.config.Auth.SecureCookies, true)
	c.JSON(http.StatusOK, gin.H{"status": "logged out"})
}

func (r *Router) refreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate refresh token
	userID, err := auth.ValidateRefreshToken(req.RefreshToken, r.config.Auth.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// Get user
	user, err := r.db.GetUserByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// Generate new access token
	accessToken, err := auth.GenerateToken(user.ID, user.Email, user.Role, user.SiteID, r.config.Auth.JWTSecret, r.config.Auth.JWTExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expires_in":   int(r.config.Auth.JWTExpiry.Seconds()),
	})
}

func (r *Router) getCurrentUser(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := r.db.GetUserByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             user.ID,
		"email":          user.Email,
		"role":           user.Role,
		"name":           user.Name.String,
		"avatar_url":     user.AvatarURL.String,
		"email_verified": user.EmailVerified,
		"created_at":     user.CreatedAt,
	})
}

func (r *Router) listOAuthProviders(c *gin.Context) {
	providers := make([]string, 0)
	for _, p := range r.config.OAuth.Providers {
		providers = append(providers, p.Name)
	}
	c.JSON(http.StatusOK, gin.H{"providers": providers})
}

func (r *Router) oauthLogin(c *gin.Context) {
	provider := c.Param("provider")

	// Find provider config
	var providerCfg *auth.OAuthProvider
	for _, p := range r.config.OAuth.Providers {
		if p.Name == provider {
			op := auth.NewOAuthProvider(p)
			providerCfg = op
			break
		}
	}

	if providerCfg == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
		return
	}

	// Generate state token
	state, err := auth.GenerateSessionToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}

	// Store state in cookie (short-lived) with SameSite=Lax for CSRF protection
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("oauth_state", state, 600, "/", "", r.config.Auth.SecureCookies, true)

	// Redirect to provider
	authURL := providerCfg.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (r *Router) oauthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")

	// Verify state
	savedState, err := c.Cookie("oauth_state")
	if err != nil || savedState != state {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	// Clear state cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("oauth_state", "", -1, "/", "", r.config.Auth.SecureCookies, true)

	// Find provider
	var providerCfg *auth.OAuthProvider
	for _, p := range r.config.OAuth.Providers {
		if p.Name == provider {
			op := auth.NewOAuthProvider(p)
			providerCfg = op
			break
		}
	}

	if providerCfg == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
		return
	}

	// Exchange code for user info
	userInfo, err := providerCfg.HandleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find or create user
	user, err := r.db.GetUserByOAuth(c.Request.Context(), provider, userInfo.ProviderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if user == nil {
		// Create new user (would need site_id from somewhere - could be in state)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found, registration required"})
		return
	}

	// Generate token
	accessToken, err := auth.GenerateToken(user.ID, user.Email, user.Role, user.SiteID, r.config.Auth.JWTSecret, r.config.Auth.JWTExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Update last login
	_ = r.db.UpdateUserLastLogin(c.Request.Context(), user.ID)

	// Set cookie and redirect to admin
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(r.config.Auth.SessionCookieKey, accessToken, int(r.config.Auth.JWTExpiry.Seconds()), "/", "", r.config.Auth.SecureCookies, true)
	c.Redirect(http.StatusTemporaryRedirect, r.config.Server.AdminPath)
}
