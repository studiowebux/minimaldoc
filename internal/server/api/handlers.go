package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
	"github.com/studiowebux/minimaldoc/internal/server/email"
)

// BootstrapRequest represents bootstrap parameters.
type BootstrapRequest struct {
	SiteName string `json:"site_name" binding:"required"`
	Domain   string `json:"domain"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
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

	result, err := Bootstrap(c.Request.Context(), r.db, r.config, req.SiteName, req.Domain, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"site_id":  result.SiteID,
		"api_key":  result.APIKey,
		"user_id":  result.UserID,
		"email":    result.Email,
		"password": result.Password,
		"message":  "Bootstrap complete. Save these credentials!",
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

	// Verify password
	if !user.PasswordHash.Valid || !auth.VerifyPassword(req.Password, user.PasswordHash.String) {
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

	// Set cookie for admin UI
	c.SetCookie(r.config.Auth.SessionCookieKey, accessToken, int(r.config.Auth.JWTExpiry.Seconds()), "/", "", false, true)

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
	c.SetCookie(r.config.Auth.SessionCookieKey, "", -1, "/", "", false, true)
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
	userID, _ := c.Get("user_id")

	user, err := r.db.GetUserByID(c.Request.Context(), userID.(string))
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
	state, _ := auth.GenerateSessionToken()

	// Store state in cookie (short-lived)
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

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
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

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
	c.SetCookie(r.config.Auth.SessionCookieKey, accessToken, int(r.config.Auth.JWTExpiry.Seconds()), "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, r.config.Server.AdminPath)
}

// Analytics handlers

type TrackRequest struct {
	SiteID      string `json:"site_id" binding:"required"`
	Path        string `json:"path" binding:"required"`
	Referrer    string `json:"referrer"`
	Country     string `json:"country"`
	DeviceType  string `json:"device_type"`
	Browser     string `json:"browser"`
	OS          string `json:"os"`
	SessionHash string `json:"session_hash"`
}

func (r *Router) trackPageView(c *gin.Context) {
	var req TrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.db.RecordPageView(c.Request.Context(),
		req.SiteID, req.Path, req.Referrer, req.Country,
		req.DeviceType, req.Browser, req.OS, req.SessionHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "tracked"})
}

func (r *Router) trackEvent(c *gin.Context) {
	// Events table not implemented yet, just acknowledge
	c.JSON(http.StatusOK, gin.H{"status": "tracked"})
}

// DurationRequest represents a page duration update
type DurationRequest struct {
	SiteID      string `json:"site_id" binding:"required"`
	Path        string `json:"path" binding:"required"`
	Duration    int    `json:"duration" binding:"required,min=1"`
	SessionHash string `json:"session_hash"`
}

func (r *Router) trackDuration(c *gin.Context) {
	var req DurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cap duration at a reasonable maximum (30 minutes = 1800 seconds)
	duration := req.Duration
	if duration > 1800 {
		duration = 1800
	}

	err := r.db.UpdatePageViewDuration(c.Request.Context(), req.SiteID, req.Path, req.SessionHash, duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update duration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (r *Router) analyticsSummary(c *gin.Context) {
	siteID, _ := c.Get("site_id")
	since := time.Now().Add(-24 * time.Hour) // Last 24 hours

	totalViews, uniqueSessions, err := r.db.GetPageViewStats(c.Request.Context(), siteID.(string), since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	topPages, _ := r.db.GetTopPages(c.Request.Context(), siteID.(string), since, 10)

	c.JSON(http.StatusOK, gin.H{
		"total_views":     totalViews,
		"unique_visitors": uniqueSessions,
		"top_pages":       topPages,
		"period":          "24h",
	})
}

func (r *Router) analyticsPages(c *gin.Context) {
	siteID, _ := c.Get("site_id")
	since := time.Now().Add(-7 * 24 * time.Hour) // Last 7 days

	pages, err := r.db.GetTopPages(c.Request.Context(), siteID.(string), since, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pages": pages, "period": "7d"})
}

// Feedback handlers

type FeedbackRequest struct {
	SiteID      string `json:"site_id" binding:"required"`
	Path        string `json:"path" binding:"required"`
	Rating      int    `json:"rating" binding:"required,min=1,max=5"`
	Feedback    string `json:"feedback"`
	SessionHash string `json:"session_hash"`
}

func (r *Router) submitFeedback(c *gin.Context) {
	var req FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.db.RecordRating(c.Request.Context(), req.SiteID, req.Path, req.Rating, req.Feedback, req.SessionHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "submitted"})
}

func (r *Router) feedbackStats(c *gin.Context) {
	siteID, _ := c.Get("site_id")

	avgRating, totalRatings, err := r.db.GetRatingStats(c.Request.Context(), siteID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"average_rating": avgRating,
		"total_ratings":  totalRatings,
	})
}

func (r *Router) feedbackList(c *gin.Context) {
	siteID, _ := c.Get("site_id")

	ratings, err := r.db.ListRatings(c.Request.Context(), siteID.(string), 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"feedback": ratings})
}

// Newsletter handlers

type SubscribeRequest struct {
	SiteID string `json:"site_id" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
}

func (r *Router) subscribe(c *gin.Context) {
	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get site name for email
	site, err := r.db.GetSiteByID(c.Request.Context(), req.SiteID)
	if err != nil || site == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid site_id"})
		return
	}

	// Check if already subscribed
	existing, _ := r.db.GetSubscriberByEmail(c.Request.Context(), req.SiteID, req.Email)
	if existing != nil {
		if existing.Verified {
			c.JSON(http.StatusOK, gin.H{"status": "already_subscribed"})
			return
		}
		// Resend verification email
		r.sendVerificationEmail(site.Name, req.SiteID, req.Email, existing.VerifyToken.String)
		c.JSON(http.StatusOK, gin.H{"status": "verification_resent"})
		return
	}

	// Create subscriber with verification token
	id, _ := auth.GenerateSessionToken()
	verifyToken, _ := auth.GenerateVerificationToken()

	err = r.db.CreateSubscriber(c.Request.Context(), id, req.SiteID, req.Email, verifyToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to subscribe"})
		return
	}

	// Send verification email
	r.sendVerificationEmail(site.Name, req.SiteID, req.Email, verifyToken)

	c.JSON(http.StatusOK, gin.H{"status": "verification_sent"})
}

func (r *Router) sendVerificationEmail(siteName, siteID, emailAddr, token string) {
	templates := email.NewTemplates(siteName, r.config.Email.BaseURL)
	msg := templates.VerificationEmail(emailAddr, siteID, token)

	if err := r.email.Send(context.Background(), msg); err != nil {
		log.Printf("Failed to send verification email to %s: %v", emailAddr, err)
	}
}

func (r *Router) verifySubscription(c *gin.Context) {
	siteID := c.Query("site_id")
	token := c.Query("token")

	if siteID == "" || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing parameters"})
		return
	}

	// Get subscriber before verification to get email
	subscriber, _ := r.db.GetSubscriberByToken(c.Request.Context(), siteID, token)
	if subscriber == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	err := r.db.VerifySubscriber(c.Request.Context(), siteID, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	// Send welcome email
	site, _ := r.db.GetSiteByID(c.Request.Context(), siteID)
	if site != nil {
		r.sendWelcomeEmail(site.Name, siteID, subscriber.Email)
	}

	c.JSON(http.StatusOK, gin.H{"status": "verified"})
}

func (r *Router) sendWelcomeEmail(siteName, siteID, emailAddr string) {
	templates := email.NewTemplates(siteName, r.config.Email.BaseURL)
	msg := templates.WelcomeEmail(emailAddr, siteID)

	if err := r.email.Send(context.Background(), msg); err != nil {
		log.Printf("Failed to send welcome email to %s: %v", emailAddr, err)
	}
}

func (r *Router) unsubscribe(c *gin.Context) {
	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.db.UnsubscribeByEmail(c.Request.Context(), req.SiteID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "unsubscribed"})
}

func (r *Router) listSubscribers(c *gin.Context) {
	siteID, _ := c.Get("site_id")

	subscribers, err := r.db.ListSubscribers(c.Request.Context(), siteID.(string), false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	count, _ := r.db.CountSubscribers(c.Request.Context(), siteID.(string), true)

	c.JSON(http.StatusOK, gin.H{
		"subscribers":      subscribers,
		"verified_count":   count,
	})
}

// Site management handlers

func (r *Router) listSites(c *gin.Context) {
	// TODO: Implement site list
	c.JSON(http.StatusOK, gin.H{"sites": []string{}})
}

func (r *Router) createSite(c *gin.Context) {
	// TODO: Implement site creation
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) getSite(c *gin.Context) {
	// TODO: Implement get site
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) updateSite(c *gin.Context) {
	// TODO: Implement site update
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) deleteSite(c *gin.Context) {
	// TODO: Implement site deletion
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) regenerateAPIKey(c *gin.Context) {
	// TODO: Implement API key regeneration
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// User management handlers

func (r *Router) listUsers(c *gin.Context) {
	// TODO: Implement user list
	c.JSON(http.StatusOK, gin.H{"users": []string{}})
}

func (r *Router) createUser(c *gin.Context) {
	// TODO: Implement user creation
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) getUser(c *gin.Context) {
	// TODO: Implement get user
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) updateUser(c *gin.Context) {
	// TODO: Implement user update
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) deleteUser(c *gin.Context) {
	// TODO: Implement user deletion
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// Admin UI handlers

func (r *Router) adminDashboard(c *gin.Context) {
	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":  "Dashboard",
		"active": "dashboard",
		"user":   user,
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
	sites, _ := r.db.ListSites(c.Request.Context())
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

	// Verify password
	if !user.PasswordHash.Valid || !auth.VerifyPassword(password, user.PasswordHash.String) {
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
	c.SetCookie(r.config.Auth.SessionCookieKey, token, int(r.config.Auth.JWTExpiry.Seconds()), "/", "", false, true)

	// Update last login
	_ = r.db.UpdateUserLastLogin(c.Request.Context(), user.ID)

	c.Redirect(http.StatusFound, r.config.Server.AdminPath)
}

func (r *Router) adminLogout(c *gin.Context) {
	c.SetCookie(r.config.Auth.SessionCookieKey, "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, r.config.Server.AdminPath+"/login")
}

func (r *Router) adminAnalytics(c *gin.Context) {
	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "analytics.html", gin.H{
		"title":  "Analytics",
		"active": "analytics",
		"user":   user,
	})
}

func (r *Router) adminFeedback(c *gin.Context) {
	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "feedback.html", gin.H{
		"title":  "Feedback",
		"active": "feedback",
		"user":   user,
	})
}

func (r *Router) adminSubscribers(c *gin.Context) {
	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "subscribers.html", gin.H{
		"title":  "Subscribers",
		"active": "subscribers",
		"user":   user,
	})
}

func (r *Router) adminSettings(c *gin.Context) {
	user, _ := c.Get("user")
	siteID, _ := c.Get("site_id")
	site, _ := r.db.GetSiteByID(c.Request.Context(), siteID.(string))

	serverURL := r.config.Email.BaseURL
	if serverURL == "" {
		serverURL = fmt.Sprintf("http://localhost:%d", r.config.Server.Port)
	}

	apiKey := ""
	if site != nil {
		apiKey = site.APIKey
	}

	c.HTML(http.StatusOK, "settings.html", gin.H{
		"title":      "Settings",
		"active":     "settings",
		"user":       user,
		"site_id":    siteID,
		"api_key":    apiKey,
		"server_url": serverURL,
	})
}
