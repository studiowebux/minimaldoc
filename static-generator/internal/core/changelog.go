package core

import (
	"strings"
	"time"
)

// ChangelogPage represents the complete changelog data
type ChangelogPage struct {
	Config      ChangelogConfig
	Releases    []Release
	LastUpdated time.Time
}

// ChangelogConfig holds configuration for the changelog
type ChangelogConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Path        string `yaml:"path"`        // output path (default: "changelog")
	RSSEnabled  bool   `yaml:"rss_enabled"` // generate RSS feed
	Repository  string `yaml:"repository"`  // GitHub repository URL for compare links
}

// DefaultChangelogConfig returns a ChangelogConfig with sensible defaults
func DefaultChangelogConfig() ChangelogConfig {
	return ChangelogConfig{
		Enabled:     false,
		Title:       "Changelog",
		Description: "All notable changes to this project",
		Path:        "changelog",
		RSSEnabled:  true,
		Repository:  "",
	}
}

// Release represents a single version release
type Release struct {
	// Identity
	Version string `yaml:"version"`
	Slug    string `yaml:"-"` // URL-friendly version (e.g., "1.2.0")

	// Metadata from frontmatter
	Date       time.Time `yaml:"date"`
	Title      string    `yaml:"title"`      // optional title for the release
	Prerelease bool      `yaml:"prerelease"` // is this a prerelease version

	// Semantic version components (for sorting)
	Major int `yaml:"-"`
	Minor int `yaml:"-"`
	Patch int `yaml:"-"`

	// Content
	Categories []ChangeCategory `yaml:"-"` // changes grouped by category
	RawMD      string           `yaml:"-"` // raw markdown content
	HTML       string           `yaml:"-"` // rendered HTML

	// Generated
	CompareURL string `yaml:"-"` // GitHub compare link to previous version
	OutputPath string `yaml:"-"` // generated HTML path
}

// ChangeCategory represents a category of changes (Added, Changed, etc.)
type ChangeCategory struct {
	Type    ChangeType
	Entries []ChangeEntry
}

// ChangeEntry represents a single change within a category
type ChangeEntry struct {
	Description string // the change description
	RawMD       string // raw markdown
	HTML        string // rendered HTML
}

// ChangeType represents the type of change following Keep a Changelog format
type ChangeType string

const (
	ChangeAdded      ChangeType = "Added"
	ChangeChanged    ChangeType = "Changed"
	ChangeDeprecated ChangeType = "Deprecated"
	ChangeRemoved    ChangeType = "Removed"
	ChangeFixed      ChangeType = "Fixed"
	ChangeSecurity   ChangeType = "Security"
)

// AllChangeTypes returns all valid change types in standard order
func AllChangeTypes() []ChangeType {
	return []ChangeType{
		ChangeAdded,
		ChangeChanged,
		ChangeDeprecated,
		ChangeRemoved,
		ChangeFixed,
		ChangeSecurity,
	}
}

// IsValid checks if the change type is valid
func (c ChangeType) IsValid() bool {
	switch c {
	case ChangeAdded, ChangeChanged, ChangeDeprecated, ChangeRemoved, ChangeFixed, ChangeSecurity:
		return true
	}
	return false
}

// ParseChangeType converts a string to a ChangeType (case-insensitive)
func ParseChangeType(s string) (ChangeType, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "added":
		return ChangeAdded, true
	case "changed":
		return ChangeChanged, true
	case "deprecated":
		return ChangeDeprecated, true
	case "removed":
		return ChangeRemoved, true
	case "fixed":
		return ChangeFixed, true
	case "security":
		return ChangeSecurity, true
	}
	return "", false
}

// Color returns the CSS color class for the change type
func (c ChangeType) Color() string {
	switch c {
	case ChangeAdded:
		return "green"
	case ChangeChanged:
		return "blue"
	case ChangeDeprecated:
		return "yellow"
	case ChangeRemoved:
		return "red"
	case ChangeFixed:
		return "purple"
	case ChangeSecurity:
		return "orange"
	}
	return "gray"
}

// Label returns a human-readable label for the change type
func (c ChangeType) Label() string {
	return string(c)
}
