// Package store — repository interface definitions.
// Pattern: repository interface per domain aggregate.
package store

import (
	"context"
	"time"
)

// Store composes all domain repositories with core database operations.
type Store interface {
	AnalyticsStore
	AuditStore
	BlogStore
	DocsStore
	EventStore
	FeedbackStore
	ForumStore
	NewsletterStore
	SessionStore
	SiteStore
	UploadStore
	UserStore

	// Core methods
	Ping(ctx context.Context) error
	Driver() string
	Close() error
}

// AnalyticsStore handles page view tracking and statistics.
type AnalyticsStore interface {
	RecordPageView(ctx context.Context, siteID, path, referrer, country, deviceType, browser, os, sessionHash string) error
	GetPageViewStats(ctx context.Context, siteID string, since time.Time) (totalViews int64, uniqueSessions int64, err error)
	GetPageViewStatsExtended(ctx context.Context, siteID string, since time.Time) (totalViews, uniqueSessions int64, avgDuration float64, err error)
	UpdatePageViewDuration(ctx context.Context, siteID, path, sessionHash string, duration int) error
	UpdatePageViewDurationAndBounce(ctx context.Context, siteID, path, sessionHash string, duration int, isBounce *bool) error
	GetBounceRate(ctx context.Context, siteID string, since time.Time) (float64, error)
	GetTopPages(ctx context.Context, siteID string, since time.Time, limit int) ([]struct {
		Path  string
		Views int64
	}, error)
	GetTrafficSources(ctx context.Context, siteID string, since time.Time, limit int) ([]TrafficSource, error)
	GetDailyViews(ctx context.Context, siteID string, days int) ([]DailyViewCount, error)
}

// AuditStore handles audit log operations.
type AuditStore interface {
	CreateAuditLog(ctx context.Context, id, siteID, userID, userEmail, action, entityType, entityID, entityName, details, ipAddress, userAgent string) error
	ListAuditLogs(ctx context.Context, siteID, action, entityType string, limit, offset int) ([]AuditLog, int, error)
	GetAuditLogStats(ctx context.Context, siteID string) (*AuditLogStats, error)
	GetAuditLogEntityTypes(ctx context.Context, siteID string) ([]string, error)
}

// BlogStore handles blog posts and comments.
type BlogStore interface {
	CreateBlogPost(ctx context.Context, id, siteID, authorID, slug, title, description, content, featuredImage, tags, category, visibility string) (*BlogPost, error)
	GetBlogPostByID(ctx context.Context, id string) (*BlogPost, error)
	GetBlogPostBySlug(ctx context.Context, siteID, slug string) (*BlogPost, error)
	UpdateBlogPost(ctx context.Context, id, slug, title, description, content, featuredImage, tags, category, visibility string) error
	DeleteBlogPost(ctx context.Context, id string) error
	ListBlogPosts(ctx context.Context, siteID string, status string, limit, offset int) ([]BlogPost, error)
	ListPublishedBlogPosts(ctx context.Context, siteID string, limit, offset int) ([]BlogPost, error)
	ListPublishedBlogPostsFiltered(ctx context.Context, siteID, category, tag, search string, limit, offset int) ([]BlogPost, int64, error)
	GetRelatedPosts(ctx context.Context, siteID, excludeID, category, tags string, limit int) ([]BlogPost, error)
	GetBlogPostStats(ctx context.Context, siteID string) (total, published, draft int64, err error)
	PublishBlogPost(ctx context.Context, id string) error
	UnpublishBlogPost(ctx context.Context, id string) error
	ScheduleBlogPost(ctx context.Context, id, scheduledAt string) error
	ClearSchedule(ctx context.Context, id string) error
	GetScheduledPosts(ctx context.Context) ([]BlogPost, error)
	PublishScheduledPosts(ctx context.Context) (int64, error)

	// Comments
	CreateBlogComment(ctx context.Context, id, siteID, postID, parentID, authorName, authorEmail, content, ipAddress, userAgent string) (*BlogComment, error)
	GetBlogCommentByID(ctx context.Context, id string) (*BlogComment, error)
	ListBlogComments(ctx context.Context, siteID string, status string, limit, offset int) ([]BlogComment, error)
	ListApprovedComments(ctx context.Context, postID string, limit, offset int) ([]BlogComment, error)
	ListPendingComments(ctx context.Context, siteID string, limit, offset int) ([]BlogComment, error)
	ModerateComment(ctx context.Context, id, status, moderatorID string) error
	DeleteBlogComment(ctx context.Context, id string) error
	GetCommentStats(ctx context.Context, siteID string) (total, pending, approved int64, err error)
}

// DocsStore handles private documentation access rules.
type DocsStore interface {
	CreateDocAccess(ctx context.Context, id, siteID, pathPattern, requiredRole, description string) (*DocAccess, error)
	GetDocAccessByID(ctx context.Context, id string) (*DocAccess, error)
	ListDocAccess(ctx context.Context, siteID string) ([]DocAccess, error)
	UpdateDocAccess(ctx context.Context, id, pathPattern, requiredRole, description string) error
	DeleteDocAccess(ctx context.Context, id string) error
	CheckDocAccess(ctx context.Context, siteID, path string) (*DocAccess, error)
}

// EventStore handles custom analytics events.
type EventStore interface {
	RecordEvent(ctx context.Context, siteID, name, category, path, value, sessionHash string) error
	GetEventStats(ctx context.Context, siteID string, since time.Time) ([]EventStat, error)
	GetEventsByName(ctx context.Context, siteID, name string, limit, offset int) ([]Event, error)
	ListRecentEvents(ctx context.Context, siteID string, limit int) ([]Event, error)
	GetTotalEventCount(ctx context.Context, siteID string, since time.Time) (int64, error)
	GetUniqueEventNames(ctx context.Context, siteID string, since time.Time) (int64, error)
}

// FeedbackStore handles user feedback and ratings.
type FeedbackStore interface {
	RecordRating(ctx context.Context, siteID, path string, rating int, feedback, sessionHash string) error
	GetRatingStats(ctx context.Context, siteID string) (avgRating float64, totalRatings int64, err error)
	GetRatingStatsExtended(ctx context.Context, siteID string) (avgRating float64, totalRatings, withComments, thisWeek int64, err error)
	GetRatingsByPath(ctx context.Context, siteID, path string) (avgRating float64, count int64, err error)
	ListRatings(ctx context.Context, siteID string, limit, offset int) ([]Rating, error)
}

// ForumStore handles all forum operations.
type ForumStore interface {
	// Categories
	CreateForumCategory(ctx context.Context, id, siteID, parentID, slug, name, description, color, icon string, position int) (*ForumCategory, error)
	GetForumCategoryByID(ctx context.Context, id string) (*ForumCategory, error)
	GetForumCategoryBySlug(ctx context.Context, siteID, slug string) (*ForumCategory, error)
	ListForumCategories(ctx context.Context, siteID string, visibilities ...string) ([]ForumCategory, error)
	UpdateForumCategory(ctx context.Context, id, slug, name, description, color, icon string, position int, isLocked bool) error
	DeleteForumCategory(ctx context.Context, id string) error

	// Topics
	CreateForumTopic(ctx context.Context, id, siteID, categoryID, authorID, slug, title, content, status string) (*ForumTopic, error)
	GetForumTopicByID(ctx context.Context, id string) (*ForumTopic, error)
	GetForumTopicBySlug(ctx context.Context, siteID, slug string) (*ForumTopic, error)
	ListForumTopics(ctx context.Context, siteID, categoryID, status, search string, limit, offset int) ([]ForumTopic, int64, error)
	ListForumTopicsByTag(ctx context.Context, siteID, tagID string, limit, offset int) ([]ForumTopic, int64, error)
	UpdateForumTopic(ctx context.Context, id, title, content string) error
	UpdateForumTopicStatus(ctx context.Context, id, status string) error
	PinForumTopic(ctx context.Context, id string, pinned bool) error
	IncrementForumTopicViews(ctx context.Context, id string) error
	DeleteForumTopic(ctx context.Context, id string) error

	// Posts
	CreateForumPost(ctx context.Context, id, siteID, topicID, parentID, authorID, content string) (*ForumPost, error)
	GetForumPostByID(ctx context.Context, id string) (*ForumPost, error)
	ListForumPosts(ctx context.Context, topicID string, limit, offset int) ([]ForumPost, int64, error)
	UpdateForumPost(ctx context.Context, id, content, editorID string) error
	MarkPostAsSolution(ctx context.Context, postID, topicID string) error
	DeleteForumPost(ctx context.Context, id string) error

	// Tags
	CreateForumTag(ctx context.Context, id, siteID, slug, name, description, color string) (*ForumTag, error)
	GetForumTagByID(ctx context.Context, id string) (*ForumTag, error)
	GetForumTagBySlug(ctx context.Context, siteID, slug string) (*ForumTag, error)
	ListForumTags(ctx context.Context, siteID string) ([]ForumTag, error)
	UpdateForumTag(ctx context.Context, id, slug, name, description, color string) error
	DeleteForumTag(ctx context.Context, id string) error
	AddTagToTopic(ctx context.Context, topicID, tagID string) error
	RemoveTagFromTopic(ctx context.Context, topicID, tagID string) error
	GetTopicTags(ctx context.Context, topicID string) ([]ForumTag, error)

	// Likes
	LikeForumTopic(ctx context.Context, id, siteID, userID, topicID string) error
	UnlikeForumTopic(ctx context.Context, userID, topicID string) error
	LikeForumPost(ctx context.Context, id, siteID, userID, postID string) error
	UnlikeForumPost(ctx context.Context, userID, postID string) error
	HasUserLikedTopic(ctx context.Context, userID, topicID string) (bool, error)
	HasUserLikedPost(ctx context.Context, userID, postID string) (bool, error)

	// Bookmarks
	CreateForumBookmark(ctx context.Context, id, siteID, userID, topicID string) error
	DeleteForumBookmark(ctx context.Context, userID, topicID string) error
	HasUserBookmarkedTopic(ctx context.Context, userID, topicID string) (bool, error)
	ListUserBookmarks(ctx context.Context, userID string, limit, offset int) ([]ForumBookmark, error)

	// Subscriptions
	CreateForumSubscription(ctx context.Context, id, siteID, userID, topicID, categoryID, level string) error
	UpdateForumSubscription(ctx context.Context, userID, topicID, categoryID, level string) error
	DeleteForumSubscription(ctx context.Context, userID, topicID, categoryID string) error
	GetTopicSubscribers(ctx context.Context, topicID string) ([]string, error)

	// Notifications
	CreateForumNotification(ctx context.Context, id, siteID, userID, notifType, title, message, topicID, postID, actorID string) error
	ListForumNotifications(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]ForumNotification, error)
	MarkNotificationRead(ctx context.Context, id string, userID string) error
	MarkAllNotificationsRead(ctx context.Context, userID string) error
	GetUnreadNotificationCount(ctx context.Context, userID string) (int64, error)

	// Flags
	CreateForumFlag(ctx context.Context, id, siteID, reporterID, topicID, postID, reason, description string) error
	GetForumFlagByID(ctx context.Context, id string) (*ForumFlag, error)
	ListForumFlags(ctx context.Context, siteID, status string, limit, offset int) ([]ForumFlag, int64, error)
	ResolveForumFlag(ctx context.Context, id, status, resolverID, resolutionNote string) error

	// Bans
	CreateForumBan(ctx context.Context, id, siteID, userID, bannedBy, reason, expiresAt string, isPermanent bool) error
	GetForumBan(ctx context.Context, siteID, userID string) (*ForumBan, error)
	IsUserBanned(ctx context.Context, siteID, userID string) (bool, error)
	ListForumBans(ctx context.Context, siteID string, limit, offset int) ([]ForumBan, int64, error)
	DeleteForumBan(ctx context.Context, siteID, userID string) error

	// User stats
	GetOrCreateForumUserStats(ctx context.Context, id, siteID, userID string) (*ForumUserStats, error)
	UpdateForumUserLastSeen(ctx context.Context, siteID, userID string) error
	IncrementForumUserStats(ctx context.Context, siteID, userID, field string, delta int) error
	GetForumLeaderboard(ctx context.Context, siteID string, limit int) ([]ForumUserStats, error)

	// Badges
	CreateForumBadge(ctx context.Context, id, siteID, slug, name, description, icon, color, criteria, tier string, isManual bool) (*ForumBadge, error)
	GetForumBadgeByID(ctx context.Context, id string) (*ForumBadge, error)
	ListForumBadges(ctx context.Context, siteID string) ([]ForumBadge, error)
	UpdateForumBadge(ctx context.Context, id, slug, name, description, icon, color, criteria, tier string, isManual bool) error
	DeleteForumBadge(ctx context.Context, id string) error
	AwardBadge(ctx context.Context, id, siteID, userID, badgeID, awardedBy string) error
	RevokeBadge(ctx context.Context, userID, badgeID string) error
	HasUserBadge(ctx context.Context, userID, badgeID string) (bool, error)
	GetUserBadges(ctx context.Context, userID string) ([]ForumUserBadge, error)

	// Reputation
	LogReputationChange(ctx context.Context, id, siteID, userID, action string, points int, topicID, postID, badgeID, note string) error
	GetReputationHistory(ctx context.Context, userID string, limit, offset int) ([]ForumReputationLog, error)

	// Stats and search
	GetForumStats(ctx context.Context, siteID string) (categories, topics, posts, users int64, err error)
	SearchForum(ctx context.Context, siteID, query string, limit, offset int) ([]ForumTopic, int64, error)
	CountUserTopicsToday(ctx context.Context, siteID, userID string) (int64, error)
	CountUserApprovedTopics(ctx context.Context, siteID, userID string) (int64, error)
	CountUserPostsToday(ctx context.Context, siteID, userID string) (int64, error)
}

// NewsletterStore handles newsletter subscriptions.
type NewsletterStore interface {
	CreateSubscriber(ctx context.Context, id, siteID, email, verifyTokenHash string) error
	CreateVerifiedSubscriber(ctx context.Context, id, siteID, email, provider, displayName string) error
	GetSubscriberByEmail(ctx context.Context, siteID, email string) (*Subscriber, error)
	GetSubscriberByToken(ctx context.Context, siteID, tokenHash string) (*Subscriber, error)
	VerifySubscriber(ctx context.Context, siteID, tokenHash string) error
	UpdateSubscriberToken(ctx context.Context, siteID, email, newTokenHash string) error
	UnsubscribeByEmail(ctx context.Context, siteID, email string) error
	ListSubscribers(ctx context.Context, siteID string, verifiedOnly bool) ([]Subscriber, error)
	CountSubscribers(ctx context.Context, siteID string, verifiedOnly bool) (int64, error)
	GetSubscriberStatsExtended(ctx context.Context, siteID string) (total, verified, pending, thisMonth int64, err error)
}

// SessionStore handles user session management.
type SessionStore interface {
	CreateSession(ctx context.Context, id, userID, tokenHash, ipAddress, userAgent string, expiresAt time.Time) error
	GetSessionByToken(ctx context.Context, tokenHash string) (*Session, error)
	DeleteSession(ctx context.Context, tokenHash string) error
	DeleteUserSessions(ctx context.Context, userID string) error
	CleanExpiredSessions(ctx context.Context) (int64, error)
	RevokeToken(ctx context.Context, jti string, expiresAt time.Time) error
	IsTokenRevoked(ctx context.Context, jti string) (bool, error)
	CleanRevokedTokens(ctx context.Context) (int64, error)
}

// SiteStore handles multi-tenant site management.
type SiteStore interface {
	CreateSite(ctx context.Context, id, name, domain, apiKeyHash string) (*Site, error)
	GetSiteByID(ctx context.Context, id string) (*Site, error)
	GetSiteByAPIKey(ctx context.Context, apiKeyHash string) (*Site, error)
	ListSites(ctx context.Context) ([]Site, error)
	UpdateSite(ctx context.Context, id, name, domain string) error
	UpdateSiteAPIKey(ctx context.Context, siteID, newAPIKeyHash string) error
	DeleteSite(ctx context.Context, id string) error
}

// UploadStore handles file upload metadata.
type UploadStore interface {
	CreateUpload(ctx context.Context, id, siteID, userID, filename, mimeType string, sizeBytes int64, storagePath, url string) (*Upload, error)
	GetUpload(ctx context.Context, id string) (*Upload, error)
	GetUploadByPath(ctx context.Context, storagePath string) (*Upload, error)
	DeleteUpload(ctx context.Context, id string) error
	ListUploads(ctx context.Context, siteID, userID string, limit, offset int) ([]Upload, error)
}

// UserStore handles user account management.
type UserStore interface {
	CreateUser(ctx context.Context, id, siteID, email, passwordHash, role, name string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, siteID, email string) (*User, error)
	GetUserByOAuth(ctx context.Context, siteID, provider, providerID string) (*User, error)
	ListUsers(ctx context.Context, siteID string) ([]User, error)
	UpdateUserLastLogin(ctx context.Context, userID string) error
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id, email, role, name string) error
	UpdateUserPassword(ctx context.Context, id, passwordHash string) error
	CountUsers(ctx context.Context, siteID string) (int, error)
	CountUsersByRole(ctx context.Context, siteID, role string) (int, error)
	CreateUserWithVerification(ctx context.Context, id, siteID, email, passwordHash, role, name, verifyToken string) (*User, error)
	GetUserByVerifyToken(ctx context.Context, token string) (*User, error)
	VerifyUserEmail(ctx context.Context, userID string) error
	CreateUserWithOAuth(ctx context.Context, id, siteID, email, oauthProvider, oauthID, name, avatarURL, role string) (*User, error)
	LinkOAuthToUser(ctx context.Context, userID, oauthProvider, oauthID string) error
	SetPasswordResetToken(ctx context.Context, userID, tokenHash string) error
	GetUserByResetToken(ctx context.Context, tokenHash string) (*User, error)
	ClearPasswordResetToken(ctx context.Context, userID string) error
}
