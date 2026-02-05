package generator

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
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
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"json": func(v interface{}) (template.JS, error) {
			bytes, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return template.JS(bytes), nil
		},
		"lower": strings.ToLower,
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
	if err := os.MkdirAll(apiDir, 0755); err != nil {
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
	if err := os.MkdirAll(specDir, 0755); err != nil {
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
	data := map[string]interface{}{
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
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
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
	specData := map[string]interface{}{
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

	if err := os.WriteFile(filepath.Join(specDir, "spec-data.json"), specDataJSON, 0644); err != nil {
		return fmt.Errorf("failed to write spec data JSON: %w", err)
	}

	// Generate endpoints file (single file, JS handles organization)
	endpointsJSON, err := json.Marshal(spec.Endpoints)
	if err != nil {
		return fmt.Errorf("failed to marshal endpoints: %w", err)
	}

	if err := os.WriteFile(filepath.Join(specDir, "endpoints.json"), endpointsJSON, 0644); err != nil {
		return fmt.Errorf("failed to write endpoints JSON: %w", err)
	}

	return nil
}

// generateAPIIndex generates an index page listing all OpenAPI specifications
func (g *OpenAPIGenerator) generateAPIIndex(apiDir string) error {
	var buf bytes.Buffer

	// Generate simple index HTML
	buf.WriteString(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>API Documentation - ` + g.site.Config.Title + `</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" href="../css/style.css">
</head>
<body>
	<div class="container">
		<h1>API Documentation</h1>
		<p>Available OpenAPI specifications:</p>
		<ul>
`)

	for _, spec := range g.site.APISpecs {
		specName := spec.Name
		// Remove common extensions
		specName = strings.TrimSuffix(specName, ".yaml")
		specName = strings.TrimSuffix(specName, ".yml")
		specName = strings.TrimSuffix(specName, ".json")

		relPath := filepath.Join("/api", specName)
		buf.WriteString(fmt.Sprintf(`			<li>
				<h3><a href="%s">%s</a></h3>
				<p>%s</p>
				<p><small>Version %s | %d endpoints</small></p>
			</li>
`, relPath, spec.Title, spec.Description, spec.Version, len(spec.Endpoints)))
	}

	buf.WriteString(`		</ul>
	</div>
</body>
</html>`)

	// Write index file
	indexPath := filepath.Join(apiDir, "index.html")
	if err := os.WriteFile(indexPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write API index: %w", err)
	}

	return nil
}

// getBasePath extracts the path component from BaseURL
func (g *OpenAPIGenerator) getBasePath() string {
	baseURL := g.site.Config.BaseURL
	if baseURL == "" {
		return ""
	}

	// Parse the URL to extract the path
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

	// Get the path part
	path := "/" + parts[1]
	path = strings.TrimSuffix(path, "/")

	if path == "/" {
		return ""
	}

	return path
}
