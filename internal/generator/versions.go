package generator

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// VersionGenerator generates versioned documentation output
type VersionGenerator struct {
	site      *core.Site
	templates *template.Template
	themeFS   embed.FS
	version   string // MinimalDoc version
}

// VersionsJSON is the structure for versions.json metadata file
type VersionsJSON struct {
	Current  string            `json:"current"`
	Default  string            `json:"default"`
	Versions []VersionInfoJSON `json:"versions"`
}

// VersionInfoJSON represents version info in JSON output
type VersionInfoJSON struct {
	Name   string `json:"name"`
	Label  string `json:"label"`
	Path   string `json:"path"`
	EOL    string `json:"eol,omitempty"`
	Active bool   `json:"active"`
}

// NewVersionGenerator creates a new version generator
func NewVersionGenerator(site *core.Site, themeFS embed.FS, version string) (*VersionGenerator, error) {
	tmpl := template.New("").Funcs(OpenAPIFuncMap()).Funcs(MCPFuncMap()).Funcs(AnalyticsFuncMap())

	tmpl, err := tmpl.ParseFS(
		themeFS,
		"themes/common/templates/*.html",
		"themes/common/templates/partials/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse common templates: %w", err)
	}

	return &VersionGenerator{
		site:      site,
		templates: tmpl,
		themeFS:   themeFS,
		version:   version,
	}, nil
}

// Generate generates versioned documentation
func (g *VersionGenerator) Generate() error {
	if !g.site.Config.Versions.Enabled || len(g.site.Config.Versions.List) == 0 {
		return nil
	}

	fmt.Println("Generating versioned documentation...")

	defaultVersion := g.site.Config.Versions.Default
	if defaultVersion == "" && len(g.site.Config.Versions.List) > 0 {
		defaultVersion = g.site.Config.Versions.List[0].Name
	}

	// Generate pages for each version
	for _, versionInfo := range g.site.Config.Versions.List {
		pages, ok := g.site.VersionedPages[versionInfo.Name]
		if !ok {
			continue
		}

		isDefault := versionInfo.Name == defaultVersion
		if err := g.generateVersionPages(versionInfo, pages, isDefault); err != nil {
			return fmt.Errorf("failed to generate version %s: %w", versionInfo.Name, err)
		}

		fmt.Printf("  Generated version %s: %d pages\n", versionInfo.Name, len(pages))
	}

	// Generate versions.json metadata
	if err := g.generateVersionsJSON(defaultVersion); err != nil {
		return fmt.Errorf("failed to generate versions.json: %w", err)
	}

	return nil
}

// generateVersionPages generates HTML pages for a specific version
func (g *VersionGenerator) generateVersionPages(versionInfo core.VersionInfo, pages []*core.Page, isDefault bool) error {
	basePath := g.getBasePath()

	for _, page := range pages {
		// Determine output path
		var outputPath string
		if isDefault {
			// Default version goes to root
			outputPath = filepath.Join(g.site.OutputRoot, page.Slug+".html")
		} else {
			// Non-default versions get prefixed
			outputPath = filepath.Join(g.site.OutputRoot, versionInfo.Path, page.Slug+".html")
		}

		// Prepare template data
		data := map[string]any{
			"Site":           g.site,
			"Page":           page,
			"Content":        template.HTML(page.HTML),
			"BasePath":       basePath,
			"Version":        g.version,
			"CurrentVersion": versionInfo,
			"IsDefault":      isDefault,
			"VersionPrefix":  g.getVersionPrefix(versionInfo, isDefault),
		}

		// Execute template
		var buf bytes.Buffer
		if err := g.templates.ExecuteTemplate(&buf, "layout.html", data); err != nil {
			return fmt.Errorf("template execution failed for %s: %w", page.Slug, err)
		}

		// Create output directory
		outputDir := filepath.Dir(outputPath)
		if err := makeWebDir(outputDir); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Write HTML file
		if err := writeWebFile(outputPath, buf.Bytes()); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	return nil
}

// generateVersionsJSON creates the versions.json metadata file
func (g *VersionGenerator) generateVersionsJSON(defaultVersion string) error {
	versionsJSON := VersionsJSON{
		Current:  g.site.CurrentVersion,
		Default:  defaultVersion,
		Versions: make([]VersionInfoJSON, 0, len(g.site.Config.Versions.List)),
	}

	for _, v := range g.site.Config.Versions.List {
		versionsJSON.Versions = append(versionsJSON.Versions, VersionInfoJSON{
			Name:   v.Name,
			Label:  v.Label,
			Path:   v.Path,
			EOL:    v.EOL,
			Active: v.Name == defaultVersion,
		})
	}

	data, err := json.MarshalIndent(versionsJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal versions.json: %w", err)
	}

	outputPath := filepath.Join(g.site.OutputRoot, "versions.json")
	if err := writeWebFile(outputPath, data); err != nil {
		return fmt.Errorf("failed to write versions.json: %w", err)
	}

	return nil
}

// getVersionPrefix returns the URL prefix for a version
func (g *VersionGenerator) getVersionPrefix(versionInfo core.VersionInfo, isDefault bool) string {
	if isDefault {
		return ""
	}
	return "/" + versionInfo.Path
}

// getBasePath extracts the path component from BaseURL
func (g *VersionGenerator) getBasePath() string {
	baseURL := g.site.Config.BaseURL
	if baseURL == "" {
		return ""
	}

	if strings.HasPrefix(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	} else if strings.HasPrefix(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	}

	parts := strings.SplitN(baseURL, "/", 2)
	if len(parts) < 2 {
		return ""
	}

	path := "/" + parts[1]
	path = strings.TrimSuffix(path, "/")

	if path == "/" {
		return ""
	}

	return path
}
