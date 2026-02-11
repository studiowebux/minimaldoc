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

// ContactGenerator generates contact page HTML
type ContactGenerator struct {
	site      *core.Site
	templates *template.Template
	themeFS   embed.FS
	version   string
}

// NewContactGenerator creates a new contact generator
func NewContactGenerator(site *core.Site, themeFS embed.FS, version string) (*ContactGenerator, error) {
	if !site.Config.Contact.Enabled {
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
		"themes/common/templates/contact/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse contact templates: %w", err)
	}

	return &ContactGenerator{
		site:      site,
		templates: tmpl,
		themeFS:   themeFS,
		version:   version,
	}, nil
}

// Generate generates the contact page
func (g *ContactGenerator) Generate() error {
	if g.site.ContactPage == nil || !g.site.Config.Contact.Enabled {
		return nil
	}

	fmt.Println("Generating contact page...")

	contactPath := g.site.Config.Contact.Path
	if contactPath == "" {
		contactPath = "contact"
	}

	outputDir := filepath.Join(g.site.OutputRoot, contactPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create contact output directory: %w", err)
	}

	// Generate main contact page
	if err := g.generateMainPage(outputDir); err != nil {
		return fmt.Errorf("failed to generate contact page: %w", err)
	}

	fmt.Println("Generated contact page")
	return nil
}

// generateMainPage generates the contact page
func (g *ContactGenerator) generateMainPage(outputDir string) error {
	contactPath := g.site.Config.Contact.Path
	if contactPath == "" {
		contactPath = "contact"
	}
	data := map[string]any{
		"Site":        g.site,
		"ContactPage": g.site.ContactPage,
		"Footer":      g.buildFooterWithLegal(),
		"BasePath":    g.getBasePath(),
		"Version":     g.version,
		"PageTitle":   g.site.ContactPage.Config.Title + " | " + g.site.Config.Title,
		"ActivePath":  "/" + contactPath + "/",
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "contact.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	outputPath := filepath.Join(outputDir, "index.html")
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write contact page: %w", err)
	}

	return nil
}

// getBasePath extracts the path component from BaseURL
func (g *ContactGenerator) getBasePath() string {
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
func (g *ContactGenerator) buildFooterWithLegal() core.FooterConfig {
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
