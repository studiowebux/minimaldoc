package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// Blog post queries

func (db *DB) CreateBlogPost(ctx context.Context, id, siteID, authorID, slug, title, description, content, featuredImage, tags, category, visibility string) (*BlogPost, error) {
	if visibility == "" {
		visibility = "public"
	}
	query := `
		INSERT INTO blog_posts (id, site_id, author_id, slug, title, description, content, featured_image, tags, category, status, visibility)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'draft', $11)
	`
	_, err := db.ExecContext(ctx, query, id, siteID, nullString(authorID), slug, title,
		nullString(description), content, nullString(featuredImage), tags, nullString(category), visibility)
	if err != nil {
		return nil, err
	}
	return db.GetBlogPostByID(ctx, id)
}

func (db *DB) GetBlogPostByID(ctx context.Context, id string) (*BlogPost, error) {
	query := `
		SELECT bp.id, bp.site_id, bp.author_id, bp.slug, bp.title, bp.description, bp.content,
			bp.featured_image, bp.tags, bp.category, bp.status, COALESCE(bp.visibility, 'public'), bp.published_at, bp.scheduled_at, bp.created_at, bp.updated_at,
			u.name, u.email
		FROM blog_posts bp
		LEFT JOIN users u ON bp.author_id = u.id
		WHERE bp.id = $1
	`
	var p BlogPost
	err := db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.SiteID, &p.AuthorID, &p.Slug, &p.Title, &p.Description, &p.Content,
		&p.FeaturedImage, &p.Tags, &p.Category, &p.Status, &p.Visibility, &p.PublishedAt, &p.ScheduledAt, &p.CreatedAt, &p.UpdatedAt,
		&p.AuthorName, &p.AuthorEmail,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (db *DB) GetBlogPostBySlug(ctx context.Context, siteID, slug string) (*BlogPost, error) {
	query := `
		SELECT bp.id, bp.site_id, bp.author_id, bp.slug, bp.title, bp.description, bp.content,
			bp.featured_image, bp.tags, bp.category, bp.status, COALESCE(bp.visibility, 'public'), bp.published_at, bp.scheduled_at, bp.created_at, bp.updated_at,
			u.name, u.email
		FROM blog_posts bp
		LEFT JOIN users u ON bp.author_id = u.id
		WHERE bp.site_id = $1 AND bp.slug = $2
	`
	var p BlogPost
	err := db.QueryRowContext(ctx, query, siteID, slug).Scan(
		&p.ID, &p.SiteID, &p.AuthorID, &p.Slug, &p.Title, &p.Description, &p.Content,
		&p.FeaturedImage, &p.Tags, &p.Category, &p.Status, &p.Visibility, &p.PublishedAt, &p.ScheduledAt, &p.CreatedAt, &p.UpdatedAt,
		&p.AuthorName, &p.AuthorEmail,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (db *DB) UpdateBlogPost(ctx context.Context, id, slug, title, description, content, featuredImage, tags, category, visibility string) error {
	if visibility == "" {
		visibility = "public"
	}
	query := `
		UPDATE blog_posts
		SET slug = $1, title = $2, description = $3, content = $4, featured_image = $5,
			tags = $6, category = $7, visibility = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
	`
	_, err := db.ExecContext(ctx, query, slug, title, nullString(description), content,
		nullString(featuredImage), tags, nullString(category), visibility, id)
	return err
}

func (db *DB) DeleteBlogPost(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM blog_posts WHERE id = $1`, id)
	return err
}

func (db *DB) ListBlogPosts(ctx context.Context, siteID string, status string, limit, offset int) ([]BlogPost, error) {
	query := `
		SELECT bp.id, bp.site_id, bp.author_id, bp.slug, bp.title, bp.description, bp.content,
			bp.featured_image, bp.tags, bp.category, bp.status, COALESCE(bp.visibility, 'public'), bp.published_at, bp.scheduled_at, bp.created_at, bp.updated_at,
			u.name, u.email
		FROM blog_posts bp
		LEFT JOIN users u ON bp.author_id = u.id
		WHERE bp.site_id = $1
	`
	args := []interface{}{siteID}
	argIndex := 2

	if status != "" {
		query += fmt.Sprintf(` AND bp.status = $%d`, argIndex)
		args = append(args, status)
		argIndex++
	}

	query += ` ORDER BY bp.created_at DESC`
	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []BlogPost
	for rows.Next() {
		var p BlogPost
		if err := rows.Scan(
			&p.ID, &p.SiteID, &p.AuthorID, &p.Slug, &p.Title, &p.Description, &p.Content,
			&p.FeaturedImage, &p.Tags, &p.Category, &p.Status, &p.Visibility, &p.PublishedAt, &p.ScheduledAt, &p.CreatedAt, &p.UpdatedAt,
			&p.AuthorName, &p.AuthorEmail,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (db *DB) ListPublishedBlogPosts(ctx context.Context, siteID string, limit, offset int) ([]BlogPost, error) {
	posts, _, err := db.ListPublishedBlogPostsFiltered(ctx, siteID, "", "", "", limit, offset)
	return posts, err
}

// ListPublishedBlogPostsFiltered returns published public posts with filtering and pagination.
func (db *DB) ListPublishedBlogPostsFiltered(ctx context.Context, siteID, category, tag, search string, limit, offset int) ([]BlogPost, int64, error) {
	// Build dynamic query - only public visibility for unauthenticated access
	baseQuery := `
		FROM blog_posts bp
		LEFT JOIN users u ON bp.author_id = u.id
		WHERE bp.site_id = $1 AND bp.status = 'published' AND COALESCE(bp.visibility, 'public') = 'public'
	`
	args := []interface{}{siteID}
	argIndex := 2

	if category != "" {
		baseQuery += fmt.Sprintf(` AND bp.category = $%d`, argIndex)
		args = append(args, category)
		argIndex++
	}

	if tag != "" {
		// Search in JSON array (SQLite compatible) - escape LIKE special chars
		baseQuery += fmt.Sprintf(` AND bp.tags LIKE $%d ESCAPE '\'`, argIndex)
		args = append(args, "%\""+escapeLike(tag)+"\"%")
		argIndex++
	}

	if search != "" {
		// Escape LIKE special chars to prevent injection
		baseQuery += fmt.Sprintf(` AND (bp.title LIKE $%d ESCAPE '\' OR bp.content LIKE $%d ESCAPE '\')`, argIndex, argIndex+1)
		searchTerm := "%" + escapeLike(search) + "%"
		args = append(args, searchTerm, searchTerm)
		argIndex += 2
	}

	// Get total count
	countQuery := `SELECT COUNT(*)` + baseQuery
	var total int64
	if err := db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get posts
	selectQuery := `
		SELECT bp.id, bp.site_id, bp.author_id, bp.slug, bp.title, bp.description, bp.content,
			bp.featured_image, bp.tags, bp.category, bp.status, COALESCE(bp.visibility, 'public'), bp.published_at, bp.scheduled_at, bp.created_at, bp.updated_at,
			u.name, u.email
	` + baseQuery + fmt.Sprintf(` ORDER BY bp.published_at DESC LIMIT $%d OFFSET $%d`, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []BlogPost
	for rows.Next() {
		var p BlogPost
		if err := rows.Scan(
			&p.ID, &p.SiteID, &p.AuthorID, &p.Slug, &p.Title, &p.Description, &p.Content,
			&p.FeaturedImage, &p.Tags, &p.Category, &p.Status, &p.Visibility, &p.PublishedAt, &p.ScheduledAt, &p.CreatedAt, &p.UpdatedAt,
			&p.AuthorName, &p.AuthorEmail,
		); err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}
	return posts, total, rows.Err()
}

// GetRelatedPosts returns posts related by category or tags (public only).
func (db *DB) GetRelatedPosts(ctx context.Context, siteID, excludeID, category, tags string, limit int) ([]BlogPost, error) {
	// Find posts with same category or overlapping tags - only public visibility
	// Use parameterized LIKE with escape to prevent injection
	query := `
		SELECT bp.id, bp.site_id, bp.author_id, bp.slug, bp.title, bp.description, bp.content,
			bp.featured_image, bp.tags, bp.category, bp.status, COALESCE(bp.visibility, 'public'), bp.published_at, bp.scheduled_at, bp.created_at, bp.updated_at,
			u.name, u.email
		FROM blog_posts bp
		LEFT JOIN users u ON bp.author_id = u.id
		WHERE bp.site_id = $1 AND bp.status = 'published' AND COALESCE(bp.visibility, 'public') = 'public' AND bp.id != $2
		AND (bp.category = $3 OR bp.tags LIKE $4 ESCAPE '\')
		ORDER BY
			CASE WHEN bp.category = $3 THEN 0 ELSE 1 END,
			bp.published_at DESC
		LIMIT $5
	`

	// Extract first tag for matching if tags is a JSON array
	firstTag := ""
	if tags != "" && tags != "[]" {
		// Simple extraction: find first quoted string
		if start := strings.Index(tags, "\""); start >= 0 {
			if end := strings.Index(tags[start+1:], "\""); end >= 0 {
				firstTag = tags[start+1 : start+1+end]
			}
		}
	}

	// Escape LIKE special chars and wrap with wildcards
	tagPattern := "%" + escapeLike(firstTag) + "%"
	rows, err := db.QueryContext(ctx, query, siteID, excludeID, category, tagPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []BlogPost
	for rows.Next() {
		var p BlogPost
		if err := rows.Scan(
			&p.ID, &p.SiteID, &p.AuthorID, &p.Slug, &p.Title, &p.Description, &p.Content,
			&p.FeaturedImage, &p.Tags, &p.Category, &p.Status, &p.Visibility, &p.PublishedAt, &p.ScheduledAt, &p.CreatedAt, &p.UpdatedAt,
			&p.AuthorName, &p.AuthorEmail,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (db *DB) GetBlogPostStats(ctx context.Context, siteID string) (total, published, draft int64, err error) {
	query := `
		SELECT
			COUNT(*),
			COUNT(CASE WHEN status = 'published' THEN 1 END),
			COUNT(CASE WHEN status = 'draft' THEN 1 END)
		FROM blog_posts WHERE site_id = $1
	`
	err = db.QueryRowContext(ctx, query, siteID).Scan(&total, &published, &draft)
	return
}

func (db *DB) PublishBlogPost(ctx context.Context, id string) error {
	query := `UPDATE blog_posts SET status = 'published', published_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	return err
}

func (db *DB) UnpublishBlogPost(ctx context.Context, id string) error {
	query := `UPDATE blog_posts SET status = 'draft', updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	return err
}

// ScheduleBlogPost sets a future publish time for a post.
func (db *DB) ScheduleBlogPost(ctx context.Context, id, scheduledAt string) error {
	query := `UPDATE blog_posts SET scheduled_at = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := db.ExecContext(ctx, query, scheduledAt, id)
	return err
}

// ClearSchedule removes the scheduled time from a post.
func (db *DB) ClearSchedule(ctx context.Context, id string) error {
	query := `UPDATE blog_posts SET scheduled_at = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	return err
}

// GetScheduledPosts returns posts that are due to be published.
func (db *DB) GetScheduledPosts(ctx context.Context) ([]BlogPost, error) {
	query := `
		SELECT bp.id, bp.site_id, bp.author_id, bp.slug, bp.title, bp.description, bp.content,
			bp.featured_image, bp.tags, bp.category, bp.status, COALESCE(bp.visibility, 'public'), bp.published_at, bp.scheduled_at, bp.created_at, bp.updated_at,
			u.name, u.email
		FROM blog_posts bp
		LEFT JOIN users u ON bp.author_id = u.id
		WHERE bp.status = 'draft' AND bp.scheduled_at IS NOT NULL AND bp.scheduled_at <= CURRENT_TIMESTAMP
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []BlogPost
	for rows.Next() {
		var p BlogPost
		if err := rows.Scan(
			&p.ID, &p.SiteID, &p.AuthorID, &p.Slug, &p.Title, &p.Description, &p.Content,
			&p.FeaturedImage, &p.Tags, &p.Category, &p.Status, &p.Visibility, &p.PublishedAt, &p.ScheduledAt, &p.CreatedAt, &p.UpdatedAt,
			&p.AuthorName, &p.AuthorEmail,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// PublishScheduledPosts publishes all posts that are due.
func (db *DB) PublishScheduledPosts(ctx context.Context) (int64, error) {
	query := `
		UPDATE blog_posts
		SET status = 'published', published_at = CURRENT_TIMESTAMP, scheduled_at = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE status = 'draft' AND scheduled_at IS NOT NULL AND scheduled_at <= CURRENT_TIMESTAMP
	`
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Blog comment queries

func (db *DB) CreateBlogComment(ctx context.Context, id, siteID, postID, parentID, authorName, authorEmail, content, ipAddress, userAgent string) (*BlogComment, error) {
	query := `
		INSERT INTO blog_comments (id, site_id, post_id, parent_id, author_name, author_email, content, ip_address, user_agent, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'pending')
	`
	_, err := db.ExecContext(ctx, query, id, siteID, postID, nullString(parentID), authorName, authorEmail, content, nullString(ipAddress), nullString(userAgent))
	if err != nil {
		return nil, err
	}
	return db.GetBlogCommentByID(ctx, id)
}

func (db *DB) GetBlogCommentByID(ctx context.Context, id string) (*BlogComment, error) {
	query := `
		SELECT bc.id, bc.site_id, bc.post_id, bc.parent_id, bc.author_name, bc.author_email,
			bc.content, bc.status, bc.ip_address, bc.user_agent, bc.created_at, bc.moderated_at, bc.moderated_by,
			u.name
		FROM blog_comments bc
		LEFT JOIN users u ON bc.moderated_by = u.id
		WHERE bc.id = $1
	`
	var c BlogComment
	err := db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.SiteID, &c.PostID, &c.ParentID, &c.AuthorName, &c.AuthorEmail,
		&c.Content, &c.Status, &c.IPAddress, &c.UserAgent, &c.CreatedAt, &c.ModeratedAt, &c.ModeratedBy,
		&c.ModeratorName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func (db *DB) ListBlogComments(ctx context.Context, siteID string, status string, limit, offset int) ([]BlogComment, error) {
	query := `
		SELECT bc.id, bc.site_id, bc.post_id, bc.parent_id, bc.author_name, bc.author_email,
			bc.content, bc.status, bc.ip_address, bc.user_agent, bc.created_at, bc.moderated_at, bc.moderated_by,
			u.name
		FROM blog_comments bc
		LEFT JOIN users u ON bc.moderated_by = u.id
		WHERE bc.site_id = $1
	`
	args := []interface{}{siteID}
	argIndex := 2

	if status != "" {
		query += fmt.Sprintf(` AND bc.status = $%d`, argIndex)
		args = append(args, status)
		argIndex++
	}

	query += ` ORDER BY bc.created_at DESC`
	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []BlogComment
	for rows.Next() {
		var c BlogComment
		if err := rows.Scan(
			&c.ID, &c.SiteID, &c.PostID, &c.ParentID, &c.AuthorName, &c.AuthorEmail,
			&c.Content, &c.Status, &c.IPAddress, &c.UserAgent, &c.CreatedAt, &c.ModeratedAt, &c.ModeratedBy,
			&c.ModeratorName,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (db *DB) ListApprovedComments(ctx context.Context, postID string, limit, offset int) ([]BlogComment, error) {
	query := `
		SELECT bc.id, bc.site_id, bc.post_id, bc.parent_id, bc.author_name, bc.author_email,
			bc.content, bc.status, bc.ip_address, bc.user_agent, bc.created_at, bc.moderated_at, bc.moderated_by,
			u.name
		FROM blog_comments bc
		LEFT JOIN users u ON bc.moderated_by = u.id
		WHERE bc.post_id = $1 AND bc.status = 'approved'
		ORDER BY bc.created_at ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []BlogComment
	for rows.Next() {
		var c BlogComment
		if err := rows.Scan(
			&c.ID, &c.SiteID, &c.PostID, &c.ParentID, &c.AuthorName, &c.AuthorEmail,
			&c.Content, &c.Status, &c.IPAddress, &c.UserAgent, &c.CreatedAt, &c.ModeratedAt, &c.ModeratedBy,
			&c.ModeratorName,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (db *DB) ListPendingComments(ctx context.Context, siteID string, limit, offset int) ([]BlogComment, error) {
	return db.ListBlogComments(ctx, siteID, "pending", limit, offset)
}

func (db *DB) ModerateComment(ctx context.Context, id, status, moderatorID string) error {
	query := `UPDATE blog_comments SET status = $1, moderated_at = CURRENT_TIMESTAMP, moderated_by = $2 WHERE id = $3`
	_, err := db.ExecContext(ctx, query, status, moderatorID, id)
	return err
}

func (db *DB) DeleteBlogComment(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM blog_comments WHERE id = $1`, id)
	return err
}

func (db *DB) GetCommentStats(ctx context.Context, siteID string) (total, pending, approved int64, err error) {
	query := `
		SELECT
			COUNT(*),
			COUNT(CASE WHEN status = 'pending' THEN 1 END),
			COUNT(CASE WHEN status = 'approved' THEN 1 END)
		FROM blog_comments WHERE site_id = $1
	`
	err = db.QueryRowContext(ctx, query, siteID).Scan(&total, &pending, &approved)
	return
}
