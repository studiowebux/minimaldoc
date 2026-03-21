package builder

import (
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/studiowebux/minimaldoc/internal/core"
	"github.com/studiowebux/minimaldoc/internal/parser"
)

const minutesPerDay = 24 * 60

// StatusBuilder handles building the status page
type StatusBuilder struct {
	parser *parser.StatusParser
}

// NewStatusBuilder creates a new status builder
func NewStatusBuilder() *StatusBuilder {
	return &StatusBuilder{
		parser: parser.NewStatusParser(),
	}
}

// Build parses and builds the status page data
func (sb *StatusBuilder) Build(docsRoot string, config core.StatusConfig) (*core.StatusPage, error) {
	statusDir := filepath.Join(docsRoot, core.StatusSourceDir)

	// Parse all status content
	statusPage, err := sb.parser.ParseStatusDir(statusDir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse status directory: %w", err)
	}

	// Merge config from site config (file config overrides defaults, site config takes precedence)
	statusPage.Config = sb.mergeConfig(statusPage.Config, config)

	// Set output paths for incidents
	for i := range statusPage.ActiveIncidents {
		statusPage.ActiveIncidents[i].OutputPath = sb.getIncidentOutputPath(
			statusPage.Config.Path,
			statusPage.ActiveIncidents[i].Slug,
			config.Path != "" && config.Path != "status",
		)
	}
	for i := range statusPage.ResolvedIncidents {
		statusPage.ResolvedIncidents[i].OutputPath = sb.getIncidentOutputPath(
			statusPage.Config.Path,
			statusPage.ResolvedIncidents[i].Slug,
			config.Path != "" && config.Path != "status",
		)
	}

	// Sort components by group and order
	sb.sortComponents(statusPage)

	// Calculate uptime data for components with mode: incidents
	allIncidents := make([]core.Incident, 0, len(statusPage.ActiveIncidents)+len(statusPage.ResolvedIncidents))
	allIncidents = append(allIncidents, statusPage.ActiveIncidents...)
	allIncidents = append(allIncidents, statusPage.ResolvedIncidents...)
	for i := range statusPage.Components {
		if statusPage.Components[i].Uptime.Mode == "incidents" {
			uptimeData := sb.CalculateUptimeFromIncidents(
				statusPage.Components[i],
				allIncidents,
				statusPage.Components[i].Uptime.PeriodDays,
			)
			statusPage.Components[i].UptimeData = uptimeData
		}
	}

	// Update component groups with uptime data
	statusPage.ComponentGroups = core.GroupComponents(statusPage.Components)

	// Update last updated time
	statusPage.LastUpdated = time.Now()

	return statusPage, nil
}

// mergeConfig merges two status configs, with override taking precedence for non-zero values
func (sb *StatusBuilder) mergeConfig(base, override core.StatusConfig) core.StatusConfig {
	result := base

	// Override enabled state
	if override.Enabled {
		result.Enabled = override.Enabled
	}

	// Override title if set
	if override.Title != "" && override.Title != "Service Status" {
		result.Title = override.Title
	}

	// Override description if set
	if override.Description != "" && override.Description != "Current status of our services" {
		result.Description = override.Description
	}

	// Override path if set
	if override.Path != "" && override.Path != "status" {
		result.Path = override.Path
	}

	// Override history settings
	if override.HistoryMonths > 0 && override.HistoryMonths != 12 {
		result.HistoryMonths = override.HistoryMonths
	}

	// RSS enabled can be explicitly disabled
	result.RSSEnabled = override.RSSEnabled

	return result
}

// getIncidentOutputPath generates the output path for an incident
func (sb *StatusBuilder) getIncidentOutputPath(statusPath, slug string, cleanURLs bool) string {
	if cleanURLs {
		return filepath.Join(statusPath, "incident", slug, "index.html")
	}
	return filepath.Join(statusPath, "incident", slug+".html")
}

// sortComponents sorts components by group and order within groups
func (sb *StatusBuilder) sortComponents(statusPage *core.StatusPage) {
	// Sort main component list by group then order
	sort.Slice(statusPage.Components, func(i, j int) bool {
		if statusPage.Components[i].Group != statusPage.Components[j].Group {
			return statusPage.Components[i].Group < statusPage.Components[j].Group
		}
		return statusPage.Components[i].Order < statusPage.Components[j].Order
	})

	// Sort components within each group
	for group := range statusPage.ComponentGroups {
		components := statusPage.ComponentGroups[group]
		sort.Slice(components, func(i, j int) bool {
			return components[i].Order < components[j].Order
		})
		statusPage.ComponentGroups[group] = components
	}
}

// CalculateUptimeFromIncidents calculates uptime data from incident history
func (sb *StatusBuilder) CalculateUptimeFromIncidents(component core.StatusComponent, incidents []core.Incident, periodDays int) *core.UptimeData {
	if periodDays <= 0 {
		periodDays = 90
	}

	now := time.Now()
	startDate := now.AddDate(0, 0, -periodDays+1)

	// Initialize daily downtime map
	dailyDowntime := make(map[string]int) // date -> downtime minutes

	// Filter incidents affecting this component
	for _, inc := range incidents {
		if !sb.incidentAffectsComponent(inc, component.ID) {
			continue
		}

		// Calculate incident duration
		incEnd := now
		if inc.ResolvedAt != nil {
			incEnd = *inc.ResolvedAt
		}

		// Skip if incident is entirely before our period
		if incEnd.Before(startDate) {
			continue
		}

		// Clip incident start to period start
		incStart := inc.CreatedAt
		if incStart.Before(startDate) {
			incStart = startDate
		}

		// Distribute downtime across affected days
		sb.distributeDowntime(dailyDowntime, incStart, incEnd)
	}

	// Build history array (newest first)
	history := make([]core.UptimeDay, 0, periodDays)
	totalDowntime := 0

	for i := 0; i < periodDays; i++ {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		downtime := dailyDowntime[dateStr]
		if downtime > minutesPerDay {
			downtime = minutesPerDay
		}
		totalDowntime += downtime

		status := sb.getStatusFromDowntime(downtime)
		history = append(history, core.UptimeDay{
			Date:            dateStr,
			Status:          status,
			DowntimeMinutes: downtime,
		})
	}

	// Calculate overall uptime percentage
	totalMinutes := periodDays * minutesPerDay
	uptimePercent := float64(totalMinutes-totalDowntime) / float64(totalMinutes) * 100

	// Calculate SLA stats
	sla := sb.CalculateSLAStats(history, component.Uptime.SLATarget)

	return &core.UptimeData{
		Mode:          "incidents",
		PeriodDays:    periodDays,
		UptimePercent: uptimePercent,
		History:       history,
		SLA:           sla,
	}
}

// incidentAffectsComponent checks if an incident affects a specific component
func (sb *StatusBuilder) incidentAffectsComponent(incident core.Incident, componentID string) bool {
	for _, affected := range incident.AffectedComponents {
		if affected == componentID {
			return true
		}
	}
	return false
}

// distributeDowntime distributes downtime minutes across affected days
func (sb *StatusBuilder) distributeDowntime(dailyDowntime map[string]int, start, end time.Time) {
	current := start

	for current.Before(end) {
		dateStr := current.Format("2006-01-02")

		// Calculate end of current day
		endOfDay := time.Date(current.Year(), current.Month(), current.Day(), 23, 59, 59, 0, current.Location())

		var dayEnd time.Time
		if end.Before(endOfDay) {
			dayEnd = end
		} else {
			dayEnd = endOfDay
		}

		// Calculate minutes of downtime on this day
		minutes := int(dayEnd.Sub(current).Minutes())
		if minutes > 0 {
			dailyDowntime[dateStr] += minutes
		}

		// Move to start of next day
		current = time.Date(current.Year(), current.Month(), current.Day()+1, 0, 0, 0, 0, current.Location())
	}
}

// getStatusFromDowntime determines status based on downtime minutes
func (sb *StatusBuilder) getStatusFromDowntime(downtimeMinutes int) string {
	switch {
	case downtimeMinutes == 0:
		return "operational"
	case downtimeMinutes < 30:
		return "degraded"
	case downtimeMinutes < 240: // 4 hours
		return "partial_outage"
	default:
		return "major_outage"
	}
}

// CalculateSLAStats calculates SLA statistics from uptime history
func (sb *StatusBuilder) CalculateSLAStats(history []core.UptimeDay, target float64) core.SLAStats {
	stats := core.SLAStats{
		Target: target,
	}

	// Calculate uptime for different periods
	stats.Current7d = sb.calculateUptimeForPeriod(history, 7)
	stats.Current30d = sb.calculateUptimeForPeriod(history, 30)
	stats.Current90d = sb.calculateUptimeForPeriod(history, 90)

	return stats
}

// calculateUptimeForPeriod calculates uptime percentage for a given number of days
func (sb *StatusBuilder) calculateUptimeForPeriod(history []core.UptimeDay, days int) float64 {
	if len(history) == 0 {
		return 100.0
	}

	// Limit to available history
	if days > len(history) {
		days = len(history)
	}

	totalMinutes := days * minutesPerDay
	totalDowntime := 0

	for i := 0; i < days; i++ {
		totalDowntime += history[i].DowntimeMinutes
	}

	return float64(totalMinutes-totalDowntime) / float64(totalMinutes) * 100
}
