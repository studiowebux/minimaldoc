package parser

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/studiowebux/minimaldoc/internal/core"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	htmlrenderer "github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

// ChangelogParser handles parsing of changelog content
type ChangelogParser struct {
	md goldmark.Markdown
}

// NewChangelogParser creates a new changelog parser
func NewChangelogParser() *ChangelogParser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			htmlrenderer.WithUnsafe(),
		),
	)

	return &ChangelogParser{md: md}
}

// ChangelogSourceDir is the directory name for changelog content
const ChangelogSourceDir = "__changelog__"

// ParseChangelogDir parses all changelog content from a directory
func (p *ChangelogParser) ParseChangelogDir(changelogDir string) (*core.ChangelogPage, error) {
	changelogPage := &core.ChangelogPage{
		Config:      core.DefaultChangelogConfig(),
		LastUpdated: time.Now(),
	}

	// Check if changelog directory exists
	if _, err := os.Stat(changelogDir); os.IsNotExist(err) {
		return changelogPage, nil // Return empty changelog if dir doesn't exist
	}

	// Note: Config is set from main site config.yaml, not from __changelog__/config.yaml

	// Parse releases
	releasesDir := filepath.Join(changelogDir, "releases")
	if _, err := os.Stat(releasesDir); err == nil {
		releases, err := p.ParseReleases(releasesDir)
		if err != nil {
			return nil, fmt.Errorf("failed to parse releases: %w", err)
		}
		changelogPage.Releases = releases
	}

	return changelogPage, nil
}

// ParseConfig parses the changelog config.yaml file
func (p *ChangelogParser) ParseConfig(path string) (core.ChangelogConfig, error) {
	config := core.DefaultChangelogConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	return config, nil
}

// ParseReleases parses all release markdown files in a directory
func (p *ChangelogParser) ParseReleases(dir string) ([]core.Release, error) {
	var releases []core.Release

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		release, err := p.ParseRelease(path)
		if err != nil {
			return nil, fmt.Errorf("failed to parse release %s: %w", path, err)
		}

		releases = append(releases, release)
	}

	return releases, nil
}

// ReleaseFrontmatter represents the YAML frontmatter for releases
type ReleaseFrontmatter struct {
	Version    string    `yaml:"version"`
	Date       time.Time `yaml:"date"`
	Title      string    `yaml:"title"`
	Prerelease bool      `yaml:"prerelease"`
}

// ParseRelease parses a single release markdown file
func (p *ChangelogParser) ParseRelease(path string) (core.Release, error) {
	var release core.Release

	// Generate slug from filename (without extension)
	filename := filepath.Base(path)
	release.Slug = strings.TrimSuffix(filename, ".md")

	// Read file
	content, err := os.ReadFile(path)
	if err != nil {
		return release, err
	}

	// Parse frontmatter
	var fm ReleaseFrontmatter
	rest, err := frontmatter.Parse(bytes.NewReader(content), &fm)
	if err != nil {
		return release, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Copy frontmatter to release
	release.Version = fm.Version
	if release.Version == "" {
		// Use slug as version if not specified
		release.Version = release.Slug
	}
	release.Date = fm.Date
	release.Title = fm.Title
	release.Prerelease = fm.Prerelease

	// Parse semantic version
	release.Major, release.Minor, release.Patch = parseSemver(release.Version)

	// Store raw markdown
	release.RawMD = string(rest)

	// Parse categories from markdown content
	release.Categories = p.parseCategories(rest)

	// Render full content to HTML
	var buf bytes.Buffer
	if err := p.md.Convert(rest, &buf); err != nil {
		return release, fmt.Errorf("failed to render markdown: %w", err)
	}
	release.HTML = buf.String()

	return release, nil
}

// categoryHeaderRegex matches h2 headers like "## Added" or "## Fixed"
var categoryHeaderRegex = regexp.MustCompile(`(?m)^##\s+(\w+)\s*$`)

// parseCategories extracts change categories from markdown content
func (p *ChangelogParser) parseCategories(content []byte) []core.ChangeCategory {
	var categories []core.ChangeCategory

	lines := strings.Split(string(content), "\n")
	var currentCategory *core.ChangeCategory
	var currentEntries []string

	for _, line := range lines {
		if matches := categoryHeaderRegex.FindStringSubmatch(line); matches != nil {
			// Save previous category if exists
			if currentCategory != nil {
				currentCategory.Entries = p.parseEntries(currentEntries)
				if len(currentCategory.Entries) > 0 {
					categories = append(categories, *currentCategory)
				}
			}

			// Check if this is a valid change type
			changeType, valid := core.ParseChangeType(matches[1])
			if valid {
				currentCategory = &core.ChangeCategory{
					Type: changeType,
				}
				currentEntries = nil
			} else {
				currentCategory = nil
			}
		} else if currentCategory != nil {
			currentEntries = append(currentEntries, line)
		}
	}

	// Save last category
	if currentCategory != nil {
		currentCategory.Entries = p.parseEntries(currentEntries)
		if len(currentCategory.Entries) > 0 {
			categories = append(categories, *currentCategory)
		}
	}

	return categories
}

// parseEntries parses list items from lines
func (p *ChangelogParser) parseEntries(lines []string) []core.ChangeEntry {
	var entries []core.ChangeEntry
	var currentEntry *core.ChangeEntry
	var currentLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if this is a new list item
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			// Save previous entry
			if currentEntry != nil {
				currentEntry.RawMD = strings.TrimSpace(strings.Join(currentLines, "\n"))
				var buf bytes.Buffer
				p.md.Convert([]byte(currentEntry.RawMD), &buf)
				currentEntry.HTML = buf.String()
				currentEntry.Description = extractDescription(currentEntry.RawMD)
				entries = append(entries, *currentEntry)
			}

			// Start new entry
			content := strings.TrimPrefix(strings.TrimPrefix(trimmed, "- "), "* ")
			currentEntry = &core.ChangeEntry{}
			currentLines = []string{content}
		} else if currentEntry != nil && trimmed != "" {
			// Continuation of current entry
			currentLines = append(currentLines, line)
		}
	}

	// Save last entry
	if currentEntry != nil {
		currentEntry.RawMD = strings.TrimSpace(strings.Join(currentLines, "\n"))
		var buf bytes.Buffer
		p.md.Convert([]byte(currentEntry.RawMD), &buf)
		currentEntry.HTML = buf.String()
		currentEntry.Description = extractDescription(currentEntry.RawMD)
		entries = append(entries, *currentEntry)
	}

	return entries
}

// extractDescription extracts a plain text description from markdown
func extractDescription(md string) string {
	// Remove markdown formatting for plain text description
	// This is a simple implementation - just return first line
	lines := strings.Split(md, "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return md
}

// parseSemver extracts major, minor, patch from a version string
func parseSemver(version string) (major, minor, patch int) {
	// Remove leading 'v' if present
	version = strings.TrimPrefix(version, "v")

	// Split by dots
	parts := strings.Split(version, ".")

	if len(parts) >= 1 {
		major, _ = strconv.Atoi(parts[0])
	}
	if len(parts) >= 2 {
		minor, _ = strconv.Atoi(parts[1])
	}
	if len(parts) >= 3 {
		// Handle prerelease suffixes like "0-beta.1"
		patchStr := parts[2]
		if idx := strings.IndexAny(patchStr, "-+"); idx != -1 {
			patchStr = patchStr[:idx]
		}
		patch, _ = strconv.Atoi(patchStr)
	}

	return major, minor, patch
}
