package core

import "time"

// StatusPage represents the complete status page data
type StatusPage struct {
	Config               StatusConfig
	Components           []StatusComponent
	ComponentGroups      map[string][]StatusComponent // grouped by Group field
	ActiveIncidents      []Incident
	ResolvedIncidents    []Incident
	ScheduledMaintenance []Maintenance
	ActiveMaintenance    []Maintenance
	OverallStatus        ComponentStatus
	LastUpdated          time.Time
	HistoryByMonth       map[string][]Incident // "2025-01" -> incidents
}

// StatusConfig holds configuration for the status page
type StatusConfig struct {
	Enabled       bool   `yaml:"enabled"`
	Title         string `yaml:"title"`
	Description   string `yaml:"description"`
	Path          string `yaml:"path"`           // output path (default: "status")
	ShowHistory   bool   `yaml:"show_history"`   // show incident history
	HistoryMonths int    `yaml:"history_months"` // months of history to show (default: 12)
	RSSEnabled    bool   `yaml:"rss_enabled"`    // generate RSS feed
}

// DefaultStatusConfig returns a StatusConfig with sensible defaults
func DefaultStatusConfig() StatusConfig {
	return StatusConfig{
		Enabled:       false,
		Title:         "Service Status",
		Description:   "Current status of our services",
		Path:          "status",
		ShowHistory:   true,
		HistoryMonths: 12,
		RSSEnabled:    true,
	}
}

// StatusComponent represents a service or component being monitored
type StatusComponent struct {
	ID          string          `yaml:"id"`
	Name        string          `yaml:"name"`
	Description string          `yaml:"description"`
	Status      ComponentStatus `yaml:"status"`
	Group       string          `yaml:"group"` // optional grouping
	Order       int             `yaml:"order"` // display order within group
	URL         string          `yaml:"url"`   // optional link to service

	// Health check configuration (for external monitoring tools)
	HealthEndpoint string `yaml:"health_endpoint"` // endpoint to check (e.g., /health)
	HealthInterval int    `yaml:"health_interval"` // check interval in seconds (0 = disabled)

	// Uptime tracking configuration
	Uptime UptimeConfig `yaml:"uptime"`

	// Computed uptime data (populated at build time for mode: incidents)
	UptimeData *UptimeData `yaml:"-" json:"uptime_data,omitempty"`
}

// UptimeConfig holds configuration for uptime tracking
type UptimeConfig struct {
	Mode       string  `yaml:"mode" json:"mode"`               // "api" or "incidents"
	Endpoint   string  `yaml:"endpoint" json:"endpoint"`       // API endpoint (for mode: api)
	SLATarget  float64 `yaml:"sla_target" json:"sla_target"`   // target SLA percentage (e.g., 99.9)
	PeriodDays int     `yaml:"period_days" json:"period_days"` // history period in days (default: 90)
}

// DefaultUptimeConfig returns an UptimeConfig with sensible defaults
func DefaultUptimeConfig() UptimeConfig {
	return UptimeConfig{
		Mode:       "",
		Endpoint:   "",
		SLATarget:  99.9,
		PeriodDays: 90,
	}
}

// UptimeDay represents a single day's uptime status
type UptimeDay struct {
	Date            string `json:"date"`             // "2025-01-28" format
	Status          string `json:"status"`           // operational, degraded, partial_outage, major_outage, maintenance
	DowntimeMinutes int    `json:"downtime_minutes"` // total downtime in minutes for this day
}

// SLAStats holds calculated SLA statistics
type SLAStats struct {
	Target     float64 `json:"target"`      // target SLA (e.g., 99.9)
	Current7d  float64 `json:"current_7d"`  // uptime % over last 7 days
	Current30d float64 `json:"current_30d"` // uptime % over last 30 days
	Current90d float64 `json:"current_90d"` // uptime % over last 90 days
}

// UptimeData holds complete uptime information for a component
type UptimeData struct {
	Mode          string      `json:"mode"`           // "api" or "incidents"
	PeriodDays    int         `json:"period_days"`    // number of days in history
	UptimePercent float64     `json:"uptime_percent"` // overall uptime percentage
	History       []UptimeDay `json:"history"`        // daily history (newest first)
	SLA           SLAStats    `json:"sla"`            // SLA statistics
}

// Incident represents a service incident or outage
type Incident struct {
	// Identity
	ID       string `yaml:"id"` // derived from filename if not set
	Slug     string `yaml:"-"`  // URL-friendly identifier
	FilePath string `yaml:"-"`  // source file path

	// Metadata from frontmatter
	Title              string           `yaml:"title"`
	Status             IncidentStatus   `yaml:"status"`
	Severity           IncidentSeverity `yaml:"severity"`
	AffectedComponents []string         `yaml:"affected_components"`
	CreatedAt          time.Time        `yaml:"created_at"`
	UpdatedAt          time.Time        `yaml:"updated_at"`
	ResolvedAt         *time.Time       `yaml:"resolved_at"`

	// Content
	Updates []IncidentUpdate `yaml:"-"` // parsed from markdown content
	RawMD   string           `yaml:"-"` // raw markdown content
	HTML    string           `yaml:"-"` // rendered HTML

	// Output
	OutputPath string `yaml:"-"` // generated HTML path
}

// IsActive returns true if the incident is not yet resolved
func (i *Incident) IsActive() bool {
	return i.Status.IsActive()
}

// Duration returns the duration of the incident (or time since creation if active)
func (i *Incident) Duration() time.Duration {
	if i.ResolvedAt != nil {
		return i.ResolvedAt.Sub(i.CreatedAt)
	}
	return time.Since(i.CreatedAt)
}

// IncidentUpdate represents a status update within an incident
type IncidentUpdate struct {
	Timestamp time.Time      `yaml:"timestamp"`
	Status    IncidentStatus `yaml:"status"`
	Message   string         `yaml:"message"` // rendered HTML
	RawMD     string         `yaml:"-"`       // raw markdown
}

// Maintenance represents scheduled maintenance
type Maintenance struct {
	// Identity
	ID       string `yaml:"id"`
	Slug     string `yaml:"-"`
	FilePath string `yaml:"-"`

	// Metadata
	Title              string            `yaml:"title"`
	Description        string            `yaml:"description"`
	AffectedComponents []string          `yaml:"affected_components"`
	ScheduledStart     time.Time         `yaml:"scheduled_start"`
	ScheduledEnd       time.Time         `yaml:"scheduled_end"`
	Status             MaintenanceStatus `yaml:"status"`

	// Content
	RawMD string `yaml:"-"`
	HTML  string `yaml:"-"`

	// Output
	OutputPath string `yaml:"-"`
}

// IsActive returns true if maintenance is currently in progress
func (m *Maintenance) IsActive() bool {
	return m.Status == MaintenanceInProgress
}

// IsUpcoming returns true if maintenance is scheduled but not started
func (m *Maintenance) IsUpcoming() bool {
	return m.Status == MaintenanceScheduled && time.Now().Before(m.ScheduledStart)
}

// Duration returns the planned duration of the maintenance
func (m *Maintenance) Duration() time.Duration {
	return m.ScheduledEnd.Sub(m.ScheduledStart)
}

// CalculateOverallStatus determines the worst status among all components
func CalculateOverallStatus(components []StatusComponent) ComponentStatus {
	if len(components) == 0 {
		return StatusOperational
	}

	worst := StatusOperational
	for _, c := range components {
		if c.Status.Severity() > worst.Severity() {
			worst = c.Status
		}
	}
	return worst
}

// GroupComponents organizes components by their Group field
func GroupComponents(components []StatusComponent) map[string][]StatusComponent {
	groups := make(map[string][]StatusComponent)
	for _, c := range components {
		group := c.Group
		if group == "" {
			group = "Services" // default group name
		}
		groups[group] = append(groups[group], c)
	}
	return groups
}

// FilterActiveIncidents returns only active (unresolved) incidents
func FilterActiveIncidents(incidents []Incident) []Incident {
	var active []Incident
	for _, inc := range incidents {
		if inc.IsActive() {
			active = append(active, inc)
		}
	}
	return active
}

// FilterResolvedIncidents returns only resolved incidents
func FilterResolvedIncidents(incidents []Incident) []Incident {
	var resolved []Incident
	for _, inc := range incidents {
		if !inc.IsActive() {
			resolved = append(resolved, inc)
		}
	}
	return resolved
}

// GroupIncidentsByMonth groups incidents by their creation month
func GroupIncidentsByMonth(incidents []Incident) map[string][]Incident {
	groups := make(map[string][]Incident)
	for _, inc := range incidents {
		key := inc.CreatedAt.Format("2006-01") // "YYYY-MM"
		groups[key] = append(groups[key], inc)
	}
	return groups
}

// FilterUpcomingMaintenance returns scheduled maintenance that hasn't started
func FilterUpcomingMaintenance(maintenance []Maintenance) []Maintenance {
	var upcoming []Maintenance
	now := time.Now()
	for _, m := range maintenance {
		if m.Status == MaintenanceScheduled && now.Before(m.ScheduledStart) {
			upcoming = append(upcoming, m)
		}
	}
	return upcoming
}

// FilterActiveMaintenance returns maintenance currently in progress
func FilterActiveMaintenance(maintenance []Maintenance) []Maintenance {
	var active []Maintenance
	for _, m := range maintenance {
		if m.Status == MaintenanceInProgress {
			active = append(active, m)
		}
	}
	return active
}
