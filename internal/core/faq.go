package core

import "sort"

// FaqPage represents the complete FAQ page data
type FaqPage struct {
	Config     FaqConfig
	Categories []FaqCategory
}

// FaqConfig holds configuration for the FAQ page
type FaqConfig struct {
	Enabled     bool          `yaml:"enabled"`
	Title       string        `yaml:"title"`
	Description string        `yaml:"description"`
	Path        string        `yaml:"path"`
	Categories  []FaqCategory `yaml:"categories"`
}

// DefaultFaqConfig returns a FaqConfig with sensible defaults
func DefaultFaqConfig() FaqConfig {
	return FaqConfig{
		Enabled:     false,
		Title:       "FAQ",
		Description: "Frequently asked questions",
		Path:        "faq",
	}
}

// FaqCategory represents a category grouping of FAQ items
type FaqCategory struct {
	Name  string    `yaml:"name"`
	Slug  string    `yaml:"-"`
	Order int       `yaml:"order"`
	Items []FaqItem `yaml:"items"`
}

// FaqItem represents a single FAQ question/answer pair
type FaqItem struct {
	Question string   `yaml:"question"`
	Answer   string   `yaml:"answer"`
	Slug     string   `yaml:"-"`
	Order    int      `yaml:"order"`
	Tags     []string `yaml:"tags"`

	// For markdown-based items
	AnswerHTML string `yaml:"-"`
	FilePath   string `yaml:"-"`
	Category   string `yaml:"category"`
}

// CollectFaqTags returns all unique tags from FAQ items
func CollectFaqTags(categories []FaqCategory) []string {
	tagMap := make(map[string]bool)
	for _, cat := range categories {
		for _, item := range cat.Items {
			for _, tag := range item.Tags {
				tagMap[tag] = true
			}
		}
	}

	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}
