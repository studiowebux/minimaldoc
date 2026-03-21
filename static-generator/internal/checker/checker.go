package checker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/studiowebux/minimaldoc/static-generator/internal/core"
)

// LinkChecker orchestrates link collection and validation
type LinkChecker struct {
	config    core.LinkCheckConfig
	site      *core.Site
	collector *LinkCollector
	validator *LinkValidator
	reporter  *Reporter
}

// NewLinkChecker creates a new link checker
func NewLinkChecker(site *core.Site, config core.LinkCheckConfig) *LinkChecker {
	basePath := extractBasePath(site.Config.BaseURL)

	return &LinkChecker{
		config:    config,
		site:      site,
		collector: NewLinkCollector(site.DocsRoot),
		validator: NewLinkValidator(
			config,
			site.OutputRoot,
			site.DocsRoot,
			basePath,
			site.Config.CleanURLs,
		),
		reporter: NewReporter(config.Mode),
	}
}

// Check runs the link checker and returns an error if links are broken and mode is "error"
func (c *LinkChecker) Check() error {
	if !c.config.Enabled || c.config.Mode == core.LinkCheckIgnore {
		return nil
	}

	fmt.Println("Checking links...")

	// Collect links from all markdown files
	if err := c.collectLinks(); err != nil {
		return fmt.Errorf("failed to collect links: %w", err)
	}

	links := c.collector.Links()
	if len(links) == 0 {
		fmt.Println("No links found to check")
		return nil
	}

	// Validate collected links
	result := c.validator.Validate(links)

	// Report results
	c.reporter.Report(result)

	// Return error if mode is "error" and there are broken links
	if c.config.Mode == core.LinkCheckError && result.HasErrors() {
		return fmt.Errorf("link check failed: %d broken links", result.BrokenCount())
	}

	return nil
}

// collectLinks walks the docs directory and collects links
func (c *LinkChecker) collectLinks() error {
	// Collect from all pages
	for _, page := range c.site.Pages {
		if err := c.collector.CollectFromFile(page.SourcePath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to collect links from %s: %v\n", page.SourcePath, err)
		}
	}

	// Also collect from special directories
	specialDirs := []string{
		core.KBSourceDir,
		core.FaqSourceDir,
		core.LegalSourceDir,
		core.ChangelogSourceDir,
		core.LandingSourceDir,
	}

	for _, dir := range specialDirs {
		dirPath := filepath.Join(c.site.DocsRoot, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			continue
		}

		_ = filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".md" {
				return nil
			}
			_ = c.collector.CollectFromFile(path) // best-effort; walk continues on per-file errors
			return nil
		})
	}

	return nil
}

// extractBasePath extracts the path component from a URL
func extractBasePath(baseURL string) string {
	if baseURL == "" {
		return ""
	}

	// Remove protocol
	url := baseURL
	for _, prefix := range []string{"https://", "http://"} {
		if len(url) > len(prefix) && url[:len(prefix)] == prefix {
			url = url[len(prefix):]
			break
		}
	}

	// Find first slash after domain
	slashIdx := -1
	for i, c := range url {
		if c == '/' {
			slashIdx = i
			break
		}
	}

	if slashIdx == -1 {
		return ""
	}

	path := url[slashIdx:]
	// Remove trailing slash
	for len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	if path == "/" {
		return ""
	}

	return path
}
