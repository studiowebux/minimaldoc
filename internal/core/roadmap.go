package core

// RoadmapColumn represents a column in the board layout
type RoadmapColumn struct {
	ID    string `yaml:"id"`
	Label string `yaml:"label"`
}

// RoadmapItem represents a single roadmap entry
type RoadmapItem struct {
	Title        string   `yaml:"title"`
	Description  string   `yaml:"description"`
	Status       string   `yaml:"status"`
	Version      string   `yaml:"version"`
	Tags         []string `yaml:"tags"`
	ShippedDate  string   `yaml:"shipped_date"`
	ChangelogURL string   `yaml:"changelog_url"`
}

// RoadmapConfig holds configuration for the roadmap page
type RoadmapConfig struct {
	Enabled      bool            `yaml:"enabled"`
	Title        string          `yaml:"title"`
	Description  string          `yaml:"description"`
	Layout       string          `yaml:"layout"` // "board" or "timeline"
	ShowVersions bool            `yaml:"show_versions"`
	Path         string          `yaml:"path"`
	Columns      []RoadmapColumn `yaml:"columns"`
	Items        []RoadmapItem   `yaml:"items"`
}

// RoadmapPage represents the roadmap page data passed to the template
type RoadmapPage struct {
	Config RoadmapConfig
}

// DefaultRoadmapConfig returns a RoadmapConfig with sensible defaults
func DefaultRoadmapConfig() RoadmapConfig {
	return RoadmapConfig{
		Enabled:      false,
		Title:        "Roadmap",
		Description:  "What we're building and where we're headed.",
		Layout:       "board",
		ShowVersions: true,
		Path:         "roadmap",
		Columns: []RoadmapColumn{
			{ID: "planned", Label: "Planned"},
			{ID: "in_progress", Label: "In Progress"},
			{ID: "shipped", Label: "Shipped"},
		},
	}
}
