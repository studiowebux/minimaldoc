package generator

import (
	"bytes"
	"embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// StatusGenerator generates status page HTML and related files
type StatusGenerator struct {
	site      *core.Site
	templates *template.Template
	themeFS   embed.FS
	version   string
}

// NewStatusGenerator creates a new status generator
func NewStatusGenerator(site *core.Site, themeFS embed.FS, version string) (*StatusGenerator, error) {
	if !site.Config.Status.Enabled {
		return nil, nil // Skip if status is not enabled
	}

	// Create template with shared status functions
	tmpl := template.New("").Funcs(StatusFuncMap()).Funcs(AnalyticsFuncMap())

	// Parse status templates from dedicated subdirectory
	var err error
	tmpl, err = tmpl.ParseFS(
		themeFS,
		"themes/common/templates/partials/landing-*.html",
		"themes/common/templates/partials/analytics.html",
		"themes/common/templates/partials/minimaldoc-widgets.html",
		"themes/common/templates/status/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse status templates: %w", err)
	}

	return &StatusGenerator{
		site:      site,
		templates: tmpl,
		themeFS:   themeFS,
		version:   version,
	}, nil
}

// Generate generates all status page files
func (g *StatusGenerator) Generate() error {
	if g.site.StatusPage == nil || !g.site.Config.Status.Enabled {
		return nil
	}

	fmt.Println("Generating status page...")

	statusPath := g.site.Config.Status.Path
	if statusPath == "" {
		statusPath = "status"
	}

	outputDir := filepath.Join(g.site.OutputRoot, statusPath)
	if err := makeWebDir(outputDir); err != nil {
		return fmt.Errorf("failed to create status output directory: %w", err)
	}

	// Generate main status page
	if err := g.generateMainPage(outputDir); err != nil {
		return fmt.Errorf("failed to generate main status page: %w", err)
	}

	// Generate incident pages
	if err := g.generateIncidentPages(outputDir); err != nil {
		return fmt.Errorf("failed to generate incident pages: %w", err)
	}

	// Generate maintenance pages
	if err := g.generateMaintenancePages(outputDir); err != nil {
		return fmt.Errorf("failed to generate maintenance pages: %w", err)
	}

	// Generate history page
	if g.site.Config.Status.ShowHistory {
		if err := g.generateHistoryPage(outputDir); err != nil {
			return fmt.Errorf("failed to generate history page: %w", err)
		}
	}

	// Generate status.json
	if err := g.generateStatusJSON(outputDir); err != nil {
		return fmt.Errorf("failed to generate status.json: %w", err)
	}

	// Generate RSS feed
	if g.site.Config.Status.RSSEnabled {
		if err := g.generateRSSFeed(outputDir); err != nil {
			return fmt.Errorf("failed to generate RSS feed: %w", err)
		}
	}

	fmt.Printf("Generated status page with %d components\n", len(g.site.StatusPage.Components))
	return nil
}

// generateMainPage generates the main status dashboard
func (g *StatusGenerator) generateMainPage(outputDir string) error {
	statusPath := g.site.Config.Status.Path
	if statusPath == "" {
		statusPath = "status"
	}
	data := map[string]any{
		"Site":       g.site,
		"StatusPage": g.site.StatusPage,
		"BasePath":   g.getBasePath(),
		"Version":    g.version,
		"PageTitle":  g.site.StatusPage.Config.Title,
		"ActivePath": "/" + statusPath + "/",
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "status.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	outputPath := filepath.Join(outputDir, "index.html")
	if err := writeWebFile(outputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write status page: %w", err)
	}

	return nil
}

// generateIncidentPages generates individual incident detail pages
func (g *StatusGenerator) generateIncidentPages(outputDir string) error {
	incidentDir := filepath.Join(outputDir, "incident")
	if err := makeWebDir(incidentDir); err != nil {
		return fmt.Errorf("failed to create incident directory: %w", err)
	}

	// Generate pages for all incidents (active and resolved) — safe copy to avoid slice mutation
	allIncidents := make([]core.Incident, 0, len(g.site.StatusPage.ActiveIncidents)+len(g.site.StatusPage.ResolvedIncidents))
	allIncidents = append(allIncidents, g.site.StatusPage.ActiveIncidents...)
	allIncidents = append(allIncidents, g.site.StatusPage.ResolvedIncidents...)

	for _, incident := range allIncidents {
		if err := g.generateIncidentPage(incidentDir, incident); err != nil {
			return fmt.Errorf("failed to generate incident %s: %w", incident.Slug, err)
		}
	}

	return nil
}

// generateIncidentPage generates a single incident detail page
func (g *StatusGenerator) generateIncidentPage(incidentDir string, incident core.Incident) error {
	statusPath := g.site.Config.Status.Path
	if statusPath == "" {
		statusPath = "status"
	}
	data := map[string]any{
		"Site":       g.site,
		"StatusPage": g.site.StatusPage,
		"Incident":   incident,
		"BasePath":   g.getBasePath(),
		"Version":    g.version,
		"PageTitle":  incident.Title + " - " + g.site.StatusPage.Config.Title,
		"ActivePath": "/" + statusPath + "/",
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "status-incident.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	outputPath := filepath.Join(incidentDir, incident.Slug+".html")
	if err := writeWebFile(outputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write incident page: %w", err)
	}

	return nil
}

// generateMaintenancePages generates individual maintenance detail pages
func (g *StatusGenerator) generateMaintenancePages(outputDir string) error {
	maintenanceDir := filepath.Join(outputDir, "maintenance")
	if err := makeWebDir(maintenanceDir); err != nil {
		return fmt.Errorf("failed to create maintenance directory: %w", err)
	}

	// Generate pages for all maintenance (scheduled and active) — safe copy
	allMaintenance := make([]core.Maintenance, 0, len(g.site.StatusPage.ScheduledMaintenance)+len(g.site.StatusPage.ActiveMaintenance))
	allMaintenance = append(allMaintenance, g.site.StatusPage.ScheduledMaintenance...)
	allMaintenance = append(allMaintenance, g.site.StatusPage.ActiveMaintenance...)

	for _, maintenance := range allMaintenance {
		if err := g.generateMaintenancePage(maintenanceDir, maintenance); err != nil {
			return fmt.Errorf("failed to generate maintenance %s: %w", maintenance.Slug, err)
		}
	}

	return nil
}

// generateMaintenancePage generates a single maintenance detail page
func (g *StatusGenerator) generateMaintenancePage(maintenanceDir string, maintenance core.Maintenance) error {
	stPath := g.site.Config.Status.Path
	if stPath == "" {
		stPath = "status"
	}
	data := map[string]any{
		"Site":        g.site,
		"StatusPage":  g.site.StatusPage,
		"Maintenance": maintenance,
		"BasePath":    g.getBasePath(),
		"Version":     g.version,
		"PageTitle":   maintenance.Title + " - " + g.site.StatusPage.Config.Title,
		"ActivePath":  "/" + stPath + "/",
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "status-maintenance.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	outputPath := filepath.Join(maintenanceDir, maintenance.Slug+".html")
	if err := writeWebFile(outputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write maintenance page: %w", err)
	}

	return nil
}

// generateHistoryPage generates the incident history page
func (g *StatusGenerator) generateHistoryPage(outputDir string) error {
	histPath := g.site.Config.Status.Path
	if histPath == "" {
		histPath = "status"
	}
	data := map[string]any{
		"Site":           g.site,
		"StatusPage":     g.site.StatusPage,
		"HistoryByMonth": g.site.StatusPage.HistoryByMonth,
		"BasePath":       g.getBasePath(),
		"Version":        g.version,
		"PageTitle":      "Incident History - " + g.site.StatusPage.Config.Title,
		"ActivePath":     "/" + histPath + "/",
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "status-history.html", data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	outputPath := filepath.Join(outputDir, "history.html")
	if err := writeWebFile(outputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write history page: %w", err)
	}

	return nil
}

// StatusJSON represents the JSON API response for status
type StatusJSON struct {
	OverallStatus        string            `json:"overall_status"`
	Components           []ComponentJSON   `json:"components"`
	ActiveIncidents      []IncidentJSON    `json:"active_incidents"`
	ScheduledMaintenance []MaintenanceJSON `json:"scheduled_maintenance"`
	LastUpdated          time.Time         `json:"last_updated"`
}

// ComponentJSON represents a component in the JSON API
type ComponentJSON struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description,omitempty"`
	Status         string            `json:"status"`
	Group          string            `json:"group,omitempty"`
	URL            string            `json:"url,omitempty"`
	HealthEndpoint string            `json:"health_endpoint,omitempty"`
	HealthInterval int               `json:"health_interval,omitempty"`
	Uptime         *UptimeConfigJSON `json:"uptime,omitempty"`
	UptimeData     *core.UptimeData  `json:"uptime_data,omitempty"`
}

// UptimeConfigJSON represents uptime configuration in the JSON API
type UptimeConfigJSON struct {
	Mode       string  `json:"mode,omitempty"`
	Endpoint   string  `json:"endpoint,omitempty"`
	SLATarget  float64 `json:"sla_target,omitempty"`
	PeriodDays int     `json:"period_days,omitempty"`
}

// IncidentJSON represents an incident in the JSON API
type IncidentJSON struct {
	ID                 string    `json:"id"`
	Title              string    `json:"title"`
	Status             string    `json:"status"`
	Severity           string    `json:"severity"`
	AffectedComponents []string  `json:"affected_components"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	URL                string    `json:"url"`
}

// MaintenanceJSON represents maintenance in the JSON API
type MaintenanceJSON struct {
	ID                 string    `json:"id"`
	Title              string    `json:"title"`
	Description        string    `json:"description,omitempty"`
	AffectedComponents []string  `json:"affected_components"`
	ScheduledStart     time.Time `json:"scheduled_start"`
	ScheduledEnd       time.Time `json:"scheduled_end"`
	Status             string    `json:"status"`
}

// generateStatusJSON generates the status.json API file
func (g *StatusGenerator) generateStatusJSON(outputDir string) error {
	statusPath := g.site.Config.Status.Path
	if statusPath == "" {
		statusPath = "status"
	}

	// Build components list
	components := make([]ComponentJSON, 0, len(g.site.StatusPage.Components))
	for _, c := range g.site.StatusPage.Components {
		comp := ComponentJSON{
			ID:             c.ID,
			Name:           c.Name,
			Description:    c.Description,
			Status:         string(c.Status),
			Group:          c.Group,
			URL:            c.URL,
			HealthEndpoint: c.HealthEndpoint,
			HealthInterval: c.HealthInterval,
		}

		// Include uptime config if mode is set
		if c.Uptime.Mode != "" {
			comp.Uptime = &UptimeConfigJSON{
				Mode:       c.Uptime.Mode,
				Endpoint:   c.Uptime.Endpoint,
				SLATarget:  c.Uptime.SLATarget,
				PeriodDays: c.Uptime.PeriodDays,
			}
		}

		// Include computed uptime data if available
		if c.UptimeData != nil {
			comp.UptimeData = c.UptimeData
		}

		components = append(components, comp)
	}

	// Build active incidents list
	incidents := make([]IncidentJSON, 0, len(g.site.StatusPage.ActiveIncidents))
	for _, inc := range g.site.StatusPage.ActiveIncidents {
		incidentURL := g.getBasePath() + "/" + statusPath + "/incident/" + inc.Slug + ".html"
		incidents = append(incidents, IncidentJSON{
			ID:                 inc.ID,
			Title:              inc.Title,
			Status:             string(inc.Status),
			Severity:           string(inc.Severity),
			AffectedComponents: inc.AffectedComponents,
			CreatedAt:          inc.CreatedAt,
			UpdatedAt:          inc.UpdatedAt,
			URL:                incidentURL,
		})
	}

	// Build scheduled maintenance list
	maintenance := make([]MaintenanceJSON, 0, len(g.site.StatusPage.ScheduledMaintenance))
	for _, m := range g.site.StatusPage.ScheduledMaintenance {
		maintenance = append(maintenance, MaintenanceJSON{
			ID:                 m.ID,
			Title:              m.Title,
			Description:        m.Description,
			AffectedComponents: m.AffectedComponents,
			ScheduledStart:     m.ScheduledStart,
			ScheduledEnd:       m.ScheduledEnd,
			Status:             string(m.Status),
		})
	}

	statusJSON := StatusJSON{
		OverallStatus:        string(g.site.StatusPage.OverallStatus),
		Components:           components,
		ActiveIncidents:      incidents,
		ScheduledMaintenance: maintenance,
		LastUpdated:          g.site.StatusPage.LastUpdated,
	}

	data, err := json.MarshalIndent(statusJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal status JSON: %w", err)
	}

	outputPath := filepath.Join(outputDir, "status.json")
	if err := writeWebFile(outputPath, data); err != nil {
		return fmt.Errorf("failed to write status.json: %w", err)
	}

	return nil
}

// RSS feed structures
type rssChannel struct {
	XMLName       xml.Name  `xml:"channel"`
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

// generateRSSFeed generates an RSS feed for status updates
func (g *StatusGenerator) generateRSSFeed(outputDir string) error {
	statusPath := g.site.Config.Status.Path
	if statusPath == "" {
		statusPath = "status"
	}

	baseURL := strings.TrimSuffix(g.site.Config.BaseURL, "/")

	// Build RSS items from incidents
	items := make([]rssItem, 0)

	// Add active incidents
	for _, inc := range g.site.StatusPage.ActiveIncidents {
		incidentURL := baseURL + "/" + statusPath + "/incident/" + inc.Slug + ".html"
		items = append(items, rssItem{
			Title:       fmt.Sprintf("[%s] %s", inc.Severity, inc.Title),
			Link:        incidentURL,
			Description: fmt.Sprintf("Status: %s. Affected: %v", inc.Status, inc.AffectedComponents),
			PubDate:     inc.CreatedAt.Format(time.RFC1123Z),
			GUID:        incidentURL,
		})
	}

	// Add recent resolved incidents (last 10)
	resolved := g.site.StatusPage.ResolvedIncidents
	if len(resolved) > 10 {
		resolved = resolved[:10]
	}
	for _, inc := range resolved {
		incidentURL := baseURL + "/" + statusPath + "/incident/" + inc.Slug + ".html"
		items = append(items, rssItem{
			Title:       fmt.Sprintf("[Resolved] %s", inc.Title),
			Link:        incidentURL,
			Description: fmt.Sprintf("Resolved on %s", inc.ResolvedAt.Format(time.RFC1123)),
			PubDate:     inc.UpdatedAt.Format(time.RFC1123Z),
			GUID:        incidentURL,
		})
	}

	// Add scheduled maintenance
	for _, m := range g.site.StatusPage.ScheduledMaintenance {
		items = append(items, rssItem{
			Title:       fmt.Sprintf("[Maintenance] %s", m.Title),
			Link:        baseURL + "/" + statusPath,
			Description: fmt.Sprintf("Scheduled: %s - %s", m.ScheduledStart.Format(time.RFC1123), m.ScheduledEnd.Format(time.RFC1123)),
			PubDate:     time.Now().Format(time.RFC1123Z),
			GUID:        fmt.Sprintf("%s/maintenance/%s", baseURL, m.ID),
		})
	}

	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:         g.site.StatusPage.Config.Title,
			Link:          baseURL + "/" + statusPath,
			Description:   g.site.StatusPage.Config.Description,
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         items,
		},
	}

	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal RSS feed: %w", err)
	}

	// Add XML declaration
	xmlData := []byte(xml.Header + string(data))

	outputPath := filepath.Join(outputDir, "feed.xml")
	if err := writeWebFile(outputPath, xmlData); err != nil {
		return fmt.Errorf("failed to write RSS feed: %w", err)
	}

	return nil
}

// getBasePath extracts the path component from BaseURL for asset linking
func (g *StatusGenerator) getBasePath() string {
	return GetBasePath(g.site.Config.BaseURL)
}
