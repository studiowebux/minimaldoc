package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// SearchEntry represents a single page in the search index
type SearchEntry struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
}

// SearchIndex represents the complete search index
type SearchIndex struct {
	Entries []SearchEntry `json:"entries"`
}

// SearchGenerator generates search index JSON
type SearchGenerator struct {
	site *core.Site
}

// NewSearchGenerator creates a new search generator
func NewSearchGenerator(site *core.Site) *SearchGenerator {
	return &SearchGenerator{site: site}
}

// Generate creates search-index.json
func (g *SearchGenerator) Generate() error {
	fmt.Println("Generating search index...")

	basePath := g.getBasePath()

	index := SearchIndex{
		Entries: []SearchEntry{},
	}

	// Add all visible pages to the index, split by sections
	for _, page := range g.site.Pages {
		if page.IsHidden() {
			continue
		}

		// Split page by sections (headings)
		sections := g.splitByHeadings(page.RawMD, page, basePath)

		// If page has sections, index each separately
		if len(sections) > 0 {
			for _, section := range sections {
				index.Entries = append(index.Entries, section)
			}
		} else {
			// Fallback: index entire page if no headings found
			content := extractPlainText(page.RawMD)
			if len(content) > 5000 {
				content = content[:5000] + "..."
			}

			entry := SearchEntry{
				Title:       page.Title(),
				Description: page.Metadata.Description,
				URL:         basePath + "/" + page.Slug + ".html",
				Content:     content,
				Tags:        page.Metadata.Tags,
			}
			index.Entries = append(index.Entries, entry)
		}
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal search index: %w", err)
	}

	// Write to file
	searchPath := filepath.Join(g.site.OutputRoot, "search-index.json")
	if err := os.WriteFile(searchPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write search index: %w", err)
	}

	fmt.Printf("Generated search index with %d entries\n", len(index.Entries))
	return nil
}

// extractPlainText removes markdown syntax and extracts plain text
func extractPlainText(md []byte) string {
	text := string(md)

	// Remove code block fences but keep content
	text = removeCodeBlockFences(text)

	// Remove admonition markers
	text = removeAdmonitions(text)

	// Remove inline code
	text = strings.ReplaceAll(text, "`", "")

	// Remove headings markers
	text = strings.ReplaceAll(text, "#", "")

	// Remove emphasis markers
	text = strings.ReplaceAll(text, "*", "")
	text = strings.ReplaceAll(text, "_", "")

	// Remove links (keep just the text)
	// Simple regex replacement for [text](url)
	for strings.Contains(text, "[") && strings.Contains(text, "]") {
		start := strings.Index(text, "[")
		end := strings.Index(text[start:], "]")
		if end == -1 {
			break
		}
		end += start

		linkText := text[start+1 : end]

		// Find and remove the URL part if exists
		if end+1 < len(text) && text[end+1] == '(' {
			urlEnd := strings.Index(text[end+1:], ")")
			if urlEnd != -1 {
				text = text[:start] + linkText + text[end+urlEnd+2:]
			} else {
				text = text[:start] + linkText + text[end+1:]
			}
		} else {
			text = text[:start] + linkText + text[end+1:]
		}
	}

	// Clean up multiple spaces
	words := strings.Fields(text)
	text = strings.Join(words, " ")

	return text
}

// removeCodeBlockFences removes the fence markers (```) but keeps code content for search
func removeCodeBlockFences(text string) string {
	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip lines that are code block fences (``` or ```language)
		if strings.HasPrefix(trimmed, "```") {
			continue
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// removeAdmonitions removes admonition markers (:::type and :::) from text
func removeAdmonitions(text string) string {
	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip lines that are admonition markers
		if strings.HasPrefix(trimmed, ":::") {
			continue
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// splitByHeadings splits markdown content into sections by headings
func (g *SearchGenerator) splitByHeadings(md []byte, page *core.Page, basePath string) []SearchEntry {
	lines := strings.Split(string(md), "\n")
	var entries []SearchEntry
	var currentHeading string
	var currentAnchor string
	var currentContent []string
	var currentLevel int
	_ = currentLevel // Silence unused warning for now

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if line is a heading
		if strings.HasPrefix(trimmed, "#") {
			// Save previous section if exists
			if currentHeading != "" && len(currentContent) > 0 {
				content := extractPlainText([]byte(strings.Join(currentContent, "\n")))
				if len(content) > 2000 {
					content = content[:2000] + "..."
				}

				url := basePath + "/" + page.Slug + ".html"
				if currentAnchor != "" {
					url += "#" + currentAnchor
				}

				entries = append(entries, SearchEntry{
					Title:       page.Title() + " - " + currentHeading,
					Description: page.Metadata.Description,
					URL:         url,
					Content:     content,
					Tags:        page.Metadata.Tags,
				})
			}

			// Start new section
			level := 0
			for i := 0; i < len(trimmed) && trimmed[i] == '#'; i++ {
				level++
			}
			currentLevel = level
			currentHeading = strings.TrimSpace(trimmed[level:])
			currentAnchor = generateAnchor(currentHeading)
			currentContent = []string{}
		} else {
			// Add line to current section content
			currentContent = append(currentContent, line)
		}
	}

	// Save last section
	if currentHeading != "" && len(currentContent) > 0 {
		content := extractPlainText([]byte(strings.Join(currentContent, "\n")))
		if len(content) > 2000 {
			content = content[:2000] + "..."
		}

		url := basePath + "/" + page.Slug + ".html"
		if currentAnchor != "" {
			url += "#" + currentAnchor
		}

		entries = append(entries, SearchEntry{
			Title:       page.Title() + " - " + currentHeading,
			Description: page.Metadata.Description,
			URL:         url,
			Content:     content,
			Tags:        page.Metadata.Tags,
		})
	}

	return entries
}

// generateAnchor creates a URL-friendly anchor from heading text
// Matches Goldmark's auto-heading-id behavior
func generateAnchor(heading string) string {
	// Convert to lowercase
	anchor := strings.ToLower(heading)

	// Replace spaces with hyphens
	anchor = strings.ReplaceAll(anchor, " ", "-")

	// Remove special characters, keep only alphanumeric and hyphens
	var result strings.Builder
	for _, r := range anchor {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// getBasePath extracts the path component from BaseURL for asset linking
// Examples:
//   - "https://example.com/docs/" → "/docs"
//   - "https://example.com/" → ""
//   - "" → ""
func (g *SearchGenerator) getBasePath() string {
	baseURL := g.site.Config.BaseURL
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
