package parser

import (
	"bytes"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	htmlrenderer "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

// MarkdownParser handles markdown to HTML conversion
type MarkdownParser struct {
	md goldmark.Markdown
}

// NewMarkdownParser creates a new markdown parser with recommended extensions
func NewMarkdownParser() *MarkdownParser {
	md := goldmark.New(
		// Extensions
		goldmark.WithExtensions(
			NewAdmonitionExtension(), // Custom admonitions (:::info, :::warning, etc.)
			extension.GFM,            // GitHub Flavored Markdown (tables, strikethrough, etc.)
			extension.Typographer,    // Smart quotes, dashes, etc.
			highlighting.NewHighlighting( // Syntax highlighting
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					html.WithClasses(true),
					html.WithLineNumbers(false),
				),
			),
		),
		// Parser options
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // Automatically generate heading IDs
			parser.WithAttribute(),     // Enable {#id .class} syntax
		),
		// Renderer options
		goldmark.WithRendererOptions(
			htmlrenderer.WithHardWraps(), // Respect line breaks
			htmlrenderer.WithXHTML(),     // Use XHTML-style self-closing tags
			// WithUnsafe allows raw HTML passthrough (iframes, video embeds, custom
			// widgets). This is safe because the static site generator only processes
			// local markdown files authored by the doc owner — the same trust model
			// used by Hugo, Jekyll, and every other static site generator.
			// Do NOT remove: it would break legitimate embedded HTML in docs.
			htmlrenderer.WithUnsafe(),
		),
	)

	return &MarkdownParser{md: md}
}

// Parse converts markdown content to HTML
func (p *MarkdownParser) Parse(content []byte) ([]byte, error) {
	return p.ParseWithContext(content, "", "")
}

// ParseWithContext converts markdown content to HTML with link transformation
// currentPagePath is the relative path from docs root (e.g., "api/getting-started.md")
// basePath is the base path to prepend to all links (e.g., "/docs")
func (p *MarkdownParser) ParseWithContext(content []byte, currentPagePath string, basePath string) ([]byte, error) {
	// Parse markdown to AST
	reader := text.NewReader(content)
	doc := p.md.Parser().Parse(reader)

	// Transform .md links to .html if we have page context
	if currentPagePath != "" {
		TransformMarkdownLinks(doc, currentPagePath, basePath)
	}

	// Add target="_blank" to external links
	TransformExternalLinks(doc)

	// Render AST to HTML
	var buf bytes.Buffer
	if err := p.md.Renderer().Render(&buf, content, doc); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
