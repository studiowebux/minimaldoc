// Package api provides the HTTP API handlers for minimaldoc-server.
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/config"
	"github.com/studiowebux/minimaldoc/internal/server/store"
)

// Router wraps the Gin engine with dependencies.
type Router struct {
	*gin.Engine
	config *config.Config
	db     *store.DB
}

// NewRouter creates a new API router.
func NewRouter(cfg *config.Config, db *store.DB) *Router {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	r := &Router{
		Engine: gin.New(),
		config: cfg,
		db:     db,
	}

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware())
	r.Use(CORSMiddleware(cfg.Server.CORSOrigins))

	// Health check
	r.GET("/health", r.healthCheck)

	// API routes
	api := r.Group(cfg.Server.APIPath)
	{
		// Public endpoints
		api.GET("/health", r.healthCheck)

		// Auth endpoints
		auth := api.Group("/auth")
		{
			auth.POST("/login", r.login)
			auth.POST("/logout", r.logout)
			auth.POST("/refresh", r.refreshToken)
			auth.GET("/me", AuthMiddleware(cfg), r.getCurrentUser)

			// OAuth endpoints
			if cfg.Auth.EnableOAuth {
				auth.GET("/providers", r.listOAuthProviders)
				auth.GET("/oauth/:provider/login", r.oauthLogin)
				auth.GET("/oauth/:provider/callback", r.oauthCallback)
			}
		}

		// Analytics endpoints
		analytics := api.Group("/analytics")
		{
			analytics.POST("/track", r.trackPageView)
			analytics.POST("/event", r.trackEvent)

			// Admin-only analytics
			analytics.GET("/summary", AuthMiddleware(cfg), r.analyticsSummary)
			analytics.GET("/pages", AuthMiddleware(cfg), r.analyticsPages)
		}

		// Feedback endpoints
		feedback := api.Group("/feedback")
		{
			feedback.POST("", r.submitFeedback)
			feedback.GET("/stats", AuthMiddleware(cfg), r.feedbackStats)
			feedback.GET("/list", AuthMiddleware(cfg), r.feedbackList)
		}

		// Newsletter endpoints
		newsletter := api.Group("/newsletter")
		{
			newsletter.POST("/subscribe", r.subscribe)
			newsletter.GET("/verify", r.verifySubscription)
			newsletter.POST("/unsubscribe", r.unsubscribe)
			newsletter.GET("/subscribers", AuthMiddleware(cfg), r.listSubscribers)
		}

		// Site management (admin only)
		sites := api.Group("/sites")
		sites.Use(AuthMiddleware(cfg), AdminMiddleware())
		{
			sites.GET("", r.listSites)
			sites.POST("", r.createSite)
			sites.GET("/:id", r.getSite)
			sites.PUT("/:id", r.updateSite)
			sites.DELETE("/:id", r.deleteSite)
			sites.POST("/:id/api-key", r.regenerateAPIKey)
		}

		// User management (admin only)
		users := api.Group("/users")
		users.Use(AuthMiddleware(cfg), AdminMiddleware())
		{
			users.GET("", r.listUsers)
			users.POST("", r.createUser)
			users.GET("/:id", r.getUser)
			users.PUT("/:id", r.updateUser)
			users.DELETE("/:id", r.deleteUser)
		}
	}

	// Admin UI routes
	admin := r.Group(cfg.Server.AdminPath)
	{
		admin.GET("", r.adminDashboard)
		admin.GET("/login", r.adminLogin)
		admin.POST("/login", r.adminLoginPost)
		admin.GET("/logout", r.adminLogout)
		admin.GET("/analytics", AuthMiddleware(cfg), r.adminAnalytics)
		admin.GET("/feedback", AuthMiddleware(cfg), r.adminFeedback)
		admin.GET("/subscribers", AuthMiddleware(cfg), r.adminSubscribers)
		admin.GET("/settings", AuthMiddleware(cfg), r.adminSettings)
	}

	// Static files for admin UI
	r.Static("/static", "web/admin/static")

	return r
}

// healthCheck returns server health status.
func (r *Router) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"database": r.db.Driver(),
	})
}
