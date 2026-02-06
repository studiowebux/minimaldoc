package core

// Site represents the entire documentation site
type Site struct {
	// Configuration
	Config SiteConfig

	// Content
	Pages         []*Page        // All pages
	RootPages     []*Page        // Top-level pages (for navigation)
	Navigation    *Navigation    // Site navigation tree
	APISpecs      []*APISpec     // OpenAPI specifications
	StatusPage    *StatusPage    // Status page data (if enabled)
	ChangelogPage *ChangelogPage // Changelog data (if enabled)
	LandingPage   *LandingPage   // Landing page data (if enabled)
	PortfolioPage *PortfolioPage // Portfolio page data (if enabled)
	ContactPage   *ContactPage   // Contact page data (if enabled)

	// Paths
	DocsRoot   string // Root directory of markdown files
	OutputRoot string // Output directory for generated site
}

// SocialLink represents a social media or external link
type SocialLink struct {
	Name string `yaml:"name"` // Display name (e.g., "GitHub", "Twitter")
	URL  string `yaml:"url"`  // Link URL
	Icon string `yaml:"icon"` // Icon identifier (github, twitter, linkedin, etc.)
}

// SiteConfig holds the site-wide configuration
type SiteConfig struct {
	// Basic info
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	BaseURL     string `yaml:"base_url"`
	Author      string `yaml:"author"`

	// Theme
	Theme     string `yaml:"theme"`      // Theme name (default: "default")
	DarkMode  bool   `yaml:"dark_mode"`  // Enable dark mode by default

	// Features
	EnableLLMS  bool `yaml:"enable_llms"`  // Generate llms.txt
	EnableSearch bool `yaml:"enable_search"` // Enable search (future)

	// Navigation
	NavDepth int `yaml:"nav_depth"` // Max depth for navigation tree (0 = unlimited)

	// Output
	CleanURLs bool `yaml:"clean_urls"` // Use /page/ instead of /page.html

	// OpenAPI
	OpenAPI OpenAPIConfig `yaml:"openapi"` // OpenAPI/Swagger configuration

	// Status
	Status StatusConfig `yaml:"status"` // Status page configuration

	// Changelog
	Changelog ChangelogConfig `yaml:"changelog"` // Changelog configuration

	// Stale Warning
	StaleWarning StaleWarningConfig `yaml:"stale_warning"` // Stale content warning configuration

	// Landing Page
	Landing LandingConfig `yaml:"landing"` // Landing page configuration

	// Portfolio
	Portfolio PortfolioConfig `yaml:"portfolio"` // Portfolio page configuration

	// Contact
	Contact ContactConfig `yaml:"contact"` // Contact page configuration

	// Footer (for landing pages)
	Footer FooterConfig `yaml:"footer"` // Footer configuration

	// Social Links
	SocialLinks []SocialLink `yaml:"social_links"` // Social media links in sidebar

	// Custom
	Custom map[string]interface{} `yaml:"custom"`
}

// StaleWarningConfig holds configuration for stale content warnings
type StaleWarningConfig struct {
	Enabled        bool   `yaml:"enabled"`
	ThresholdDays  int    `yaml:"threshold_days"`   // Days before content is considered stale
	Message        string `yaml:"message"`          // Custom message (optional)
	ShowUpdateDate bool   `yaml:"show_update_date"` // Show the actual last update date
}

// DefaultStaleWarningConfig returns a StaleWarningConfig with sensible defaults
func DefaultStaleWarningConfig() StaleWarningConfig {
	return StaleWarningConfig{
		Enabled:        false,
		ThresholdDays:  365,
		Message:        "",
		ShowUpdateDate: true,
	}
}

// DefaultSiteConfig returns a SiteConfig with sensible defaults
func DefaultSiteConfig() SiteConfig {
	return SiteConfig{
		Title:        "Documentation",
		Description:  "Documentation site powered by Minimal Doc",
		Theme:        "default",
		DarkMode:     false,
		EnableLLMS:   true,
		EnableSearch: false,
		NavDepth:     0,
		CleanURLs:    false,
		OpenAPI:      DefaultOpenAPIConfig(),
		Status:       DefaultStatusConfig(),
		Changelog:    DefaultChangelogConfig(),
		StaleWarning: DefaultStaleWarningConfig(),
		Landing:      DefaultLandingConfig(),
		Portfolio:    DefaultPortfolioConfig(),
		Contact:      DefaultContactConfig(),
		Footer:       DefaultFooterConfig(),
		Custom:       make(map[string]interface{}),
	}
}

// Navigation represents the site navigation tree
type Navigation struct {
	Items []*NavItem
}

// NavItem represents a single navigation item
type NavItem struct {
	Title      string     // Display title
	Path       string     // URL path
	Order      int        // Sort order
	Active     bool       // Is current page
	IsExternal bool       // True if external URL (opens in new tab)
	Children   []*NavItem // Nested items
	Page       *Page      // Associated page
}

// NewSite creates a new Site instance
func NewSite(docsRoot, outputRoot string, config SiteConfig) *Site {
	return &Site{
		Config:     config,
		Pages:      []*Page{},
		RootPages:  []*Page{},
		Navigation: &Navigation{Items: []*NavItem{}},
		APISpecs:   []*APISpec{},
		DocsRoot:   docsRoot,
		OutputRoot: outputRoot,
	}
}
