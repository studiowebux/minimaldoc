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
	Title        string `yaml:"title"`
	Description  string `yaml:"description"`
	BaseURL      string `yaml:"base_url"`
	Author       string `yaml:"author"`
	Theme        string `yaml:"theme"`
	DarkMode     bool   `yaml:"dark_mode"`
	EnableLLMS   bool   `yaml:"enable_llms"`
	EnableSearch bool   `yaml:"enable_search"`
	CleanURLs    bool   `yaml:"clean_urls"`
	Entrypoint   string `yaml:"entrypoint"`

	OpenAPI struct {
		Enabled           bool     `yaml:"enabled"`
		SpecFiles         []string `yaml:"spec_files"`
		SpecURLs          []string `yaml:"spec_urls"`
		DefaultView       string   `yaml:"default_view"`
		SyncOnBuild       bool     `yaml:"sync_on_build"`
		CacheDir          string   `yaml:"cache_dir"`
		EnableTesting     bool     `yaml:"enable_testing"`
		EnableExport      bool     `yaml:"enable_export"`
		EnableCodeSamples bool     `yaml:"enable_code_samples"`
		LazyLoadChunkSize int      `yaml:"lazy_load_chunk_size"`
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

	Landing       core.LandingConfig      `yaml:"landing"`
	ThemeConfig   core.ThemeConfig        `yaml:"theme_config"`
	Portfolio     core.PortfolioConfig    `yaml:"portfolio"`
	Contact       core.ContactConfig      `yaml:"contact"`
	Faq           core.FaqConfig          `yaml:"faq"`
	Legal         core.LegalConfig        `yaml:"legal"`
	KnowledgeBase core.KBConfig           `yaml:"knowledgebase"`
	Waitlist      core.WaitlistConfig     `yaml:"waitlist"`
	Roadmap       core.RoadmapConfig      `yaml:"roadmap"`
	Footer        core.FooterConfig       `yaml:"footer"`
	LinkCheck     core.LinkCheckConfig    `yaml:"link_check"`
	Versions      core.VersionConfig      `yaml:"versions"`
	I18n          core.I18nConfig         `yaml:"i18n"`
	PDFExport     core.PDFExportConfig    `yaml:"pdf_export"`
	ClaudeAssist  core.ClaudeAssistConfig `yaml:"claude_assist"`
	Analytics     core.AnalyticsConfig    `yaml:"analytics"`

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

// MergeWithCLI merges config file with CLI flags.
// CLI flags take precedence over config file values.
func (cfg *FileConfig) MergeWithCLI(cliConfig core.SiteConfig, cliFlags map[string]bool) core.SiteConfig {
	result := cliConfig

	// Basic site settings (CLI flag check)
	mergeString(&result.Title, cfg.Title, cliFlags, "title")
	mergeString(&result.Description, cfg.Description, cliFlags, "description")
	mergeString(&result.BaseURL, cfg.BaseURL, cliFlags, "base-url")
	mergeString(&result.Theme, cfg.Theme, cliFlags, "theme")
	mergeBool(&result.EnableLLMS, cfg.EnableLLMS, cliFlags, "llms")
	mergeBool(&result.CleanURLs, cfg.CleanURLs, cliFlags, "clean-urls")
	mergeStringNoFlag(&result.Entrypoint, cfg.Entrypoint)

	// OpenAPI config
	mergeBool(&result.OpenAPI.Enabled, cfg.OpenAPI.Enabled, cliFlags, "openapi")
	if result.OpenAPI.Enabled {
		mergeStringSlice(&result.OpenAPI.SpecFiles, cfg.OpenAPI.SpecFiles, cliFlags, "openapi-dir")
		mergeStringSliceNoFlag(&result.OpenAPI.SpecURLs, cfg.OpenAPI.SpecURLs)
		mergeStringNoFlag(&result.OpenAPI.DefaultView, cfg.OpenAPI.DefaultView)
		mergeBoolNoFlag(&result.OpenAPI.SyncOnBuild, cfg.OpenAPI.SyncOnBuild)
		mergeStringNoFlag(&result.OpenAPI.CacheDir, cfg.OpenAPI.CacheDir)
		mergeBoolNoFlag(&result.OpenAPI.EnableTesting, cfg.OpenAPI.EnableTesting)
		mergeBoolNoFlag(&result.OpenAPI.EnableExport, cfg.OpenAPI.EnableExport)
		mergeBoolNoFlag(&result.OpenAPI.EnableCodeSamples, cfg.OpenAPI.EnableCodeSamples)
		mergeIntNoFlag(&result.OpenAPI.LazyLoadChunkSize, cfg.OpenAPI.LazyLoadChunkSize)
	}

	// Status config
	mergeBool(&result.Status.Enabled, cfg.Status.Enabled, cliFlags, "status")
	if result.Status.Enabled {
		mergeString(&result.Status.Title, cfg.Status.Title, cliFlags, "status-title")
		mergeString(&result.Status.Description, cfg.Status.Description, cliFlags, "status-description")
		mergeString(&result.Status.Path, cfg.Status.Path, cliFlags, "status-path")
		mergeBoolNoFlag(&result.Status.ShowHistory, cfg.Status.ShowHistory)
		mergeIntNoFlag(&result.Status.HistoryMonths, cfg.Status.HistoryMonths)
		mergeBoolNoFlag(&result.Status.RSSEnabled, cfg.Status.RSSEnabled)
	}

	// Changelog config
	mergeBool(&result.Changelog.Enabled, cfg.Changelog.Enabled, cliFlags, "changelog")
	if result.Changelog.Enabled {
		mergeString(&result.Changelog.Title, cfg.Changelog.Title, cliFlags, "changelog-title")
		mergeString(&result.Changelog.Path, cfg.Changelog.Path, cliFlags, "changelog-path")
		mergeStringNoFlag(&result.Changelog.Description, cfg.Changelog.Description)
		mergeStringNoFlag(&result.Changelog.Repository, cfg.Changelog.Repository)
		mergeBoolNoFlag(&result.Changelog.RSSEnabled, cfg.Changelog.RSSEnabled)
	}

	// Stale Warning config
	mergeBool(&result.StaleWarning.Enabled, cfg.StaleWarning.Enabled, cliFlags, "stale-warning")
	if result.StaleWarning.Enabled {
		mergeInt(&result.StaleWarning.ThresholdDays, cfg.StaleWarning.ThresholdDays, cliFlags, "stale-threshold")
		mergeStringNoFlag(&result.StaleWarning.Message, cfg.StaleWarning.Message)
		mergeBoolNoFlag(&result.StaleWarning.ShowUpdateDate, cfg.StaleWarning.ShowUpdateDate)
	}

	// Complex configs that merge as whole structs when enabled
	if cfg.Landing.Enabled {
		result.Landing = cfg.Landing
	}
	if cfg.ThemeConfig.Name != "" || cfg.ThemeConfig.HasCustomColors() || cfg.ThemeConfig.HasCustomFonts() {
		result.ThemeConfig = mergeThemeConfig(result.ThemeConfig, cfg.ThemeConfig)
	}
	if cfg.Portfolio.Enabled {
		result.Portfolio = cfg.Portfolio
	}
	if cfg.Contact.Enabled {
		result.Contact = cfg.Contact
	}
	if cfg.Faq.Enabled {
		result.Faq = cfg.Faq
	}
	if cfg.Legal.Enabled {
		result.Legal = cfg.Legal
	}
	if cfg.KnowledgeBase.Enabled {
		result.KnowledgeBase = cfg.KnowledgeBase
	}
	if cfg.Waitlist.Enabled {
		result.Waitlist = cfg.Waitlist
	}
	if cfg.Roadmap.Enabled {
		result.Roadmap = cfg.Roadmap
	}
	if cfg.Footer.Copyright != "" || len(cfg.Footer.Links) > 0 || len(cfg.Footer.Social) > 0 {
		result.Footer = cfg.Footer
	}
	if cfg.LinkCheck.Mode != "" {
		result.LinkCheck = cfg.LinkCheck
	}
	if cfg.Versions.Enabled {
		result.Versions = cfg.Versions
	}
	if cfg.I18n.Enabled {
		result.I18n = cfg.I18n
	}
	if cfg.PDFExport.Enabled {
		result.PDFExport = cfg.PDFExport
	}
	if cfg.ClaudeAssist.Enabled {
		result.ClaudeAssist = cfg.ClaudeAssist
		if result.ClaudeAssist.Label == "" {
			result.ClaudeAssist.Label = "Ask Claude"
		}
	}
	if cfg.Analytics.Enabled {
		result.Analytics = cfg.Analytics
	}

	// Social Links
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

// mergeThemeConfig merges theme config, with file config values overriding defaults.
// Uses reflection-based merging for nested structs with string fields.
func mergeThemeConfig(base, override core.ThemeConfig) core.ThemeConfig {
	result := base

	mergeStringNoFlag(&result.Name, override.Name)

	// Merge colors (light and dark)
	result.Colors.Light = mergeColorSet(result.Colors.Light, override.Colors.Light)
	result.Colors.Dark = mergeColorSet(result.Colors.Dark, override.Colors.Dark)

	// Merge fonts and hero using reflection
	mergeStringFields(&result.Fonts, override.Fonts)
	mergeStringFields(&result.Hero, override.Hero)

	return result
}

// mergeColorSet merges color sets, with override values replacing base values.
// Uses reflection to merge all non-empty string fields automatically.
func mergeColorSet(base, override core.ThemeColorSet) core.ThemeColorSet {
	result := base
	mergeStringFields(&result, override)
	return result
}
