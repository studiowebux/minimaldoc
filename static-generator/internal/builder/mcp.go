package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/studiowebux/minimaldoc/static-generator/internal/core"
	"github.com/studiowebux/minimaldoc/static-generator/internal/parser"
)

// MCPBuilder discovers and parses MCP server manifest files
type MCPBuilder struct {
	parser *parser.MCPParser
}

// NewMCPBuilder creates a new MCP builder
func NewMCPBuilder() *MCPBuilder {
	return &MCPBuilder{
		parser: parser.NewMCPParser(),
	}
}

// Build discovers manifest files and parses them into MCPSpec structs
func (b *MCPBuilder) Build(docsRoot string, config core.MCPConfig) ([]*core.MCPSpec, error) {
	files, err := b.discoverFiles(docsRoot, config)
	if err != nil {
		return nil, fmt.Errorf("discover MCP manifests: %w", err)
	}

	var specs []*core.MCPSpec
	for _, f := range files {
		spec, err := b.parser.ParseFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse MCP manifest %s: %v\n", f, err)
			continue
		}
		specs = append(specs, spec)
	}

	return specs, nil
}

// discoverFiles returns all manifest file paths to parse
func (b *MCPBuilder) discoverFiles(docsRoot string, config core.MCPConfig) ([]string, error) {
	var files []string

	if len(config.SpecFiles) > 0 {
		// Explicit list — support globs
		for _, pattern := range config.SpecFiles {
			if strings.Contains(pattern, "*") {
				matches, err := filepath.Glob(filepath.Join(docsRoot, pattern))
				if err != nil {
					return nil, fmt.Errorf("invalid glob %s: %w", pattern, err)
				}
				files = append(files, matches...)
			} else {
				full := filepath.Join(docsRoot, pattern)
				if _, err := os.Stat(full); err == nil {
					files = append(files, full)
				}
			}
		}
		return files, nil
	}

	// Auto-discover *.mcp.json files under docs root
	err := filepath.WalkDir(docsRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".mcp.json") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
