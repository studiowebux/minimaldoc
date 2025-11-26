package core

// Site represents the entire documentation site
type Site struct {
	// Configuration
	Config SiteConfig

	// Content
	Pages      []*Page      // All pages
	RootPages  []*Page      // Top-level pages (for navigation)
	Navigation *Navigation  // Site navigation tree
	APISpecs   []*APISpec   // OpenAPI specifications

	// Paths
	DocsRoot   string // Root directory of markdown files
	OutputRoot string // Output directory for generated site
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

	// Custom
	Custom map[string]interface{} `yaml:"custom"`
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
		Custom:       make(map[string]interface{}),
	}
}

// Navigation represents the site navigation tree
type Navigation struct {
	Items []*NavItem
}

// NavItem represents a single navigation item
type NavItem struct {
	Title    string     // Display title
	Path     string     // URL path
	Order    int        // Sort order
	Active   bool       // Is current page
	Children []*NavItem // Nested items
	Page     *Page      // Associated page
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
