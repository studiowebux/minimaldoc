package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
	"github.com/studiowebux/minimaldoc/internal/parser"
)

const FaqSourceDir = "__faq__"

// FaqBuilder handles building the FAQ page
type FaqBuilder struct {
	frontmatterParser *parser.FrontmatterParser
	markdownParser    *parser.MarkdownParser
}

// NewFaqBuilder creates a new FAQ builder
func NewFaqBuilder() *FaqBuilder {
	return &FaqBuilder{
		frontmatterParser: parser.NewFrontmatterParser(),
		markdownParser:    parser.NewMarkdownParser(),
	}
}

// Build creates the FAQ page from config and markdown files
func (fb *FaqBuilder) Build(docsRoot string, config core.FaqConfig, basePath string) (*core.FaqPage, error) {
	if !config.Enabled {
		return nil, nil
	}

	// Start with categories from config
	categoryMap := make(map[string]*core.FaqCategory)
	for i := range config.Categories {
		cat := &config.Categories[i]
		cat.Slug = generateSlug(cat.Name)
		// Assign slugs to config items
		for j := range cat.Items {
			cat.Items[j].Slug = generateSlug(cat.Items[j].Question)
		}
		categoryMap[cat.Name] = cat
	}

	// Parse markdown files from __faq__/ directory
	faqDir := filepath.Join(docsRoot, FaqSourceDir)
	if _, err := os.Stat(faqDir); err == nil {
		mdItems, err := fb.parseMarkdownFaqs(faqDir, basePath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse FAQ markdown files: %w", err)
		}

		// Merge markdown items into categories
		for _, item := range mdItems {
			categoryName := item.Category
			if categoryName == "" {
				categoryName = "General"
			}

			cat, exists := categoryMap[categoryName]
			if !exists {
				cat = &core.FaqCategory{
					Name: categoryName,
					Slug: generateSlug(categoryName),
				}
				categoryMap[categoryName] = cat
			}
			cat.Items = append(cat.Items, item)
		}
	}

	// Convert map to slice and sort
	var categories []core.FaqCategory
	for _, cat := range categoryMap {
		// Sort items within category by order
		sort.Slice(cat.Items, func(i, j int) bool {
			if cat.Items[i].Order != cat.Items[j].Order {
				return cat.Items[i].Order < cat.Items[j].Order
			}
			return cat.Items[i].Question < cat.Items[j].Question
		})
		categories = append(categories, *cat)
	}

	// Sort categories by order then name
	sort.Slice(categories, func(i, j int) bool {
		if categories[i].Order != categories[j].Order {
			return categories[i].Order < categories[j].Order
		}
		return categories[i].Name < categories[j].Name
	})

	page := &core.FaqPage{
		Config:     config,
		Categories: categories,
	}

	return page, nil
}

// parseMarkdownFaqs parses all markdown files in the FAQ directory
func (fb *FaqBuilder) parseMarkdownFaqs(faqDir string, basePath string) ([]core.FaqItem, error) {
	var items []core.FaqItem

	err := filepath.WalkDir(faqDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		item, err := fb.parseMarkdownFaq(path, faqDir, basePath)
		if err != nil {
			fmt.Printf("Warning: failed to parse FAQ file %s: %v\n", path, err)
			return nil
		}

		items = append(items, *item)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

// parseMarkdownFaq parses a single FAQ markdown file
func (fb *FaqBuilder) parseMarkdownFaq(filePath, faqDir, basePath string) (*core.FaqItem, error) {
	// Parse frontmatter
	meta, content, err := fb.frontmatterParser.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("frontmatter parse error: %w", err)
	}

	// Parse markdown to HTML
	html, err := fb.markdownParser.ParseWithContext(content, "", basePath)
	if err != nil {
		return nil, fmt.Errorf("markdown parse error: %w", err)
	}

	// Extract category from directory structure or frontmatter
	relPath, _ := filepath.Rel(faqDir, filePath)
	dirCategory := filepath.Dir(relPath)
	if dirCategory == "." {
		dirCategory = ""
	}

	category := meta.Category
	if category == "" && dirCategory != "" {
		// Use directory name as category, capitalize first letter
		category = capitalizeFirst(dirCategory)
	}

	// Generate slug from filename
	base := filepath.Base(filePath)
	slug := base[:len(base)-len(filepath.Ext(base))]

	// Use frontmatter question or fallback to title
	question := meta.Question
	if question == "" {
		question = meta.Title
	}

	item := &core.FaqItem{
		Question:   question,
		Answer:     string(content),
		AnswerHTML: string(html),
		Slug:       generateSlug(question),
		Order:      meta.MenuOrder,
		Tags:       meta.Tags,
		FilePath:   filePath,
		Category:   category,
	}

	// Use slug from filename if question slug would be empty
	if item.Slug == "" {
		item.Slug = slug
	}

	return item, nil
}

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// generateSlug creates a URL-friendly slug from text
func generateSlug(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters, keep only alphanumeric and hyphens
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	// Clean up multiple hyphens
	slug = result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}
