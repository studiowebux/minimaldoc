package core

// MCPConfig holds MCP server documentation configuration
type MCPConfig struct {
	Enabled   bool     `yaml:"enabled"`
	SpecFiles []string `yaml:"spec_files"` // JSON manifest file paths/globs relative to docs root
	Path      string   `yaml:"path"`       // Output URL path (default: "mcp")
}

// MCPSpec represents a parsed MCP server manifest
type MCPSpec struct {
	Name        string
	Description string
	Tools       []*MCPTool
	Resources   []*MCPResource
	Prompts     []*MCPPrompt
	OutputPath  string
}

// MCPTool represents a single MCP tool
type MCPTool struct {
	Name        string
	Description string
	InputSchema *MCPSchema
}

// MCPSchema represents a JSON Schema for tool input parameters
type MCPSchema struct {
	Type        string
	Description string
	Properties  map[string]*MCPSchema
	Required    []string
	Enum        []any
	Items       *MCPSchema // for array types
}

// MCPResource represents a single MCP resource
type MCPResource struct {
	URI         string
	Name        string
	Description string
	MimeType    string
}

// MCPPrompt represents a single MCP prompt
type MCPPrompt struct {
	Name        string
	Description string
	Arguments   []*MCPPromptArgument
}

// MCPPromptArgument represents a single prompt argument
type MCPPromptArgument struct {
	Name        string
	Description string
	Required    bool
}

// DefaultMCPConfig returns an MCPConfig with sensible defaults
func DefaultMCPConfig() MCPConfig {
	return MCPConfig{
		Enabled:   false,
		SpecFiles: []string{},
		Path:      "mcp",
	}
}
