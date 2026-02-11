package core

// KBPage represents the complete Knowledge Base page data
type KBPage struct {
	Config        KBConfig
	Categories    []KBCategory
	Articles      []KBArticle // All articles (flat list for search)
	TotalArticles int
}

// KBConfig holds configuration for the Knowledge Base
type KBConfig struct {
	Enabled     bool                     `yaml:"enabled"`
	Title       string                   `yaml:"title"`
	Description string                   `yaml:"description"`
	Path        string                   `yaml:"path"`
	Search      KBSearchConfig           `yaml:"search"`
	Categories  map[string]KBCategoryDef `yaml:"categories"` // Override category display names
}

// KBSearchConfig holds search configuration for the Knowledge Base
type KBSearchConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Placeholder string `yaml:"placeholder"`
}

// KBCategoryDef holds category definition overrides from config
type KBCategoryDef struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Icon        string `yaml:"icon"`
	Order       int    `yaml:"order"`
}

// DefaultKBConfig returns a KBConfig with sensible defaults
func DefaultKBConfig() KBConfig {
	return KBConfig{
		Enabled:     false,
		Title:       "Knowledge Base",
		Description: "Find answers and solutions",
		Path:        "kb",
		Search: KBSearchConfig{
			Enabled:     true,
			Placeholder: "Search knowledge base...",
		},
		Categories: make(map[string]KBCategoryDef),
	}
}

// KBCategory represents a category in the Knowledge Base
type KBCategory struct {
	// Identity
	Slug string `yaml:"-"`

	// Display
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Icon        string `yaml:"icon"`
	Order       int    `yaml:"order"`

	// Content
	Articles     []KBArticle `yaml:"-"`
	ArticleCount int         `yaml:"-"`

	// Output
	OutputPath string `yaml:"-"`
}

// KBArticle represents a single article in the Knowledge Base
type KBArticle struct {
	// Identity
	Slug     string `yaml:"-"`
	FilePath string `yaml:"-"`

	// Metadata from frontmatter
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
	Order       int      `yaml:"order"`

	// Category
	CategorySlug string `yaml:"-"`
	CategoryName string `yaml:"-"`

	// Content
	RawMD string `yaml:"-"`
	HTML  string `yaml:"-"`

	// Navigation
	Prev *KBArticle `yaml:"-"`
	Next *KBArticle `yaml:"-"`

	// Related
	RelatedArticles []KBArticle `yaml:"-"`

	// Output
	OutputPath string `yaml:"-"`
}

// CollectKBTags returns all unique tags from KB articles
func CollectKBTags(articles []KBArticle) []string {
	tagMap := make(map[string]bool)
	for _, article := range articles {
		for _, tag := range article.Tags {
			tagMap[tag] = true
		}
	}

	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	return tags
}

// FindRelatedArticles finds articles related by tags or same category
func FindRelatedArticles(current KBArticle, allArticles []KBArticle, limit int) []KBArticle {
	type scored struct {
		article KBArticle
		score   int
	}

	var candidates []scored
	currentTags := make(map[string]bool)
	for _, tag := range current.Tags {
		currentTags[tag] = true
	}

	for _, article := range allArticles {
		// Skip self
		if article.Slug == current.Slug && article.CategorySlug == current.CategorySlug {
			continue
		}

		score := 0

		// Same category: +1 point
		if article.CategorySlug == current.CategorySlug {
			score += 1
		}

		// Matching tags: +2 points per tag
		for _, tag := range article.Tags {
			if currentTags[tag] {
				score += 2
			}
		}

		if score > 0 {
			candidates = append(candidates, scored{article: article, score: score})
		}
	}

	// Sort by score descending
	for i := 0; i < len(candidates)-1; i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].score > candidates[i].score {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	// Return top N
	var related []KBArticle
	for i := 0; i < len(candidates) && i < limit; i++ {
		related = append(related, candidates[i].article)
	}

	return related
}
