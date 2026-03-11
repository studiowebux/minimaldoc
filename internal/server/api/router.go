// Package api provides the HTTP API handlers for minimaldoc-server.
package api

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/config"
	"github.com/studiowebux/minimaldoc/internal/server/email"
	"github.com/studiowebux/minimaldoc/internal/server/ratelimit"
	"github.com/studiowebux/minimaldoc/internal/server/storage"
	"github.com/studiowebux/minimaldoc/internal/server/store"
)

// Router wraps the Gin engine with dependencies.
type Router struct {
	*gin.Engine
	config  *config.Config
	db      *store.DB
	email   email.Sender
	storage storage.Storage

	// Rate limiters
	loginLimiter  *ratelimit.Limiter
	apiLimiter    *ratelimit.Limiter
	submitLimiter *ratelimit.Limiter
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

	// Initialize rate limiters if enabled
	if cfg.RateLimit.Enabled {
		r.apiLimiter = ratelimit.New(cfg.RateLimit.APILimit, cfg.RateLimit.APIWindow)
		r.submitLimiter = ratelimit.New(cfg.RateLimit.SubmitLimit, cfg.RateLimit.SubmitWindow)
	}

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware())
	r.Use(SecurityHeadersMiddleware())
	r.Use(CORSMiddleware(cfg.Server.CORSOrigins))

	// Health check
	r.GET("/health", r.healthCheck)

	// Public API routes
	api := r.Group(cfg.Server.APIPath)

	// Apply API rate limiting to all public routes
	if r.apiLimiter != nil {
		api.Use(r.apiLimiter.Middleware(ratelimit.IPKeyFunc))
	}

	{
		api.GET("/health", r.healthCheck)

		// Analytics - public tracking only (uses general API limit)
		api.POST("/analytics/track", r.trackPageView)
		api.POST("/analytics/duration", r.trackDuration)
		api.POST("/analytics/event", r.trackEvent)

		// Feedback - stricter rate limit for submissions
		if r.submitLimiter != nil {
			api.POST("/feedback", r.submitLimiter.Middleware(ratelimit.IPKeyFunc), r.submitFeedback)
		} else {
			api.POST("/feedback", r.submitFeedback)
		}

		// Newsletter - stricter rate limit for subscriptions
		newsletter := api.Group("/newsletter")
		{
			if r.submitLimiter != nil {
				newsletter.POST("/subscribe", r.submitLimiter.Middleware(ratelimit.IPKeyFunc), r.subscribe)
			} else {
				newsletter.POST("/subscribe", r.subscribe)
			}
			newsletter.GET("/verify", r.verifySubscription)
			newsletter.POST("/unsubscribe", r.unsubscribe)
		}

		// Blog - public reading
		blog := api.Group("/blog")
		{
			blog.GET("/posts", r.listPublishedPosts)
			blog.GET("/posts/:slug", r.getPublishedPost)
			blog.GET("/posts/:slug/comments", r.listApprovedCommentsPublic)
			// Stricter rate limit for comment submissions
			if r.submitLimiter != nil {
				blog.POST("/posts/:slug/comments", r.submitLimiter.Middleware(ratelimit.IPKeyFunc), r.submitCommentPublic)
			} else {
				blog.POST("/posts/:slug/comments", r.submitCommentPublic)
			}
			blog.GET("/posts/:slug/meta", r.getPostMeta)
			blog.GET("/posts/:slug/related", r.getRelatedPosts)
			blog.GET("/feed.xml", r.getRSSFeed)
		}

		// Private docs - public access check (optionally authenticated)
		docs := api.Group("/docs")
		{
			docs.GET("/check", r.checkDocAccess)
			// Content endpoint requires auth via middleware
			docs.GET("/content/*path", AuthMiddleware(cfg), r.getDocContent)
		}

		// Sitemap with blog posts
		api.GET("/sitemap.xml", r.getSitemap)
	}

	// Client JS/CSS for static sites integration
	r.StaticFile("/minimaldoc.js", "web/client/minimaldoc-client.js")
	r.StaticFile("/minimaldoc.css", "web/client/minimaldoc-client.css")

	return r
}

// NewAdminRouter creates the admin API and UI router.
// This handles: authentication, dashboard, management APIs.
// Should be on a separate port, not publicly accessible.
func NewAdminRouter(cfg *config.Config, db *store.DB, emailSender email.Sender, store storage.Storage) *Router {
	gin.SetMode(gin.ReleaseMode)

	r := &Router{
		Engine:  gin.New(),
		config:  cfg,
		db:      db,
		email:   emailSender,
		storage: store,
	}

	// Initialize rate limiters if enabled
	if cfg.RateLimit.Enabled {
		r.loginLimiter = ratelimit.New(cfg.RateLimit.LoginLimit, cfg.RateLimit.LoginWindow)
	}

	// Load templates
	tmpl := template.Must(template.ParseGlob("web/admin/templates/*.html"))
	r.SetHTMLTemplate(tmpl)

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware())
	r.Use(SecurityHeadersMiddleware())

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
			// Apply strict rate limiting to login
			if r.loginLimiter != nil {
				auth.POST("/login", r.loginLimiter.Middleware(ratelimit.IPKeyFunc), r.login)
			} else {
				auth.POST("/login", r.login)
			}
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

		// Blog management
		blog := api.Group("/blog")
		blog.Use(AuthMiddleware(cfg))
		{
			blog.GET("/posts", r.listAllPosts)
			blog.POST("/posts", AuthorOrAboveMiddleware(), r.createPost)
			blog.GET("/posts/:id", r.getPost)
			blog.PUT("/posts/:id", AuthorOrAboveMiddleware(), r.updatePost)
			blog.DELETE("/posts/:id", AdminMiddleware(), r.deletePost)
			blog.POST("/posts/:id/publish", AuthorOrAboveMiddleware(), r.publishPost)
			blog.POST("/posts/:id/unpublish", AuthorOrAboveMiddleware(), r.unpublishPost)
			blog.POST("/posts/:id/schedule", AuthorOrAboveMiddleware(), r.schedulePost)
			blog.POST("/posts/:id/unschedule", AuthorOrAboveMiddleware(), r.unschedulePost)

			// Comment management
			blog.GET("/comments", EditorOrAboveMiddleware(), r.listCommentsAdmin)
			blog.PUT("/comments/:id/approve", EditorOrAboveMiddleware(), r.approveComment)
			blog.PUT("/comments/:id/reject", EditorOrAboveMiddleware(), r.rejectComment)
			blog.PUT("/comments/:id/spam", EditorOrAboveMiddleware(), r.markSpam)
			blog.DELETE("/comments/:id", AdminMiddleware(), r.deleteComment)
		}

		// Doc access management (admin only)
		docs := api.Group("/docs")
		docs.Use(AuthMiddleware(cfg), AdminMiddleware())
		{
			docs.GET("/rules", r.listDocAccessRules)
			docs.POST("/rules", r.createDocAccessRule)
			docs.PUT("/rules/:id", r.updateDocAccessRule)
			docs.DELETE("/rules/:id", r.deleteDocAccessRule)
		}

		// Upload management
		uploads := api.Group("/uploads")
		uploads.Use(AuthMiddleware(cfg), AuthorOrAboveMiddleware())
		{
			uploads.POST("", r.uploadImage)
			uploads.GET("", r.listUploads)
			uploads.DELETE("/:id", r.deleteImage)
		}
	}

	// Admin UI routes
	admin := r.Group(cfg.Server.AdminPath)
	{
		admin.GET("", AdminUIAuthMiddleware(cfg), r.adminDashboard)
		admin.GET("/login", r.adminLogin)
		// Apply login rate limiting to UI login
		if r.loginLimiter != nil {
			admin.POST("/login", r.loginLimiter.Middleware(ratelimit.IPKeyFunc), r.adminLoginPost)
		} else {
			admin.POST("/login", r.adminLoginPost)
		}
		admin.GET("/logout", r.adminLogout)
		admin.GET("/analytics", AdminUIAuthMiddleware(cfg), r.adminAnalytics)
		admin.GET("/feedback", AdminUIAuthMiddleware(cfg), r.adminFeedback)
		admin.GET("/subscribers", AdminUIAuthMiddleware(cfg), r.adminSubscribers)
		admin.GET("/settings", AdminUIAuthMiddleware(cfg), AdminMiddleware(), r.adminSettings)

		// Blog UI routes
		admin.GET("/blog", AdminUIAuthMiddleware(cfg), r.adminBlog)
		admin.GET("/blog/new", AdminUIAuthMiddleware(cfg), AuthorOrAboveMiddleware(), r.adminBlogEditor)
		admin.GET("/blog/edit/:id", AdminUIAuthMiddleware(cfg), AuthorOrAboveMiddleware(), r.adminBlogEditor)
		admin.GET("/comments", AdminUIAuthMiddleware(cfg), EditorOrAboveMiddleware(), r.adminComments)

		// Doc access UI routes
		admin.GET("/doc-access", AdminUIAuthMiddleware(cfg), AdminMiddleware(), r.adminDocAccess)

		// User management UI routes (admin only)
		admin.GET("/users", AdminUIAuthMiddleware(cfg), AdminMiddleware(), r.adminUsers)

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
			fragments.GET("/event-stats", r.fragmentEventStats)
			fragments.GET("/events-by-name", r.fragmentEventsByName)
			fragments.GET("/recent-events", r.fragmentRecentEvents)
			fragments.GET("/feedback-stats", r.fragmentFeedbackStats)
			fragments.GET("/feedback-list", r.fragmentFeedbackList)
			fragments.GET("/subscriber-stats", r.fragmentSubscriberStats)
			fragments.GET("/subscriber-list", r.fragmentSubscriberList)
			fragments.GET("/site-info", AdminMiddleware(), r.fragmentSiteInfo)

			// Blog fragments
			fragments.GET("/blog-stats", r.fragmentBlogStats)
			fragments.GET("/blog-post-list", r.fragmentBlogPostList)
			fragments.GET("/blog-post-editor", r.fragmentBlogPostEditor)
			fragments.GET("/blog-post-editor/:id", r.fragmentBlogPostEditor)
			fragments.POST("/blog-post-preview", r.fragmentBlogPostPreview)
			fragments.GET("/comment-list", r.fragmentCommentList)
			fragments.GET("/comment-moderation", r.fragmentCommentModerationQueue)

			// Doc access fragments
			fragments.GET("/doc-access-list", AdminMiddleware(), r.fragmentDocAccessList)
			fragments.GET("/doc-access-form", AdminMiddleware(), r.fragmentDocAccessForm)
			fragments.GET("/doc-access-form/:id", AdminMiddleware(), r.fragmentDocAccessForm)

			// User management fragments
			fragments.GET("/user-stats", AdminMiddleware(), r.fragmentUserStats)
			fragments.GET("/user-list", AdminMiddleware(), r.fragmentUserList)
			fragments.GET("/user-form", AdminMiddleware(), r.fragmentUserForm)
			fragments.GET("/user-form/:id", AdminMiddleware(), r.fragmentUserForm)
		}
	}

	// Static files for admin UI
	r.Static("/static", "web/admin/static")

	// Serve uploaded files (local storage)
	if cfg.Storage.Provider == "local" {
		r.Static("/uploads", cfg.Storage.LocalPath)
	}

	return r
}

// Stop gracefully stops background tasks like rate limiter cleanup.
func (r *Router) Stop() {
	if r.loginLimiter != nil {
		r.loginLimiter.Stop()
	}
	if r.apiLimiter != nil {
		r.apiLimiter.Stop()
	}
	if r.submitLimiter != nil {
		r.submitLimiter.Stop()
	}
}

// healthCheck returns server health status.
func (r *Router) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"database": r.db.Driver(),
	})
}
