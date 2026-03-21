package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
	"gopkg.in/yaml.v3"
)

// validThemes lists the built-in theme names.
var validThemes = map[string]bool{
	"default": true,
	"yellow":  true,
}

// SocialLinkConfig represents a social link in config.yaml
type SocialLinkConfig struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Icon string `yaml:"icon"`
}

// FileConfig represents the structure of config.yaml.
// Bool and int fields use pointers to distinguish "not set" from "explicitly set to false/0".
type FileConfig struct {
	Title        string `yaml:"title"`
	Description  string `yaml:"description"`
	BaseURL      string `yaml:"base_url"`
	Author       string `yaml:"author"`
	Theme        string `yaml:"theme"`
	DarkMode     *bool  `yaml:"dark_mode"`
	EnableLLMS   *bool  `yaml:"enable_llms"`
	EnableSearch *bool  `yaml:"enable_search"`
	CleanURLs    *bool  `yaml:"clean_urls"`
	Entrypoint   string `yaml:"entrypoint"`
	NavDepth     *int   `yaml:"nav_depth"`

	OpenAPI struct {
		Enabled           *bool    `yaml:"enabled"`
		SpecFiles         []string `yaml:"spec_files"`
		SpecURLs          []string `yaml:"spec_urls"`
		DefaultView       string   `yaml:"default_view"`
		SyncOnBuild       *bool    `yaml:"sync_on_build"`
		CacheDir          string   `yaml:"cache_dir"`
		EnableTesting     *bool    `yaml:"enable_testing"`
		EnableExport      *bool    `yaml:"enable_export"`
		EnableCodeSamples *bool    `yaml:"enable_code_samples"`
		LazyLoadChunkSize *int     `yaml:"lazy_load_chunk_size"`
	} `yaml:"openapi"`

	Status struct {
		Enabled       *bool  `yaml:"enabled"`
		Title         string `yaml:"title"`
		Description   string `yaml:"description"`
		Path          string `yaml:"path"`
		ShowHistory   *bool  `yaml:"show_history"`
		HistoryMonths *int   `yaml:"history_months"`
		RSSEnabled    *bool  `yaml:"rss_enabled"`
	} `yaml:"status"`

	Changelog struct {
		Enabled     *bool  `yaml:"enabled"`
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
		Path        string `yaml:"path"`
		RSSEnabled  *bool  `yaml:"rss_enabled"`
		Repository  string `yaml:"repository"`
	} `yaml:"changelog"`

	StaleWarning struct {
		Enabled        *bool  `yaml:"enabled"`
		ThresholdDays  *int   `yaml:"threshold_days"`
		Message        string `yaml:"message"`
		ShowUpdateDate *bool  `yaml:"show_update_date"`
	} `yaml:"stale_warning"`

	Landing       core.LandingConfig      `yaml:"landing"`
	ThemeConfig   core.ThemeConfig        `yaml:"theme_config"`
	Portfolio     core.PortfolioConfig    `yaml:"portfolio"`
	Contact       core.ContactConfig      `yaml:"contact"`
	Faq           core.FaqConfig          `yaml:"faq"`
	Legal         core.LegalConfig        `yaml:"legal"`
	KnowledgeBase core.KBConfig           `yaml:"knowledgebase"`
	Roadmap       core.RoadmapConfig      `yaml:"roadmap"`
	Footer        core.FooterConfig       `yaml:"footer"`
	LinkCheck     core.LinkCheckConfig    `yaml:"link_check"`
	Versions      core.VersionConfig      `yaml:"versions"`
	I18n          core.I18nConfig         `yaml:"i18n"`
	PDFExport     core.PDFExportConfig    `yaml:"pdf_export"`
	ClaudeAssist  core.ClaudeAssistConfig `yaml:"claude_assist"`
	Analytics     core.AnalyticsConfig    `yaml:"analytics"`
	MCP           core.MCPConfig          `yaml:"mcp"`

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
	data, err := os.ReadFile(configPath) // #nosec G304 -- path from trusted user configuration
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML with strict mode to catch typos in field names
	var cfg FileConfig
	dec := yaml.NewDecoder(strings.NewReader(string(data)))
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config.yaml: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config.yaml: %w", err)
	}

	return &cfg, nil
}

// Validate checks FileConfig values for correctness after YAML unmarshal.
func (cfg *FileConfig) Validate() error {
	var errs []string

	// base_url must be a valid http/https URL if set
	if cfg.BaseURL != "" {
		u, err := url.Parse(cfg.BaseURL)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			errs = append(errs, fmt.Sprintf("base_url: must be a valid http or https URL, got %q", cfg.BaseURL))
		}
	}

	// theme must be a known built-in theme if set
	if cfg.Theme != "" && !validThemes[cfg.Theme] {
		names := make([]string, 0, len(validThemes))
		for k := range validThemes {
			names = append(names, k)
		}
		errs = append(errs, fmt.Sprintf("theme: unknown theme %q (available: %s)", cfg.Theme, strings.Join(names, ", ")))
	}

	// OpenAPI bounds
	if cfg.OpenAPI.LazyLoadChunkSize != nil && (*cfg.OpenAPI.LazyLoadChunkSize <= 0 || *cfg.OpenAPI.LazyLoadChunkSize > 10000000) {
		errs = append(errs, fmt.Sprintf("openapi.lazy_load_chunk_size: must be 1-10000000, got %d", *cfg.OpenAPI.LazyLoadChunkSize))
	}
	if len(cfg.OpenAPI.SpecURLs) > 20 {
		errs = append(errs, fmt.Sprintf("openapi.spec_urls: max 20 entries, got %d", len(cfg.OpenAPI.SpecURLs)))
	}
	for i, specURL := range cfg.OpenAPI.SpecURLs {
		u, err := url.Parse(specURL)
		if err != nil {
			errs = append(errs, fmt.Sprintf("openapi.spec_urls[%d]: invalid URL: %s", i, specURL))
		} else if u.Scheme != "http" && u.Scheme != "https" {
			errs = append(errs, fmt.Sprintf("openapi.spec_urls[%d]: must use http or https scheme, got %q", i, u.Scheme))
		}
	}
	if len(cfg.OpenAPI.SpecFiles) > 50 {
		errs = append(errs, fmt.Sprintf("openapi.spec_files: max 50 entries, got %d", len(cfg.OpenAPI.SpecFiles)))
	}

	// Status bounds
	if cfg.Status.HistoryMonths != nil && (*cfg.Status.HistoryMonths < 0 || *cfg.Status.HistoryMonths > 120) {
		errs = append(errs, fmt.Sprintf("status.history_months: must be 0-120, got %d", *cfg.Status.HistoryMonths))
	}

	// Stale warning bounds
	if cfg.StaleWarning.ThresholdDays != nil && (*cfg.StaleWarning.ThresholdDays < 0 || *cfg.StaleWarning.ThresholdDays > 3650) {
		errs = append(errs, fmt.Sprintf("stale_warning.threshold_days: must be 0-3650, got %d", *cfg.StaleWarning.ThresholdDays))
	}

	// Contact email must look like an email if set
	if cfg.Contact.Email != "" && !strings.Contains(cfg.Contact.Email, "@") {
		errs = append(errs, fmt.Sprintf("contact.email: invalid email address %q", cfg.Contact.Email))
	}

	// Path fields must be relative, no traversal
	pathFields := map[string]string{
		"status.path":        cfg.Status.Path,
		"changelog.path":     cfg.Changelog.Path,
		"faq.path":           cfg.Faq.Path,
		"contact.path":       cfg.Contact.Path,
		"portfolio.path":     cfg.Portfolio.Path,
		"knowledgebase.path": cfg.KnowledgeBase.Path,
		"legal.path":         cfg.Legal.Path,
		"roadmap.path":       cfg.Roadmap.Path,
		"mcp.path":           cfg.MCP.Path,
		"openapi.cache_dir":  cfg.OpenAPI.CacheDir,
	}
	for name, val := range pathFields {
		if val != "" && !isSafePath(val) {
			errs = append(errs, fmt.Sprintf("%s: must be a relative path without traversal, got %q", name, val))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return nil
}

// isSafePath checks that a path is relative, contains no traversal, and has no absolute prefix.
func isSafePath(p string) bool {
	if filepath.IsAbs(p) {
		return false
	}
	cleaned := filepath.Clean(p)
	if strings.HasPrefix(cleaned, "..") {
		return false
	}
	if strings.Contains(cleaned, ".."+string(filepath.Separator)) {
		return false
	}
	return true
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
	mergeBool(&result.DarkMode, cfg.DarkMode, cliFlags, "dark-mode")
	mergeBool(&result.EnableLLMS, cfg.EnableLLMS, cliFlags, "llms")
	mergeBool(&result.EnableSearch, cfg.EnableSearch, cliFlags, "search")
	mergeBool(&result.CleanURLs, cfg.CleanURLs, cliFlags, "clean-urls")
	mergeStringNoFlag(&result.Entrypoint, cfg.Entrypoint)
	if cfg.NavDepth != nil {
		result.NavDepth = *cfg.NavDepth
	}

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
		if result.Portfolio.Path == "" {
			result.Portfolio.Path = core.DefaultPortfolioConfig().Path
		}
	}
	if cfg.Contact.Enabled {
		result.Contact = cfg.Contact
		if result.Contact.Path == "" {
			result.Contact.Path = core.DefaultContactConfig().Path
		}
	}
	if cfg.Faq.Enabled {
		result.Faq = cfg.Faq
		if result.Faq.Path == "" {
			result.Faq.Path = core.DefaultFaqConfig().Path
		}
	}
	if cfg.Legal.Enabled {
		result.Legal = cfg.Legal
		if result.Legal.Path == "" {
			result.Legal.Path = core.DefaultLegalConfig().Path
		}
	}
	if cfg.KnowledgeBase.Enabled {
		result.KnowledgeBase = cfg.KnowledgeBase
		if result.KnowledgeBase.Path == "" {
			result.KnowledgeBase.Path = core.DefaultKBConfig().Path
		}
	}
	if cfg.Roadmap.Enabled {
		result.Roadmap = cfg.Roadmap
		if result.Roadmap.Path == "" {
			result.Roadmap.Path = core.DefaultRoadmapConfig().Path
		}
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

	// MCP config: bare bool cannot distinguish "false" from "unset",
	// so only override the default/CLI when the file explicitly enables MCP.
	if cfg.MCP.Enabled {
		result.MCP.Enabled = true
	}
	if result.MCP.Enabled {
		mergeStringSlice(&result.MCP.SpecFiles, cfg.MCP.SpecFiles, cliFlags, "mcp-dir")
		mergeStringNoFlag(&result.MCP.Path, cfg.MCP.Path)
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
