package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestLiveness(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	w := performRequest(r.Engine, http.MethodGet, "/healthz", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["status"] != "alive" {
		t.Errorf("status = %q, want %q", resp["status"], "alive")
	}
}

func TestReadiness(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	w := performRequest(r.Engine, http.MethodGet, "/readyz", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["status"] != "ready" {
		t.Errorf("status = %q, want %q", resp["status"], "ready")
	}
}

func TestHealthCheck(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	w := performRequest(r.Engine, http.MethodGet, "/api/health", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["status"] != "ok" {
		t.Errorf("status = %q, want %q", resp["status"], "ok")
	}
	if resp["database"] != "sqlite" {
		t.Errorf("database = %q, want %q", resp["database"], "sqlite")
	}
}
