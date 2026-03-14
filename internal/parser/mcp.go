package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// MCPParser handles parsing of MCP server manifest files
type MCPParser struct{}

// NewMCPParser creates a new MCP parser
func NewMCPParser() *MCPParser {
	return &MCPParser{}
}

// mcpManifest is the raw JSON structure of an MCP server manifest file
type mcpManifest struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Tools       []mcpRawTool `json:"tools"`
	Resources   []mcpRawRes  `json:"resources"`
	Prompts     []mcpRawProm `json:"prompts"`
}

type mcpRawTool struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	InputSchema *mcpRawSchema `json:"inputSchema"`
}

type mcpRawSchema struct {
	Type        string                   `json:"type"`
	Description string                   `json:"description"`
	Properties  map[string]*mcpRawSchema `json:"properties"`
	Required    []string                 `json:"required"`
	Enum        []any                    `json:"enum"`
	Items       *mcpRawSchema            `json:"items"`
}

type mcpRawRes struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`
}

type mcpRawProm struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Arguments   []mcpRawPromArg `json:"arguments"`
}

type mcpRawPromArg struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// ParseFile parses an MCP server manifest JSON file into an MCPSpec
func (p *MCPParser) ParseFile(filePath string) (*core.MCPSpec, error) {
	data, err := os.ReadFile(filePath) // #nosec G304 -- path comes from glob discovery within the docs root, not user input
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}

	var raw mcpManifest
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse manifest JSON: %w", err)
	}

	spec := &core.MCPSpec{
		Name:        raw.Name,
		Description: raw.Description,
	}

	// Use filename as name fallback
	if spec.Name == "" {
		base := filepath.Base(filePath)
		ext := filepath.Ext(base)
		spec.Name = base[:len(base)-len(ext)]
	}

	for i := range raw.Tools {
		spec.Tools = append(spec.Tools, convertTool(&raw.Tools[i]))
	}

	for i := range raw.Resources {
		spec.Resources = append(spec.Resources, convertResource(&raw.Resources[i]))
	}

	for i := range raw.Prompts {
		spec.Prompts = append(spec.Prompts, convertPrompt(&raw.Prompts[i]))
	}

	return spec, nil
}

func convertTool(r *mcpRawTool) *core.MCPTool {
	t := &core.MCPTool{
		Name:        r.Name,
		Description: r.Description,
	}
	if r.InputSchema != nil {
		t.InputSchema = convertSchema(r.InputSchema)
	}
	return t
}

func convertSchema(r *mcpRawSchema) *core.MCPSchema {
	if r == nil {
		return nil
	}
	s := &core.MCPSchema{
		Type:        r.Type,
		Description: r.Description,
		Required:    r.Required,
		Enum:        r.Enum,
	}
	if len(r.Properties) > 0 {
		s.Properties = make(map[string]*core.MCPSchema, len(r.Properties))
		for k, v := range r.Properties {
			s.Properties[k] = convertSchema(v)
		}
	}
	if r.Items != nil {
		s.Items = convertSchema(r.Items)
	}
	return s
}

func convertResource(r *mcpRawRes) *core.MCPResource {
	return &core.MCPResource{
		URI:         r.URI,
		Name:        r.Name,
		Description: r.Description,
		MimeType:    r.MimeType,
	}
}

func convertPrompt(r *mcpRawProm) *core.MCPPrompt {
	p := &core.MCPPrompt{
		Name:        r.Name,
		Description: r.Description,
	}
	for i := range r.Arguments {
		p.Arguments = append(p.Arguments, &core.MCPPromptArgument{
			Name:        r.Arguments[i].Name,
			Description: r.Arguments[i].Description,
			Required:    r.Arguments[i].Required,
		})
	}
	return p
}
