package parser

import (
	"path/filepath"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
	"github.com/yuin/goldmark/ast"
)

// TransformMarkdownLinks walks the AST and converts relative .md links to .html
// It applies slug transformation logic to match the navigation URL structure
// basePath is prepended to all transformed links (e.g., "/docs")
func TransformMarkdownLinks(doc ast.Node, currentPagePath string, basePath string) {
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) { // walk never errors; closure always returns nil
		if !entering {
			return ast.WalkContinue, nil
		}

		// Only process link nodes
		link, ok := n.(*ast.Link)
		if !ok {
			return ast.WalkContinue, nil
		}

		destination := string(link.Destination)

		// Skip non-.md links and external URLs
		if !isInternalMarkdownLink(destination) {
			return ast.WalkContinue, nil
		}

		// Transform the link
		transformed := transformLink(destination, currentPagePath, basePath)
		link.Destination = []byte(transformed)

		return ast.WalkContinue, nil
	})
}

// isInternalMarkdownLink checks if a link is a relative .md link
func isInternalMarkdownLink(link string) bool {
	// Skip external URLs
	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") || strings.HasPrefix(link, "//") {
		return false
	}

	// Skip absolute paths (starting with /)
	if strings.HasPrefix(link, "/") {
		return false
	}

	// Check if it's a .md link (before any fragment)
	linkWithoutFragment := strings.Split(link, "#")[0]
	return strings.HasSuffix(linkWithoutFragment, ".md")
}

// transformLink converts a .md link to .html with slug transformation
// basePath is prepended to the final link (e.g., "/docs")
func transformLink(link string, currentPagePath string, basePath string) string {
	// Split link into path and fragment
	parts := strings.SplitN(link, "#", 2)
	mdPath := parts[0]
	var fragment string
	if len(parts) > 1 {
		fragment = parts[1]
	}

	// Resolve relative path to absolute path from docs root
	resolvedPath := resolveRelativePath(mdPath, currentPagePath)

	// Apply slug transformation (remove .md, number prefixes, etc.)
	slug := core.GenerateSlugFromPath(resolvedPath)

	// Construct final .html link with basePath prepended
	htmlLink := basePath + "/" + slug + ".html"
	if fragment != "" {
		htmlLink += "#" + fragment
	}

	return htmlLink
}

// resolveRelativePath converts a relative link to a path from docs root
func resolveRelativePath(linkPath string, currentPagePath string) string {
	// Get directory of current page
	currentDir := filepath.Dir(currentPagePath)

	// Join, clean, then normalise to forward slashes for URL construction.
	resolved := filepath.ToSlash(filepath.Clean(filepath.Join(currentDir, linkPath)))

	return resolved
}

// TransformExternalLinks adds target="_blank" and rel="noopener noreferrer" to external links
func TransformExternalLinks(doc ast.Node) {
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) { // walk never errors; closure always returns nil
		if !entering {
			return ast.WalkContinue, nil
		}

		link, ok := n.(*ast.Link)
		if !ok {
			return ast.WalkContinue, nil
		}

		destination := string(link.Destination)

		// Check if external (http:// or https://)
		if strings.HasPrefix(destination, "http://") || strings.HasPrefix(destination, "https://") {
			link.SetAttributeString("target", "_blank")
			link.SetAttributeString("rel", "noopener noreferrer")
		}

		return ast.WalkContinue, nil
	})
}
