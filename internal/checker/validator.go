package checker

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// LinkValidator validates collected links
type LinkValidator struct {
	config     core.LinkCheckConfig
	outputRoot string
	docsRoot   string
	basePath   string
	cleanURLs  bool

	// Cache for heading IDs in HTML files
	headingCache map[string]map[string]bool

	// HTTP client for external links
	httpClient *http.Client
}

// NewLinkValidator creates a new link validator
func NewLinkValidator(config core.LinkCheckConfig, outputRoot, docsRoot, basePath string, cleanURLs bool) *LinkValidator {
	return &LinkValidator{
		config:       config,
		outputRoot:   outputRoot,
		docsRoot:     docsRoot,
		basePath:     basePath,
		cleanURLs:    cleanURLs,
		headingCache: make(map[string]map[string]bool),
		httpClient: &http.Client{
			Timeout: time.Duration(config.ExternalTimeout) * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 3 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

// Validate checks all collected links and returns results
func (v *LinkValidator) Validate(links []core.CollectedLink) *core.LinkCheckResult {
	result := &core.LinkCheckResult{
		TotalLinks:    len(links),
		BrokenLinks:   []core.BrokenLink{},
		SkippedLinks:  0,
		ExternalLinks: 0,
	}

	for _, link := range links {
		// Check if link should be ignored
		if v.shouldIgnore(link.URL) {
			result.SkippedLinks++
			continue
		}

		// Validate based on link type
		switch link.LinkType {
		case core.LinkTypeInternalPage:
			if broken := v.validateInternalPage(link); broken != nil {
				result.BrokenLinks = append(result.BrokenLinks, *broken)
			}

		case core.LinkTypeInternalAnchor:
			if broken := v.validateAnchor(link); broken != nil {
				result.BrokenLinks = append(result.BrokenLinks, *broken)
			}

		case core.LinkTypeInternalAsset:
			if broken := v.validateAsset(link); broken != nil {
				result.BrokenLinks = append(result.BrokenLinks, *broken)
			}

		case core.LinkTypeExternal:
			result.ExternalLinks++
			if v.config.CheckExternal {
				if broken := v.validateExternal(link); broken != nil {
					result.BrokenLinks = append(result.BrokenLinks, *broken)
				}
			}

		case core.LinkTypeEmail, core.LinkTypeOther:
			// Skip validation for these types
			result.SkippedLinks++
		}
	}

	return result
}

// shouldIgnore checks if a link matches ignore patterns or allowed broken list
func (v *LinkValidator) shouldIgnore(url string) bool {
	// Check allowed broken links
	for _, allowed := range v.config.AllowedBroken {
		if url == allowed {
			return true
		}
	}

	// Check ignore patterns (simple glob matching)
	for _, pattern := range v.config.IgnorePatterns {
		if matchesGlob(url, pattern) {
			return true
		}
	}

	return false
}

// matchesGlob performs simple glob matching with * wildcard
func matchesGlob(s, pattern string) bool {
	// Convert glob pattern to regex
	regexPattern := "^" + regexp.QuoteMeta(pattern) + "$"
	regexPattern = strings.ReplaceAll(regexPattern, `\*`, ".*")

	matched, _ := regexp.MatchString(regexPattern, s)
	return matched
}

// validateInternalPage validates an internal page link
func (v *LinkValidator) validateInternalPage(link core.CollectedLink) *core.BrokenLink {
	url := link.URL

	// Split URL and fragment
	parts := strings.SplitN(url, "#", 2)
	path := parts[0]
	var fragment string
	if len(parts) > 1 {
		fragment = parts[1]
	}

	// Handle relative .md links
	if strings.HasSuffix(path, ".md") {
		path = v.resolveMarkdownLink(path, link.SourceFile)
	}

	// Handle absolute links (starting with /)
	if strings.HasPrefix(path, "/") {
		// Remove base path if present
		if v.basePath != "" && strings.HasPrefix(path, v.basePath) {
			path = strings.TrimPrefix(path, v.basePath)
		}
		path = strings.TrimPrefix(path, "/")
	}

	// Resolve to output path
	outputPath := v.resolveOutputPath(path)

	// Check if file exists
	fullPath := filepath.Join(v.outputRoot, outputPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		suggestion := v.findSimilarFile(outputPath)
		return &core.BrokenLink{
			Link:       link,
			Reason:     "file not found",
			Suggestion: suggestion,
		}
	}

	// Validate fragment if present
	if fragment != "" {
		if !v.validateFragment(fullPath, fragment) {
			suggestion := v.findSimilarAnchor(fullPath, fragment)
			return &core.BrokenLink{
				Link:       link,
				Reason:     fmt.Sprintf("anchor #%s not found", fragment),
				Suggestion: suggestion,
			}
		}
	}

	return nil
}

// validateAnchor validates an anchor-only link (#section)
func (v *LinkValidator) validateAnchor(link core.CollectedLink) *core.BrokenLink {
	fragment := strings.TrimPrefix(link.URL, "#")

	// Find the output file for the source file
	sourcePath := link.SourceFile
	slug := core.GenerateSlugFromPath(sourcePath)
	outputPath := v.resolveOutputPath(slug + ".html")
	fullPath := filepath.Join(v.outputRoot, outputPath)

	if !v.validateFragment(fullPath, fragment) {
		suggestion := v.findSimilarAnchor(fullPath, fragment)
		return &core.BrokenLink{
			Link:       link,
			Reason:     fmt.Sprintf("anchor #%s not found in current page", fragment),
			Suggestion: suggestion,
		}
	}

	return nil
}

// validateAsset validates an asset link (image, file, etc.)
func (v *LinkValidator) validateAsset(link core.CollectedLink) *core.BrokenLink {
	url := link.URL

	// Remove query string
	path := strings.Split(url, "?")[0]

	// Handle absolute paths
	if strings.HasPrefix(path, "/") {
		if v.basePath != "" && strings.HasPrefix(path, v.basePath) {
			path = strings.TrimPrefix(path, v.basePath)
		}
		path = strings.TrimPrefix(path, "/")
	} else {
		// Relative path - resolve from source file location
		sourceDir := filepath.Dir(link.SourceFile)
		path = filepath.Join(sourceDir, path)
		path = filepath.Clean(path)
	}

	// Check in output directory
	fullPath := filepath.Join(v.outputRoot, path)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// Also check in docs directory (might be copied during build)
		docsPath := filepath.Join(v.docsRoot, path)
		if _, err := os.Stat(docsPath); os.IsNotExist(err) {
			return &core.BrokenLink{
				Link:   link,
				Reason: "asset not found",
			}
		}
	}

	return nil
}

// validateExternal validates an external URL
func (v *LinkValidator) validateExternal(link core.CollectedLink) *core.BrokenLink {
	req, err := http.NewRequest("HEAD", link.URL, nil)
	if err != nil {
		return &core.BrokenLink{
			Link:   link,
			Reason: fmt.Sprintf("invalid URL: %v", err),
		}
	}

	req.Header.Set("User-Agent", "MinimalDoc Link Checker/1.0")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return &core.BrokenLink{
			Link:   link,
			Reason: fmt.Sprintf("request failed: %v", err),
		}
	}
	defer resp.Body.Close()

	// Consider 2xx and 3xx as success
	if resp.StatusCode >= 400 {
		return &core.BrokenLink{
			Link:   link,
			Reason: fmt.Sprintf("HTTP %d", resp.StatusCode),
		}
	}

	return nil
}

// resolveMarkdownLink resolves a relative .md link to output path
func (v *LinkValidator) resolveMarkdownLink(mdPath, sourcePath string) string {
	// Get directory of source file
	sourceDir := filepath.Dir(sourcePath)

	// Join and clean
	resolved := filepath.Join(sourceDir, mdPath)
	resolved = filepath.Clean(resolved)

	// Convert to slug
	slug := core.GenerateSlugFromPath(resolved)

	return "/" + slug + ".html"
}

// resolveOutputPath converts a URL path to the actual output file path
func (v *LinkValidator) resolveOutputPath(urlPath string) string {
	// Remove leading slash
	path := strings.TrimPrefix(urlPath, "/")

	// Handle clean URLs
	if v.cleanURLs {
		if !strings.HasSuffix(path, ".html") && !strings.Contains(filepath.Base(path), ".") {
			path = filepath.Join(path, "index.html")
		}
	}

	// Handle directory requests
	if strings.HasSuffix(path, "/") {
		path = filepath.Join(path, "index.html")
	}

	return path
}

// validateFragment checks if a fragment/anchor exists in an HTML file
func (v *LinkValidator) validateFragment(htmlPath, fragment string) bool {
	// Check cache first
	if headings, ok := v.headingCache[htmlPath]; ok {
		return headings[fragment]
	}

	// Parse HTML file for id attributes
	headings := v.extractHeadingIDs(htmlPath)
	v.headingCache[htmlPath] = headings

	return headings[fragment]
}

// extractHeadingIDs extracts all id attributes from an HTML file
func (v *LinkValidator) extractHeadingIDs(htmlPath string) map[string]bool {
	ids := make(map[string]bool)

	content, err := os.ReadFile(htmlPath)
	if err != nil {
		return ids
	}

	// Simple regex to find id attributes
	idRegex := regexp.MustCompile(`id=["']([^"']+)["']`)
	matches := idRegex.FindAllSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			ids[string(match[1])] = true
		}
	}

	return ids
}

// findSimilarFile suggests a similar file name
func (v *LinkValidator) findSimilarFile(path string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	dirPath := filepath.Join(v.outputRoot, dir)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return ""
	}

	baseLower := strings.ToLower(base)
	for _, entry := range entries {
		if strings.ToLower(entry.Name()) == baseLower && entry.Name() != base {
			return fmt.Sprintf("did you mean %s?", filepath.Join(dir, entry.Name()))
		}
	}

	return ""
}

// findSimilarAnchor suggests a similar anchor name
func (v *LinkValidator) findSimilarAnchor(htmlPath, fragment string) string {
	headings := v.headingCache[htmlPath]
	if headings == nil {
		headings = v.extractHeadingIDs(htmlPath)
		v.headingCache[htmlPath] = headings
	}

	fragmentLower := strings.ToLower(fragment)
	for id := range headings {
		if strings.ToLower(id) == fragmentLower && id != fragment {
			return fmt.Sprintf("did you mean #%s?", id)
		}
	}

	return ""
}
