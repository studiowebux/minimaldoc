package parser

import (
	"bytes"
	stdhtml "html"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// AdmonitionKind represents the type of admonition
var AdmonitionKind = ast.NewNodeKind("Admonition")

// Admonition is a custom block node for admonitions (:::info, :::warning, etc.)
type Admonition struct {
	ast.BaseBlock
	AdmonType string
	Title     string
}

// Kind implements ast.Node.Kind
func (n *Admonition) Kind() ast.NodeKind {
	return AdmonitionKind
}

// Dump implements ast.Node.Dump
func (n *Admonition) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"AdmonType": n.AdmonType,
		"Title":     n.Title,
	}, nil)
}

// admonitionParser is a block parser for admonitions
type admonitionParser struct{}

var admonitionPattern = regexp.MustCompile(`^:::(\w+)\s*(.*)$`)

// NewAdmonitionParser returns a new admonition parser
func NewAdmonitionParser() parser.BlockParser {
	return &admonitionParser{}
}

// Trigger implements parser.BlockParser.Trigger
func (p *admonitionParser) Trigger() []byte {
	return []byte{':'}
}

// Open implements parser.BlockParser.Open
func (p *admonitionParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()

	// Check if line starts with :::
	if !bytes.HasPrefix(bytes.TrimSpace(line), []byte(":::")) {
		return nil, parser.NoChildren
	}

	// Trim the line to remove trailing whitespace/newlines
	trimmedLine := bytes.TrimSpace(line)
	matches := admonitionPattern.FindSubmatch(trimmedLine)

	if matches == nil {
		return nil, parser.NoChildren
	}

	admonType := string(matches[1])
	title := strings.TrimSpace(string(matches[2]))

	if title == "" {
		// Default titles based on type
		titleMap := map[string]string{
			"info":     "Info",
			"warning":  "Warning",
			"question": "Question",
			"danger":   "Danger",
			"success":  "Success",
			"note":     "Note",
		}
		if defaultTitle, ok := titleMap[admonType]; ok {
			title = defaultTitle
		} else {
			title = admonType
		}
	}

	reader.Advance(len(line))

	return &Admonition{
		AdmonType: admonType,
		Title:     title,
	}, parser.HasChildren
}

// Continue implements parser.BlockParser.Continue
func (p *admonitionParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()

	// Check for closing :::
	trimmed := bytes.TrimSpace(line)
	if len(trimmed) >= 3 && bytes.HasPrefix(trimmed, []byte(":::")) && (len(trimmed) == 3 || !isAlphanumeric(trimmed[3])) {
		reader.Advance(segment.Len())
		return parser.Close
	}

	return parser.Continue | parser.HasChildren
}

// isAlphanumeric checks if a byte is alphanumeric
func isAlphanumeric(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

// Close implements parser.BlockParser.Close
func (p *admonitionParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// Nothing to do on close
}

// CanInterruptParagraph implements parser.BlockParser.CanInterruptParagraph
func (p *admonitionParser) CanInterruptParagraph() bool {
	return true
}

// CanAcceptIndentedLine implements parser.BlockParser.CanAcceptIndentedLine
func (p *admonitionParser) CanAcceptIndentedLine() bool {
	return false
}

// admonitionHTMLRenderer is a renderer for Admonition nodes
type admonitionHTMLRenderer struct {
	html.Config
}

// NewAdmonitionHTMLRenderer returns a new admonition HTML renderer
func NewAdmonitionHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &admonitionHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs
func (r *admonitionHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(AdmonitionKind, r.renderAdmonition)
}

// renderAdmonition renders an Admonition node to HTML
func (r *admonitionHTMLRenderer) renderAdmonition(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*Admonition)

	if entering {
		// Opening div with class
		_, _ = w.WriteString(`<div class="admonition `)
		_, _ = w.WriteString(n.AdmonType)
		_, _ = w.WriteString(`">`)

		// Title (escaped to prevent XSS)
		_, _ = w.WriteString(`<div class="admonition-title">`)
		_, _ = w.WriteString(stdhtml.EscapeString(n.Title))
		_, _ = w.WriteString(`</div>`)

		// Content wrapper
		_, _ = w.WriteString(`<div class="admonition-content">`)
	} else {
		// Close content wrapper and main div
		_, _ = w.WriteString(`</div></div>`)
	}

	return ast.WalkContinue, nil
}

// AdmonitionExtension is a Goldmark extension for admonitions
type AdmonitionExtension struct{}

// Extend implements goldmark.Extender
func (e *AdmonitionExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(NewAdmonitionParser(), 999),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewAdmonitionHTMLRenderer(), 999),
		),
	)
}

// NewAdmonitionExtension returns a new AdmonitionExtension
func NewAdmonitionExtension() goldmark.Extender {
	return &AdmonitionExtension{}
}
