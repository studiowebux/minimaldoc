package parser

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	neturl "net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/studiowebux/minimaldoc/internal/core"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	htmlrenderer "github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

// OpenAPIParser handles parsing of OpenAPI specifications.
// Descriptions in OpenAPI specs may come from remote URLs, so the rendered
// HTML is sanitized with bluemonday to prevent XSS from untrusted specs.
type OpenAPIParser struct {
	cacheDir  string
	md        goldmark.Markdown  // Markdown renderer for descriptions
	sanitizer *bluemonday.Policy // HTML sanitizer for rendered descriptions
}

// NewOpenAPIParser creates a new OpenAPI parser
func NewOpenAPIParser(cacheDir string) *OpenAPIParser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			htmlrenderer.WithUnsafe(),
		),
	)

	return &OpenAPIParser{
		cacheDir:  cacheDir,
		md:        md,
		sanitizer: bluemonday.UGCPolicy(),
	}
}

// ParseFile parses an OpenAPI spec from a local file
func (p *OpenAPIParser) ParseFile(filePath string) (*core.APISpec, error) {
	specBytes, err := os.ReadFile(filePath) // #nosec G304 -- file path comes from site config, not user input
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI spec from %s: %w", filePath, err)
	}

	doc, err := libopenapi.NewDocumentWithConfiguration(specBytes, &datamodel.DocumentConfiguration{
		AllowFileReferences:   true,  // Needed for multi-file specs ($ref: "./schemas/user.yaml")
		AllowRemoteReferences: false, // Disabled: local specs should not trigger HTTP requests (SSRF risk)
		BasePath:              filepath.Dir(filePath),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec from %s: %w", filePath, err)
	}

	v3Model, errs := doc.BuildV3Model()
	if v3Model == nil {
		return nil, fmt.Errorf("failed to build OpenAPI 3.x model from %s: %v", filePath, errs)
	}

	spec := p.convertSpec(&v3Model.Model)
	spec.Name = filepath.Base(filePath)
	spec.FilePath = filePath

	return spec, nil
}

// ParseURL fetches and parses an OpenAPI spec from a URL
func (p *OpenAPIParser) ParseURL(url string) (*core.APISpec, error) {
	if err := os.MkdirAll(p.cacheDir, 0755); err != nil { // #nosec G301 -- cache directory is not sensitive
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	parsed, parseErr := neturl.Parse(url)
	if parseErr != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return nil, fmt.Errorf("invalid OpenAPI spec URL (must be http or https): %s", url)
	}

	// Block requests to private/loopback IPs to prevent SSRF
	hostname := parsed.Hostname()
	ips, lookupErr := net.LookupIP(hostname)
	if lookupErr != nil {
		return nil, fmt.Errorf("DNS lookup failed for OpenAPI spec URL %s: %w", url, lookupErr)
	}
	for _, ip := range ips {
		if isPrivateOrLoopback(ip) {
			return nil, fmt.Errorf("OpenAPI spec URL %s resolves to private/loopback address", url)
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url) // #nosec G107 -- scheme validated, private IPs blocked
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OpenAPI spec from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch OpenAPI spec from %s: status %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI spec from %s: %w", url, err)
	}

	doc, err := libopenapi.NewDocument(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec from %s: %w", url, err)
	}

	v3Model, errs := doc.BuildV3Model()
	if v3Model == nil {
		return nil, fmt.Errorf("failed to build OpenAPI 3.x model from %s: %v", url, errs)
	}

	spec := p.convertSpec(&v3Model.Model)
	spec.Name = p.NameFromURL(url)
	spec.URL = url
	spec.LastFetched = time.Now()
	spec.ETag = resp.Header.Get("ETag")
	spec.LastModified = resp.Header.Get("Last-Modified")

	cachePath := filepath.Join(p.cacheDir, spec.Name+".json")
	if err := os.WriteFile(cachePath, body, 0644); err != nil { // #nosec G306 -- OpenAPI spec cache is not sensitive
		fmt.Fprintf(os.Stderr, "Warning: failed to cache OpenAPI spec: %v\n", err)
	}

	return spec, nil
}

// convertSpec converts a libopenapi v3 Document to our internal APISpec
func (p *OpenAPIParser) convertSpec(model *v3high.Document) *core.APISpec {
	spec := &core.APISpec{
		Title:          model.Info.Title,
		Description:    p.renderMarkdown(model.Info.Description),
		Version:        model.Info.Version,
		OpenAPIVersion: model.Version,
		Servers:        p.convertServers(model.Servers),
		Tags:           p.convertTags(model.Tags),
		Endpoints:      []*core.APIEndpoint{},
	}

	if model.Components != nil {
		spec.SecuritySchemes = p.convertSecuritySchemes(model.Components.SecuritySchemes)
		spec.Schemas = p.convertSchemas(model.Components.Schemas)
	}

	if model.Paths != nil && model.Paths.PathItems != nil {
		for path, pathItem := range model.Paths.PathItems.FromOldest() {
			spec.Endpoints = append(spec.Endpoints, p.convertPathItem(path, pathItem, spec.Name)...)
		}
	}

	spec.OrganizedByPath = p.organizeByPath(spec.Endpoints)
	spec.OrganizedByTag = p.organizeByTag(spec.Endpoints, spec.Tags)
	spec.FlatEndpoints = spec.Endpoints

	return spec
}

// convertSchemas converts OpenAPI component schemas
func (p *OpenAPIParser) convertSchemas(schemas *orderedmap.Map[string, *base.SchemaProxy]) map[string]*core.APISchema {
	if schemas == nil {
		return nil
	}

	result := make(map[string]*core.APISchema)
	for name, schemaProxy := range schemas.FromOldest() {
		if schemaProxy == nil {
			continue
		}
		result[name] = p.convertSchemaProxy(schemaProxy)
	}
	return result
}

// convertServers converts OpenAPI servers
func (p *OpenAPIParser) convertServers(servers []*v3high.Server) []core.APIServer {
	result := make([]core.APIServer, 0, len(servers))
	for _, server := range servers {
		s := core.APIServer{
			URL:         server.URL,
			Description: p.renderMarkdown(server.Description),
			Variables:   make(map[string]core.APIServerVariable),
		}

		if server.Variables != nil {
			for name, variable := range server.Variables.FromOldest() {
				sv := core.APIServerVariable{
					Description: p.renderMarkdown(variable.Description),
					Enum:        variable.Enum,
				}
				if variable.Default != "" {
					sv.Default = variable.Default
				}
				s.Variables[name] = sv
			}
		}

		result = append(result, s)
	}
	return result
}

// convertTags converts OpenAPI tags
func (p *OpenAPIParser) convertTags(tags []*base.Tag) []core.APITag {
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
func (p *OpenAPIParser) convertSecuritySchemes(schemes *orderedmap.Map[string, *v3high.SecurityScheme]) map[string]*core.APISecurityScheme {
	if schemes == nil {
		return make(map[string]*core.APISecurityScheme)
	}

	result := make(map[string]*core.APISecurityScheme)
	for name, scheme := range schemes.FromOldest() {
		if scheme == nil {
			continue
		}

		s := &core.APISecurityScheme{
			Type:             scheme.Type,
			Description:      p.renderMarkdown(scheme.Description),
			Name:             scheme.Name,
			In:               scheme.In,
			Scheme:           scheme.Scheme,
			BearerFormat:     scheme.BearerFormat,
			OpenIDConnectURL: scheme.OpenIdConnectUrl,
		}

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
func (p *OpenAPIParser) convertOAuthFlow(flow *v3high.OAuthFlow) *core.APIOAuthFlow {
	f := &core.APIOAuthFlow{
		AuthorizationURL: flow.AuthorizationUrl,
		TokenURL:         flow.TokenUrl,
		RefreshURL:       flow.RefreshUrl,
	}

	if flow.Scopes != nil {
		f.Scopes = make(map[string]string)
		for scope, desc := range flow.Scopes.FromOldest() {
			f.Scopes[scope] = desc
		}
	}

	return f
}

// convertPathItem converts a path item to endpoints
func (p *OpenAPIParser) convertPathItem(path string, pathItem *v3high.PathItem, specName string) []*core.APIEndpoint {
	endpoints := []*core.APIEndpoint{}

	operations := map[core.HTTPMethod]*v3high.Operation{
		core.HTTPGet:     pathItem.Get,
		core.HTTPPost:    pathItem.Post,
		core.HTTPPut:     pathItem.Put,
		core.HTTPPatch:   pathItem.Patch,
		core.HTTPDelete:  pathItem.Delete,
		core.HTTPHead:    pathItem.Head,
		core.HTTPOptions: pathItem.Options,
		core.HTTPTrace:   pathItem.Trace,
	}

	for method, operation := range operations {
		if operation == nil {
			continue
		}

		endpoint := &core.APIEndpoint{
			OperationID: operation.OperationId,
			Path:        path,
			Method:      string(method),
			Summary:     p.renderMarkdown(operation.Summary),
			Description: p.renderMarkdown(operation.Description),
			Tags:        operation.Tags,
			Deprecated:  derefBool(operation.Deprecated),
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
func (p *OpenAPIParser) convertParameters(params []*v3high.Parameter) []*core.APIParameter {
	result := make([]*core.APIParameter, 0, len(params))
	for _, param := range params {
		if param == nil {
			continue
		}

		apiParam := &core.APIParameter{
			Name:        param.Name,
			In:          param.In,
			Description: p.renderMarkdown(param.Description),
			Required:    derefBool(param.Required),
			Deprecated:  param.Deprecated,
			Schema:      p.convertSchemaProxy(param.Schema),
			Example:     decodeYAMLNode(param.Example),
		}

		if param.Examples != nil {
			apiParam.Examples = make(map[string]*core.APIExample)
			for name, example := range param.Examples.FromOldest() {
				if example != nil {
					apiParam.Examples[name] = &core.APIExample{
						Summary:     p.renderMarkdown(example.Summary),
						Description: p.renderMarkdown(example.Description),
						Value:       decodeYAMLNode(example.Value),
					}
				}
			}
		}

		result = append(result, apiParam)
	}
	return result
}

// convertRequestBody converts an OpenAPI request body
func (p *OpenAPIParser) convertRequestBody(body *v3high.RequestBody) *core.APIRequestBody {
	if body == nil {
		return nil
	}

	rb := &core.APIRequestBody{
		Description: p.renderMarkdown(body.Description),
		Required:    derefBool(body.Required),
		Content:     make(map[string]*core.APIMediaType),
	}

	if body.Content != nil {
		for mediaType, mediaTypeValue := range body.Content.FromOldest() {
			rb.Content[mediaType] = p.convertMediaType(mediaTypeValue)
		}
	}

	return rb
}

// convertResponses converts OpenAPI responses
func (p *OpenAPIParser) convertResponses(responses *v3high.Responses) map[string]*core.APIResponse {
	result := make(map[string]*core.APIResponse)
	if responses == nil {
		return result
	}

	if responses.Codes != nil {
		for status, response := range responses.Codes.FromOldest() {
			if response == nil {
				continue
			}
			result[status] = p.convertResponse(response)
		}
	}

	if responses.Default != nil {
		result["default"] = p.convertResponse(responses.Default)
	}

	return result
}

// convertResponse converts a single OpenAPI response
func (p *OpenAPIParser) convertResponse(response *v3high.Response) *core.APIResponse {
	r := &core.APIResponse{
		Description: p.renderMarkdown(response.Description),
		Headers:     make(map[string]*core.APIParameter),
		Content:     make(map[string]*core.APIMediaType),
	}

	if response.Headers != nil {
		for name, header := range response.Headers.FromOldest() {
			if header != nil {
				r.Headers[name] = &core.APIParameter{
					Name:        name,
					In:          string(core.ParamInHeader),
					Description: p.renderMarkdown(header.Description),
					Required:    header.Required,
					Schema:      p.convertSchemaProxy(header.Schema),
				}
			}
		}
	}

	if response.Content != nil {
		for mediaType, mediaTypeValue := range response.Content.FromOldest() {
			r.Content[mediaType] = p.convertMediaType(mediaTypeValue)
		}
	}

	return r
}

// convertMediaType converts a media type
func (p *OpenAPIParser) convertMediaType(mt *v3high.MediaType) *core.APIMediaType {
	if mt == nil {
		return nil
	}

	m := &core.APIMediaType{
		Schema:  p.convertSchemaProxy(mt.Schema),
		Example: decodeYAMLNode(mt.Example),
	}

	if mt.Examples != nil {
		m.Examples = make(map[string]*core.APIExample)
		for name, example := range mt.Examples.FromOldest() {
			if example != nil {
				m.Examples[name] = &core.APIExample{
					Summary:     example.Summary,
					Description: example.Description,
					Value:       decodeYAMLNode(example.Value),
				}
			}
		}
	}

	return m
}

// convertSchemaProxy converts a SchemaProxy to our internal APISchema
func (p *OpenAPIParser) convertSchemaProxy(proxy *base.SchemaProxy) *core.APISchema {
	if proxy == nil {
		return nil
	}

	// If there's a $ref, keep the reference instead of inlining
	// This prevents massive duplication in the output JSON
	if proxy.IsReference() {
		return &core.APISchema{
			Ref: proxy.GetReference(),
		}
	}

	schema := proxy.Schema()
	if schema == nil {
		return nil
	}

	// Extract type as string
	typeStr := ""
	if len(schema.Type) > 0 {
		typeStr = schema.Type[0]
	}

	s := &core.APISchema{
		Type:        typeStr,
		Format:      schema.Format,
		Description: p.renderMarkdown(schema.Description),
		Required:    schema.Required,
		Enum:        decodeYAMLNodes(schema.Enum),
		Default:     decodeYAMLNode(schema.Default),
		Example:     decodeYAMLNode(schema.Example),
		Nullable:    derefBool(schema.Nullable),
		ReadOnly:    derefBool(schema.ReadOnly),
		WriteOnly:   derefBool(schema.WriteOnly),
		Deprecated:  derefBool(schema.Deprecated),
		Pattern:     schema.Pattern,
		UniqueItems: derefBool(schema.UniqueItems),
	}

	// Numeric constraints
	if schema.Minimum != nil {
		min := *schema.Minimum
		s.Minimum = &min
	}
	if schema.Maximum != nil {
		max := *schema.Maximum
		s.Maximum = &max
	}

	// String constraints
	if schema.MinLength != nil {
		ml := int(*schema.MinLength)
		s.MinLength = &ml
	}
	if schema.MaxLength != nil {
		ml := int(*schema.MaxLength)
		s.MaxLength = &ml
	}

	// Array constraints
	if schema.MinItems != nil {
		mi := int(*schema.MinItems)
		s.MinItems = &mi
	}
	if schema.MaxItems != nil {
		mi := int(*schema.MaxItems)
		s.MaxItems = &mi
	}

	// Object constraints
	if schema.MinProperties != nil {
		mp := int(*schema.MinProperties)
		s.MinProperties = &mp
	}
	if schema.MaxProperties != nil {
		mp := int(*schema.MaxProperties)
		s.MaxProperties = &mp
	}

	// Properties
	if schema.Properties != nil {
		s.Properties = make(map[string]*core.APISchema)
		for name, propProxy := range schema.Properties.FromOldest() {
			s.Properties[name] = p.convertSchemaProxy(propProxy)
		}
	}

	// Items (for arrays)
	if schema.Items != nil && schema.Items.IsA() {
		s.Items = p.convertSchemaProxy(schema.Items.A)
	}

	// Composition
	if len(schema.AllOf) > 0 {
		s.AllOf = make([]core.APISchema, 0, len(schema.AllOf))
		for _, sp := range schema.AllOf {
			if converted := p.convertSchemaProxy(sp); converted != nil {
				s.AllOf = append(s.AllOf, *converted)
			}
		}
	}
	if len(schema.OneOf) > 0 {
		s.OneOf = make([]core.APISchema, 0, len(schema.OneOf))
		for _, sp := range schema.OneOf {
			if converted := p.convertSchemaProxy(sp); converted != nil {
				s.OneOf = append(s.OneOf, *converted)
			}
		}
	}
	if len(schema.AnyOf) > 0 {
		s.AnyOf = make([]core.APISchema, 0, len(schema.AnyOf))
		for _, sp := range schema.AnyOf {
			if converted := p.convertSchemaProxy(sp); converted != nil {
				s.AnyOf = append(s.AnyOf, *converted)
			}
		}
	}
	if schema.Not != nil {
		s.Not = p.convertSchemaProxy(schema.Not)
	}

	return s
}

// renderMarkdown converts markdown description to HTML.
// Output is sanitized with bluemonday because OpenAPI specs may come from
// remote URLs — their descriptions are not fully trusted content.
func (p *OpenAPIParser) renderMarkdown(markdown string) string {
	if markdown == "" {
		return ""
	}

	var buf bytes.Buffer
	if err := p.md.Convert([]byte(markdown), &buf); err != nil {
		return markdown
	}

	// Sanitize rendered HTML to strip dangerous tags (script, event handlers, etc.)
	html := p.sanitizer.Sanitize(strings.TrimSpace(buf.String()))

	// Only trim <p> tags if it's a single paragraph (no newlines, starts/ends with <p>)
	if strings.HasPrefix(html, "<p>") && strings.HasSuffix(html, "</p>") && strings.Count(html, "<p>") == 1 {
		html = strings.TrimPrefix(html, "<p>")
		html = strings.TrimSuffix(html, "</p>")
	}

	return html
}

// convertSecurity converts security requirements
func (p *OpenAPIParser) convertSecurity(security []*base.SecurityRequirement) []map[string][]string {
	if security == nil {
		return nil
	}

	result := make([]map[string][]string, 0, len(security))
	for _, req := range security {
		if req == nil || req.Requirements == nil {
			continue
		}
		m := make(map[string][]string)
		for name, scopes := range req.Requirements.FromOldest() {
			m[name] = scopes
		}
		result = append(result, m)
	}
	return result
}

// organizeByPath organizes endpoints into a hierarchical path structure
func (p *OpenAPIParser) organizeByPath(endpoints []*core.APIEndpoint) []*core.APIPathGroup {
	root := &pathNode{
		children: make(map[string]*pathNode),
	}

	for _, endpoint := range endpoints {
		parts := strings.Split(strings.Trim(endpoint.Path, "/"), "/")
		if len(parts) == 0 || (len(parts) == 1 && parts[0] == "") {
			parts = []string{"/"}
		}
		p.addToPathNode(root, parts, endpoint, 0)
	}

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
		node.children[part].endpoints = append(node.children[part].endpoints, endpoint)
	} else {
		p.addToPathNode(node.children[part], parts, endpoint, depth+1)
	}
}

// pathNodesToGroups converts pathNode tree to APIPathGroup slices
func (p *OpenAPIParser) pathNodesToGroups(node *pathNode, pathPrefix string) []*core.APIPathGroup {
	var groups []*core.APIPathGroup

	keys := make([]string, 0, len(node.children))
	for key := range node.children {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		child := node.children[key]

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

// organizeByTag organizes endpoints by tags
func (p *OpenAPIParser) organizeByTag(endpoints []*core.APIEndpoint, tags []core.APITag) []*core.APITagGroup {
	tagMap := make(map[string]*core.APITagGroup)

	for _, tag := range tags {
		tagMap[tag.Name] = &core.APITagGroup{
			Tag:       tag,
			Endpoints: []*core.APIEndpoint{},
		}
	}

	for _, endpoint := range endpoints {
		for _, tagName := range endpoint.Tags {
			if tagMap[tagName] == nil {
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

	result := make([]*core.APITagGroup, 0, len(tagMap))
	for _, group := range tagMap {
		if len(group.Endpoints) > 0 {
			result = append(result, group)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Tag.Name < result[j].Tag.Name
	})

	return result
}

// NameFromURL extracts a filename from a URL
func (p *OpenAPIParser) NameFromURL(url string) string {
	parts := strings.Split(url, "/")
	name := parts[len(parts)-1]

	if idx := strings.LastIndex(name, "."); idx != -1 {
		name = name[:idx]
	}

	if name == "" {
		name = fmt.Sprintf("spec-%x", time.Now().Unix())
	}

	return name
}

// derefBool safely dereferences a *bool, returning false if nil
func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// decodeYAMLNode converts a *yaml.Node to a Go value
func decodeYAMLNode(node *yaml.Node) any {
	if node == nil {
		return nil
	}
	var val any
	if err := node.Decode(&val); err != nil {
		return nil
	}
	return val
}

// decodeYAMLNodes converts a slice of *yaml.Node to a slice of Go values
func decodeYAMLNodes(nodes []*yaml.Node) []any {
	if nodes == nil {
		return nil
	}
	result := make([]any, 0, len(nodes))
	for _, n := range nodes {
		result = append(result, decodeYAMLNode(n))
	}
	return result
}

// privateNetworks defines CIDR ranges blocked for SSRF prevention.
var privateNetworks = func() []*net.IPNet {
	cidrs := []string{
		"127.0.0.0/8", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16",
		"169.254.0.0/16", "::1/128", "fc00::/7", "fe80::/10",
	}
	nets := make([]*net.IPNet, 0, len(cidrs))
	for _, cidr := range cidrs {
		_, n, _ := net.ParseCIDR(cidr)
		nets = append(nets, n)
	}
	return nets
}()

// isPrivateOrLoopback returns true if the IP falls within a private/loopback range.
func isPrivateOrLoopback(ip net.IP) bool {
	for _, n := range privateNetworks {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}
