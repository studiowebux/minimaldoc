package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/studiowebux/minimaldoc/static-generator/internal/core"
	"github.com/studiowebux/minimaldoc/static-generator/internal/parser"
)

// KBBuilder handles building the Knowledge Base pages
type KBBuilder struct {
	frontmatterParser *parser.FrontmatterParser
	markdownParser    *parser.MarkdownParser
}

// NewKBBuilder creates a new Knowledge Base builder
func NewKBBuilder() *KBBuilder {
	return &KBBuilder{
		frontmatterParser: parser.NewFrontmatterParser(),
		markdownParser:    parser.NewMarkdownParser(),
	}
}

// Build creates the Knowledge Base from the __kb__ directory
func (kb *KBBuilder) Build(docsRoot string, config core.KBConfig, basePath string) (*core.KBPage, error) {
	if !config.Enabled {
		return nil, nil
	}

	kbDir := filepath.Join(docsRoot, core.KBSourceDir)
	if _, err := os.Stat(kbDir); os.IsNotExist(err) {
		// No KB directory, return empty page
		return &core.KBPage{
			Config:     config,
			Categories: []core.KBCategory{},
			Articles:   []core.KBArticle{},
		}, nil
	}

	// Scan for categories (subdirectories)
	categories, err := kb.scanCategories(kbDir, config, basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to scan KB categories: %w", err)
	}

	// Collect all articles flat
	var allArticles []core.KBArticle
	for _, cat := range categories {
		allArticles = append(allArticles, cat.Articles...)
	}

	// Compute prev/next within each category
	for i := range categories {
		kb.computePrevNext(&categories[i])
	}

	// Compute related articles
	for i := range categories {
		for j := range categories[i].Articles {
			categories[i].Articles[j].RelatedArticles = core.FindRelatedArticles(
				categories[i].Articles[j],
				allArticles,
				3,
			)
		}
	}

	page := &core.KBPage{
		Config:        config,
		Categories:    categories,
		Articles:      allArticles,
		TotalArticles: len(allArticles),
	}

	return page, nil
}

// scanCategories scans the KB directory for category subdirectories
func (kb *KBBuilder) scanCategories(kbDir string, config core.KBConfig, basePath string) ([]core.KBCategory, error) {
	entries, err := os.ReadDir(kbDir)
	if err != nil {
		return nil, err
	}

	var categories []core.KBCategory

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		slug := entry.Name()
		catDir := filepath.Join(kbDir, slug)

		// Parse articles in category
		articles, err := kb.parseArticles(catDir, slug, basePath)
		if err != nil {
			fmt.Printf("Warning: failed to parse KB category %s: %v\n", slug, err)
			continue
		}

		// Skip empty categories
		if len(articles) == 0 {
			continue
		}

		// Build category
		cat := core.KBCategory{
			Slug:         slug,
			Name:         kb.getCategoryName(slug, config),
			Description:  kb.getCategoryDescription(slug, config),
			Icon:         kb.getCategoryIcon(slug, config),
			Order:        kb.getCategoryOrder(slug, config),
			Articles:     articles,
			ArticleCount: len(articles),
		}

		categories = append(categories, cat)
	}

	// Sort categories by order, then by name
	sort.Slice(categories, func(i, j int) bool {
		if categories[i].Order != categories[j].Order {
			return categories[i].Order < categories[j].Order
		}
		return categories[i].Name < categories[j].Name
	})

	return categories, nil
}

// parseArticles parses all markdown files in a category directory
func (kb *KBBuilder) parseArticles(catDir, categorySlug, basePath string) ([]core.KBArticle, error) {
	entries, err := os.ReadDir(catDir)
	if err != nil {
		return nil, err
	}

	var articles []core.KBArticle

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		filePath := filepath.Join(catDir, entry.Name())
		article, err := kb.parseArticle(filePath, categorySlug, basePath)
		if err != nil {
			fmt.Printf("Warning: failed to parse KB article %s: %v\n", filePath, err)
			continue
		}

		articles = append(articles, *article)
	}

	// Sort articles by order, then by title
	sort.Slice(articles, func(i, j int) bool {
		if articles[i].Order != articles[j].Order {
			return articles[i].Order < articles[j].Order
		}
		return articles[i].Title < articles[j].Title
	})

	return articles, nil
}

// parseArticle parses a single KB article markdown file
func (kb *KBBuilder) parseArticle(filePath, categorySlug, basePath string) (*core.KBArticle, error) {
	// Parse frontmatter
	meta, content, err := kb.frontmatterParser.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("frontmatter parse error: %w", err)
	}

	// Parse markdown to HTML
	html, err := kb.markdownParser.ParseWithContext(content, "", basePath)
	if err != nil {
		return nil, fmt.Errorf("markdown parse error: %w", err)
	}

	// Generate slug from filename
	base := filepath.Base(filePath)
	slug := base[:len(base)-len(filepath.Ext(base))]

	// Strip numeric prefix from slug (e.g., "01-quick-start" -> "quick-start")
	slug = stripNumericPrefix(slug)

	// Extract order from filename or frontmatter
	order := extractOrder(base)
	if meta.MenuOrder >= 0 {
		order = meta.MenuOrder
	}

	article := &core.KBArticle{
		Slug:         slug,
		FilePath:     filePath,
		Title:        meta.Title,
		Description:  meta.Description,
		Tags:         meta.Tags,
		Order:        order,
		CategorySlug: categorySlug,
		CategoryName: capitalizeFirst(strings.ReplaceAll(categorySlug, "-", " ")),
		RawMD:        string(content),
		HTML:         string(html),
	}

	// Fallback title from filename if not in frontmatter
	if article.Title == "" {
		article.Title = capitalizeFirst(strings.ReplaceAll(slug, "-", " "))
	}

	return article, nil
}

// computePrevNext sets prev/next links for articles within a category
func (kb *KBBuilder) computePrevNext(cat *core.KBCategory) {
	for i := range cat.Articles {
		if i > 0 {
			cat.Articles[i].Prev = &cat.Articles[i-1]
		}
		if i < len(cat.Articles)-1 {
			cat.Articles[i].Next = &cat.Articles[i+1]
		}
	}
}

// getCategoryName returns the display name for a category
func (kb *KBBuilder) getCategoryName(slug string, config core.KBConfig) string {
	if def, ok := config.Categories[slug]; ok && def.Name != "" {
		return def.Name
	}
	return capitalizeFirst(strings.ReplaceAll(slug, "-", " "))
}

// getCategoryDescription returns the description for a category
func (kb *KBBuilder) getCategoryDescription(slug string, config core.KBConfig) string {
	if def, ok := config.Categories[slug]; ok {
		return def.Description
	}
	return ""
}

// getCategoryIcon returns the icon for a category
func (kb *KBBuilder) getCategoryIcon(slug string, config core.KBConfig) string {
	if def, ok := config.Categories[slug]; ok && def.Icon != "" {
		return def.Icon
	}
	return "folder"
}

// getCategoryOrder returns the order for a category
func (kb *KBBuilder) getCategoryOrder(slug string, config core.KBConfig) int {
	if def, ok := config.Categories[slug]; ok {
		return def.Order
	}
	return 999
}

// stripNumericPrefix removes numeric prefix from filename slug
// e.g., "01-quick-start" -> "quick-start", "10-advanced" -> "advanced"
func stripNumericPrefix(slug string) string {
	parts := strings.SplitN(slug, "-", 2)
	if len(parts) == 2 {
		// Check if first part is numeric
		isNumeric := true
		for _, r := range parts[0] {
			if r < '0' || r > '9' {
				isNumeric = false
				break
			}
		}
		if isNumeric {
			return parts[1]
		}
	}
	return slug
}
