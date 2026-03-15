package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
	"github.com/studiowebux/minimaldoc/internal/parser"
)

// NavigationBuilder builds the site navigation tree from pages
type NavigationBuilder struct{}

// pathGroup is a helper structure for building the navigation tree
type pathGroup struct {
	path     string
	pages    []*core.Page
	children map[string]*pathGroup
}

// NewNavigationBuilder creates a new navigation builder
func NewNavigationBuilder() *NavigationBuilder {
	return &NavigationBuilder{}
}

// Build creates a navigation tree from a list of pages
// If TOC.md exists, it will be used to define the navigation structure
func (b *NavigationBuilder) Build(pages []*core.Page, docsDir string, maxDepth int) *core.Navigation {
	// Check if TOC.md exists
	tocPath := filepath.Join(docsDir, "TOC.md")
	if _, err := os.Stat(tocPath); err == nil {
		// Use TOC.md-based navigation
		return b.buildFromTOC(pages, docsDir, tocPath, maxDepth)
	}

	// Fall back to folder-based navigation
	// Filter out hidden pages
	visiblePages := make([]*core.Page, 0)
	for _, page := range pages {
		if !page.IsHidden() {
			visiblePages = append(visiblePages, page)
		}
	}

	// Build tree structure
	rootItems := b.buildTree(visiblePages, maxDepth)

	return &core.Navigation{Items: rootItems}
}

// buildFromTOC creates navigation from TOC.md file
func (b *NavigationBuilder) buildFromTOC(pages []*core.Page, docsDir string, tocPath string, maxDepth int) *core.Navigation {
	tocParser := parser.NewTOCFileParser(docsDir)
	entries, err := tocParser.Parse(tocPath)
	if err != nil {
		// Fall back to folder-based navigation on error
		return b.Build(pages, docsDir, maxDepth)
	}

	// Create a map of file paths to pages for quick lookup
	pageMap := make(map[string]*core.Page)
	for _, page := range pages {
		if !page.IsHidden() {
			// Store both absolute and relative paths
			pageMap[page.SourcePath] = page
			relPath, _ := filepath.Rel(docsDir, page.SourcePath)
			pageMap[relPath] = page
			// Also store with normalized path
			pageMap[filepath.ToSlash(relPath)] = page
		}
	}

	// Convert TOC entries to NavItems
	rootItems := b.tocEntriesToNavItems(entries, pageMap, 0, maxDepth)

	return &core.Navigation{Items: rootItems}
}

// tocEntriesToNavItems converts TOC entries to navigation items
func (b *NavigationBuilder) tocEntriesToNavItems(entries []*parser.TOCEntry, pageMap map[string]*core.Page, depth int, maxDepth int) []*core.NavItem {
	if maxDepth > 0 && depth >= maxDepth {
		return nil
	}

	items := make([]*core.NavItem, 0, len(entries))

	for i, entry := range entries {
		item := &core.NavItem{
			Title:      entry.Title,
			Order:      i, // Use position in TOC as order
			IsExternal: entry.IsExternal,
			Children:   []*core.NavItem{},
		}

		// If entry has a file path, find the corresponding page
		if entry.FilePath != "" {
			if entry.IsExternal {
				// External URL - use as-is
				item.Path = entry.FilePath
			} else {
				// Internal file path - try to find the page in the map
				if page, found := pageMap[entry.FilePath]; found {
					item.Path = "/" + page.Slug + ".html"
					item.Page = page
				} else {
					// Try normalized path
					normalized := filepath.ToSlash(entry.FilePath)
					if page, found := pageMap[normalized]; found {
						item.Path = "/" + page.Slug + ".html"
						item.Page = page
					}
				}
			}
		}

		// Process children recursively
		if len(entry.Children) > 0 {
			item.Children = b.tocEntriesToNavItems(entry.Children, pageMap, depth+1, maxDepth)
		}

		items = append(items, item)
	}

	return items
}

// buildTree creates a hierarchical navigation tree
func (b *NavigationBuilder) buildTree(pages []*core.Page, maxDepth int) []*core.NavItem {

	root := &pathGroup{
		path:     "",
		pages:    []*core.Page{},
		children: make(map[string]*pathGroup),
	}

	// Organize pages into tree structure based on file paths
	for _, page := range pages {
		parts := strings.Split(page.RelPath, string(filepath.Separator))
		current := root

		// Navigate/create path groups
		for i := 0; i < len(parts)-1; i++ {
			part := parts[i]
			if _, exists := current.children[part]; !exists {
				current.children[part] = &pathGroup{
					path:     part,
					pages:    []*core.Page{},
					children: make(map[string]*pathGroup),
				}
			}
			current = current.children[part]
		}

		// Add page to its parent group
		current.pages = append(current.pages, page)
	}

	// Convert tree to NavItems
	return b.groupToNavItems(root, 0, maxDepth)
}

// groupToNavItems converts a pathGroup tree to NavItems
func (b *NavigationBuilder) groupToNavItems(group *pathGroup, depth int, maxDepth int) []*core.NavItem {
	if maxDepth > 0 && depth >= maxDepth {
		return nil
	}

	items := []*core.NavItem{}

	// Add pages at this level
	for _, page := range group.pages {
		item := &core.NavItem{
			Title:    page.Title(),
			Path:     "/" + page.Slug + ".html",
			Order:    page.Order,
			Children: []*core.NavItem{},
			Page:     page,
		}

		items = append(items, item)
	}

	// Add child groups
	for _, child := range group.children {
		childItems := b.groupToNavItems(child, depth+1, maxDepth)

		// If child group has items, create a section
		if len(childItems) > 0 {
			// Use first page as section header, or create a synthetic one
			sectionTitle := cleanPathSegment(child.path)

			sectionItem := &core.NavItem{
				Title:    sectionTitle,
				Path:     "",
				Order:    extractOrder(child.path),
				Children: childItems,
			}

			items = append(items, sectionItem)
		}
	}

	// Sort items by order
	sort.Slice(items, func(i, j int) bool {
		if items[i].Order != items[j].Order {
			return items[i].Order < items[j].Order
		}
		return items[i].Title < items[j].Title
	})

	// Warn about duplicate labels at this level
	seen := make(map[string]string) // label -> first directory path
	for _, child := range group.children {
		label := cleanPathSegment(child.path)
		if prev, exists := seen[label]; exists {
			fmt.Fprintf(os.Stderr, "Warning: navigation label collision — directories %q and %q both produce label %q. Rename one to avoid duplicate sidebar entries.\n", prev, child.path, label)
		} else {
			seen[label] = child.path
		}
	}

	return items
}

// extractOrder extracts the order number from a filename or directory name
// Examples: "01-intro" -> 1, "02-guide" -> 2, "guide" -> 999
func extractOrder(name string) int {
	re := regexp.MustCompile(`^(\d+)[-_]`)
	matches := re.FindStringSubmatch(filepath.Base(name))

	if len(matches) > 1 {
		order, err := strconv.Atoi(matches[1])
		if err == nil {
			return order
		}
	}

	return 999 // Default order for items without numbers
}

// cleanPathSegment removes number prefixes and cleans up a path segment for display
// Examples: "01-getting-started" -> "Getting Started", "02_API" -> "API"
func cleanPathSegment(segment string) string {
	// Remove number prefix
	re := regexp.MustCompile(`^(\d+)[-_]`)
	clean := re.ReplaceAllString(segment, "")

	// Replace separators with spaces
	clean = strings.ReplaceAll(clean, "-", " ")
	clean = strings.ReplaceAll(clean, "_", " ")

	// Title case
	words := strings.Fields(clean)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}

	return strings.Join(words, " ")
}
