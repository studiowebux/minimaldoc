// Package api provides the HTTP API handlers for minimaldoc-server.
package api

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/assets"
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
	db      store.Store
	email   email.Sender
	storage storage.Storage

	// Rate limiters
	loginLimiter  *ratelimit.Limiter
	apiLimiter    *ratelimit.Limiter
	submitLimiter *ratelimit.Limiter
}

// themeStaticCSS returns the embedded common theme CSS as an http.FileSystem.
// This eliminates duplication -- tokens.css lives only in internal/assets/themes/.
func themeStaticCSS() http.FileSystem {
	sub, err := fs.Sub(assets.ThemeFS, "themes/common/static/css")
	if err != nil {
		panic("embedded theme CSS not found: " + err.Error())
	}
	return http.FS(sub)
}

// NewPublicRouter creates the public-facing API router.
// This handles: tracking, feedback submission, newsletter subscribe.
// No authentication required, no admin access.
func NewPublicRouter(cfg *config.Config, db store.Store, emailSender email.Sender) *Router {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.RedirectTrailingSlash = true

	r := &Router{
		Engine: engine,
		config: cfg,
		db:     db,
		email:  emailSender,
	}

	// Initialize rate limiters if enabled
	if cfg.RateLimit.Enabled {
		r.loginLimiter = ratelimit.New(cfg.RateLimit.LoginLimit, cfg.RateLimit.LoginWindow)
		r.apiLimiter = ratelimit.New(cfg.RateLimit.APILimit, cfg.RateLimit.APIWindow)
		r.submitLimiter = ratelimit.New(cfg.RateLimit.SubmitLimit, cfg.RateLimit.SubmitWindow)
	}

	// Global middleware
	r.Use(gin.Recovery())
	if cfg.Telemetry.Enabled {
		r.Use(TracingMiddleware())
	}
	r.Use(LoggerMiddleware())
	r.Use(SecurityHeadersMiddleware())
	r.Use(CORSMiddleware(cfg.Server.CORSOrigins))

	// Health checks
	r.GET("/health", r.healthCheck)
	r.GET("/healthz", r.liveness)
	r.GET("/readyz", r.readiness)

	// Public API routes
	api := r.Group(cfg.Server.APIPath)

	// Apply API rate limiting to all public routes
	if r.apiLimiter != nil {
		api.Use(r.apiLimiter.Middleware(ratelimit.IPKeyFunc))
	}

	{
		api.GET("/health", r.healthCheck)

		// Public auth endpoints (for forum users, etc.)
		auth := api.Group("/auth")
		{
			// Login with rate limiting
			if r.loginLimiter != nil {
				auth.POST("/login", r.loginLimiter.Middleware(ratelimit.IPKeyFunc), r.login)
			} else {
				auth.POST("/login", r.login)
			}
			auth.POST("/logout", r.logout)
			auth.POST("/refresh", r.refreshToken)
			auth.GET("/me", AuthMiddleware(cfg), r.getCurrentUser)

			// Registration (if local auth enabled)
			if cfg.Auth.EnableLocal {
				if r.loginLimiter != nil {
					auth.POST("/register", r.loginLimiter.Middleware(ratelimit.IPKeyFunc), r.register)
				} else {
					auth.POST("/register", r.register)
				}
				auth.GET("/verify", r.verifyEmail)
			}

			// OAuth (if enabled)
			if cfg.Auth.EnableOAuth {
				auth.GET("/providers", r.listOAuthProviders)
				auth.GET("/oauth/:provider/login", r.publicOAuthLogin)
				auth.GET("/oauth/:provider/callback", r.publicOAuthCallback)
			}
		}

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
			// Stricter rate limit for comment submissions (with optional auth for logged-in users)
			if r.submitLimiter != nil {
				blog.POST("/posts/:slug/comments", OptionalAuthMiddleware(cfg), r.submitLimiter.Middleware(ratelimit.IPKeyFunc), r.submitCommentPublic)
			} else {
				blog.POST("/posts/:slug/comments", OptionalAuthMiddleware(cfg), r.submitCommentPublic)
			}
			blog.GET("/posts/:slug/meta", r.getPostMeta)
			blog.GET("/posts/:slug/related", r.getRelatedPosts)
			blog.GET("/feed.xml", r.getRSSFeed)

			// Blog HTML fragments for HTMX
			fragments := blog.Group("/fragments")
			{
				fragments.GET("/posts", r.fragmentBlogPostCards)
				fragments.GET("/search", r.fragmentBlogSearchDropdown)
				fragments.GET("/posts/:slug", r.fragmentBlogArticle)
				fragments.GET("/posts/:slug/related", r.fragmentBlogRelated)
				fragments.GET("/posts/:slug/comments", OptionalAuthMiddleware(cfg), r.fragmentBlogComments)
			}
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

		// Forum - public reading
		forum := api.Group("/forum")
		{
			forum.GET("/categories", r.listForumCategories)
			forum.GET("/categories/:slug", r.getForumCategory)
			forum.GET("/topics", r.listForumTopics)
			forum.GET("/topics/by-slug/:slug", r.getForumTopic)
			forum.GET("/topics/by-slug/:slug/posts", r.getForumTopicPosts)
			forum.GET("/search", r.searchForum)
			forum.GET("/tags", r.listForumTags)
			forum.GET("/tags/:slug/topics", r.getForumTopicsByTag)
			forum.GET("/leaderboard", r.getForumLeaderboard)

			// Authenticated forum actions
			authForum := forum.Group("")
			authForum.Use(AuthMiddleware(cfg))
			{
				authForum.POST("/topics", r.createForumTopic)
				authForum.POST("/topics/by-slug/:slug/posts", r.createForumPost)
				authForum.PUT("/topics/:id", r.updateForumTopic)
				authForum.PUT("/posts/:id", r.updateForumPost)
				authForum.DELETE("/topics/:id", r.deleteForumTopic)
				authForum.DELETE("/posts/:id", r.deleteForumPost)
				authForum.POST("/topics/:id/like", r.likeForumTopic)
				authForum.POST("/posts/:id/like", r.likeForumPost)
				authForum.POST("/topics/:id/bookmark", r.bookmarkForumTopic)
				authForum.POST("/flag", r.flagForumContent)
				authForum.GET("/notifications", r.getForumNotifications)
				authForum.PUT("/notifications/:id/read", r.markNotificationRead)
				authForum.PUT("/notifications/read-all", r.markAllNotificationsRead)
				authForum.GET("/bookmarks", r.getUserBookmarks)
			}

			// Forum HTML fragments for HTMX (with optional auth to show user-specific UI)
			fragments := forum.Group("/fragments")
			fragments.Use(OptionalAuthMiddleware(cfg))
			{
				fragments.GET("/categories", r.fragmentForumCategories)
				fragments.GET("/category/:slug/header", r.fragmentForumCategoryHeader)
				fragments.GET("/category-options", r.fragmentForumCategoryOptions)
				fragments.GET("/category-select", r.fragmentForumCategorySelect)
				fragments.GET("/topics", r.fragmentForumTopics)
				fragments.GET("/topics/:slug/detail", r.fragmentForumTopicDetail)
				fragments.GET("/search", r.fragmentForumSearchDropdown)
				fragments.GET("/leaderboard", r.fragmentForumLeaderboard)
				fragments.GET("/tags", r.fragmentForumTags)
				fragments.GET("/tag-buttons", r.fragmentForumTagButtons)
				fragments.GET("/stats", r.fragmentForumStats)
				fragments.GET("/topics/:slug/posts", r.fragmentForumPosts)
				fragments.GET("/topics/:slug/reply-form", r.fragmentForumReplyForm)
			}
		}
	}

	// Client JS/CSS for static sites integration
	r.StaticFile("/minimaldoc.js", "web/client/minimaldoc-client.js")
	r.StaticFile("/minimaldoc.css", "web/client/minimaldoc-client.css")

	// Shared theme CSS from embedded FS (single source of truth)
	r.StaticFS("/theme/css", themeStaticCSS())

	// Static files for public pages (admin.css, auth.css, public.css)
	r.Static("/static", "web/admin/static")

	// Load public templates
	tmpl := template.Must(template.ParseGlob("web/public/templates/*.html"))
	r.SetHTMLTemplate(tmpl)

	// Public blog and forum pages
	r.GET("/blog/", r.publicBlogPage)
	r.GET("/blog", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/blog/") })
	r.GET("/blog/:slug/", r.publicBlogArticlePage)
	r.GET("/blog/:slug", r.publicBlogArticlePage)
	r.GET("/forum", r.publicForumPage)
	r.GET("/forum/", r.publicForumPage)
	r.GET("/forum/new", r.publicForumNewTopicPage)
	r.GET("/forum/new/", r.publicForumNewTopicPage)
	r.GET("/forum/topic/:slug", r.publicForumTopicPage)
	r.GET("/forum/category/:slug", r.publicForumCategoryPage)

	// Public auth UI pages (if local auth enabled)
	if cfg.Auth.EnableLocal || cfg.Auth.EnableOAuth {
		// Login page
		r.GET("/login", r.publicLoginPage)
		if r.loginLimiter != nil {
			r.POST("/login", r.loginLimiter.Middleware(ratelimit.IPKeyFunc), r.publicLoginSubmit)
		} else {
			r.POST("/login", r.publicLoginSubmit)
		}

		// Registration page (only if local auth enabled)
		if cfg.Auth.EnableLocal {
			r.GET("/register", r.publicRegisterPage)
			if r.loginLimiter != nil {
				r.POST("/register", r.loginLimiter.Middleware(ratelimit.IPKeyFunc), r.publicRegisterSubmit)
			} else {
				r.POST("/register", r.publicRegisterSubmit)
			}
		}

		// Logout
		r.GET("/logout", r.publicLogout)
	}

	return r
}

// NewAdminRouter creates the admin API and UI router.
// This handles: authentication, dashboard, management APIs.
// Should be on a separate port, not publicly accessible.
func NewAdminRouter(cfg *config.Config, db store.Store, emailSender email.Sender, store storage.Storage) *Router {
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
	if cfg.Telemetry.Enabled {
		r.Use(TracingMiddleware())
	}
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

		// Forum management
		forum := api.Group("/forum")
		forum.Use(AuthMiddleware(cfg))
		{
			forum.GET("/stats", r.getForumStats)
			forum.GET("/topics", r.adminListForumTopics)
			forum.GET("/categories", r.adminListForumCategories)
			forum.GET("/categories/:id", r.adminGetForumCategory)

			// Category management (editor+)
			forum.POST("/categories", EditorOrAboveMiddleware(), r.adminCreateForumCategory)
			forum.PUT("/categories/:id", EditorOrAboveMiddleware(), r.adminUpdateForumCategory)
			forum.DELETE("/categories/:id", AdminMiddleware(), r.adminDeleteForumCategory)

			// Moderation (editor+)
			forum.POST("/topics/:id/pin", EditorOrAboveMiddleware(), r.adminPinForumTopic)
			forum.POST("/topics/:id/lock", EditorOrAboveMiddleware(), r.adminLockForumTopic)
			forum.POST("/topics/:id/close", EditorOrAboveMiddleware(), r.adminCloseForumTopic)
			forum.POST("/topics/:id/open", EditorOrAboveMiddleware(), r.adminOpenForumTopic)
			forum.POST("/posts/:id/solution", EditorOrAboveMiddleware(), r.adminMarkSolution)
			forum.DELETE("/topics/:id", EditorOrAboveMiddleware(), r.adminDeleteForumTopic)
			forum.DELETE("/posts/:id", EditorOrAboveMiddleware(), r.adminDeleteForumPost)

			// Flags (editor+)
			forum.GET("/flags", EditorOrAboveMiddleware(), r.adminListForumFlags)
			forum.PUT("/flags/:id", EditorOrAboveMiddleware(), r.adminResolveForumFlag)

			// Bans (admin only)
			forum.POST("/users/ban", AdminMiddleware(), r.adminBanUser)
			forum.DELETE("/users/:id/ban", AdminMiddleware(), r.adminUnbanUser)
			forum.GET("/bans", AdminMiddleware(), r.adminListForumBans)

			// Tags (editor+)
			forum.GET("/tags", r.adminListForumTags)
			forum.GET("/tags/:id", r.adminGetForumTag)
			forum.POST("/tags", EditorOrAboveMiddleware(), r.adminCreateForumTag)
			forum.PUT("/tags/:id", EditorOrAboveMiddleware(), r.adminUpdateForumTag)
			forum.DELETE("/tags/:id", AdminMiddleware(), r.adminDeleteForumTag)
		}

		// Upload management
		uploads := api.Group("/uploads")
		uploads.Use(AuthMiddleware(cfg), AuthorOrAboveMiddleware())
		{
			uploads.POST("", r.uploadImage)
			uploads.GET("", r.listUploads)
			uploads.DELETE("/:id", r.deleteImage)
		}

		// Audit log API (admin only)
		api.GET("/audit-logs", AuthMiddleware(cfg), AdminMiddleware(), r.listAuditLogsAPI)
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

		// Forum UI routes
		admin.GET("/forum", AdminUIAuthMiddleware(cfg), EditorOrAboveMiddleware(), r.adminForum)
		admin.GET("/forum/categories", AdminUIAuthMiddleware(cfg), EditorOrAboveMiddleware(), r.adminForumCategories)
		admin.GET("/forum/topics", AdminUIAuthMiddleware(cfg), EditorOrAboveMiddleware(), r.adminForumTopics)
		admin.GET("/forum/flags", AdminUIAuthMiddleware(cfg), EditorOrAboveMiddleware(), r.adminForumFlags)
		admin.GET("/forum/bans", AdminUIAuthMiddleware(cfg), AdminMiddleware(), r.adminForumBans)
		admin.GET("/forum/tags", AdminUIAuthMiddleware(cfg), EditorOrAboveMiddleware(), r.adminForumTags)

		// User management UI routes (admin only)
		admin.GET("/users", AdminUIAuthMiddleware(cfg), AdminMiddleware(), r.adminUsers)

		// Audit log UI routes (admin only)
		admin.GET("/audit-log", AdminUIAuthMiddleware(cfg), AdminMiddleware(), r.adminAuditLog)

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

			// Audit log fragments (admin only)
			fragments.GET("/audit-stats", AdminMiddleware(), r.fragmentAuditStats)
			fragments.GET("/audit-log-list", AdminMiddleware(), r.fragmentAuditLogList)
		}
	}

	// Shared theme CSS from embedded FS (single source of truth)
	r.StaticFS("/theme/css", themeStaticCSS())

	// Static files for admin UI (admin.css, auth.css, public.css)
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

// liveness returns 200 if the server process is running.
func (r *Router) liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

// readiness returns 200 if the server can serve requests (DB is reachable).
func (r *Router) readiness(c *gin.Context) {
	if err := r.db.Ping(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": "database unreachable", "code": ErrDBUnreachable})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
