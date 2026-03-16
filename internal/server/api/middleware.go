package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
	"github.com/studiowebux/minimaldoc/internal/server/config"
	"github.com/studiowebux/minimaldoc/internal/server/store"
)

// TracingMiddleware creates a span per HTTP request with method, path, and status attributes.
func TracingMiddleware() gin.HandlerFunc {
	tracer := otel.Tracer("minimaldoc-server")

	return func(c *gin.Context) {
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
		ctx, span := tracer.Start(c.Request.Context(), spanName,
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", c.FullPath()),
			attribute.Int("http.status_code", c.Writer.Status()),
		)
	}
}

// generateNonce creates a 16-byte base64-encoded cryptographic nonce for CSP.
func generateNonce() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

// SecurityHeadersMiddleware adds common security headers to responses.
// Generates a per-request nonce for CSP and stores it in gin context as "csp_nonce".
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate per-request nonce for inline scripts/styles
		nonce := generateNonce()
		c.Set("csp_nonce", nonce)

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		// Enable XSS filter in browsers
		c.Header("X-XSS-Protection", "1; mode=block")
		// Referrer policy - don't leak full URL
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// Permissions policy - disable unnecessary features
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		// Content Security Policy with nonce — no unsafe-inline
		csp := fmt.Sprintf(
			"default-src 'self'; script-src 'self' 'nonce-%s'; style-src 'self' 'nonce-%s'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'",
			nonce, nonce,
		)
		c.Header("Content-Security-Policy", csp)

		c.Next()
	}
}

// RequestIDMiddleware generates or propagates a unique request ID.
// Reads X-Request-ID from the incoming request (set by reverse proxy);
// generates a new UUID if absent. Stores in gin context and response header.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set("request_id", reqID)
		c.Header("X-Request-ID", reqID)
		c.Next()
	}
}

// LoggerMiddleware logs HTTP requests with request ID correlation.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		reqID, _ := c.Get("request_id")

		slog.Info("request", "method", method, "path", path, "status", status, "latency", latency, "request_id", reqID)
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing.
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	// Check if wildcard is configured
	hasWildcard := false
	for _, o := range allowedOrigins {
		if o == "*" {
			hasWildcard = true
			break
		}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			c.Next()
			return
		}

		// Check if origin is allowed
		allowed := false
		for _, o := range allowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if hasWildcard {
				// With wildcard, use "*" and don't allow credentials
				// This is the safe configuration for public APIs
				c.Header("Access-Control-Allow-Origin", "*")
				// Don't set Allow-Credentials with wildcard
			} else {
				// With explicit origins, echo back the origin and allow credentials
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-API-Key")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// AuthMiddleware validates JWT tokens, session cookies, or API keys.
// API keys are validated against the database; valid keys set site_id in context.
func AuthMiddleware(cfg *config.Config, db store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check API key header first
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" {
			keyHash := auth.HashAPIKey(apiKey)
			site, err := db.GetSiteByAPIKey(c.Request.Context(), keyHash)
			if err != nil || site == nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
				c.Abort()
				return
			}
			c.Set("site_id", site.ID)
			c.Set("auth_method", "api_key")
			c.Next()
			return
		}

		var token string

		// Check Authorization header
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// Check session cookie
		if token == "" {
			cookie, err := c.Cookie(cfg.Auth.SessionCookieKey)
			if err == nil {
				token = cookie
			}
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		// Validate JWT
		claims, err := auth.ValidateToken(token, cfg.Auth.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Check if token has been revoked (logout invalidation)
		if claims.ID != "" {
			revoked, err := db.IsTokenRevoked(c.Request.Context(), claims.ID)
			if err != nil {
				slog.Error("failed to check token revocation", "error", err)
			}
			if revoked {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
				c.Abort()
				return
			}
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("site_id", claims.SiteID)
		c.Set("auth_method", "jwt")
		c.Set("token_jti", claims.ID)
		c.Set("token_exp", claims.ExpiresAt.Time)

		c.Next()
	}
}

// OptionalAuthMiddleware tries to authenticate the user but doesn't require it.
// If a valid token is present, it sets user info in context.
// If not, it continues without setting user info (for anonymous access).
func OptionalAuthMiddleware(cfg *config.Config, db store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// Check Authorization header
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// Check session cookie
		if token == "" {
			if sessionCookie, err := c.Cookie(cfg.Auth.SessionCookieKey); err == nil && sessionCookie != "" {
				token = sessionCookie
			}
		}

		// If no token, continue as anonymous
		if token == "" {
			c.Next()
			return
		}

		// Token present — validate it. Invalid/revoked tokens get 401, not anonymous fallback.
		claims, err := auth.ValidateToken(token, cfg.Auth.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Check if token has been revoked
		if claims.ID != "" {
			if revoked, _ := db.IsTokenRevoked(c.Request.Context(), claims.ID); revoked {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
				c.Abort()
				return
			}
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("site_id", claims.SiteID)

		c.Next()
	}
}

// wantsHTML checks if the request prefers HTML response.
func wantsHTML(c *gin.Context) bool {
	accept := c.GetHeader("Accept")
	// Check if Accept header prefers HTML or if it's a browser navigation
	return strings.Contains(accept, "text/html") ||
		(!strings.Contains(accept, "application/json") && c.Request.Method == "GET")
}

// renderForbidden renders an appropriate forbidden response based on Accept header.
func renderForbidden(c *gin.Context, message string) {
	if wantsHTML(c) {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"Title":   "Access Denied",
			"Code":    403,
			"Message": message,
		})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": message})
	}
	c.Abort()
}

// AdminMiddleware ensures the user has admin role.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != "admin" {
			renderForbidden(c, "Admin access required. Your current role does not have permission to view this page.")
			return
		}
		c.Next()
	}
}

// AdminUIAuthMiddleware validates JWT for admin UI routes.
// Redirects to login page instead of returning JSON error.
func AdminUIAuthMiddleware(cfg *config.Config, db store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// Check session cookie
		cookie, err := c.Cookie(cfg.Auth.SessionCookieKey)
		if err == nil {
			token = cookie
		}

		if token == "" {
			c.Redirect(http.StatusFound, cfg.Server.AdminPath+"/login")
			c.Abort()
			return
		}

		// Validate JWT
		claims, err := auth.ValidateToken(token, cfg.Auth.JWTSecret)
		if err != nil {
			c.Redirect(http.StatusFound, cfg.Server.AdminPath+"/login")
			c.Abort()
			return
		}

		// Check if token has been revoked
		if claims.ID != "" {
			if revoked, _ := db.IsTokenRevoked(c.Request.Context(), claims.ID); revoked {
				c.Redirect(http.StatusFound, cfg.Server.AdminPath+"/login")
				c.Abort()
				return
			}
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("site_id", claims.SiteID)
		c.Set("user", claims)

		c.Next()
	}
}

// SiteMiddleware validates and sets site context from user's JWT.
// For API key-based requests, handlers validate the key directly.
func SiteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		siteID, exists := c.Get("site_id")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "site context required"})
			c.Abort()
			return
		}

		c.Set("site_id", siteID)
		c.Next()
	}
}

// EditorOrAboveMiddleware requires admin or editor role.
func EditorOrAboveMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			renderForbidden(c, "Access denied. Authentication required.")
			return
		}
		r := role.(string)
		if r != "admin" && r != "editor" {
			renderForbidden(c, "Editor or admin access required. Your current role does not have permission to view this page.")
			return
		}
		c.Next()
	}
}

// AuthorOrAboveMiddleware requires admin, editor, or author role.
func AuthorOrAboveMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			renderForbidden(c, "Access denied. Authentication required.")
			return
		}
		r := role.(string)
		if r != "admin" && r != "editor" && r != "author" {
			renderForbidden(c, "Author or above access required. Your current role does not have permission to view this page.")
			return
		}
		c.Next()
	}
}

// CanEditPostMiddleware checks if user can edit a specific post.
// Admin/editor can edit any post, author can only edit own posts.
// Requires post_id param and author_id to be fetched before this middleware.
func CanEditPostMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := getUserRole(c)

		// Admin and editor can edit any post
		if role == "admin" || role == "editor" {
			c.Next()
			return
		}

		// Author can only edit own posts
		if role == "author" {
			userID, _ := getUserID(c)
			postAuthorID, exists := c.Get("post_author_id")
			if !exists || postAuthorID != userID {
				renderForbidden(c, "You can only edit your own posts.")
				return
			}
			c.Next()
			return
		}

		renderForbidden(c, "Access denied.")
	}
}
