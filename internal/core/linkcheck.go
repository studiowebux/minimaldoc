package core

// LinkCheckMode represents how link check errors are handled
type LinkCheckMode string

const (
	LinkCheckError  LinkCheckMode = "error"  // Fail build on broken links
	LinkCheckWarn   LinkCheckMode = "warn"   // Warn but continue build
	LinkCheckIgnore LinkCheckMode = "ignore" // Skip link checking
)

// LinkCheckConfig holds configuration for the link checker
type LinkCheckConfig struct {
	Enabled       bool          `yaml:"enabled"`
	Mode          LinkCheckMode `yaml:"mode"`
	CheckExternal bool          `yaml:"check_external"`
	ExternalTimeout int         `yaml:"external_timeout"` // Seconds
	IgnorePatterns []string     `yaml:"ignore_patterns"`  // Glob patterns to skip
	AllowedBroken  []string     `yaml:"allowed_broken"`   // Known broken links to ignore
}

// DefaultLinkCheckConfig returns a LinkCheckConfig with sensible defaults
func DefaultLinkCheckConfig() LinkCheckConfig {
	return LinkCheckConfig{
		Enabled:        true,
		Mode:           LinkCheckWarn,
		CheckExternal:  false,
		ExternalTimeout: 5,
		IgnorePatterns: []string{},
		AllowedBroken:  []string{},
	}
}

// CollectedLink represents a link found in a source file
type CollectedLink struct {
	URL        string // The link URL
	SourceFile string // Source markdown file path
	Line       int    // Line number in source file
	Column     int    // Column number
	Text       string // Link text (for context)
	LinkType   LinkType
}

// LinkType categorizes links for validation
type LinkType int

const (
	LinkTypeInternalPage   LinkType = iota // /page.html, ./page.md
	LinkTypeInternalAnchor                  // #section
	LinkTypeInternalAsset                   // /images/logo.png
	LinkTypeExternal                        // https://example.com
	LinkTypeEmail                           // mailto:user@example.com
	LinkTypeOther                           // tel:, javascript:, etc.
)

// String returns a human-readable link type
func (lt LinkType) String() string {
	switch lt {
	case LinkTypeInternalPage:
		return "internal page"
	case LinkTypeInternalAnchor:
		return "anchor"
	case LinkTypeInternalAsset:
		return "asset"
	case LinkTypeExternal:
		return "external"
	case LinkTypeEmail:
		return "email"
	default:
		return "other"
	}
}

// BrokenLink represents a link that failed validation
type BrokenLink struct {
	Link   CollectedLink
	Reason string
	Suggestion string // Optional: "did you mean X?"
}

// LinkCheckResult holds the results of link checking
type LinkCheckResult struct {
	TotalLinks   int
	BrokenLinks  []BrokenLink
	SkippedLinks int // Links that matched ignore patterns
	ExternalLinks int // External links (checked or not)
}

// HasErrors returns true if there are broken links
func (r *LinkCheckResult) HasErrors() bool {
	return len(r.BrokenLinks) > 0
}

// BrokenCount returns the number of broken links
func (r *LinkCheckResult) BrokenCount() int {
	return len(r.BrokenLinks)
}
