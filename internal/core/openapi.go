package core

import "time"

// APISpec represents a single OpenAPI/Swagger specification
type APISpec struct {
	// Metadata
	Name        string    // Spec name (filename or derived from URL)
	Title       string    // API title from spec info.title
	Description string    // API description from spec info.description
	Version     string    // API version from spec info.version
	FilePath    string    // Local file path (if from file)
	URL         string    // Remote URL (if fetched)
	LastFetched time.Time // Last fetch time for remote specs
	ETag        string    // ETag for cache validation
	LastModified string   // Last-Modified header for cache validation

	// Specification data
	OpenAPIVersion string              // OpenAPI version (3.0.x, 3.1.x)
	Servers        []APIServer         // API servers
	Endpoints      []*APIEndpoint      // All endpoints/operations
	SecuritySchemes map[string]*APISecurityScheme // Security schemes defined in spec
	Tags           []APITag            // Tags for organization

	// Organization
	OrganizedByPath []*APIPathGroup   // Endpoints organized by path hierarchy
	OrganizedByTag  []*APITagGroup    // Endpoints organized by tags
	FlatEndpoints   []*APIEndpoint    // Flat list for search

	// Output paths
	OutputPath string // Generated HTML file path
	ChunkFiles []string // Lazy-load chunk file paths
}

// APIServer represents an OpenAPI server definition
type APIServer struct {
	URL         string                 // Server URL
	Description string                 // Server description
	Variables   map[string]APIServerVariable // Server variables
}

// APIServerVariable represents a server URL variable
type APIServerVariable struct {
	Default     string   // Default value
	Description string   // Variable description
	Enum        []string // Allowed values
}

// APIEndpoint represents a single API operation (path + method)
type APIEndpoint struct {
	// Identity
	OperationID string // Unique operation ID
	Path        string // URL path (e.g., /users/{id})
	Method      string // HTTP method (GET, POST, etc.)

	// Documentation
	Summary     string   // Short summary
	Description string   // Long description
	Tags        []string // Associated tags
	Deprecated  bool     // Is deprecated

	// Request
	Parameters   []*APIParameter  // Path, query, header, cookie params
	RequestBody  *APIRequestBody  // Request body

	// Responses
	Responses    map[string]*APIResponse // Status code -> response

	// Security
	Security     []map[string][]string // Security requirements

	// Spec reference
	SpecName     string // Parent spec name
}

// APIParameter represents a parameter (path, query, header, cookie)
type APIParameter struct {
	Name        string      // Parameter name
	In          string      // Location: path, query, header, cookie
	Description string      // Parameter description
	Required    bool        // Is required
	Deprecated  bool        // Is deprecated
	Schema      *APISchema  // Parameter schema
	Example     interface{} // Example value
	Examples    map[string]*APIExample // Named examples
}

// APIRequestBody represents a request body
type APIRequestBody struct {
	Description string                  // Request body description
	Required    bool                    // Is required
	Content     map[string]*APIMediaType // Media type -> content
}

// APIResponse represents a response
type APIResponse struct {
	Description string                  // Response description
	Headers     map[string]*APIParameter // Response headers
	Content     map[string]*APIMediaType // Media type -> content
}

// APIMediaType represents a media type (e.g., application/json)
type APIMediaType struct {
	Schema   *APISchema           // Schema for this media type
	Example  interface{}          // Example value
	Examples map[string]*APIExample // Named examples
	Encoding map[string]*APIEncoding // Encoding for multipart
}

// APISchema represents a JSON schema
type APISchema struct {
	Type        string                 // Schema type (string, number, object, array, etc.)
	Format      string                 // Format (date-time, email, uuid, etc.)
	Description string                 // Schema description
	Properties  map[string]*APISchema  // Object properties
	Items       *APISchema             // Array items schema
	Required    []string               // Required properties
	Enum        []interface{}          // Enumeration values
	Default     interface{}            // Default value
	Example     interface{}            // Example value
	Nullable    bool                   // Can be null
	ReadOnly    bool                   // Read-only property
	WriteOnly   bool                   // Write-only property

	// Validation
	Minimum          *float64 // Minimum value
	Maximum          *float64 // Maximum value
	MinLength        *int     // Minimum string length
	MaxLength        *int     // Maximum string length
	Pattern          string   // Regex pattern
	MinItems         *int     // Minimum array items
	MaxItems         *int     // Maximum array items
	UniqueItems      bool     // Array items must be unique
	MinProperties    *int     // Minimum object properties
	MaxProperties    *int     // Maximum object properties

	// Composition
	AllOf []APISchema // All of these schemas
	OneOf []APISchema // One of these schemas
	AnyOf []APISchema // Any of these schemas
	Not   *APISchema  // Not this schema

	// Reference
	Ref string // $ref value
}

// APIExample represents a named example
type APIExample struct {
	Summary     string      // Example summary
	Description string      // Example description
	Value       interface{} // Example value
}

// APIEncoding represents encoding for multipart/form-data
type APIEncoding struct {
	ContentType   string               // Content-Type for this field
	Headers       map[string]*APIParameter // Headers
	Style         string               // Serialization style
	Explode       bool                 // Explode parameter
	AllowReserved bool                 // Allow reserved characters
}

// APISecurityScheme represents a security scheme
type APISecurityScheme struct {
	Type             string // Type: apiKey, http, oauth2, openIdConnect
	Description      string // Description
	Name             string // Header/query/cookie name (for apiKey)
	In               string // Location: header, query, cookie (for apiKey)
	Scheme           string // Auth scheme: basic, bearer, etc. (for http)
	BearerFormat     string // Bearer token format (for http bearer)
	Flows            *APIOAuthFlows // OAuth flows (for oauth2)
	OpenIDConnectURL string // OpenID Connect URL (for openIdConnect)
}

// APIOAuthFlows represents OAuth 2.0 flows
type APIOAuthFlows struct {
	Implicit          *APIOAuthFlow // Implicit flow
	Password          *APIOAuthFlow // Password flow
	ClientCredentials *APIOAuthFlow // Client credentials flow
	AuthorizationCode *APIOAuthFlow // Authorization code flow
}

// APIOAuthFlow represents a single OAuth flow
type APIOAuthFlow struct {
	AuthorizationURL string            // Authorization URL
	TokenURL         string            // Token URL
	RefreshURL       string            // Refresh URL
	Scopes           map[string]string // Available scopes (scope -> description)
}

// APITag represents an OpenAPI tag
type APITag struct {
	Name        string // Tag name
	Description string // Tag description
}

// APIPathGroup represents endpoints grouped by path hierarchy
type APIPathGroup struct {
	Path      string          // Path segment (e.g., /users)
	Endpoints []*APIEndpoint  // Endpoints at this path
	Children  []*APIPathGroup // Sub-paths
}

// APITagGroup represents endpoints grouped by tag
type APITagGroup struct {
	Tag       APITag         // Tag info
	Endpoints []*APIEndpoint // Endpoints with this tag
}

// APISpecSource represents the source of an API spec
type APISpecSource struct {
	Type     string // "file" or "url"
	Location string // File path or URL
	CachePath string // Cache file path for URL sources
}

// OpenAPIConfig holds OpenAPI-specific configuration
type OpenAPIConfig struct {
	Enabled         bool     // Enable OpenAPI support
	SpecFiles       []string // Local spec file paths/globs
	SpecURLs        []string // Remote spec URLs to fetch
	DefaultView     string   // Default organization view: "path", "tag", "flat"
	SyncOnBuild     bool     // Sync remote specs on every build
	CacheDir        string   // Cache directory for remote specs
	EnableTesting   bool     // Enable API testing UI
	EnableExport    bool     // Enable export to restcli/curl
	LazyLoadChunkSize int    // Chunk size for lazy loading (bytes)
}

// DefaultOpenAPIConfig returns an OpenAPIConfig with sensible defaults
func DefaultOpenAPIConfig() OpenAPIConfig {
	return OpenAPIConfig{
		Enabled:           false,
		SpecFiles:         []string{},
		SpecURLs:          []string{},
		DefaultView:       "path",
		SyncOnBuild:       false,
		CacheDir:          ".openapi-cache",
		EnableTesting:     true,
		EnableExport:      true,
		LazyLoadChunkSize: 51200, // 50KB chunks
	}
}
