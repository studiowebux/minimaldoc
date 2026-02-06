package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/studiowebux/minimaldoc/internal/core"
	"gopkg.in/yaml.v3"
)

// SocialLinkConfig represents a social link in config.yaml
type SocialLinkConfig struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Icon string `yaml:"icon"`
}

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

	Status struct {
		Enabled       bool   `yaml:"enabled"`
		Title         string `yaml:"title"`
		Description   string `yaml:"description"`
		Path          string `yaml:"path"`
		ShowHistory   bool   `yaml:"show_history"`
		HistoryMonths int    `yaml:"history_months"`
		RSSEnabled    bool   `yaml:"rss_enabled"`
	} `yaml:"status"`

	Changelog struct {
		Enabled     bool   `yaml:"enabled"`
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
		Path        string `yaml:"path"`
		RSSEnabled  bool   `yaml:"rss_enabled"`
		Repository  string `yaml:"repository"`
	} `yaml:"changelog"`

	StaleWarning struct {
		Enabled        bool   `yaml:"enabled"`
		ThresholdDays  int    `yaml:"threshold_days"`
		Message        string `yaml:"message"`
		ShowUpdateDate bool   `yaml:"show_update_date"`
	} `yaml:"stale_warning"`

	Landing   core.LandingConfig   `yaml:"landing"`
	Portfolio core.PortfolioConfig `yaml:"portfolio"`
	Contact   core.ContactConfig   `yaml:"contact"`
	Footer    core.FooterConfig    `yaml:"footer"`

	SocialLinks []SocialLinkConfig `yaml:"social_links"`
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

	// Merge Status config
	if !cliFlags["status"] {
		result.Status.Enabled = cfg.Status.Enabled
	}

	// If Status is enabled (either via CLI or config), merge Status settings
	if result.Status.Enabled {
		if !cliFlags["status-title"] && cfg.Status.Title != "" {
			result.Status.Title = cfg.Status.Title
		}
		if !cliFlags["status-description"] && cfg.Status.Description != "" {
			result.Status.Description = cfg.Status.Description
		}
		if !cliFlags["status-path"] && cfg.Status.Path != "" {
			result.Status.Path = cfg.Status.Path
		}
		if cfg.Status.ShowHistory {
			result.Status.ShowHistory = cfg.Status.ShowHistory
		}
		if cfg.Status.HistoryMonths > 0 {
			result.Status.HistoryMonths = cfg.Status.HistoryMonths
		}
		if cfg.Status.RSSEnabled {
			result.Status.RSSEnabled = cfg.Status.RSSEnabled
		}
	}

	// Merge Changelog config
	if !cliFlags["changelog"] {
		result.Changelog.Enabled = cfg.Changelog.Enabled
	}

	// If Changelog is enabled (either via CLI or config), merge Changelog settings
	if result.Changelog.Enabled {
		if !cliFlags["changelog-title"] && cfg.Changelog.Title != "" {
			result.Changelog.Title = cfg.Changelog.Title
		}
		if !cliFlags["changelog-path"] && cfg.Changelog.Path != "" {
			result.Changelog.Path = cfg.Changelog.Path
		}
		if cfg.Changelog.Description != "" {
			result.Changelog.Description = cfg.Changelog.Description
		}
		if cfg.Changelog.Repository != "" {
			result.Changelog.Repository = cfg.Changelog.Repository
		}
		if cfg.Changelog.RSSEnabled {
			result.Changelog.RSSEnabled = cfg.Changelog.RSSEnabled
		}
	}

	// Merge Stale Warning config
	if !cliFlags["stale-warning"] {
		result.StaleWarning.Enabled = cfg.StaleWarning.Enabled
	}

	// If Stale Warning is enabled, merge settings
	if result.StaleWarning.Enabled {
		if !cliFlags["stale-threshold"] && cfg.StaleWarning.ThresholdDays > 0 {
			result.StaleWarning.ThresholdDays = cfg.StaleWarning.ThresholdDays
		}
		if cfg.StaleWarning.Message != "" {
			result.StaleWarning.Message = cfg.StaleWarning.Message
		}
		if cfg.StaleWarning.ShowUpdateDate {
			result.StaleWarning.ShowUpdateDate = cfg.StaleWarning.ShowUpdateDate
		}
	}

	// Merge Landing config
	if cfg.Landing.Enabled {
		result.Landing = cfg.Landing
	}

	// Merge Portfolio config
	if cfg.Portfolio.Enabled {
		result.Portfolio = cfg.Portfolio
	}

	// Merge Contact config
	if cfg.Contact.Enabled {
		result.Contact = cfg.Contact
	}

	// Merge Footer config
	if cfg.Footer.Copyright != "" || len(cfg.Footer.Links) > 0 || len(cfg.Footer.Social) > 0 {
		result.Footer = cfg.Footer
	}

	// Merge Social Links
	if len(cfg.SocialLinks) > 0 {
		result.SocialLinks = make([]core.SocialLink, len(cfg.SocialLinks))
		for i, link := range cfg.SocialLinks {
			result.SocialLinks[i] = core.SocialLink{
				Name: link.Name,
				URL:  link.URL,
				Icon: link.Icon,
			}
		}
	}

	return result
}
