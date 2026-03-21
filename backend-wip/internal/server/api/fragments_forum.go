package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/markdown"
	"github.com/studiowebux/minimaldoc/internal/server/store"
)

// Forum categories list fragment
func (r *Router) fragmentForumCategories(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)

	categories, err := r.db.ListForumCategories(c.Request.Context(), siteID)
	if err != nil {
		respondHTMLError(c, "Failed to load categories")
		return
	}

	if len(categories) == 0 {
		respondHTML(c, `<div class="forum-empty">No categories yet</div>`)
		return
	}

	html := ""
	for _, cat := range categories {
		color := "#3b82f6"
		if cat.Color.Valid && cat.Color.String != "" {
			color = cat.Color.String
		}
		description := ""
		if cat.Description.Valid {
			description = cat.Description.String
		}
		html += fmt.Sprintf(`
			<div class="forum-category">
				<div class="forum-category-color" style="background:%s"></div>
				<div class="forum-category-info">
					<h3 class="forum-category-name"><a href="/forum/category/%s/">%s</a></h3>
					<p class="forum-category-description">%s</p>
				</div>
				<div class="forum-category-stats">
					<div class="forum-category-stat">
						<span class="forum-category-stat-value">%d</span>
						<span class="forum-category-stat-label">Topics</span>
					</div>
					<div class="forum-category-stat">
						<span class="forum-category-stat-value">%d</span>
						<span class="forum-category-stat-label">Posts</span>
					</div>
				</div>
			</div>`,
			color,
			escapeHTML(cat.Slug),
			escapeHTML(cat.Name),
			escapeHTML(description),
			cat.TopicCount,
			cat.PostCount,
		)
	}

	respondHTML(c, html)
}

// Forum topics list fragment
func (r *Router) fragmentForumTopics(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	categoryID := c.Query("category_id")
	search := c.Query("q")
	limit := 10

	var topics []store.ForumTopic
	var err error

	if search != "" {
		topics, _, err = r.db.SearchForum(c.Request.Context(), siteID, search, limit, 0)
	} else {
		topics, _, err = r.db.ListForumTopics(c.Request.Context(), siteID, categoryID, "", "", limit, 0)
	}

	if err != nil {
		respondHTMLError(c, "Failed to load topics")
		return
	}

	if len(topics) == 0 {
		respondHTML(c, `<div class="forum-empty">No discussions yet. Be the first to start one!</div>`)
		return
	}

	html := ""
	for _, topic := range topics {
		badges := ""
		pinnedClass := ""
		if topic.IsPinned {
			badges += `<span class="forum-topic-badge forum-topic-badge--pinned">Pinned</span>`
			pinnedClass = " forum-topic--pinned"
		}
		if topic.IsSolved {
			badges += `<span class="forum-topic-badge forum-topic-badge--solved">Solved</span>`
		}
		if topic.Status == "locked" {
			badges += `<span class="forum-topic-badge forum-topic-badge--locked">Locked</span>`
		}

		authorName := "Anonymous"
		if topic.AuthorName.Valid {
			authorName = topic.AuthorName.String
		}

		categoryName := ""
		if topic.CategoryName.Valid && topic.CategoryName.String != "" {
			categoryName = fmt.Sprintf(`<span class="forum-topic-category">%s</span>`, escapeHTML(topic.CategoryName.String))
		}

		avatar := forumGetInitials(authorName)

		html += fmt.Sprintf(`
			<div class="forum-topic%s">
				<div class="forum-topic-avatar">%s</div>
				<div class="forum-topic-main">
					<h4 class="forum-topic-title">
						<a href="/forum/topic/%s/">%s</a>
						%s
					</h4>
					<div class="forum-topic-meta">
						<span>%s</span>
						%s
						<span>%s</span>
					</div>
				</div>
				<div class="forum-topic-stats">
					<div class="forum-topic-stat">
						<span class="forum-topic-stat-value">%d</span>
						<span class="forum-topic-stat-label">Replies</span>
					</div>
					<div class="forum-topic-stat">
						<span class="forum-topic-stat-value">%d</span>
						<span class="forum-topic-stat-label">Views</span>
					</div>
				</div>
			</div>`,
			pinnedClass,
			avatar,
			escapeHTML(topic.Slug),
			escapeHTML(topic.Title),
			badges,
			escapeHTML(authorName),
			categoryName,
			forumFormatTimeAgo(topic.CreatedAt),
			topic.PostCount,
			topic.ViewCount,
		)
	}

	respondHTML(c, html)
}

// Forum search results fragment (dropdown)
func (r *Router) fragmentForumSearchDropdown(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	query := c.Query("q")

	if len(query) < 2 {
		respondHTML(c, "")
		return
	}

	topics, _, err := r.db.SearchForum(c.Request.Context(), siteID, query, 5, 0)
	if err != nil {
		respondHTMLError(c, "Search failed")
		return
	}

	if len(topics) == 0 {
		respondHTML(c, `<div class="forum-search-empty">No results found</div>`)
		return
	}

	html := ""
	for _, topic := range topics {
		authorName := "Anonymous"
		if topic.AuthorName.Valid {
			authorName = topic.AuthorName.String
		}

		html += fmt.Sprintf(`
			<a href="/forum/topic/%s/" class="forum-search-result">
				<span class="forum-search-result-title">%s</span>
				<span class="forum-search-result-meta">%s - %d replies</span>
			</a>`,
			escapeHTML(topic.Slug),
			escapeHTML(topic.Title),
			escapeHTML(authorName),
			topic.PostCount,
		)
	}

	respondHTML(c, html)
}

// Forum leaderboard fragment
func (r *Router) fragmentForumLeaderboard(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)

	stats, err := r.db.GetForumLeaderboard(c.Request.Context(), siteID, 5)
	if err != nil {
		respondHTMLError(c, "Failed to load leaderboard")
		return
	}

	if len(stats) == 0 {
		respondHTML(c, `<div class="forum-empty" style="padding: 0.5rem 0;">No contributors yet</div>`)
		return
	}

	html := ""
	for i, user := range stats {
		rankClass := ""
		switch i {
		case 0:
			rankClass = " forum-leaderboard-rank--gold"
		case 1:
			rankClass = " forum-leaderboard-rank--silver"
		case 2:
			rankClass = " forum-leaderboard-rank--bronze"
		}

		userName := "User"
		if user.UserName.Valid {
			userName = user.UserName.String
		}

		html += fmt.Sprintf(`
			<div class="forum-leaderboard-item">
				<span class="forum-leaderboard-rank%s">%d</span>
				<span class="forum-leaderboard-user">%s</span>
				<span class="forum-leaderboard-reputation">%d pts</span>
			</div>`,
			rankClass,
			i+1,
			escapeHTML(userName),
			user.Reputation,
		)
	}

	respondHTML(c, html)
}

// Forum tags fragment
func (r *Router) fragmentForumTags(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)

	tags, err := r.db.ListForumTags(c.Request.Context(), siteID)
	if err != nil {
		respondHTMLError(c, "Failed to load tags")
		return
	}

	if len(tags) == 0 {
		respondHTML(c, `<div class="forum-empty" style="padding: 0;">No tags yet</div>`)
		return
	}

	html := ""
	limit := 10
	if len(tags) < limit {
		limit = len(tags)
	}

	for i := 0; i < limit; i++ {
		tag := tags[i]
		style := ""
		if tag.Color.Valid && tag.Color.String != "" {
			style = fmt.Sprintf(` style="border-color:%s;color:%s"`, tag.Color.String, tag.Color.String)
		}
		html += fmt.Sprintf(`<a href="/forum/tag/%s/" class="forum-tag-btn"%s>%s (%d)</a>`,
			escapeHTML(tag.Slug),
			style,
			escapeHTML(tag.Name),
			tag.UsageCount,
		)
	}

	respondHTML(c, html)
}

// Forum stats fragment
func (r *Router) fragmentForumStats(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)

	_, topicCount, postCount, memberCount, err := r.db.GetForumStats(c.Request.Context(), siteID)
	if err != nil {
		respondHTML(c, `<span>--</span><span>--</span><span>--</span>`)
		return
	}

	html := fmt.Sprintf(`
		<div class="forum-hero-stat">
			<span class="forum-hero-stat-value">%d</span>
			<span>Topics</span>
		</div>
		<div class="forum-hero-stat">
			<span class="forum-hero-stat-value">%d</span>
			<span>Posts</span>
		</div>
		<div class="forum-hero-stat">
			<span class="forum-hero-stat-value">%d</span>
			<span>Members</span>
		</div>`,
		topicCount,
		postCount,
		memberCount,
	)

	respondHTML(c, html)
}

// Forum topic detail fragment (posts list)
func (r *Router) fragmentForumPosts(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	topicSlug := c.Param("slug")

	topic, err := r.db.GetForumTopicBySlug(c.Request.Context(), siteID, topicSlug)
	if err != nil || topic == nil {
		respondHTMLError(c, "Topic not found")
		return
	}

	posts, _, err := r.db.ListForumPosts(c.Request.Context(), topic.ID, 100, 0)
	if err != nil {
		respondHTMLError(c, "Failed to load posts")
		return
	}

	if len(posts) == 0 {
		respondHTML(c, `<p style="color:var(--text-muted);text-align:center;padding:2rem;">No replies yet. Be the first to respond!</p>`)
		return
	}

	// Check if user is authenticated
	userID, _ := getUserID(c)
	isAuthenticated := userID != ""

	html := ""
	for _, post := range posts {
		solutionClass := ""
		solutionBadge := ""
		if post.IsSolution {
			solutionClass = " forum-post-solution"
			solutionBadge = `<span class="forum-solution-badge"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg> Solution</span>`
		}

		authorName := "Anonymous"
		if post.AuthorName.Valid {
			authorName = post.AuthorName.String
		}

		avatar := forumGetInitials(authorName)

		// Render markdown content
		contentHTML := post.Content
		renderer := markdown.NewRenderer()
		if rendered, err := renderer.Render(post.Content); err == nil {
			contentHTML = rendered
		}

		// Build like action based on auth status
		var likeAction string
		if isAuthenticated {
			liked, _ := r.db.HasUserLikedPost(c.Request.Context(), userID, post.ID)
			likeClass := ""
			likeFill := "none"
			if liked {
				likeClass = " forum-action-btn--active"
				likeFill = "currentColor"
			}
			likeAction = fmt.Sprintf(`
				<button class="forum-action-btn%s" hx-post="/api/forum/posts/%s/like" hx-swap="outerHTML">
					<svg viewBox="0 0 24 24" fill="%s" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/></svg>
					<span>%d</span>
				</button>`, likeClass, post.ID, likeFill, post.LikeCount)
		} else {
			likeAction = fmt.Sprintf(`
				<span class="forum-stat-display">
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/></svg>
					<span>%d</span>
				</span>`, post.LikeCount)
		}

		html += fmt.Sprintf(`
			<div class="forum-post%s" id="post-%s">
				<div class="forum-post-header">
					<div class="forum-post-author">
						<div class="forum-avatar">%s</div>
						<div class="forum-post-author-info">
							<span class="forum-post-author-name">%s</span>
							<span class="forum-post-date">%s</span>
						</div>
					</div>
					%s
				</div>
				<div class="forum-post-body">%s</div>
				<div class="forum-post-actions">%s</div>
			</div>`,
			solutionClass,
			post.ID,
			avatar,
			escapeHTML(authorName),
			forumFormatTimeAgo(post.CreatedAt),
			solutionBadge,
			contentHTML,
			likeAction,
		)
	}

	respondHTML(c, html)
}

// Forum topic detail fragment (header, content, actions)
func (r *Router) fragmentForumTopicDetail(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	topicSlug := c.Param("slug")

	topic, err := r.db.GetForumTopicBySlug(c.Request.Context(), siteID, topicSlug)
	if err != nil || topic == nil {
		respondHTML(c, `<div class="forum-error">Topic not found</div>`)
		return
	}

	// Build badges
	badges := ""
	if topic.IsPinned {
		badges += `<span class="forum-topic-badge forum-topic-badge--pinned">Pinned</span>`
	}
	if topic.IsSolved {
		badges += `<span class="forum-topic-badge forum-topic-badge--solved">Solved</span>`
	}
	if topic.Status == "locked" {
		badges += `<span class="forum-topic-badge forum-topic-badge--locked">Locked</span>`
	}

	authorName := "Anonymous"
	if topic.AuthorName.Valid {
		authorName = topic.AuthorName.String
	}

	categoryLink := ""
	if topic.CategoryName.Valid && topic.CategoryName.String != "" {
		categorySlug := ""
		if topic.CategorySlug.Valid {
			categorySlug = topic.CategorySlug.String
		}
		categoryLink = fmt.Sprintf(` in <a href="/forum/category/%s/">%s</a>`,
			escapeHTML(categorySlug), escapeHTML(topic.CategoryName.String))
	}

	// Render content
	contentHTML := escapeHTML(topic.Content)
	renderer := markdown.NewRenderer()
	if rendered, err := renderer.Render(topic.Content); err == nil {
		contentHTML = rendered
	}

	html := `<a href="/forum/" class="forum-topic-back">&larr; Back to Forum</a>`
	html += `<header class="forum-topic-header">`
	html += fmt.Sprintf(`<h1 class="forum-topic-title">%s</h1>`, escapeHTML(topic.Title))
	if badges != "" {
		html += fmt.Sprintf(`<div class="forum-topic-badges">%s</div>`, badges)
	}
	html += `<div class="forum-topic-meta">`
	html += fmt.Sprintf(`<span>Posted by <strong>%s</strong></span>`, escapeHTML(authorName))
	html += fmt.Sprintf(`<span>%s</span>`, forumFormatTimeAgo(topic.CreatedAt))
	html += fmt.Sprintf(`<span>%d views</span>`, topic.ViewCount)
	html += categoryLink
	html += `</div></header>`

	html += fmt.Sprintf(`<div class="forum-topic-content">%s</div>`, contentHTML)

	// Actions - show like count always, but action buttons only for authenticated users
	html += `<div class="forum-topic-actions">`

	// Check if user is authenticated
	userID, _ := getUserID(c)
	isAuthenticated := userID != ""

	if isAuthenticated {
		// Check current like/bookmark status
		liked, _ := r.db.HasUserLikedTopic(c.Request.Context(), userID, topic.ID)
		bookmarked, _ := r.db.HasUserBookmarkedTopic(c.Request.Context(), userID, topic.ID)

		likeClass := ""
		likeFill := "none"
		if liked {
			likeClass = " forum-action-btn--active"
			likeFill = "currentColor"
		}
		html += fmt.Sprintf(`
			<button class="forum-action-btn%s" hx-post="/api/forum/topics/%s/like" hx-swap="outerHTML">
				<svg viewBox="0 0 24 24" fill="%s" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/></svg>
				<span>%d</span> Likes
			</button>`, likeClass, topic.ID, likeFill, topic.LikeCount)

		bookmarkClass := ""
		bookmarkFill := "none"
		bookmarkText := "Bookmark"
		if bookmarked {
			bookmarkClass = " forum-action-btn--active"
			bookmarkFill = "currentColor"
			bookmarkText = "Bookmarked"
		}
		html += fmt.Sprintf(`
			<button class="forum-action-btn%s" hx-post="/api/forum/topics/%s/bookmark" hx-swap="outerHTML">
				<svg viewBox="0 0 24 24" fill="%s" stroke="currentColor" stroke-width="2"><path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"/></svg> %s
			</button>`, bookmarkClass, topic.ID, bookmarkFill, bookmarkText)

		html += `
			<button class="forum-action-btn" onclick="document.getElementById('reply-form').scrollIntoView({behavior:'smooth'})">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg> Reply
			</button>`
	} else {
		// Show read-only stats for unauthenticated users
		html += fmt.Sprintf(`
			<span class="forum-stat-display">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/></svg>
				<span>%d</span> Likes
			</span>`, topic.LikeCount)
	}
	html += `</div>`

	// Update page title
	html += fmt.Sprintf(`<script>document.title = %q;</script>`, escapeHTML(topic.Title)+" | Forum")

	respondHTML(c, html)
}

// Forum category options fragment (for new topic form) - returns just options
func (r *Router) fragmentForumCategoryOptions(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)

	categories, err := r.db.ListForumCategories(c.Request.Context(), siteID)
	if err != nil {
		respondHTML(c, `<option value="">Failed to load categories</option>`)
		return
	}

	html := `<option value="">Select a category (optional)</option>`
	for _, cat := range categories {
		html += fmt.Sprintf(`<option value="%s">%s</option>`, cat.ID, escapeHTML(cat.Name))
	}

	respondHTML(c, html)
}

// Forum category select fragment (for new topic form) - returns full select element
func (r *Router) fragmentForumCategorySelect(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)

	categories, err := r.db.ListForumCategories(c.Request.Context(), siteID)
	if err != nil {
		respondHTML(c, `<select id="category_id" name="category_id" required><option value="">Failed to load</option></select>`)
		return
	}

	html := `<select id="category_id" name="category_id" required>`
	html += `<option value="">Select a category</option>`
	for _, cat := range categories {
		html += fmt.Sprintf(`<option value="%s">%s</option>`, cat.ID, escapeHTML(cat.Name))
	}
	html += `</select>`

	respondHTML(c, html)
}

// Forum category header fragment (for category page)
func (r *Router) fragmentForumCategoryHeader(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	slug := c.Param("slug")

	category, err := r.db.GetForumCategoryBySlug(c.Request.Context(), siteID, slug)
	if err != nil || category == nil {
		respondHTML(c, `<div class="forum-error">Category not found</div>`)
		return
	}

	html := fmt.Sprintf(`<h1 class="forum-category-title">%s</h1>`, escapeHTML(category.Name))
	if category.Description.Valid && category.Description.String != "" {
		html += fmt.Sprintf(`<p class="forum-category-description">%s</p>`, escapeHTML(category.Description.String))
	}
	html += fmt.Sprintf(`<p class="forum-category-stats">%d topics</p>`, category.TopicCount)

	respondHTML(c, html)
}

// Forum tag buttons fragment (for new topic form)
func (r *Router) fragmentForumTagButtons(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)

	tags, err := r.db.ListForumTags(c.Request.Context(), siteID)
	if err != nil || len(tags) == 0 {
		respondHTML(c, `<span style="color: var(--text-muted); font-size: 0.875rem;">No tags available</span>`)
		return
	}

	html := ""
	for _, tag := range tags {
		style := ""
		if tag.Color.Valid && tag.Color.String != "" {
			style = fmt.Sprintf(` style="border-color:%s"`, tag.Color.String)
		}
		html += fmt.Sprintf(`<button type="button" class="forum-tag-btn" data-tag="%s" onclick="toggleTag(this)"%s>%s</button>`,
			escapeHTML(tag.Slug),
			style,
			escapeHTML(tag.Name),
		)
	}

	respondHTML(c, html)
}

// Forum reply form fragment
func (r *Router) fragmentForumReplyForm(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	topicSlug := c.Param("slug")

	topic, err := r.db.GetForumTopicBySlug(c.Request.Context(), siteID, topicSlug)
	if err != nil || topic == nil {
		respondHTMLError(c, "Topic not found")
		return
	}

	if topic.Status == "locked" || topic.Status == "closed" {
		respondHTML(c, `<div class="forum-login-notice">This topic is closed for new replies.</div>`)
		return
	}

	// Check if user is authenticated
	userID, _ := getUserID(c)
	if userID == "" {
		topicPageURL := fmt.Sprintf("/forum/topic/%s/", topicSlug)
		loginURL := fmt.Sprintf("/login?site_id=%s&redirect=%s", siteID, topicPageURL)
		respondHTML(c, fmt.Sprintf(`<div class="forum-login-notice"><a href="%s">Sign in</a> to reply to this topic.</div>`, loginURL))
		return
	}

	html := fmt.Sprintf(`
		<form id="reply-form" class="forum-reply-form" hx-post="/api/forum/topics/by-slug/%s/posts" hx-ext="json-enc" hx-target="#forum-posts" hx-swap="beforeend" hx-on::after-request="this.reset()">
			<h3>Post a Reply</h3>
			<textarea name="content" class="forum-reply-textarea" placeholder="Write your reply... Markdown is supported." required></textarea>
			<div class="forum-reply-actions">
				<button type="submit" class="forum-btn forum-btn-primary">Post Reply</button>
			</div>
		</form>`,
		escapeHTML(topicSlug),
	)

	respondHTML(c, html)
}

// Helper to get initials from name
func forumGetInitials(name string) string {
	if name == "" {
		return "?"
	}
	initials := ""
	words := forumSplitWords(name)
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
	return initials
}

func forumSplitWords(s string) []string {
	var words []string
	word := ""
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' {
			if word != "" {
				words = append(words, word)
				word = ""
			}
		} else {
			word += string(r)
		}
	}
	if word != "" {
		words = append(words, word)
	}
	return words
}

// Format time as relative string
func forumFormatTimeAgo(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// Try other formats
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
