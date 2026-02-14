package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/markdown"
)

// Blog stats fragment
func (r *Router) fragmentBlogStats(c *gin.Context) {
	siteID, _ := getSiteID(c)

	total, published, draft, _ := r.db.GetBlogPostStats(c.Request.Context(), siteID)
	totalComments, pendingComments, _, _ := r.db.GetCommentStats(c.Request.Context(), siteID)

	html := buildStatCards([]StatCard{
		{Value: total, Label: "Total Posts"},
		{Value: published, Label: "Published"},
		{Value: draft, Label: "Drafts"},
		{Value: pendingComments, Label: "Pending Comments"},
	})

	// Add warning if pending comments
	if pendingComments > 0 {
		html += fmt.Sprintf(`
			<div class="alert alert-warning">
				%d comment(s) awaiting moderation
			</div>
		`, pendingComments)
	}

	// Add total comments
	html += buildStatCards([]StatCard{
		{Value: totalComments, Label: "Total Comments", Secondary: true},
	})

	respondHTML(c, html)
}

// Blog post list fragment
func (r *Router) fragmentBlogPostList(c *gin.Context) {
	siteID, _ := getSiteID(c)
	status := c.Query("status")

	posts, err := r.db.ListBlogPosts(c.Request.Context(), siteID, status, 50, 0)
	if err != nil {
		respondHTMLError(c, "Failed to load posts")
		return
	}

	if len(posts) == 0 {
		respondHTML(c, `
			<div class="empty-state">
				<p>No posts found</p>
				<a href="/admin/blog/new" class="btn btn-primary">Create your first post</a>
			</div>
		`)
		return
	}

	html := `<table class="data-table">
		<thead>
			<tr>
				<th>Title</th>
				<th>Status</th>
				<th>Author</th>
				<th>Created</th>
				<th>Actions</th>
			</tr>
		</thead>
		<tbody>`

	for _, post := range posts {
		statusClass := "status-draft"
		if post.Status == "published" {
			statusClass = "status-published"
		} else if post.Status == "archived" {
			statusClass = "status-archived"
		}

		authorName := "Unknown"
		if post.AuthorName.Valid {
			authorName = post.AuthorName.String
		}

		html += fmt.Sprintf(`
			<tr>
				<td>
					<a href="/admin/blog/edit/%s" class="post-title">%s</a>
					<span class="post-slug">/%s</span>
				</td>
				<td><span class="%s">%s</span></td>
				<td>%s</td>
				<td>%s</td>
				<td class="actions">
					<a href="/admin/blog/edit/%s" class="btn btn-sm">Edit</a>
				</td>
			</tr>`,
			post.ID,
			escapeHTML(post.Title),
			escapeHTML(post.Slug),
			statusClass,
			post.Status,
			escapeHTML(authorName),
			post.CreatedAt,
			post.ID,
		)
	}

	html += `</tbody></table>`

	respondHTML(c, html)
}

// Blog post editor fragment (for new/edit form)
func (r *Router) fragmentBlogPostEditor(c *gin.Context) {
	id := c.Param("id")

	var html string
	if id == "" {
		// New post form
		html = `
			<form id="post-form" hx-post="/api/blog/posts" hx-headers='{"Content-Type": "application/json"}' hx-ext="json-enc">
				<div class="form-group">
					<label for="title">Title</label>
					<input type="text" id="title" name="title" required placeholder="Post title">
				</div>
				<div class="form-group">
					<label for="slug">Slug</label>
					<input type="text" id="slug" name="slug" required placeholder="url-friendly-slug">
				</div>
				<div class="form-group">
					<label for="description">Description</label>
					<textarea id="description" name="description" rows="2" placeholder="Brief description for SEO"></textarea>
				</div>
				<div class="form-group">
					<label for="category">Category</label>
					<input type="text" id="category" name="category" placeholder="Category">
				</div>
				<div class="form-group">
					<label for="tags">Tags (JSON array)</label>
					<input type="text" id="tags" name="tags" placeholder='["tag1", "tag2"]' value="[]">
				</div>
				<div class="form-group">
					<label for="featured_image">Featured Image URL</label>
					<input type="text" id="featured_image" name="featured_image" placeholder="https://...">
				</div>
				<div class="form-group">
					<label for="content">Content (Markdown)</label>
					<div class="editor-toolbar">
						<button type="button" onclick="insertMarkdown('**', '**')" title="Bold">B</button>
						<button type="button" onclick="insertMarkdown('*', '*')" title="Italic">I</button>
						<button type="button" onclick="insertMarkdown('## ', '')" title="Heading">H</button>
						<button type="button" onclick="insertMarkdown('[', '](url)')" title="Link">Link</button>
						<button type="button" onclick="insertMarkdown('\x60', '\x60')" title="Code">Code</button>
						<button type="button" onclick="insertMarkdown('\x60\x60\x60\n', '\n\x60\x60\x60')" title="Code Block">Block</button>
					</div>
					<textarea id="content" name="content" required rows="20" placeholder="Write your post in Markdown..."></textarea>
				</div>
				<div class="form-actions">
					<button type="submit" class="btn btn-primary">Save Draft</button>
					<button type="button" class="btn btn-secondary" hx-post="/admin/fragments/blog-post-preview" hx-include="#content" hx-target="#preview-pane">Preview</button>
				</div>
			</form>
			<div id="preview-pane" class="preview-pane"></div>
		`
	} else {
		// Edit existing post
		post, err := r.db.GetBlogPostByID(c.Request.Context(), id)
		if err != nil || post == nil {
			respondHTMLError(c, "Post not found")
			return
		}

		description := ""
		if post.Description.Valid {
			description = post.Description.String
		}
		category := ""
		if post.Category.Valid {
			category = post.Category.String
		}
		featuredImage := ""
		if post.FeaturedImage.Valid {
			featuredImage = post.FeaturedImage.String
		}

		html = fmt.Sprintf(`
			<form id="post-form" hx-put="/api/blog/posts/%s" hx-headers='{"Content-Type": "application/json"}' hx-ext="json-enc">
				<div class="form-group">
					<label for="title">Title</label>
					<input type="text" id="title" name="title" required value="%s">
				</div>
				<div class="form-group">
					<label for="slug">Slug</label>
					<input type="text" id="slug" name="slug" required value="%s">
				</div>
				<div class="form-group">
					<label for="description">Description</label>
					<textarea id="description" name="description" rows="2">%s</textarea>
				</div>
				<div class="form-group">
					<label for="category">Category</label>
					<input type="text" id="category" name="category" value="%s">
				</div>
				<div class="form-group">
					<label for="tags">Tags (JSON array)</label>
					<input type="text" id="tags" name="tags" value="%s">
				</div>
				<div class="form-group">
					<label for="featured_image">Featured Image URL</label>
					<input type="text" id="featured_image" name="featured_image" value="%s">
				</div>
				<div class="form-group">
					<label for="content">Content (Markdown)</label>
					<div class="editor-toolbar">
						<button type="button" onclick="insertMarkdown('**', '**')" title="Bold">B</button>
						<button type="button" onclick="insertMarkdown('*', '*')" title="Italic">I</button>
						<button type="button" onclick="insertMarkdown('## ', '')" title="Heading">H</button>
						<button type="button" onclick="insertMarkdown('[', '](url)')" title="Link">Link</button>
						<button type="button" onclick="insertMarkdown('\x60', '\x60')" title="Code">Code</button>
						<button type="button" onclick="insertMarkdown('\x60\x60\x60\n', '\n\x60\x60\x60')" title="Code Block">Block</button>
					</div>
					<textarea id="content" name="content" required rows="20">%s</textarea>
				</div>
				<div class="form-actions">
					<button type="submit" class="btn btn-primary">Save</button>
					<button type="button" class="btn btn-secondary" hx-post="/admin/fragments/blog-post-preview" hx-include="#content" hx-target="#preview-pane">Preview</button>
		`,
			post.ID,
			escapeHTML(post.Title),
			escapeHTML(post.Slug),
			escapeHTML(description),
			escapeHTML(category),
			escapeHTML(post.Tags),
			escapeHTML(featuredImage),
			escapeHTML(post.Content),
		)

		// Add publish/unpublish buttons
		if post.Status == "draft" {
			html += fmt.Sprintf(`
					<button type="button" class="btn btn-success" hx-post="/api/blog/posts/%s/publish" hx-confirm="Publish this post?">Publish</button>
			`, post.ID)
		} else if post.Status == "published" {
			html += fmt.Sprintf(`
					<button type="button" class="btn btn-warning" hx-post="/api/blog/posts/%s/unpublish" hx-confirm="Unpublish this post?">Unpublish</button>
			`, post.ID)
		}

		html += `
				</div>
			</form>
			<div id="preview-pane" class="preview-pane"></div>
		`
	}

	respondHTML(c, html)
}

// Blog post preview fragment
func (r *Router) fragmentBlogPostPreview(c *gin.Context) {
	content := c.PostForm("content")
	if content == "" {
		respondHTMLEmpty(c, "No content to preview")
		return
	}

	renderer := markdown.NewRenderer()
	html, err := renderer.Render(content)
	if err != nil {
		respondHTMLError(c, fmt.Sprintf("Render error: %s", err.Error()))
		return
	}

	respondHTML(c, fmt.Sprintf(`<div class="markdown-preview">%s</div>`, html))
}

// Comment list fragment
func (r *Router) fragmentCommentList(c *gin.Context) {
	siteID, _ := getSiteID(c)
	status := c.Query("status")

	comments, err := r.db.ListBlogComments(c.Request.Context(), siteID, status, 50, 0)
	if err != nil {
		respondHTMLError(c, "Failed to load comments")
		return
	}

	if len(comments) == 0 {
		respondHTMLEmpty(c, "No comments found")
		return
	}

	html := `<table class="data-table">
		<thead>
			<tr>
				<th>Author</th>
				<th>Content</th>
				<th>Status</th>
				<th>Created</th>
				<th>Actions</th>
			</tr>
		</thead>
		<tbody>`

	for _, comment := range comments {
		statusClass := "status-pending"
		switch comment.Status {
		case "approved":
			statusClass = "status-approved"
		case "rejected":
			statusClass = "status-rejected"
		case "spam":
			statusClass = "status-spam"
		}

		// Truncate content for display
		contentPreview := comment.Content
		if len(contentPreview) > 100 {
			contentPreview = contentPreview[:100] + "..."
		}

		html += fmt.Sprintf(`
			<tr>
				<td>
					<span class="comment-author">%s</span>
					<span class="comment-email">%s</span>
				</td>
				<td class="comment-content">%s</td>
				<td><span class="%s">%s</span></td>
				<td>%s</td>
				<td class="actions">`,
			escapeHTML(comment.AuthorName),
			escapeHTML(comment.AuthorEmail),
			escapeHTML(contentPreview),
			statusClass,
			comment.Status,
			comment.CreatedAt,
		)

		// Add action buttons based on status
		if comment.Status == "pending" {
			html += fmt.Sprintf(`
					<button class="btn btn-sm btn-success" hx-put="/api/blog/comments/%s/approve" hx-swap="outerHTML" hx-target="closest tr">Approve</button>
					<button class="btn btn-sm btn-warning" hx-put="/api/blog/comments/%s/reject" hx-swap="outerHTML" hx-target="closest tr">Reject</button>
					<button class="btn btn-sm btn-danger" hx-put="/api/blog/comments/%s/spam" hx-swap="outerHTML" hx-target="closest tr">Spam</button>
			`, comment.ID, comment.ID, comment.ID)
		}
		html += fmt.Sprintf(`
					<button class="btn btn-sm btn-danger" hx-delete="/api/blog/comments/%s" hx-confirm="Delete this comment?" hx-swap="outerHTML" hx-target="closest tr">Delete</button>
				</td>
			</tr>`, comment.ID)
	}

	html += `</tbody></table>`

	respondHTML(c, html)
}

// Comment moderation queue fragment
func (r *Router) fragmentCommentModerationQueue(c *gin.Context) {
	siteID, _ := getSiteID(c)

	comments, err := r.db.ListPendingComments(c.Request.Context(), siteID, 20, 0)
	if err != nil {
		respondHTMLError(c, "Failed to load comments")
		return
	}

	if len(comments) == 0 {
		respondHTMLEmpty(c, "No comments pending moderation")
		return
	}

	html := `<div class="moderation-queue">`

	for _, comment := range comments {
		html += fmt.Sprintf(`
			<div class="moderation-item" id="comment-%s">
				<div class="comment-header">
					<span class="comment-author">%s</span>
					<span class="comment-email">%s</span>
					<span class="comment-date">%s</span>
				</div>
				<div class="comment-body">%s</div>
				<div class="comment-actions">
					<button class="btn btn-sm btn-success" hx-put="/api/blog/comments/%s/approve" hx-swap="outerHTML" hx-target="#comment-%s">Approve</button>
					<button class="btn btn-sm btn-warning" hx-put="/api/blog/comments/%s/reject" hx-swap="outerHTML" hx-target="#comment-%s">Reject</button>
					<button class="btn btn-sm btn-danger" hx-put="/api/blog/comments/%s/spam" hx-swap="outerHTML" hx-target="#comment-%s">Spam</button>
				</div>
			</div>
		`,
			comment.ID,
			escapeHTML(comment.AuthorName),
			escapeHTML(comment.AuthorEmail),
			comment.CreatedAt,
			escapeHTML(comment.Content),
			comment.ID, comment.ID,
			comment.ID, comment.ID,
			comment.ID, comment.ID,
		)
	}

	html += `</div>`

	respondHTML(c, html)
}
