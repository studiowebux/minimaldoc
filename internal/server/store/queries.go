package store

import (
	"database/sql"
	"strings"
)

// escapeLike escapes special characters in LIKE patterns to prevent SQL injection.
// The escape character is backslash (\).
// Use with: WHERE column LIKE $1 ESCAPE '\'
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`) // Escape backslash first
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

// Site represents a registered site.
type Site struct {
	ID        string
	Name      string
	Domain    sql.NullString
	APIKey    string
	Config    string
	CreatedAt string // SQLite returns string
	UpdatedAt string
}

// User represents a user account.
type User struct {
	ID            string
	SiteID        string
	Email         string
	PasswordHash  sql.NullString
	Role          string
	OAuthProvider sql.NullString
	OAuthID       sql.NullString
	Name          sql.NullString
	AvatarURL     sql.NullString
	EmailVerified bool
	VerifyToken   sql.NullString
	CreatedAt     string
	UpdatedAt     string
	LastLoginAt   sql.NullString
}

// Session represents a user session.
type Session struct {
	ID        string
	UserID    string
	TokenHash string
	IPAddress sql.NullString
	UserAgent sql.NullString
	ExpiresAt string
	CreatedAt string
}

// PageView represents an analytics page view.
type PageView struct {
	ID          int64
	SiteID      string
	Path        string
	Referrer    sql.NullString
	Country     sql.NullString
	DeviceType  sql.NullString
	Browser     sql.NullString
	OS          sql.NullString
	SessionHash sql.NullString
	CreatedAt   string
}

// Rating represents page feedback.
type Rating struct {
	ID          int64
	SiteID      string
	Path        string
	Rating      int
	Feedback    sql.NullString
	SessionHash sql.NullString
	CreatedAt   string
}

// Subscriber represents a newsletter subscriber.
type Subscriber struct {
	ID               string
	SiteID           string
	Email            string
	Verified         bool
	VerifyToken      sql.NullString
	VerifySentAt     sql.NullString
	SubscribedAt     string
	UnsubscribedAt   sql.NullString
	OAuthProvider    string
	OAuthDisplayName string
	VerifiedVia      string
}

// BlogPost represents a blog post.
type BlogPost struct {
	ID            string         `json:"id"`
	SiteID        string         `json:"site_id"`
	AuthorID      sql.NullString `json:"author_id"`
	Slug          string         `json:"slug"`
	Title         string         `json:"title"`
	Description   sql.NullString `json:"description"`
	Content       string         `json:"content"`
	FeaturedImage sql.NullString `json:"featured_image"`
	Tags          string         `json:"tags"`
	Category      sql.NullString `json:"category"`
	Status        string         `json:"status"`
	Visibility    string         `json:"visibility"`
	PublishedAt   sql.NullString `json:"published_at"`
	ScheduledAt   sql.NullString `json:"scheduled_at"`
	CreatedAt     string         `json:"created_at"`
	UpdatedAt     string         `json:"updated_at"`
	AuthorName    sql.NullString `json:"author_name"`
	AuthorEmail   sql.NullString `json:"author_email"`
}

// BlogComment represents a blog comment.
type BlogComment struct {
	ID           string
	SiteID       string
	PostID       string
	ParentID     sql.NullString
	AuthorName   string
	AuthorEmail  string
	Content      string
	Status       string
	IPAddress    sql.NullString
	UserAgent    sql.NullString
	CreatedAt    string
	ModeratedAt  sql.NullString
	ModeratedBy  sql.NullString
	// Joined fields
	ModeratorName sql.NullString
	Replies       []BlogComment
}

// DocAccess represents an access rule for protected documentation paths.
type DocAccess struct {
	ID           string
	SiteID       string
	PathPattern  string
	RequiredRole string
	Description  sql.NullString
	CreatedAt    string
	UpdatedAt    string
}

// Upload represents an uploaded file.
type Upload struct {
	ID          string
	SiteID      string
	UserID      string
	Filename    string
	MimeType    string
	SizeBytes   int64
	StoragePath string
	URL         string
	CreatedAt   string
}

// Event represents a custom analytics event.
type Event struct {
	ID          int64
	SiteID      string
	Name        string
	Category    sql.NullString
	Path        sql.NullString
	Value       sql.NullString
	SessionHash sql.NullString
	CreatedAt   string
}

// ForumCategory represents a forum category.
type ForumCategory struct {
	ID          string
	SiteID      string
	ParentID    sql.NullString
	Slug        string
	Name        string
	Description sql.NullString
	Color       sql.NullString
	Icon        sql.NullString
	Position    int
	IsLocked    bool
	CreatedAt   string
	UpdatedAt   string
	// Joined/computed fields
	TopicCount int64
	PostCount  int64
	Children   []ForumCategory
}

// ForumTopic represents a forum discussion thread.
type ForumTopic struct {
	ID             string
	SiteID         string
	CategoryID     sql.NullString
	AuthorID       sql.NullString
	Slug           string
	Title          string
	Content        string
	Status         string
	IsPinned       bool
	IsSolved       bool
	SolutionPostID sql.NullString
	ViewCount      int64
	LikeCount      int64
	PostCount      int64
	LastPostAt     sql.NullString
	LastPostBy     sql.NullString
	CreatedAt      string
	UpdatedAt      string
	// Joined fields
	AuthorName     sql.NullString
	AuthorEmail    sql.NullString
	AuthorAvatar   sql.NullString
	CategoryName   sql.NullString
	CategorySlug   sql.NullString
	LastPostAuthor sql.NullString
	Tags           []ForumTag
	IsLiked        bool
	IsBookmarked   bool
}

// ForumPost represents a forum reply.
type ForumPost struct {
	ID         string
	SiteID     string
	TopicID    string
	ParentID   sql.NullString
	AuthorID   sql.NullString
	Content    string
	LikeCount  int64
	IsSolution bool
	EditedAt   sql.NullString
	EditedBy   sql.NullString
	CreatedAt  string
	UpdatedAt  string
	// Joined fields
	AuthorName   sql.NullString
	AuthorEmail  sql.NullString
	AuthorAvatar sql.NullString
	EditorName   sql.NullString
	IsLiked      bool
	Replies      []ForumPost
}

// ForumTag represents a forum tag.
type ForumTag struct {
	ID          string
	SiteID      string
	Slug        string
	Name        string
	Description sql.NullString
	Color       sql.NullString
	UsageCount  int64
	CreatedAt   string
}

// ForumLike represents a like on a topic or post.
type ForumLike struct {
	ID        string
	SiteID    string
	UserID    string
	TopicID   sql.NullString
	PostID    sql.NullString
	CreatedAt string
}

// ForumBookmark represents a saved topic.
type ForumBookmark struct {
	ID        string
	SiteID    string
	UserID    string
	TopicID   string
	CreatedAt string
	// Joined fields
	TopicTitle sql.NullString
	TopicSlug  sql.NullString
}

// ForumSubscription represents a topic/category subscription.
type ForumSubscription struct {
	ID         string
	SiteID     string
	UserID     string
	TopicID    sql.NullString
	CategoryID sql.NullString
	Level      string
	CreatedAt  string
}

// ForumMention represents an @mention.
type ForumMention struct {
	ID          string
	SiteID      string
	UserID      string
	MentionedBy string
	TopicID     sql.NullString
	PostID      sql.NullString
	IsRead      bool
	CreatedAt   string
	// Joined fields
	MentionerName sql.NullString
	TopicTitle    sql.NullString
}

// ForumNotification represents an in-app notification.
type ForumNotification struct {
	ID        string
	SiteID    string
	UserID    string
	Type      string
	Title     string
	Message   sql.NullString
	TopicID   sql.NullString
	PostID    sql.NullString
	ActorID   sql.NullString
	IsRead    bool
	CreatedAt string
	// Joined fields
	ActorName   sql.NullString
	ActorAvatar sql.NullString
	TopicSlug   sql.NullString
}

// ForumFlag represents a content report.
type ForumFlag struct {
	ID             string
	SiteID         string
	ReporterID     sql.NullString
	TopicID        sql.NullString
	PostID         sql.NullString
	Reason         string
	Description    sql.NullString
	Status         string
	ResolvedBy     sql.NullString
	ResolvedAt     sql.NullString
	ResolutionNote sql.NullString
	CreatedAt      string
	// Joined fields
	ReporterName sql.NullString
	ResolverName sql.NullString
	TopicTitle   sql.NullString
	PostContent  sql.NullString
}

// ForumBan represents a user ban.
type ForumBan struct {
	ID          string
	SiteID      string
	UserID      string
	BannedBy    sql.NullString
	Reason      string
	ExpiresAt   sql.NullString
	IsPermanent bool
	CreatedAt   string
	// Joined fields
	UserName    sql.NullString
	UserEmail   sql.NullString
	BannerName  sql.NullString
}

// ForumUserStats represents user reputation and stats.
type ForumUserStats struct {
	ID                string
	SiteID            string
	UserID            string
	Reputation        int64
	TopicCount        int64
	PostCount         int64
	LikeReceivedCount int64
	LikeGivenCount    int64
	SolutionCount     int64
	LastSeenAt        sql.NullString
	CreatedAt         string
	UpdatedAt         string
	// Joined fields
	UserName   sql.NullString
	UserEmail  sql.NullString
	UserAvatar sql.NullString
	Badges     []ForumUserBadge
}

// ForumBadge represents a badge definition.
type ForumBadge struct {
	ID          string
	SiteID      string
	Slug        string
	Name        string
	Description sql.NullString
	Icon        sql.NullString
	Color       sql.NullString
	Criteria    sql.NullString
	Tier        string
	IsManual    bool
	CreatedAt   string
}

// ForumUserBadge represents a badge awarded to a user.
type ForumUserBadge struct {
	ID        string
	SiteID    string
	UserID    string
	BadgeID   string
	AwardedBy sql.NullString
	AwardedAt string
	// Joined fields
	BadgeName  sql.NullString
	BadgeIcon  sql.NullString
	BadgeColor sql.NullString
	BadgeTier  sql.NullString
}

// ForumReputationLog represents a reputation point entry.
type ForumReputationLog struct {
	ID        string
	SiteID    string
	UserID    string
	Action    string
	Points    int
	TopicID   sql.NullString
	PostID    sql.NullString
	BadgeID   sql.NullString
	Note      sql.NullString
	CreatedAt string
}

// AuditLog represents an audit log entry for admin actions.
type AuditLog struct {
	ID         string
	SiteID     string
	UserID     sql.NullString
	UserEmail  string
	Action     string
	EntityType string
	EntityID   sql.NullString
	EntityName sql.NullString
	Details    sql.NullString
	IPAddress  sql.NullString
	UserAgent  sql.NullString
	CreatedAt  string
}

// Helper for nullable strings
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
