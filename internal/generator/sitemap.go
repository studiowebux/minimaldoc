package generator

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// URLSet represents the root element of a sitemap
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// URL represents a single URL in the sitemap
type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float32 `xml:"priority,omitempty"`
}

// SitemapGenerator generates sitemap.xml
type SitemapGenerator struct {
	site *core.Site
}

// NewSitemapGenerator creates a new sitemap generator
func NewSitemapGenerator(site *core.Site) *SitemapGenerator {
	return &SitemapGenerator{site: site}
}

// Generate creates sitemap.xml
func (g *SitemapGenerator) Generate() error {
	if g.site.Config.BaseURL == "" {
		fmt.Println("Skipping sitemap generation (no baseURL configured)")
		return nil
	}

	// Validate that baseURL is absolute (contains protocol)
	baseURL := strings.TrimSuffix(g.site.Config.BaseURL, "/")
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		return fmt.Errorf("base_url must be an absolute URL (e.g., https://example.com), got: %s", g.site.Config.BaseURL)
	}

	fmt.Println("Generating sitemap.xml...")

	urlSet := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  []URL{},
	}

	// Add all pages
	for _, page := range g.site.Pages {
		if page.IsHidden() {
			continue
		}

		// Get file modification time
		fileInfo, err := os.Stat(page.SourcePath)
		var lastMod string
		if err == nil {
			lastMod = fileInfo.ModTime().Format(time.RFC3339)
		}

		// Determine priority (index gets highest priority)
		priority := float32(0.5)
		if page.Slug == "index" {
			priority = 1.0
		} else if strings.Count(page.Slug, "/") == 0 {
			priority = 0.8 // Top-level pages
		}

		// Determine change frequency
		changeFreq := "monthly"
		if page.Slug == "index" {
			changeFreq = "weekly"
		}

		// Remove leading slash from slug to avoid double slashes
		slug := strings.TrimPrefix(page.Slug, "/")

		url := URL{
			Loc:        fmt.Sprintf("%s/%s.html", baseURL, slug),
			LastMod:    lastMod,
			ChangeFreq: changeFreq,
			Priority:   priority,
		}

		urlSet.URLs = append(urlSet.URLs, url)
	}

	// Marshal to XML
	output, err := xml.MarshalIndent(urlSet, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sitemap: %w", err)
	}

	// Add XML declaration
	xmlContent := []byte(xml.Header + string(output))

	// Write to file
	sitemapPath := filepath.Join(g.site.OutputRoot, "sitemap.xml")
	if err := os.WriteFile(sitemapPath, xmlContent, 0644); err != nil {
		return fmt.Errorf("failed to write sitemap: %w", err)
	}

	fmt.Printf("Generated sitemap.xml with %d URLs\n", len(urlSet.URLs))
	return nil
}
