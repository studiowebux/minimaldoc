package parser

import (
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark/ast"
)

// TransformMarkdownLinks walks the AST and converts relative .md links to .html
// It applies slug transformation logic to match the navigation URL structure
func TransformMarkdownLinks(doc ast.Node, currentPagePath string) {
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
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
		transformed := transformLink(destination, currentPagePath)
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
func transformLink(link string, currentPagePath string) string {
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
	slug := generateSlugFromPath(resolvedPath)

	// Construct final .html link as absolute path (starting with /)
	htmlLink := "/" + slug + ".html"
	if fragment != "" {
		htmlLink += "#" + fragment
	}

	return htmlLink
}

// resolveRelativePath converts a relative link to a path from docs root
func resolveRelativePath(linkPath string, currentPagePath string) string {
	// Get directory of current page
	currentDir := filepath.Dir(currentPagePath)

	// Join and clean the path
	resolved := filepath.Join(currentDir, linkPath)
	resolved = filepath.Clean(resolved)

	// Convert back to forward slashes (filepath.Join uses OS separator)
	resolved = strings.ReplaceAll(resolved, "\\", "/")

	return resolved
}

// generateSlugFromPath creates a URL-friendly slug from a file path
// Matches the logic in core.generateSlug
func generateSlugFromPath(relPath string) string {
	// Remove extension
	slug := strings.TrimSuffix(relPath, filepath.Ext(relPath))

	// Remove numbered prefixes from each segment
	parts := strings.Split(slug, "/")
	for i, part := range parts {
		// Remove leading numbers and separators (e.g., "01-", "02_")
		part = strings.TrimLeft(part, "0123456789-_")
		parts[i] = part
	}

	slug = strings.Join(parts, "/")

	// Convert to lowercase and replace spaces
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")

	return slug
}
