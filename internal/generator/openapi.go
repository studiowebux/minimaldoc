package generator

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// OpenAPIGenerator generates static HTML and JSON for OpenAPI specifications
type OpenAPIGenerator struct {
	site      *core.Site
	templates *template.Template
	themeFS   embed.FS
	version   string
}

// NewOpenAPIGenerator creates a new OpenAPI generator
func NewOpenAPIGenerator(site *core.Site, themeFS embed.FS, version string) (*OpenAPIGenerator, error) {
	if !site.Config.OpenAPI.Enabled {
		return nil, nil // Skip if OpenAPI is not enabled
	}

	// Create template with shared OpenAPI functions
	tmpl := template.New("").Funcs(OpenAPIFuncMap()).Funcs(AnalyticsFuncMap())

	// Parse OpenAPI templates from common (all structure is shared, themes only provide CSS)
	tmpl, err := tmpl.ParseFS(
		themeFS,
		"themes/common/templates/openapi-page.html",
		"themes/common/templates/partials/openapi-nav.html",
	)
	if err != nil {
		// Templates don't exist, return nil generator
		fmt.Fprintf(os.Stderr, "Warning: failed to parse OpenAPI templates: %v\n", err)
		return nil, nil
	}

	return &OpenAPIGenerator{
		site:      site,
		templates: tmpl,
		themeFS:   themeFS,
		version:   version,
	}, nil
}

// Generate generates HTML and JSON files for all OpenAPI specifications
func (g *OpenAPIGenerator) Generate() error {
	if g == nil || !g.site.Config.OpenAPI.Enabled {
		return nil // Skip if disabled
	}

	if len(g.site.APISpecs) == 0 {
		return nil // Nothing to generate
	}

	fmt.Println("Generating OpenAPI documentation...")

	// Create API directory
	apiDir := filepath.Join(g.site.OutputRoot, "api")
	if err := makeWebDir(apiDir); err != nil {
		return fmt.Errorf("failed to create API directory: %w", err)
	}

	// Generate each spec
	for _, spec := range g.site.APISpecs {
		if err := g.generateSpec(spec, apiDir); err != nil {
			return fmt.Errorf("failed to generate spec %s: %w", spec.Name, err)
		}
	}

	// Generate API index page (listing all specs)
	if err := g.generateAPIIndex(apiDir); err != nil {
		return fmt.Errorf("failed to generate API index: %w", err)
	}

	fmt.Printf("Generated %d OpenAPI specifications\n", len(g.site.APISpecs))
	return nil
}

// generateSpec generates HTML and JSON for a single OpenAPI specification
func (g *OpenAPIGenerator) generateSpec(spec *core.APISpec, apiDir string) error {
	// Create spec directory (strip file extension from name)
	specName := spec.Name
	// Remove common extensions
	specName = strings.TrimSuffix(specName, ".yaml")
	specName = strings.TrimSuffix(specName, ".yml")
	specName = strings.TrimSuffix(specName, ".json")

	specDir := filepath.Join(apiDir, specName)
	if err := makeWebDir(specDir); err != nil {
		return fmt.Errorf("failed to create spec directory: %w", err)
	}

	// Set output path
	spec.OutputPath = filepath.Join(specDir, "index.html")

	// Generate HTML page
	if err := g.generateSpecHTML(spec, specDir); err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Generate JSON data files for lazy loading
	if err := g.generateSpecJSON(spec, specDir); err != nil {
		return fmt.Errorf("failed to generate JSON: %w", err)
	}

	return nil
}

// generateSpecHTML generates the HTML page for an OpenAPI spec
func (g *OpenAPIGenerator) generateSpecHTML(spec *core.APISpec, specDir string) error {
	// Prepare template data
	data := map[string]any{
		"Site":              g.site,
		"Spec":              spec,
		"BasePath":          g.getBasePath(),
		"DefaultView":       g.site.Config.OpenAPI.DefaultView,
		"EnableCodeSamples": g.site.Config.OpenAPI.EnableCodeSamples,
		"Version":           g.version,
	}

	var buf bytes.Buffer

	// If we have templates, use them
	if g.templates != nil {
		if err := g.templates.ExecuteTemplate(&buf, "openapi-page.html", data); err != nil {
			return fmt.Errorf("template execution failed: %w", err)
		}
	} else {
		// Fallback: generate basic HTML
		buf.WriteString(g.generateFallbackHTML(spec))
	}

	// Write HTML file
	outputPath := filepath.Join(specDir, "index.html")
	if err := writeWebFile(outputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write HTML file: %w", err)
	}

	return nil
}

// generateFallbackHTML generates a basic HTML fallback when templates don't exist
func (g *OpenAPIGenerator) generateFallbackHTML(spec *core.APISpec) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>%s - %s</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
	<h1>%s</h1>
	<p>%s</p>
	<p>Version: %s</p>
	<p><em>OpenAPI templates not yet implemented. See generated JSON data.</em></p>
	<script>
		// Load spec data
		fetch('./spec-data.json')
			.then(r => r.json())
			.then(data => {
				console.log('OpenAPI spec loaded:', data);
			});
	</script>
</body>
</html>`, spec.Title, g.site.Config.Title, spec.Title, spec.Description, spec.Version)
}

// generateSpecJSON generates JSON data files for an OpenAPI spec
func (g *OpenAPIGenerator) generateSpecJSON(spec *core.APISpec, specDir string) error {
	// Generate main spec data file with metadata and schemas (schemas stored once)
	specData := map[string]any{
		"name":            spec.Name,
		"title":           spec.Title,
		"description":     spec.Description,
		"version":         spec.Version,
		"openapi":         spec.OpenAPIVersion,
		"servers":         spec.Servers,
		"tags":            spec.Tags,
		"securitySchemes": spec.SecuritySchemes,
		"schemas":         spec.Schemas,
		"endpointCount":   len(spec.Endpoints),
		"schemaCount":     len(spec.Schemas),
	}

	specDataJSON, err := json.Marshal(specData)
	if err != nil {
		return fmt.Errorf("failed to marshal spec data: %w", err)
	}

	if err := writeWebFile(filepath.Join(specDir, "spec-data.json"), specDataJSON); err != nil {
		return fmt.Errorf("failed to write spec data JSON: %w", err)
	}

	// Generate endpoints file (single file, JS handles organization)
	endpointsJSON, err := json.Marshal(spec.Endpoints)
	if err != nil {
		return fmt.Errorf("failed to marshal endpoints: %w", err)
	}

	if err := writeWebFile(filepath.Join(specDir, "endpoints.json"), endpointsJSON); err != nil {
		return fmt.Errorf("failed to write endpoints JSON: %w", err)
	}

	return nil
}

// generateAPIIndex generates an index page listing all OpenAPI specifications
func (g *OpenAPIGenerator) generateAPIIndex(apiDir string) error {
	basePath := g.getBasePath()

	// If only one API spec, redirect directly to it
	if len(g.site.APISpecs) == 1 {
		spec := g.site.APISpecs[0]
		// Strip extension from name (same as generateSpec does)
		specSlug := strings.TrimSuffix(spec.Name, ".yaml")
		specSlug = strings.TrimSuffix(specSlug, ".yml")
		specSlug = strings.TrimSuffix(specSlug, ".json")
		safeBase := html.EscapeString(basePath)
		safeSlug := html.EscapeString(specSlug)
		redirectHTML := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="refresh" content="0; url=%s/api/%s/">
<link rel="canonical" href="%s/api/%s/">
</head>
<body>
<p>Redirecting to <a href="%s/api/%s/">API Documentation</a>...</p>
</body>
</html>`, safeBase, safeSlug, safeBase, safeSlug, safeBase, safeSlug)

		indexPath := filepath.Join(apiDir, "index.html")
		if err := writeWebFile(indexPath, []byte(redirectHTML)); err != nil {
			return fmt.Errorf("failed to write API redirect: %w", err)
		}
		return nil
	}

	var buf bytes.Buffer

	// Build navigation links
	var navLinks strings.Builder
	for _, nav := range g.site.Config.Landing.Nav {
		navLinks.WriteString(fmt.Sprintf(`<a href="%s">%s</a>`, nav.URL, nav.Text))
	}
	if g.site.Config.OpenAPI.Enabled {
		navLinks.WriteString(fmt.Sprintf(`<a href="%s/api/" class="active">API</a>`, basePath))
	}
	if g.site.Config.Portfolio.Enabled && g.site.Config.Portfolio.Path != "" {
		navLinks.WriteString(fmt.Sprintf(`<a href="%s/%s/">Portfolio</a>`, basePath, g.site.Config.Portfolio.Path))
	}
	if g.site.Config.Faq.Enabled && g.site.Config.Faq.Path != "" {
		navLinks.WriteString(fmt.Sprintf(`<a href="%s/%s/">FAQ</a>`, basePath, g.site.Config.Faq.Path))
	}
	if g.site.Config.Contact.Enabled && g.site.Config.Contact.Path != "" {
		navLinks.WriteString(fmt.Sprintf(`<a href="%s/%s/">Contact</a>`, basePath, g.site.Config.Contact.Path))
	}

	// Build footer links with legal pages
	footer := BuildFooter(g.site, g.version)
	var footerLinks strings.Builder
	if len(footer.Links) > 0 {
		footerLinks.WriteString(`<div class="footer-links">`)
		for _, group := range footer.Links {
			footerLinks.WriteString(fmt.Sprintf(`<div class="footer-link-group"><h4 class="footer-group-title">%s</h4><ul class="footer-group-list">`, group.Title))
			for _, item := range group.Items {
				footerLinks.WriteString(fmt.Sprintf(`<li><a href="%s">%s</a></li>`, item.URL, item.Text))
			}
			footerLinks.WriteString(`</ul></div>`)
		}
		footerLinks.WriteString(`</div>`)
	}

	// Generate index HTML with proper styling
	buf.WriteString(`<!DOCTYPE html>
<html lang="en" data-theme="light" data-base-path="` + html.EscapeString(basePath) + `">
<head>
<script>
(function(){var t=localStorage.getItem('theme');if(t==='dark'||(t!=='light'&&window.matchMedia('(prefers-color-scheme:dark)').matches)){document.documentElement.setAttribute('data-theme','dark');}})();
</script>
<style>
:root{--bg-primary:#fafafa;--text-primary:#1a1a1a;--text-secondary:#4a4a4a;--link-color:#2563eb}
:root[data-theme="dark"]{--bg-primary:#1a1a1a;--text-primary:#fff;--text-secondary:#e8e8e8;--link-color:#7bb3ff}
html,body{background-color:var(--bg-primary);color:var(--text-primary)}
a{color:var(--link-color)}
</style>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=5.0">
<title>API Documentation | ` + html.EscapeString(g.site.Config.Title) + `</title>
<meta name="description" content="OpenAPI documentation for ` + html.EscapeString(g.site.Config.Title) + `">
<link rel="stylesheet" href="` + html.EscapeString(basePath) + `/css/tokens.css?v=` + html.EscapeString(g.version) + `">
<link rel="stylesheet" href="` + html.EscapeString(basePath) + `/css/base.css?v=` + html.EscapeString(g.version) + `">
<link rel="stylesheet" href="` + html.EscapeString(basePath) + `/css/main.css?v=` + html.EscapeString(g.version) + `">
<link rel="stylesheet" href="` + html.EscapeString(basePath) + `/css/landing.css?v=` + html.EscapeString(g.version) + `">
<link rel="stylesheet" href="` + html.EscapeString(basePath) + `/css/portfolio.css?v=` + html.EscapeString(g.version) + `">
</head>
<body>
<a href="#main" class="skip-link">Skip to main content</a>
<div class="landing-page">
<header class="landing-header">
<div class="landing-header-content">
<a href="` + html.EscapeString(basePath) + `/" class="landing-logo">` + html.EscapeString(g.site.Config.Title) + `</a>
<button class="landing-menu-toggle" id="landing-menu-toggle" aria-label="Toggle menu" aria-expanded="false">
<span></span>
</button>
<div class="landing-nav-backdrop" id="landing-nav-backdrop"></div>
<nav class="landing-nav" id="landing-nav">
` + navLinks.String() + `
<button id="theme-toggle" class="theme-toggle" aria-label="Toggle theme">
<span class="theme-icon" aria-hidden="true"></span>
</button>
</nav>
</div>
</header>
<main id="main" class="landing-main">
<section class="portfolio-header">
<h1 style="font-size:2.5rem">API Documentation</h1>
</section>
<section class="portfolio-grid">
`)

	for _, spec := range g.site.APISpecs {
		specName := spec.Name
		// Remove common extensions
		specName = strings.TrimSuffix(specName, ".yaml")
		specName = strings.TrimSuffix(specName, ".yml")
		specName = strings.TrimSuffix(specName, ".json")

		// Strip HTML and truncate description for the card
		desc := stripHTML(spec.Description)
		if len(desc) > 150 {
			desc = desc[:147] + "..."
		}
		desc = template.HTMLEscapeString(desc)

		relPath := html.EscapeString(basePath + "/api/" + specName)
		buf.WriteString(fmt.Sprintf(`<a href="%s" class="project-card" style="text-decoration:none;color:inherit;cursor:pointer;display:block">
<div class="project-content">
<h2 class="project-title">%s</h2>
<p class="project-description">%s</p>
<div class="project-tags">
<span class="project-tag">v%s</span>
<span class="project-tag">%d endpoints</span>
</div>
</div>
</a>
`, relPath, html.EscapeString(spec.Title), desc, html.EscapeString(spec.Version), len(spec.Endpoints)))
	}

	buf.WriteString(`</section>
</main>
<footer class="landing-footer">
<div class="footer-content">
` + footerLinks.String() + `
<p class="footer-copyright">` + html.EscapeString(footer.Copyright) + `</p>
</div>
</footer>
</div>
<script defer src="` + html.EscapeString(basePath) + `/js/theme-toggle.js?v=` + html.EscapeString(g.version) + `"></script>
<script>
(function(){
    var toggle = document.getElementById('landing-menu-toggle');
    var nav = document.getElementById('landing-nav');
    var backdrop = document.getElementById('landing-nav-backdrop');
    if (!toggle || !nav || !backdrop) return;
    function open() {
        nav.classList.add('open'); backdrop.classList.add('open'); toggle.classList.add('open');
        toggle.setAttribute('aria-expanded', 'true'); document.body.style.overflow = 'hidden';
    }
    function close() {
        nav.classList.remove('open'); backdrop.classList.remove('open'); toggle.classList.remove('open');
        toggle.setAttribute('aria-expanded', 'false'); document.body.style.overflow = '';
    }
    toggle.addEventListener('click', function() { nav.classList.contains('open') ? close() : open(); });
    backdrop.addEventListener('click', close);
    nav.querySelectorAll('a').forEach(function(link) { link.addEventListener('click', close); });
    window.addEventListener('resize', function() { if (window.innerWidth > 768 && nav.classList.contains('open')) close(); });
})();
</script>
</body>
</html>`)

	// Write index file
	indexPath := filepath.Join(apiDir, "index.html")
	if err := writeWebFile(indexPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write API index: %w", err)
	}

	return nil
}

// getBasePath extracts the path component from BaseURL
func (g *OpenAPIGenerator) getBasePath() string {
	return GetBasePath(g.site.Config.BaseURL)
}
