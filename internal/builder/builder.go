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
	openapiParser     *parser.OpenAPIParser

	// Builders
	navBuilder *NavigationBuilder
}

// NewBuilder creates a new site builder
func NewBuilder(site *core.Site) *Builder {
	cacheDir := filepath.Join(site.OutputRoot, site.Config.OpenAPI.CacheDir)
	return &Builder{
		site:              site,
		frontmatterParser: parser.NewFrontmatterParser(),
		markdownParser:    parser.NewMarkdownParser(),
		tocParser:         parser.NewTOCParser(),
		openapiParser:     parser.NewOpenAPIParser(cacheDir),
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

	// 3. Discover and parse OpenAPI specs (if enabled)
	if b.site.Config.OpenAPI.Enabled {
		if err := b.discoverAndParseOpenAPISpecs(); err != nil {
			return fmt.Errorf("failed to parse OpenAPI specs: %w", err)
		}
	}

	// 4. Build navigation
	b.site.Navigation = b.navBuilder.Build(b.site.Pages, b.site.DocsRoot, b.site.Config.NavDepth)

	// 5. Compute prev/next links
	b.computePrevNext()

	fmt.Printf("Discovered %d pages\n", len(b.site.Pages))
	if b.site.Config.OpenAPI.Enabled {
		fmt.Printf("Discovered %d OpenAPI specs\n", len(b.site.APISpecs))
	}
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

	// 2. Parse markdown to HTML with link transformation
	html, err := b.markdownParser.ParseWithContext(content, page.RelPath)
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

// discoverAndParseOpenAPISpecs discovers and parses all OpenAPI specifications
func (b *Builder) discoverAndParseOpenAPISpecs() error {
	// 1. Discover local OpenAPI files
	localFiles, err := b.discoverOpenAPIFiles()
	if err != nil {
		return fmt.Errorf("failed to discover OpenAPI files: %w", err)
	}

	// 2. Parse local files
	for _, filePath := range localFiles {
		spec, err := b.openapiParser.ParseFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse OpenAPI spec %s: %v\n", filePath, err)
			continue
		}
		b.site.APISpecs = append(b.site.APISpecs, spec)
	}

	// 3. Parse remote URLs
	if len(b.site.Config.OpenAPI.SpecURLs) > 0 {
		cacheDir := filepath.Join(b.site.OutputRoot, b.site.Config.OpenAPI.CacheDir)

		for _, url := range b.site.Config.OpenAPI.SpecURLs {
			// Generate cache file name from URL
			cacheName := b.openapiParser.NameFromURL(url) + ".json"
			cachePath := filepath.Join(cacheDir, cacheName)

			// Check if cache exists
			cacheExists := false
			if _, err := os.Stat(cachePath); err == nil {
				cacheExists = true
			}

			// Decide whether to fetch from URL or use cache
			var spec *core.APISpec
			var err error

			if b.site.Config.OpenAPI.SyncOnBuild || !cacheExists {
				// Fetch from URL (either sync is enabled or no cache exists)
				spec, err = b.openapiParser.ParseURL(url)
				if err != nil {
					// If fetch fails and cache exists, try loading from cache
					if cacheExists {
						fmt.Fprintf(os.Stderr, "Warning: failed to fetch OpenAPI spec from %s, using cache: %v\n", url, err)
						spec, err = b.openapiParser.ParseFile(cachePath)
						if err != nil {
							fmt.Fprintf(os.Stderr, "Warning: failed to load cached spec: %v\n", err)
							continue
						}
					} else {
						fmt.Fprintf(os.Stderr, "Warning: failed to fetch OpenAPI spec from %s: %v\n", url, err)
						continue
					}
				}
			} else {
				// Use cached spec
				spec, err = b.openapiParser.ParseFile(cachePath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to load cached spec from %s, trying to fetch: %v\n", cachePath, err)
					spec, err = b.openapiParser.ParseURL(url)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Warning: failed to fetch OpenAPI spec from %s: %v\n", url, err)
						continue
					}
				}
			}

			b.site.APISpecs = append(b.site.APISpecs, spec)
		}
	}

	return nil
}

// discoverOpenAPIFiles walks the docs directory and discovers OpenAPI spec files
func (b *Builder) discoverOpenAPIFiles() ([]string, error) {
	var specFiles []string

	// 1. Check explicitly configured spec files
	for _, pattern := range b.site.Config.OpenAPI.SpecFiles {
		// If it's a glob pattern, expand it
		if strings.Contains(pattern, "*") {
			matches, err := filepath.Glob(filepath.Join(b.site.DocsRoot, pattern))
			if err != nil {
				return nil, fmt.Errorf("invalid glob pattern %s: %w", pattern, err)
			}
			specFiles = append(specFiles, matches...)
		} else {
			// It's a direct file path
			fullPath := filepath.Join(b.site.DocsRoot, pattern)
			if _, err := os.Stat(fullPath); err == nil {
				specFiles = append(specFiles, fullPath)
			}
		}
	}

	// 2. Auto-discover OpenAPI files in docs directory (if spec files list is empty)
	if len(b.site.Config.OpenAPI.SpecFiles) == 0 {
		err := filepath.WalkDir(b.site.DocsRoot, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Skip directories
			if d.IsDir() {
				return nil
			}

			// Check if file looks like an OpenAPI spec
			if b.isOpenAPIFile(path) {
				specFiles = append(specFiles, path)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return specFiles, nil
}

// isOpenAPIFile checks if a file is likely an OpenAPI specification
func (b *Builder) isOpenAPIFile(path string) bool {
	ext := filepath.Ext(path)

	// Check extension
	if ext != ".yaml" && ext != ".yml" && ext != ".json" {
		return false
	}

	// Check filename patterns
	base := strings.ToLower(filepath.Base(path))
	patterns := []string{
		"openapi",
		"swagger",
		"api-spec",
		"api_spec",
	}

	for _, pattern := range patterns {
		if strings.Contains(base, pattern) {
			return true
		}
	}

	return false
}
