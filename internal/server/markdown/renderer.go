// Package markdown provides markdown rendering services for the server.
package markdown

import (
	"github.com/studiowebux/minimaldoc/internal/parser"
)

// Renderer provides markdown to HTML conversion.
type Renderer struct {
	parser *parser.MarkdownParser
}

// NewRenderer creates a new markdown renderer.
func NewRenderer() *Renderer {
	return &Renderer{
		parser: parser.NewMarkdownParser(),
	}
}

// Render converts markdown content to HTML.
func (r *Renderer) Render(content string) (string, error) {
	html, err := r.parser.Parse([]byte(content))
	if err != nil {
		return "", err
	}
	return string(html), nil
}

// RenderBytes converts markdown content to HTML bytes.
func (r *Renderer) RenderBytes(content []byte) ([]byte, error) {
	return r.parser.Parse(content)
}
