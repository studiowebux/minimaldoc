// Package markdown provides markdown rendering services for the server.
// All rendered output is sanitized with bluemonday because server-side
// markdown comes from user-generated content (blog posts, forum topics,
// comments) and may contain malicious HTML/JS.
package markdown

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/studiowebux/minimaldoc/internal/parser"
)

// Renderer provides markdown to HTML conversion with XSS sanitization.
type Renderer struct {
	parser    *parser.MarkdownParser
	sanitizer *bluemonday.Policy
}

// NewRenderer creates a new markdown renderer with UGC sanitization policy.
func NewRenderer() *Renderer {
	return &Renderer{
		parser:    parser.NewMarkdownParser(),
		sanitizer: bluemonday.UGCPolicy(),
	}
}

// Render converts markdown content to sanitized HTML.
func (r *Renderer) Render(content string) (string, error) {
	html, err := r.parser.Parse([]byte(content))
	if err != nil {
		return "", err
	}
	return r.sanitizer.Sanitize(string(html)), nil
}

// RenderBytes converts markdown content to sanitized HTML bytes.
func (r *Renderer) RenderBytes(content []byte) ([]byte, error) {
	html, err := r.parser.Parse(content)
	if err != nil {
		return nil, err
	}
	return r.sanitizer.SanitizeBytes(html), nil
}
