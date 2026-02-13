// Package api provides the HTTP API handlers for minimaldoc-server.
package api

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/config"
	"github.com/studiowebux/minimaldoc/internal/server/email"
	"github.com/studiowebux/minimaldoc/internal/server/store"
)

// Router wraps the Gin engine with dependencies.
type Router struct {
	*gin.Engine
	config *config.Config
	db     *store.DB
	email  email.Sender
}

// NewPublicRouter creates the public-facing API router.
// This handles: tracking, feedback submission, newsletter subscribe.
// No authentication required, no admin access.
func NewPublicRouter(cfg *config.Config, db *store.DB, emailSender email.Sender) *Router {
	gin.SetMode(gin.ReleaseMode)

	r := &Router{
		Engine: gin.New(),
		config: cfg,
		db:     db,
		email:  emailSender,
	}

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware())
	r.Use(CORSMiddleware(cfg.Server.CORSOrigins))

	// Health check
	r.GET("/health", r.healthCheck)

	// Public API routes
	api := r.Group(cfg.Server.APIPath)
	{
		api.GET("/health", r.healthCheck)

		// Analytics - public tracking only
		api.POST("/analytics/track", r.trackPageView)
		api.POST("/analytics/duration", r.trackDuration)
		api.POST("/analytics/event", r.trackEvent)

		// Feedback - public submission only
		api.POST("/feedback", r.submitFeedback)

		// Newsletter - public subscription
		api.POST("/newsletter/subscribe", r.subscribe)
		api.GET("/newsletter/verify", r.verifySubscription)
		api.POST("/newsletter/unsubscribe", r.unsubscribe)
	}

	// Client JS/CSS for static sites integration
	r.StaticFile("/minimaldoc.js", "web/client/minimaldoc-client.js")
	r.StaticFile("/minimaldoc.css", "web/client/minimaldoc-client.css")

	return r
}

// NewAdminRouter creates the admin API and UI router.
// This handles: authentication, dashboard, management APIs.
// Should be on a separate port, not publicly accessible.
func NewAdminRouter(cfg *config.Config, db *store.DB, emailSender email.Sender) *Router {
	gin.SetMode(gin.ReleaseMode)

	r := &Router{
		Engine: gin.New(),
		config: cfg,
		db:     db,
		email:  emailSender,
	}

	// Load templates
	tmpl := template.Must(template.ParseGlob("web/admin/templates/*.html"))
	r.SetHTMLTemplate(tmpl)

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware())

	// Health check
	r.GET("/health", r.healthCheck)

	// Admin API routes
	api := r.Group(cfg.Server.APIPath)
	{
		api.GET("/health", r.healthCheck)
		api.POST("/bootstrap", r.bootstrap)

		// Auth endpoints
		auth := api.Group("/auth")
		{
			auth.POST("/login", r.login)
			auth.POST("/logout", r.logout)
			auth.POST("/refresh", r.refreshToken)
			auth.GET("/me", AuthMiddleware(cfg), r.getCurrentUser)

			if cfg.Auth.EnableOAuth {
				auth.GET("/providers", r.listOAuthProviders)
				auth.GET("/oauth/:provider/login", r.oauthLogin)
				auth.GET("/oauth/:provider/callback", r.oauthCallback)
			}
		}

		// Analytics - admin viewing
		analytics := api.Group("/analytics")
		analytics.Use(AuthMiddleware(cfg))
		{
			analytics.GET("/summary", r.analyticsSummary)
			analytics.GET("/pages", r.analyticsPages)
		}

		// Feedback - admin viewing
		feedback := api.Group("/feedback")
		feedback.Use(AuthMiddleware(cfg))
		{
			feedback.GET("/stats", r.feedbackStats)
			feedback.GET("/list", r.feedbackList)
		}

		// Newsletter - admin management
		newsletter := api.Group("/newsletter")
		newsletter.Use(AuthMiddleware(cfg))
		{
			newsletter.GET("/subscribers", r.listSubscribers)
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
		admin.GET("", AdminUIAuthMiddleware(cfg), r.adminDashboard)
		admin.GET("/login", r.adminLogin)
		admin.POST("/login", r.adminLoginPost)
		admin.GET("/logout", r.adminLogout)
		admin.GET("/analytics", AdminUIAuthMiddleware(cfg), r.adminAnalytics)
		admin.GET("/feedback", AdminUIAuthMiddleware(cfg), r.adminFeedback)
		admin.GET("/subscribers", AdminUIAuthMiddleware(cfg), r.adminSubscribers)
		admin.GET("/settings", AdminUIAuthMiddleware(cfg), r.adminSettings)

		// HTMX fragment endpoints
		fragments := admin.Group("/fragments")
		fragments.Use(AdminUIAuthMiddleware(cfg))
		{
			fragments.GET("/dashboard-stats", r.fragmentDashboardStats)
			fragments.GET("/recent-pages", r.fragmentRecentPages)
			fragments.GET("/recent-feedback", r.fragmentRecentFeedback)
			fragments.GET("/analytics-stats", r.fragmentAnalyticsStats)
			fragments.GET("/top-pages", r.fragmentTopPages)
			fragments.GET("/traffic-sources", r.fragmentTrafficSources)
			fragments.GET("/views-chart", r.fragmentViewsChart)
			fragments.GET("/feedback-stats", r.fragmentFeedbackStats)
			fragments.GET("/feedback-list", r.fragmentFeedbackList)
			fragments.GET("/subscriber-stats", r.fragmentSubscriberStats)
			fragments.GET("/subscriber-list", r.fragmentSubscriberList)
			fragments.GET("/site-info", r.fragmentSiteInfo)
		}
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
