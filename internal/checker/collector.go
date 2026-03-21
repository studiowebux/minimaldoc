package checker

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// LinkCollector collects links from markdown files
type LinkCollector struct {
	docsRoot string
	links    []core.CollectedLink
}

// NewLinkCollector creates a new link collector
func NewLinkCollector(docsRoot string) *LinkCollector {
	return &LinkCollector{
		docsRoot: docsRoot,
		links:    []core.CollectedLink{},
	}
}

// Patterns for link extraction
var (
	// Markdown links: [text](url) or [text](url "title")
	markdownLinkRegex = regexp.MustCompile(`\[([^\]]*)\]\(([^)\s]+)(?:\s+"[^"]*")?\)`)

	// HTML links: <a href="url">
	htmlLinkRegex = regexp.MustCompile(`<a\s+[^>]*href=["']([^"']+)["'][^>]*>`)

	// Image links: ![alt](url)
	imageRegex = regexp.MustCompile(`!\[([^\]]*)\]\(([^)\s]+)(?:\s+"[^"]*")?\)`)
)

// CollectFromFile extracts all links from a markdown file
func (c *LinkCollector) CollectFromFile(filePath string) error {
	file, err := os.Open(filePath) // #nosec G304 -- path from trusted user configuration
	if err != nil {
		return err
	}
	defer file.Close()

	relPath, _ := filepath.Rel(c.docsRoot, filePath)

	scanner := bufio.NewScanner(file)
	lineNum := 0
	inCodeBlock := false

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Track code blocks to skip link extraction inside them
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		if inCodeBlock {
			continue
		}

		// Extract markdown links
		c.extractMarkdownLinks(line, relPath, lineNum)

		// Extract HTML links
		c.extractHTMLLinks(line, relPath, lineNum)

		// Extract image links
		c.extractImageLinks(line, relPath, lineNum)
	}

	return scanner.Err()
}

// extractMarkdownLinks extracts [text](url) style links
func (c *LinkCollector) extractMarkdownLinks(line, filePath string, lineNum int) {
	matches := markdownLinkRegex.FindAllStringSubmatchIndex(line, -1)

	for _, match := range matches {
		if len(match) >= 6 {
			// Skip image links (preceded by !) to avoid double-collecting
			if match[0] > 0 && line[match[0]-1] == '!' {
				continue
			}
			text := line[match[2]:match[3]]
			url := line[match[4]:match[5]]
			col := match[4] + 1 // 1-based column

			c.addLink(url, filePath, lineNum, col, text)
		}
	}
}

// extractHTMLLinks extracts <a href="url"> style links
func (c *LinkCollector) extractHTMLLinks(line, filePath string, lineNum int) {
	matches := htmlLinkRegex.FindAllStringSubmatchIndex(line, -1)

	for _, match := range matches {
		if len(match) >= 4 {
			url := line[match[2]:match[3]]
			col := match[2] + 1

			c.addLink(url, filePath, lineNum, col, "")
		}
	}
}

// extractImageLinks extracts ![alt](url) style links
func (c *LinkCollector) extractImageLinks(line, filePath string, lineNum int) {
	matches := imageRegex.FindAllStringSubmatchIndex(line, -1)

	for _, match := range matches {
		if len(match) >= 6 {
			alt := line[match[2]:match[3]]
			url := line[match[4]:match[5]]
			col := match[4] + 1

			link := core.CollectedLink{
				URL:        url,
				SourceFile: filePath,
				Line:       lineNum,
				Column:     col,
				Text:       alt,
				LinkType:   classifyLink(url),
			}

			// Override type for images - they're assets
			if link.LinkType == core.LinkTypeInternalPage {
				link.LinkType = core.LinkTypeInternalAsset
			}

			c.links = append(c.links, link)
		}
	}
}

// addLink adds a link to the collection
func (c *LinkCollector) addLink(url, filePath string, line, col int, text string) {
	link := core.CollectedLink{
		URL:        url,
		SourceFile: filePath,
		Line:       line,
		Column:     col,
		Text:       text,
		LinkType:   classifyLink(url),
	}

	c.links = append(c.links, link)
}

// classifyLink determines the type of a link
func classifyLink(url string) core.LinkType {
	// Empty or whitespace
	if strings.TrimSpace(url) == "" {
		return core.LinkTypeOther
	}

	// External links
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "//") {
		return core.LinkTypeExternal
	}

	// Email links
	if strings.HasPrefix(url, "mailto:") {
		return core.LinkTypeEmail
	}

	// Other protocols
	if strings.Contains(url, ":") && !strings.HasPrefix(url, "/") {
		return core.LinkTypeOther
	}

	// Anchor-only links
	if strings.HasPrefix(url, "#") {
		return core.LinkTypeInternalAnchor
	}

	// Check if it's an asset (image, file, etc.)
	ext := strings.ToLower(filepath.Ext(strings.Split(url, "#")[0]))
	assetExtensions := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".svg": true, ".webp": true,
		".pdf": true, ".zip": true, ".tar": true, ".gz": true,
		".css": true, ".js": true, ".json": true, ".xml": true,
		".mp4": true, ".webm": true, ".mp3": true, ".wav": true,
	}

	if assetExtensions[ext] {
		return core.LinkTypeInternalAsset
	}

	// Default to internal page
	return core.LinkTypeInternalPage
}

// Links returns all collected links
func (c *LinkCollector) Links() []core.CollectedLink {
	return c.links
}

// CollectFromPages collects links from all pages in the site
func (c *LinkCollector) CollectFromPages(pages []*core.Page) {
	for _, page := range pages {
		if err := c.CollectFromFile(page.SourcePath); err != nil {
			// Log warning but continue
			continue
		}
	}
}
