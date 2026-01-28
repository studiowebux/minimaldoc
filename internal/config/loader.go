package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/studiowebux/minimaldoc/internal/core"
	"gopkg.in/yaml.v3"
)

// FileConfig represents the structure of config.yaml
type FileConfig struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	BaseURL     string `yaml:"base_url"`
	Author      string `yaml:"author"`
	Theme       string `yaml:"theme"`
	DarkMode    bool   `yaml:"dark_mode"`
	EnableLLMS  bool   `yaml:"enable_llms"`
	EnableSearch bool  `yaml:"enable_search"`
	CleanURLs   bool   `yaml:"clean_urls"`

	OpenAPI struct {
		Enabled            bool     `yaml:"enabled"`
		SpecFiles          []string `yaml:"spec_files"`
		SpecURLs           []string `yaml:"spec_urls"`
		DefaultView        string   `yaml:"default_view"`
		SyncOnBuild        bool     `yaml:"sync_on_build"`
		CacheDir           string   `yaml:"cache_dir"`
		EnableTesting      bool     `yaml:"enable_testing"`
		EnableExport       bool     `yaml:"enable_export"`
		EnableCodeSamples  bool     `yaml:"enable_code_samples"`
		LazyLoadChunkSize  int      `yaml:"lazy_load_chunk_size"`
	} `yaml:"openapi"`
}

// LoadConfig loads config.yaml from the docs directory if it exists
// Returns nil if no config file is found (not an error)
func LoadConfig(docsDir string) (*FileConfig, error) {
	configPath := filepath.Join(docsDir, "config.yaml")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// No config file - not an error
		return nil, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var cfg FileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config.yaml: %w", err)
	}

	return &cfg, nil
}

// MergeWithCLI merges config file with CLI flags
// CLI flags take precedence over config file values
func (cfg *FileConfig) MergeWithCLI(cliConfig core.SiteConfig, cliFlags map[string]bool) core.SiteConfig {
	result := cliConfig

	// Only use config file values if CLI flag was not explicitly set
	if !cliFlags["title"] && cfg.Title != "" {
		result.Title = cfg.Title
	}
	if !cliFlags["description"] && cfg.Description != "" {
		result.Description = cfg.Description
	}
	if !cliFlags["base-url"] && cfg.BaseURL != "" {
		result.BaseURL = cfg.BaseURL
	}
	if !cliFlags["theme"] && cfg.Theme != "" {
		result.Theme = cfg.Theme
	}
	if !cliFlags["llms"] && cfg.EnableLLMS {
		result.EnableLLMS = cfg.EnableLLMS
	}
	if !cliFlags["clean-urls"] && cfg.CleanURLs {
		result.CleanURLs = cfg.CleanURLs
	}

	// Merge OpenAPI config
	if !cliFlags["openapi"] {
		result.OpenAPI.Enabled = cfg.OpenAPI.Enabled
	}

	// If OpenAPI is enabled (either via CLI or config), merge OpenAPI settings
	if result.OpenAPI.Enabled {
		if !cliFlags["openapi-dir"] && len(cfg.OpenAPI.SpecFiles) > 0 {
			result.OpenAPI.SpecFiles = cfg.OpenAPI.SpecFiles
		}
		if len(cfg.OpenAPI.SpecURLs) > 0 {
			result.OpenAPI.SpecURLs = cfg.OpenAPI.SpecURLs
		}
		if cfg.OpenAPI.DefaultView != "" {
			result.OpenAPI.DefaultView = cfg.OpenAPI.DefaultView
		}
		if cfg.OpenAPI.SyncOnBuild {
			result.OpenAPI.SyncOnBuild = cfg.OpenAPI.SyncOnBuild
		}
		if cfg.OpenAPI.CacheDir != "" {
			result.OpenAPI.CacheDir = cfg.OpenAPI.CacheDir
		}
		if cfg.OpenAPI.EnableTesting {
			result.OpenAPI.EnableTesting = cfg.OpenAPI.EnableTesting
		}
		if cfg.OpenAPI.EnableExport {
			result.OpenAPI.EnableExport = cfg.OpenAPI.EnableExport
		}
		if cfg.OpenAPI.EnableCodeSamples {
			result.OpenAPI.EnableCodeSamples = cfg.OpenAPI.EnableCodeSamples
		}
		if cfg.OpenAPI.LazyLoadChunkSize > 0 {
			result.OpenAPI.LazyLoadChunkSize = cfg.OpenAPI.LazyLoadChunkSize
		}
	}

	return result
}
