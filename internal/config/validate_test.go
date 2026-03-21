package config

import (
	"strings"
	"testing"
)

func TestValidate_ValidConfig(t *testing.T) {
	cfg := FileConfig{
		BaseURL: "https://example.com",
		Theme:   "default",
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_EmptyConfig(t *testing.T) {
	cfg := FileConfig{}
	if err := cfg.Validate(); err != nil {
		t.Errorf("empty config should be valid, got: %v", err)
	}
}

func TestValidate_InvalidBaseURL(t *testing.T) {
	cfg := FileConfig{BaseURL: "not-a-url"}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for invalid base_url")
	}
	if !strings.Contains(err.Error(), "base_url") {
		t.Errorf("error should mention base_url, got: %v", err)
	}
}

func TestValidate_FTPBaseURL(t *testing.T) {
	cfg := FileConfig{BaseURL: "ftp://example.com"}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for ftp base_url")
	}
}

func TestValidate_UnknownTheme(t *testing.T) {
	cfg := FileConfig{Theme: "nonexistent"}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for unknown theme")
	}
	if !strings.Contains(err.Error(), "theme") {
		t.Errorf("error should mention theme, got: %v", err)
	}
}

func TestValidate_ValidThemes(t *testing.T) {
	for _, theme := range []string{"default", "yellow"} {
		cfg := FileConfig{Theme: theme}
		if err := cfg.Validate(); err != nil {
			t.Errorf("theme %q should be valid, got: %v", theme, err)
		}
	}
}

func TestValidate_LazyLoadChunkSize_Negative(t *testing.T) {
	v := -1
	cfg := FileConfig{}
	cfg.OpenAPI.LazyLoadChunkSize = &v
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for negative lazy_load_chunk_size")
	}
}

func TestValidate_LazyLoadChunkSize_TooLarge(t *testing.T) {
	v := 20000000
	cfg := FileConfig{}
	cfg.OpenAPI.LazyLoadChunkSize = &v
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for oversized lazy_load_chunk_size")
	}
}

func TestValidate_LazyLoadChunkSize_Zero(t *testing.T) {
	v := 0
	cfg := FileConfig{}
	cfg.OpenAPI.LazyLoadChunkSize = &v
	if err := cfg.Validate(); err == nil {
		t.Errorf("zero lazy_load_chunk_size should be invalid")
	}
}

func TestValidate_SpecURLs_TooMany(t *testing.T) {
	cfg := FileConfig{}
	cfg.OpenAPI.SpecURLs = make([]string, 21)
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for too many spec_urls")
	}
}

func TestValidate_HistoryMonths_Negative(t *testing.T) {
	v := -1
	cfg := FileConfig{}
	cfg.Status.HistoryMonths = &v
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for negative history_months")
	}
}

func TestValidate_ThresholdDays_Negative(t *testing.T) {
	v := -1
	cfg := FileConfig{}
	cfg.StaleWarning.ThresholdDays = &v
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for negative threshold_days")
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	v := -1
	cfg := FileConfig{
		BaseURL: "ftp://bad",
		Theme:   "nope",
	}
	cfg.StaleWarning.ThresholdDays = &v
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected multiple errors")
	}
	if !strings.Contains(err.Error(), "base_url") || !strings.Contains(err.Error(), "theme") || !strings.Contains(err.Error(), "threshold_days") {
		t.Errorf("expected all three errors, got: %v", err)
	}
}
