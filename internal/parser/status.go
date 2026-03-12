package parser

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/studiowebux/minimaldoc/internal/core"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	htmlrenderer "github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

// StatusParser handles parsing of status page content
type StatusParser struct {
	md goldmark.Markdown
}

// NewStatusParser creates a new status parser
func NewStatusParser() *StatusParser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			htmlrenderer.WithUnsafe(),
		),
	)

	return &StatusParser{md: md}
}

// ParseStatusDir parses all status content from a directory
func (p *StatusParser) ParseStatusDir(statusDir string) (*core.StatusPage, error) {
	statusPage := &core.StatusPage{
		Config:      core.DefaultStatusConfig(),
		LastUpdated: time.Now(),
	}

	// Check if status directory exists
	if _, err := os.Stat(statusDir); os.IsNotExist(err) {
		return statusPage, nil // Return empty status page if dir doesn't exist
	}

	// Note: Config is set from main site config.yaml, not from __status__/config.yaml

	// Parse components.yaml
	componentsPath := filepath.Join(statusDir, "components.yaml")
	if _, err := os.Stat(componentsPath); err == nil {
		components, err := p.ParseComponents(componentsPath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse components: %w", err)
		}
		statusPage.Components = components
	}

	// Parse incidents
	incidentsDir := filepath.Join(statusDir, "incidents")
	if _, err := os.Stat(incidentsDir); err == nil {
		incidents, err := p.ParseIncidents(incidentsDir)
		if err != nil {
			return nil, fmt.Errorf("failed to parse incidents: %w", err)
		}

		// Sort by creation date (newest first)
		sort.Slice(incidents, func(i, j int) bool {
			return incidents[i].CreatedAt.After(incidents[j].CreatedAt)
		})

		statusPage.ActiveIncidents = core.FilterActiveIncidents(incidents)
		statusPage.ResolvedIncidents = core.FilterResolvedIncidents(incidents)
		statusPage.HistoryByMonth = core.GroupIncidentsByMonth(incidents)
	}

	// Parse maintenance
	maintenanceDir := filepath.Join(statusDir, "maintenance")
	if _, err := os.Stat(maintenanceDir); err == nil {
		maintenance, err := p.ParseMaintenance(maintenanceDir)
		if err != nil {
			return nil, fmt.Errorf("failed to parse maintenance: %w", err)
		}

		// Sort by scheduled start (soonest first)
		sort.Slice(maintenance, func(i, j int) bool {
			return maintenance[i].ScheduledStart.Before(maintenance[j].ScheduledStart)
		})

		statusPage.ScheduledMaintenance = core.FilterUpcomingMaintenance(maintenance)
		statusPage.ActiveMaintenance = core.FilterActiveMaintenance(maintenance)
	}

	// Calculate overall status
	statusPage.OverallStatus = core.CalculateOverallStatus(statusPage.Components)
	statusPage.ComponentGroups = core.GroupComponents(statusPage.Components)

	return statusPage, nil
}

// ParseConfig parses the status config.yaml file
func (p *StatusParser) ParseConfig(path string) (core.StatusConfig, error) {
	config := core.DefaultStatusConfig()

	data, err := os.ReadFile(path) // #nosec G304 -- path from trusted user configuration
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	return config, nil
}

// ParseComponents parses the components.yaml file
func (p *StatusParser) ParseComponents(path string) ([]core.StatusComponent, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- path from trusted user configuration
	if err != nil {
		return nil, err
	}

	var components []core.StatusComponent
	if err := yaml.Unmarshal(data, &components); err != nil {
		return nil, fmt.Errorf("failed to parse components YAML: %w", err)
	}

	// Set defaults
	for i := range components {
		if components[i].Status == "" {
			components[i].Status = core.StatusOperational
		}
		// Set uptime defaults
		if components[i].Uptime.PeriodDays == 0 {
			components[i].Uptime.PeriodDays = 90
		}
		if components[i].Uptime.SLATarget == 0 {
			components[i].Uptime.SLATarget = 99.9
		}
	}

	return components, nil
}

// ParseIncidents parses all incident markdown files in a directory
func (p *StatusParser) ParseIncidents(dir string) ([]core.Incident, error) {
	var incidents []core.Incident

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		incident, err := p.ParseIncident(path)
		if err != nil {
			return nil, fmt.Errorf("failed to parse incident %s: %w", path, err)
		}

		incidents = append(incidents, incident)
	}

	return incidents, nil
}

// IncidentFrontmatter represents the YAML frontmatter for incidents
type IncidentFrontmatter struct {
	ID                 string                `yaml:"id"`
	Title              string                `yaml:"title"`
	Status             core.IncidentStatus   `yaml:"status"`
	Severity           core.IncidentSeverity `yaml:"severity"`
	AffectedComponents []string              `yaml:"affected_components"`
	CreatedAt          time.Time             `yaml:"created_at"`
	UpdatedAt          time.Time             `yaml:"updated_at"`
	ResolvedAt         *time.Time            `yaml:"resolved_at"`
}

// ParseIncident parses a single incident markdown file
func (p *StatusParser) ParseIncident(path string) (core.Incident, error) {
	var incident core.Incident
	incident.FilePath = path

	// Generate slug from filename
	filename := filepath.Base(path)
	incident.Slug = strings.TrimSuffix(filename, ".md")

	// Read file
	content, err := os.ReadFile(path) // #nosec G304 -- path from trusted user configuration
	if err != nil {
		return incident, err
	}

	// Parse frontmatter
	var fm IncidentFrontmatter
	rest, err := frontmatter.Parse(bytes.NewReader(content), &fm)
	if err != nil {
		return incident, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Copy frontmatter to incident
	incident.ID = fm.ID
	if incident.ID == "" {
		incident.ID = incident.Slug
	}
	incident.Title = fm.Title
	incident.Status = fm.Status
	incident.Severity = fm.Severity
	incident.AffectedComponents = fm.AffectedComponents
	incident.CreatedAt = fm.CreatedAt
	incident.UpdatedAt = fm.UpdatedAt
	incident.ResolvedAt = fm.ResolvedAt

	// Set defaults
	if incident.Status == "" {
		incident.Status = core.IncidentInvestigating
	}
	if incident.Severity == "" {
		incident.Severity = core.SeverityMinor
	}
	if incident.UpdatedAt.IsZero() {
		incident.UpdatedAt = incident.CreatedAt
	}

	// Store raw markdown
	incident.RawMD = string(rest)

	// Parse updates from markdown content
	incident.Updates = p.parseIncidentUpdates(rest, incident.CreatedAt)

	// Render full content to HTML
	var buf bytes.Buffer
	if err := p.md.Convert(rest, &buf); err != nil {
		return incident, fmt.Errorf("failed to render markdown: %w", err)
	}
	incident.HTML = buf.String()

	return incident, nil
}

// updateHeaderRegex matches h2 headers with timestamps like "## Update - 12:30 UTC" or "## 2025-01-28 12:30 UTC"
var updateHeaderRegex = regexp.MustCompile(`(?m)^##\s+(?:Update\s*-?\s*)?(.+)$`)

// parseIncidentUpdates extracts updates from markdown content based on h2 headers
func (p *StatusParser) parseIncidentUpdates(content []byte, baseDate time.Time) []core.IncidentUpdate {
	var updates []core.IncidentUpdate

	lines := strings.Split(string(content), "\n")
	var currentUpdate *core.IncidentUpdate
	var currentContent []string

	for _, line := range lines {
		if matches := updateHeaderRegex.FindStringSubmatch(line); matches != nil {
			// Save previous update if exists
			if currentUpdate != nil {
				currentUpdate.RawMD = strings.TrimSpace(strings.Join(currentContent, "\n"))
				var buf bytes.Buffer
				_ = p.md.Convert([]byte(currentUpdate.RawMD), &buf) // goldmark Convert rarely errors on valid markdown
				currentUpdate.Message = buf.String()
				updates = append(updates, *currentUpdate)
			}

			// Start new update
			timestamp := p.parseTimestamp(matches[1], baseDate)
			currentUpdate = &core.IncidentUpdate{
				Timestamp: timestamp,
				Status:    p.detectStatusFromContent(matches[1]),
			}
			currentContent = nil
		} else if currentUpdate != nil {
			currentContent = append(currentContent, line)
		}
	}

	// Save last update
	if currentUpdate != nil {
		currentUpdate.RawMD = strings.TrimSpace(strings.Join(currentContent, "\n"))
		var buf bytes.Buffer
		_ = p.md.Convert([]byte(currentUpdate.RawMD), &buf) // goldmark Convert rarely errors on valid markdown
		currentUpdate.Message = buf.String()
		updates = append(updates, *currentUpdate)
	}

	return updates
}

// parseTimestamp attempts to parse various timestamp formats from update headers
// baseDate is used when only time is provided (e.g., "12:30 UTC")
func (p *StatusParser) parseTimestamp(header string, baseDate time.Time) time.Time {
	header = strings.TrimSpace(header)

	// Formats with full date
	fullDateFormats := []string{
		"2006-01-02 15:04 MST",
		"2006-01-02 15:04",
		"January 2, 2006 15:04 MST",
		"Jan 2, 2006 15:04",
		time.RFC3339,
	}

	for _, format := range fullDateFormats {
		if t, err := time.Parse(format, header); err == nil {
			return t
		}
	}

	// Time-only formats - use baseDate for the date portion
	timeOnlyFormats := []string{
		"15:04 MST",
		"15:04",
	}

	for _, format := range timeOnlyFormats {
		if t, err := time.Parse(format, header); err == nil {
			// Combine base date with parsed time
			return time.Date(
				baseDate.Year(), baseDate.Month(), baseDate.Day(),
				t.Hour(), t.Minute(), t.Second(), 0,
				baseDate.Location(),
			)
		}
	}

	// If no format matches, return base date
	return baseDate
}

// detectStatusFromContent tries to detect incident status from update header/content
func (p *StatusParser) detectStatusFromContent(content string) core.IncidentStatus {
	lower := strings.ToLower(content)

	switch {
	case strings.Contains(lower, "resolved"):
		return core.IncidentResolved
	case strings.Contains(lower, "monitoring"):
		return core.IncidentMonitoring
	case strings.Contains(lower, "identified"):
		return core.IncidentIdentified
	default:
		return core.IncidentInvestigating
	}
}

// ParseMaintenance parses all maintenance markdown files in a directory
func (p *StatusParser) ParseMaintenance(dir string) ([]core.Maintenance, error) {
	var maintenance []core.Maintenance

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		m, err := p.ParseMaintenanceFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to parse maintenance %s: %w", path, err)
		}

		maintenance = append(maintenance, m)
	}

	return maintenance, nil
}

// MaintenanceFrontmatter represents the YAML frontmatter for maintenance
type MaintenanceFrontmatter struct {
	ID                 string                 `yaml:"id"`
	Title              string                 `yaml:"title"`
	Description        string                 `yaml:"description"`
	AffectedComponents []string               `yaml:"affected_components"`
	ScheduledStart     time.Time              `yaml:"scheduled_start"`
	ScheduledEnd       time.Time              `yaml:"scheduled_end"`
	Status             core.MaintenanceStatus `yaml:"status"`
}

// ParseMaintenanceFile parses a single maintenance markdown file
func (p *StatusParser) ParseMaintenanceFile(path string) (core.Maintenance, error) {
	var m core.Maintenance
	m.FilePath = path

	// Generate slug from filename
	filename := filepath.Base(path)
	m.Slug = strings.TrimSuffix(filename, ".md")

	// Read file
	content, err := os.ReadFile(path) // #nosec G304 -- path from trusted user configuration
	if err != nil {
		return m, err
	}

	// Parse frontmatter
	var fm MaintenanceFrontmatter
	rest, err := frontmatter.Parse(bytes.NewReader(content), &fm)
	if err != nil {
		return m, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Copy frontmatter to maintenance
	m.ID = fm.ID
	if m.ID == "" {
		m.ID = m.Slug
	}
	m.Title = fm.Title
	m.Description = fm.Description
	m.AffectedComponents = fm.AffectedComponents
	m.ScheduledStart = fm.ScheduledStart
	m.ScheduledEnd = fm.ScheduledEnd
	m.Status = fm.Status

	// Set defaults
	if m.Status == "" {
		m.Status = core.MaintenanceScheduled
	}

	// Store raw markdown
	m.RawMD = string(rest)

	// Render content to HTML
	var buf bytes.Buffer
	if err := p.md.Convert(rest, &buf); err != nil {
		return m, fmt.Errorf("failed to render markdown: %w", err)
	}
	m.HTML = buf.String()

	return m, nil
}
