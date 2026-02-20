package store

import (
	"database/sql"
	"testing"
)

func TestCreateSite(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		siteName   string
		domain     string
		apiKeyHash string
	}{
		{
			name:       "creates site with all fields",
			id:         "site-001",
			siteName:   "My Docs",
			domain:     "docs.example.com",
			apiKeyHash: "hashed-key-abc",
		},
		{
			name:       "creates site with empty domain",
			id:         "site-002",
			siteName:   "No Domain Site",
			domain:     "",
			apiKeyHash: "hashed-key-def",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			ctx := testContext()

			site, err := db.CreateSite(ctx, tt.id, tt.siteName, tt.domain, tt.apiKeyHash)

			if err != nil {
				t.Fatalf("CreateSite returned error: %v", err)
			}
			if site == nil {
				t.Fatal("CreateSite returned nil site")
			}
			if site.ID != tt.id {
				t.Errorf("ID = %q, want %q", site.ID, tt.id)
			}
			if site.Name != tt.siteName {
				t.Errorf("Name = %q, want %q", site.Name, tt.siteName)
			}
			if tt.domain != "" {
				if !site.Domain.Valid || site.Domain.String != tt.domain {
					t.Errorf("Domain = %v, want Valid=%t String=%q", site.Domain, true, tt.domain)
				}
			} else {
				if site.Domain.Valid {
					t.Errorf("Domain.Valid = true, want false for empty domain")
				}
			}
			if site.APIKey != tt.apiKeyHash {
				t.Errorf("APIKey = %q, want %q", site.APIKey, tt.apiKeyHash)
			}
			if site.CreatedAt == "" {
				t.Error("CreatedAt is empty")
			}
			if site.UpdatedAt == "" {
				t.Error("UpdatedAt is empty")
			}
		})
	}
}

func TestGetSiteByID(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		siteName   string
		domain     string
		apiKeyHash string
		lookupID   string
		wantNil    bool
	}{
		{
			name:       "returns existing site",
			id:         "site-get-1",
			siteName:   "Get Test",
			domain:     "get.example.com",
			apiKeyHash: "hash-get-1",
			lookupID:   "site-get-1",
			wantNil:    false,
		},
		{
			name:       "returns nil for nonexistent site",
			id:         "site-get-2",
			siteName:   "Exists",
			domain:     "exists.example.com",
			apiKeyHash: "hash-get-2",
			lookupID:   "nonexistent-id",
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			ctx := testContext()

			created, err := db.CreateSite(ctx, tt.id, tt.siteName, tt.domain, tt.apiKeyHash)
			if err != nil {
				t.Fatalf("setup CreateSite: %v", err)
			}

			site, err := db.GetSiteByID(ctx, tt.lookupID)

			if err != nil {
				t.Fatalf("GetSiteByID returned error: %v", err)
			}
			if tt.wantNil {
				if site != nil {
					t.Errorf("expected nil, got site with ID %q", site.ID)
				}
				return
			}
			if site == nil {
				t.Fatal("GetSiteByID returned nil for existing site")
			}
			if site.ID != created.ID {
				t.Errorf("ID = %q, want %q", site.ID, created.ID)
			}
			if site.Name != created.Name {
				t.Errorf("Name = %q, want %q", site.Name, created.Name)
			}
			if site.Domain != created.Domain {
				t.Errorf("Domain = %v, want %v", site.Domain, created.Domain)
			}
			if site.APIKey != created.APIKey {
				t.Errorf("APIKey = %q, want %q", site.APIKey, created.APIKey)
			}
			if site.CreatedAt != created.CreatedAt {
				t.Errorf("CreatedAt = %q, want %q", site.CreatedAt, created.CreatedAt)
			}
		})
	}
}

func TestGetSiteByAPIKey(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		siteName      string
		domain        string
		apiKeyHash    string
		lookupHash    string
		wantNil       bool
	}{
		{
			name:       "returns site by matching API key hash",
			id:         "site-apikey-1",
			siteName:   "API Key Test",
			domain:     "apikey.example.com",
			apiKeyHash: "unique-hash-123",
			lookupHash: "unique-hash-123",
			wantNil:    false,
		},
		{
			name:       "returns nil for unknown API key hash",
			id:         "site-apikey-2",
			siteName:   "Another Site",
			domain:     "another.example.com",
			apiKeyHash: "known-hash-456",
			lookupHash: "unknown-hash-999",
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			ctx := testContext()

			_, err := db.CreateSite(ctx, tt.id, tt.siteName, tt.domain, tt.apiKeyHash)
			if err != nil {
				t.Fatalf("setup CreateSite: %v", err)
			}

			site, err := db.GetSiteByAPIKey(ctx, tt.lookupHash)

			if err != nil {
				t.Fatalf("GetSiteByAPIKey returned error: %v", err)
			}
			if tt.wantNil {
				if site != nil {
					t.Errorf("expected nil, got site with ID %q", site.ID)
				}
				return
			}
			if site == nil {
				t.Fatal("GetSiteByAPIKey returned nil for matching hash")
			}
			if site.ID != tt.id {
				t.Errorf("ID = %q, want %q", site.ID, tt.id)
			}
			if site.APIKey != tt.apiKeyHash {
				t.Errorf("APIKey = %q, want %q", site.APIKey, tt.apiKeyHash)
			}
		})
	}
}

func TestListSites(t *testing.T) {
	tests := []struct {
		name      string
		sites     []struct{ id, name, domain, hash string }
		wantCount int
	}{
		{
			name:      "returns empty list when no sites exist",
			sites:     nil,
			wantCount: 0,
		},
		{
			name: "returns single site",
			sites: []struct{ id, name, domain, hash string }{
				{"list-1", "Site One", "one.example.com", "hash-1"},
			},
			wantCount: 1,
		},
		{
			name: "returns multiple sites",
			sites: []struct{ id, name, domain, hash string }{
				{"list-a", "Alpha", "alpha.example.com", "hash-a"},
				{"list-b", "Beta", "beta.example.com", "hash-b"},
				{"list-c", "Gamma", "gamma.example.com", "hash-c"},
			},
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			ctx := testContext()

			for _, s := range tt.sites {
				_, err := db.CreateSite(ctx, s.id, s.name, s.domain, s.hash)
				if err != nil {
					t.Fatalf("setup CreateSite(%q): %v", s.id, err)
				}
			}

			sites, err := db.ListSites(ctx)

			if err != nil {
				t.Fatalf("ListSites returned error: %v", err)
			}
			if len(sites) != tt.wantCount {
				t.Errorf("len(sites) = %d, want %d", len(sites), tt.wantCount)
			}
		})
	}
}

func TestUpdateSite(t *testing.T) {
	tests := []struct {
		name      string
		newName   string
		newDomain string
	}{
		{
			name:      "updates name and domain",
			newName:   "Updated Name",
			newDomain: "updated.example.com",
		},
		{
			name:      "clears domain",
			newName:   "No Domain",
			newDomain: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			ctx := testContext()

			original, err := db.CreateSite(ctx, "site-upd", "Original", "original.example.com", "hash-upd")
			if err != nil {
				t.Fatalf("setup CreateSite: %v", err)
			}

			err = db.UpdateSite(ctx, original.ID, tt.newName, tt.newDomain)

			if err != nil {
				t.Fatalf("UpdateSite returned error: %v", err)
			}

			updated, err := db.GetSiteByID(ctx, original.ID)
			if err != nil {
				t.Fatalf("GetSiteByID after update: %v", err)
			}
			if updated == nil {
				t.Fatal("GetSiteByID returned nil after update")
			}
			if updated.Name != tt.newName {
				t.Errorf("Name = %q, want %q", updated.Name, tt.newName)
			}
			if tt.newDomain != "" {
				want := sql.NullString{String: tt.newDomain, Valid: true}
				if updated.Domain != want {
					t.Errorf("Domain = %v, want %v", updated.Domain, want)
				}
			} else {
				if updated.Domain.Valid {
					t.Errorf("Domain.Valid = true, want false after clearing domain")
				}
			}
			if updated.APIKey != original.APIKey {
				t.Errorf("APIKey changed: got %q, want %q", updated.APIKey, original.APIKey)
			}
		})
	}
}

func TestDeleteSite(t *testing.T) {
	tests := []struct {
		name     string
		deleteID string
		createID string
	}{
		{
			name:     "deletes existing site",
			createID: "site-del-1",
			deleteID: "site-del-1",
		},
		{
			name:     "no error when deleting nonexistent site",
			createID: "site-del-2",
			deleteID: "nonexistent-del",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			ctx := testContext()

			_, err := db.CreateSite(ctx, tt.createID, "Delete Test", "del.example.com", "hash-del")
			if err != nil {
				t.Fatalf("setup CreateSite: %v", err)
			}

			err = db.DeleteSite(ctx, tt.deleteID)

			if err != nil {
				t.Fatalf("DeleteSite returned error: %v", err)
			}

			site, err := db.GetSiteByID(ctx, tt.deleteID)
			if err != nil {
				t.Fatalf("GetSiteByID after delete: %v", err)
			}
			if tt.createID == tt.deleteID && site != nil {
				t.Errorf("expected nil after delete, got site with ID %q", site.ID)
			}
		})
	}
}

func TestUpdateSiteAPIKey(t *testing.T) {
	tests := []struct {
		name       string
		newKeyHash string
	}{
		{
			name:       "rotates API key hash",
			newKeyHash: "rotated-hash-new",
		},
		{
			name:       "sets empty API key hash",
			newKeyHash: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			ctx := testContext()

			original, err := db.CreateSite(ctx, "site-key", "Key Site", "key.example.com", "original-hash")
			if err != nil {
				t.Fatalf("setup CreateSite: %v", err)
			}

			err = db.UpdateSiteAPIKey(ctx, original.ID, tt.newKeyHash)

			if err != nil {
				t.Fatalf("UpdateSiteAPIKey returned error: %v", err)
			}

			updated, err := db.GetSiteByID(ctx, original.ID)
			if err != nil {
				t.Fatalf("GetSiteByID after key update: %v", err)
			}
			if updated == nil {
				t.Fatal("GetSiteByID returned nil after key update")
			}
			if updated.APIKey != tt.newKeyHash {
				t.Errorf("APIKey = %q, want %q", updated.APIKey, tt.newKeyHash)
			}
			if updated.Name != original.Name {
				t.Errorf("Name changed unexpectedly: got %q, want %q", updated.Name, original.Name)
			}
		})
	}
}
