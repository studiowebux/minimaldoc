package store

import (
	"testing"
)

func TestCreateBlogPost(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)
	authorID := createTestUser(t, db, siteID)

	// Act
	post, err := db.CreateBlogPost(ctx, "post-1", siteID, authorID,
		"hello-world", "Hello World", "A greeting post",
		"# Hello\nWorld", "", `["go","test"]`, "tech", "public")

	// Assert
	if err != nil {
		t.Fatalf("CreateBlogPost returned error: %v", err)
	}
	if post == nil {
		t.Fatal("CreateBlogPost returned nil post")
	}
	if post.ID != "post-1" {
		t.Errorf("expected ID %q, got %q", "post-1", post.ID)
	}
	if post.SiteID != siteID {
		t.Errorf("expected SiteID %q, got %q", siteID, post.SiteID)
	}
	if post.Slug != "hello-world" {
		t.Errorf("expected Slug %q, got %q", "hello-world", post.Slug)
	}
	if post.Title != "Hello World" {
		t.Errorf("expected Title %q, got %q", "Hello World", post.Title)
	}
	if post.Content != "# Hello\nWorld" {
		t.Errorf("expected Content %q, got %q", "# Hello\nWorld", post.Content)
	}
	if post.Tags != `["go","test"]` {
		t.Errorf("expected Tags %q, got %q", `["go","test"]`, post.Tags)
	}
	if post.Status != "draft" {
		t.Errorf("expected Status %q, got %q", "draft", post.Status)
	}
	if post.Visibility != "public" {
		t.Errorf("expected Visibility %q, got %q", "public", post.Visibility)
	}
	if post.CreatedAt == "" {
		t.Error("expected CreatedAt to be set")
	}
}

func TestGetBlogPostBySlug(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)
	authorID := createTestUser(t, db, siteID)

	_, err := db.CreateBlogPost(ctx, "post-slug-1", siteID, authorID,
		"my-slug", "Slug Post", "desc", "content", "", "[]", "", "public")
	if err != nil {
		t.Fatalf("setup CreateBlogPost: %v", err)
	}

	// Act
	post, err := db.GetBlogPostBySlug(ctx, siteID, "my-slug")

	// Assert
	if err != nil {
		t.Fatalf("GetBlogPostBySlug returned error: %v", err)
	}
	if post == nil {
		t.Fatal("GetBlogPostBySlug returned nil for existing slug")
	}
	if post.ID != "post-slug-1" {
		t.Errorf("expected ID %q, got %q", "post-slug-1", post.ID)
	}
	if post.Slug != "my-slug" {
		t.Errorf("expected Slug %q, got %q", "my-slug", post.Slug)
	}
	if post.Title != "Slug Post" {
		t.Errorf("expected Title %q, got %q", "Slug Post", post.Title)
	}

	// Non-existent slug returns nil
	missing, err := db.GetBlogPostBySlug(ctx, siteID, "no-such-slug")
	if err != nil {
		t.Fatalf("GetBlogPostBySlug for missing slug returned error: %v", err)
	}
	if missing != nil {
		t.Error("expected nil for non-existent slug")
	}
}

func TestListBlogPosts(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)
	authorID := createTestUser(t, db, siteID)

	_, err := db.CreateBlogPost(ctx, "list-post-1", siteID, authorID,
		"post-one", "Post One", "", "content one", "", "[]", "", "public")
	if err != nil {
		t.Fatalf("setup post 1: %v", err)
	}

	_, err = db.CreateBlogPost(ctx, "list-post-2", siteID, authorID,
		"post-two", "Post Two", "", "content two", "", "[]", "", "public")
	if err != nil {
		t.Fatalf("setup post 2: %v", err)
	}

	// Publish one post
	if err := db.PublishBlogPost(ctx, "list-post-2"); err != nil {
		t.Fatalf("setup publish: %v", err)
	}

	// Act: list all
	all, err := db.ListBlogPosts(ctx, siteID, "", 10, 0)

	// Assert
	if err != nil {
		t.Fatalf("ListBlogPosts all returned error: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(all))
	}

	// Act: list only drafts
	drafts, err := db.ListBlogPosts(ctx, siteID, "draft", 10, 0)

	// Assert
	if err != nil {
		t.Fatalf("ListBlogPosts draft returned error: %v", err)
	}
	if len(drafts) != 1 {
		t.Fatalf("expected 1 draft, got %d", len(drafts))
	}
	if drafts[0].ID != "list-post-1" {
		t.Errorf("expected draft post ID %q, got %q", "list-post-1", drafts[0].ID)
	}

	// Act: list only published
	published, err := db.ListBlogPosts(ctx, siteID, "published", 10, 0)

	// Assert
	if err != nil {
		t.Fatalf("ListBlogPosts published returned error: %v", err)
	}
	if len(published) != 1 {
		t.Fatalf("expected 1 published, got %d", len(published))
	}
	if published[0].ID != "list-post-2" {
		t.Errorf("expected published post ID %q, got %q", "list-post-2", published[0].ID)
	}
}

func TestUpdateBlogPost(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)
	authorID := createTestUser(t, db, siteID)

	_, err := db.CreateBlogPost(ctx, "update-post-1", siteID, authorID,
		"original-slug", "Original Title", "original desc",
		"original content", "", "[]", "general", "public")
	if err != nil {
		t.Fatalf("setup CreateBlogPost: %v", err)
	}

	// Act
	err = db.UpdateBlogPost(ctx, "update-post-1",
		"new-slug", "New Title", "new desc",
		"new content", "image.png", `["updated"]`, "tech", "private")

	// Assert
	if err != nil {
		t.Fatalf("UpdateBlogPost returned error: %v", err)
	}

	post, err := db.GetBlogPostByID(ctx, "update-post-1")
	if err != nil {
		t.Fatalf("GetBlogPostByID after update: %v", err)
	}
	if post == nil {
		t.Fatal("post is nil after update")
	}
	if post.Slug != "new-slug" {
		t.Errorf("expected Slug %q, got %q", "new-slug", post.Slug)
	}
	if post.Title != "New Title" {
		t.Errorf("expected Title %q, got %q", "New Title", post.Title)
	}
	if post.Content != "new content" {
		t.Errorf("expected Content %q, got %q", "new content", post.Content)
	}
	if post.Tags != `["updated"]` {
		t.Errorf("expected Tags %q, got %q", `["updated"]`, post.Tags)
	}
	if post.Visibility != "private" {
		t.Errorf("expected Visibility %q, got %q", "private", post.Visibility)
	}
}

func TestDeleteBlogPost(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)
	authorID := createTestUser(t, db, siteID)

	_, err := db.CreateBlogPost(ctx, "delete-post-1", siteID, authorID,
		"delete-me", "Delete Me", "", "bye", "", "[]", "", "public")
	if err != nil {
		t.Fatalf("setup CreateBlogPost: %v", err)
	}

	// Act
	err = db.DeleteBlogPost(ctx, "delete-post-1")

	// Assert
	if err != nil {
		t.Fatalf("DeleteBlogPost returned error: %v", err)
	}

	post, err := db.GetBlogPostByID(ctx, "delete-post-1")
	if err != nil {
		t.Fatalf("GetBlogPostByID after delete: %v", err)
	}
	if post != nil {
		t.Error("expected nil after deletion, got a post")
	}
}

func TestPublishUnpublishPost(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)
	authorID := createTestUser(t, db, siteID)

	_, err := db.CreateBlogPost(ctx, "pub-post-1", siteID, authorID,
		"pub-slug", "Pub Post", "", "content", "", "[]", "", "public")
	if err != nil {
		t.Fatalf("setup CreateBlogPost: %v", err)
	}

	// Verify initial draft status
	post, _ := db.GetBlogPostByID(ctx, "pub-post-1")
	if post.Status != "draft" {
		t.Fatalf("expected initial status %q, got %q", "draft", post.Status)
	}

	// Act: publish
	err = db.PublishBlogPost(ctx, "pub-post-1")

	// Assert
	if err != nil {
		t.Fatalf("PublishBlogPost returned error: %v", err)
	}

	post, _ = db.GetBlogPostByID(ctx, "pub-post-1")
	if post.Status != "published" {
		t.Errorf("expected status %q after publish, got %q", "published", post.Status)
	}
	if !post.PublishedAt.Valid {
		t.Error("expected PublishedAt to be set after publish")
	}

	// Act: unpublish
	err = db.UnpublishBlogPost(ctx, "pub-post-1")

	// Assert
	if err != nil {
		t.Fatalf("UnpublishBlogPost returned error: %v", err)
	}

	post, _ = db.GetBlogPostByID(ctx, "pub-post-1")
	if post.Status != "draft" {
		t.Errorf("expected status %q after unpublish, got %q", "draft", post.Status)
	}
}

func TestBlogPostStats(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)
	authorID := createTestUser(t, db, siteID)

	// Create 3 drafts
	for i, id := range []string{"stats-post-1", "stats-post-2", "stats-post-3"} {
		slug := "stats-slug-" + string(rune('a'+i))
		_, err := db.CreateBlogPost(ctx, id, siteID, authorID,
			slug, "Stats Post", "", "content", "", "[]", "", "public")
		if err != nil {
			t.Fatalf("setup post %s: %v", id, err)
		}
	}

	// Publish 2 of them
	if err := db.PublishBlogPost(ctx, "stats-post-1"); err != nil {
		t.Fatalf("publish post 1: %v", err)
	}
	if err := db.PublishBlogPost(ctx, "stats-post-2"); err != nil {
		t.Fatalf("publish post 2: %v", err)
	}

	// Act
	total, published, draft, err := db.GetBlogPostStats(ctx, siteID)

	// Assert
	if err != nil {
		t.Fatalf("GetBlogPostStats returned error: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if published != 2 {
		t.Errorf("expected published 2, got %d", published)
	}
	if draft != 1 {
		t.Errorf("expected draft 1, got %d", draft)
	}
}

func TestCreateBlogComment(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)
	authorID := createTestUser(t, db, siteID)

	_, err := db.CreateBlogPost(ctx, "comment-post-1", siteID, authorID,
		"comment-slug", "Comment Post", "", "content", "", "[]", "", "public")
	if err != nil {
		t.Fatalf("setup post: %v", err)
	}

	// Act
	comment, err := db.CreateBlogComment(ctx,
		"comment-1", siteID, "comment-post-1", "",
		"Jane Doe", "jane@example.com", "Great post!",
		"192.168.1.1", "TestAgent/1.0")

	// Assert
	if err != nil {
		t.Fatalf("CreateBlogComment returned error: %v", err)
	}
	if comment == nil {
		t.Fatal("CreateBlogComment returned nil comment")
	}
	if comment.ID != "comment-1" {
		t.Errorf("expected ID %q, got %q", "comment-1", comment.ID)
	}
	if comment.PostID != "comment-post-1" {
		t.Errorf("expected PostID %q, got %q", "comment-post-1", comment.PostID)
	}
	if comment.AuthorName != "Jane Doe" {
		t.Errorf("expected AuthorName %q, got %q", "Jane Doe", comment.AuthorName)
	}
	if comment.AuthorEmail != "jane@example.com" {
		t.Errorf("expected AuthorEmail %q, got %q", "jane@example.com", comment.AuthorEmail)
	}
	if comment.Content != "Great post!" {
		t.Errorf("expected Content %q, got %q", "Great post!", comment.Content)
	}
	if comment.Status != "pending" {
		t.Errorf("expected Status %q, got %q", "pending", comment.Status)
	}
}

func TestListApprovedComments(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)
	authorID := createTestUser(t, db, siteID)

	_, err := db.CreateBlogPost(ctx, "approved-post-1", siteID, authorID,
		"approved-slug", "Approved Post", "", "content", "", "[]", "", "public")
	if err != nil {
		t.Fatalf("setup post: %v", err)
	}

	// Create 3 comments
	for _, id := range []string{"appr-c1", "appr-c2", "appr-c3"} {
		_, err := db.CreateBlogComment(ctx, id, siteID, "approved-post-1", "",
			"Commenter", "c@example.com", "Comment "+id, "", "")
		if err != nil {
			t.Fatalf("setup comment %s: %v", id, err)
		}
	}

	// Approve 2 of 3
	if err := db.ModerateComment(ctx, "appr-c1", "approved", authorID); err != nil {
		t.Fatalf("approve c1: %v", err)
	}
	if err := db.ModerateComment(ctx, "appr-c3", "approved", authorID); err != nil {
		t.Fatalf("approve c3: %v", err)
	}

	// Act
	approved, err := db.ListApprovedComments(ctx, "approved-post-1", 10, 0)

	// Assert
	if err != nil {
		t.Fatalf("ListApprovedComments returned error: %v", err)
	}
	if len(approved) != 2 {
		t.Fatalf("expected 2 approved comments, got %d", len(approved))
	}

	// Verify the pending comment is excluded
	for _, c := range approved {
		if c.ID == "appr-c2" {
			t.Error("pending comment appr-c2 should not appear in approved list")
		}
	}
}

func TestModerateComment(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)
	authorID := createTestUser(t, db, siteID)

	_, err := db.CreateBlogPost(ctx, "mod-post-1", siteID, authorID,
		"mod-slug", "Mod Post", "", "content", "", "[]", "", "public")
	if err != nil {
		t.Fatalf("setup post: %v", err)
	}

	_, err = db.CreateBlogComment(ctx, "mod-comment-1", siteID, "mod-post-1", "",
		"Moderatee", "m@example.com", "Needs review", "", "")
	if err != nil {
		t.Fatalf("setup comment: %v", err)
	}

	// Verify initial pending status
	comment, _ := db.GetBlogCommentByID(ctx, "mod-comment-1")
	if comment.Status != "pending" {
		t.Fatalf("expected initial status %q, got %q", "pending", comment.Status)
	}

	// Act
	err = db.ModerateComment(ctx, "mod-comment-1", "approved", authorID)

	// Assert
	if err != nil {
		t.Fatalf("ModerateComment returned error: %v", err)
	}

	comment, err = db.GetBlogCommentByID(ctx, "mod-comment-1")
	if err != nil {
		t.Fatalf("GetBlogCommentByID after moderate: %v", err)
	}
	if comment.Status != "approved" {
		t.Errorf("expected status %q, got %q", "approved", comment.Status)
	}
	if !comment.ModeratedAt.Valid {
		t.Error("expected ModeratedAt to be set after moderation")
	}
	if !comment.ModeratedBy.Valid || comment.ModeratedBy.String != authorID {
		t.Errorf("expected ModeratedBy %q, got %q", authorID, comment.ModeratedBy.String)
	}
}
