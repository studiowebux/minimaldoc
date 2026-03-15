package api

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"

	"github.com/studiowebux/minimaldoc/internal/server/markdown"
	"github.com/studiowebux/minimaldoc/internal/server/store"
)

// commentSanitizer strips all HTML from comments — plaintext only.
var commentSanitizer = bluemonday.StrictPolicy()

// sanitizeComment strips all HTML and trims whitespace from user comment input.
func sanitizeComment(s string) string {
	s = commentSanitizer.Sanitize(s)
	s = strings.TrimSpace(s)
	return s
}

// BlogPostRequest represents the request body for creating/updating a blog post.
type BlogPostRequest struct {
	Slug          string `json:"slug" binding:"required"`
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description"`
	Content       string `json:"content" binding:"required"`
	FeaturedImage string `json:"featured_image"`
	Tags          string `json:"tags"` // JSON array string
	Category      string `json:"category"`
	Visibility    string `json:"visibility"`   // public, authenticated, role_viewer, role_author, role_editor, role_admin
	ScheduledAt   string `json:"scheduled_at"` // ISO 8601 datetime for scheduled publishing
}

// BlogPostResponse is the JSON-friendly response for a blog post.
type BlogPostResponse struct {
	ID            string `json:"id"`
	Slug          string `json:"slug"`
	Title         string `json:"title"`
	Description   string `json:"description,omitempty"`
	Content       string `json:"content,omitempty"`
	ContentHTML   string `json:"content_html,omitempty"`
	FeaturedImage string `json:"featured_image,omitempty"`
	Tags          string `json:"tags,omitempty"`
	Category      string `json:"category,omitempty"`
	Status        string `json:"status"`
	Visibility    string `json:"visibility"`
	PublishedAt   string `json:"published_at,omitempty"`
	ScheduledAt   string `json:"scheduled_at,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	AuthorName    string `json:"author_name,omitempty"`
	ReadingTime   int    `json:"reading_time,omitempty"`
}

// toResponse converts a store.BlogPost to BlogPostResponse.
func blogPostToResponse(p *store.BlogPost, includeContent bool) BlogPostResponse {
	r := BlogPostResponse{
		ID:         p.ID,
		Slug:       p.Slug,
		Title:      p.Title,
		Status:     p.Status,
		Visibility: p.Visibility,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
	if p.Description.Valid {
		r.Description = p.Description.String
	}
	if includeContent {
		r.Content = p.Content
		// Render markdown to HTML
		renderer := markdown.NewRenderer()
		html, _ := renderer.Render(p.Content)
		r.ContentHTML = html
	}
	if p.FeaturedImage.Valid {
		r.FeaturedImage = p.FeaturedImage.String
	}
	if p.Tags != "" && p.Tags != "[]" {
		r.Tags = p.Tags
	}
	if p.Category.Valid {
		r.Category = p.Category.String
	}
	if p.PublishedAt.Valid {
		r.PublishedAt = p.PublishedAt.String
	}
	if p.ScheduledAt.Valid {
		r.ScheduledAt = p.ScheduledAt.String
	}
	if p.AuthorName.Valid {
		r.AuthorName = p.AuthorName.String
	}
	r.ReadingTime = calculateReadingTime(p.Content)
	return r
}

// RSS Feed structures
type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Category    string `xml:"category,omitempty"`
}

// BlogCommentRequest represents a comment submission.
// AuthorName and AuthorEmail are optional for authenticated users (auto-filled from profile).
type BlogCommentRequest struct {
	AuthorName  string `json:"author_name" form:"author_name"`
	AuthorEmail string `json:"author_email" form:"author_email"`
	Content     string `json:"content" form:"content" binding:"required,min=1,max=5000"`
	ParentID    string `json:"parent_id" form:"parent_id"`
}

// Public handlers

// listPublishedPosts returns published posts for the public API with pagination and filtering.
func (r *Router) listPublishedPosts(c *gin.Context) {
	siteID := c.Query("site_id")
	if siteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "site_id required"})
		return
	}

	// Pagination
	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
			if limit > 100 {
				limit = 100
			}
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	// Filtering
	category := c.Query("category")
	tag := c.Query("tag")
	search := c.Query("q")

	posts, total, err := r.db.ListPublishedBlogPostsFiltered(c.Request.Context(), siteID, category, tag, search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Convert to response type
	postsResp := make([]BlogPostResponse, 0, len(posts))
	for i := range posts {
		postsResp = append(postsResp, blogPostToResponse(&posts[i], false))
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":  postsResp,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// getPublishedPost returns a single published post by slug.
func (r *Router) getPublishedPost(c *gin.Context) {
	siteID := c.Query("site_id")
	if siteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "site_id required"})
		return
	}
	slug := c.Param("slug")

	post, err := r.db.GetBlogPostBySlug(c.Request.Context(), siteID, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if post == nil || post.Status != "published" {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	// Check visibility - public API only serves public posts
	if post.Visibility != "" && post.Visibility != "public" {
		c.JSON(http.StatusForbidden, gin.H{"error": "authentication required", "visibility": post.Visibility})
		return
	}

	resp := blogPostToResponse(post, true)
	c.JSON(http.StatusOK, gin.H{"post": resp})
}

// listApprovedCommentsPublic returns approved comments for a post.
func (r *Router) listApprovedCommentsPublic(c *gin.Context) {
	siteID := c.Query("site_id")
	if siteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "site_id required"})
		return
	}
	slug := c.Param("slug")

	post, err := r.db.GetBlogPostBySlug(c.Request.Context(), siteID, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if post == nil || post.Status != "published" {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	comments, err := r.db.ListApprovedComments(c.Request.Context(), post.ID, 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

// submitCommentPublic allows visitors to submit comments.
// Authenticated users auto-fill name/email from their profile.
func (r *Router) submitCommentPublic(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	if siteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "site_id required"})
		return
	}
	slug := c.Param("slug")

	post, err := r.db.GetBlogPostBySlug(c.Request.Context(), siteID, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if post == nil || post.Status != "published" {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	var req BlogCommentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var authorName, authorEmail string
	var isAuthenticated bool

	// Check if user is authenticated (set by OptionalAuthMiddleware)
	userID, exists := c.Get("user_id")
	if exists && userID != "" {
		// Fetch user from database to get name and email
		user, err := r.db.GetUserByID(c.Request.Context(), userID.(string))
		if err == nil && user != nil {
			isAuthenticated = true
			authorEmail = user.Email
			if user.Name.Valid && user.Name.String != "" {
				authorName = user.Name.String
			} else {
				// Fall back to email prefix if no name set
				authorName = strings.Split(user.Email, "@")[0]
			}
		}
	}

	// For anonymous users, require name and email from form
	if !isAuthenticated {
		if req.AuthorName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		if req.AuthorEmail == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
			return
		}
		authorName = req.AuthorName
		authorEmail = req.AuthorEmail
	}

	// Sanitize user input to prevent XSS
	authorName = sanitizeComment(authorName)
	content := sanitizeComment(req.Content)

	// Validate sanitized content isn't empty
	if authorName == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content after sanitization"})
		return
	}

	id := uuid.New().String()
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	_, err = r.db.CreateBlogComment(c.Request.Context(), id, siteID, post.ID, req.ParentID,
		authorName, authorEmail, content, ipAddress, userAgent)
	if err != nil {
		if c.GetHeader("HX-Request") == "true" {
			c.Header("Content-Type", "text/html")
			c.String(http.StatusInternalServerError, `<div class="blog-comment-error">Failed to submit comment. Please try again.</div>`)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit comment"})
		return
	}

	// Return HTML for HTMX requests
	if c.GetHeader("HX-Request") == "true" {
		// Re-render comments section with success message
		comments, _ := r.db.ListApprovedComments(c.Request.Context(), post.ID, 50, 0)
		html := fmt.Sprintf(`<section class="blog-comments" id="comments"><h3>Comments (%d)</h3>`, len(comments))
		html += `<div class="blog-comment-success">Thank you! Your comment has been submitted and is awaiting moderation.</div>`

		if len(comments) == 0 {
			html += `<p class="blog-comments-empty">No comments yet. Be the first to share your thoughts!</p>`
		} else {
			for _, comment := range comments {
				avatar := blogGetInitials(comment.AuthorName)
				html += fmt.Sprintf(`
					<div class="blog-comment">
						<div class="blog-comment-avatar">%s</div>
						<div class="blog-comment-body">
							<div class="blog-comment-header">
								<span class="blog-comment-author">%s</span>
								<span class="blog-comment-date">%s</span>
							</div>
							<div class="blog-comment-content">%s</div>
						</div>
					</div>`,
					avatar,
					escapeHTML(comment.AuthorName),
					blogFormatTimeAgo(comment.CreatedAt),
					escapeHTML(comment.Content),
				)
			}
		}

		// Comment form - show simplified form for authenticated users
		html += fmt.Sprintf(`
			<form class="blog-comment-form" hx-post="/api/blog/posts/%s/comments" hx-target="#comments" hx-swap="outerHTML" hx-on::after-request="this.reset()">
				<h4>Leave a Comment</h4>`, escapeHTML(slug))

		if isAuthenticated {
			html += fmt.Sprintf(`<p class="blog-comment-logged-in">Commenting as <strong>%s</strong></p>`, escapeHTML(authorName))
		} else {
			html += `
				<div class="blog-comment-form-row">
					<input type="text" name="author_name" placeholder="Your name" required>
					<input type="email" name="author_email" placeholder="Your email" required>
				</div>`
		}

		html += `
				<textarea name="content" placeholder="Write your comment..." required></textarea>
				<button type="submit" class="blog-btn blog-btn-primary">Post Comment</button>
				<p class="blog-comment-notice">Comments are moderated and may take a moment to appear.</p>
			</form>
		</section>`

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "pending_moderation"})
}

// Admin handlers

// listAllPosts returns all posts for admin view.
func (r *Router) listAllPosts(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	status := c.Query("status")

	posts, err := r.db.ListBlogPosts(c.Request.Context(), siteID, status, 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// createPost creates a new blog post.
func (r *Router) createPost(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	authorID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req BlogPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check slug uniqueness
	existing, err := r.db.GetBlogPostBySlug(c.Request.Context(), siteID, req.Slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "slug already exists"})
		return
	}

	id := uuid.New().String()
	tags := req.Tags
	if tags == "" {
		tags = "[]"
	}

	post, err := r.db.CreateBlogPost(c.Request.Context(), id, siteID, authorID,
		req.Slug, req.Title, req.Description, req.Content, req.FeaturedImage, tags, req.Category, req.Visibility)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	// Audit log: blog post created
	r.logAuditAction(c, "create", "blog_post", post.ID, req.Title, "")

	c.JSON(http.StatusCreated, gin.H{"post": post})
}

// getPost returns a single post by ID (scoped to authenticated user's site).
func (r *Router) getPost(c *gin.Context) {
	id := c.Param("id")
	siteID, err := getSiteID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	post, err := r.db.GetBlogPostByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if post == nil || post.SiteID != siteID {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	// Set post author for middleware checks
	c.Set("post_author_id", post.AuthorID.String)

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// updatePost updates an existing post.
func (r *Router) updatePost(c *gin.Context) {
	id := c.Param("id")
	siteID, err := getSiteID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	post, err := r.db.GetBlogPostByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	// Check authorization for authors
	role, err := getUserRole(c)
	if err == nil && role == "author" {
		userID, err := getUserID(c)
		if err != nil || post.AuthorID.String != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "can only edit own posts"})
			return
		}
	}

	var req BlogPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check slug uniqueness if changed
	if req.Slug != post.Slug {
		existing, err := r.db.GetBlogPostBySlug(c.Request.Context(), siteID, req.Slug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
		if existing != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "slug already exists"})
			return
		}
	}

	tags := req.Tags
	if tags == "" {
		tags = "[]"
	}

	err = r.db.UpdateBlogPost(c.Request.Context(), id, req.Slug, req.Title, req.Description,
		req.Content, req.FeaturedImage, tags, req.Category, req.Visibility)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update post"})
		return
	}

	// Audit log: blog post updated
	r.logAuditAction(c, "update", "blog_post", id, req.Title, "")

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// deletePost deletes a post (enforces author ownership and site scoping).
func (r *Router) deletePost(c *gin.Context) {
	id := c.Param("id")
	siteID, err := getSiteID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get post info for ownership check and audit log
	post, err := r.db.GetBlogPostByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if post == nil || post.SiteID != siteID {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	if !canAuthorEditPost(c, post) {
		denyAuthorEdit(c, "delete")
		return
	}

	postTitle := post.Title
	err = r.db.DeleteBlogPost(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post"})
		return
	}

	// Audit log: blog post deleted
	r.logAuditAction(c, "delete", "blog_post", id, postTitle, "")

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// publishPost publishes a draft post.
func (r *Router) publishPost(c *gin.Context) {
	id := c.Param("id")

	post, err := r.db.GetBlogPostByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	if !canAuthorEditPost(c, post) {
		denyAuthorEdit(c, "publish")
		return
	}

	err = r.db.PublishBlogPost(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish post"})
		return
	}

	// Audit log: blog post published
	r.logAuditAction(c, "publish", "blog_post", id, post.Title, "")

	c.JSON(http.StatusOK, gin.H{"status": "published"})
}

// unpublishPost unpublishes a post.
func (r *Router) unpublishPost(c *gin.Context) {
	id := c.Param("id")

	post, err := r.db.GetBlogPostByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	if !canAuthorEditPost(c, post) {
		denyAuthorEdit(c, "unpublish")
		return
	}

	err = r.db.UnpublishBlogPost(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unpublish post"})
		return
	}

	// Audit log: blog post unpublished
	r.logAuditAction(c, "unpublish", "blog_post", id, post.Title, "")

	c.JSON(http.StatusOK, gin.H{"status": "unpublished"})
}

// ScheduleRequest represents the request to schedule a post.
type ScheduleRequest struct {
	ScheduledAt string `json:"scheduled_at" binding:"required"` // ISO 8601 datetime
}

// schedulePost schedules a post for future publishing.
func (r *Router) schedulePost(c *gin.Context) {
	id := c.Param("id")

	post, err := r.db.GetBlogPostByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	if !canAuthorEditPost(c, post) {
		denyAuthorEdit(c, "schedule")
		return
	}

	var req ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate datetime
	scheduledTime, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid datetime format, use ISO 8601 (RFC3339)"})
		return
	}

	if scheduledTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scheduled time must be in the future"})
		return
	}

	err = r.db.ScheduleBlogPost(c.Request.Context(), id, req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to schedule post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "scheduled", "scheduled_at": req.ScheduledAt})
}

// unschedulePost removes the scheduled publish time.
func (r *Router) unschedulePost(c *gin.Context) {
	id := c.Param("id")

	post, err := r.db.GetBlogPostByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	if !canAuthorEditPost(c, post) {
		denyAuthorEdit(c, "unschedule")
		return
	}

	err = r.db.ClearSchedule(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unschedule post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "unscheduled"})
}

// Comment management handlers

// listCommentsAdmin returns all comments for admin view.
func (r *Router) listCommentsAdmin(c *gin.Context) {
	siteID, _ := getSiteID(c)
	status := c.Query("status")

	comments, err := r.db.ListBlogComments(c.Request.Context(), siteID, status, 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

// approveComment approves a pending comment.
func (r *Router) approveComment(c *gin.Context) {
	id := c.Param("id")
	moderatorID, _ := getUserID(c)

	err := r.db.ModerateComment(c.Request.Context(), id, "approved", moderatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve comment"})
		return
	}

	// Audit log: comment approved
	r.logAuditAction(c, "approve", "blog_comment", id, "", "")

	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}

// rejectComment rejects a comment.
func (r *Router) rejectComment(c *gin.Context) {
	id := c.Param("id")
	moderatorID, _ := getUserID(c)

	err := r.db.ModerateComment(c.Request.Context(), id, "rejected", moderatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reject comment"})
		return
	}

	// Audit log: comment rejected
	r.logAuditAction(c, "reject", "blog_comment", id, "", "")

	c.JSON(http.StatusOK, gin.H{"status": "rejected"})
}

// markSpam marks a comment as spam.
func (r *Router) markSpam(c *gin.Context) {
	id := c.Param("id")
	moderatorID, _ := getUserID(c)

	err := r.db.ModerateComment(c.Request.Context(), id, "spam", moderatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark as spam"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "spam"})
}

// deleteComment deletes a comment (admin only).
func (r *Router) deleteComment(c *gin.Context) {
	id := c.Param("id")

	err := r.db.DeleteBlogComment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete comment"})
		return
	}

	// Audit log: comment deleted
	r.logAuditAction(c, "delete", "blog_comment", id, "", "")

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// Admin UI handlers

func (r *Router) adminBlog(c *gin.Context) {
	claims, _ := getUserClaims(c)
	c.HTML(http.StatusOK, "blog.html", gin.H{
		"Title":       "Blog",
		"CurrentPage": "blog",
		"User":        claims,
	})
}

func (r *Router) adminBlogEditor(c *gin.Context) {
	claims, _ := getUserClaims(c)
	id := c.Param("id")

	var post *store.BlogPost
	title := "New Post"
	if id != "" {
		p, err := r.db.GetBlogPostByID(c.Request.Context(), id)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Title": "Error",
				"Error": "Failed to load post",
			})
			return
		}
		if p == nil {
			c.HTML(http.StatusNotFound, "error.html", gin.H{
				"Title": "Not Found",
				"Error": "Post not found",
			})
			return
		}

		// Check authorization for authors - they can only edit their own posts
		role, roleErr := getUserRole(c)
		userID, userErr := getUserID(c)
		if roleErr == nil && userErr == nil && role == "author" {
			if p.AuthorID.String != userID {
				c.HTML(http.StatusForbidden, "error.html", gin.H{
					"Title": "Forbidden",
					"Error": "You can only edit your own posts",
				})
				return
			}
		}
		post = p
		title = "Edit: " + p.Title
	}

	c.HTML(http.StatusOK, "blog-editor.html", gin.H{
		"Title":       title,
		"CurrentPage": "blog",
		"User":        claims,
		"Post":        post,
	})
}

func (r *Router) adminComments(c *gin.Context) {
	claims, _ := getUserClaims(c)
	c.HTML(http.StatusOK, "comments.html", gin.H{
		"Title":       "Comments",
		"CurrentPage": "comments",
		"User":        claims,
	})
}

// RSS Feed handler

// getRSSFeed returns an RSS 2.0 feed of published posts.
func (r *Router) getRSSFeed(c *gin.Context) {
	siteID := c.Query("site_id")
	if siteID == "" {
		c.String(http.StatusBadRequest, "site_id required")
		return
	}

	site, err := r.db.GetSiteByID(c.Request.Context(), siteID)
	if err != nil || site == nil {
		c.String(http.StatusNotFound, "site not found")
		return
	}

	posts, _, err := r.db.ListPublishedBlogPostsFiltered(c.Request.Context(), siteID, "", "", "", 20, 0)
	if err != nil {
		c.String(http.StatusInternalServerError, "database error")
		return
	}

	baseURL := c.Query("base_url")
	if baseURL == "" {
		baseURL = "https://" + site.Domain.String
	}

	items := make([]RSSItem, 0, len(posts))
	for _, post := range posts {
		pubDate := post.CreatedAt
		if post.PublishedAt.Valid {
			pubDate = post.PublishedAt.String
		}
		// Parse and format as RFC1123
		if t, err := time.Parse(time.RFC3339, pubDate); err == nil {
			pubDate = t.Format(time.RFC1123Z)
		}

		desc := post.Description.String
		if desc == "" {
			desc = generateExcerpt(post.Content, 300)
		}

		category := ""
		if post.Category.Valid {
			category = post.Category.String
		}

		items = append(items, RSSItem{
			Title:       post.Title,
			Link:        fmt.Sprintf("%s/blog/%s", baseURL, post.Slug),
			Description: desc,
			PubDate:     pubDate,
			GUID:        post.ID,
			Category:    category,
		})
	}

	feed := RSSFeed{
		Version: "2.0",
		Channel: RSSChannel{
			Title:         site.Name + " Blog",
			Link:          baseURL + "/blog",
			Description:   site.Name + " blog feed",
			Language:      "en",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         items,
		},
	}

	c.Header("Content-Type", "application/rss+xml; charset=utf-8")
	c.XML(http.StatusOK, feed)
}

// getPostMeta returns social sharing metadata for a post.
func (r *Router) getPostMeta(c *gin.Context) {
	siteID := c.Query("site_id")
	if siteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "site_id required"})
		return
	}
	slug := c.Param("slug")

	post, err := r.db.GetBlogPostBySlug(c.Request.Context(), siteID, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if post == nil || post.Status != "published" {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	baseURL := c.Query("base_url")
	if baseURL == "" {
		baseURL = c.Request.Host
	}

	description := post.Description.String
	if description == "" {
		description = generateExcerpt(post.Content, 160)
	}

	image := ""
	if post.FeaturedImage.Valid {
		image = post.FeaturedImage.String
	}

	authorName := ""
	if post.AuthorName.Valid {
		authorName = post.AuthorName.String
	}

	c.JSON(http.StatusOK, gin.H{
		"title":        post.Title,
		"description":  description,
		"url":          fmt.Sprintf("%s/blog/%s", baseURL, post.Slug),
		"image":        image,
		"author":       authorName,
		"published_at": post.PublishedAt.String,
		"reading_time": calculateReadingTime(post.Content),
		"og": gin.H{
			"title":       post.Title,
			"description": description,
			"type":        "article",
			"url":         fmt.Sprintf("%s/blog/%s", baseURL, post.Slug),
			"image":       image,
		},
		"twitter": gin.H{
			"card":        "summary_large_image",
			"title":       post.Title,
			"description": description,
			"image":       image,
		},
	})
}

// getRelatedPosts returns posts related to the given post.
func (r *Router) getRelatedPosts(c *gin.Context) {
	siteID := c.Query("site_id")
	if siteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "site_id required"})
		return
	}
	slug := c.Param("slug")

	post, err := r.db.GetBlogPostBySlug(c.Request.Context(), siteID, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	// Get related posts by category (excluding current post)
	related, err := r.db.GetRelatedPosts(c.Request.Context(), siteID, post.ID, post.Category.String, post.Tags, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"related": related})
}

// Helper functions

// calculateReadingTime estimates reading time in minutes.
// Assumes average reading speed of 200 words per minute.
func calculateReadingTime(content string) int {
	// Strip markdown syntax for more accurate word count
	stripped := stripMarkdown(content)
	words := len(strings.Fields(stripped))
	minutes := words / 200
	if minutes < 1 {
		minutes = 1
	}
	return minutes
}

// generateExcerpt creates a plain text excerpt from markdown content.
func generateExcerpt(content string, maxLen int) string {
	// Strip markdown
	text := stripMarkdown(content)

	// Trim to maxLen
	if len(text) > maxLen {
		// Find last space before maxLen
		text = text[:maxLen]
		if lastSpace := strings.LastIndex(text, " "); lastSpace > 0 {
			text = text[:lastSpace]
		}
		text += "..."
	}

	return strings.TrimSpace(text)
}

// stripMarkdown removes markdown formatting from text.
func stripMarkdown(content string) string {
	// Remove code blocks
	re := regexp.MustCompile("```[\\s\\S]*?```")
	content = re.ReplaceAllString(content, "")

	// Remove inline code
	re = regexp.MustCompile("`[^`]+`")
	content = re.ReplaceAllString(content, "")

	// Remove images
	re = regexp.MustCompile(`!\[[^\]]*\]\([^)]+\)`)
	content = re.ReplaceAllString(content, "")

	// Remove links but keep text
	re = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	content = re.ReplaceAllString(content, "$1")

	// Remove headers
	re = regexp.MustCompile(`(?m)^#+\s*`)
	content = re.ReplaceAllString(content, "")

	// Remove bold/italic
	re = regexp.MustCompile(`[*_]{1,3}([^*_]+)[*_]{1,3}`)
	content = re.ReplaceAllString(content, "$1")

	// Remove blockquotes
	re = regexp.MustCompile(`(?m)^>\s*`)
	content = re.ReplaceAllString(content, "")

	// Remove horizontal rules
	re = regexp.MustCompile(`(?m)^[-*_]{3,}\s*$`)
	content = re.ReplaceAllString(content, "")

	// Collapse whitespace
	re = regexp.MustCompile(`\s+`)
	content = re.ReplaceAllString(content, " ")

	return strings.TrimSpace(content)
}

// getSitemap returns a dynamic sitemap.xml including blog posts.
func (r *Router) getSitemap(c *gin.Context) {
	siteID := c.Query("site_id")
	baseURL := c.Query("base_url")
	if siteID == "" || baseURL == "" {
		c.String(http.StatusBadRequest, "site_id and base_url required")
		return
	}

	baseURL = strings.TrimSuffix(baseURL, "/")

	type URL struct {
		Loc        string `xml:"loc"`
		LastMod    string `xml:"lastmod,omitempty"`
		ChangeFreq string `xml:"changefreq,omitempty"`
	}

	type URLSet struct {
		XMLName xml.Name `xml:"urlset"`
		XMLNS   string   `xml:"xmlns,attr"`
		URLs    []URL    `xml:"url"`
	}

	urlset := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  []URL{},
	}

	// Add blog list page
	urlset.URLs = append(urlset.URLs, URL{
		Loc:        baseURL + "/blog/",
		ChangeFreq: "daily",
	})

	// Get published posts
	posts, _, err := r.db.ListPublishedBlogPostsFiltered(c.Request.Context(), siteID, "", "", "", 1000, 0)
	if err == nil {
		for _, post := range posts {
			lastMod := post.UpdatedAt
			if post.PublishedAt.Valid {
				lastMod = post.PublishedAt.String
			}
			// Parse and format the date
			if t, err := time.Parse("2006-01-02 15:04:05", lastMod); err == nil {
				lastMod = t.Format("2006-01-02")
			} else if t, err := time.Parse(time.RFC3339, lastMod); err == nil {
				lastMod = t.Format("2006-01-02")
			}

			urlset.URLs = append(urlset.URLs, URL{
				Loc:        baseURL + "/blog/?article=" + post.Slug,
				LastMod:    lastMod,
				ChangeFreq: "monthly",
			})
		}
	}

	c.Header("Content-Type", "application/xml")
	c.XML(http.StatusOK, urlset)
}
