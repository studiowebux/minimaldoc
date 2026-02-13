package api

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Fragment handlers return HTML snippets for HTMX

// Dashboard fragments

func (r *Router) fragmentDashboardStats(c *gin.Context) {
	siteID, _ := c.Get("site_id")
	since := time.Now().Add(-24 * time.Hour)

	totalViews, uniqueSessions, _ := r.db.GetPageViewStats(c.Request.Context(), siteID.(string), since)
	avgRating, totalRatings, _ := r.db.GetRatingStats(c.Request.Context(), siteID.(string))
	subscriberCount, _ := r.db.CountSubscribers(c.Request.Context(), siteID.(string), true)

	html := fmt.Sprintf(`
		<div class="stat-card">
			<div class="stat-value">%d</div>
			<div class="stat-label">Page Views (Today)</div>
		</div>
		<div class="stat-card">
			<div class="stat-value">%d</div>
			<div class="stat-label">Unique Visitors</div>
		</div>
		<div class="stat-card">
			<div class="stat-value">%.1f</div>
			<div class="stat-label">Avg. Rating (%d)</div>
		</div>
		<div class="stat-card">
			<div class="stat-value">%d</div>
			<div class="stat-label">Subscribers</div>
		</div>
	`, totalViews, uniqueSessions, avgRating, totalRatings, subscriberCount)

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

func (r *Router) fragmentRecentPages(c *gin.Context) {
	siteID, _ := c.Get("site_id")
	since := time.Now().Add(-24 * time.Hour)

	pages, err := r.db.GetTopPages(c.Request.Context(), siteID.(string), since, 10)
	if err != nil {
		c.String(http.StatusOK, `<p class="error">Failed to load pages</p>`)
		return
	}

	if len(pages) == 0 {
		c.String(http.StatusOK, `<p class="empty">No page views recorded yet</p>`)
		return
	}

	html := `<table class="data-table">
		<thead>
			<tr>
				<th>Page</th>
				<th>Views</th>
			</tr>
		</thead>
		<tbody>`

	for _, page := range pages {
		html += fmt.Sprintf(`
			<tr>
				<td class="path">%s</td>
				<td class="count">%d</td>
			</tr>`, template.HTMLEscapeString(page.Path), page.Views)
	}

	html += `</tbody></table>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

func (r *Router) fragmentRecentFeedback(c *gin.Context) {
	siteID, _ := c.Get("site_id")

	ratings, err := r.db.ListRatings(c.Request.Context(), siteID.(string), 5, 0)
	if err != nil {
		c.String(http.StatusOK, `<p class="error">Failed to load feedback</p>`)
		return
	}

	if len(ratings) == 0 {
		c.String(http.StatusOK, `<p class="empty">No feedback received yet</p>`)
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
			</div>`, stars, template.HTMLEscapeString(rating.Path), template.HTMLEscapeString(feedback))
	}
	html += `</div>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// Analytics fragments

func (r *Router) fragmentAnalyticsStats(c *gin.Context) {
	siteID, _ := c.Get("site_id")
	period := c.DefaultQuery("period", "7d")

	since := parsePeriod(period)
	totalViews, uniqueSessions, avgDuration, _ := r.db.GetPageViewStatsExtended(c.Request.Context(), siteID.(string), since)

	// Format avg duration as mm:ss
	avgTimeStr := "--"
	if avgDuration > 0 {
		mins := int(avgDuration) / 60
		secs := int(avgDuration) % 60
		avgTimeStr = fmt.Sprintf("%d:%02d", mins, secs)
	}

	html := fmt.Sprintf(`
		<div class="stat-card">
			<div class="stat-value">%d</div>
			<div class="stat-label">Page Views</div>
		</div>
		<div class="stat-card">
			<div class="stat-value">%d</div>
			<div class="stat-label">Unique Visitors</div>
		</div>
		<div class="stat-card">
			<div class="stat-value">%s</div>
			<div class="stat-label">Avg. Time on Page</div>
		</div>
		<div class="stat-card">
			<div class="stat-value">--</div>
			<div class="stat-label">Bounce Rate</div>
		</div>
	`, totalViews, uniqueSessions, avgTimeStr)

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

func (r *Router) fragmentTopPages(c *gin.Context) {
	siteID, _ := c.Get("site_id")
	period := c.DefaultQuery("period", "7d")

	since := parsePeriod(period)
	pages, err := r.db.GetTopPages(c.Request.Context(), siteID.(string), since, 20)
	if err != nil {
		c.String(http.StatusOK, `<p class="error">Failed to load pages</p>`)
		return
	}

	if len(pages) == 0 {
		c.String(http.StatusOK, `<p class="empty">No page views recorded</p>`)
		return
	}

	html := `<table class="data-table">
		<thead>
			<tr>
				<th>Page</th>
				<th>Views</th>
			</tr>
		</thead>
		<tbody>`

	for _, page := range pages {
		html += fmt.Sprintf(`
			<tr>
				<td class="path">%s</td>
				<td class="count">%d</td>
			</tr>`, template.HTMLEscapeString(page.Path), page.Views)
	}

	html += `</tbody></table>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

func (r *Router) fragmentTrafficSources(c *gin.Context) {
	siteID, _ := c.Get("site_id")
	period := c.DefaultQuery("period", "7d")

	since := parsePeriod(period)
	sources, err := r.db.GetTrafficSources(c.Request.Context(), siteID.(string), since, 10)
	if err != nil {
		c.String(http.StatusOK, `<p class="error">Failed to load traffic sources</p>`)
		return
	}

	if len(sources) == 0 {
		c.String(http.StatusOK, `<p class="empty">No traffic data recorded</p>`)
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
			</div>`, template.HTMLEscapeString(s.Source), pct, s.Visits)
	}
	html += `</div>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

func (r *Router) fragmentViewsChart(c *gin.Context) {
	siteID, _ := c.Get("site_id")
	period := c.DefaultQuery("period", "7d")

	days := 7
	switch period {
	case "24h":
		days = 1
	case "30d":
		days = 30
	}

	dailyViews, err := r.db.GetDailyViews(c.Request.Context(), siteID.(string), days)
	if err != nil {
		c.String(http.StatusOK, `<p class="error">Failed to load chart data</p>`)
		return
	}

	if len(dailyViews) == 0 {
		c.String(http.StatusOK, `<p class="empty">No data available for chart</p>`)
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
			</div>`, template.HTMLEscapeString(d.Date), d.Views, height, dateShort)
	}
	html += `</div>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// Feedback fragments

func (r *Router) fragmentFeedbackStats(c *gin.Context) {
	siteID, _ := c.Get("site_id")

	avgRating, totalRatings, withComments, thisWeek, _ := r.db.GetRatingStatsExtended(c.Request.Context(), siteID.(string))

	html := fmt.Sprintf(`
		<div class="stat-card">
			<div class="stat-value">%.1f</div>
			<div class="stat-label">Average Rating</div>
		</div>
		<div class="stat-card">
			<div class="stat-value">%d</div>
			<div class="stat-label">Total Feedback</div>
		</div>
		<div class="stat-card">
			<div class="stat-value">%d</div>
			<div class="stat-label">With Comments</div>
		</div>
		<div class="stat-card">
			<div class="stat-value">%d</div>
			<div class="stat-label">This Week</div>
		</div>
	`, avgRating, totalRatings, withComments, thisWeek)

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

func (r *Router) fragmentFeedbackList(c *gin.Context) {
	siteID, _ := c.Get("site_id")

	ratings, err := r.db.ListRatings(c.Request.Context(), siteID.(string), 50, 0)
	if err != nil {
		c.String(http.StatusOK, `<p class="error">Failed to load feedback</p>`)
		return
	}

	if len(ratings) == 0 {
		c.String(http.StatusOK, `<p class="empty">No feedback found</p>`)
		return
	}

	html := `<table class="data-table">
		<thead>
			<tr>
				<th>Page</th>
				<th>Rating</th>
				<th>Comment</th>
				<th>Date</th>
			</tr>
		</thead>
		<tbody>`

	for _, rating := range ratings {
		comment := "-"
		if rating.Feedback.Valid && rating.Feedback.String != "" {
			comment = rating.Feedback.String
		}
		html += fmt.Sprintf(`
			<tr>
				<td class="path">%s</td>
				<td class="rating">%s</td>
				<td class="comment">%s</td>
				<td class="date">%s</td>
			</tr>`, template.HTMLEscapeString(rating.Path), renderStars(rating.Rating),
			template.HTMLEscapeString(comment), rating.CreatedAt)
	}

	html += `</tbody></table>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// Subscriber fragments

func (r *Router) fragmentSubscriberStats(c *gin.Context) {
	siteID, _ := c.Get("site_id")

	total, verified, pending, thisMonth, _ := r.db.GetSubscriberStatsExtended(c.Request.Context(), siteID.(string))

	html := fmt.Sprintf(`
		<div class="stat-card">
			<div class="stat-value">%d</div>
			<div class="stat-label">Total Subscribers</div>
		</div>
		<div class="stat-card">
			<div class="stat-value">%d</div>
			<div class="stat-label">Verified</div>
		</div>
		<div class="stat-card">
			<div class="stat-value">%d</div>
			<div class="stat-label">Pending</div>
		</div>
		<div class="stat-card">
			<div class="stat-value">%d</div>
			<div class="stat-label">This Month</div>
		</div>
	`, total, verified, pending, thisMonth)

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

func (r *Router) fragmentSubscriberList(c *gin.Context) {
	siteID, _ := c.Get("site_id")
	statusFilter := c.Query("status-filter")

	verifiedOnly := statusFilter == "verified"

	subscribers, err := r.db.ListSubscribers(c.Request.Context(), siteID.(string), verifiedOnly)
	if err != nil {
		c.String(http.StatusOK, `<p class="error">Failed to load subscribers</p>`)
		return
	}

	// Filter pending if needed
	if statusFilter == "pending" {
		var filtered []struct {
			Email       string
			Verified    bool
			SubscribedAt string
		}
		for _, s := range subscribers {
			if !s.Verified {
				filtered = append(filtered, struct {
					Email       string
					Verified    bool
					SubscribedAt string
				}{s.Email, s.Verified, s.SubscribedAt})
			}
		}
		if len(filtered) == 0 {
			c.String(http.StatusOK, `<p class="empty">No pending subscribers</p>`)
			return
		}

		html := `<table class="data-table">
			<thead>
				<tr>
					<th>Email</th>
					<th>Status</th>
					<th>Subscribed</th>
				</tr>
			</thead>
			<tbody>`

		for _, s := range filtered {
			html += fmt.Sprintf(`
				<tr>
					<td class="email">%s</td>
					<td class="status"><span class="status-pending">Pending</span></td>
					<td class="date">%s</td>
				</tr>`, template.HTMLEscapeString(s.Email), s.SubscribedAt)
		}

		html += `</tbody></table>`
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
		return
	}

	if len(subscribers) == 0 {
		c.String(http.StatusOK, `<p class="empty">No subscribers found</p>`)
		return
	}

	html := `<table class="data-table">
		<thead>
			<tr>
				<th>Email</th>
				<th>Status</th>
				<th>Subscribed</th>
			</tr>
		</thead>
		<tbody>`

	for _, s := range subscribers {
		status := "Pending"
		statusClass := "status-pending"
		if s.Verified {
			status = "Verified"
			statusClass = "status-verified"
		}
		html += fmt.Sprintf(`
			<tr>
				<td class="email">%s</td>
				<td class="status"><span class="%s">%s</span></td>
				<td class="date">%s</td>
			</tr>`, template.HTMLEscapeString(s.Email), statusClass, status, s.SubscribedAt)
	}

	html += `</tbody></table>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// Settings fragments

func (r *Router) fragmentSiteInfo(c *gin.Context) {
	siteID, _ := c.Get("site_id")

	site, err := r.db.GetSiteByID(c.Request.Context(), siteID.(string))
	if err != nil || site == nil {
		c.String(http.StatusOK, `<p class="error">Failed to load site info</p>`)
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
	`, template.HTMLEscapeString(site.Name), template.HTMLEscapeString(domain), site.CreatedAt)

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
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
