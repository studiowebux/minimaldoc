package generator

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// HTMLGenerator generates HTML files from pages
type HTMLGenerator struct {
	site      *core.Site
	templates *template.Template
	themeFS   embed.FS
	version   string
}

// NewHTMLGenerator creates a new HTML generator
func NewHTMLGenerator(site *core.Site, themeFS embed.FS, version string) (*HTMLGenerator, error) {
	// Create template with shared functions (includes OpenAPI helpers for layout compatibility)
	tmpl := template.New("").Funcs(OpenAPIFuncMap())

	// Parse templates from embedded filesystem using configured theme
	themeName := site.Config.Theme
	if themeName == "" {
		themeName = "default"
	}

	// Parse common templates (all structure is in common, themes only provide CSS)
	var err error

	tmpl, err = tmpl.ParseFS(
		themeFS,
		"themes/common/templates/*.html",
		"themes/common/templates/partials/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse common templates: %w", err)
	}

	return &HTMLGenerator{
		site:      site,
		templates: tmpl,
		themeFS:   themeFS,
		version:   version,
	}, nil
}

// Generate generates HTML files for all pages
func (g *HTMLGenerator) Generate() error {
	fmt.Println("Generating HTML files...")

	// Create output directory
	if err := os.MkdirAll(g.site.OutputRoot, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate each page
	for _, page := range g.site.Pages {
		if err := g.generatePage(page); err != nil {
			return fmt.Errorf("failed to generate page %s: %w", page.SourcePath, err)
		}
	}

	// Copy static assets
	if err := g.copyStaticAssets(); err != nil {
		return fmt.Errorf("failed to copy static assets: %w", err)
	}

	fmt.Printf("Generated %d pages\n", len(g.site.Pages))
	return nil
}

// generatePage generates a single HTML page
func (g *HTMLGenerator) generatePage(page *core.Page) error {
	// Get current version info for template
	var currentVersion *core.VersionInfo
	isDefault := true
	if g.site.Config.Versions.Enabled && len(g.site.Config.Versions.List) > 0 {
		defaultVersionName := g.site.Config.Versions.Default
		if defaultVersionName == "" {
			defaultVersionName = g.site.Config.Versions.List[0].Name
		}
		currentVersion = g.site.Config.Versions.GetVersion(defaultVersionName)
	}

	// Prepare template data
	data := map[string]any{
		"Site":           g.site,
		"Page":           page,
		"Content":        template.HTML(page.HTML),
		"BasePath":       g.getBasePath(),
		"Version":        g.version,
		"CurrentVersion": currentVersion,
		"IsDefault":      isDefault,
	}

	// Execute template
	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "layout.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	// Create output directory
	outputDir := filepath.Dir(page.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write HTML file
	if err := os.WriteFile(page.OutputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// copyStaticAssets copies CSS, JS, and other static assets to the output directory
func (g *HTMLGenerator) copyStaticAssets() error {
	fmt.Println("Copying static assets...")

	// Get theme name
	themeName := g.site.Config.Theme
	if themeName == "" {
		themeName = "default"
	}

	// Helper function to copy files from a static directory
	copyFromPath := func(staticPath string) error {
		return fs.WalkDir(g.themeFS, staticPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				// If common directory doesn't exist, skip it silently
				if staticPath == "themes/common/static" {
					return nil
				}
				return err
			}

			if d.IsDir() {
				return nil
			}

			// Skip OpenAPI files if OpenAPI is not enabled
			if !g.site.Config.OpenAPI.Enabled {
				filename := filepath.Base(path)
				// Skip OpenAPI-specific CSS and JS files
				if filename == "openapi.css" ||
					filename == "openapi-explorer.js" ||
					filename == "api-tester.js" ||
					filename == "oauth-handler.js" ||
					filename == "export.js" ||
					filename == "code-copy.js" ||
					filename == "sidebar-resize.js" {
					return nil
				}
			}

			// Read from embedded FS
			content, err := g.themeFS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			// Determine output path (remove "themes/*/static/" prefix)
			relPath := filepath.Join(strings.TrimPrefix(path, staticPath+"/"))
			outPath := filepath.Join(g.site.OutputRoot, relPath)

			// Create output directory
			if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

			// Write to destination
			if err := os.WriteFile(outPath, content, 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", outPath, err)
			}

			return nil
		})
	}

	// Copy common static assets first
	if err := copyFromPath("themes/common/static"); err != nil {
		return err
	}

	// Copy theme-specific static assets (will override common files if they exist)
	themePath := fmt.Sprintf("themes/%s/static", themeName)
	if err := copyFromPath(themePath); err != nil {
		return err
	}

	return nil
}

// getBasePath extracts the path component from BaseURL for asset linking
// Examples:
//   - "https://example.com/docs/" → "/docs"
//   - "https://example.com/" → ""
//   - "" → ""
func (g *HTMLGenerator) getBasePath() string {
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
