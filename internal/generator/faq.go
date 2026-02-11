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

// FaqGenerator generates the FAQ page HTML
type FaqGenerator struct {
	site      *core.Site
	templates *template.Template
	themeFS   embed.FS
	version   string
}

// NewFaqGenerator creates a new FAQ generator
func NewFaqGenerator(site *core.Site, themeFS embed.FS, version string) (*FaqGenerator, error) {
	if !site.Config.Faq.Enabled {
		return nil, nil
	}

	tmpl := template.New("").Funcs(BaseFuncMap())

	var err error
	tmpl, err = tmpl.ParseFS(
		themeFS,
		"themes/common/templates/partials/landing-*.html",
		"themes/common/templates/faq/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse FAQ templates: %w", err)
	}

	return &FaqGenerator{
		site:      site,
		templates: tmpl,
		themeFS:   themeFS,
		version:   version,
	}, nil
}

// Generate generates the FAQ page
func (g *FaqGenerator) Generate() error {
	if g.site.FaqPage == nil || !g.site.Config.Faq.Enabled {
		return nil
	}

	fmt.Println("Generating FAQ page...")

	faqPath := g.site.Config.Faq.Path
	if faqPath == "" {
		faqPath = "faq"
	}

	outputDir := filepath.Join(g.site.OutputRoot, faqPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create FAQ output directory: %w", err)
	}

	// Generate main FAQ page
	if err := g.generateMainPage(outputDir); err != nil {
		return fmt.Errorf("failed to generate FAQ page: %w", err)
	}

	// Count total questions
	totalQuestions := 0
	for _, cat := range g.site.FaqPage.Categories {
		totalQuestions += len(cat.Items)
	}

	fmt.Printf("Generated FAQ page: %d categories, %d questions\n",
		len(g.site.FaqPage.Categories), totalQuestions)
	return nil
}

// generateMainPage generates the FAQ page
func (g *FaqGenerator) generateMainPage(outputDir string) error {
	faqPath := g.site.Config.Faq.Path
	if faqPath == "" {
		faqPath = "faq"
	}
	data := map[string]any{
		"Site":       g.site,
		"FaqPage":    g.site.FaqPage,
		"Footer":     g.buildFooterWithLegal(),
		"BasePath":   g.getBasePath(),
		"Version":    g.version,
		"PageTitle":  g.site.FaqPage.Config.Title + " | " + g.site.Config.Title,
		"ActivePath": "/" + faqPath + "/",
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "faq.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	outputPath := filepath.Join(outputDir, "index.html")
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write FAQ page: %w", err)
	}

	return nil
}

// getBasePath extracts the path component from BaseURL
func (g *FaqGenerator) getBasePath() string {
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
func (g *FaqGenerator) buildFooterWithLegal() core.FooterConfig {
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
