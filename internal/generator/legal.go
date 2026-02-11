package generator

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
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

	tmpl := template.New("").Funcs(template.FuncMap{
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict requires an even number of arguments")
			}
			dict := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"lower":     strings.ToLower,
		"hasPrefix": strings.HasPrefix,
	})

	var err error
	tmpl, err = tmpl.ParseFS(
		themeFS,
		"themes/common/templates/partials/landing-*.html",
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
	if err := os.MkdirAll(outputDir, 0755); err != nil {
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
		"Footer":    g.buildFooterWithLegal(),
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
	if err := os.MkdirAll(pageDir, 0755); err != nil {
		return fmt.Errorf("failed to create page directory: %w", err)
	}

	outputPath := filepath.Join(pageDir, "index.html")
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
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

// buildFooterWithLegal creates a footer config with auto-generated legal links
func (g *LegalGenerator) buildFooterWithLegal() core.FooterConfig {
	footer := g.site.Config.Footer

	if g.site.Config.Legal.Enabled && len(g.site.LegalPages) > 0 {
		legalPath := g.site.Config.Legal.Path
		if legalPath == "" {
			legalPath = "legal"
		}

		groupTitle := g.site.Config.Legal.FooterGroup
		if groupTitle == "" {
			groupTitle = "Legal"
		}

		var legalLinks []core.FooterLink
		for _, page := range g.site.LegalPages {
			legalLinks = append(legalLinks, core.FooterLink{
				Text: page.Title,
				URL:  "/" + legalPath + "/" + page.Slug + "/",
			})
		}

		legalGroup := core.FooterLinkGroup{
			Title: groupTitle,
			Items: legalLinks,
		}
		footer.Links = append(footer.Links, legalGroup)
	}

	return footer
}
