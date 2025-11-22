package parser

import (
	"bytes"
	"os"

	"github.com/adrg/frontmatter"
	"github.com/studiowebux/minimaldoc/internal/core"
)

// FrontmatterParser handles extraction of YAML frontmatter from markdown
type FrontmatterParser struct{}

// NewFrontmatterParser creates a new frontmatter parser
func NewFrontmatterParser() *FrontmatterParser {
	return &FrontmatterParser{}
}

// Parse extracts frontmatter from markdown content
// Returns the metadata and the remaining markdown content (without frontmatter)
func (p *FrontmatterParser) Parse(content []byte) (core.Metadata, []byte, error) {
	meta := core.DefaultMetadata()

	// Parse frontmatter
	rest, err := frontmatter.Parse(bytes.NewReader(content), &meta)
	if err != nil {
		// If no frontmatter found or parse error, return defaults with full content
		return meta, content, nil
	}

	// rest is already []byte, not a Reader
	return meta, rest, nil
}

// ParseFile reads and parses a markdown file
func (p *FrontmatterParser) ParseFile(path string) (core.Metadata, []byte, error) {
	meta := core.DefaultMetadata()

	// Read file
	content, err := os.ReadFile(path)
	if err != nil {
		return meta, nil, err
	}

	// Parse frontmatter
	rest, err := frontmatter.Parse(bytes.NewReader(content), &meta)
	if err != nil {
		// If no frontmatter found or parse error, return defaults with full content
		return meta, content, nil
	}

	// rest is already []byte, not a Reader
	return meta, rest, nil
}
