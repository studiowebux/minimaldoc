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

// LandingGenerator generates landing page HTML
type LandingGenerator struct {
	site      *core.Site
	templates *template.Template
	themeFS   embed.FS
	version   string
}

// NewLandingGenerator creates a new landing generator
func NewLandingGenerator(site *core.Site, themeFS embed.FS, version string) (*LandingGenerator, error) {
	if !site.Config.Landing.Enabled {
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
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"add": func(a, b int) int {
			return a + b
		},
		"hasPrefix": strings.HasPrefix,
	})

	var err error
	tmpl, err = tmpl.ParseFS(
		themeFS,
		"themes/common/templates/partials/landing-*.html",
		"themes/common/templates/landing/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse landing templates: %w", err)
	}

	return &LandingGenerator{
		site:      site,
		templates: tmpl,
		themeFS:   themeFS,
		version:   version,
	}, nil
}

// Generate generates the landing page
func (g *LandingGenerator) Generate() error {
	if g.site.LandingPage == nil || !g.site.Config.Landing.Enabled {
		return nil
	}

	fmt.Println("Generating landing page...")

	// Generate main landing page as index.html
	if err := g.generateMainPage(); err != nil {
		return fmt.Errorf("failed to generate landing page: %w", err)
	}

	fmt.Println("Generated landing page")
	return nil
}

// generateMainPage generates the main landing page
func (g *LandingGenerator) generateMainPage() error {
	// Build footer with auto-generated legal links
	footer := g.buildFooterWithLegal()

	data := map[string]any{
		"Site":        g.site,
		"LandingPage": g.site.LandingPage,
		"Footer":      footer,
		"BasePath":    g.getBasePath(),
		"Version":     g.version,
		"PageTitle":   g.site.Config.Title,
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "landing.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	outputPath := filepath.Join(g.site.OutputRoot, "index.html")
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write landing page: %w", err)
	}

	return nil
}

// getBasePath extracts the path component from BaseURL
func (g *LandingGenerator) getBasePath() string {
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
func (g *LandingGenerator) buildFooterWithLegal() core.FooterConfig {
	footer := g.site.Config.Footer

	// If legal pages are enabled, auto-generate a "Legal" footer group
	if g.site.Config.Legal.Enabled && len(g.site.LegalPages) > 0 {
		legalPath := g.site.Config.Legal.Path
		if legalPath == "" {
			legalPath = "legal"
		}

		groupTitle := g.site.Config.Legal.FooterGroup
		if groupTitle == "" {
			groupTitle = "Legal"
		}

		// Build legal links (no basePath - template will add it)
		var legalLinks []core.FooterLink
		for _, page := range g.site.LegalPages {
			legalLinks = append(legalLinks, core.FooterLink{
				Text: page.Title,
				URL:  "/" + legalPath + "/" + page.Slug + "/",
			})
		}

		// Add legal group to footer
		legalGroup := core.FooterLinkGroup{
			Title: groupTitle,
			Items: legalLinks,
		}
		footer.Links = append(footer.Links, legalGroup)
	}

	return footer
}
