package store

import (
	"testing"
)

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	user, err := db.CreateUser(ctx, "user-create-1", siteID, "alice@example.com", "$2a$12$fakehash", "admin", "Alice")
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	if user.ID != "user-create-1" {
		t.Errorf("expected ID %q, got %q", "user-create-1", user.ID)
	}
	if user.SiteID != siteID {
		t.Errorf("expected SiteID %q, got %q", siteID, user.SiteID)
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected Email %q, got %q", "alice@example.com", user.Email)
	}
	if !user.PasswordHash.Valid || user.PasswordHash.String != "$2a$12$fakehash" {
		t.Errorf("expected PasswordHash %q, got %v", "$2a$12$fakehash", user.PasswordHash)
	}
	if user.Role != "admin" {
		t.Errorf("expected Role %q, got %q", "admin", user.Role)
	}
	if !user.Name.Valid || user.Name.String != "Alice" {
		t.Errorf("expected Name %q, got %v", "Alice", user.Name)
	}
	if user.EmailVerified {
		t.Error("expected EmailVerified to be false for new user")
	}
	if user.CreatedAt == "" {
		t.Error("expected CreatedAt to be set")
	}
	if user.UpdatedAt == "" {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestGetUserByID(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	_, err := db.CreateUser(ctx, "user-get-1", siteID, "bob@example.com", "$2a$12$fakehash", "viewer", "Bob")
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	user, err := db.GetUserByID(ctx, "user-get-1")
	if err != nil {
		t.Fatalf("GetUserByID returned error: %v", err)
	}
	if user == nil {
		t.Fatal("GetUserByID returned nil for existing user")
	}

	if user.ID != "user-get-1" {
		t.Errorf("expected ID %q, got %q", "user-get-1", user.ID)
	}
	if user.Email != "bob@example.com" {
		t.Errorf("expected Email %q, got %q", "bob@example.com", user.Email)
	}
	if user.Role != "viewer" {
		t.Errorf("expected Role %q, got %q", "viewer", user.Role)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()

	user, err := db.GetUserByID(ctx, "nonexistent-id")
	if err != nil {
		t.Fatalf("GetUserByID returned error: %v", err)
	}

	if user != nil {
		t.Errorf("expected nil for nonexistent user, got %+v", user)
	}
}

func TestGetUserByEmail(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	_, err := db.CreateUser(ctx, "user-email-1", siteID, "carol@example.com", "$2a$12$fakehash", "editor", "Carol")
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	user, err := db.GetUserByEmail(ctx, siteID, "carol@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail returned error: %v", err)
	}
	if user == nil {
		t.Fatal("GetUserByEmail returned nil for existing user")
	}

	if user.ID != "user-email-1" {
		t.Errorf("expected ID %q, got %q", "user-email-1", user.ID)
	}
	if user.Email != "carol@example.com" {
		t.Errorf("expected Email %q, got %q", "carol@example.com", user.Email)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	user, err := db.GetUserByEmail(ctx, siteID, "nobody@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail returned error: %v", err)
	}

	if user != nil {
		t.Errorf("expected nil for nonexistent email, got %+v", user)
	}
}

func TestListUsers(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	_, err := db.CreateUser(ctx, "user-list-1", siteID, "dan@example.com", "$2a$12$fakehash", "admin", "Dan")
	if err != nil {
		t.Fatalf("CreateUser (1) returned error: %v", err)
	}

	_, err = db.CreateUser(ctx, "user-list-2", siteID, "eve@example.com", "$2a$12$fakehash", "viewer", "Eve")
	if err != nil {
		t.Fatalf("CreateUser (2) returned error: %v", err)
	}

	_, err = db.CreateUser(ctx, "user-list-3", siteID, "frank@example.com", "$2a$12$fakehash", "editor", "Frank")
	if err != nil {
		t.Fatalf("CreateUser (3) returned error: %v", err)
	}

	users, err := db.ListUsers(ctx, siteID)
	if err != nil {
		t.Fatalf("ListUsers returned error: %v", err)
	}

	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}
}

func TestListUsers_Empty(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	users, err := db.ListUsers(ctx, siteID)
	if err != nil {
		t.Fatalf("ListUsers returned error: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("expected 0 users for empty site, got %d", len(users))
	}
}

func TestUpdateUser(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	_, err := db.CreateUser(ctx, "user-update-1", siteID, "grace@example.com", "$2a$12$fakehash", "viewer", "Grace")
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	err = db.UpdateUser(ctx, "user-update-1", "grace.new@example.com", "admin", "Grace Updated")
	if err != nil {
		t.Fatalf("UpdateUser returned error: %v", err)
	}

	user, err := db.GetUserByID(ctx, "user-update-1")
	if err != nil {
		t.Fatalf("GetUserByID returned error: %v", err)
	}
	if user == nil {
		t.Fatal("GetUserByID returned nil after update")
	}

	if user.Email != "grace.new@example.com" {
		t.Errorf("expected Email %q, got %q", "grace.new@example.com", user.Email)
	}
	if user.Role != "admin" {
		t.Errorf("expected Role %q, got %q", "admin", user.Role)
	}
	if !user.Name.Valid || user.Name.String != "Grace Updated" {
		t.Errorf("expected Name %q, got %v", "Grace Updated", user.Name)
	}
}

func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	_, err := db.CreateUser(ctx, "user-delete-1", siteID, "heidi@example.com", "$2a$12$fakehash", "viewer", "Heidi")
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	err = db.DeleteUser(ctx, "user-delete-1")
	if err != nil {
		t.Fatalf("DeleteUser returned error: %v", err)
	}

	user, err := db.GetUserByID(ctx, "user-delete-1")
	if err != nil {
		t.Fatalf("GetUserByID returned error after delete: %v", err)
	}

	if user != nil {
		t.Errorf("expected nil after delete, got %+v", user)
	}
}

func TestUpdateUserLastLogin(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	_, err := db.CreateUser(ctx, "user-login-1", siteID, "ivan@example.com", "$2a$12$fakehash", "viewer", "Ivan")
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	err = db.UpdateUserLastLogin(ctx, "user-login-1")
	if err != nil {
		t.Fatalf("UpdateUserLastLogin returned error: %v", err)
	}

	user, err := db.GetUserByID(ctx, "user-login-1")
	if err != nil {
		t.Fatalf("GetUserByID returned error: %v", err)
	}
	if user == nil {
		t.Fatal("GetUserByID returned nil after last login update")
	}

	if !user.LastLoginAt.Valid {
		t.Error("expected LastLoginAt to be set after UpdateUserLastLogin")
	}
}

func TestCountUsers(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	count, err := db.CountUsers(ctx, siteID)
	if err != nil {
		t.Fatalf("CountUsers returned error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected count 0 for empty site, got %d", count)
	}

	_, err = db.CreateUser(ctx, "user-count-1", siteID, "judy@example.com", "$2a$12$fakehash", "admin", "Judy")
	if err != nil {
		t.Fatalf("CreateUser (1) returned error: %v", err)
	}

	_, err = db.CreateUser(ctx, "user-count-2", siteID, "karl@example.com", "$2a$12$fakehash", "viewer", "Karl")
	if err != nil {
		t.Fatalf("CreateUser (2) returned error: %v", err)
	}

	count, err = db.CountUsers(ctx, siteID)
	if err != nil {
		t.Fatalf("CountUsers returned error: %v", err)
	}

	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
}
