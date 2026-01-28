package parser

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/studiowebux/minimaldoc/internal/core"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	htmlrenderer "github.com/yuin/goldmark/renderer/html"
)

// OpenAPIParser handles parsing of OpenAPI specifications
type OpenAPIParser struct {
	cacheDir string
	doc      *openapi3.T        // Current document being parsed (for resolving $refs)
	md       goldmark.Markdown  // Markdown renderer for descriptions
}

// NewOpenAPIParser creates a new OpenAPI parser
func NewOpenAPIParser(cacheDir string) *OpenAPIParser {
	// Create a simple markdown renderer for descriptions
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,          // GitHub Flavored Markdown
			extension.Typographer,  // Smart quotes, dashes
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			htmlrenderer.WithUnsafe(), // Allow raw HTML
		),
	)

	return &OpenAPIParser{
		cacheDir: cacheDir,
		md:       md,
	}
}

// ParseFile parses an OpenAPI spec from a local file
func (p *OpenAPIParser) ParseFile(filePath string) (*core.APISpec, error) {
	// Load the spec
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	doc, err := loader.LoadFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI spec from %s: %w", filePath, err)
	}

	// Validate the spec
	ctx := context.Background()
	if err := doc.Validate(ctx); err != nil {
		return nil, fmt.Errorf("invalid OpenAPI spec in %s: %w", filePath, err)
	}

	// Store document for reference resolution
	p.doc = doc

	// Convert to our internal structure
	spec := p.convertSpec(doc)
	spec.Name = filepath.Base(filePath)
	spec.FilePath = filePath

	return spec, nil
}

// ParseURL fetches and parses an OpenAPI spec from a URL
func (p *OpenAPIParser) ParseURL(url string) (*core.APISpec, error) {
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(p.cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Fetch the spec
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OpenAPI spec from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch OpenAPI spec from %s: status %d", url, resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI spec from %s: %w", url, err)
	}

	// Parse the spec
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	doc, err := loader.LoadFromData(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec from %s: %w", url, err)
	}

	// Validate the spec
	ctx := context.Background()
	if err := doc.Validate(ctx); err != nil {
		return nil, fmt.Errorf("invalid OpenAPI spec from %s: %w", url, err)
	}

	// Store document for reference resolution
	p.doc = doc

	// Convert to our internal structure
	spec := p.convertSpec(doc)
	spec.Name = p.NameFromURL(url)
	spec.URL = url
	spec.LastFetched = time.Now()
	spec.ETag = resp.Header.Get("ETag")
	spec.LastModified = resp.Header.Get("Last-Modified")

	// Cache the spec
	cachePath := filepath.Join(p.cacheDir, spec.Name+".json")
	if err := os.WriteFile(cachePath, body, 0644); err != nil {
		// Non-fatal: continue even if caching fails
		fmt.Fprintf(os.Stderr, "Warning: failed to cache OpenAPI spec: %v\n", err)
	}

	return spec, nil
}

// convertSpec converts an openapi3.T to our internal APISpec
func (p *OpenAPIParser) convertSpec(doc *openapi3.T) *core.APISpec {
	spec := &core.APISpec{
		Title:          doc.Info.Title,
		Description:    p.renderMarkdown(doc.Info.Description),
		Version:        doc.Info.Version,
		OpenAPIVersion: doc.OpenAPI,
		Servers:        p.convertServers(doc.Servers),
		Tags:           p.convertTags(doc.Tags),
		SecuritySchemes: p.convertSecuritySchemes(doc.Components.SecuritySchemes),
		Schemas:        p.convertSchemas(doc.Components.Schemas),
		Endpoints:      []*core.APIEndpoint{},
	}

	// Parse all paths
	for path, pathItem := range doc.Paths.Map() {
		spec.Endpoints = append(spec.Endpoints, p.convertPathItem(path, pathItem, spec.Name)...)
	}

	// Organize endpoints
	spec.OrganizedByPath = p.organizeByPath(spec.Endpoints)
	spec.OrganizedByTag = p.organizeByTag(spec.Endpoints, spec.Tags)
	spec.FlatEndpoints = spec.Endpoints

	return spec
}

// convertSchemas converts OpenAPI component schemas
func (p *OpenAPIParser) convertSchemas(schemas openapi3.Schemas) map[string]*core.APISchema {
	if schemas == nil {
		return nil
	}

	result := make(map[string]*core.APISchema)
	for name, schemaRef := range schemas {
		if schemaRef == nil {
			continue
		}
		result[name] = p.convertSchema(schemaRef)
	}
	return result
}

// convertServers converts OpenAPI servers
func (p *OpenAPIParser) convertServers(servers openapi3.Servers) []core.APIServer {
	result := make([]core.APIServer, 0, len(servers))
	for _, server := range servers {
		s := core.APIServer{
			URL:         server.URL,
			Description: p.renderMarkdown(server.Description),
			Variables:   make(map[string]core.APIServerVariable),
		}

		for name, variable := range server.Variables {
			s.Variables[name] = core.APIServerVariable{
				Default:     variable.Default,
				Description: p.renderMarkdown(variable.Description),
				Enum:        variable.Enum,
			}
		}

		result = append(result, s)
	}
	return result
}

// convertTags converts OpenAPI tags
func (p *OpenAPIParser) convertTags(tags openapi3.Tags) []core.APITag {
	result := make([]core.APITag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, core.APITag{
			Name:        tag.Name,
			Description: p.renderMarkdown(tag.Description),
		})
	}
	return result
}

// convertSecuritySchemes converts OpenAPI security schemes
func (p *OpenAPIParser) convertSecuritySchemes(schemes openapi3.SecuritySchemes) map[string]*core.APISecurityScheme {
	result := make(map[string]*core.APISecurityScheme)
	for name, schemeRef := range schemes {
		if schemeRef == nil || schemeRef.Value == nil {
			continue
		}
		scheme := schemeRef.Value

		s := &core.APISecurityScheme{
			Type:             scheme.Type,
			Description:      p.renderMarkdown(scheme.Description),
			Name:             scheme.Name,
			In:               scheme.In,
			Scheme:           scheme.Scheme,
			BearerFormat:     scheme.BearerFormat,
			OpenIDConnectURL: scheme.OpenIdConnectUrl,
		}

		// Convert OAuth flows
		if scheme.Flows != nil {
			s.Flows = &core.APIOAuthFlows{}

			if scheme.Flows.Implicit != nil {
				s.Flows.Implicit = p.convertOAuthFlow(scheme.Flows.Implicit)
			}
			if scheme.Flows.Password != nil {
				s.Flows.Password = p.convertOAuthFlow(scheme.Flows.Password)
			}
			if scheme.Flows.ClientCredentials != nil {
				s.Flows.ClientCredentials = p.convertOAuthFlow(scheme.Flows.ClientCredentials)
			}
			if scheme.Flows.AuthorizationCode != nil {
				s.Flows.AuthorizationCode = p.convertOAuthFlow(scheme.Flows.AuthorizationCode)
			}
		}

		result[name] = s
	}
	return result
}

// convertOAuthFlow converts an OAuth flow
func (p *OpenAPIParser) convertOAuthFlow(flow *openapi3.OAuthFlow) *core.APIOAuthFlow {
	return &core.APIOAuthFlow{
		AuthorizationURL: flow.AuthorizationURL,
		TokenURL:         flow.TokenURL,
		RefreshURL:       flow.RefreshURL,
		Scopes:           flow.Scopes,
	}
}

// convertPathItem converts a path item to endpoints
func (p *OpenAPIParser) convertPathItem(path string, pathItem *openapi3.PathItem, specName string) []*core.APIEndpoint {
	endpoints := []*core.APIEndpoint{}

	operations := map[string]*openapi3.Operation{
		"GET":     pathItem.Get,
		"POST":    pathItem.Post,
		"PUT":     pathItem.Put,
		"PATCH":   pathItem.Patch,
		"DELETE":  pathItem.Delete,
		"HEAD":    pathItem.Head,
		"OPTIONS": pathItem.Options,
		"TRACE":   pathItem.Trace,
	}

	for method, operation := range operations {
		if operation == nil {
			continue
		}

		endpoint := &core.APIEndpoint{
			OperationID: operation.OperationID,
			Path:        path,
			Method:      method,
			Summary:     p.renderMarkdown(operation.Summary),
			Description: p.renderMarkdown(operation.Description),
			Tags:        operation.Tags,
			Deprecated:  operation.Deprecated,
			Parameters:  p.convertParameters(operation.Parameters),
			RequestBody: p.convertRequestBody(operation.RequestBody),
			Responses:   p.convertResponses(operation.Responses),
			Security:    p.convertSecurity(operation.Security),
			SpecName:    specName,
		}

		// Add path-level parameters
		endpoint.Parameters = append(endpoint.Parameters, p.convertParameters(pathItem.Parameters)...)

		endpoints = append(endpoints, endpoint)
	}

	return endpoints
}

// convertParameters converts OpenAPI parameters
func (p *OpenAPIParser) convertParameters(params openapi3.Parameters) []*core.APIParameter {
	result := make([]*core.APIParameter, 0, len(params))
	for _, paramRef := range params {
		if paramRef == nil || paramRef.Value == nil {
			continue
		}
		param := paramRef.Value

		apiParam := &core.APIParameter{
			Name:        param.Name,
			In:          param.In,
			Description: p.renderMarkdown(param.Description),
			Required:    param.Required,
			Deprecated:  param.Deprecated,
			Schema:      p.convertSchema(param.Schema),
			Example:     param.Example,
		}

		// Convert examples
		if len(param.Examples) > 0 {
			apiParam.Examples = make(map[string]*core.APIExample)
			for name, exampleRef := range param.Examples {
				if exampleRef != nil && exampleRef.Value != nil {
					apiParam.Examples[name] = &core.APIExample{
						Summary:     p.renderMarkdown(exampleRef.Value.Summary),
						Description: p.renderMarkdown(exampleRef.Value.Description),
						Value:       exampleRef.Value.Value,
					}
				}
			}
		}

		result = append(result, apiParam)
	}
	return result
}

// convertRequestBody converts an OpenAPI request body
func (p *OpenAPIParser) convertRequestBody(bodyRef *openapi3.RequestBodyRef) *core.APIRequestBody {
	if bodyRef == nil || bodyRef.Value == nil {
		return nil
	}
	body := bodyRef.Value

	rb := &core.APIRequestBody{
		Description: p.renderMarkdown(body.Description),
		Required:    body.Required,
		Content:     make(map[string]*core.APIMediaType),
	}

	for mediaType, mediaTypeValue := range body.Content {
		rb.Content[mediaType] = p.convertMediaType(mediaTypeValue)
	}

	return rb
}

// convertResponses converts OpenAPI responses
func (p *OpenAPIParser) convertResponses(responses *openapi3.Responses) map[string]*core.APIResponse {
	result := make(map[string]*core.APIResponse)
	if responses == nil {
		return result
	}

	for status, responseRef := range responses.Map() {
		if responseRef == nil || responseRef.Value == nil {
			continue
		}
		response := responseRef.Value

		r := &core.APIResponse{
			Description: p.renderMarkdown(*response.Description),
			Headers:     make(map[string]*core.APIParameter),
			Content:     make(map[string]*core.APIMediaType),
		}

		// Convert headers
		for name, headerRef := range response.Headers {
			if headerRef != nil && headerRef.Value != nil {
				r.Headers[name] = &core.APIParameter{
					Name:        name,
					In:          "header",
					Description: p.renderMarkdown(headerRef.Value.Description),
					Required:    headerRef.Value.Required,
					Schema:      p.convertSchema(headerRef.Value.Schema),
				}
			}
		}

		// Convert content
		for mediaType, mediaTypeValue := range response.Content {
			r.Content[mediaType] = p.convertMediaType(mediaTypeValue)
		}

		result[status] = r
	}

	return result
}

// convertMediaType converts a media type
func (p *OpenAPIParser) convertMediaType(mt *openapi3.MediaType) *core.APIMediaType {
	if mt == nil {
		return nil
	}

	m := &core.APIMediaType{
		Schema:  p.convertSchema(mt.Schema),
		Example: mt.Example,
	}

	// Convert examples
	if len(mt.Examples) > 0 {
		m.Examples = make(map[string]*core.APIExample)
		for name, exampleRef := range mt.Examples {
			if exampleRef != nil && exampleRef.Value != nil {
				m.Examples[name] = &core.APIExample{
					Summary:     exampleRef.Value.Summary,
					Description: exampleRef.Value.Description,
					Value:       exampleRef.Value.Value,
				}
			}
		}
	}

	return m
}

// convertSchema converts an OpenAPI schema
func (p *OpenAPIParser) convertSchema(schemaRef *openapi3.SchemaRef) *core.APISchema {
	if schemaRef == nil {
		return nil
	}

	// If there's a $ref, try to resolve it
	if schemaRef.Ref != "" {
		resolved := p.resolveSchemaRef(schemaRef.Ref)
		if resolved != nil && resolved.Value != nil {
			// Use the resolved schema but keep track of the original ref
			schemaRef = resolved
		}
	}

	// If Value is still nil after resolution attempt, return just the reference
	if schemaRef.Value == nil {
		if schemaRef.Ref != "" {
			return &core.APISchema{
				Ref: schemaRef.Ref,
			}
		}
		return nil
	}
	schema := schemaRef.Value

	// Extract type as string
	typeStr := ""
	if schema.Type != nil && len(*schema.Type) > 0 {
		typeStr = (*schema.Type)[0]
	}

	s := &core.APISchema{
		Type:        typeStr,
		Format:      schema.Format,
		Description: p.renderMarkdown(schema.Description),
		Required:    schema.Required,
		Enum:        schema.Enum,
		Default:     schema.Default,
		Example:     schema.Example,
		Nullable:    schema.Nullable,
		ReadOnly:    schema.ReadOnly,
		WriteOnly:   schema.WriteOnly,
		Deprecated:  schema.Deprecated,
		Pattern:     schema.Pattern,
		UniqueItems: schema.UniqueItems,
	}

	// Numeric constraints
	if schema.Min != nil {
		min := *schema.Min
		s.Minimum = &min
	}
	if schema.Max != nil {
		max := *schema.Max
		s.Maximum = &max
	}

	// String constraints
	if schema.MinLength > 0 {
		ml := int(schema.MinLength)
		s.MinLength = &ml
	}
	if schema.MaxLength != nil {
		ml := int(*schema.MaxLength)
		s.MaxLength = &ml
	}

	// Array constraints
	if schema.MinItems > 0 {
		mi := int(schema.MinItems)
		s.MinItems = &mi
	}
	if schema.MaxItems != nil {
		mi := int(*schema.MaxItems)
		s.MaxItems = &mi
	}

	// Object constraints
	if schema.MinProps > 0 {
		mp := int(schema.MinProps)
		s.MinProperties = &mp
	}
	if schema.MaxProps != nil {
		mp := int(*schema.MaxProps)
		s.MaxProperties = &mp
	}

	// Properties
	if len(schema.Properties) > 0 {
		s.Properties = make(map[string]*core.APISchema)
		for name, propRef := range schema.Properties {
			s.Properties[name] = p.convertSchema(propRef)
		}
	}

	// Items (for arrays)
	if schema.Items != nil {
		s.Items = p.convertSchema(schema.Items)
	}

	// Composition
	if len(schema.AllOf) > 0 {
		s.AllOf = make([]core.APISchema, 0, len(schema.AllOf))
		for _, schemaRef := range schema.AllOf {
			if converted := p.convertSchema(schemaRef); converted != nil {
				s.AllOf = append(s.AllOf, *converted)
			}
		}
	}
	if len(schema.OneOf) > 0 {
		s.OneOf = make([]core.APISchema, 0, len(schema.OneOf))
		for _, schemaRef := range schema.OneOf {
			if converted := p.convertSchema(schemaRef); converted != nil {
				s.OneOf = append(s.OneOf, *converted)
			}
		}
	}
	if len(schema.AnyOf) > 0 {
		s.AnyOf = make([]core.APISchema, 0, len(schema.AnyOf))
		for _, schemaRef := range schema.AnyOf {
			if converted := p.convertSchema(schemaRef); converted != nil {
				s.AnyOf = append(s.AnyOf, *converted)
			}
		}
	}
	if schema.Not != nil {
		s.Not = p.convertSchema(schema.Not)
	}

	return s
}

// resolveSchemaRef resolves a schema reference like "#/components/schemas/Error"
func (p *OpenAPIParser) resolveSchemaRef(ref string) *openapi3.SchemaRef {
	if p.doc == nil || p.doc.Components == nil {
		return nil
	}

	// Handle #/components/schemas/SchemaName format
	if strings.HasPrefix(ref, "#/components/schemas/") {
		schemaName := strings.TrimPrefix(ref, "#/components/schemas/")
		if schemaRef, ok := p.doc.Components.Schemas[schemaName]; ok {
			return schemaRef
		}
	}

	// Handle #/components/responses/ResponseName format
	if strings.HasPrefix(ref, "#/components/responses/") {
		responseName := strings.TrimPrefix(ref, "#/components/responses/")
		if responseRef, ok := p.doc.Components.Responses[responseName]; ok {
			if responseRef.Value != nil && responseRef.Value.Content != nil {
				// Return the first content schema
				for _, mediaType := range responseRef.Value.Content {
					if mediaType.Schema != nil {
						return mediaType.Schema
					}
				}
			}
		}
	}

	// Handle #/components/parameters/ParameterName format
	if strings.HasPrefix(ref, "#/components/parameters/") {
		paramName := strings.TrimPrefix(ref, "#/components/parameters/")
		if paramRef, ok := p.doc.Components.Parameters[paramName]; ok {
			if paramRef.Value != nil && paramRef.Value.Schema != nil {
				return paramRef.Value.Schema
			}
		}
	}

	return nil
}

// renderMarkdown converts markdown description to HTML
func (p *OpenAPIParser) renderMarkdown(markdown string) string {
	if markdown == "" {
		return ""
	}

	var buf bytes.Buffer
	if err := p.md.Convert([]byte(markdown), &buf); err != nil {
		// If markdown rendering fails, return the original text
		return markdown
	}

	// Return the rendered HTML
	html := strings.TrimSpace(buf.String())

	// Only trim <p> tags if it's a single paragraph (no newlines, starts/ends with <p>)
	if strings.HasPrefix(html, "<p>") && strings.HasSuffix(html, "</p>") && strings.Count(html, "<p>") == 1 {
		html = strings.TrimPrefix(html, "<p>")
		html = strings.TrimSuffix(html, "</p>")
	}

	return html
}

// convertSecurity converts security requirements
func (p *OpenAPIParser) convertSecurity(security *openapi3.SecurityRequirements) []map[string][]string {
	if security == nil {
		return nil
	}

	result := make([]map[string][]string, 0, len(*security))
	for _, req := range *security {
		m := make(map[string][]string)
		for name, scopes := range req {
			m[name] = scopes
		}
		result = append(result, m)
	}
	return result
}

// organizeByPath organizes endpoints into a hierarchical path structure
func (p *OpenAPIParser) organizeByPath(endpoints []*core.APIEndpoint) []*core.APIPathGroup {
	// Build a tree structure using map-based approach
	root := &pathNode{
		children: make(map[string]*pathNode),
	}

	// Build the tree
	for _, endpoint := range endpoints {
		parts := strings.Split(strings.Trim(endpoint.Path, "/"), "/")
		if len(parts) == 0 || (len(parts) == 1 && parts[0] == "") {
			// Root path "/"
			parts = []string{"/"}
		}
		p.addToPathNode(root, parts, endpoint, 0)
	}

	// Convert to APIPathGroup structure
	return p.pathNodesToGroups(root, "")
}

// pathNode is an internal structure for building the path tree
type pathNode struct {
	endpoints []*core.APIEndpoint
	children  map[string]*pathNode
}

// addToPathNode recursively adds an endpoint to the path node tree
func (p *OpenAPIParser) addToPathNode(node *pathNode, parts []string, endpoint *core.APIEndpoint, depth int) {
	if depth >= len(parts) {
		node.endpoints = append(node.endpoints, endpoint)
		return
	}

	part := parts[depth]
	if node.children[part] == nil {
		node.children[part] = &pathNode{
			endpoints: []*core.APIEndpoint{},
			children:  make(map[string]*pathNode),
		}
	}

	if depth == len(parts)-1 {
		// Last part - add endpoint here
		node.children[part].endpoints = append(node.children[part].endpoints, endpoint)
	} else {
		// Continue recursively
		p.addToPathNode(node.children[part], parts, endpoint, depth+1)
	}
}

// pathNodesToGroups converts pathNode tree to APIPathGroup slices
func (p *OpenAPIParser) pathNodesToGroups(node *pathNode, pathPrefix string) []*core.APIPathGroup {
	var groups []*core.APIPathGroup

	// Sort keys for consistent output
	keys := make([]string, 0, len(node.children))
	for key := range node.children {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		child := node.children[key]

		// Build the full path
		fullPath := pathPrefix + "/" + key

		group := &core.APIPathGroup{
			Path:      fullPath,
			Endpoints: child.endpoints,
			Children:  p.pathNodesToGroups(child, fullPath),
		}

		groups = append(groups, group)
	}

	return groups
}

// flattenPathTree converts a path tree map to a sorted slice
func (p *OpenAPIParser) flattenPathTree(tree map[string]*core.APIPathGroup) []*core.APIPathGroup {
	result := make([]*core.APIPathGroup, 0, len(tree))
	for _, group := range tree {
		result = append(result, group)
	}

	// Sort by path
	sort.Slice(result, func(i, j int) bool {
		return result[i].Path < result[j].Path
	})

	return result
}

// organizeByTag organizes endpoints by tags
func (p *OpenAPIParser) organizeByTag(endpoints []*core.APIEndpoint, tags []core.APITag) []*core.APITagGroup {
	tagMap := make(map[string]*core.APITagGroup)

	// Create tag groups from defined tags
	for _, tag := range tags {
		tagMap[tag.Name] = &core.APITagGroup{
			Tag:       tag,
			Endpoints: []*core.APIEndpoint{},
		}
	}

	// Add endpoints to their tag groups
	for _, endpoint := range endpoints {
		for _, tagName := range endpoint.Tags {
			if tagMap[tagName] == nil {
				// Create tag group for undefined tags
				tagMap[tagName] = &core.APITagGroup{
					Tag: core.APITag{
						Name:        tagName,
						Description: "",
					},
					Endpoints: []*core.APIEndpoint{},
				}
			}
			tagMap[tagName].Endpoints = append(tagMap[tagName].Endpoints, endpoint)
		}

		// If endpoint has no tags, add to "Untagged" group
		if len(endpoint.Tags) == 0 {
			if tagMap["Untagged"] == nil {
				tagMap["Untagged"] = &core.APITagGroup{
					Tag: core.APITag{
						Name:        "Untagged",
						Description: "Endpoints without tags",
					},
					Endpoints: []*core.APIEndpoint{},
				}
			}
			tagMap["Untagged"].Endpoints = append(tagMap["Untagged"].Endpoints, endpoint)
		}
	}

	// Convert map to sorted slice
	result := make([]*core.APITagGroup, 0, len(tagMap))
	for _, group := range tagMap {
		if len(group.Endpoints) > 0 {
			result = append(result, group)
		}
	}

	// Sort by tag name
	sort.Slice(result, func(i, j int) bool {
		return result[i].Tag.Name < result[j].Tag.Name
	})

	return result
}

// NameFromURL extracts a filename from a URL
func (p *OpenAPIParser) NameFromURL(url string) string {
	// Extract filename from URL
	parts := strings.Split(url, "/")
	name := parts[len(parts)-1]

	// Remove extension
	if idx := strings.LastIndex(name, "."); idx != -1 {
		name = name[:idx]
	}

	// If empty, use a hash of the URL
	if name == "" {
		name = fmt.Sprintf("spec-%x", time.Now().Unix())
	}

	return name
}
