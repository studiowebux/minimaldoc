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

// PortfolioGenerator generates portfolio page HTML
type PortfolioGenerator struct {
	site      *core.Site
	templates *template.Template
	themeFS   embed.FS
	version   string
}

// NewPortfolioGenerator creates a new portfolio generator
func NewPortfolioGenerator(site *core.Site, themeFS embed.FS, version string) (*PortfolioGenerator, error) {
	if !site.Config.Portfolio.Enabled {
		return nil, nil
	}

	tmpl := template.New("").Funcs(template.FuncMap{
		"hasPrefix": strings.HasPrefix,
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
		"join": strings.Join,
		"formatDate": func(t any) string {
			if tm, ok := t.(interface{ Format(string) string }); ok {
				return tm.Format("January 2006")
			}
			return ""
		},
	})

	var err error
	tmpl, err = tmpl.ParseFS(
		themeFS,
		"themes/common/templates/partials/landing-*.html",
		"themes/common/templates/portfolio/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse portfolio templates: %w", err)
	}

	return &PortfolioGenerator{
		site:      site,
		templates: tmpl,
		themeFS:   themeFS,
		version:   version,
	}, nil
}

// Generate generates the portfolio pages
func (g *PortfolioGenerator) Generate() error {
	if g.site.PortfolioPage == nil || !g.site.Config.Portfolio.Enabled {
		return nil
	}

	fmt.Println("Generating portfolio page...")

	portfolioPath := g.site.Config.Portfolio.Path
	if portfolioPath == "" {
		portfolioPath = "portfolio"
	}

	outputDir := filepath.Join(g.site.OutputRoot, portfolioPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create portfolio output directory: %w", err)
	}

	// Generate main portfolio page
	if err := g.generateMainPage(outputDir); err != nil {
		return fmt.Errorf("failed to generate portfolio page: %w", err)
	}

	// Generate individual project pages
	if err := g.generateProjectPages(outputDir); err != nil {
		return fmt.Errorf("failed to generate project pages: %w", err)
	}

	fmt.Printf("Generated portfolio page with %d projects\n", len(g.site.PortfolioPage.Projects))
	return nil
}

// generateMainPage generates the main portfolio listing page
func (g *PortfolioGenerator) generateMainPage(outputDir string) error {
	portfolioPath := g.site.Config.Portfolio.Path
	if portfolioPath == "" {
		portfolioPath = "portfolio"
	}
	data := map[string]any{
		"Site":          g.site,
		"PortfolioPage": g.site.PortfolioPage,
		"Footer":        g.buildFooterWithLegal(),
		"BasePath":      g.getBasePath(),
		"Version":       g.version,
		"PageTitle":     g.site.PortfolioPage.Config.Title + " | " + g.site.Config.Title,
		"ActivePath":    "/" + portfolioPath + "/",
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "portfolio.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	outputPath := filepath.Join(outputDir, "index.html")
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write portfolio page: %w", err)
	}

	return nil
}

// generateProjectPages generates individual project pages
func (g *PortfolioGenerator) generateProjectPages(outputDir string) error {
	projectDir := filepath.Join(outputDir, "project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	for _, project := range g.site.PortfolioPage.Projects {
		if err := g.generateProjectPage(projectDir, project); err != nil {
			return fmt.Errorf("failed to generate project %s: %w", project.Slug, err)
		}
	}

	return nil
}

// generateProjectPage generates a single project detail page
func (g *PortfolioGenerator) generateProjectPage(projectDir string, project core.Project) error {
	pPath := g.site.Config.Portfolio.Path
	if pPath == "" {
		pPath = "portfolio"
	}
	data := map[string]any{
		"Site":          g.site,
		"PortfolioPage": g.site.PortfolioPage,
		"Project":       project,
		"Footer":        g.buildFooterWithLegal(),
		"BasePath":      g.getBasePath(),
		"Version":       g.version,
		"PageTitle":     project.Title + " | " + g.site.PortfolioPage.Config.Title + " | " + g.site.Config.Title,
		"ActivePath":    "/" + pPath + "/",
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "project.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	outputPath := filepath.Join(projectDir, project.Slug+".html")
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write project page: %w", err)
	}

	return nil
}

// getBasePath extracts the path component from BaseURL
func (g *PortfolioGenerator) getBasePath() string {
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
func (g *PortfolioGenerator) buildFooterWithLegal() core.FooterConfig {
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
