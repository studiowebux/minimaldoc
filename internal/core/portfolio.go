package core

import (
	"sort"
	"time"
)

// PortfolioPage represents the complete portfolio page data
type PortfolioPage struct {
	Config           PortfolioConfig
	Projects         []Project
	FeaturedProjects []Project
	Tags             []string // All unique tags for filtering
}

// PortfolioConfig holds configuration for the portfolio page
type PortfolioConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
}

// DefaultPortfolioConfig returns a PortfolioConfig with sensible defaults
func DefaultPortfolioConfig() PortfolioConfig {
	return PortfolioConfig{
		Enabled:     false,
		Title:       "Portfolio",
		Description: "Projects and experiments",
		Path:        "portfolio",
	}
}

// Project represents a portfolio project
type Project struct {
	// Identity
	ID       string `yaml:"-"`
	Slug     string `yaml:"-"`
	FilePath string `yaml:"-"`

	// Metadata from frontmatter
	Title       string       `yaml:"title"`
	Description string       `yaml:"description"`
	Image       string       `yaml:"image"`
	Tags        []string     `yaml:"tags"`
	Links       []SimpleLink `yaml:"links"`
	Date        time.Time    `yaml:"date"`
	Featured    bool         `yaml:"featured"`
	Order       int          `yaml:"order"`

	// Content
	RawMD string `yaml:"-"`
	HTML  string `yaml:"-"`

	// Output
	OutputPath string `yaml:"-"`
}

// FilterFeaturedProjects returns only featured projects
func FilterFeaturedProjects(projects []Project) []Project {
	var featured []Project
	for _, p := range projects {
		if p.Featured {
			featured = append(featured, p)
		}
	}
	return featured
}

// CollectProjectTags returns all unique tags from projects
func CollectProjectTags(projects []Project) []string {
	tagMap := make(map[string]bool)
	for _, p := range projects {
		for _, tag := range p.Tags {
			tagMap[tag] = true
		}
	}

	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}
