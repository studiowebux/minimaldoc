package generator

import (
	"bytes"
	"embed"
	"encoding/json"
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
	// Create template with custom functions
	tmpl := template.New("").Funcs(template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict requires an even number of arguments")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"json": func(v interface{}) (template.JS, error) {
			bytes, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return template.JS(bytes), nil
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"lower": strings.ToLower,
		"replace": func(input, old, new string) string {
			return strings.ReplaceAll(input, old, new)
		},
		"stripSpecExt": func(name string) string {
			// Remove common OpenAPI spec extensions
			name = strings.TrimSuffix(name, ".yaml")
			name = strings.TrimSuffix(name, ".yml")
			name = strings.TrimSuffix(name, ".json")
			return name
		},
		"endpointID": func(endpoint *core.APIEndpoint) string {
			if endpoint.OperationID != "" {
				return endpoint.OperationID
			}
			// Fallback: use Method-Path with non-alphanumeric chars replaced
			id := endpoint.Method + "-" + endpoint.Path
			// Replace non-alphanumeric characters with dashes
			result := ""
			lastWasDash := false
			for _, r := range id {
				if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
					result += string(r)
					lastWasDash = false
				} else if !lastWasDash {
					result += "-"
					lastWasDash = true
				}
			}
			// Trim trailing dash
			return strings.TrimSuffix(result, "-")
		},
	})

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
	// Prepare template data
	data := map[string]interface{}{
		"Site":     g.site,
		"Page":     page,
		"Content":  template.HTML(page.HTML),
		"BasePath": g.getBasePath(),
		"Version":  g.version,
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
