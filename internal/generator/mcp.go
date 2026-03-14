package generator

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// MCPGenerator generates static HTML from parsed MCP server manifests
type MCPGenerator struct {
	site      *core.Site
	templates *template.Template
	themeFS   embed.FS
	version   string
}

// MCPFuncMap returns template functions specific to MCP documentation
func MCPFuncMap() template.FuncMap {
	return template.FuncMap{
		"mcpIsRequired": func(name string, required []string) bool {
			for _, r := range required {
				if r == name {
					return true
				}
			}
			return false
		},
		"mcpSlug": func(name string) string {
			re := regexp.MustCompile(`[^a-zA-Z0-9-]`)
			slug := re.ReplaceAllString(strings.ToLower(name), "-")
			// Collapse multiple dashes
			multi := regexp.MustCompile(`-+`)
			return strings.Trim(multi.ReplaceAllString(slug, "-"), "-")
		},
	}
}

// NewMCPGenerator creates a new MCP generator. Returns nil if MCP is disabled or no specs exist.
func NewMCPGenerator(site *core.Site, themeFS embed.FS, version string) (*MCPGenerator, error) {
	if !site.Config.MCP.Enabled {
		return nil, nil
	}

	tmpl := template.New("").
		Funcs(BaseFuncMap()).
		Funcs(MCPFuncMap()).
		Funcs(AnalyticsFuncMap())

	var err error
	tmpl, err = tmpl.ParseFS(themeFS, "themes/common/templates/mcp-page.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to parse MCP templates: %v\n", err)
		return nil, nil
	}

	return &MCPGenerator{
		site:      site,
		templates: tmpl,
		themeFS:   themeFS,
		version:   version,
	}, nil
}

// Generate generates HTML files for all MCP specs
func (g *MCPGenerator) Generate() error {
	if g == nil || !g.site.Config.MCP.Enabled {
		return nil
	}

	if len(g.site.MCPSpecs) == 0 {
		return nil
	}

	fmt.Println("Generating MCP server documentation...")

	mcpDir := filepath.Join(g.site.OutputRoot, g.site.Config.MCP.Path)
	if err := makeWebDir(mcpDir); err != nil {
		return fmt.Errorf("create MCP output directory: %w", err)
	}

	basePath := g.getBasePath()

	for _, spec := range g.site.MCPSpecs {
		if err := g.generateSpec(spec, mcpDir, basePath); err != nil {
			return fmt.Errorf("generate MCP spec %s: %w", spec.Name, err)
		}
	}

	if err := g.generateIndex(mcpDir, basePath); err != nil {
		return fmt.Errorf("generate MCP index: %w", err)
	}

	fmt.Printf("Generated %d MCP spec(s)\n", len(g.site.MCPSpecs))
	return nil
}

// generateSpec renders the HTML page for a single MCP spec
func (g *MCPGenerator) generateSpec(spec *core.MCPSpec, mcpDir, basePath string) error {
	slug := specSlug(spec.Name)
	specDir := filepath.Join(mcpDir, slug)
	if err := makeWebDir(specDir); err != nil {
		return fmt.Errorf("create spec directory: %w", err)
	}

	outputPath := filepath.Join(specDir, "index.html")
	spec.OutputPath = outputPath

	data := map[string]any{
		"Site":     g.site,
		"Spec":     spec,
		"BasePath": basePath,
		"Version":  g.version,
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "mcp-page.html", data); err != nil {
		return fmt.Errorf("template execution: %w", err)
	}

	if err := writeWebFile(outputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("write HTML: %w", err)
	}

	return nil
}

// generateIndex writes either a redirect (single spec) or a listing page (multiple specs)
func (g *MCPGenerator) generateIndex(mcpDir, basePath string) error {
	indexPath := filepath.Join(mcpDir, "index.html")

	if len(g.site.MCPSpecs) == 1 {
		spec := g.site.MCPSpecs[0]
		slug := specSlug(spec.Name)
		redirect := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="refresh" content="0; url=%s/%s/%s/">
<link rel="canonical" href="%s/%s/%s/">
</head>
<body><p>Redirecting to <a href="%s/%s/%s/">MCP Documentation</a>...</p></body>
</html>`,
			basePath, g.site.Config.MCP.Path, slug,
			basePath, g.site.Config.MCP.Path, slug,
			basePath, g.site.Config.MCP.Path, slug)
		return writeWebFile(indexPath, []byte(redirect))
	}

	// Multiple specs — simple listing page
	var buf strings.Builder
	buf.WriteString(`<!DOCTYPE html>
<html lang="en" data-theme="light">
<head>
<script>
(function(){var t=localStorage.getItem('theme');if(t==='dark'||(t!=='light'&&window.matchMedia('(prefers-color-scheme:dark)').matches)){document.documentElement.setAttribute('data-theme','dark');}})();
</script>
<style>
:root{--bg-primary:#fafafa;--text-primary:#1a1a1a;--link-color:#2563eb}
:root[data-theme="dark"]{--bg-primary:#1a1a1a;--text-primary:#fff;--link-color:#7bb3ff}
html,body{background-color:var(--bg-primary);color:var(--text-primary);font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Arial,sans-serif;margin:0;padding:2rem}
a{color:var(--link-color)}
</style>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>MCP Documentation | ` + g.site.Config.Title + `</title>
</head>
<body>
<h1>MCP Server Documentation</h1>
<p><a href="` + basePath + `/">← Back to Docs</a></p>
<ul>
`)
	for _, spec := range g.site.MCPSpecs {
		slug := specSlug(spec.Name)
		href := fmt.Sprintf("%s/%s/%s/", basePath, g.site.Config.MCP.Path, slug)
		buf.WriteString(fmt.Sprintf(`<li><a href="%s">%s</a> — %d tools</li>`,
			href, spec.Name, len(spec.Tools)))
	}
	buf.WriteString(`</ul>
<script defer src="` + basePath + `/js/theme-toggle.js?v=` + g.version + `"></script>
</body></html>`)

	return writeWebFile(indexPath, []byte(buf.String()))
}

// getBasePath extracts the path component from BaseURL
func (g *MCPGenerator) getBasePath() string {
	baseURL := g.site.Config.BaseURL
	if baseURL == "" {
		return ""
	}
	for _, prefix := range []string{"https://", "http://"} {
		if strings.HasPrefix(baseURL, prefix) {
			baseURL = strings.TrimPrefix(baseURL, prefix)
			break
		}
	}
	parts := strings.SplitN(baseURL, "/", 2)
	if len(parts) < 2 {
		return ""
	}
	path := "/" + strings.TrimSuffix(parts[1], "/")
	if path == "/" {
		return ""
	}
	return path
}

// specSlug converts an MCP spec name to a URL-safe slug
func specSlug(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9-]`)
	slug := re.ReplaceAllString(strings.ToLower(name), "-")
	multi := regexp.MustCompile(`-+`)
	return strings.Trim(multi.ReplaceAllString(slug, "-"), "-")
}
