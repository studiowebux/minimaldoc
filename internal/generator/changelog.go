package generator

import (
	"bytes"
	"embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// ChangelogGenerator generates changelog HTML and related files
type ChangelogGenerator struct {
	site      *core.Site
	templates *template.Template
	themeFS   embed.FS
	version   string
}

// NewChangelogGenerator creates a new changelog generator
func NewChangelogGenerator(site *core.Site, themeFS embed.FS, version string) (*ChangelogGenerator, error) {
	if !site.Config.Changelog.Enabled {
		return nil, nil // Skip if changelog is not enabled
	}

	// Create template with shared changelog functions
	tmpl := template.New("").Funcs(ChangelogFuncMap()).Funcs(AnalyticsFuncMap())

	// Parse changelog templates from dedicated subdirectory
	var err error
	tmpl, err = tmpl.ParseFS(
		themeFS,
		"themes/common/templates/partials/landing-*.html",
		"themes/common/templates/partials/analytics.html",
		"themes/common/templates/partials/minimaldoc-widgets.html",
		"themes/common/templates/changelog/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse changelog templates: %w", err)
	}

	return &ChangelogGenerator{
		site:      site,
		templates: tmpl,
		themeFS:   themeFS,
		version:   version,
	}, nil
}

// Generate generates all changelog files
func (g *ChangelogGenerator) Generate() error {
	if g.site.ChangelogPage == nil || !g.site.Config.Changelog.Enabled {
		return nil
	}

	fmt.Println("Generating changelog...")

	changelogPath := g.site.Config.Changelog.Path
	if changelogPath == "" {
		changelogPath = "changelog"
	}

	outputDir := filepath.Join(g.site.OutputRoot, changelogPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create changelog output directory: %w", err)
	}

	// Generate main changelog page
	if err := g.generateIndexPage(outputDir); err != nil {
		return fmt.Errorf("failed to generate changelog index page: %w", err)
	}

	// Generate individual release pages
	if err := g.generateReleasePages(outputDir); err != nil {
		return fmt.Errorf("failed to generate release pages: %w", err)
	}

	// Generate changelog.json
	if err := g.generateChangelogJSON(outputDir); err != nil {
		return fmt.Errorf("failed to generate changelog.json: %w", err)
	}

	// Generate RSS feed
	if g.site.Config.Changelog.RSSEnabled {
		if err := g.generateRSSFeed(outputDir); err != nil {
			return fmt.Errorf("failed to generate RSS feed: %w", err)
		}
	}

	fmt.Printf("Generated changelog with %d releases\n", len(g.site.ChangelogPage.Releases))
	return nil
}

// generateIndexPage generates the main changelog page
func (g *ChangelogGenerator) generateIndexPage(outputDir string) error {
	changelogPath := g.site.Config.Changelog.Path
	if changelogPath == "" {
		changelogPath = "changelog"
	}
	data := map[string]any{
		"Site":          g.site,
		"ChangelogPage": g.site.ChangelogPage,
		"BasePath":      g.getBasePath(),
		"Version":       g.version,
		"PageTitle":     g.site.ChangelogPage.Config.Title,
		"ActivePath":    "/" + changelogPath + "/",
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "changelog.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	outputPath := filepath.Join(outputDir, "index.html")
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write changelog page: %w", err)
	}

	return nil
}

// generateReleasePages generates individual release detail pages
func (g *ChangelogGenerator) generateReleasePages(outputDir string) error {
	for _, release := range g.site.ChangelogPage.Releases {
		if err := g.generateReleasePage(outputDir, release); err != nil {
			return fmt.Errorf("failed to generate release %s: %w", release.Version, err)
		}
	}
	return nil
}

// generateReleasePage generates a single release detail page
func (g *ChangelogGenerator) generateReleasePage(outputDir string, release core.Release) error {
	clPath := g.site.Config.Changelog.Path
	if clPath == "" {
		clPath = "changelog"
	}
	data := map[string]any{
		"Site":          g.site,
		"ChangelogPage": g.site.ChangelogPage,
		"Release":       release,
		"BasePath":      g.getBasePath(),
		"Version":       g.version,
		"PageTitle":     fmt.Sprintf("%s - %s", release.Version, g.site.ChangelogPage.Config.Title),
		"ActivePath":    "/" + clPath + "/",
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "changelog-release.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	outputPath := filepath.Join(outputDir, release.Slug+".html")
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write release page: %w", err)
	}

	return nil
}

// ChangelogJSON represents the JSON API response for changelog
type ChangelogJSON struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Releases    []ReleaseJSON `json:"releases"`
	LastUpdated time.Time     `json:"last_updated"`
}

// ReleaseJSON represents a release in the JSON API
type ReleaseJSON struct {
	Version    string         `json:"version"`
	Date       string         `json:"date"`
	Title      string         `json:"title,omitempty"`
	Prerelease bool           `json:"prerelease,omitempty"`
	URL        string         `json:"url"`
	CompareURL string         `json:"compare_url,omitempty"`
	Categories []CategoryJSON `json:"categories,omitempty"`
}

// CategoryJSON represents a change category in the JSON API
type CategoryJSON struct {
	Type    string   `json:"type"`
	Entries []string `json:"entries"`
}

// generateChangelogJSON generates the changelog.json API file
func (g *ChangelogGenerator) generateChangelogJSON(outputDir string) error {
	changelogPath := g.site.Config.Changelog.Path
	if changelogPath == "" {
		changelogPath = "changelog"
	}

	releases := make([]ReleaseJSON, 0, len(g.site.ChangelogPage.Releases))
	for _, r := range g.site.ChangelogPage.Releases {
		releaseURL := g.getBasePath() + "/" + changelogPath + "/" + r.Slug + ".html"

		categories := make([]CategoryJSON, 0, len(r.Categories))
		for _, cat := range r.Categories {
			entries := make([]string, 0, len(cat.Entries))
			for _, e := range cat.Entries {
				entries = append(entries, e.Description)
			}
			categories = append(categories, CategoryJSON{
				Type:    string(cat.Type),
				Entries: entries,
			})
		}

		releases = append(releases, ReleaseJSON{
			Version:    r.Version,
			Date:       r.Date.Format("2006-01-02"),
			Title:      r.Title,
			Prerelease: r.Prerelease,
			URL:        releaseURL,
			CompareURL: r.CompareURL,
			Categories: categories,
		})
	}

	changelogJSON := ChangelogJSON{
		Title:       g.site.ChangelogPage.Config.Title,
		Description: g.site.ChangelogPage.Config.Description,
		Releases:    releases,
		LastUpdated: g.site.ChangelogPage.LastUpdated,
	}

	data, err := json.MarshalIndent(changelogJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal changelog JSON: %w", err)
	}

	outputPath := filepath.Join(outputDir, "changelog.json")
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write changelog.json: %w", err)
	}

	return nil
}

// RSS feed structures for changelog
type changelogRSSChannel struct {
	XMLName       xml.Name           `xml:"channel"`
	Title         string             `xml:"title"`
	Link          string             `xml:"link"`
	Description   string             `xml:"description"`
	LastBuildDate string             `xml:"lastBuildDate"`
	Items         []changelogRSSItem `xml:"item"`
}

type changelogRSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

type changelogRSSFeed struct {
	XMLName xml.Name            `xml:"rss"`
	Version string              `xml:"version,attr"`
	Channel changelogRSSChannel `xml:"channel"`
}

// generateRSSFeed generates an RSS feed for changelog releases
func (g *ChangelogGenerator) generateRSSFeed(outputDir string) error {
	changelogPath := g.site.Config.Changelog.Path
	if changelogPath == "" {
		changelogPath = "changelog"
	}

	baseURL := strings.TrimSuffix(g.site.Config.BaseURL, "/")

	// Build RSS items from releases
	items := make([]changelogRSSItem, 0, len(g.site.ChangelogPage.Releases))

	for _, release := range g.site.ChangelogPage.Releases {
		releaseURL := baseURL + "/" + changelogPath + "/" + release.Slug + ".html"

		// Build description from categories
		var descParts []string
		for _, cat := range release.Categories {
			if len(cat.Entries) > 0 {
				descParts = append(descParts, fmt.Sprintf("%s: %d changes", cat.Type, len(cat.Entries)))
			}
		}
		description := strings.Join(descParts, ", ")
		if description == "" {
			description = "New release"
		}

		title := release.Version
		if release.Title != "" {
			title = fmt.Sprintf("%s - %s", release.Version, release.Title)
		}
		if release.Prerelease {
			title = fmt.Sprintf("[Prerelease] %s", title)
		}

		items = append(items, changelogRSSItem{
			Title:       title,
			Link:        releaseURL,
			Description: description,
			PubDate:     release.Date.Format(time.RFC1123Z),
			GUID:        releaseURL,
		})
	}

	feed := changelogRSSFeed{
		Version: "2.0",
		Channel: changelogRSSChannel{
			Title:         g.site.ChangelogPage.Config.Title,
			Link:          baseURL + "/" + changelogPath,
			Description:   g.site.ChangelogPage.Config.Description,
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         items,
		},
	}

	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal RSS feed: %w", err)
	}

	// Add XML declaration
	xmlData := []byte(xml.Header + string(data))

	outputPath := filepath.Join(outputDir, "feed.xml")
	if err := os.WriteFile(outputPath, xmlData, 0644); err != nil {
		return fmt.Errorf("failed to write RSS feed: %w", err)
	}

	return nil
}

// getBasePath extracts the path component from BaseURL for asset linking
func (g *ChangelogGenerator) getBasePath() string {
	baseURL := g.site.Config.BaseURL
	if baseURL == "" {
		return ""
	}

	// Remove protocol
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

	// Get the path part
	path := "/" + parts[1]
	path = strings.TrimSuffix(path, "/")

	if path == "/" {
		return ""
	}

	return path
}
