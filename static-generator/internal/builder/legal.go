package builder

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/studiowebux/minimaldoc/static-generator/internal/core"
	"github.com/studiowebux/minimaldoc/static-generator/internal/parser"
)

// LegalBuilder handles building legal pages
type LegalBuilder struct {
	frontmatterParser *parser.FrontmatterParser
	markdownParser    *parser.MarkdownParser
}

// NewLegalBuilder creates a new legal builder
func NewLegalBuilder() *LegalBuilder {
	return &LegalBuilder{
		frontmatterParser: parser.NewFrontmatterParser(),
		markdownParser:    parser.NewMarkdownParser(),
	}
}

// Build creates legal pages from markdown files in the legal source directory
func (lb *LegalBuilder) Build(docsRoot string, config core.LegalConfig, basePath string) ([]*core.LegalPage, error) {
	if !config.Enabled {
		return nil, nil
	}

	legalDir := filepath.Join(docsRoot, core.LegalSourceDir)

	// Parse all legal page files
	pages, err := lb.parsePages(legalDir, basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse legal pages: %w", err)
	}

	// Sort pages by order
	sort.Slice(pages, func(i, j int) bool {
		if pages[i].Order != pages[j].Order {
			return pages[i].Order < pages[j].Order
		}
		return pages[i].Title < pages[j].Title
	})

	return pages, nil
}

// parsePages parses all markdown files in the legal directory
func (lb *LegalBuilder) parsePages(legalDir string, basePath string) ([]*core.LegalPage, error) {
	var pages []*core.LegalPage

	// Glob for markdown files
	pattern := filepath.Join(legalDir, "*.md")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob legal files: %w", err)
	}

	for _, filePath := range files {
		page, err := lb.parsePage(filePath, basePath)
		if err != nil {
			fmt.Printf("Warning: failed to parse legal page %s: %v\n", filePath, err)
			continue
		}
		pages = append(pages, page)
	}

	return pages, nil
}

// parsePage parses a single legal page markdown file
func (lb *LegalBuilder) parsePage(filePath string, basePath string) (*core.LegalPage, error) {
	// Parse frontmatter
	meta, content, err := lb.frontmatterParser.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("frontmatter parse error: %w", err)
	}

	// Generate slug from filename
	base := filepath.Base(filePath)
	slug := base[:len(base)-len(filepath.Ext(base))]

	// Parse markdown to HTML
	html, err := lb.markdownParser.ParseWithContext(content, "", basePath)
	if err != nil {
		return nil, fmt.Errorf("markdown parse error: %w", err)
	}

	page := &core.LegalPage{
		Slug:        slug,
		FilePath:    filePath,
		Title:       meta.Title,
		Description: meta.Description,
		Order:       meta.MenuOrder,
		RawMD:       string(content),
		HTML:        string(html),
	}

	// Default title to capitalized slug if not set
	if page.Title == "" {
		page.Title = capitalizeFirst(slug)
	}

	return page, nil
}
