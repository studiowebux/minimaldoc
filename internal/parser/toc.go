package parser

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/toc"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// TOCParser generates table of contents from markdown
type TOCParser struct {
	md goldmark.Markdown
}

// NewTOCParser creates a new TOC parser
func NewTOCParser() *TOCParser {
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	return &TOCParser{md: md}
}

// Parse generates a table of contents from markdown content
func (p *TOCParser) Parse(content []byte) (*core.TOC, error) {
	// Parse markdown to AST
	doc := p.md.Parser().Parse(text.NewReader(content))

	// Extract TOC
	tree, err := toc.Inspect(doc, content)
	if err != nil {
		return nil, err
	}

	// Convert goldmark TOC to our core.TOC structure
	coreTOC := convertToCoreTOC(tree)

	return coreTOC, nil
}

// convertToCoreTOC converts goldmark's TOC to our core.TOC structure
func convertToCoreTOC(tree *toc.TOC) *core.TOC {
	if tree == nil {
		return &core.TOC{Items: []*core.TOCItem{}}
	}

	items := make([]*core.TOCItem, len(tree.Items))
	for i, item := range tree.Items {
		items[i] = convertTOCItemWithLevel(item, 1)
	}

	return &core.TOC{Items: items}
}

// convertTOCItemWithLevel recursively converts a TOC item and its children, tracking level
func convertTOCItemWithLevel(item *toc.Item, level int) *core.TOCItem {
	tocItem := &core.TOCItem{
		ID:       string(item.ID),
		Title:    string(item.Title),
		Level:    level,
		Children: make([]*core.TOCItem, len(item.Items)),
	}

	// Recursively convert children with incremented level
	for i, child := range item.Items {
		tocItem.Children[i] = convertTOCItemWithLevel(child, level+1)
	}

	return tocItem
}
