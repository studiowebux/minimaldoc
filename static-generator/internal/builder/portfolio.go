package builder

import (
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/studiowebux/minimaldoc/static-generator/internal/core"
	"github.com/studiowebux/minimaldoc/static-generator/internal/parser"
)

// PortfolioBuilder handles building the portfolio page
type PortfolioBuilder struct {
	frontmatterParser *parser.FrontmatterParser
	markdownParser    *parser.MarkdownParser
}

// NewPortfolioBuilder creates a new portfolio builder
func NewPortfolioBuilder() *PortfolioBuilder {
	return &PortfolioBuilder{
		frontmatterParser: parser.NewFrontmatterParser(),
		markdownParser:    parser.NewMarkdownParser(),
	}
}

// Build creates the portfolio page from markdown files
func (pb *PortfolioBuilder) Build(docsRoot string, config core.PortfolioConfig, basePath string) (*core.PortfolioPage, error) {
	if !config.Enabled {
		return nil, nil
	}

	portfolioDir := filepath.Join(docsRoot, core.PortfolioSourceDir)

	// Parse all portfolio projects
	projects, err := pb.parseProjects(portfolioDir, basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse portfolio projects: %w", err)
	}

	// Sort projects by date (newest first) then by order
	sort.Slice(projects, func(i, j int) bool {
		if projects[i].Order != projects[j].Order {
			return projects[i].Order < projects[j].Order
		}
		return projects[j].Date.Before(projects[i].Date)
	})

	page := &core.PortfolioPage{
		Config:           config,
		Projects:         projects,
		FeaturedProjects: core.FilterFeaturedProjects(projects),
		Tags:             core.CollectProjectTags(projects),
	}

	return page, nil
}

// parseProjects parses all markdown files in the portfolio directory
func (pb *PortfolioBuilder) parseProjects(portfolioDir string, basePath string) ([]core.Project, error) {
	var projects []core.Project

	// Glob for markdown files
	pattern := filepath.Join(portfolioDir, "*.md")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob portfolio files: %w", err)
	}

	for _, filePath := range files {
		project, err := pb.parseProject(filePath, basePath)
		if err != nil {
			fmt.Printf("Warning: failed to parse portfolio project %s: %v\n", filePath, err)
			continue
		}
		projects = append(projects, *project)
	}

	return projects, nil
}

// parseProject parses a single project markdown file
func (pb *PortfolioBuilder) parseProject(filePath string, basePath string) (*core.Project, error) {
	// Parse frontmatter
	meta, content, err := pb.frontmatterParser.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("frontmatter parse error: %w", err)
	}

	// Generate slug from filename
	base := filepath.Base(filePath)
	slug := base[:len(base)-len(filepath.Ext(base))]

	// Parse markdown to HTML
	html, err := pb.markdownParser.ParseWithContext(content, "", basePath)
	if err != nil {
		return nil, fmt.Errorf("markdown parse error: %w", err)
	}

	// Parse date
	var projectDate time.Time
	if meta.Date != "" {
		parsedDate, parseErr := time.Parse("2006-01-02", meta.Date)
		if parseErr == nil {
			projectDate = parsedDate
		}
	}

	project := &core.Project{
		ID:          slug,
		Slug:        slug,
		FilePath:    filePath,
		Title:       meta.Title,
		Description: meta.Description,
		Image:       meta.Image,
		Tags:        meta.Tags,
		Links:       convertLinksToSimpleLinks(meta.Links),
		Date:        projectDate,
		Featured:    meta.Featured,
		Order:       meta.MenuOrder,
		RawMD:       string(content),
		HTML:        string(html),
	}

	return project, nil
}

// convertLinksToSimpleLinks converts metadata links to SimpleLink format
func convertLinksToSimpleLinks(links []core.MetadataLink) []core.SimpleLink {
	result := make([]core.SimpleLink, len(links))
	for i, l := range links {
		result[i] = core.SimpleLink(l)
	}
	return result
}
