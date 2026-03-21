package core

// LegalConfig holds configuration for legal pages
type LegalConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Path        string `yaml:"path"`
	FooterGroup string `yaml:"footer_group"`
}

// DefaultLegalConfig returns a LegalConfig with sensible defaults
func DefaultLegalConfig() LegalConfig {
	return LegalConfig{
		Enabled:     false,
		Path:        "legal",
		FooterGroup: "Legal",
	}
}

// LegalPage represents a single legal page (privacy, terms, etc.)
type LegalPage struct {
	// Identity
	Slug     string `yaml:"-"`
	FilePath string `yaml:"-"`

	// Metadata from frontmatter
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Order       int    `yaml:"order"`

	// Content
	RawMD string `yaml:"-"`
	HTML  string `yaml:"-"`

	// Output
	OutputPath string `yaml:"-"`
}
