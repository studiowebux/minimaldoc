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
	ID             string
	SiteID         string
	Email          string
	Verified       bool
	VerifyToken    sql.NullString
	VerifySentAt   sql.NullString
	SubscribedAt   string
	UnsubscribedAt sql.NullString
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

// Helper for nullable strings
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
