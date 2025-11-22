package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
	"github.com/studiowebux/minimaldoc/internal/parser"
)

// Builder orchestrates the site generation process
type Builder struct {
	site *core.Site

	// Parsers
	frontmatterParser *parser.FrontmatterParser
	markdownParser    *parser.MarkdownParser
	tocParser         *parser.TOCParser

	// Builders
	navBuilder *NavigationBuilder
}

// NewBuilder creates a new site builder
func NewBuilder(site *core.Site) *Builder {
	return &Builder{
		site:              site,
		frontmatterParser: parser.NewFrontmatterParser(),
		markdownParser:    parser.NewMarkdownParser(),
		tocParser:         parser.NewTOCParser(),
		navBuilder:        NewNavigationBuilder(),
	}
}

// Build orchestrates the entire build process
func (b *Builder) Build() error {
	fmt.Println("Building site...")

	// 1. Discover and parse all markdown files
	if err := b.discoverPages(); err != nil {
		return fmt.Errorf("failed to discover pages: %w", err)
	}

	// 2. Parse each page
	for _, page := range b.site.Pages {
		if err := b.parsePage(page); err != nil {
			return fmt.Errorf("failed to parse page %s: %w", page.SourcePath, err)
		}
	}

	// 3. Build navigation
	b.site.Navigation = b.navBuilder.Build(b.site.Pages, b.site.DocsRoot, b.site.Config.NavDepth)

	// 4. Compute prev/next links
	b.computePrevNext()

	fmt.Printf("Discovered %d pages\n", len(b.site.Pages))
	return nil
}

// discoverPages walks the docs directory and discovers all markdown files
func (b *Builder) discoverPages() error {
	return filepath.WalkDir(b.site.DocsRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		// Skip TOC.md (navigation definition file)
		if d.Name() == "TOC.md" {
			return nil
		}

		// Create page
		page := core.NewPage(path, b.site.DocsRoot)

		// Extract order from filename
		page.Order = extractOrder(filepath.Base(path))

		b.site.Pages = append(b.site.Pages, page)

		return nil
	})
}

// parsePage parses a single page's markdown content
func (b *Builder) parsePage(page *core.Page) error {
	// 1. Parse frontmatter
	meta, content, err := b.frontmatterParser.ParseFile(page.SourcePath)
	if err != nil {
		return fmt.Errorf("frontmatter parse error: %w", err)
	}

	page.Metadata = meta
	page.RawMD = content

	// Override order if specified in metadata (-1 is the default, so >= 0 means explicit)
	if meta.MenuOrder >= 0 {
		page.Order = meta.MenuOrder
	}

	// 2. Parse markdown to HTML
	html, err := b.markdownParser.Parse(content)
	if err != nil {
		return fmt.Errorf("markdown parse error: %w", err)
	}

	page.HTML = html

	// 3. Generate table of contents
	toc, err := b.tocParser.Parse(content)
	if err != nil {
		return fmt.Errorf("TOC parse error: %w", err)
	}

	page.TOC = toc

	// 4. Determine output path
	page.OutputPath = b.getOutputPath(page)

	return nil
}

// getOutputPath determines the output file path for a page
func (b *Builder) getOutputPath(page *core.Page) string {
	slug := page.Slug

	if b.site.Config.CleanURLs {
		// Clean URLs: /docs/guide/ instead of /docs/guide.html
		// Special case: root index page should stay as index.html
		if slug == "index" {
			return filepath.Join(b.site.OutputRoot, "index.html")
		}

		// For other pages, trim trailing /index and create directory structure
		if strings.HasSuffix(slug, "/index") {
			slug = strings.TrimSuffix(slug, "/index")
		}
		return filepath.Join(b.site.OutputRoot, slug, "index.html")
	}

	// Standard URLs: /docs/guide.html
	return filepath.Join(b.site.OutputRoot, slug+".html")
}

// computePrevNext computes the previous and next page links
func (b *Builder) computePrevNext() {
	// Get ordered list of pages from navigation
	orderedPages := b.flattenNavigation(b.site.Navigation.Items)

	// Set prev/next links
	for i, page := range orderedPages {
		if i > 0 {
			page.Prev = orderedPages[i-1]
		}
		if i < len(orderedPages)-1 {
			page.Next = orderedPages[i+1]
		}
	}
}

// flattenNavigation extracts pages in order from navigation tree
func (b *Builder) flattenNavigation(items []*core.NavItem) []*core.Page {
	var pages []*core.Page

	for _, item := range items {
		// Add page if it's a leaf node
		if item.Page != nil {
			pages = append(pages, item.Page)
		}

		// Recursively process children
		if len(item.Children) > 0 {
			pages = append(pages, b.flattenNavigation(item.Children)...)
		}
	}

	return pages
}
