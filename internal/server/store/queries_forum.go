package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// Category queries

func (db *DB) CreateForumCategory(ctx context.Context, id, siteID, parentID, slug, name, description, color, icon string, position int) (*ForumCategory, error) {
	query := `
		INSERT INTO forum_categories (id, site_id, parent_id, slug, name, description, color, icon, position)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := db.ExecContext(ctx, query, id, siteID, nullString(parentID), slug, name,
		nullString(description), nullString(color), nullString(icon), position)
	if err != nil {
		return nil, err
	}
	return db.GetForumCategoryByID(ctx, id)
}

func (db *DB) GetForumCategoryByID(ctx context.Context, id string) (*ForumCategory, error) {
	query := `
		SELECT id, site_id, parent_id, slug, name, description, color, icon, position, is_locked, COALESCE(visibility, 'public'), created_at, updated_at
		FROM forum_categories WHERE id = $1
	`
	var c ForumCategory
	var isLocked int
	err := db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.SiteID, &c.ParentID, &c.Slug, &c.Name, &c.Description,
		&c.Color, &c.Icon, &c.Position, &isLocked, &c.Visibility, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	c.IsLocked = isLocked != 0
	return &c, err
}

func (db *DB) GetForumCategoryBySlug(ctx context.Context, siteID, slug string) (*ForumCategory, error) {
	query := `
		SELECT id, site_id, parent_id, slug, name, description, color, icon, position, is_locked, COALESCE(visibility, 'public'), created_at, updated_at
		FROM forum_categories WHERE site_id = $1 AND slug = $2
	`
	var c ForumCategory
	var isLocked int
	err := db.QueryRowContext(ctx, query, siteID, slug).Scan(
		&c.ID, &c.SiteID, &c.ParentID, &c.Slug, &c.Name, &c.Description,
		&c.Color, &c.Icon, &c.Position, &isLocked, &c.Visibility, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	c.IsLocked = isLocked != 0
	return &c, err
}

// ListForumCategories returns categories filtered by visibility.
// Pass nil for maxVisibility to return all (admin). Pass allowed visibilities for public/member access.
func (db *DB) ListForumCategories(ctx context.Context, siteID string, visibilities ...string) ([]ForumCategory, error) {
	query := `
		SELECT fc.id, fc.site_id, fc.parent_id, fc.slug, fc.name, fc.description, fc.color, fc.icon, fc.position, fc.is_locked, COALESCE(fc.visibility, 'public'), fc.created_at, fc.updated_at,
			COALESCE((SELECT COUNT(*) FROM forum_topics WHERE category_id = fc.id), 0) as topic_count
		FROM forum_categories fc
		WHERE fc.site_id = $1
	`
	args := []interface{}{siteID}
	if len(visibilities) > 0 {
		placeholders := make([]string, len(visibilities))
		for i, v := range visibilities {
			args = append(args, v)
			placeholders[i] = fmt.Sprintf("$%d", i+2)
		}
		query += ` AND COALESCE(fc.visibility, 'public') IN (` + strings.Join(placeholders, ",") + `)`
	}
	query += ` ORDER BY fc.position, fc.name`

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []ForumCategory
	for rows.Next() {
		var c ForumCategory
		var isLocked int
		if err := rows.Scan(
			&c.ID, &c.SiteID, &c.ParentID, &c.Slug, &c.Name, &c.Description,
			&c.Color, &c.Icon, &c.Position, &isLocked, &c.Visibility, &c.CreatedAt, &c.UpdatedAt, &c.TopicCount,
		); err != nil {
			return nil, err
		}
		c.IsLocked = isLocked != 0
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (db *DB) UpdateForumCategory(ctx context.Context, id, slug, name, description, color, icon string, position int, isLocked bool) error {
	locked := 0
	if isLocked {
		locked = 1
	}
	query := `
		UPDATE forum_categories
		SET slug = $1, name = $2, description = $3, color = $4, icon = $5, position = $6, is_locked = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
	`
	_, err := db.ExecContext(ctx, query, slug, name, nullString(description), nullString(color), nullString(icon), position, locked, id)
	return err
}

func (db *DB) DeleteForumCategory(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM forum_categories WHERE id = $1`, id)
	return err
}

// Topic queries

func (db *DB) CreateForumTopic(ctx context.Context, id, siteID, categoryID, authorID, slug, title, content, status string) (*ForumTopic, error) {
	if status == "" {
		status = "published"
	}
	query := `
		INSERT INTO forum_topics (id, site_id, category_id, author_id, slug, title, content, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := db.ExecContext(ctx, query, id, siteID, nullString(categoryID), nullString(authorID), slug, title, content, status)
	if err != nil {
		return nil, err
	}
	return db.GetForumTopicByID(ctx, id)
}

func (db *DB) GetForumTopicByID(ctx context.Context, id string) (*ForumTopic, error) {
	query := `
		SELECT ft.id, ft.site_id, ft.category_id, ft.author_id, ft.slug, ft.title, ft.content,
			ft.status, ft.is_pinned, ft.is_solved, ft.solution_post_id, ft.view_count, ft.like_count, ft.post_count,
			ft.last_post_at, ft.last_post_by, ft.created_at, ft.updated_at,
			u.name, u.email, u.avatar_url,
			fc.name, fc.slug
		FROM forum_topics ft
		LEFT JOIN users u ON ft.author_id = u.id
		LEFT JOIN forum_categories fc ON ft.category_id = fc.id
		WHERE ft.id = $1
	`
	var t ForumTopic
	var isPinned, isSolved int
	err := db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.SiteID, &t.CategoryID, &t.AuthorID, &t.Slug, &t.Title, &t.Content,
		&t.Status, &isPinned, &isSolved, &t.SolutionPostID, &t.ViewCount, &t.LikeCount, &t.PostCount,
		&t.LastPostAt, &t.LastPostBy, &t.CreatedAt, &t.UpdatedAt,
		&t.AuthorName, &t.AuthorEmail, &t.AuthorAvatar,
		&t.CategoryName, &t.CategorySlug,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	t.IsPinned = isPinned != 0
	t.IsSolved = isSolved != 0
	return &t, err
}

func (db *DB) GetForumTopicBySlug(ctx context.Context, siteID, slug string) (*ForumTopic, error) {
	query := `
		SELECT ft.id, ft.site_id, ft.category_id, ft.author_id, ft.slug, ft.title, ft.content,
			ft.status, ft.is_pinned, ft.is_solved, ft.solution_post_id, ft.view_count, ft.like_count, ft.post_count,
			ft.last_post_at, ft.last_post_by, ft.created_at, ft.updated_at,
			u.name, u.email, u.avatar_url,
			fc.name, fc.slug
		FROM forum_topics ft
		LEFT JOIN users u ON ft.author_id = u.id
		LEFT JOIN forum_categories fc ON ft.category_id = fc.id
		WHERE ft.site_id = $1 AND ft.slug = $2
	`
	var t ForumTopic
	var isPinned, isSolved int
	err := db.QueryRowContext(ctx, query, siteID, slug).Scan(
		&t.ID, &t.SiteID, &t.CategoryID, &t.AuthorID, &t.Slug, &t.Title, &t.Content,
		&t.Status, &isPinned, &isSolved, &t.SolutionPostID, &t.ViewCount, &t.LikeCount, &t.PostCount,
		&t.LastPostAt, &t.LastPostBy, &t.CreatedAt, &t.UpdatedAt,
		&t.AuthorName, &t.AuthorEmail, &t.AuthorAvatar,
		&t.CategoryName, &t.CategorySlug,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	t.IsPinned = isPinned != 0
	t.IsSolved = isSolved != 0
	return &t, err
}

func (db *DB) ListForumTopics(ctx context.Context, siteID, categoryID, status, search string, limit, offset int) ([]ForumTopic, int64, error) {
	baseQuery := `
		FROM forum_topics ft
		LEFT JOIN users u ON ft.author_id = u.id
		LEFT JOIN forum_categories fc ON ft.category_id = fc.id
		WHERE ft.site_id = $1
	`
	args := []any{siteID}
	argIndex := 2

	if categoryID != "" {
		baseQuery += fmt.Sprintf(` AND ft.category_id = $%d`, argIndex)
		args = append(args, categoryID)
		argIndex++
	}

	if status != "" {
		baseQuery += fmt.Sprintf(` AND ft.status = $%d`, argIndex)
		args = append(args, status)
		argIndex++
	}

	if search != "" {
		baseQuery += fmt.Sprintf(` AND (ft.title LIKE $%d ESCAPE '\' OR ft.content LIKE $%d ESCAPE '\')`, argIndex, argIndex+1)
		searchTerm := "%" + escapeLike(search) + "%"
		args = append(args, searchTerm, searchTerm)
		argIndex += 2
	}

	// Count query
	countQuery := `SELECT COUNT(*)` + baseQuery
	var total int64
	if err := db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Select query with pagination
	selectQuery := `
		SELECT ft.id, ft.site_id, ft.category_id, ft.author_id, ft.slug, ft.title, ft.content,
			ft.status, ft.is_pinned, ft.is_solved, ft.solution_post_id, ft.view_count, ft.like_count, ft.post_count,
			ft.last_post_at, ft.last_post_by, ft.created_at, ft.updated_at,
			u.name, u.email, u.avatar_url,
			fc.name, fc.slug
	` + baseQuery + fmt.Sprintf(` ORDER BY ft.is_pinned DESC, ft.last_post_at DESC NULLS LAST, ft.created_at DESC LIMIT $%d OFFSET $%d`, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var topics []ForumTopic
	for rows.Next() {
		var t ForumTopic
		var isPinned, isSolved int
		if err := rows.Scan(
			&t.ID, &t.SiteID, &t.CategoryID, &t.AuthorID, &t.Slug, &t.Title, &t.Content,
			&t.Status, &isPinned, &isSolved, &t.SolutionPostID, &t.ViewCount, &t.LikeCount, &t.PostCount,
			&t.LastPostAt, &t.LastPostBy, &t.CreatedAt, &t.UpdatedAt,
			&t.AuthorName, &t.AuthorEmail, &t.AuthorAvatar,
			&t.CategoryName, &t.CategorySlug,
		); err != nil {
			return nil, 0, err
		}
		t.IsPinned = isPinned != 0
		t.IsSolved = isSolved != 0
		topics = append(topics, t)
	}
	return topics, total, rows.Err()
}

func (db *DB) ListForumTopicsByTag(ctx context.Context, siteID, tagID string, limit, offset int) ([]ForumTopic, int64, error) {
	baseQuery := `
		FROM forum_topics ft
		LEFT JOIN users u ON ft.author_id = u.id
		LEFT JOIN forum_categories fc ON ft.category_id = fc.id
		INNER JOIN forum_topic_tags ftt ON ft.id = ftt.topic_id
		WHERE ft.site_id = $1 AND ftt.tag_id = $2
	`
	args := []any{siteID, tagID}

	// Count query
	countQuery := `SELECT COUNT(*)` + baseQuery
	var total int64
	if err := db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Select query
	selectQuery := `
		SELECT ft.id, ft.site_id, ft.category_id, ft.author_id, ft.slug, ft.title, ft.content,
			ft.status, ft.is_pinned, ft.is_solved, ft.solution_post_id, ft.view_count, ft.like_count, ft.post_count,
			ft.last_post_at, ft.last_post_by, ft.created_at, ft.updated_at,
			u.name, u.email, u.avatar_url,
			fc.name, fc.slug
	` + baseQuery + ` ORDER BY ft.is_pinned DESC, ft.last_post_at DESC NULLS LAST LIMIT $3 OFFSET $4`
	args = append(args, limit, offset)

	rows, err := db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var topics []ForumTopic
	for rows.Next() {
		var t ForumTopic
		var isPinned, isSolved int
		if err := rows.Scan(
			&t.ID, &t.SiteID, &t.CategoryID, &t.AuthorID, &t.Slug, &t.Title, &t.Content,
			&t.Status, &isPinned, &isSolved, &t.SolutionPostID, &t.ViewCount, &t.LikeCount, &t.PostCount,
			&t.LastPostAt, &t.LastPostBy, &t.CreatedAt, &t.UpdatedAt,
			&t.AuthorName, &t.AuthorEmail, &t.AuthorAvatar,
			&t.CategoryName, &t.CategorySlug,
		); err != nil {
			return nil, 0, err
		}
		t.IsPinned = isPinned != 0
		t.IsSolved = isSolved != 0
		topics = append(topics, t)
	}
	return topics, total, rows.Err()
}

func (db *DB) UpdateForumTopic(ctx context.Context, id, title, content string) error {
	query := `UPDATE forum_topics SET title = $1, content = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	_, err := db.ExecContext(ctx, query, title, content, id)
	return err
}

func (db *DB) UpdateForumTopicStatus(ctx context.Context, id, status string) error {
	query := `UPDATE forum_topics SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := db.ExecContext(ctx, query, status, id)
	return err
}

func (db *DB) PinForumTopic(ctx context.Context, id string, pinned bool) error {
	pin := 0
	if pinned {
		pin = 1
	}
	query := `UPDATE forum_topics SET is_pinned = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := db.ExecContext(ctx, query, pin, id)
	return err
}

func (db *DB) IncrementForumTopicViews(ctx context.Context, id string) error {
	query := `UPDATE forum_topics SET view_count = view_count + 1 WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	return err
}

func (db *DB) DeleteForumTopic(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM forum_topics WHERE id = $1`, id)
	return err
}

// Post queries

func (db *DB) CreateForumPost(ctx context.Context, id, siteID, topicID, parentID, authorID, content string) (*ForumPost, error) {
	query := `
		INSERT INTO forum_posts (id, site_id, topic_id, parent_id, author_id, content)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := db.ExecContext(ctx, query, id, siteID, topicID, nullString(parentID), nullString(authorID), content)
	if err != nil {
		return nil, err
	}

	// Update topic post count and last post info
	updateQuery := `
		UPDATE forum_topics
		SET post_count = post_count + 1, last_post_at = CURRENT_TIMESTAMP, last_post_by = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err = db.ExecContext(ctx, updateQuery, nullString(authorID), topicID)
	if err != nil {
		return nil, err
	}

	return db.GetForumPostByID(ctx, id)
}

func (db *DB) GetForumPostByID(ctx context.Context, id string) (*ForumPost, error) {
	query := `
		SELECT fp.id, fp.site_id, fp.topic_id, fp.parent_id, fp.author_id, fp.content,
			fp.like_count, fp.is_solution, fp.edited_at, fp.edited_by, fp.created_at, fp.updated_at,
			u.name, u.email, u.avatar_url,
			eu.name
		FROM forum_posts fp
		LEFT JOIN users u ON fp.author_id = u.id
		LEFT JOIN users eu ON fp.edited_by = eu.id
		WHERE fp.id = $1
	`
	var p ForumPost
	var isSolution int
	err := db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.SiteID, &p.TopicID, &p.ParentID, &p.AuthorID, &p.Content,
		&p.LikeCount, &isSolution, &p.EditedAt, &p.EditedBy, &p.CreatedAt, &p.UpdatedAt,
		&p.AuthorName, &p.AuthorEmail, &p.AuthorAvatar,
		&p.EditorName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	p.IsSolution = isSolution != 0
	return &p, err
}

func (db *DB) ListForumPosts(ctx context.Context, topicID string, limit, offset int) ([]ForumPost, int64, error) {
	// Count query
	countQuery := `SELECT COUNT(*) FROM forum_posts WHERE topic_id = $1`
	var total int64
	if err := db.QueryRowContext(ctx, countQuery, topicID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Select query
	query := `
		SELECT fp.id, fp.site_id, fp.topic_id, fp.parent_id, fp.author_id, fp.content,
			fp.like_count, fp.is_solution, fp.edited_at, fp.edited_by, fp.created_at, fp.updated_at,
			u.name, u.email, u.avatar_url,
			eu.name
		FROM forum_posts fp
		LEFT JOIN users u ON fp.author_id = u.id
		LEFT JOIN users eu ON fp.edited_by = eu.id
		WHERE fp.topic_id = $1
		ORDER BY fp.is_solution DESC, fp.created_at ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := db.QueryContext(ctx, query, topicID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []ForumPost
	for rows.Next() {
		var p ForumPost
		var isSolution int
		if err := rows.Scan(
			&p.ID, &p.SiteID, &p.TopicID, &p.ParentID, &p.AuthorID, &p.Content,
			&p.LikeCount, &isSolution, &p.EditedAt, &p.EditedBy, &p.CreatedAt, &p.UpdatedAt,
			&p.AuthorName, &p.AuthorEmail, &p.AuthorAvatar,
			&p.EditorName,
		); err != nil {
			return nil, 0, err
		}
		p.IsSolution = isSolution != 0
		posts = append(posts, p)
	}
	return posts, total, rows.Err()
}

func (db *DB) UpdateForumPost(ctx context.Context, id, content, editorID string) error {
	query := `UPDATE forum_posts SET content = $1, edited_at = CURRENT_TIMESTAMP, edited_by = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	_, err := db.ExecContext(ctx, query, content, nullString(editorID), id)
	return err
}

func (db *DB) MarkPostAsSolution(ctx context.Context, postID, topicID string) error {
	// Clear any existing solution
	clearQuery := `UPDATE forum_posts SET is_solution = 0 WHERE topic_id = $1`
	if _, err := db.ExecContext(ctx, clearQuery, topicID); err != nil {
		return err
	}

	// Mark new solution
	markQuery := `UPDATE forum_posts SET is_solution = 1 WHERE id = $1`
	if _, err := db.ExecContext(ctx, markQuery, postID); err != nil {
		return err
	}

	// Update topic
	updateQuery := `UPDATE forum_topics SET is_solved = 1, solution_post_id = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := db.ExecContext(ctx, updateQuery, postID, topicID)
	return err
}

func (db *DB) DeleteForumPost(ctx context.Context, id string) error {
	// Get topic ID first to update count
	var topicID string
	if err := db.QueryRowContext(ctx, `SELECT topic_id FROM forum_posts WHERE id = $1`, id).Scan(&topicID); err != nil {
		return err
	}

	// Delete post
	if _, err := db.ExecContext(ctx, `DELETE FROM forum_posts WHERE id = $1`, id); err != nil {
		return err
	}

	// Update topic post count
	updateQuery := `UPDATE forum_topics SET post_count = post_count - 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := db.ExecContext(ctx, updateQuery, topicID)
	return err
}

// Tag queries

func (db *DB) CreateForumTag(ctx context.Context, id, siteID, slug, name, description, color string) (*ForumTag, error) {
	query := `INSERT INTO forum_tags (id, site_id, slug, name, description, color) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.ExecContext(ctx, query, id, siteID, slug, name, nullString(description), nullString(color))
	if err != nil {
		return nil, err
	}
	return db.GetForumTagByID(ctx, id)
}

func (db *DB) GetForumTagByID(ctx context.Context, id string) (*ForumTag, error) {
	query := `SELECT id, site_id, slug, name, description, color, usage_count, created_at FROM forum_tags WHERE id = $1`
	var t ForumTag
	err := db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.SiteID, &t.Slug, &t.Name, &t.Description, &t.Color, &t.UsageCount, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

func (db *DB) GetForumTagBySlug(ctx context.Context, siteID, slug string) (*ForumTag, error) {
	query := `SELECT id, site_id, slug, name, description, color, usage_count, created_at FROM forum_tags WHERE site_id = $1 AND slug = $2`
	var t ForumTag
	err := db.QueryRowContext(ctx, query, siteID, slug).Scan(&t.ID, &t.SiteID, &t.Slug, &t.Name, &t.Description, &t.Color, &t.UsageCount, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

func (db *DB) ListForumTags(ctx context.Context, siteID string) ([]ForumTag, error) {
	query := `SELECT id, site_id, slug, name, description, color, usage_count, created_at FROM forum_tags WHERE site_id = $1 ORDER BY usage_count DESC, name`
	rows, err := db.QueryContext(ctx, query, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []ForumTag
	for rows.Next() {
		var t ForumTag
		if err := rows.Scan(&t.ID, &t.SiteID, &t.Slug, &t.Name, &t.Description, &t.Color, &t.UsageCount, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func (db *DB) UpdateForumTag(ctx context.Context, id, slug, name, description, color string) error {
	query := `UPDATE forum_tags SET slug = $1, name = $2, description = $3, color = $4 WHERE id = $5`
	_, err := db.ExecContext(ctx, query, slug, name, nullString(description), nullString(color), id)
	return err
}

func (db *DB) DeleteForumTag(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM forum_tags WHERE id = $1`, id)
	return err
}

func (db *DB) AddTagToTopic(ctx context.Context, topicID, tagID string) error {
	query := `INSERT INTO forum_topic_tags (topic_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.ExecContext(ctx, query, topicID, tagID)
	if err != nil {
		return err
	}
	// Update usage count
	updateQuery := `UPDATE forum_tags SET usage_count = usage_count + 1 WHERE id = $1`
	_, err = db.ExecContext(ctx, updateQuery, tagID)
	return err
}

func (db *DB) RemoveTagFromTopic(ctx context.Context, topicID, tagID string) error {
	query := `DELETE FROM forum_topic_tags WHERE topic_id = $1 AND tag_id = $2`
	result, err := db.ExecContext(ctx, query, topicID, tagID)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected > 0 {
		updateQuery := `UPDATE forum_tags SET usage_count = usage_count - 1 WHERE id = $1 AND usage_count > 0`
		_, err = db.ExecContext(ctx, updateQuery, tagID)
	}
	return err
}

func (db *DB) GetTopicTags(ctx context.Context, topicID string) ([]ForumTag, error) {
	query := `
		SELECT ft.id, ft.site_id, ft.slug, ft.name, ft.description, ft.color, ft.usage_count, ft.created_at
		FROM forum_tags ft
		INNER JOIN forum_topic_tags ftt ON ft.id = ftt.tag_id
		WHERE ftt.topic_id = $1
		ORDER BY ft.name
	`
	rows, err := db.QueryContext(ctx, query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []ForumTag
	for rows.Next() {
		var t ForumTag
		if err := rows.Scan(&t.ID, &t.SiteID, &t.Slug, &t.Name, &t.Description, &t.Color, &t.UsageCount, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

// Like queries

func (db *DB) LikeForumTopic(ctx context.Context, id, siteID, userID, topicID string) error {
	query := `INSERT INTO forum_likes (id, site_id, user_id, topic_id) VALUES ($1, $2, $3, $4)`
	_, err := db.ExecContext(ctx, query, id, siteID, userID, topicID)
	if err != nil {
		return err
	}
	// Update like count
	updateQuery := `UPDATE forum_topics SET like_count = like_count + 1 WHERE id = $1`
	_, err = db.ExecContext(ctx, updateQuery, topicID)
	return err
}

func (db *DB) UnlikeForumTopic(ctx context.Context, userID, topicID string) error {
	query := `DELETE FROM forum_likes WHERE user_id = $1 AND topic_id = $2`
	result, err := db.ExecContext(ctx, query, userID, topicID)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected > 0 {
		updateQuery := `UPDATE forum_topics SET like_count = like_count - 1 WHERE id = $1 AND like_count > 0`
		_, err = db.ExecContext(ctx, updateQuery, topicID)
	}
	return err
}

func (db *DB) LikeForumPost(ctx context.Context, id, siteID, userID, postID string) error {
	query := `INSERT INTO forum_likes (id, site_id, user_id, post_id) VALUES ($1, $2, $3, $4)`
	_, err := db.ExecContext(ctx, query, id, siteID, userID, postID)
	if err != nil {
		return err
	}
	// Update like count
	updateQuery := `UPDATE forum_posts SET like_count = like_count + 1 WHERE id = $1`
	_, err = db.ExecContext(ctx, updateQuery, postID)
	return err
}

func (db *DB) UnlikeForumPost(ctx context.Context, userID, postID string) error {
	query := `DELETE FROM forum_likes WHERE user_id = $1 AND post_id = $2`
	result, err := db.ExecContext(ctx, query, userID, postID)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected > 0 {
		updateQuery := `UPDATE forum_posts SET like_count = like_count - 1 WHERE id = $1 AND like_count > 0`
		_, err = db.ExecContext(ctx, updateQuery, postID)
	}
	return err
}

func (db *DB) HasUserLikedTopic(ctx context.Context, userID, topicID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM forum_likes WHERE user_id = $1 AND topic_id = $2)`
	var exists bool
	err := db.QueryRowContext(ctx, query, userID, topicID).Scan(&exists)
	return exists, err
}

func (db *DB) HasUserLikedPost(ctx context.Context, userID, postID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM forum_likes WHERE user_id = $1 AND post_id = $2)`
	var exists bool
	err := db.QueryRowContext(ctx, query, userID, postID).Scan(&exists)
	return exists, err
}

// Bookmark queries

func (db *DB) CreateForumBookmark(ctx context.Context, id, siteID, userID, topicID string) error {
	query := `INSERT INTO forum_bookmarks (id, site_id, user_id, topic_id) VALUES ($1, $2, $3, $4)`
	_, err := db.ExecContext(ctx, query, id, siteID, userID, topicID)
	return err
}

func (db *DB) DeleteForumBookmark(ctx context.Context, userID, topicID string) error {
	query := `DELETE FROM forum_bookmarks WHERE user_id = $1 AND topic_id = $2`
	_, err := db.ExecContext(ctx, query, userID, topicID)
	return err
}

func (db *DB) HasUserBookmarkedTopic(ctx context.Context, userID, topicID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM forum_bookmarks WHERE user_id = $1 AND topic_id = $2)`
	var exists bool
	err := db.QueryRowContext(ctx, query, userID, topicID).Scan(&exists)
	return exists, err
}

func (db *DB) ListUserBookmarks(ctx context.Context, userID string, limit, offset int) ([]ForumBookmark, error) {
	query := `
		SELECT fb.id, fb.site_id, fb.user_id, fb.topic_id, fb.created_at, ft.title, ft.slug
		FROM forum_bookmarks fb
		INNER JOIN forum_topics ft ON fb.topic_id = ft.id
		WHERE fb.user_id = $1
		ORDER BY fb.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []ForumBookmark
	for rows.Next() {
		var b ForumBookmark
		if err := rows.Scan(&b.ID, &b.SiteID, &b.UserID, &b.TopicID, &b.CreatedAt, &b.TopicTitle, &b.TopicSlug); err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, rows.Err()
}

// Subscription queries

func (db *DB) CreateForumSubscription(ctx context.Context, id, siteID, userID, topicID, categoryID, level string) error {
	query := `INSERT INTO forum_subscriptions (id, site_id, user_id, topic_id, category_id, level) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.ExecContext(ctx, query, id, siteID, userID, nullString(topicID), nullString(categoryID), level)
	return err
}

func (db *DB) UpdateForumSubscription(ctx context.Context, userID, topicID, categoryID, level string) error {
	if topicID != "" {
		query := `UPDATE forum_subscriptions SET level = $1 WHERE user_id = $2 AND topic_id = $3`
		_, err := db.ExecContext(ctx, query, level, userID, topicID)
		return err
	}
	query := `UPDATE forum_subscriptions SET level = $1 WHERE user_id = $2 AND category_id = $3`
	_, err := db.ExecContext(ctx, query, level, userID, categoryID)
	return err
}

func (db *DB) DeleteForumSubscription(ctx context.Context, userID, topicID, categoryID string) error {
	if topicID != "" {
		query := `DELETE FROM forum_subscriptions WHERE user_id = $1 AND topic_id = $2`
		_, err := db.ExecContext(ctx, query, userID, topicID)
		return err
	}
	query := `DELETE FROM forum_subscriptions WHERE user_id = $1 AND category_id = $2`
	_, err := db.ExecContext(ctx, query, userID, categoryID)
	return err
}

func (db *DB) GetTopicSubscribers(ctx context.Context, topicID string) ([]string, error) {
	query := `SELECT user_id FROM forum_subscriptions WHERE topic_id = $1 AND level = 'watching'`
	rows, err := db.QueryContext(ctx, query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, rows.Err()
}

// Notification queries

func (db *DB) CreateForumNotification(ctx context.Context, id, siteID, userID, notifType, title, message, topicID, postID, actorID string) error {
	query := `
		INSERT INTO forum_notifications (id, site_id, user_id, type, title, message, topic_id, post_id, actor_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := db.ExecContext(ctx, query, id, siteID, userID, notifType, title, nullString(message), nullString(topicID), nullString(postID), nullString(actorID))
	return err
}

func (db *DB) ListForumNotifications(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]ForumNotification, error) {
	query := `
		SELECT fn.id, fn.site_id, fn.user_id, fn.type, fn.title, fn.message, fn.topic_id, fn.post_id, fn.actor_id, fn.is_read, fn.created_at,
			u.name, u.avatar_url, ft.slug
		FROM forum_notifications fn
		LEFT JOIN users u ON fn.actor_id = u.id
		LEFT JOIN forum_topics ft ON fn.topic_id = ft.id
		WHERE fn.user_id = $1
	`
	args := []any{userID}
	argIndex := 2

	if unreadOnly {
		query += fmt.Sprintf(` AND fn.is_read = $%d`, argIndex)
		args = append(args, 0)
		argIndex++
	}

	query += fmt.Sprintf(` ORDER BY fn.created_at DESC LIMIT $%d OFFSET $%d`, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []ForumNotification
	for rows.Next() {
		var n ForumNotification
		var isRead int
		if err := rows.Scan(
			&n.ID, &n.SiteID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.TopicID, &n.PostID, &n.ActorID, &isRead, &n.CreatedAt,
			&n.ActorName, &n.ActorAvatar, &n.TopicSlug,
		); err != nil {
			return nil, err
		}
		n.IsRead = isRead != 0
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

func (db *DB) MarkNotificationRead(ctx context.Context, id string, userID string) error {
	query := `UPDATE forum_notifications SET is_read = 1 WHERE id = $1 AND user_id = $2`
	_, err := db.ExecContext(ctx, query, id, userID)
	return err
}

func (db *DB) MarkAllNotificationsRead(ctx context.Context, userID string) error {
	query := `UPDATE forum_notifications SET is_read = 1 WHERE user_id = $1`
	_, err := db.ExecContext(ctx, query, userID)
	return err
}

func (db *DB) GetUnreadNotificationCount(ctx context.Context, userID string) (int64, error) {
	query := `SELECT COUNT(*) FROM forum_notifications WHERE user_id = $1 AND is_read = 0`
	var count int64
	err := db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

// Flag queries

func (db *DB) CreateForumFlag(ctx context.Context, id, siteID, reporterID, topicID, postID, reason, description string) error {
	query := `
		INSERT INTO forum_flags (id, site_id, reporter_id, topic_id, post_id, reason, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.ExecContext(ctx, query, id, siteID, nullString(reporterID), nullString(topicID), nullString(postID), reason, nullString(description))
	return err
}

func (db *DB) GetForumFlagByID(ctx context.Context, id string) (*ForumFlag, error) {
	query := `
		SELECT ff.id, ff.site_id, ff.reporter_id, ff.topic_id, ff.post_id, ff.reason, ff.description,
			ff.status, ff.resolved_by, ff.resolved_at, ff.resolution_note, ff.created_at,
			ur.name, ures.name, ft.title, fp.content
		FROM forum_flags ff
		LEFT JOIN users ur ON ff.reporter_id = ur.id
		LEFT JOIN users ures ON ff.resolved_by = ures.id
		LEFT JOIN forum_topics ft ON ff.topic_id = ft.id
		LEFT JOIN forum_posts fp ON ff.post_id = fp.id
		WHERE ff.id = $1
	`
	var f ForumFlag
	err := db.QueryRowContext(ctx, query, id).Scan(
		&f.ID, &f.SiteID, &f.ReporterID, &f.TopicID, &f.PostID, &f.Reason, &f.Description,
		&f.Status, &f.ResolvedBy, &f.ResolvedAt, &f.ResolutionNote, &f.CreatedAt,
		&f.ReporterName, &f.ResolverName, &f.TopicTitle, &f.PostContent,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &f, err
}

func (db *DB) ListForumFlags(ctx context.Context, siteID, status string, limit, offset int) ([]ForumFlag, int64, error) {
	baseQuery := `
		FROM forum_flags ff
		LEFT JOIN users ur ON ff.reporter_id = ur.id
		LEFT JOIN users ures ON ff.resolved_by = ures.id
		LEFT JOIN forum_topics ft ON ff.topic_id = ft.id
		LEFT JOIN forum_posts fp ON ff.post_id = fp.id
		WHERE ff.site_id = $1
	`
	args := []any{siteID}
	argIndex := 2

	if status != "" {
		baseQuery += fmt.Sprintf(` AND ff.status = $%d`, argIndex)
		args = append(args, status)
		argIndex++
	}

	// Count
	var total int64
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*)`+baseQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Select
	selectQuery := `
		SELECT ff.id, ff.site_id, ff.reporter_id, ff.topic_id, ff.post_id, ff.reason, ff.description,
			ff.status, ff.resolved_by, ff.resolved_at, ff.resolution_note, ff.created_at,
			ur.name, ures.name, ft.title, fp.content
	` + baseQuery + fmt.Sprintf(` ORDER BY ff.created_at DESC LIMIT $%d OFFSET $%d`, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var flags []ForumFlag
	for rows.Next() {
		var f ForumFlag
		if err := rows.Scan(
			&f.ID, &f.SiteID, &f.ReporterID, &f.TopicID, &f.PostID, &f.Reason, &f.Description,
			&f.Status, &f.ResolvedBy, &f.ResolvedAt, &f.ResolutionNote, &f.CreatedAt,
			&f.ReporterName, &f.ResolverName, &f.TopicTitle, &f.PostContent,
		); err != nil {
			return nil, 0, err
		}
		flags = append(flags, f)
	}
	return flags, total, rows.Err()
}

func (db *DB) ResolveForumFlag(ctx context.Context, id, status, resolverID, resolutionNote string) error {
	query := `
		UPDATE forum_flags
		SET status = $1, resolved_by = $2, resolved_at = CURRENT_TIMESTAMP, resolution_note = $3
		WHERE id = $4
	`
	_, err := db.ExecContext(ctx, query, status, resolverID, nullString(resolutionNote), id)
	return err
}

// Ban queries

func (db *DB) CreateForumBan(ctx context.Context, id, siteID, userID, bannedBy, reason, expiresAt string, isPermanent bool) error {
	perm := 0
	if isPermanent {
		perm = 1
	}
	query := `
		INSERT INTO forum_bans (id, site_id, user_id, banned_by, reason, expires_at, is_permanent)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.ExecContext(ctx, query, id, siteID, userID, nullString(bannedBy), reason, nullString(expiresAt), perm)
	return err
}

func (db *DB) GetForumBan(ctx context.Context, siteID, userID string) (*ForumBan, error) {
	query := `
		SELECT fb.id, fb.site_id, fb.user_id, fb.banned_by, fb.reason, fb.expires_at, fb.is_permanent, fb.created_at,
			u.name, u.email, ub.name
		FROM forum_bans fb
		LEFT JOIN users u ON fb.user_id = u.id
		LEFT JOIN users ub ON fb.banned_by = ub.id
		WHERE fb.site_id = $1 AND fb.user_id = $2
	`
	var b ForumBan
	var isPermanent int
	err := db.QueryRowContext(ctx, query, siteID, userID).Scan(
		&b.ID, &b.SiteID, &b.UserID, &b.BannedBy, &b.Reason, &b.ExpiresAt, &isPermanent, &b.CreatedAt,
		&b.UserName, &b.UserEmail, &b.BannerName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	b.IsPermanent = isPermanent != 0
	return &b, err
}

func (db *DB) IsUserBanned(ctx context.Context, siteID, userID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM forum_bans
			WHERE site_id = $1 AND user_id = $2
			AND (is_permanent = 1 OR expires_at > CURRENT_TIMESTAMP)
		)
	`
	var banned bool
	err := db.QueryRowContext(ctx, query, siteID, userID).Scan(&banned)
	return banned, err
}

func (db *DB) ListForumBans(ctx context.Context, siteID string, limit, offset int) ([]ForumBan, int64, error) {
	// Count
	var total int64
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM forum_bans WHERE site_id = $1`, siteID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Select
	query := `
		SELECT fb.id, fb.site_id, fb.user_id, fb.banned_by, fb.reason, fb.expires_at, fb.is_permanent, fb.created_at,
			u.name, u.email, ub.name
		FROM forum_bans fb
		LEFT JOIN users u ON fb.user_id = u.id
		LEFT JOIN users ub ON fb.banned_by = ub.id
		WHERE fb.site_id = $1
		ORDER BY fb.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := db.QueryContext(ctx, query, siteID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var bans []ForumBan
	for rows.Next() {
		var b ForumBan
		var isPermanent int
		if err := rows.Scan(
			&b.ID, &b.SiteID, &b.UserID, &b.BannedBy, &b.Reason, &b.ExpiresAt, &isPermanent, &b.CreatedAt,
			&b.UserName, &b.UserEmail, &b.BannerName,
		); err != nil {
			return nil, 0, err
		}
		b.IsPermanent = isPermanent != 0
		bans = append(bans, b)
	}
	return bans, total, rows.Err()
}

func (db *DB) DeleteForumBan(ctx context.Context, siteID, userID string) error {
	query := `DELETE FROM forum_bans WHERE site_id = $1 AND user_id = $2`
	_, err := db.ExecContext(ctx, query, siteID, userID)
	return err
}

// User stats queries

func (db *DB) GetOrCreateForumUserStats(ctx context.Context, id, siteID, userID string) (*ForumUserStats, error) {
	// Try to get existing
	query := `
		SELECT fus.id, fus.site_id, fus.user_id, fus.reputation, fus.topic_count, fus.post_count,
			fus.like_received_count, fus.like_given_count, fus.solution_count, fus.last_seen_at, fus.created_at, fus.updated_at,
			u.name, u.email, u.avatar_url
		FROM forum_user_stats fus
		LEFT JOIN users u ON fus.user_id = u.id
		WHERE fus.site_id = $1 AND fus.user_id = $2
	`
	var s ForumUserStats
	err := db.QueryRowContext(ctx, query, siteID, userID).Scan(
		&s.ID, &s.SiteID, &s.UserID, &s.Reputation, &s.TopicCount, &s.PostCount,
		&s.LikeReceivedCount, &s.LikeGivenCount, &s.SolutionCount, &s.LastSeenAt, &s.CreatedAt, &s.UpdatedAt,
		&s.UserName, &s.UserEmail, &s.UserAvatar,
	)
	if err == nil {
		return &s, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	// Create new
	insertQuery := `INSERT INTO forum_user_stats (id, site_id, user_id) VALUES ($1, $2, $3)`
	if _, err := db.ExecContext(ctx, insertQuery, id, siteID, userID); err != nil {
		return nil, err
	}

	// Return newly created
	return db.GetOrCreateForumUserStats(ctx, id, siteID, userID)
}

func (db *DB) UpdateForumUserLastSeen(ctx context.Context, siteID, userID string) error {
	query := `UPDATE forum_user_stats SET last_seen_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE site_id = $1 AND user_id = $2`
	_, err := db.ExecContext(ctx, query, siteID, userID)
	return err
}

func (db *DB) IncrementForumUserStats(ctx context.Context, siteID, userID, field string, delta int) error {
	validFields := map[string]bool{
		"topic_count":         true,
		"post_count":          true,
		"like_received_count": true,
		"like_given_count":    true,
		"solution_count":      true,
		"reputation":          true,
	}
	if !validFields[field] {
		return fmt.Errorf("invalid field: %s", field)
	}
	query := fmt.Sprintf(`UPDATE forum_user_stats SET %s = %s + $1, updated_at = CURRENT_TIMESTAMP WHERE site_id = $2 AND user_id = $3`, field, field)
	_, err := db.ExecContext(ctx, query, delta, siteID, userID)
	return err
}

func (db *DB) GetForumLeaderboard(ctx context.Context, siteID string, limit int) ([]ForumUserStats, error) {
	query := `
		SELECT fus.id, fus.site_id, fus.user_id, fus.reputation, fus.topic_count, fus.post_count,
			fus.like_received_count, fus.like_given_count, fus.solution_count, fus.last_seen_at, fus.created_at, fus.updated_at,
			u.name, u.email, u.avatar_url
		FROM forum_user_stats fus
		LEFT JOIN users u ON fus.user_id = u.id
		WHERE fus.site_id = $1
		ORDER BY fus.reputation DESC
		LIMIT $2
	`
	rows, err := db.QueryContext(ctx, query, siteID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []ForumUserStats
	for rows.Next() {
		var s ForumUserStats
		if err := rows.Scan(
			&s.ID, &s.SiteID, &s.UserID, &s.Reputation, &s.TopicCount, &s.PostCount,
			&s.LikeReceivedCount, &s.LikeGivenCount, &s.SolutionCount, &s.LastSeenAt, &s.CreatedAt, &s.UpdatedAt,
			&s.UserName, &s.UserEmail, &s.UserAvatar,
		); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}

// Badge queries

func (db *DB) CreateForumBadge(ctx context.Context, id, siteID, slug, name, description, icon, color, criteria, tier string, isManual bool) (*ForumBadge, error) {
	manual := 0
	if isManual {
		manual = 1
	}
	query := `
		INSERT INTO forum_badges (id, site_id, slug, name, description, icon, color, criteria, tier, is_manual)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := db.ExecContext(ctx, query, id, siteID, slug, name, nullString(description), nullString(icon), nullString(color), nullString(criteria), tier, manual)
	if err != nil {
		return nil, err
	}
	return db.GetForumBadgeByID(ctx, id)
}

func (db *DB) GetForumBadgeByID(ctx context.Context, id string) (*ForumBadge, error) {
	query := `SELECT id, site_id, slug, name, description, icon, color, criteria, tier, is_manual, created_at FROM forum_badges WHERE id = $1`
	var b ForumBadge
	var isManual int
	err := db.QueryRowContext(ctx, query, id).Scan(&b.ID, &b.SiteID, &b.Slug, &b.Name, &b.Description, &b.Icon, &b.Color, &b.Criteria, &b.Tier, &isManual, &b.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	b.IsManual = isManual != 0
	return &b, err
}

func (db *DB) ListForumBadges(ctx context.Context, siteID string) ([]ForumBadge, error) {
	query := `SELECT id, site_id, slug, name, description, icon, color, criteria, tier, is_manual, created_at FROM forum_badges WHERE site_id = $1 ORDER BY tier, name`
	rows, err := db.QueryContext(ctx, query, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []ForumBadge
	for rows.Next() {
		var b ForumBadge
		var isManual int
		if err := rows.Scan(&b.ID, &b.SiteID, &b.Slug, &b.Name, &b.Description, &b.Icon, &b.Color, &b.Criteria, &b.Tier, &isManual, &b.CreatedAt); err != nil {
			return nil, err
		}
		b.IsManual = isManual != 0
		badges = append(badges, b)
	}
	return badges, rows.Err()
}

func (db *DB) UpdateForumBadge(ctx context.Context, id, slug, name, description, icon, color, criteria, tier string, isManual bool) error {
	manual := 0
	if isManual {
		manual = 1
	}
	query := `UPDATE forum_badges SET slug = $1, name = $2, description = $3, icon = $4, color = $5, criteria = $6, tier = $7, is_manual = $8 WHERE id = $9`
	_, err := db.ExecContext(ctx, query, slug, name, nullString(description), nullString(icon), nullString(color), nullString(criteria), tier, manual, id)
	return err
}

func (db *DB) DeleteForumBadge(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM forum_badges WHERE id = $1`, id)
	return err
}

// User badge queries

func (db *DB) AwardBadge(ctx context.Context, id, siteID, userID, badgeID, awardedBy string) error {
	query := `INSERT INTO forum_user_badges (id, site_id, user_id, badge_id, awarded_by) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.ExecContext(ctx, query, id, siteID, userID, badgeID, nullString(awardedBy))
	return err
}

func (db *DB) RevokeBadge(ctx context.Context, userID, badgeID string) error {
	query := `DELETE FROM forum_user_badges WHERE user_id = $1 AND badge_id = $2`
	_, err := db.ExecContext(ctx, query, userID, badgeID)
	return err
}

func (db *DB) HasUserBadge(ctx context.Context, userID, badgeID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM forum_user_badges WHERE user_id = $1 AND badge_id = $2)`
	var exists bool
	err := db.QueryRowContext(ctx, query, userID, badgeID).Scan(&exists)
	return exists, err
}

func (db *DB) GetUserBadges(ctx context.Context, userID string) ([]ForumUserBadge, error) {
	query := `
		SELECT fub.id, fub.site_id, fub.user_id, fub.badge_id, fub.awarded_by, fub.awarded_at,
			fb.name, fb.icon, fb.color, fb.tier
		FROM forum_user_badges fub
		INNER JOIN forum_badges fb ON fub.badge_id = fb.id
		WHERE fub.user_id = $1
		ORDER BY fub.awarded_at DESC
	`
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []ForumUserBadge
	for rows.Next() {
		var b ForumUserBadge
		if err := rows.Scan(&b.ID, &b.SiteID, &b.UserID, &b.BadgeID, &b.AwardedBy, &b.AwardedAt, &b.BadgeName, &b.BadgeIcon, &b.BadgeColor, &b.BadgeTier); err != nil {
			return nil, err
		}
		badges = append(badges, b)
	}
	return badges, rows.Err()
}

// Reputation log queries

func (db *DB) LogReputationChange(ctx context.Context, id, siteID, userID, action string, points int, topicID, postID, badgeID, note string) error {
	query := `
		INSERT INTO forum_reputation_log (id, site_id, user_id, action, points, topic_id, post_id, badge_id, note)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := db.ExecContext(ctx, query, id, siteID, userID, action, points, nullString(topicID), nullString(postID), nullString(badgeID), nullString(note))
	return err
}

func (db *DB) GetReputationHistory(ctx context.Context, userID string, limit, offset int) ([]ForumReputationLog, error) {
	query := `
		SELECT id, site_id, user_id, action, points, topic_id, post_id, badge_id, note, created_at
		FROM forum_reputation_log
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []ForumReputationLog
	for rows.Next() {
		var l ForumReputationLog
		if err := rows.Scan(&l.ID, &l.SiteID, &l.UserID, &l.Action, &l.Points, &l.TopicID, &l.PostID, &l.BadgeID, &l.Note, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

// Forum stats

func (db *DB) GetForumStats(ctx context.Context, siteID string) (categories, topics, posts, users int64, err error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM forum_categories WHERE site_id = $1),
			(SELECT COUNT(*) FROM forum_topics WHERE site_id = $1),
			(SELECT COUNT(*) FROM forum_posts WHERE site_id = $1),
			(SELECT COUNT(*) FROM forum_user_stats WHERE site_id = $1)
	`
	err = db.QueryRowContext(ctx, query, siteID).Scan(&categories, &topics, &posts, &users)
	return
}

// Search

func (db *DB) SearchForum(ctx context.Context, siteID, query string, limit, offset int) ([]ForumTopic, int64, error) {
	return db.ListForumTopics(ctx, siteID, "", "", query, limit, offset)
}

// CountUserTopicsToday returns the number of topics created by a user today.
func (db *DB) CountUserTopicsToday(ctx context.Context, siteID, userID string) (int64, error) {
	query := `SELECT COUNT(*) FROM forum_topics WHERE site_id = $1 AND author_id = $2 AND created_at >= date('now')`
	var count int64
	err := db.QueryRowContext(ctx, query, siteID, userID).Scan(&count)
	return count, err
}

// CountUserApprovedTopics returns the number of published topics by a user.
func (db *DB) CountUserApprovedTopics(ctx context.Context, siteID, userID string) (int64, error) {
	query := `SELECT COUNT(*) FROM forum_topics WHERE site_id = $1 AND author_id = $2 AND status = 'published'`
	var count int64
	err := db.QueryRowContext(ctx, query, siteID, userID).Scan(&count)
	return count, err
}

// CountUserPostsToday returns the number of posts created by a user today.
func (db *DB) CountUserPostsToday(ctx context.Context, siteID, userID string) (int64, error) {
	query := `SELECT COUNT(*) FROM forum_posts WHERE site_id = $1 AND author_id = $2 AND created_at >= date('now')`
	var count int64
	err := db.QueryRowContext(ctx, query, siteID, userID).Scan(&count)
	return count, err
}
