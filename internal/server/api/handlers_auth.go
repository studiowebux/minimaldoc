package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
	"github.com/studiowebux/minimaldoc/internal/server/email"
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
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	// Check bootstrap token if configured
	if r.config.Auth.BootstrapToken != "" {
		if req.BootstrapToken != r.config.Auth.BootstrapToken {
			respondUnauthorized(c, ErrInvalidBootstrapToken, "invalid bootstrap token")
			return
		}
	}

	// Require password to be provided (never auto-generate and return)
	if req.Password == "" {
		respondBadRequest(c, ErrPasswordRequired, "password is required")
		return
	}
	if len(req.Password) < 8 {
		respondBadRequest(c, ErrPasswordTooShort, "password must be at least 8 characters")
		return
	}

	result, err := Bootstrap(c.Request.Context(), r.db, r.config, req.SiteName, req.Domain, req.Email, req.Password)
	if err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
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
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	// Get user from database
	user, err := r.db.GetUserByEmail(c.Request.Context(), req.SiteID, req.Email)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if user == nil {
		respondUnauthorized(c, ErrInvalidCredentials, "invalid credentials")
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
		respondUnauthorized(c, ErrInvalidCredentials, "invalid credentials")
		return
	}

	// Generate tokens
	accessToken, err := auth.GenerateToken(user.ID, user.Email, user.Role, user.SiteID, r.config.Auth.JWTSecret, r.config.Auth.JWTExpiry)
	if err != nil {
		respondInternalError(c, ErrTokenGenerationFailed, "failed to generate token")
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, r.config.Auth.JWTSecret, r.config.Auth.RefreshExpiry)
	if err != nil {
		respondInternalError(c, ErrTokenGenerationFailed, "failed to generate refresh token")
		return
	}

	// Update last login
	_ = r.db.UpdateUserLastLogin(c.Request.Context(), user.ID)

	// Audit log: login (set context manually since auth middleware hasn't run)
	c.Set("site_id", user.SiteID)
	c.Set("user_id", user.ID)
	c.Set("user_email", user.Email)
	r.logAuditAction(c, "login", "session", "", user.Email, "")

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
	// Audit log: logout (best effort - user may not be authenticated)
	r.logAuditAction(c, "logout", "session", "", "", "")

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
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	// Validate refresh token
	userID, err := auth.ValidateRefreshToken(req.RefreshToken, r.config.Auth.JWTSecret)
	if err != nil {
		respondUnauthorized(c, ErrInvalidRefreshToken, "invalid refresh token")
		return
	}

	// Get user
	user, err := r.db.GetUserByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		respondUnauthorized(c, ErrUserNotFound, "user not found")
		return
	}

	// Generate new access token
	accessToken, err := auth.GenerateToken(user.ID, user.Email, user.Role, user.SiteID, r.config.Auth.JWTSecret, r.config.Auth.JWTExpiry)
	if err != nil {
		respondInternalError(c, ErrTokenGenerationFailed, "failed to generate token")
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
		respondUnauthorized(c, ErrUnauthorized, "unauthorized")
		return
	}

	user, err := r.db.GetUserByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		respondNotFound(c, ErrUserNotFound, "user not found")
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
		respondBadRequest(c, ErrUnknownProvider, "unknown provider")
		return
	}

	// Generate state token
	state, err := auth.GenerateSessionToken()
	if err != nil {
		respondInternalError(c, ErrTokenGenerationFailed, "failed to generate state")
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
		respondBadRequest(c, ErrInvalidState, "invalid state")
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
		respondBadRequest(c, ErrUnknownProvider, "unknown provider")
		return
	}

	// Exchange code for user info
	userInfo, err := providerCfg.HandleCallback(c.Request.Context(), code)
	if err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	// Find or create user
	user, err := r.db.GetUserByOAuth(c.Request.Context(), provider, userInfo.ProviderID)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	if user == nil {
		// Create new user (would need site_id from somewhere - could be in state)
		respondBadRequest(c, ErrRegistrationRequired, "user not found, registration required")
		return
	}

	// Generate token
	accessToken, err := auth.GenerateToken(user.ID, user.Email, user.Role, user.SiteID, r.config.Auth.JWTSecret, r.config.Auth.JWTExpiry)
	if err != nil {
		respondInternalError(c, ErrTokenGenerationFailed, "failed to generate token")
		return
	}

	// Update last login
	_ = r.db.UpdateUserLastLogin(c.Request.Context(), user.ID)

	// Set cookie and redirect to admin
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(r.config.Auth.SessionCookieKey, accessToken, int(r.config.Auth.JWTExpiry.Seconds()), "/", "", r.config.Auth.SecureCookies, true)
	c.Redirect(http.StatusTemporaryRedirect, r.config.Server.AdminPath)
}

// RegisterRequest represents public registration parameters.
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name"`
	SiteID   string `json:"site_id" binding:"required"`
}

// register handles public user registration.
func (r *Router) register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	// Validate site exists
	site, err := r.db.GetSiteByID(c.Request.Context(), req.SiteID)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if site == nil {
		respondBadRequest(c, ErrSiteInvalid, "invalid site")
		return
	}

	// Check if email already exists
	existing, err := r.db.GetUserByEmail(c.Request.Context(), req.SiteID, req.Email)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if existing != nil {
		respondError(c, http.StatusConflict, ErrUserAlreadyExists, "email already registered")
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password, r.config.Auth.BCryptCost)
	if err != nil {
		respondInternalError(c, ErrPasswordHashFailed, "failed to hash password")
		return
	}

	// Generate verification token
	verifyToken, err := auth.GenerateSessionToken()
	if err != nil {
		respondInternalError(c, ErrTokenGenerationFailed, "failed to generate token")
		return
	}

	// Create user with role=viewer, email_verified=false
	userID := uuid.New().String()
	name := req.Name
	if name == "" {
		name = strings.Split(req.Email, "@")[0]
	}

	_, err = r.db.CreateUserWithVerification(c.Request.Context(), userID, req.SiteID, req.Email, hashedPassword, "viewer", name, verifyToken)
	if err != nil {
		respondInternalError(c, ErrUserCreationFailed, "failed to create user")
		return
	}

	// Send verification email
	if r.email != nil {
		verifyURL := r.config.Email.BaseURL + r.config.Server.APIPath + "/auth/verify?token=" + verifyToken
		msg := &email.Message{
			To:      req.Email,
			Subject: "Verify your email address",
			HTMLBody: `<h2>Welcome!</h2>
<p>Hi ` + name + `,</p>
<p>Please verify your email address by clicking the link below:</p>
<p><a href="` + verifyURL + `">Verify Email</a></p>
<p>Or copy this URL: ` + verifyURL + `</p>
<p>This link will expire in 24 hours.</p>`,
			TextBody: "Hi " + name + ",\n\nPlease verify your email address by visiting:\n" + verifyURL + "\n\nThis link will expire in 24 hours.",
		}
		_ = r.email.Send(c.Request.Context(), msg)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful. Please check your email to verify your account.",
		"user_id": userID,
	})
}

// verifyEmail handles email verification.
func (r *Router) verifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		// Check if request wants HTML
		if wantsHTML(c) {
			c.HTML(http.StatusBadRequest, "verify.html", gin.H{
				"success": false,
				"error":   "Verification token required.",
			})
			return
		}
		respondBadRequest(c, ErrInvalidVerifyToken, "verification token required")
		return
	}

	// Find user by verification token
	user, err := r.db.GetUserByVerifyToken(c.Request.Context(), token)
	if err != nil {
		if wantsHTML(c) {
			c.HTML(http.StatusInternalServerError, "verify.html", gin.H{
				"success": false,
				"error":   "An error occurred. Please try again.",
			})
			return
		}
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if user == nil {
		if wantsHTML(c) {
			c.HTML(http.StatusBadRequest, "verify.html", gin.H{
				"success": false,
				"error":   "Invalid or expired verification token.",
			})
			return
		}
		respondBadRequest(c, ErrInvalidVerifyToken, "invalid or expired verification token")
		return
	}

	// Mark email as verified
	err = r.db.VerifyUserEmail(c.Request.Context(), user.ID)
	if err != nil {
		if wantsHTML(c) {
			c.HTML(http.StatusInternalServerError, "verify.html", gin.H{
				"success": false,
				"error":   "Failed to verify email. Please try again.",
			})
			return
		}
		respondInternalError(c, ErrEmailVerifyFailed, "failed to verify email")
		return
	}

	// Return success
	if wantsHTML(c) {
		c.HTML(http.StatusOK, "verify.html", gin.H{
			"success": true,
			"site_id": user.SiteID,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully. You can now log in.",
	})
}

// publicOAuthLogin initiates OAuth flow for public users.
func (r *Router) publicOAuthLogin(c *gin.Context) {
	provider := c.Param("provider")
	siteID := c.Query("site_id")

	if siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}

	// Validate site exists
	site, err := r.db.GetSiteByID(c.Request.Context(), siteID)
	if err != nil || site == nil {
		respondBadRequest(c, ErrSiteInvalid, "invalid site")
		return
	}

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
		respondBadRequest(c, ErrUnknownProvider, "unknown provider")
		return
	}

	// Generate state token with site_id embedded
	state, err := auth.GenerateSessionToken()
	if err != nil {
		respondInternalError(c, ErrTokenGenerationFailed, "failed to generate state")
		return
	}

	// Store state and site_id in cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("oauth_state", state, 600, "/", "", r.config.Auth.SecureCookies, true)
	c.SetCookie("oauth_site_id", siteID, 600, "/", "", r.config.Auth.SecureCookies, true)

	// Store intent cookie for newsletter subscribe flow
	intent := c.Query("intent")
	if intent == "subscribe" && r.config.Auth.AllowNewsletterOAuth {
		c.SetCookie("oauth_intent", "subscribe", 600, "/", "", r.config.Auth.SecureCookies, true)
	}

	// Redirect to provider
	authURL := providerCfg.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// publicOAuthCallback handles OAuth callback for public users (creates accounts).
func (r *Router) publicOAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")

	// Verify state
	savedState, err := c.Cookie("oauth_state")
	if err != nil || savedState != state {
		respondBadRequest(c, ErrInvalidState, "invalid state")
		return
	}

	// Get site_id from cookie
	siteID, err := c.Cookie("oauth_site_id")
	if err != nil || siteID == "" {
		respondBadRequest(c, ErrMissingSiteContext, "missing site context")
		return
	}

	// Read and clear intent cookie
	intent, _ := c.Cookie("oauth_intent")

	// Clear cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("oauth_state", "", -1, "/", "", r.config.Auth.SecureCookies, true)
	c.SetCookie("oauth_site_id", "", -1, "/", "", r.config.Auth.SecureCookies, true)
	c.SetCookie("oauth_intent", "", -1, "/", "", r.config.Auth.SecureCookies, true)

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
		respondBadRequest(c, ErrUnknownProvider, "unknown provider")
		return
	}

	// Exchange code for user info
	userInfo, err := providerCfg.HandleCallback(c.Request.Context(), code)
	if err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	// Branch on intent: subscribe creates a verified subscriber, not a user account
	if intent == "subscribe" && r.config.Auth.AllowNewsletterOAuth {
		r.handleOAuthNewsletterSubscribe(c, siteID, provider, userInfo)
		return
	}

	// Default flow: find or create user account
	user, err := r.db.GetUserByOAuth(c.Request.Context(), provider, userInfo.ProviderID)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	if user == nil {
		// Try to find by email
		user, err = r.db.GetUserByEmail(c.Request.Context(), siteID, userInfo.Email)
		if err != nil {
			respondInternalError(c, ErrDatabaseError, "database error")
			return
		}

		if user == nil {
			// Create new user with OAuth
			userID := uuid.New().String()
			name := userInfo.Name
			if name == "" {
				name = strings.Split(userInfo.Email, "@")[0]
			}

			user, err = r.db.CreateUserWithOAuth(c.Request.Context(), userID, siteID, userInfo.Email, provider, userInfo.ProviderID, name, userInfo.AvatarURL, "viewer")
			if err != nil {
				respondInternalError(c, ErrUserCreationFailed, "failed to create user")
				return
			}
		} else {
			// Link OAuth to existing account
			err = r.db.LinkOAuthToUser(c.Request.Context(), user.ID, provider, userInfo.ProviderID)
			if err != nil {
				respondInternalError(c, ErrOAuthLinkFailed, "failed to link OAuth")
				return
			}
		}
	}

	// Generate token
	accessToken, err := auth.GenerateToken(user.ID, user.Email, user.Role, user.SiteID, r.config.Auth.JWTSecret, r.config.Auth.JWTExpiry)
	if err != nil {
		respondInternalError(c, ErrTokenGenerationFailed, "failed to generate token")
		return
	}

	// Update last login
	_ = r.db.UpdateUserLastLogin(c.Request.Context(), user.ID)

	// Set cookie and redirect
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(r.config.Auth.SessionCookieKey, accessToken, int(r.config.Auth.JWTExpiry.Seconds()), "/", "", r.config.Auth.SecureCookies, true)

	// Redirect to forum or home (could use a redirect_uri cookie for flexibility)
	c.Redirect(http.StatusTemporaryRedirect, "/forum/")
}

// handleOAuthNewsletterSubscribe creates a verified subscriber via OAuth.
// No session cookie is set — the user is subscribing, not logging in.
func (r *Router) handleOAuthNewsletterSubscribe(c *gin.Context, siteID, provider string, userInfo *auth.UserInfo) {
	displayName := userInfo.Name
	if displayName == "" {
		displayName = strings.Split(userInfo.Email, "@")[0]
	}

	subscriberID := uuid.New().String()
	err := r.db.CreateVerifiedSubscriber(c.Request.Context(), subscriberID, siteID, userInfo.Email, provider, displayName)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "newsletter-subscribe-result.html", gin.H{
			"success": false,
			"error":   "Failed to subscribe. Please try again.",
		})
		return
	}

	c.HTML(http.StatusOK, "newsletter-subscribe-result.html", gin.H{
		"success":  true,
		"provider": provider,
		"email":    userInfo.Email,
	})
}

// Public Auth UI Handlers

// publicLoginPage renders the public login page.
func (r *Router) publicLoginPage(c *gin.Context) {
	siteID := c.Query("site_id")
	if siteID == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Site ID required. Add ?site_id=... to the URL.",
		})
		return
	}

	// Get site info
	site, err := r.db.GetSiteByID(c.Request.Context(), siteID)
	if err != nil || site == nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Invalid site.",
		})
		return
	}

	// Get OAuth providers if enabled
	var providers []string
	if r.config.Auth.EnableOAuth {
		for _, p := range r.config.OAuth.Providers {
			providers = append(providers, p.Name)
		}
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"site_id":             siteID,
		"site_name":           site.Name,
		"oauth_providers":     providers,
		"enable_registration": r.config.Auth.EnableLocal,
		"message":             c.Query("message"),
		"redirect":            c.Query("redirect"),
	})
}

// publicLoginSubmit handles the login form submission.
func (r *Router) publicLoginSubmit(c *gin.Context) {
	siteID := c.PostForm("site_id")
	emailAddr := c.PostForm("email")
	password := c.PostForm("password")

	if siteID == "" || emailAddr == "" || password == "" {
		r.renderLoginError(c, siteID, "All fields are required.")
		return
	}

	// Get user
	user, err := r.db.GetUserByEmail(c.Request.Context(), siteID, emailAddr)
	if err != nil {
		r.renderLoginError(c, siteID, "An error occurred. Please try again.")
		return
	}
	if user == nil {
		r.renderLoginError(c, siteID, "Invalid email or password.")
		return
	}

	// Verify password
	passwordHash := user.PasswordHash.String
	if !user.PasswordHash.Valid {
		passwordHash = "$2a$12$000000000000000000000.0000000000000000000000000000000"
	}
	if !auth.VerifyPassword(password, passwordHash) {
		r.renderLoginError(c, siteID, "Invalid email or password.")
		return
	}

	// Check if email is verified
	if !user.EmailVerified {
		r.renderLoginError(c, siteID, "Please verify your email before logging in.")
		return
	}

	// Generate token
	accessToken, err := auth.GenerateToken(user.ID, user.Email, user.Role, user.SiteID, r.config.Auth.JWTSecret, r.config.Auth.JWTExpiry)
	if err != nil {
		r.renderLoginError(c, siteID, "An error occurred. Please try again.")
		return
	}

	// Update last login
	_ = r.db.UpdateUserLastLogin(c.Request.Context(), user.ID)

	// Set cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(r.config.Auth.SessionCookieKey, accessToken, int(r.config.Auth.JWTExpiry.Seconds()), "/", "", r.config.Auth.SecureCookies, true)

	// Redirect to specified URL or forum
	redirectURL := c.PostForm("redirect")
	if redirectURL == "" {
		redirectURL = "/forum/"
	}
	c.Redirect(http.StatusFound, redirectURL)
}

func (r *Router) renderLoginError(c *gin.Context, siteID, errMsg string) {
	site, _ := r.db.GetSiteByID(c.Request.Context(), siteID)
	siteName := "Site"
	if site != nil {
		siteName = site.Name
	}

	var providers []string
	if r.config.Auth.EnableOAuth {
		for _, p := range r.config.OAuth.Providers {
			providers = append(providers, p.Name)
		}
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"site_id":             siteID,
		"site_name":           siteName,
		"oauth_providers":     providers,
		"enable_registration": r.config.Auth.EnableLocal,
		"error":               errMsg,
	})
}

// publicRegisterPage renders the public registration page.
func (r *Router) publicRegisterPage(c *gin.Context) {
	siteID := c.Query("site_id")
	if siteID == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Site ID required. Add ?site_id=... to the URL.",
		})
		return
	}

	// Get site info
	site, err := r.db.GetSiteByID(c.Request.Context(), siteID)
	if err != nil || site == nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Invalid site.",
		})
		return
	}

	// Get OAuth providers if enabled
	var providers []string
	if r.config.Auth.EnableOAuth {
		for _, p := range r.config.OAuth.Providers {
			providers = append(providers, p.Name)
		}
	}

	c.HTML(http.StatusOK, "register.html", gin.H{
		"site_id":         siteID,
		"site_name":       site.Name,
		"oauth_providers": providers,
	})
}

// publicRegisterSubmit handles the registration form submission.
func (r *Router) publicRegisterSubmit(c *gin.Context) {
	siteID := c.PostForm("site_id")
	name := c.PostForm("name")
	emailAddr := c.PostForm("email")
	password := c.PostForm("password")

	if siteID == "" || emailAddr == "" || password == "" {
		r.renderRegisterError(c, siteID, "Email and password are required.")
		return
	}

	if len(password) < 8 {
		r.renderRegisterError(c, siteID, "Password must be at least 8 characters.")
		return
	}

	// Validate site
	site, err := r.db.GetSiteByID(c.Request.Context(), siteID)
	if err != nil || site == nil {
		r.renderRegisterError(c, siteID, "Invalid site.")
		return
	}

	// Check if email exists
	existing, err := r.db.GetUserByEmail(c.Request.Context(), siteID, emailAddr)
	if err != nil {
		r.renderRegisterError(c, siteID, "An error occurred. Please try again.")
		return
	}
	if existing != nil {
		r.renderRegisterError(c, siteID, "This email is already registered.")
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password, r.config.Auth.BCryptCost)
	if err != nil {
		r.renderRegisterError(c, siteID, "An error occurred. Please try again.")
		return
	}

	// Generate verification token
	verifyToken, err := auth.GenerateSessionToken()
	if err != nil {
		r.renderRegisterError(c, siteID, "An error occurred. Please try again.")
		return
	}

	// Create user
	userID := uuid.New().String()
	if name == "" {
		name = strings.Split(emailAddr, "@")[0]
	}

	_, err = r.db.CreateUserWithVerification(c.Request.Context(), userID, siteID, emailAddr, hashedPassword, "viewer", name, verifyToken)
	if err != nil {
		r.renderRegisterError(c, siteID, "Failed to create account. Please try again.")
		return
	}

	// Send verification email
	if r.email != nil {
		verifyURL := r.config.Email.BaseURL + r.config.Server.APIPath + "/auth/verify?token=" + verifyToken
		msg := &email.Message{
			To:      emailAddr,
			Subject: "Verify your email address",
			HTMLBody: `<h2>Welcome!</h2>
<p>Hi ` + name + `,</p>
<p>Please verify your email address by clicking the link below:</p>
<p><a href="` + verifyURL + `">Verify Email</a></p>
<p>Or copy this URL: ` + verifyURL + `</p>
<p>This link will expire in 24 hours.</p>`,
			TextBody: "Hi " + name + ",\n\nPlease verify your email address by visiting:\n" + verifyURL + "\n\nThis link will expire in 24 hours.",
		}
		_ = r.email.Send(c.Request.Context(), msg)
	}

	// Redirect to login with success message
	c.Redirect(http.StatusFound, "/login?site_id="+siteID+"&message=Registration+successful.+Please+check+your+email+to+verify+your+account.")
}

func (r *Router) renderRegisterError(c *gin.Context, siteID, errMsg string) {
	site, _ := r.db.GetSiteByID(c.Request.Context(), siteID)
	siteName := "Site"
	if site != nil {
		siteName = site.Name
	}

	var providers []string
	if r.config.Auth.EnableOAuth {
		for _, p := range r.config.OAuth.Providers {
			providers = append(providers, p.Name)
		}
	}

	c.HTML(http.StatusOK, "register.html", gin.H{
		"site_id":         siteID,
		"site_name":       siteName,
		"oauth_providers": providers,
		"error":           errMsg,
	})
}

// publicLogout handles logout and redirects.
func (r *Router) publicLogout(c *gin.Context) {
	siteID := c.Query("site_id")

	// Clear cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(r.config.Auth.SessionCookieKey, "", -1, "/", "", r.config.Auth.SecureCookies, true)

	// Redirect to login or home
	if siteID != "" {
		c.Redirect(http.StatusFound, "/login?site_id="+siteID)
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}

// publicBlogPage serves the public blog page.
func (r *Router) publicBlogPage(c *gin.Context) {
	c.HTML(http.StatusOK, "blog.html", r.getPublicPageData(c, "blog"))
}

// publicForumPage serves the public forum page.
func (r *Router) publicForumPage(c *gin.Context) {
	c.HTML(http.StatusOK, "forum.html", r.getPublicPageData(c, "forum"))
}

// publicForumNewTopicPage serves the new topic form page.
func (r *Router) publicForumNewTopicPage(c *gin.Context) {
	data := r.getPublicPageData(c, "forum")
	// Require authentication
	if data["authenticated"] != true {
		siteID, _ := data["site_id"].(string)
		c.Redirect(http.StatusFound, "/login?site_id="+siteID)
		return
	}
	// Load categories server-side
	siteID, _ := data["site_id"].(string)
	categories, _ := r.db.ListForumCategories(c.Request.Context(), siteID)
	data["categories"] = categories
	c.HTML(http.StatusOK, "forum-new.html", data)
}

// publicForumTopicPage serves a single topic view page.
func (r *Router) publicForumTopicPage(c *gin.Context) {
	data := r.getPublicPageData(c, "forum")
	data["slug"] = c.Param("slug")
	c.HTML(http.StatusOK, "forum-topic.html", data)
}

// publicBlogArticlePage serves a single blog article page.
func (r *Router) publicBlogArticlePage(c *gin.Context) {
	data := r.getPublicPageData(c, "blog")
	data["slug"] = c.Param("slug")
	c.HTML(http.StatusOK, "blog-article.html", data)
}

// publicForumCategoryPage serves a forum category page.
func (r *Router) publicForumCategoryPage(c *gin.Context) {
	data := r.getPublicPageData(c, "forum")
	data["slug"] = c.Param("slug")
	c.HTML(http.StatusOK, "forum-category.html", data)
}
