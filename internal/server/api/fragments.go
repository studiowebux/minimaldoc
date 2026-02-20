package api

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Fragment handlers return HTML snippets for HTMX

// Dashboard fragments

func (r *Router) fragmentDashboardStats(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}
	since := time.Now().Add(-24 * time.Hour)

	totalViews, uniqueSessions, err := r.db.GetPageViewStats(c.Request.Context(), siteID, since)
	if err != nil {
		log.Printf("Failed to get page view stats: %v", err)
	}
	avgRating, totalRatings, err := r.db.GetRatingStats(c.Request.Context(), siteID)
	if err != nil {
		log.Printf("Failed to get rating stats: %v", err)
	}
	subscriberCount, err := r.db.CountSubscribers(c.Request.Context(), siteID, true)
	if err != nil {
		log.Printf("Failed to count subscribers: %v", err)
	}

	html := buildStatCards([]StatCard{
		{Value: totalViews, Label: "Page Views (Today)"},
		{Value: uniqueSessions, Label: "Unique Visitors"},
		{Value: fmt.Sprintf("%.1f", avgRating), Label: fmt.Sprintf("Avg. Rating (%d)", totalRatings)},
		{Value: subscriberCount, Label: "Subscribers"},
	})

	respondHTML(c, html)
}

func (r *Router) fragmentRecentPages(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}
	since := time.Now().Add(-24 * time.Hour)

	pages, err := r.db.GetTopPages(c.Request.Context(), siteID, since, 10)
	if err != nil {
		respondHTMLError(c, "Failed to load pages")
		return
	}

	if len(pages) == 0 {
		respondHTMLEmpty(c, "No page views recorded yet")
		return
	}

	rows := make([]TableRow, len(pages))
	for i, page := range pages {
		rows[i] = TableRow{Cells: []string{escapeHTML(page.Path), fmt.Sprintf("%d", page.Views)}}
	}

	html := buildDataTable([]TableColumn{
		{Header: "Page", Class: "path"},
		{Header: "Views", Class: "count"},
	}, rows)

	respondHTML(c, html)
}

func (r *Router) fragmentRecentFeedback(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}

	ratings, err := r.db.ListRatings(c.Request.Context(), siteID, 5, 0)
	if err != nil {
		respondHTMLError(c, "Failed to load feedback")
		return
	}

	if len(ratings) == 0 {
		respondHTMLEmpty(c, "No feedback received yet")
		return
	}

	html := `<div class="feedback-list">`
	for _, rating := range ratings {
		stars := renderStars(rating.Rating)
		feedback := rating.Feedback.String
		if !rating.Feedback.Valid {
			feedback = ""
		}
		html += fmt.Sprintf(`
			<div class="feedback-item">
				<div class="feedback-header">
					<span class="stars">%s</span>
					<span class="path">%s</span>
				</div>
				<p class="feedback-text">%s</p>
			</div>`, stars, escapeHTML(rating.Path), escapeHTML(feedback))
	}
	html += `</div>`

	respondHTML(c, html)
}

// Analytics fragments

func (r *Router) fragmentAnalyticsStats(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}
	period := c.DefaultQuery("period", "7d")

	since := parsePeriod(period)
	totalViews, uniqueSessions, avgDuration, err := r.db.GetPageViewStatsExtended(c.Request.Context(), siteID, since)
	if err != nil {
		log.Printf("Failed to get page view stats: %v", err)
	}
	bounceRate, err := r.db.GetBounceRate(c.Request.Context(), siteID, since)
	if err != nil {
		log.Printf("Failed to get bounce rate: %v", err)
	}

	// Format avg duration as mm:ss
	avgTimeStr := "--"
	if avgDuration > 0 {
		mins := int(avgDuration) / 60
		secs := int(avgDuration) % 60
		avgTimeStr = fmt.Sprintf("%d:%02d", mins, secs)
	}

	// Format bounce rate
	bounceStr := "--"
	if bounceRate > 0 {
		bounceStr = fmt.Sprintf("%.1f%%", bounceRate)
	}

	html := buildStatCards([]StatCard{
		{Value: totalViews, Label: "Page Views"},
		{Value: uniqueSessions, Label: "Unique Visitors"},
		{Value: avgTimeStr, Label: "Avg. Time on Page"},
		{Value: bounceStr, Label: "Bounce Rate"},
	})

	respondHTML(c, html)
}

func (r *Router) fragmentTopPages(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}
	period := c.DefaultQuery("period", "7d")

	since := parsePeriod(period)
	pages, err := r.db.GetTopPages(c.Request.Context(), siteID, since, 20)
	if err != nil {
		respondHTMLError(c, "Failed to load pages")
		return
	}

	if len(pages) == 0 {
		respondHTMLEmpty(c, "No page views recorded")
		return
	}

	rows := make([]TableRow, len(pages))
	for i, page := range pages {
		rows[i] = TableRow{Cells: []string{escapeHTML(page.Path), fmt.Sprintf("%d", page.Views)}}
	}

	html := buildDataTable([]TableColumn{
		{Header: "Page", Class: "path"},
		{Header: "Views", Class: "count"},
	}, rows)

	respondHTML(c, html)
}

func (r *Router) fragmentTrafficSources(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}
	period := c.DefaultQuery("period", "7d")

	since := parsePeriod(period)
	sources, err := r.db.GetTrafficSources(c.Request.Context(), siteID, since, 10)
	if err != nil {
		respondHTMLError(c, "Failed to load traffic sources")
		return
	}

	if len(sources) == 0 {
		respondHTMLEmpty(c, "No traffic data recorded")
		return
	}

	// Calculate total for percentages
	var total int64
	for _, s := range sources {
		total += s.Visits
	}

	html := `<div class="traffic-sources">`
	for _, s := range sources {
		pct := float64(0)
		if total > 0 {
			pct = float64(s.Visits) / float64(total) * 100
		}
		html += fmt.Sprintf(`
			<div class="traffic-source">
				<span class="source-name">%s</span>
				<div class="source-bar-wrapper">
					<div class="source-bar" style="width: %.1f%%"></div>
				</div>
				<span class="source-count">%d</span>
			</div>`, escapeHTML(s.Source), pct, s.Visits)
	}
	html += `</div>`

	respondHTML(c, html)
}

func (r *Router) fragmentViewsChart(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}
	period := c.DefaultQuery("period", "7d")

	days := 7
	switch period {
	case "24h":
		days = 1
	case "30d":
		days = 30
	}

	dailyViews, err := r.db.GetDailyViews(c.Request.Context(), siteID, days)
	if err != nil {
		respondHTMLError(c, "Failed to load chart data")
		return
	}

	if len(dailyViews) == 0 {
		respondHTMLEmpty(c, "No data available for chart")
		return
	}

	// Find max for scaling
	var maxViews int64
	for _, d := range dailyViews {
		if d.Views > maxViews {
			maxViews = d.Views
		}
	}

	html := `<div class="views-chart">`
	for _, d := range dailyViews {
		height := float64(0)
		if maxViews > 0 {
			height = float64(d.Views) / float64(maxViews) * 100
		}
		// Format date as short form
		dateShort := d.Date
		if len(d.Date) >= 10 {
			dateShort = d.Date[5:10] // MM-DD
		}
		html += fmt.Sprintf(`
			<div class="chart-bar-wrapper" title="%s: %d views">
				<div class="chart-bar" style="height: %.1f%%"></div>
				<span class="chart-label">%s</span>
			</div>`, escapeHTML(d.Date), d.Views, height, dateShort)
	}
	html += `</div>`

	respondHTML(c, html)
}

// Feedback fragments

func (r *Router) fragmentFeedbackStats(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}

	avgRating, totalRatings, withComments, thisWeek, err := r.db.GetRatingStatsExtended(c.Request.Context(), siteID)
	if err != nil {
		log.Printf("Failed to get rating stats: %v", err)
	}

	html := buildStatCards([]StatCard{
		{Value: fmt.Sprintf("%.1f", avgRating), Label: "Average Rating"},
		{Value: totalRatings, Label: "Total Feedback"},
		{Value: withComments, Label: "With Comments"},
		{Value: thisWeek, Label: "This Week"},
	})

	respondHTML(c, html)
}

func (r *Router) fragmentFeedbackList(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}

	ratings, err := r.db.ListRatings(c.Request.Context(), siteID, 50, 0)
	if err != nil {
		respondHTMLError(c, "Failed to load feedback")
		return
	}

	if len(ratings) == 0 {
		respondHTMLEmpty(c, "No feedback found")
		return
	}

	rows := make([]TableRow, len(ratings))
	for i, rating := range ratings {
		comment := "-"
		if rating.Feedback.Valid && rating.Feedback.String != "" {
			comment = rating.Feedback.String
		}
		rows[i] = TableRow{Cells: []string{
			escapeHTML(rating.Path),
			renderStars(rating.Rating),
			escapeHTML(comment),
			rating.CreatedAt,
		}}
	}

	html := buildDataTable([]TableColumn{
		{Header: "Page", Class: "path"},
		{Header: "Rating", Class: "rating"},
		{Header: "Comment", Class: "comment"},
		{Header: "Date", Class: "date"},
	}, rows)

	respondHTML(c, html)
}

// Subscriber fragments

func (r *Router) fragmentSubscriberStats(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}

	total, verified, pending, thisMonth, err := r.db.GetSubscriberStatsExtended(c.Request.Context(), siteID)
	if err != nil {
		log.Printf("Failed to get subscriber stats: %v", err)
	}

	html := buildStatCards([]StatCard{
		{Value: total, Label: "Total Subscribers"},
		{Value: verified, Label: "Verified"},
		{Value: pending, Label: "Pending"},
		{Value: thisMonth, Label: "This Month"},
	})

	respondHTML(c, html)
}

func (r *Router) fragmentSubscriberList(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}
	statusFilter := c.Query("status-filter")

	verifiedOnly := statusFilter == "verified"

	subscribers, err := r.db.ListSubscribers(c.Request.Context(), siteID, verifiedOnly)
	if err != nil {
		respondHTMLError(c, "Failed to load subscribers")
		return
	}

	// Filter pending if needed
	if statusFilter == "pending" {
		var rows []TableRow
		for _, s := range subscribers {
			if !s.Verified {
				rows = append(rows, TableRow{Cells: []string{
					escapeHTML(s.Email),
					`<span class="status-pending">Pending</span>`,
					s.SubscribedAt,
				}})
			}
		}
		if len(rows) == 0 {
			respondHTMLEmpty(c, "No pending subscribers")
			return
		}

		html := buildDataTable([]TableColumn{
			{Header: "Email", Class: "email"},
			{Header: "Status", Class: "status"},
			{Header: "Subscribed", Class: "date"},
		}, rows)
		respondHTML(c, html)
		return
	}

	if len(subscribers) == 0 {
		respondHTMLEmpty(c, "No subscribers found")
		return
	}

	rows := make([]TableRow, len(subscribers))
	for i, s := range subscribers {
		status := `<span class="status-pending">Pending</span>`
		if s.Verified {
			if s.VerifiedVia != "" {
				status = `<span class="status-verified">Verified via ` + escapeHTML(s.VerifiedVia) + `</span>`
			} else {
				status = `<span class="status-verified">Verified</span>`
			}
		}
		rows[i] = TableRow{Cells: []string{
			escapeHTML(s.Email),
			status,
			s.SubscribedAt,
		}}
	}

	html := buildDataTable([]TableColumn{
		{Header: "Email", Class: "email"},
		{Header: "Status", Class: "status"},
		{Header: "Subscribed", Class: "date"},
	}, rows)

	respondHTML(c, html)
}

// Settings fragments

func (r *Router) fragmentSiteInfo(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}

	site, err := r.db.GetSiteByID(c.Request.Context(), siteID)
	if err != nil || site == nil {
		respondHTMLError(c, "Failed to load site info")
		return
	}

	domain := ""
	if site.Domain.Valid {
		domain = site.Domain.String
	}

	html := fmt.Sprintf(`
		<div class="form-group">
			<label>Site Name</label>
			<input type="text" value="%s" readonly>
		</div>
		<div class="form-group">
			<label>Domain</label>
			<input type="text" value="%s" readonly>
		</div>
		<div class="form-group">
			<label>Created</label>
			<input type="text" value="%s" readonly>
		</div>
	`, escapeHTML(site.Name), escapeHTML(domain), site.CreatedAt)

	respondHTML(c, html)
}

// Helper functions

func parsePeriod(period string) time.Time {
	switch period {
	case "24h":
		return time.Now().Add(-24 * time.Hour)
	case "7d":
		return time.Now().Add(-7 * 24 * time.Hour)
	case "30d":
		return time.Now().Add(-30 * 24 * time.Hour)
	default:
		return time.Now().Add(-7 * 24 * time.Hour)
	}
}

func renderStars(rating int) string {
	stars := ""
	for i := 0; i < 5; i++ {
		if i < rating {
			stars += "&#9733;"
		} else {
			stars += "&#9734;"
		}
	}
	return stars
}

// Event fragments

func (r *Router) fragmentEventStats(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}
	period := c.DefaultQuery("period", "7d")

	since := parsePeriod(period)
	totalEvents, err := r.db.GetTotalEventCount(c.Request.Context(), siteID, since)
	if err != nil {
		log.Printf("Failed to get total event count: %v", err)
	}
	uniqueNames, err := r.db.GetUniqueEventNames(c.Request.Context(), siteID, since)
	if err != nil {
		log.Printf("Failed to get unique event names: %v", err)
	}

	html := buildStatCards([]StatCard{
		{Value: totalEvents, Label: "Total Events"},
		{Value: uniqueNames, Label: "Unique Events"},
	})

	respondHTML(c, html)
}

func (r *Router) fragmentEventsByName(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}
	period := c.DefaultQuery("period", "7d")

	since := parsePeriod(period)
	stats, err := r.db.GetEventStats(c.Request.Context(), siteID, since)
	if err != nil {
		respondHTMLError(c, "Failed to load event stats")
		return
	}

	if len(stats) == 0 {
		respondHTMLEmpty(c, "No events recorded")
		return
	}

	rows := make([]TableRow, len(stats))
	for i, stat := range stats {
		rows[i] = TableRow{Cells: []string{escapeHTML(stat.Name), fmt.Sprintf("%d", stat.Count)}}
	}

	html := buildDataTable([]TableColumn{
		{Header: "Event Name", Class: "name"},
		{Header: "Count", Class: "count"},
	}, rows)

	respondHTML(c, html)
}

func (r *Router) fragmentRecentEvents(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondHTMLError(c, "Unauthorized")
		return
	}

	events, err := r.db.ListRecentEvents(c.Request.Context(), siteID, 20)
	if err != nil {
		respondHTMLError(c, "Failed to load recent events")
		return
	}

	if len(events) == 0 {
		respondHTMLEmpty(c, "No events recorded yet")
		return
	}

	rows := make([]TableRow, len(events))
	for i, event := range events {
		category := "-"
		if event.Category.Valid && event.Category.String != "" {
			category = event.Category.String
		}
		path := "-"
		if event.Path.Valid && event.Path.String != "" {
			path = event.Path.String
		}
		value := "-"
		if event.Value.Valid && event.Value.String != "" {
			value = event.Value.String
		}
		rows[i] = TableRow{Cells: []string{
			escapeHTML(event.Name),
			escapeHTML(category),
			escapeHTML(path),
			escapeHTML(value),
			event.CreatedAt,
		}}
	}

	html := buildDataTable([]TableColumn{
		{Header: "Event", Class: "name"},
		{Header: "Category", Class: "category"},
		{Header: "Path", Class: "path"},
		{Header: "Value", Class: "value"},
		{Header: "Date", Class: "date"},
	}, rows)

	respondHTML(c, html)
}
