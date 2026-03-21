package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestBootstrap(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	body := `{"site_name":"Test Site","email":"admin@test.com","password":"securepassword123"}`
	w := performRequest(r.Engine, http.MethodPost, "/api/bootstrap", strings.NewReader(body))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["email"] != "admin@test.com" {
		t.Errorf("email = %q, want %q", resp["email"], "admin@test.com")
	}
	if resp["site_id"] == nil || resp["site_id"] == "" {
		t.Error("site_id should not be empty")
	}
	if resp["api_key"] == nil || resp["api_key"] == "" {
		t.Error("api_key should not be empty")
	}
}

func TestBootstrapDuplicate(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	body := `{"site_name":"Test Site","email":"admin@test.com","password":"securepassword123"}`

	// First bootstrap succeeds
	w := performRequest(r.Engine, http.MethodPost, "/api/bootstrap", strings.NewReader(body))
	if w.Code != http.StatusOK {
		t.Fatalf("first bootstrap failed: status = %d", w.Code)
	}

	// Second bootstrap fails
	w = performRequest(r.Engine, http.MethodPost, "/api/bootstrap", strings.NewReader(body))
	if w.Code != http.StatusBadRequest {
		t.Errorf("second bootstrap: status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestBootstrapMissingPassword(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	body := `{"site_name":"Test Site","email":"admin@test.com"}`
	w := performRequest(r.Engine, http.MethodPost, "/api/bootstrap", strings.NewReader(body))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestLogin(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	// Bootstrap first
	bootstrap := `{"site_name":"Test Site","email":"admin@test.com","password":"securepassword123"}`
	bw := performRequest(r.Engine, http.MethodPost, "/api/bootstrap", strings.NewReader(bootstrap))
	if bw.Code != http.StatusOK {
		t.Fatalf("bootstrap failed: %s", bw.Body.String())
	}

	var bResp map[string]interface{}
	json.Unmarshal(bw.Body.Bytes(), &bResp)
	siteID := bResp["site_id"].(string)

	// Login
	login := `{"email":"admin@test.com","password":"securepassword123","site_id":"` + siteID + `"}`
	w := performRequest(r.Engine, http.MethodPost, "/api/auth/login", strings.NewReader(login))

	if w.Code != http.StatusOK {
		t.Fatalf("login: status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["access_token"] == nil || resp["access_token"] == "" {
		t.Error("access_token should not be empty")
	}
	if resp["refresh_token"] == nil || resp["refresh_token"] == "" {
		t.Error("refresh_token should not be empty")
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	// Bootstrap first
	bootstrap := `{"site_name":"Test Site","email":"admin@test.com","password":"securepassword123"}`
	bw := performRequest(r.Engine, http.MethodPost, "/api/bootstrap", strings.NewReader(bootstrap))
	if bw.Code != http.StatusOK {
		t.Fatalf("bootstrap failed: %s", bw.Body.String())
	}

	var bResp map[string]interface{}
	json.Unmarshal(bw.Body.Bytes(), &bResp)
	siteID := bResp["site_id"].(string)

	// Login with wrong password
	login := `{"email":"admin@test.com","password":"wrongpassword","site_id":"` + siteID + `"}`
	w := performRequest(r.Engine, http.MethodPost, "/api/auth/login", strings.NewReader(login))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestGetCurrentUser(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	// Bootstrap and login
	bootstrap := `{"site_name":"Test Site","email":"admin@test.com","password":"securepassword123"}`
	bw := performRequest(r.Engine, http.MethodPost, "/api/bootstrap", strings.NewReader(bootstrap))
	var bResp map[string]interface{}
	json.Unmarshal(bw.Body.Bytes(), &bResp)
	siteID := bResp["site_id"].(string)

	login := `{"email":"admin@test.com","password":"securepassword123","site_id":"` + siteID + `"}`
	lw := performRequest(r.Engine, http.MethodPost, "/api/auth/login", strings.NewReader(login))
	if lw.Code != http.StatusOK {
		t.Fatalf("login failed: %s", lw.Body.String())
	}
	var lResp map[string]interface{}
	json.Unmarshal(lw.Body.Bytes(), &lResp)
	token := lResp["access_token"].(string)

	// Get current user
	w := performRequestWithAuth(r.Engine, http.MethodGet, "/api/auth/me", token, nil)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["email"] != "admin@test.com" {
		t.Errorf("email = %q, want %q", resp["email"], "admin@test.com")
	}
}

func TestGetCurrentUserUnauthorized(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	w := performRequest(r.Engine, http.MethodGet, "/api/auth/me", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
