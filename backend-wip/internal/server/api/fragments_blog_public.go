package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/markdown"
)

// Blog post cards fragment (for blog.html list view)
func (r *Router) fragmentBlogPostCards(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	search := c.Query("q")
	category := c.Query("category")
	tag := c.Query("tag")
	limit := 12

	posts, _, err := r.db.ListPublishedBlogPostsFiltered(c.Request.Context(), siteID, category, tag, search, limit, 0)
	if err != nil {
		respondHTMLError(c, "Failed to load posts")
		return
	}

	if len(posts) == 0 {
		msg := "No posts found"
		if search != "" {
			msg = fmt.Sprintf("No posts found for \"%s\"", escapeHTML(search))
		}
		respondHTML(c, fmt.Sprintf(`<div class="blog-empty"><div class="blog-empty-icon">:(</div>%s</div>`, msg))
		return
	}

	html := ""
	for _, post := range posts {
		html += `<article class="blog-card">`

		// Featured image
		if post.FeaturedImage.Valid && post.FeaturedImage.String != "" {
			html += fmt.Sprintf(`<img src="%s" alt="" class="blog-card-image" loading="lazy">`,
				escapeHTML(post.FeaturedImage.String))
		}

		html += `<div class="blog-card-body">`

		// Category
		if post.Category.Valid && post.Category.String != "" {
			html += fmt.Sprintf(`<span class="blog-card-category">%s</span>`,
				escapeHTML(post.Category.String))
		}

		// Title
		html += fmt.Sprintf(`<h2 class="blog-card-title"><a href="/blog/%s/">%s</a></h2>`,
			escapeHTML(post.Slug),
			escapeHTML(post.Title))

		// Description/excerpt
		if post.Description.Valid && post.Description.String != "" {
			html += fmt.Sprintf(`<p class="blog-card-excerpt">%s</p>`,
				escapeHTML(post.Description.String))
		}

		// Meta (date, reading time)
		html += `<div class="blog-card-meta">`
		if post.PublishedAt.Valid {
			html += fmt.Sprintf(`<span class="blog-card-date">%s</span>`,
				blogFormatDate(post.PublishedAt.String))
		}
		readingTime := blogEstimateReadingTime(post.Content)
		if readingTime > 0 {
			html += fmt.Sprintf(`<span class="blog-card-reading">%d min read</span>`, readingTime)
		}
		html += `</div>`

		html += `</div></article>`
	}

	respondHTML(c, html)
}

// Blog search dropdown fragment (instant search results)
func (r *Router) fragmentBlogSearchDropdown(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	query := c.Query("q")

	if len(query) < 2 {
		respondHTML(c, "")
		return
	}

	posts, _, err := r.db.ListPublishedBlogPostsFiltered(c.Request.Context(), siteID, "", "", query, 5, 0)
	if err != nil {
		respondHTMLError(c, "Search failed")
		return
	}

	if len(posts) == 0 {
		respondHTML(c, `<div class="blog-search-empty">No results found</div>`)
		return
	}

	html := ""
	for _, post := range posts {
		category := ""
		if post.Category.Valid && post.Category.String != "" {
			category = post.Category.String
		}

		html += fmt.Sprintf(`
			<a href="/blog/%s/" class="blog-search-result">
				<span class="blog-search-result-title">%s</span>
				<span class="blog-search-result-meta">%s</span>
			</a>`,
			escapeHTML(post.Slug),
			escapeHTML(post.Title),
			escapeHTML(category),
		)
	}

	respondHTML(c, html)
}

// Blog article fragment (full article content for blog-article.html)
func (r *Router) fragmentBlogArticle(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	slug := c.Param("slug")

	if slug == "" {
		respondHTMLError(c, "Article not found")
		return
	}

	post, err := r.db.GetBlogPostBySlug(c.Request.Context(), siteID, slug)
	if err != nil || post == nil {
		respondHTML(c, `<div class="blog-article-error">Article not found</div>`)
		return
	}

	// Only show published public posts
	if post.Status != "published" || post.Visibility != "public" {
		respondHTML(c, `<div class="blog-article-error">Article not found</div>`)
		return
	}

	html := `<header class="blog-article-header">`
	html += `<a href="/blog/" class="blog-article-back">&larr; Back to Blog</a>`

	// Category
	if post.Category.Valid && post.Category.String != "" {
		html += fmt.Sprintf(`<span class="blog-article-category">%s</span>`,
			escapeHTML(post.Category.String))
	}

	// Title
	html += fmt.Sprintf(`<h1 class="blog-article-title">%s</h1>`, escapeHTML(post.Title))

	// Meta
	html += `<div class="blog-article-meta">`
	if post.PublishedAt.Valid {
		html += fmt.Sprintf(`<time datetime="%s">%s</time>`,
			escapeHTML(post.PublishedAt.String),
			blogFormatDateLong(post.PublishedAt.String))
	}
	readingTime := blogEstimateReadingTime(post.Content)
	if readingTime > 0 {
		html += fmt.Sprintf(`<span>%d min read</span>`, readingTime)
	}
	if post.AuthorName.Valid && post.AuthorName.String != "" {
		html += fmt.Sprintf(`<span>by %s</span>`, escapeHTML(post.AuthorName.String))
	}
	html += `</div></header>`

	// Featured image
	if post.FeaturedImage.Valid && post.FeaturedImage.String != "" {
		html += fmt.Sprintf(`<img src="%s" alt="" class="blog-article-featured">`,
			escapeHTML(post.FeaturedImage.String))
	}

	// Content (render markdown)
	contentHTML := post.Content
	renderer := markdown.NewRenderer()
	if rendered, err := renderer.Render(post.Content); err == nil {
		contentHTML = rendered
	}
	html += fmt.Sprintf(`<div class="blog-article-content">%s</div>`, contentHTML)

	// Inject title update script
	html += fmt.Sprintf(`<script>document.title = %q;</script>`,
		escapeHTML(post.Title)+" | Blog")

	respondHTML(c, html)
}

// Blog related posts fragment
func (r *Router) fragmentBlogRelated(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	slug := c.Param("slug")

	if slug == "" {
		respondHTML(c, "")
		return
	}

	post, err := r.db.GetBlogPostBySlug(c.Request.Context(), siteID, slug)
	if err != nil || post == nil {
		respondHTML(c, "")
		return
	}

	category := ""
	if post.Category.Valid {
		category = post.Category.String
	}

	related, err := r.db.GetRelatedPosts(c.Request.Context(), siteID, post.ID, category, post.Tags, 3)
	if err != nil || len(related) == 0 {
		respondHTML(c, "")
		return
	}

	html := `<section class="blog-related"><h3>Related Posts</h3><div class="blog-related-grid">`

	for _, p := range related {
		html += `<article class="blog-related-card">`
		if p.FeaturedImage.Valid && p.FeaturedImage.String != "" {
			html += fmt.Sprintf(`<img src="%s" alt="" class="blog-related-image" loading="lazy">`,
				escapeHTML(p.FeaturedImage.String))
		}
		html += fmt.Sprintf(`<h4><a href="/blog/%s/">%s</a></h4>`,
			escapeHTML(p.Slug),
			escapeHTML(p.Title))
		html += `</article>`
	}

	html += `</div></section>`

	respondHTML(c, html)
}

// Blog comments fragment (public view)
// Shows simplified form for authenticated users.
func (r *Router) fragmentBlogComments(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	slug := c.Param("slug")

	if slug == "" {
		respondHTML(c, "")
		return
	}

	post, err := r.db.GetBlogPostBySlug(c.Request.Context(), siteID, slug)
	if err != nil || post == nil {
		respondHTML(c, "")
		return
	}

	// Check if user is authenticated (set by OptionalAuthMiddleware)
	var userName string
	var isAuthenticated bool
	userID, exists := c.Get("user_id")
	if exists && userID != "" {
		user, err := r.db.GetUserByID(c.Request.Context(), userID.(string))
		if err == nil && user != nil {
			isAuthenticated = true
			if user.Name.Valid && user.Name.String != "" {
				userName = user.Name.String
			} else {
				userName = strings.Split(user.Email, "@")[0]
			}
		}
	}

	comments, err := r.db.ListApprovedComments(c.Request.Context(), post.ID, 50, 0)
	if err != nil {
		respondHTMLError(c, "Failed to load comments")
		return
	}

	html := fmt.Sprintf(`<section class="blog-comments" id="comments"><h3>Comments (%d)</h3>`, len(comments))

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
		html += fmt.Sprintf(`<p class="blog-comment-logged-in">Commenting as <strong>%s</strong></p>`, escapeHTML(userName))
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

	respondHTML(c, html)
}

// Helper: estimate reading time based on word count
func blogEstimateReadingTime(content string) int {
	words := len(strings.Fields(content))
	// Average reading speed: 200 words per minute
	minutes := words / 200
	if minutes < 1 && words > 0 {
		return 1
	}
	return minutes
}

// Helper: format date as "Jan 2, 2006"
func blogFormatDate(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			return dateStr
		}
	}
	return t.Format("Jan 2, 2006")
}

// Helper: format date as "January 2, 2006"
func blogFormatDateLong(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			return dateStr
		}
	}
	return t.Format("January 2, 2006")
}

// Helper: format time as relative string
func blogFormatTimeAgo(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			return dateStr
		}
	}

	diff := time.Since(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d mins ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("Jan 2, 2006")
	}
}

// Helper: get initials from name
func blogGetInitials(name string) string {
	if name == "" {
		return "?"
	}
	words := strings.Fields(name)
	initials := ""
	for i, word := range words {
		if i >= 2 {
			break
		}
		if len(word) > 0 {
			initials += string([]rune(word)[0])
		}
	}
	if initials == "" {
		return "?"
	}
	return strings.ToUpper(initials)
}
