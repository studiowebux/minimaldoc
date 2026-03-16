package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/studiowebux/minimaldoc/internal/server/markdown"
	"github.com/studiowebux/minimaldoc/internal/server/store"
)

// Request types

type ForumCategoryRequest struct {
	Slug        string `json:"slug" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
	Position    int    `json:"position"`
	ParentID    string `json:"parent_id"`
	IsLocked    bool   `json:"is_locked"`
}

type ForumTopicRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	CategoryID string `json:"category_id"`
	Tags       string `json:"tags"` // comma-separated
}

type ForumPostRequest struct {
	Content  string `json:"content" binding:"required"`
	ParentID string `json:"parent_id"`
}

type ForumFlagRequest struct {
	TopicID     string `json:"topic_id"`
	PostID      string `json:"post_id"`
	Reason      string `json:"reason" binding:"required"`
	Description string `json:"description"`
}

type ForumBanRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	Reason      string `json:"reason" binding:"required"`
	ExpiresAt   string `json:"expires_at"`
	IsPermanent bool   `json:"is_permanent"`
}

type ForumTagRequest struct {
	Slug        string `json:"slug" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// Response types

type ForumTopicResponse struct {
	ID           string   `json:"id"`
	Slug         string   `json:"slug"`
	Title        string   `json:"title"`
	Content      string   `json:"content,omitempty"`
	ContentHTML  string   `json:"content_html,omitempty"`
	Status       string   `json:"status"`
	IsPinned     bool     `json:"is_pinned"`
	IsSolved     bool     `json:"is_solved"`
	ViewCount    int64    `json:"view_count"`
	LikeCount    int64    `json:"like_count"`
	PostCount    int64    `json:"post_count"`
	CategoryID   string   `json:"category_id,omitempty"`
	CategoryName string   `json:"category_name,omitempty"`
	CategorySlug string   `json:"category_slug,omitempty"`
	AuthorID     string   `json:"author_id,omitempty"`
	AuthorName   string   `json:"author_name,omitempty"`
	AuthorAvatar string   `json:"author_avatar,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	LastPostAt   string   `json:"last_post_at,omitempty"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
	IsLiked      bool     `json:"is_liked,omitempty"`
	IsBookmarked bool     `json:"is_bookmarked,omitempty"`
}

type ForumPostResponse struct {
	ID           string `json:"id"`
	TopicID      string `json:"topic_id"`
	Content      string `json:"content"`
	ContentHTML  string `json:"content_html"`
	LikeCount    int64  `json:"like_count"`
	IsSolution   bool   `json:"is_solution"`
	ParentID     string `json:"parent_id,omitempty"`
	AuthorID     string `json:"author_id,omitempty"`
	AuthorName   string `json:"author_name,omitempty"`
	AuthorAvatar string `json:"author_avatar,omitempty"`
	EditedAt     string `json:"edited_at,omitempty"`
	EditorName   string `json:"editor_name,omitempty"`
	CreatedAt    string `json:"created_at"`
	IsLiked      bool   `json:"is_liked,omitempty"`
}

type ForumCategoryResponse struct {
	ID          string                  `json:"id"`
	Slug        string                  `json:"slug"`
	Name        string                  `json:"name"`
	Description string                  `json:"description,omitempty"`
	Color       string                  `json:"color,omitempty"`
	Icon        string                  `json:"icon,omitempty"`
	Position    int                     `json:"position"`
	IsLocked    bool                    `json:"is_locked"`
	TopicCount  int64                   `json:"topic_count"`
	PostCount   int64                   `json:"post_count"`
	ParentID    string                  `json:"parent_id,omitempty"`
	Children    []ForumCategoryResponse `json:"children,omitempty"`
	CreatedAt   string                  `json:"created_at"`
}

type ForumTagResponse struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	TopicCount  int64  `json:"topic_count"`
}

type ForumUserStatsResponse struct {
	UserID     string `json:"user_id"`
	UserName   string `json:"user_name"`
	UserAvatar string `json:"user_avatar,omitempty"`
	Reputation int64  `json:"reputation"`
	TopicCount int64  `json:"topic_count"`
	PostCount  int64  `json:"post_count"`
	LikesGiven int64  `json:"likes_given"`
	LikesRecv  int64  `json:"likes_received"`
}

// Helper functions

func forumTopicToResponse(t *store.ForumTopic, includeContent bool) ForumTopicResponse {
	r := ForumTopicResponse{
		ID:           t.ID,
		Slug:         t.Slug,
		Title:        t.Title,
		Status:       t.Status,
		IsPinned:     t.IsPinned,
		IsSolved:     t.IsSolved,
		ViewCount:    t.ViewCount,
		LikeCount:    t.LikeCount,
		PostCount:    t.PostCount,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
		IsLiked:      t.IsLiked,
		IsBookmarked: t.IsBookmarked,
	}

	if includeContent {
		r.Content = t.Content
		renderer := markdown.NewRenderer()
		html, _ := renderer.Render(t.Content)
		r.ContentHTML = html
	}

	if t.CategoryID.Valid {
		r.CategoryID = t.CategoryID.String
	}
	if t.CategoryName.Valid {
		r.CategoryName = t.CategoryName.String
	}
	if t.CategorySlug.Valid {
		r.CategorySlug = t.CategorySlug.String
	}
	if t.AuthorID.Valid {
		r.AuthorID = t.AuthorID.String
	}
	if t.AuthorName.Valid {
		r.AuthorName = t.AuthorName.String
	}
	if t.AuthorAvatar.Valid {
		r.AuthorAvatar = t.AuthorAvatar.String
	}
	if t.LastPostAt.Valid {
		r.LastPostAt = t.LastPostAt.String
	}

	// Convert tags
	for _, tag := range t.Tags {
		r.Tags = append(r.Tags, tag.Name)
	}

	return r
}

func forumPostToResponse(p *store.ForumPost) ForumPostResponse {
	r := ForumPostResponse{
		ID:         p.ID,
		TopicID:    p.TopicID,
		Content:    p.Content,
		LikeCount:  p.LikeCount,
		IsSolution: p.IsSolution,
		CreatedAt:  p.CreatedAt,
		IsLiked:    p.IsLiked,
	}

	// Render markdown
	renderer := markdown.NewRenderer()
	html, _ := renderer.Render(p.Content)
	r.ContentHTML = html

	if p.ParentID.Valid {
		r.ParentID = p.ParentID.String
	}
	if p.AuthorID.Valid {
		r.AuthorID = p.AuthorID.String
	}
	if p.AuthorName.Valid {
		r.AuthorName = p.AuthorName.String
	}
	if p.AuthorAvatar.Valid {
		r.AuthorAvatar = p.AuthorAvatar.String
	}
	if p.EditedAt.Valid {
		r.EditedAt = p.EditedAt.String
	}
	if p.EditorName.Valid {
		r.EditorName = p.EditorName.String
	}

	return r
}

func generateTopicSlug(title string) string {
	slug := slugify(title)
	if len(slug) > 60 {
		slug = slug[:60]
	}
	return slug + "-" + uuid.New().String()[:8]
}

func forumCategoryToResponse(c *store.ForumCategory) ForumCategoryResponse {
	r := ForumCategoryResponse{
		ID:         c.ID,
		Slug:       c.Slug,
		Name:       c.Name,
		Position:   c.Position,
		IsLocked:   c.IsLocked,
		TopicCount: c.TopicCount,
		PostCount:  c.PostCount,
		CreatedAt:  c.CreatedAt,
	}
	if c.Description.Valid {
		r.Description = c.Description.String
	}
	if c.Color.Valid {
		r.Color = c.Color.String
	}
	if c.Icon.Valid {
		r.Icon = c.Icon.String
	}
	if c.ParentID.Valid {
		r.ParentID = c.ParentID.String
	}
	for _, child := range c.Children {
		r.Children = append(r.Children, forumCategoryToResponse(&child))
	}
	return r
}

func forumTagToResponse(t *store.ForumTag) ForumTagResponse {
	r := ForumTagResponse{
		ID:         t.ID,
		Slug:       t.Slug,
		Name:       t.Name,
		TopicCount: t.UsageCount,
	}
	if t.Description.Valid {
		r.Description = t.Description.String
	}
	if t.Color.Valid {
		r.Color = t.Color.String
	}
	return r
}

func forumUserStatsToResponse(s *store.ForumUserStats) ForumUserStatsResponse {
	r := ForumUserStatsResponse{
		UserID:     s.UserID,
		Reputation: s.Reputation,
		TopicCount: s.TopicCount,
		PostCount:  s.PostCount,
		LikesGiven: s.LikeGivenCount,
		LikesRecv:  s.LikeReceivedCount,
	}
	if s.UserName.Valid {
		r.UserName = s.UserName.String
	}
	if s.UserAvatar.Valid {
		r.UserAvatar = s.UserAvatar.String
	}
	return r
}

// Public API handlers

func (r *Router) listForumCategories(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	if siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}

	// Filter categories by visibility based on auth status
	var categories []store.ForumCategory
	var err error
	_, userErr := getUserID(c)
	role, _ := getUserRole(c)
	if role == "admin" || role == "editor" {
		categories, err = r.db.ListForumCategories(c.Request.Context(), siteID)
	} else if userErr == nil {
		categories, err = r.db.ListForumCategories(c.Request.Context(), siteID, "public", "members_only")
	} else {
		categories, err = r.db.ListForumCategories(c.Request.Context(), siteID, "public")
	}
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	var resp []ForumCategoryResponse
	for _, cat := range categories {
		resp = append(resp, forumCategoryToResponse(&cat))
	}
	c.JSON(http.StatusOK, gin.H{"categories": resp})
}

func (r *Router) getForumCategory(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	if siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}
	slug := c.Param("slug")

	category, err := r.db.GetForumCategoryBySlug(c.Request.Context(), siteID, slug)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if category == nil {
		respondNotFound(c, ErrCategoryNotFound, "category not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

func (r *Router) listForumTopics(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	if siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}

	categoryID := c.Query("category_id")
	status := c.Query("status")
	search := c.Query("q")

	limit, offset := parsePagination(c)

	topics, total, err := r.db.ListForumTopics(c.Request.Context(), siteID, categoryID, status, search, limit, offset)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	// Check if user is authenticated to add like/bookmark status
	userID, _ := getUserID(c)

	topicsResp := make([]ForumTopicResponse, 0, len(topics))
	for i := range topics {
		t := &topics[i]
		if userID != "" {
			t.IsLiked, _ = r.db.HasUserLikedTopic(c.Request.Context(), userID, t.ID)
			t.IsBookmarked, _ = r.db.HasUserBookmarkedTopic(c.Request.Context(), userID, t.ID)
		}
		// Load tags
		tags, _ := r.db.GetTopicTags(c.Request.Context(), t.ID)
		t.Tags = tags
		topicsResp = append(topicsResp, forumTopicToResponse(t, false))
	}

	c.JSON(http.StatusOK, gin.H{
		"topics": topicsResp,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (r *Router) getForumTopic(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	if siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}
	slug := c.Param("slug")

	topic, err := r.db.GetForumTopicBySlug(c.Request.Context(), siteID, slug)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if topic == nil {
		respondNotFound(c, ErrTopicNotFound, "topic not found")
		return
	}

	// Increment view count
	_ = r.db.IncrementForumTopicViews(c.Request.Context(), topic.ID)

	// Check user interactions
	userID, _ := getUserID(c)
	if userID != "" {
		topic.IsLiked, _ = r.db.HasUserLikedTopic(c.Request.Context(), userID, topic.ID)
		topic.IsBookmarked, _ = r.db.HasUserBookmarkedTopic(c.Request.Context(), userID, topic.ID)
	}

	// Load tags
	tags, _ := r.db.GetTopicTags(c.Request.Context(), topic.ID)
	topic.Tags = tags

	c.JSON(http.StatusOK, gin.H{"topic": forumTopicToResponse(topic, true)})
}

func (r *Router) getForumTopicPosts(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	if siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}
	slug := c.Param("slug")

	topic, err := r.db.GetForumTopicBySlug(c.Request.Context(), siteID, slug)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if topic == nil {
		respondNotFound(c, ErrTopicNotFound, "topic not found")
		return
	}

	limit, offset := parsePagination(c)

	posts, total, err := r.db.ListForumPosts(c.Request.Context(), topic.ID, limit, offset)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	userID, _ := getUserID(c)

	postsResp := make([]ForumPostResponse, 0, len(posts))
	for i := range posts {
		p := &posts[i]
		if userID != "" {
			p.IsLiked, _ = r.db.HasUserLikedPost(c.Request.Context(), userID, p.ID)
		}
		postsResp = append(postsResp, forumPostToResponse(p))
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":  postsResp,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (r *Router) searchForum(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	if siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}

	query := c.Query("q")
	if query == "" {
		respondBadRequest(c, ErrMissingParams, "search query required")
		return
	}

	limit, offset := parsePagination(c)

	topics, total, err := r.db.SearchForum(c.Request.Context(), siteID, query, limit, offset)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	topicsResp := make([]ForumTopicResponse, 0, len(topics))
	for i := range topics {
		topicsResp = append(topicsResp, forumTopicToResponse(&topics[i], false))
	}

	c.JSON(http.StatusOK, gin.H{
		"topics": topicsResp,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (r *Router) listForumTags(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	if siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}

	tags, err := r.db.ListForumTags(c.Request.Context(), siteID)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	var resp []ForumTagResponse
	for _, tag := range tags {
		resp = append(resp, forumTagToResponse(&tag))
	}
	c.JSON(http.StatusOK, gin.H{"tags": resp})
}

func (r *Router) getForumTopicsByTag(c *gin.Context) {
	siteID := r.getSiteIDWithFallback(c)
	if siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}
	slug := c.Param("slug")

	tag, err := r.db.GetForumTagBySlug(c.Request.Context(), siteID, slug)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if tag == nil {
		respondNotFound(c, ErrTagNotFound, "tag not found")
		return
	}

	limit, offset := parsePagination(c)

	topics, total, err := r.db.ListForumTopicsByTag(c.Request.Context(), siteID, tag.ID, limit, offset)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	topicsResp := make([]ForumTopicResponse, 0, len(topics))
	for i := range topics {
		topicsResp = append(topicsResp, forumTopicToResponse(&topics[i], false))
	}

	c.JSON(http.StatusOK, gin.H{
		"tag":    tag,
		"topics": topicsResp,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// Authenticated API handlers

func (r *Router) createForumTopic(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil || siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		respondUnauthorized(c, ErrAuthRequired, "authentication required")
		return
	}

	// Check if user is banned
	banned, err := r.db.IsUserBanned(c.Request.Context(), siteID, userID)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if banned {
		respondError(c, http.StatusForbidden, ErrBanned, "you are banned from the forum")
		return
	}

	// Check rate limit
	count, _ := r.db.CountUserTopicsToday(c.Request.Context(), siteID, userID)
	if count >= int64(r.config.Forum.MaxTopicsPerDay) {
		respondError(c, http.StatusTooManyRequests, ErrDailyLimitReached, "daily topic limit reached")
		return
	}

	var req ForumTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	// Sanitize content
	content := sanitizeComment(req.Content)
	if content == "" {
		respondBadRequest(c, ErrInvalidContent, "invalid content")
		return
	}

	id := uuid.New().String()
	slug := generateTopicSlug(req.Title)

	// Determine status based on moderation config
	status := "published"
	if r.config.Forum.ModerationMode == "all" {
		status = "pending"
	} else if r.config.Forum.ModerationMode == "first_post" {
		// Check if user has any approved topics
		count, _ := r.db.CountUserApprovedTopics(c.Request.Context(), siteID, userID)
		if count == 0 {
			status = "pending"
		}
	}

	topic, err := r.db.CreateForumTopic(c.Request.Context(), id, siteID, req.CategoryID, userID, slug, req.Title, content, status)
	if err != nil {
		respondInternalError(c, ErrTopicCreationFailed, "failed to create topic")
		return
	}

	// Add tags (comma-separated string)
	if req.Tags != "" {
		tagSlugs := strings.Split(req.Tags, ",")
		for _, tagSlug := range tagSlugs {
			tagSlug = strings.TrimSpace(tagSlug)
			if tagSlug == "" {
				continue
			}
			tag, _ := r.db.GetForumTagBySlug(c.Request.Context(), siteID, tagSlug)
			if tag != nil {
				_ = r.db.AddTagToTopic(c.Request.Context(), topic.ID, tag.ID)
			}
		}
	}

	// Update user stats
	statsID := uuid.New().String()
	_, _ = r.db.GetOrCreateForumUserStats(c.Request.Context(), statsID, siteID, userID)
	_ = r.db.IncrementForumUserStats(c.Request.Context(), siteID, userID, "topic_count", 1)

	// Add reputation
	if r.config.Forum.ReputationEnabled {
		_ = r.db.IncrementForumUserStats(c.Request.Context(), siteID, userID, "reputation", r.config.Forum.RepTopicCreate)
		logID := uuid.New().String()
		_ = r.db.LogReputationChange(c.Request.Context(), logID, siteID, userID, "topic_create", r.config.Forum.RepTopicCreate, topic.ID, "", "", "")
	}

	response := gin.H{"topic": forumTopicToResponse(topic, true)}
	if status == "pending" {
		response["message"] = "Your topic has been submitted and is awaiting moderation."
		response["slug"] = "" // Don't redirect to topic page if pending
	} else {
		response["slug"] = topic.Slug
	}
	c.JSON(http.StatusCreated, response)
}

func (r *Router) createForumPost(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil || siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}
	slug := c.Param("slug")

	userID, err := getUserID(c)
	if err != nil {
		respondUnauthorized(c, ErrAuthRequired, "authentication required")
		return
	}

	// Check if user is banned
	banned, err := r.db.IsUserBanned(c.Request.Context(), siteID, userID)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if banned {
		respondError(c, http.StatusForbidden, ErrBanned, "you are banned from the forum")
		return
	}

	topic, err := r.db.GetForumTopicBySlug(c.Request.Context(), siteID, slug)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if topic == nil {
		respondNotFound(c, ErrTopicNotFound, "topic not found")
		return
	}

	if topic.Status == "locked" || topic.Status == "closed" {
		respondError(c, http.StatusForbidden, ErrTopicClosed, "topic is closed for new replies")
		return
	}

	// Check rate limit
	count, _ := r.db.CountUserPostsToday(c.Request.Context(), siteID, userID)
	if count >= int64(r.config.Forum.MaxPostsPerDay) {
		respondError(c, http.StatusTooManyRequests, ErrDailyLimitReached, "daily post limit reached")
		return
	}

	var req ForumPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	content := sanitizeComment(req.Content)
	if content == "" {
		respondBadRequest(c, ErrInvalidContent, "invalid content")
		return
	}

	id := uuid.New().String()
	post, err := r.db.CreateForumPost(c.Request.Context(), id, siteID, topic.ID, req.ParentID, userID, content)
	if err != nil {
		respondInternalError(c, ErrInternalError, "failed to create post")
		return
	}

	// Update user stats
	statsID := uuid.New().String()
	_, _ = r.db.GetOrCreateForumUserStats(c.Request.Context(), statsID, siteID, userID)
	_ = r.db.IncrementForumUserStats(c.Request.Context(), siteID, userID, "post_count", 1)

	// Add reputation
	if r.config.Forum.ReputationEnabled {
		_ = r.db.IncrementForumUserStats(c.Request.Context(), siteID, userID, "reputation", r.config.Forum.RepPostCreate)
		logID := uuid.New().String()
		_ = r.db.LogReputationChange(c.Request.Context(), logID, siteID, userID, "post_create", r.config.Forum.RepPostCreate, topic.ID, post.ID, "", "")
	}

	// Create notification for topic author
	if topic.AuthorID.Valid && topic.AuthorID.String != userID {
		notifID := uuid.New().String()
		_ = r.db.CreateForumNotification(c.Request.Context(), notifID, siteID, topic.AuthorID.String, "reply", "New reply to your topic", req.Content[:min(100, len(req.Content))], topic.ID, post.ID, userID)
	}

	// Notify subscribers
	subscribers, _ := r.db.GetTopicSubscribers(c.Request.Context(), topic.ID)
	for _, subID := range subscribers {
		if subID != userID {
			notifID := uuid.New().String()
			_ = r.db.CreateForumNotification(c.Request.Context(), notifID, siteID, subID, "topic_update", "New reply in topic you're watching", req.Content[:min(100, len(req.Content))], topic.ID, post.ID, userID)
		}
	}

	// Return HTML for HTMX or JSON for API
	if c.GetHeader("HX-Request") == "true" {
		// Get author info
		authorName := "Anonymous"
		authorInitials := "?"
		user, _ := r.db.GetUserByID(c.Request.Context(), userID)
		if user != nil && user.Name.Valid && user.Name.String != "" {
			authorName = user.Name.String
			if len(authorName) > 0 {
				authorInitials = string([]rune(authorName)[0])
			}
		}

		// Render markdown
		contentHTML := escapeHTML(post.Content)
		mdRenderer := markdown.NewRenderer()
		if rendered, err := mdRenderer.Render(post.Content); err == nil {
			contentHTML = rendered
		}

		html := fmt.Sprintf(`
			<div class="forum-post" id="post-%s">
				<div class="forum-post-header">
					<div class="forum-post-avatar">%s</div>
					<div class="forum-post-author">
						<span class="forum-post-author-name">%s</span>
						<span class="forum-post-author-date">just now</span>
					</div>
				</div>
				<div class="forum-post-content">%s</div>
			</div>`,
			post.ID,
			authorInitials,
			escapeHTML(authorName),
			contentHTML,
		)
		respondHTML(c, html)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"post": forumPostToResponse(post)})
}

func (r *Router) updateForumTopic(c *gin.Context) {
	id := c.Param("id")

	userID, err := getUserID(c)
	if err != nil {
		respondUnauthorized(c, ErrAuthRequired, "authentication required")
		return
	}

	topic, err := r.db.GetForumTopicByID(c.Request.Context(), id)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if topic == nil {
		respondNotFound(c, ErrTopicNotFound, "topic not found")
		return
	}

	// Check ownership
	if !topic.AuthorID.Valid || topic.AuthorID.String != userID {
		respondError(c, http.StatusForbidden, ErrOwnPostsOnly, "can only edit own topics")
		return
	}

	var req ForumTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	content := sanitizeComment(req.Content)
	if err := r.db.UpdateForumTopic(c.Request.Context(), id, req.Title, content); err != nil {
		respondInternalError(c, ErrInternalError, "failed to update topic")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (r *Router) updateForumPost(c *gin.Context) {
	id := c.Param("id")

	userID, err := getUserID(c)
	if err != nil {
		respondUnauthorized(c, ErrAuthRequired, "authentication required")
		return
	}

	post, err := r.db.GetForumPostByID(c.Request.Context(), id)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if post == nil {
		respondNotFound(c, ErrNotFound, "post not found")
		return
	}

	// Check ownership
	if !post.AuthorID.Valid || post.AuthorID.String != userID {
		respondError(c, http.StatusForbidden, ErrOwnPostsOnly, "can only edit own posts")
		return
	}

	var req ForumPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	content := sanitizeComment(req.Content)
	if err := r.db.UpdateForumPost(c.Request.Context(), id, content, userID); err != nil {
		respondInternalError(c, ErrInternalError, "failed to update post")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (r *Router) deleteForumTopic(c *gin.Context) {
	id := c.Param("id")

	userID, err := getUserID(c)
	if err != nil {
		respondUnauthorized(c, ErrAuthRequired, "authentication required")
		return
	}

	topic, err := r.db.GetForumTopicByID(c.Request.Context(), id)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if topic == nil {
		respondNotFound(c, ErrTopicNotFound, "topic not found")
		return
	}

	// Check ownership
	if !topic.AuthorID.Valid || topic.AuthorID.String != userID {
		respondError(c, http.StatusForbidden, ErrOwnPostsOnly, "can only delete own topics")
		return
	}

	if err := r.db.DeleteForumTopic(c.Request.Context(), id); err != nil {
		respondInternalError(c, ErrInternalError, "failed to delete topic")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (r *Router) deleteForumPost(c *gin.Context) {
	id := c.Param("id")

	userID, err := getUserID(c)
	if err != nil {
		respondUnauthorized(c, ErrAuthRequired, "authentication required")
		return
	}

	post, err := r.db.GetForumPostByID(c.Request.Context(), id)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if post == nil {
		respondNotFound(c, ErrNotFound, "post not found")
		return
	}

	// Check ownership
	if !post.AuthorID.Valid || post.AuthorID.String != userID {
		respondError(c, http.StatusForbidden, ErrOwnPostsOnly, "can only delete own posts")
		return
	}

	if err := r.db.DeleteForumPost(c.Request.Context(), id); err != nil {
		respondInternalError(c, ErrInternalError, "failed to delete post")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (r *Router) likeForumTopic(c *gin.Context) {
	id := c.Param("id")
	siteID, err := getSiteID(c)
	if err != nil || siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		respondUnauthorized(c, ErrAuthRequired, "authentication required")
		return
	}

	topic, err := r.db.GetForumTopicByID(c.Request.Context(), id)
	if err != nil || topic == nil {
		respondNotFound(c, ErrTopicNotFound, "topic not found")
		return
	}

	// Check if already liked
	liked, _ := r.db.HasUserLikedTopic(c.Request.Context(), userID, id)
	if liked {
		// Unlike
		if err := r.db.UnlikeForumTopic(c.Request.Context(), userID, id); err != nil {
			respondInternalError(c, ErrInternalError, "failed to unlike")
			return
		}
		// Return HTML for HTMX or JSON for API
		if c.GetHeader("HX-Request") == "true" {
			newCount := topic.LikeCount - 1
			if newCount < 0 {
				newCount = 0
			}
			respondHTML(c, fmt.Sprintf(`<button class="forum-action-btn" hx-post="/api/forum/topics/%s/like" hx-swap="outerHTML">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/></svg>
				<span>%d</span> Likes
			</button>`, id, newCount))
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "unliked"})
		return
	}

	// Like
	likeID := uuid.New().String()
	if err := r.db.LikeForumTopic(c.Request.Context(), likeID, siteID, userID, id); err != nil {
		respondInternalError(c, ErrInternalError, "failed to like")
		return
	}

	// Update stats and reputation for topic author
	if r.config.Forum.ReputationEnabled && topic.AuthorID.Valid && topic.AuthorID.String != userID {
		_ = r.db.IncrementForumUserStats(c.Request.Context(), siteID, topic.AuthorID.String, "like_received_count", 1)
		_ = r.db.IncrementForumUserStats(c.Request.Context(), siteID, topic.AuthorID.String, "reputation", r.config.Forum.RepLikeReceived)
		logID := uuid.New().String()
		_ = r.db.LogReputationChange(c.Request.Context(), logID, siteID, topic.AuthorID.String, "like_received", r.config.Forum.RepLikeReceived, topic.ID, "", "", "")

		// Notification
		notifID := uuid.New().String()
		_ = r.db.CreateForumNotification(c.Request.Context(), notifID, siteID, topic.AuthorID.String, "like", "Someone liked your topic", "", topic.ID, "", userID)
	}

	_ = r.db.IncrementForumUserStats(c.Request.Context(), siteID, userID, "like_given_count", 1)

	// Return HTML for HTMX or JSON for API
	if c.GetHeader("HX-Request") == "true" {
		newCount := topic.LikeCount + 1
		respondHTML(c, fmt.Sprintf(`<button class="forum-action-btn forum-action-btn--active" hx-post="/api/forum/topics/%s/like" hx-swap="outerHTML">
			<svg viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/></svg>
			<span>%d</span> Likes
		</button>`, id, newCount))
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "liked"})
}

func (r *Router) likeForumPost(c *gin.Context) {
	id := c.Param("id")
	siteID, err := getSiteID(c)
	if err != nil || siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		respondUnauthorized(c, ErrAuthRequired, "authentication required")
		return
	}

	post, err := r.db.GetForumPostByID(c.Request.Context(), id)
	if err != nil || post == nil {
		respondNotFound(c, ErrNotFound, "post not found")
		return
	}

	isHTMX := c.GetHeader("HX-Request") == "true"
	var nowLiked bool

	// Check if already liked
	liked, _ := r.db.HasUserLikedPost(c.Request.Context(), userID, id)
	if liked {
		// Unlike
		if err := r.db.UnlikeForumPost(c.Request.Context(), userID, id); err != nil {
			respondInternalError(c, ErrInternalError, "failed to unlike")
			return
		}
		nowLiked = false
	} else {
		// Like
		likeID := uuid.New().String()
		if err := r.db.LikeForumPost(c.Request.Context(), likeID, siteID, userID, id); err != nil {
			respondInternalError(c, ErrInternalError, "failed to like")
			return
		}
		nowLiked = true

		// Update stats and reputation for post author
		if r.config.Forum.ReputationEnabled && post.AuthorID.Valid && post.AuthorID.String != userID {
			_ = r.db.IncrementForumUserStats(c.Request.Context(), siteID, post.AuthorID.String, "like_received_count", 1)
			_ = r.db.IncrementForumUserStats(c.Request.Context(), siteID, post.AuthorID.String, "reputation", r.config.Forum.RepLikeReceived)
			logID := uuid.New().String()
			_ = r.db.LogReputationChange(c.Request.Context(), logID, siteID, post.AuthorID.String, "like_received", r.config.Forum.RepLikeReceived, post.TopicID, post.ID, "", "")

			// Notification
			notifID := uuid.New().String()
			_ = r.db.CreateForumNotification(c.Request.Context(), notifID, siteID, post.AuthorID.String, "like", "Someone liked your post", "", post.TopicID, post.ID, userID)
		}

		_ = r.db.IncrementForumUserStats(c.Request.Context(), siteID, userID, "like_given_count", 1)
	}

	// Return HTML for HTMX requests
	if isHTMX {
		// Get updated like count
		updatedPost, _ := r.db.GetForumPostByID(c.Request.Context(), id)
		likeCount := int64(0)
		if updatedPost != nil {
			likeCount = updatedPost.LikeCount
		}

		activeClass := ""
		likeFill := "none"
		if nowLiked {
			activeClass = " forum-action-btn--active"
			likeFill = "currentColor"
		}

		html := fmt.Sprintf(`<button class="forum-action-btn%s" hx-post="/api/forum/posts/%s/like" hx-swap="outerHTML">
			<svg viewBox="0 0 24 24" fill="%s" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/></svg>
			<span>%d</span>
		</button>`, activeClass, id, likeFill, likeCount)

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
		return
	}

	if nowLiked {
		c.JSON(http.StatusOK, gin.H{"status": "liked"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "unliked"})
	}
}

func (r *Router) bookmarkForumTopic(c *gin.Context) {
	id := c.Param("id")
	siteID, err := getSiteID(c)
	if err != nil || siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		respondUnauthorized(c, ErrAuthRequired, "authentication required")
		return
	}

	// Check if already bookmarked
	bookmarked, _ := r.db.HasUserBookmarkedTopic(c.Request.Context(), userID, id)
	if bookmarked {
		// Remove bookmark
		if err := r.db.DeleteForumBookmark(c.Request.Context(), userID, id); err != nil {
			respondInternalError(c, ErrInternalError, "failed to remove bookmark")
			return
		}
		if c.GetHeader("HX-Request") == "true" {
			respondHTML(c, fmt.Sprintf(`<button class="forum-action-btn" hx-post="/api/forum/topics/%s/bookmark" hx-swap="outerHTML">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"/></svg> Bookmark
			</button>`, id))
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "unbookmarked"})
		return
	}

	// Add bookmark
	bookmarkID := uuid.New().String()
	if err := r.db.CreateForumBookmark(c.Request.Context(), bookmarkID, siteID, userID, id); err != nil {
		respondInternalError(c, ErrInternalError, "failed to bookmark")
		return
	}

	if c.GetHeader("HX-Request") == "true" {
		respondHTML(c, fmt.Sprintf(`<button class="forum-action-btn forum-action-btn--active" hx-post="/api/forum/topics/%s/bookmark" hx-swap="outerHTML">
			<svg viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2"><path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"/></svg> Bookmarked
		</button>`, id))
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "bookmarked"})
}

func (r *Router) flagForumContent(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil || siteID == "" {
		respondBadRequest(c, ErrSiteIDRequired, "site_id required")
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		respondUnauthorized(c, ErrAuthRequired, "authentication required")
		return
	}

	var req ForumFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	if req.TopicID == "" && req.PostID == "" {
		respondBadRequest(c, ErrMissingParams, "topic_id or post_id required")
		return
	}

	id := uuid.New().String()
	if err := r.db.CreateForumFlag(c.Request.Context(), id, siteID, userID, req.TopicID, req.PostID, req.Reason, req.Description); err != nil {
		respondInternalError(c, ErrInternalError, "failed to submit flag")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "flagged"})
}

func (r *Router) getForumNotifications(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		respondUnauthorized(c, ErrAuthRequired, "authentication required")
		return
	}

	unreadOnly := c.Query("unread") == "true"
	limit, offset := parsePagination(c)

	notifications, err := r.db.ListForumNotifications(c.Request.Context(), userID, unreadOnly, limit, offset)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	unreadCount, _ := r.db.GetUnreadNotificationCount(c.Request.Context(), userID)

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"unread_count":  unreadCount,
	})
}

func (r *Router) markNotificationRead(c *gin.Context) {
	id := c.Param("id")

	if err := r.db.MarkNotificationRead(c.Request.Context(), id); err != nil {
		respondInternalError(c, ErrInternalError, "failed to mark as read")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "read"})
}

func (r *Router) markAllNotificationsRead(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		respondUnauthorized(c, ErrAuthRequired, "authentication required")
		return
	}

	if err := r.db.MarkAllNotificationsRead(c.Request.Context(), userID); err != nil {
		respondInternalError(c, ErrInternalError, "failed to mark all as read")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "all_read"})
}

func (r *Router) getUserBookmarks(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		respondUnauthorized(c, ErrAuthRequired, "authentication required")
		return
	}

	limit, offset := parsePagination(c)

	bookmarks, err := r.db.ListUserBookmarks(c.Request.Context(), userID, limit, offset)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookmarks": bookmarks})
}

// Admin API handlers

func (r *Router) adminListForumCategories(c *gin.Context) {
	siteID, _ := getSiteID(c)

	categories, err := r.db.ListForumCategories(c.Request.Context(), siteID)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	var resp []ForumCategoryResponse
	for _, cat := range categories {
		resp = append(resp, forumCategoryToResponse(&cat))
	}
	c.JSON(http.StatusOK, gin.H{"categories": resp})
}

func (r *Router) adminGetForumCategory(c *gin.Context) {
	siteID, _ := getSiteID(c)
	id := c.Param("id")

	category, err := r.db.GetForumCategoryByID(c.Request.Context(), id)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if category == nil || category.SiteID != siteID {
		respondNotFound(c, ErrCategoryNotFound, "category not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

func (r *Router) adminCreateForumCategory(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondUnauthorized(c, ErrUnauthorized, "unauthorized")
		return
	}

	var req ForumCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	// Check slug uniqueness
	existing, err := r.db.GetForumCategoryBySlug(c.Request.Context(), siteID, req.Slug)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if existing != nil {
		respondError(c, http.StatusConflict, ErrConflict, "slug already exists")
		return
	}

	id := uuid.New().String()
	category, err := r.db.CreateForumCategory(c.Request.Context(), id, siteID, req.ParentID, req.Slug, req.Name, req.Description, req.Color, req.Icon, req.Position)
	if err != nil {
		respondInternalError(c, ErrInternalError, "failed to create category")
		return
	}

	// Audit log: forum category created
	r.logAuditAction(c, "create", "forum_category", id, req.Name, "")

	c.JSON(http.StatusCreated, gin.H{"category": category})
}

func (r *Router) adminUpdateForumCategory(c *gin.Context) {
	id := c.Param("id")

	var req ForumCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	if err := r.db.UpdateForumCategory(c.Request.Context(), id, req.Slug, req.Name, req.Description, req.Color, req.Icon, req.Position, req.IsLocked); err != nil {
		respondInternalError(c, ErrInternalError, "failed to update category")
		return
	}

	// Audit log: forum category updated
	r.logAuditAction(c, "update", "forum_category", id, req.Name, "")

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (r *Router) adminDeleteForumCategory(c *gin.Context) {
	id := c.Param("id")

	// Get category name for audit log before deletion
	category, _ := r.db.GetForumCategoryByID(c.Request.Context(), id)
	categoryName := ""
	if category != nil {
		categoryName = category.Name
	}

	if err := r.db.DeleteForumCategory(c.Request.Context(), id); err != nil {
		respondInternalError(c, ErrInternalError, "failed to delete category")
		return
	}

	// Audit log: forum category deleted
	r.logAuditAction(c, "delete", "forum_category", id, categoryName, "")

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (r *Router) adminPinForumTopic(c *gin.Context) {
	id := c.Param("id")

	topic, err := r.db.GetForumTopicByID(c.Request.Context(), id)
	if err != nil || topic == nil {
		respondNotFound(c, ErrTopicNotFound, "topic not found")
		return
	}

	// Toggle pin
	if err := r.db.PinForumTopic(c.Request.Context(), id, !topic.IsPinned); err != nil {
		respondInternalError(c, ErrInternalError, "failed to pin topic")
		return
	}

	status := "pinned"
	if topic.IsPinned {
		status = "unpinned"
	}
	c.JSON(http.StatusOK, gin.H{"status": status})
}

func (r *Router) adminLockForumTopic(c *gin.Context) {
	id := c.Param("id")

	if err := r.db.UpdateForumTopicStatus(c.Request.Context(), id, "locked"); err != nil {
		respondInternalError(c, ErrInternalError, "failed to lock topic")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "locked"})
}

func (r *Router) adminCloseForumTopic(c *gin.Context) {
	id := c.Param("id")

	if err := r.db.UpdateForumTopicStatus(c.Request.Context(), id, "closed"); err != nil {
		respondInternalError(c, ErrInternalError, "failed to close topic")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "closed"})
}

func (r *Router) adminOpenForumTopic(c *gin.Context) {
	id := c.Param("id")

	if err := r.db.UpdateForumTopicStatus(c.Request.Context(), id, "open"); err != nil {
		respondInternalError(c, ErrInternalError, "failed to open topic")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "opened"})
}

func (r *Router) adminMarkSolution(c *gin.Context) {
	postID := c.Param("id")
	siteID, _ := getSiteID(c)

	post, err := r.db.GetForumPostByID(c.Request.Context(), postID)
	if err != nil || post == nil {
		respondNotFound(c, ErrNotFound, "post not found")
		return
	}

	if err := r.db.MarkPostAsSolution(c.Request.Context(), postID, post.TopicID); err != nil {
		respondInternalError(c, ErrInternalError, "failed to mark solution")
		return
	}

	// Award reputation
	if r.config.Forum.ReputationEnabled && post.AuthorID.Valid {
		_ = r.db.IncrementForumUserStats(c.Request.Context(), siteID, post.AuthorID.String, "solution_count", 1)
		_ = r.db.IncrementForumUserStats(c.Request.Context(), siteID, post.AuthorID.String, "reputation", r.config.Forum.RepSolutionMarked)
		logID := uuid.New().String()
		_ = r.db.LogReputationChange(c.Request.Context(), logID, siteID, post.AuthorID.String, "solution_marked", r.config.Forum.RepSolutionMarked, post.TopicID, post.ID, "", "")

		// Notification
		notifID := uuid.New().String()
		_ = r.db.CreateForumNotification(c.Request.Context(), notifID, siteID, post.AuthorID.String, "solution", "Your answer was marked as the solution!", "", post.TopicID, post.ID, "")
	}

	c.JSON(http.StatusOK, gin.H{"status": "marked_solution"})
}

func (r *Router) adminDeleteForumTopic(c *gin.Context) {
	id := c.Param("id")

	// Get topic title for audit log before deletion
	topic, _ := r.db.GetForumTopicByID(c.Request.Context(), id)
	topicTitle := ""
	if topic != nil {
		topicTitle = topic.Title
	}

	if err := r.db.DeleteForumTopic(c.Request.Context(), id); err != nil {
		respondInternalError(c, ErrInternalError, "failed to delete topic")
		return
	}

	// Audit log: forum topic deleted
	r.logAuditAction(c, "delete", "forum_topic", id, topicTitle, "")

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (r *Router) adminDeleteForumPost(c *gin.Context) {
	id := c.Param("id")

	if err := r.db.DeleteForumPost(c.Request.Context(), id); err != nil {
		respondInternalError(c, ErrInternalError, "failed to delete post")
		return
	}

	// Audit log: forum post deleted
	r.logAuditAction(c, "delete", "forum_post", id, "", "")

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (r *Router) adminListForumTopics(c *gin.Context) {
	siteID, _ := getSiteID(c)

	categoryID := c.Query("category_id")
	status := c.Query("status")
	search := c.Query("q")

	limit, offset := parsePagination(c)

	topics, total, err := r.db.ListForumTopics(c.Request.Context(), siteID, categoryID, status, search, limit, offset)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	topicsResp := make([]ForumTopicResponse, 0, len(topics))
	for i := range topics {
		t := &topics[i]
		tags, _ := r.db.GetTopicTags(c.Request.Context(), t.ID)
		t.Tags = tags
		topicsResp = append(topicsResp, forumTopicToResponse(t, false))
	}

	c.JSON(http.StatusOK, gin.H{
		"topics": topicsResp,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (r *Router) adminListForumFlags(c *gin.Context) {
	siteID, _ := getSiteID(c)
	status := c.Query("status")

	limit, offset := parsePagination(c)

	flags, total, err := r.db.ListForumFlags(c.Request.Context(), siteID, status, limit, offset)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"flags": flags,
		"total": total,
	})
}

func (r *Router) adminResolveForumFlag(c *gin.Context) {
	id := c.Param("id")
	userID, _ := getUserID(c)

	var req struct {
		Status         string `json:"status" binding:"required"`
		ResolutionNote string `json:"resolution_note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	if err := r.db.ResolveForumFlag(c.Request.Context(), id, req.Status, userID, req.ResolutionNote); err != nil {
		respondInternalError(c, ErrInternalError, "failed to resolve flag")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "resolved"})
}

func (r *Router) adminBanUser(c *gin.Context) {
	siteID, _ := getSiteID(c)
	bannerID, _ := getUserID(c)

	var req ForumBanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	id := uuid.New().String()
	if err := r.db.CreateForumBan(c.Request.Context(), id, siteID, req.UserID, bannerID, req.Reason, req.ExpiresAt, req.IsPermanent); err != nil {
		respondInternalError(c, ErrInternalError, "failed to ban user")
		return
	}

	// Audit log: user banned
	r.logAuditAction(c, "create", "forum_ban", req.UserID, "", "reason: "+req.Reason)

	c.JSON(http.StatusOK, gin.H{"status": "banned"})
}

func (r *Router) adminUnbanUser(c *gin.Context) {
	siteID, _ := getSiteID(c)
	userID := c.Param("id")

	if err := r.db.DeleteForumBan(c.Request.Context(), siteID, userID); err != nil {
		respondInternalError(c, ErrInternalError, "failed to unban user")
		return
	}

	// Audit log: user unbanned
	r.logAuditAction(c, "delete", "forum_ban", userID, "", "")

	c.JSON(http.StatusOK, gin.H{"status": "unbanned"})
}

func (r *Router) adminListForumBans(c *gin.Context) {
	siteID, _ := getSiteID(c)

	limit, offset := parsePagination(c)

	bans, total, err := r.db.ListForumBans(c.Request.Context(), siteID, limit, offset)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bans":  bans,
		"total": total,
	})
}

func (r *Router) adminListForumTags(c *gin.Context) {
	siteID, _ := getSiteID(c)

	tags, err := r.db.ListForumTags(c.Request.Context(), siteID)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	var resp []ForumTagResponse
	for _, tag := range tags {
		resp = append(resp, forumTagToResponse(&tag))
	}
	c.JSON(http.StatusOK, gin.H{"tags": resp})
}

func (r *Router) adminGetForumTag(c *gin.Context) {
	siteID, _ := getSiteID(c)
	id := c.Param("id")

	tag, err := r.db.GetForumTagByID(c.Request.Context(), id)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if tag == nil || tag.SiteID != siteID {
		respondNotFound(c, ErrTagNotFound, "tag not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"tag": tag})
}

func (r *Router) adminCreateForumTag(c *gin.Context) {
	siteID, _ := getSiteID(c)

	var req ForumTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	// Check uniqueness
	existing, _ := r.db.GetForumTagBySlug(c.Request.Context(), siteID, req.Slug)
	if existing != nil {
		respondError(c, http.StatusConflict, ErrTagExists, "tag already exists")
		return
	}

	id := uuid.New().String()
	tag, err := r.db.CreateForumTag(c.Request.Context(), id, siteID, req.Slug, req.Name, req.Description, req.Color)
	if err != nil {
		respondInternalError(c, ErrInternalError, "failed to create tag")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"tag": tag})
}

func (r *Router) adminUpdateForumTag(c *gin.Context) {
	id := c.Param("id")

	var req ForumTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	if err := r.db.UpdateForumTag(c.Request.Context(), id, req.Slug, req.Name, req.Description, req.Color); err != nil {
		respondInternalError(c, ErrInternalError, "failed to update tag")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (r *Router) adminDeleteForumTag(c *gin.Context) {
	id := c.Param("id")

	if err := r.db.DeleteForumTag(c.Request.Context(), id); err != nil {
		respondInternalError(c, ErrInternalError, "failed to delete tag")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (r *Router) getForumStats(c *gin.Context) {
	siteID, _ := getSiteID(c)

	categories, topics, posts, users, err := r.db.GetForumStats(c.Request.Context(), siteID)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"topics":     topics,
		"posts":      posts,
		"users":      users,
	})
}

func (r *Router) getForumLeaderboard(c *gin.Context) {
	siteID := c.Query("site_id")
	if siteID == "" {
		siteID = c.GetHeader("X-API-Key")
	}
	if siteID == "" {
		siteID, _ = getSiteID(c)
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 50 {
			limit = v
		}
	}

	stats, err := r.db.GetForumLeaderboard(c.Request.Context(), siteID, limit)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	var resp []ForumUserStatsResponse
	for _, s := range stats {
		resp = append(resp, forumUserStatsToResponse(&s))
	}
	c.JSON(http.StatusOK, gin.H{"leaderboard": resp})
}

// Admin UI handlers

func (r *Router) adminForum(c *gin.Context) {
	claims, _ := getUserClaims(c)
	c.HTML(http.StatusOK, "forum.html", gin.H{
		"Title":       "Forum",
		"CurrentPage": "forum",
		"User":        claims,
		"Nonce":       cspNonce(c),
	})
}

func (r *Router) adminForumCategories(c *gin.Context) {
	claims, _ := getUserClaims(c)
	c.HTML(http.StatusOK, "forum-categories.html", gin.H{
		"Title":       "Forum Categories",
		"CurrentPage": "forum",
		"User":        claims,
		"Nonce":       cspNonce(c),
	})
}

func (r *Router) adminForumTopics(c *gin.Context) {
	claims, _ := getUserClaims(c)
	c.HTML(http.StatusOK, "forum-topics.html", gin.H{
		"Title":       "Forum Topics",
		"CurrentPage": "forum",
		"User":        claims,
		"Nonce":       cspNonce(c),
	})
}

func (r *Router) adminForumFlags(c *gin.Context) {
	claims, _ := getUserClaims(c)
	c.HTML(http.StatusOK, "forum-flags.html", gin.H{
		"Title":       "Forum Flags",
		"CurrentPage": "forum",
		"User":        claims,
		"Nonce":       cspNonce(c),
	})
}

func (r *Router) adminForumBans(c *gin.Context) {
	claims, _ := getUserClaims(c)
	c.HTML(http.StatusOK, "forum-bans.html", gin.H{
		"Title":       "Forum Bans",
		"CurrentPage": "forum",
		"User":        claims,
		"Nonce":       cspNonce(c),
	})
}

func (r *Router) adminForumTags(c *gin.Context) {
	claims, _ := getUserClaims(c)
	c.HTML(http.StatusOK, "forum-tags.html", gin.H{
		"Title":       "Forum Tags",
		"CurrentPage": "forum",
		"User":        claims,
		"Nonce":       cspNonce(c),
	})
}

// Helper function for slugifying
func slugify(s string) string {
	// Simple slug generation - lowercase, replace spaces with dashes
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	// Remove non-alphanumeric except dashes
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}
