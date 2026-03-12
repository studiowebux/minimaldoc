package generator

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// LegalGenerator generates legal page HTML
type LegalGenerator struct {
	site      *core.Site
	templates *template.Template
	themeFS   embed.FS
	version   string
}

// NewLegalGenerator creates a new legal generator
func NewLegalGenerator(site *core.Site, themeFS embed.FS, version string) (*LegalGenerator, error) {
	if !site.Config.Legal.Enabled {
		return nil, nil
	}

	tmpl := template.New("").Funcs(BaseFuncMap()).Funcs(AnalyticsFuncMap())

	var err error
	tmpl, err = tmpl.ParseFS(
		themeFS,
		"themes/common/templates/partials/landing-*.html",
		"themes/common/templates/partials/analytics.html",
		"themes/common/templates/partials/minimaldoc-widgets.html",
		"themes/common/templates/legal/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse legal templates: %w", err)
	}

	return &LegalGenerator{
		site:      site,
		templates: tmpl,
		themeFS:   themeFS,
		version:   version,
	}, nil
}

// Generate generates all legal pages
func (g *LegalGenerator) Generate() error {
	if len(g.site.LegalPages) == 0 || !g.site.Config.Legal.Enabled {
		return nil
	}

	fmt.Println("Generating legal pages...")

	legalPath := g.site.Config.Legal.Path
	if legalPath == "" {
		legalPath = "legal"
	}

	outputDir := filepath.Join(g.site.OutputRoot, legalPath)
	if err := makeWebDir(outputDir); err != nil {
		return fmt.Errorf("failed to create legal output directory: %w", err)
	}

	// Generate each legal page
	for _, page := range g.site.LegalPages {
		if err := g.generatePage(outputDir, page); err != nil {
			return fmt.Errorf("failed to generate legal page %s: %w", page.Slug, err)
		}
	}

	fmt.Printf("Generated %d legal pages\n", len(g.site.LegalPages))
	return nil
}

// generatePage generates a single legal page
func (g *LegalGenerator) generatePage(outputDir string, page *core.LegalPage) error {
	data := map[string]any{
		"Site":      g.site,
		"LegalPage": page,
		"Footer":    BuildFooter(g.site, g.version),
		"BasePath":  g.getBasePath(),
		"Version":   g.version,
		"PageTitle": page.Title + " | " + g.site.Config.Title,
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "legal.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	// Create page directory
	pageDir := filepath.Join(outputDir, page.Slug)
	if err := makeWebDir(pageDir); err != nil {
		return fmt.Errorf("failed to create page directory: %w", err)
	}

	outputPath := filepath.Join(pageDir, "index.html")
	if err := writeWebFile(outputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write legal page: %w", err)
	}

	page.OutputPath = outputPath
	return nil
}

// getBasePath extracts the path component from BaseURL
func (g *LegalGenerator) getBasePath() string {
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
