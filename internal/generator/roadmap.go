package generator

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// RoadmapVersionGroup holds items grouped by version for template rendering
type RoadmapVersionGroup struct {
	Version string
	Status  string // aggregate: "in_progress", "planned", or "shipped"
	Label   string // human-readable status label
	Items   []core.RoadmapItem
}

// RoadmapGenerator generates the roadmap page HTML
type RoadmapGenerator struct {
	site      *core.Site
	templates *template.Template
	themeFS   embed.FS
	version   string
}

// NewRoadmapGenerator creates a new roadmap generator
func NewRoadmapGenerator(site *core.Site, themeFS embed.FS, version string) (*RoadmapGenerator, error) {
	if !site.Config.Roadmap.Enabled {
		return nil, nil
	}

	tmpl := template.New("").Funcs(BaseFuncMap()).Funcs(AnalyticsFuncMap())

	var err error
	tmpl, err = tmpl.ParseFS(
		themeFS,
		"themes/common/templates/partials/landing-*.html",
		"themes/common/templates/partials/analytics.html",
		"themes/common/templates/partials/minimaldoc-widgets.html",
		"themes/common/templates/roadmap/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse roadmap templates: %w", err)
	}

	return &RoadmapGenerator{
		site:      site,
		templates: tmpl,
		themeFS:   themeFS,
		version:   version,
	}, nil
}

// Generate generates the roadmap page
func (g *RoadmapGenerator) Generate() error {
	if g.site.RoadmapPage == nil || !g.site.Config.Roadmap.Enabled {
		return nil
	}

	fmt.Println("Generating roadmap page...")

	roadmapPath := g.site.Config.Roadmap.Path
	if roadmapPath == "" {
		roadmapPath = "roadmap"
	}

	outputDir := filepath.Join(g.site.OutputRoot, roadmapPath)
	if err := makeWebDir(outputDir); err != nil {
		return fmt.Errorf("failed to create roadmap output directory: %w", err)
	}

	// Group items by version (time-based)
	groups := g.groupItemsByVersion()

	// Collect all unique tags
	tags := g.collectTags()

	data := map[string]any{
		"Site":       g.site,
		"Config":     g.site.RoadmapPage.Config,
		"BasePath":   g.getBasePath(),
		"Version":    g.version,
		"Groups":     groups,
		"Tags":       tags,
		"PageTitle":  g.site.Config.Roadmap.Title,
		"ActivePath": "/" + roadmapPath + "/",
		"Footer":     BuildFooter(g.site, g.version),
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "roadmap.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	outputPath := filepath.Join(outputDir, "index.html")
	if err := writeWebFile(outputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write roadmap page: %w", err)
	}

	fmt.Printf("Generated roadmap page with %d items in %d versions\n", len(g.site.Config.Roadmap.Items), len(groups))
	return nil
}

// groupItemsByVersion organizes items by version, ordered:
// in_progress groups first, then planned, then shipped.
// Within each tier, versions appear in order of first occurrence in config.
func (g *RoadmapGenerator) groupItemsByVersion() []RoadmapVersionGroup {
	cfg := g.site.Config.Roadmap

	// Build column label lookup from config
	columnLabels := make(map[string]string)
	for _, col := range cfg.Columns {
		columnLabels[col.ID] = col.Label
	}

	// Group items by version, preserving order of first appearance
	type versionEntry struct {
		version string
		items   []core.RoadmapItem
	}
	var ordered []versionEntry
	seen := make(map[string]int) // version -> index in ordered

	for _, item := range cfg.Items {
		v := item.Version
		if v == "" {
			v = "Unversioned"
		}
		if idx, ok := seen[v]; ok {
			ordered[idx].items = append(ordered[idx].items, item)
		} else {
			seen[v] = len(ordered)
			ordered = append(ordered, versionEntry{version: v, items: []core.RoadmapItem{item}})
		}
	}

	// Build groups with aggregate status
	var inProgress, planned, shipped []RoadmapVersionGroup

	for _, entry := range ordered {
		status := deriveGroupStatus(entry.items)
		label := columnLabels[status]
		if label == "" {
			label = status
		}

		group := RoadmapVersionGroup{
			Version: entry.version,
			Status:  status,
			Label:   label,
			Items:   entry.items,
		}

		// Sort into tiers: in_progress first, planned, shipped last
		switch status {
		case "in_progress":
			inProgress = append(inProgress, group)
		case "planned":
			planned = append(planned, group)
		default:
			shipped = append(shipped, group)
		}
	}

	// Combine: in_progress -> planned -> shipped
	result := make([]RoadmapVersionGroup, 0, len(ordered))
	result = append(result, inProgress...)
	result = append(result, planned...)
	result = append(result, shipped...)

	return result
}

// deriveGroupStatus returns the aggregate status for a version group.
// Priority: in_progress > planned > shipped.
func deriveGroupStatus(items []core.RoadmapItem) string {
	hasPlanned := false
	for _, item := range items {
		if item.Status == "in_progress" {
			return "in_progress"
		}
		if item.Status == "planned" {
			hasPlanned = true
		}
	}
	if hasPlanned {
		return "planned"
	}
	return "shipped"
}

// collectTags returns all unique tags across all items
func (g *RoadmapGenerator) collectTags() []string {
	seen := make(map[string]bool)
	var tags []string

	for _, item := range g.site.Config.Roadmap.Items {
		for _, tag := range item.Tags {
			if !seen[tag] {
				seen[tag] = true
				tags = append(tags, tag)
			}
		}
	}

	return tags
}

// getBasePath extracts the path component from BaseURL
func (g *RoadmapGenerator) getBasePath() string {
	return GetBasePath(g.site.Config.BaseURL)
}
