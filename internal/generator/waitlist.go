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

// WaitlistGenerator generates the waitlist landing page HTML
type WaitlistGenerator struct {
	site      *core.Site
	templates *template.Template
	themeFS   embed.FS
	version   string
}

// NewWaitlistGenerator creates a new waitlist generator
func NewWaitlistGenerator(site *core.Site, themeFS embed.FS, version string) (*WaitlistGenerator, error) {
	if !site.Config.Waitlist.Enabled {
		return nil, nil
	}

	tmpl := template.New("").Funcs(BaseFuncMap())

	var err error
	tmpl, err = tmpl.ParseFS(
		themeFS,
		"themes/common/templates/waitlist/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse waitlist templates: %w", err)
	}

	return &WaitlistGenerator{
		site:      site,
		templates: tmpl,
		themeFS:   themeFS,
		version:   version,
	}, nil
}

// Generate generates the waitlist page as index.html
func (g *WaitlistGenerator) Generate() error {
	if g.site.WaitlistPage == nil || !g.site.Config.Waitlist.Enabled {
		return nil
	}

	fmt.Println("Generating waitlist page...")

	data := map[string]any{
		"Site":     g.site,
		"Config":   g.site.WaitlistPage.Config,
		"BasePath": g.getBasePath(),
		"Version":  g.version,
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "waitlist.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	outputPath := filepath.Join(g.site.OutputRoot, "index.html")
	if err := writeWebFile(outputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write waitlist page: %w", err)
	}

	fmt.Println("Generated waitlist page")
	return nil
}

// getBasePath extracts the path component from BaseURL
func (g *WaitlistGenerator) getBasePath() string {
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
