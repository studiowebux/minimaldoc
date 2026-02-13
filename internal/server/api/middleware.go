package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
	"github.com/studiowebux/minimaldoc/internal/server/config"
)

// LoggerMiddleware logs HTTP requests.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method

		log.Printf("%s %s %d %v", method, path, status, latency)
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing.
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, o := range allowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-API-Key")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// AuthMiddleware validates JWT tokens or session cookies.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// Check API key header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" {
			// TODO: Validate API key from database
			// For now, just continue
			c.Next()
			return
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

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("site_id", claims.SiteID)

		c.Next()
	}
}

// AdminMiddleware ensures the user has admin role.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AdminUIAuthMiddleware validates JWT for admin UI routes.
// Redirects to login page instead of returning JSON error.
func AdminUIAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
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

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("site_id", claims.SiteID)
		c.Set("user", claims)

		c.Next()
	}
}

// SiteMiddleware validates and sets site context from API key or user.
func SiteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Site can be determined from:
		// 1. X-API-Key header (for client-side tracking)
		// 2. User's associated site (for admin operations)

		siteID, exists := c.Get("site_id")
		if !exists {
			// Try to get from API key
			// TODO: Implement API key lookup
			c.JSON(http.StatusBadRequest, gin.H{"error": "site context required"})
			c.Abort()
			return
		}

		c.Set("site_id", siteID)
		c.Next()
	}
}
