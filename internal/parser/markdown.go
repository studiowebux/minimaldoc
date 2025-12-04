package parser

import (
	"bytes"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	htmlrenderer "github.com/yuin/goldmark/renderer/html"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
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
			NewAdmonitionExtension(),       // Custom admonitions (:::info, :::warning, etc.)
			extension.GFM,                  // GitHub Flavored Markdown (tables, strikethrough, etc.)
			extension.Typographer,          // Smart quotes, dashes, etc.
			highlighting.NewHighlighting(   // Syntax highlighting
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
			htmlrenderer.WithUnsafe(),    // Allow raw HTML (for embeds, etc.)
		),
	)

	return &MarkdownParser{md: md}
}

// Parse converts markdown content to HTML
func (p *MarkdownParser) Parse(content []byte) ([]byte, error) {
	return p.ParseWithContext(content, "")
}

// ParseWithContext converts markdown content to HTML with link transformation
// currentPagePath is the relative path from docs root (e.g., "api/getting-started.md")
func (p *MarkdownParser) ParseWithContext(content []byte, currentPagePath string) ([]byte, error) {
	// Parse markdown to AST
	reader := text.NewReader(content)
	doc := p.md.Parser().Parse(reader)

	// Transform .md links to .html if we have page context
	if currentPagePath != "" {
		TransformMarkdownLinks(doc, currentPagePath)
	}

	// Render AST to HTML
	var buf bytes.Buffer
	if err := p.md.Renderer().Render(&buf, content, doc); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ParseWithExternalLinks converts markdown to HTML and marks external links with target="_blank"
// This is done as a post-processing step
func (p *MarkdownParser) ParseWithExternalLinks(content []byte, baseURL string) ([]byte, error) {
	html, err := p.Parse(content)
	if err != nil {
		return nil, err
	}

	// Post-process to add target="_blank" to external links
	// This will be handled by a custom transformer in the future
	// For now, we'll use a simple string replacement approach
	html = addExternalLinkAttributes(html, baseURL)

	return html, nil
}

// addExternalLinkAttributes adds target="_blank" and rel="noopener noreferrer" to external links
// This is a simple implementation - a more robust solution would use an AST transformer
func addExternalLinkAttributes(html []byte, baseURL string) []byte {
	// Simple regex-based replacement for external links
	// TODO: Implement proper AST-based transformer for better reliability
	htmlStr := string(html)

	// For now, return as-is
	// In a production implementation, we'd use regex or HTML parsing
	// to find <a href="http..."> tags and add target="_blank"

	return []byte(htmlStr)
}
