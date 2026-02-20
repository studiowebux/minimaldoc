package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// bootstrapAndGetSiteID performs bootstrap and returns the site ID.
func bootstrapAndGetSiteID(t *testing.T, r *Router) string {
	t.Helper()

	body := `{"site_name":"Test Site","email":"admin@test.com","password":"securepassword123"}`
	w := performRequest(r.Engine, http.MethodPost, "/api/bootstrap", strings.NewReader(body))
	if w.Code != http.StatusOK {
		t.Fatalf("bootstrap failed: %s", w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp["site_id"].(string)
}

func TestSubscribe(t *testing.T) {
	r, _, mockEmail := setupTestRouter(t)
	siteID := bootstrapAndGetSiteID(t, r)

	body := `{"site_id":"` + siteID + `","email":"subscriber@test.com"}`
	w := performRequest(r.Engine, http.MethodPost, "/api/newsletter/subscribe", strings.NewReader(body))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "verification_sent" {
		t.Errorf("status = %q, want %q", resp["status"], "verification_sent")
	}

	// Verify email was sent
	msgs := mockEmail.Messages()
	if len(msgs) == 0 {
		t.Error("expected verification email to be sent")
	}
}

func TestSubscribeDuplicate(t *testing.T) {
	r, _, _ := setupTestRouter(t)
	siteID := bootstrapAndGetSiteID(t, r)

	body := `{"site_id":"` + siteID + `","email":"subscriber@test.com"}`

	// First subscribe
	performRequest(r.Engine, http.MethodPost, "/api/newsletter/subscribe", strings.NewReader(body))

	// Second subscribe — should resend verification
	w := performRequest(r.Engine, http.MethodPost, "/api/newsletter/subscribe", strings.NewReader(body))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "verification_resent" {
		t.Errorf("status = %q, want %q", resp["status"], "verification_resent")
	}
}

func TestSubscribeInvalidEmail(t *testing.T) {
	r, _, _ := setupTestRouter(t)
	siteID := bootstrapAndGetSiteID(t, r)

	body := `{"site_id":"` + siteID + `","email":"not-an-email"}`
	w := performRequest(r.Engine, http.MethodPost, "/api/newsletter/subscribe", strings.NewReader(body))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUnsubscribe(t *testing.T) {
	r, db, _ := setupTestRouter(t)
	siteID := bootstrapAndGetSiteID(t, r)

	// Subscribe and verify directly in DB
	body := `{"site_id":"` + siteID + `","email":"subscriber@test.com"}`
	performRequest(r.Engine, http.MethodPost, "/api/newsletter/subscribe", strings.NewReader(body))

	sub, _ := db.GetSubscriberByEmail(testContext(), siteID, "subscriber@test.com")
	if sub != nil && sub.VerifyToken.Valid {
		db.VerifySubscriber(testContext(), siteID, sub.VerifyToken.String)
	}

	// Unsubscribe
	w := performRequest(r.Engine, http.MethodPost, "/api/newsletter/unsubscribe", strings.NewReader(body))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "unsubscribed" {
		t.Errorf("status = %q, want %q", resp["status"], "unsubscribed")
	}
}
