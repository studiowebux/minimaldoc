package generator

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// KBGenerator generates the Knowledge Base pages
type KBGenerator struct {
	site      *core.Site
	templates *template.Template
	themeFS   embed.FS
	version   string
}

// NewKBGenerator creates a new Knowledge Base generator
func NewKBGenerator(site *core.Site, themeFS embed.FS, version string) (*KBGenerator, error) {
	if !site.Config.KnowledgeBase.Enabled {
		return nil, nil
	}

	tmpl := template.New("").Funcs(KBFuncMap()).Funcs(AnalyticsFuncMap())

	var err error
	tmpl, err = tmpl.ParseFS(
		themeFS,
		"themes/common/templates/partials/landing-*.html",
		"themes/common/templates/partials/analytics.html",
		"themes/common/templates/partials/minimaldoc-widgets.html",
		"themes/common/templates/kb/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse KB templates: %w", err)
	}

	return &KBGenerator{
		site:      site,
		templates: tmpl,
		themeFS:   themeFS,
		version:   version,
	}, nil
}

// Generate generates all Knowledge Base pages
func (g *KBGenerator) Generate() error {
	if g.site.KBPage == nil || !g.site.Config.KnowledgeBase.Enabled {
		return nil
	}

	fmt.Println("Generating Knowledge Base...")

	kbPath := g.site.Config.KnowledgeBase.Path
	if kbPath == "" {
		kbPath = "kb"
	}

	outputDir := filepath.Join(g.site.OutputRoot, kbPath)
	if err := makeWebDir(outputDir); err != nil {
		return fmt.Errorf("failed to create KB output directory: %w", err)
	}

	// Generate landing page
	if err := g.generateLanding(outputDir); err != nil {
		return fmt.Errorf("failed to generate KB landing: %w", err)
	}

	// Generate category pages
	for _, cat := range g.site.KBPage.Categories {
		if err := g.generateCategory(outputDir, cat); err != nil {
			return fmt.Errorf("failed to generate KB category %s: %w", cat.Slug, err)
		}
	}

	// Generate KB search index
	if g.site.Config.KnowledgeBase.Search.Enabled {
		if err := g.generateSearchIndex(outputDir); err != nil {
			return fmt.Errorf("failed to generate KB search index: %w", err)
		}
	}

	fmt.Printf("Generated Knowledge Base: %d categories, %d articles\n",
		len(g.site.KBPage.Categories), g.site.KBPage.TotalArticles)
	return nil
}

// generateLanding generates the KB landing page
func (g *KBGenerator) generateLanding(outputDir string) error {
	kbPath := g.site.Config.KnowledgeBase.Path
	if kbPath == "" {
		kbPath = "kb"
	}

	data := map[string]any{
		"Site":       g.site,
		"KBPage":     g.site.KBPage,
		"Footer":     BuildFooter(g.site, g.version),
		"BasePath":   g.getBasePath(),
		"KBPath":     kbPath,
		"Version":    g.version,
		"PageTitle":  g.site.KBPage.Config.Title + " | " + g.site.Config.Title,
		"ActivePath": "/" + kbPath + "/",
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "kb-landing.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	outputPath := filepath.Join(outputDir, "index.html")
	if err := writeWebFile(outputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write KB landing: %w", err)
	}

	return nil
}

// generateCategory generates a category page and its articles
func (g *KBGenerator) generateCategory(outputDir string, cat core.KBCategory) error {
	kbPath := g.site.Config.KnowledgeBase.Path
	if kbPath == "" {
		kbPath = "kb"
	}

	catDir := filepath.Join(outputDir, cat.Slug)
	if err := makeWebDir(catDir); err != nil {
		return err
	}

	// Generate category index page
	data := map[string]any{
		"Site":       g.site,
		"KBPage":     g.site.KBPage,
		"Category":   cat,
		"Footer":     BuildFooter(g.site, g.version),
		"BasePath":   g.getBasePath(),
		"KBPath":     kbPath,
		"Version":    g.version,
		"PageTitle":  cat.Name + " | " + g.site.KBPage.Config.Title + " | " + g.site.Config.Title,
		"ActivePath": "/" + kbPath + "/",
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "kb-category.html", data); err != nil {
		return fmt.Errorf("category template execution failed: %w", err)
	}

	outputPath := filepath.Join(catDir, "index.html")
	if err := writeWebFile(outputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write category page: %w", err)
	}

	// Generate article pages
	for _, article := range cat.Articles {
		if err := g.generateArticle(catDir, cat, article); err != nil {
			return fmt.Errorf("failed to generate article %s: %w", article.Slug, err)
		}
	}

	return nil
}

// generateArticle generates a single article page
func (g *KBGenerator) generateArticle(catDir string, cat core.KBCategory, article core.KBArticle) error {
	kbPath := g.site.Config.KnowledgeBase.Path
	if kbPath == "" {
		kbPath = "kb"
	}

	data := map[string]any{
		"Site":       g.site,
		"KBPage":     g.site.KBPage,
		"Category":   cat,
		"Article":    article,
		"Footer":     BuildFooter(g.site, g.version),
		"BasePath":   g.getBasePath(),
		"KBPath":     kbPath,
		"Version":    g.version,
		"PageTitle":  article.Title + " | " + cat.Name + " | " + g.site.Config.Title,
		"ActivePath": "/" + kbPath + "/",
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "kb-article.html", data); err != nil {
		return fmt.Errorf("article template execution failed: %w", err)
	}

	outputPath := filepath.Join(catDir, article.Slug+".html")
	if err := writeWebFile(outputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write article page: %w", err)
	}

	return nil
}

// generateSearchIndex generates the KB-specific search index
func (g *KBGenerator) generateSearchIndex(outputDir string) error {
	type searchEntry struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		URL         string   `json:"url"`
		Category    string   `json:"category"`
		Tags        []string `json:"tags"`
		Content     string   `json:"content"`
	}

	kbPath := g.site.Config.KnowledgeBase.Path
	if kbPath == "" {
		kbPath = "kb"
	}

	var entries []searchEntry
	for _, cat := range g.site.KBPage.Categories {
		for _, article := range cat.Articles {
			// Strip HTML and truncate content for search
			content := stripHTML(article.HTML)
			if len(content) > 500 {
				content = content[:500]
			}

			entries = append(entries, searchEntry{
				Title:       article.Title,
				Description: article.Description,
				URL:         g.getBasePath() + "/" + kbPath + "/" + cat.Slug + "/" + article.Slug + ".html",
				Category:    cat.Name,
				Tags:        article.Tags,
				Content:     content,
			})
		}
	}

	// Write as JSON
	var buf bytes.Buffer
	buf.WriteString("[\n")
	for i, entry := range entries {
		if i > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString(fmt.Sprintf(`{"title":%q,"description":%q,"url":%q,"category":%q,"tags":[`,
			entry.Title, entry.Description, entry.URL, entry.Category))
		for j, tag := range entry.Tags {
			if j > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(fmt.Sprintf("%q", tag))
		}
		buf.WriteString(fmt.Sprintf(`],"content":%q}`, entry.Content))
	}
	buf.WriteString("\n]")

	outputPath := filepath.Join(outputDir, "kb-search.json")
	if err := writeWebFile(outputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write KB search index: %w", err)
	}

	return nil
}

// stripHTML removes HTML tags from a string
func stripHTML(html string) string {
	var result strings.Builder
	inTag := false
	for _, r := range html {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			result.WriteRune(' ')
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}
	// Clean up multiple spaces
	text := result.String()
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}
	return strings.TrimSpace(text)
}

// getBasePath extracts the path component from BaseURL
func (g *KBGenerator) getBasePath() string {
	return GetBasePath(g.site.Config.BaseURL)
}

// KBFuncMap returns template functions specific to KB pages.
func KBFuncMap() template.FuncMap {
	return ExtendFuncMap(template.FuncMap{
		"kbArticleURL": func(basePath, kbPath string, article core.KBArticle) string {
			return basePath + "/" + kbPath + "/" + article.CategorySlug + "/" + article.Slug + ".html"
		},
		"kbCategoryURL": func(basePath, kbPath string, cat core.KBCategory) string {
			return basePath + "/" + kbPath + "/" + cat.Slug + "/"
		},
	})
}
