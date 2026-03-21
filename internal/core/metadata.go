package core

// Metadata represents the frontmatter data from a markdown file
type Metadata struct {
	// Basic metadata
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
	Author      string   `yaml:"author"`
	Date        string   `yaml:"date"`

	// Menu/Navigation overrides
	MenuTitle string `yaml:"menu_title"` // Override display name in navigation
	MenuOrder int    `yaml:"menu_order"` // Override order (takes precedence over file numbering)
	Hidden    bool   `yaml:"hidden"`     // Hide from navigation

	// Layout overrides
	FullWidth bool `yaml:"full_width"` // Hide sidebar, use full-width layout
	NoHeader  bool `yaml:"no_header"`  // Hide page header (title, buttons, description)

	// SEO metadata
	SEO SEO `yaml:"seo"`

	// OpenAPI integration
	OpenAPISpec   string `yaml:"openapi_spec"`   // Reference to OpenAPI spec file/URL
	OpenAPIPath   string `yaml:"openapi_path"`   // Specific endpoint path to embed
	OpenAPIMethod string `yaml:"openapi_method"` // Specific HTTP method to embed

	// Stale warning overrides
	StaleWarning       *bool `yaml:"stale_warning"`        // Override site setting (nil = use site default)
	StaleThresholdDays *int  `yaml:"stale_threshold_days"` // Override threshold (nil = use site default)

	// Portfolio-specific fields
	Image    string         `yaml:"image"`    // Project image
	Links    []MetadataLink `yaml:"links"`    // Project links
	Featured bool           `yaml:"featured"` // Featured project

	// FAQ-specific fields
	Question string `yaml:"question"` // FAQ question text
	Category string `yaml:"category"` // FAQ category name

	// Version-specific fields
	Versions     []string `yaml:"versions"`      // Versions this page appears in (empty = all)
	Since        string   `yaml:"since"`         // Version this feature was introduced
	DeprecatedIn string   `yaml:"deprecated_in"` // Version where feature was deprecated
	RemovedIn    string   `yaml:"removed_in"`    // Version where feature was removed
	VersionNote  string   `yaml:"version_note"`  // Version to show note for

	// Custom fields (for extensibility)
	Custom map[string]any `yaml:"custom"`
}

// MetadataLink represents a link in metadata
type MetadataLink struct {
	Text string `yaml:"text"`
	URL  string `yaml:"url"`
}

// SEO represents SEO-specific metadata
type SEO struct {
	Title       string   `yaml:"title"`       // SEO title (og:title)
	Description string   `yaml:"description"` // SEO description (og:description)
	Keywords    []string `yaml:"keywords"`    // Meta keywords
	Image       string   `yaml:"image"`       // og:image
	Author      string   `yaml:"author"`      // Author
	Canonical   string   `yaml:"canonical"`   // Canonical URL
	NoIndex     bool     `yaml:"noindex"`     // Don't index this page
	NoFollow    bool     `yaml:"nofollow"`    // Don't follow links on this page
}

// DefaultMetadata returns a Metadata instance with sensible defaults
func DefaultMetadata() Metadata {
	return Metadata{
		Title:       "Untitled",
		Description: "",
		Tags:        []string{},
		MenuOrder:   -1, // -1 means "not set", use filename-based order
		Hidden:      false,
		SEO:         SEO{},
		Custom:      make(map[string]any),
	}
}
