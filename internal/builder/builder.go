package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	navBuilder        *NavigationBuilder
	statusBuilder     *StatusBuilder
	changelogBuilder  *ChangelogBuilder
	landingBuilder    *LandingBuilder
	portfolioBuilder  *PortfolioBuilder
	contactBuilder    *ContactBuilder
	faqBuilder        *FaqBuilder
	legalBuilder      *LegalBuilder
	kbBuilder         *KBBuilder
	versionBuilder    *VersionBuilder
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
		statusBuilder:     NewStatusBuilder(),
		changelogBuilder:  NewChangelogBuilder(),
		landingBuilder:    NewLandingBuilder(),
		portfolioBuilder:  NewPortfolioBuilder(),
		contactBuilder:    NewContactBuilder(),
		faqBuilder:        NewFaqBuilder(),
		legalBuilder:      NewLegalBuilder(),
		kbBuilder:         NewKBBuilder(),
		versionBuilder:    NewVersionBuilder(),
	}
}

// Build orchestrates the entire build process
func (b *Builder) Build() error {
	fmt.Println("Building site...")

	// 1. Discover and parse all markdown files
	if err := b.discoverPages(); err != nil {
		return fmt.Errorf("failed to discover pages: %w", err)
	}

	// 2. Handle custom entrypoint (if configured)
	b.handleEntrypoint()

	// 3. Parse each page
	for _, page := range b.site.Pages {
		if err := b.parsePage(page); err != nil {
			return fmt.Errorf("failed to parse page %s: %w", page.SourcePath, err)
		}
	}

	// 4. Discover and parse OpenAPI specs (if enabled)
	if b.site.Config.OpenAPI.Enabled {
		if err := b.discoverAndParseOpenAPISpecs(); err != nil {
			return fmt.Errorf("failed to parse OpenAPI specs: %w", err)
		}
	}

	// 4. Build status page (if enabled)
	if b.site.Config.Status.Enabled {
		statusPage, err := b.statusBuilder.Build(b.site.DocsRoot, b.site.Config.Status)
		if err != nil {
			return fmt.Errorf("failed to build status page: %w", err)
		}
		b.site.StatusPage = statusPage
	}

	// 5. Build changelog (if enabled)
	if b.site.Config.Changelog.Enabled {
		changelogPage, err := b.changelogBuilder.Build(b.site.DocsRoot, b.site.Config.Changelog)
		if err != nil {
			return fmt.Errorf("failed to build changelog: %w", err)
		}
		b.site.ChangelogPage = changelogPage
	}

	// 6. Build landing page (if enabled)
	if b.site.Config.Landing.Enabled {
		basePath := b.getBasePath()
		landingPage, err := b.landingBuilder.Build(b.site.DocsRoot, b.site.Config.Landing, basePath)
		if err != nil {
			return fmt.Errorf("failed to build landing page: %w", err)
		}
		b.site.LandingPage = landingPage
	}

	// 7. Build portfolio page (if enabled)
	if b.site.Config.Portfolio.Enabled {
		basePath := b.getBasePath()
		portfolioPage, err := b.portfolioBuilder.Build(b.site.DocsRoot, b.site.Config.Portfolio, basePath)
		if err != nil {
			return fmt.Errorf("failed to build portfolio page: %w", err)
		}
		b.site.PortfolioPage = portfolioPage
	}

	// 8. Build contact page (if enabled)
	if b.site.Config.Contact.Enabled {
		contactPage, err := b.contactBuilder.Build(b.site.Config.Contact)
		if err != nil {
			return fmt.Errorf("failed to build contact page: %w", err)
		}
		b.site.ContactPage = contactPage
	}

	// 9. Build FAQ page (if enabled)
	if b.site.Config.Faq.Enabled {
		basePath := b.getBasePath()
		faqPage, err := b.faqBuilder.Build(b.site.DocsRoot, b.site.Config.Faq, basePath)
		if err != nil {
			return fmt.Errorf("failed to build FAQ page: %w", err)
		}
		b.site.FaqPage = faqPage
	}

	// 10. Build legal pages (if enabled)
	if b.site.Config.Legal.Enabled {
		basePath := b.getBasePath()
		legalPages, err := b.legalBuilder.Build(b.site.DocsRoot, b.site.Config.Legal, basePath)
		if err != nil {
			return fmt.Errorf("failed to build legal pages: %w", err)
		}
		b.site.LegalPages = legalPages
	}

	// 11. Build knowledge base (if enabled)
	if b.site.Config.KnowledgeBase.Enabled {
		basePath := b.getBasePath()
		kbPage, err := b.kbBuilder.Build(b.site.DocsRoot, b.site.Config.KnowledgeBase, basePath)
		if err != nil {
			return fmt.Errorf("failed to build knowledge base: %w", err)
		}
		b.site.KBPage = kbPage
	}

	// 12. Build versioned page sets (if enabled)
	if b.site.Config.Versions.Enabled {
		if err := b.versionBuilder.Build(b.site); err != nil {
			return fmt.Errorf("failed to build versioned pages: %w", err)
		}
		fmt.Printf("Built %d version(s)\n", len(b.site.Config.Versions.List))
	}

	// 13. Build navigation
	b.site.Navigation = b.navBuilder.Build(b.site.Pages, b.site.DocsRoot, b.site.Config.NavDepth)

	// 14. Compute prev/next links
	b.computePrevNext()

	fmt.Printf("Discovered %d pages\n", len(b.site.Pages))
	if b.site.Config.OpenAPI.Enabled {
		fmt.Printf("Discovered %d OpenAPI specs\n", len(b.site.APISpecs))
	}
	if b.site.Config.Status.Enabled && b.site.StatusPage != nil {
		fmt.Printf("Status page: %d components, %d active incidents\n",
			len(b.site.StatusPage.Components),
			len(b.site.StatusPage.ActiveIncidents))
	}
	if b.site.Config.Changelog.Enabled && b.site.ChangelogPage != nil {
		fmt.Printf("Changelog: %d releases\n", len(b.site.ChangelogPage.Releases))
	}
	return nil
}

// discoverPages walks the docs directory and discovers all markdown files
func (b *Builder) discoverPages() error {
	// Single-file mode: only process the entrypoint file
	if b.site.Config.SingleFileMode && b.site.Config.Entrypoint != "" {
		filePath := filepath.Join(b.site.DocsRoot, b.site.Config.Entrypoint)
		if _, err := os.Stat(filePath); err != nil {
			return fmt.Errorf("entrypoint file not found: %s", filePath)
		}

		page := core.NewPage(filePath, b.site.DocsRoot)
		page.Order = 0
		b.site.Pages = append(b.site.Pages, page)
		return nil
	}

	return filepath.WalkDir(b.site.DocsRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the status directory entirely (it has its own build process)
		if d.IsDir() && d.Name() == core.StatusSourceDir {
			return filepath.SkipDir
		}

		// Skip the changelog directory entirely (it has its own build process)
		if d.IsDir() && d.Name() == core.ChangelogSourceDir {
			return filepath.SkipDir
		}

		// Skip the portfolio directory entirely (it has its own build process)
		if d.IsDir() && d.Name() == core.PortfolioSourceDir {
			return filepath.SkipDir
		}

		// Skip the faq directory entirely (it has its own build process)
		if d.IsDir() && d.Name() == core.FaqSourceDir {
			return filepath.SkipDir
		}

		// Skip the legal directory entirely (it has its own build process)
		if d.IsDir() && d.Name() == core.LegalSourceDir {
			return filepath.SkipDir
		}

		// Skip the landing directory entirely (it has its own build process)
		if d.IsDir() && d.Name() == core.LandingSourceDir {
			return filepath.SkipDir
		}

		// Skip the knowledge base directory entirely (it has its own build process)
		if d.IsDir() && d.Name() == core.KBSourceDir {
			return filepath.SkipDir
		}

		// Skip the versions directory entirely (it has its own build process)
		if d.IsDir() && d.Name() == core.VersionSourceDir {
			return filepath.SkipDir
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
	basePath := b.getBasePath()
	html, err := b.markdownParser.ParseWithContext(content, page.RelPath, basePath)
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

	// 5. Calculate stale warning (if enabled)
	b.calculateStaleWarning(page)

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

// getBasePath extracts the path component from BaseURL for asset linking
// Examples:
//   - "https://example.com/docs/" → "/docs"
//   - "https://example.com/" → ""
//   - "" → ""
func (b *Builder) getBasePath() string {
	baseURL := b.site.Config.BaseURL
	if baseURL == "" {
		return ""
	}

	// Parse the URL to extract the path
	// Remove protocol and domain, keep only the path
	if strings.HasPrefix(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	} else if strings.HasPrefix(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	}

	// Find the first / after the domain
	parts := strings.SplitN(baseURL, "/", 2)
	if len(parts) < 2 {
		return ""
	}

	// Get the path part and ensure it starts with / and doesn't end with /
	path := "/" + parts[1]
	path = strings.TrimSuffix(path, "/")

	// If path is just "/", return empty string
	if path == "/" {
		return ""
	}

	return path
}

// calculateStaleWarning calculates whether a page is stale and sets the relevant fields
func (b *Builder) calculateStaleWarning(page *core.Page) {
	config := b.site.Config.StaleWarning

	// Check if stale warning is enabled at site level
	if !config.Enabled {
		return
	}

	// Check per-page override to disable
	if page.Metadata.StaleWarning != nil && !*page.Metadata.StaleWarning {
		return
	}

	// Determine threshold (per-page override or site default)
	threshold := config.ThresholdDays
	if page.Metadata.StaleThresholdDays != nil && *page.Metadata.StaleThresholdDays > 0 {
		threshold = *page.Metadata.StaleThresholdDays
	}

	// Calculate days since last update
	now := time.Now()
	daysSince := int(now.Sub(page.ModTime).Hours() / 24)
	page.DaysSinceUpdate = daysSince

	// Check if content is stale
	if daysSince >= threshold {
		page.IsStale = true
		page.StaleAge = formatAge(daysSince)
	}
}

// formatAge formats the number of days into a human-readable string
func formatAge(days int) string {
	switch {
	case days < 30:
		if days == 1 {
			return "1 day"
		}
		return fmt.Sprintf("%d days", days)
	case days < 60:
		return "about a month"
	case days < 365:
		months := days / 30
		if months == 1 {
			return "about a month"
		}
		return fmt.Sprintf("%d months", months)
	case days < 730:
		return "over a year"
	case days < 1095:
		return "about 2 years"
	default:
		years := days / 365
		return fmt.Sprintf("over %d years", years)
	}
}

// handleEntrypoint handles custom entrypoint configuration
// When entrypoint is set, the specified file becomes the homepage (index.html)
func (b *Builder) handleEntrypoint() {
	entrypoint := b.site.Config.Entrypoint
	if entrypoint == "" {
		return
	}

	// Normalize the entrypoint path
	entrypoint = strings.TrimPrefix(entrypoint, "./")

	// Find the page matching the entrypoint
	for _, page := range b.site.Pages {
		if page.RelPath == entrypoint {
			// Change the slug to "index" so it outputs as index.html
			page.Slug = "index"
			fmt.Printf("Using %s as homepage\n", entrypoint)
			return
		}
	}

	// Warn if entrypoint file not found
	fmt.Printf("Warning: entrypoint file %s not found\n", entrypoint)
}
