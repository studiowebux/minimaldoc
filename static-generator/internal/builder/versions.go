package builder

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/studiowebux/minimaldoc/static-generator/internal/core"
	"github.com/studiowebux/minimaldoc/static-generator/internal/parser"
)

// VersionBuilder handles building version-specific page sets
type VersionBuilder struct {
	frontmatterParser *parser.FrontmatterParser
	markdownParser    *parser.MarkdownParser
	tocParser         *parser.TOCParser
}

// NewVersionBuilder creates a new version builder
func NewVersionBuilder() *VersionBuilder {
	return &VersionBuilder{
		frontmatterParser: parser.NewFrontmatterParser(),
		markdownParser:    parser.NewMarkdownParser(),
		tocParser:         parser.NewTOCParser(),
	}
}

// Build builds version-specific page sets from shared pages and version overrides
func (vb *VersionBuilder) Build(site *core.Site) error {
	if !site.Config.Versions.Enabled || len(site.Config.Versions.List) == 0 {
		return nil
	}

	// Initialize the versioned pages map
	site.VersionedPages = make(map[string][]*Page)

	// For each configured version, build its page set
	for _, version := range site.Config.Versions.List {
		pages, err := vb.buildVersionPages(site, version)
		if err != nil {
			return err
		}
		site.VersionedPages[version.Name] = pages
	}

	return nil
}

// buildVersionPages builds the page set for a specific version
func (vb *VersionBuilder) buildVersionPages(site *core.Site, version core.VersionInfo) ([]*core.Page, error) {
	var pages []*core.Page

	// 1. Start with shared pages that match this version
	for _, page := range site.Pages {
		if vb.pageMatchesVersion(page, version.Name) {
			// Clone the page for this version
			versionPage := vb.clonePage(page, version)
			pages = append(pages, versionPage)
		}
	}

	// 2. Look for version-specific overrides in __versions__/{version}/
	versionDir := filepath.Join(site.DocsRoot, core.VersionSourceDir, version.Name)
	if _, err := os.Stat(versionDir); err == nil {
		overrides, err := vb.discoverVersionOverrides(versionDir, site.DocsRoot, version)
		if err != nil {
			return nil, err
		}

		// Merge overrides with existing pages (override wins)
		pages = vb.mergeOverrides(pages, overrides)
	}

	return pages, nil
}

// pageMatchesVersion checks if a page should appear in a given version
func (vb *VersionBuilder) pageMatchesVersion(page *core.Page, versionName string) bool {
	// If no versions specified, page appears in all versions
	if len(page.Metadata.Versions) == 0 {
		return true
	}

	// Check if version is in the list
	for _, v := range page.Metadata.Versions {
		if v == versionName || v == "all" {
			return true
		}
	}

	// Check removed_in - if page is removed in this or earlier version, exclude
	if page.Metadata.RemovedIn != "" {
		if vb.versionCompare(versionName, page.Metadata.RemovedIn) >= 0 {
			return false
		}
	}

	return false
}

// versionCompare compares two version strings
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func (vb *VersionBuilder) versionCompare(v1, v2 string) int {
	// Simple string comparison for now (works for v1, v2, v3, etc.)
	// For semantic versioning, this would need more sophisticated parsing
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	if v1 < v2 {
		return -1
	}
	if v1 > v2 {
		return 1
	}
	return 0
}

// clonePage creates a copy of a page for a specific version
func (vb *VersionBuilder) clonePage(page *core.Page, version core.VersionInfo) *core.Page {
	clone := &core.Page{
		SourcePath:      page.SourcePath,
		RelPath:         page.RelPath,
		Slug:            page.Slug,
		OutputPath:      page.OutputPath,
		ModTime:         page.ModTime,
		Metadata:        cloneMetadata(page.Metadata),
		RawMD:           page.RawMD,
		HTML:            page.HTML,
		Order:           page.Order,
		TOC:             page.TOC,
		IsStale:         page.IsStale,
		StaleAge:        page.StaleAge,
		DaysSinceUpdate: page.DaysSinceUpdate,
	}

	// Update output path for versioned output
	// For non-default versions, prepend version path
	// e.g., "getting-started.html" -> "v1/getting-started.html"
	if version.Path != "" {
		clone.OutputPath = filepath.Join(version.Path, page.OutputPath)
	}

	return clone
}

// discoverVersionOverrides finds version-specific override pages
func (vb *VersionBuilder) discoverVersionOverrides(versionDir, docsRoot string, version core.VersionInfo) ([]*core.Page, error) {
	var pages []*core.Page

	err := filepath.WalkDir(versionDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		// Create page with relative path from version directory
		relPath, _ := filepath.Rel(versionDir, path)
		page := core.NewPage(path, versionDir)
		page.RelPath = relPath

		// Parse frontmatter
		meta, content, err := vb.frontmatterParser.ParseFile(path)
		if err != nil {
			return err
		}
		page.Metadata = meta
		page.RawMD = content

		// Override order if specified
		if meta.MenuOrder >= 0 {
			page.Order = meta.MenuOrder
		} else {
			page.Order = extractOrder(filepath.Base(path))
		}

		// Parse markdown with link transformation
		html, err := vb.markdownParser.ParseWithContext(content, relPath, "")
		if err != nil {
			return err
		}
		page.HTML = html

		// Generate TOC
		toc, err := vb.tocParser.Parse(content)
		if err != nil {
			return err
		}
		page.TOC = toc

		// Set output path for versioned output
		if version.Path != "" {
			page.OutputPath = filepath.Join(version.Path, page.Slug+".html")
		} else {
			page.OutputPath = page.Slug + ".html"
		}

		pages = append(pages, page)
		return nil
	})

	return pages, err
}

// mergeOverrides merges version-specific overrides with shared pages
func (vb *VersionBuilder) mergeOverrides(shared, overrides []*core.Page) []*core.Page {
	// Create a map of shared pages by slug for quick lookup
	pageMap := make(map[string]*core.Page)
	for _, page := range shared {
		pageMap[page.Slug] = page
	}

	// Override or add pages from overrides
	for _, override := range overrides {
		pageMap[override.Slug] = override
	}

	// Convert map back to slice
	result := make([]*core.Page, 0, len(pageMap))
	for _, page := range pageMap {
		result = append(result, page)
	}

	return result
}

// Page is an alias to avoid import cycle issues
type Page = core.Page
