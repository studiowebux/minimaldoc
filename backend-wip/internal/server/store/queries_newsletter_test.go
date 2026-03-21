package store

import (
	"testing"
)

func TestCreateSubscriber(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	err := db.CreateSubscriber(ctx, "sub-1", siteID, "alice@example.com", "token-abc")

	if err != nil {
		t.Fatalf("CreateSubscriber returned error: %v", err)
	}

	sub, err := db.GetSubscriberByEmail(ctx, siteID, "alice@example.com")
	if err != nil {
		t.Fatalf("GetSubscriberByEmail returned error: %v", err)
	}
	if sub == nil {
		t.Fatal("expected subscriber, got nil")
	}
	if sub.ID != "sub-1" {
		t.Errorf("expected ID %q, got %q", "sub-1", sub.ID)
	}
	if sub.SiteID != siteID {
		t.Errorf("expected SiteID %q, got %q", siteID, sub.SiteID)
	}
	if sub.Email != "alice@example.com" {
		t.Errorf("expected Email %q, got %q", "alice@example.com", sub.Email)
	}
	if sub.Verified {
		t.Error("expected Verified to be false for new subscriber")
	}
	if !sub.VerifyToken.Valid || sub.VerifyToken.String != "token-abc" {
		t.Errorf("expected VerifyToken %q, got %v", "token-abc", sub.VerifyToken)
	}
	if sub.SubscribedAt == "" {
		t.Error("expected SubscribedAt to be set")
	}
}

func TestGetSubscriberByEmail(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	err := db.CreateSubscriber(ctx, "sub-2", siteID, "bob@example.com", "token-bob")
	if err != nil {
		t.Fatalf("CreateSubscriber returned error: %v", err)
	}

	sub, err := db.GetSubscriberByEmail(ctx, siteID, "bob@example.com")

	if err != nil {
		t.Fatalf("GetSubscriberByEmail returned error: %v", err)
	}
	if sub == nil {
		t.Fatal("expected subscriber, got nil")
	}
	if sub.Email != "bob@example.com" {
		t.Errorf("expected Email %q, got %q", "bob@example.com", sub.Email)
	}

	// Not found returns nil, nil.
	sub, err = db.GetSubscriberByEmail(ctx, siteID, "nonexistent@example.com")
	if err != nil {
		t.Fatalf("expected nil error for missing subscriber, got: %v", err)
	}
	if sub != nil {
		t.Errorf("expected nil subscriber for missing email, got %+v", sub)
	}
}

func TestGetSubscriberByToken(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	err := db.CreateSubscriber(ctx, "sub-3", siteID, "carol@example.com", "token-carol")
	if err != nil {
		t.Fatalf("CreateSubscriber returned error: %v", err)
	}

	sub, err := db.GetSubscriberByToken(ctx, siteID, "token-carol")

	if err != nil {
		t.Fatalf("GetSubscriberByToken returned error: %v", err)
	}
	if sub == nil {
		t.Fatal("expected subscriber, got nil")
	}
	if sub.Email != "carol@example.com" {
		t.Errorf("expected Email %q, got %q", "carol@example.com", sub.Email)
	}

	// Not found returns nil, nil.
	sub, err = db.GetSubscriberByToken(ctx, siteID, "no-such-token")
	if err != nil {
		t.Fatalf("expected nil error for missing token, got: %v", err)
	}
	if sub != nil {
		t.Errorf("expected nil subscriber for missing token, got %+v", sub)
	}
}

func TestVerifySubscriber(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	err := db.CreateSubscriber(ctx, "sub-4", siteID, "dave@example.com", "token-dave")
	if err != nil {
		t.Fatalf("CreateSubscriber returned error: %v", err)
	}

	err = db.VerifySubscriber(ctx, siteID, "token-dave")

	if err != nil {
		t.Fatalf("VerifySubscriber returned error: %v", err)
	}

	sub, err := db.GetSubscriberByEmail(ctx, siteID, "dave@example.com")
	if err != nil {
		t.Fatalf("GetSubscriberByEmail returned error: %v", err)
	}
	if sub == nil {
		t.Fatal("expected subscriber, got nil")
	}
	if !sub.Verified {
		t.Error("expected Verified to be true after verification")
	}
	if sub.VerifyToken.Valid {
		t.Error("expected VerifyToken to be NULL after verification")
	}

	// Token lookup should return nil after verification (token cleared).
	sub, err = db.GetSubscriberByToken(ctx, siteID, "token-dave")
	if err != nil {
		t.Fatalf("GetSubscriberByToken returned error: %v", err)
	}
	if sub != nil {
		t.Error("expected nil subscriber after token was cleared by verification")
	}

	// Verifying a non-existent token returns sql.ErrNoRows.
	err = db.VerifySubscriber(ctx, siteID, "no-such-token")
	if err == nil {
		t.Error("expected error when verifying non-existent token, got nil")
	}
}

func TestUnsubscribeByEmail(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	err := db.CreateSubscriber(ctx, "sub-5", siteID, "eve@example.com", "token-eve")
	if err != nil {
		t.Fatalf("CreateSubscriber returned error: %v", err)
	}

	err = db.VerifySubscriber(ctx, siteID, "token-eve")
	if err != nil {
		t.Fatalf("VerifySubscriber returned error: %v", err)
	}

	err = db.UnsubscribeByEmail(ctx, siteID, "eve@example.com")

	if err != nil {
		t.Fatalf("UnsubscribeByEmail returned error: %v", err)
	}

	// Row still exists but has unsubscribed_at set.
	sub, err := db.GetSubscriberByEmail(ctx, siteID, "eve@example.com")
	if err != nil {
		t.Fatalf("GetSubscriberByEmail returned error: %v", err)
	}
	if sub == nil {
		t.Fatal("expected subscriber row to still exist after unsubscribe")
	}
	if !sub.UnsubscribedAt.Valid {
		t.Error("expected UnsubscribedAt to be set after unsubscribe")
	}

	// ListSubscribers excludes unsubscribed.
	subs, err := db.ListSubscribers(ctx, siteID, false)
	if err != nil {
		t.Fatalf("ListSubscribers returned error: %v", err)
	}
	if len(subs) != 0 {
		t.Errorf("expected 0 active subscribers after unsubscribe, got %d", len(subs))
	}
}

func TestListSubscribers(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	// Create three subscribers: two verified, one unverified.
	err := db.CreateSubscriber(ctx, "sub-6a", siteID, "frank@example.com", "token-frank")
	if err != nil {
		t.Fatalf("CreateSubscriber returned error: %v", err)
	}
	err = db.VerifySubscriber(ctx, siteID, "token-frank")
	if err != nil {
		t.Fatalf("VerifySubscriber returned error: %v", err)
	}

	err = db.CreateSubscriber(ctx, "sub-6b", siteID, "grace@example.com", "token-grace")
	if err != nil {
		t.Fatalf("CreateSubscriber returned error: %v", err)
	}
	err = db.VerifySubscriber(ctx, siteID, "token-grace")
	if err != nil {
		t.Fatalf("VerifySubscriber returned error: %v", err)
	}

	err = db.CreateSubscriber(ctx, "sub-6c", siteID, "heidi@example.com", "token-heidi")
	if err != nil {
		t.Fatalf("CreateSubscriber returned error: %v", err)
	}

	// All active (verified + unverified).
	all, err := db.ListSubscribers(ctx, siteID, false)
	if err != nil {
		t.Fatalf("ListSubscribers(all) returned error: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("expected 3 subscribers (all), got %d", len(all))
	}

	// Verified only.
	verified, err := db.ListSubscribers(ctx, siteID, true)
	if err != nil {
		t.Fatalf("ListSubscribers(verifiedOnly) returned error: %v", err)
	}
	if len(verified) != 2 {
		t.Errorf("expected 2 verified subscribers, got %d", len(verified))
	}
	for _, s := range verified {
		if !s.Verified {
			t.Errorf("expected only verified subscribers, got unverified: %s", s.Email)
		}
	}
}

func TestCountSubscribers(t *testing.T) {
	db := setupTestDB(t)
	ctx := testContext()
	siteID := createTestSite(t, db)

	// Create three subscribers: one verified, two unverified.
	err := db.CreateSubscriber(ctx, "sub-7a", siteID, "ivan@example.com", "token-ivan")
	if err != nil {
		t.Fatalf("CreateSubscriber returned error: %v", err)
	}
	err = db.VerifySubscriber(ctx, siteID, "token-ivan")
	if err != nil {
		t.Fatalf("VerifySubscriber returned error: %v", err)
	}

	err = db.CreateSubscriber(ctx, "sub-7b", siteID, "judy@example.com", "token-judy")
	if err != nil {
		t.Fatalf("CreateSubscriber returned error: %v", err)
	}

	err = db.CreateSubscriber(ctx, "sub-7c", siteID, "karl@example.com", "token-karl")
	if err != nil {
		t.Fatalf("CreateSubscriber returned error: %v", err)
	}

	// Count all active.
	countAll, err := db.CountSubscribers(ctx, siteID, false)
	if err != nil {
		t.Fatalf("CountSubscribers(all) returned error: %v", err)
	}
	if countAll != 3 {
		t.Errorf("expected count 3 (all), got %d", countAll)
	}

	// Count verified only.
	countVerified, err := db.CountSubscribers(ctx, siteID, true)
	if err != nil {
		t.Fatalf("CountSubscribers(verifiedOnly) returned error: %v", err)
	}
	if countVerified != 1 {
		t.Errorf("expected count 1 (verified only), got %d", countVerified)
	}

	// Unsubscribe one unverified; counts should decrease.
	err = db.UnsubscribeByEmail(ctx, siteID, "judy@example.com")
	if err != nil {
		t.Fatalf("UnsubscribeByEmail returned error: %v", err)
	}

	countAll, err = db.CountSubscribers(ctx, siteID, false)
	if err != nil {
		t.Fatalf("CountSubscribers(all) after unsubscribe returned error: %v", err)
	}
	if countAll != 2 {
		t.Errorf("expected count 2 after unsubscribe, got %d", countAll)
	}

	countVerified, err = db.CountSubscribers(ctx, siteID, true)
	if err != nil {
		t.Fatalf("CountSubscribers(verifiedOnly) after unsubscribe returned error: %v", err)
	}
	if countVerified != 1 {
		t.Errorf("expected verified count still 1 after unsubscribing unverified, got %d", countVerified)
	}
}
