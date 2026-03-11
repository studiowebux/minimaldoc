package core

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Page represents a single documentation page
type Page struct {
	// File information
	SourcePath string    // Original markdown file path
	RelPath    string    // Relative path from docs root
	Slug       string    // URL-friendly slug (e.g., "getting-started")
	OutputPath string    // Output HTML file path
	ModTime    time.Time // File modification time

	// Content
	Metadata Metadata // Parsed frontmatter
	RawMD    []byte   // Raw markdown content (without frontmatter)
	HTML     []byte   // Rendered HTML content

	// Navigation
	Order    int     // Computed order (from filename prefix or metadata)
	Parent   *Page   // Parent page (for nested navigation)
	Children []*Page // Child pages
	Next     *Page   // Next page in sequence
	Prev     *Page   // Previous page in sequence

	// Table of Contents
	TOC *TOC // Parsed table of contents from headings

	// Stale content warning (computed at build time)
	IsStale         bool   // True if content is older than threshold
	StaleAge        string // Human-readable age: "2 years", "18 months"
	DaysSinceUpdate int    // Exact days since last update
}

// TOC represents the table of contents for a page
type TOC struct {
	Items []*TOCItem
}

// TOCItem represents a single heading in the table of contents
type TOCItem struct {
	ID       string     // Anchor ID (e.g., "getting-started")
	Title    string     // Heading text
	Level    int        // Heading level (1-6)
	Children []*TOCItem // Nested headings
}

// NewPage creates a new Page from a file path
func NewPage(sourcePath string, docsRoot string) *Page {
	relPath, _ := filepath.Rel(docsRoot, sourcePath)

	// Get file modification time
	var modTime time.Time
	if info, err := os.Stat(sourcePath); err == nil {
		modTime = info.ModTime()
	}

	return &Page{
		SourcePath: sourcePath,
		RelPath:    relPath,
		Slug:       GenerateSlugFromPath(relPath),
		Metadata:   DefaultMetadata(),
		Children:   []*Page{},
		ModTime:    modTime,
	}
}

// GenerateSlugFromPath creates a URL-friendly slug from a file path
// Examples:
//   - "01-getting-started.md" -> "getting-started"
//   - "docs/02-api/index.md" -> "docs/api/index"
func GenerateSlugFromPath(relPath string) string {
	// Remove extension
	slug := strings.TrimSuffix(relPath, filepath.Ext(relPath))

	// Remove numbered prefixes from each segment
	parts := strings.Split(slug, string(filepath.Separator))
	for i, part := range parts {
		// Remove leading numbers and separators (e.g., "01-", "02_")
		part = strings.TrimLeft(part, "0123456789-_")
		parts[i] = part
	}

	slug = strings.Join(parts, "/")

	// Convert to lowercase and replace spaces
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")

	return slug
}

// Title returns the display title for the page
func (p *Page) Title() string {
	// Priority: MenuTitle > Metadata.Title > generated from filename
	if p.Metadata.MenuTitle != "" {
		return p.Metadata.MenuTitle
	}
	if p.Metadata.Title != "" {
		return p.Metadata.Title
	}
	// Fallback: use filename (cleaned up) as title
	return titleFromFilename(p.SourcePath)
}

// titleFromFilename generates a title from a filename
func titleFromFilename(sourcePath string) string {
	// Get base filename without extension
	base := filepath.Base(sourcePath)
	base = strings.TrimSuffix(base, filepath.Ext(base))

	// Remove number prefix
	base = strings.TrimLeft(base, "0123456789-_")

	// Replace separators with spaces
	base = strings.ReplaceAll(base, "-", " ")
	base = strings.ReplaceAll(base, "_", " ")

	// Title case
	words := strings.Fields(base)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}

	return strings.Join(words, " ")
}

// IsHidden returns whether the page should be hidden from navigation
func (p *Page) IsHidden() bool {
	return p.Metadata.Hidden
}
