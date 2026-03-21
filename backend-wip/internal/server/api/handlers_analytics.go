package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TrackRequest represents a page view tracking request.
type TrackRequest struct {
	SiteID      string `json:"site_id" binding:"required,max=64"`
	Path        string `json:"path" binding:"required,max=512"`
	Referrer    string `json:"referrer" binding:"max=1024"`
	Country     string `json:"country" binding:"max=8"`
	DeviceType  string `json:"device_type" binding:"max=16"`
	Browser     string `json:"browser" binding:"max=64"`
	OS          string `json:"os" binding:"max=64"`
	SessionHash string `json:"session_hash" binding:"max=64"`
}

func (r *Router) trackPageView(c *gin.Context) {
	var req TrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	// Validate site exists to prevent data pollution
	site, err := r.db.GetSiteByID(c.Request.Context(), req.SiteID)
	if err != nil || site == nil {
		respondBadRequest(c, ErrSiteInvalid, "invalid site_id")
		return
	}

	err = r.db.RecordPageView(c.Request.Context(),
		req.SiteID, req.Path, req.Referrer, req.Country,
		req.DeviceType, req.Browser, req.OS, req.SessionHash)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "failed to record")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "tracked"})
}

// EventRequest represents a custom event tracking request.
type EventRequest struct {
	SiteID      string `json:"site_id" binding:"required,max=64"`
	Name        string `json:"name" binding:"required,max=128"`
	Category    string `json:"category" binding:"max=64"`
	Path        string `json:"path" binding:"max=512"`
	Value       string `json:"value" binding:"max=1024"`
	SessionHash string `json:"session_hash" binding:"max=64"`
}

func (r *Router) trackEvent(c *gin.Context) {
	var req EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrInvalidRequest, "invalid request")
		return
	}

	// Validate site exists to prevent orphaned events
	site, err := r.db.GetSiteByID(c.Request.Context(), req.SiteID)
	if err != nil || site == nil {
		respondBadRequest(c, ErrSiteInvalid, "invalid site_id")
		return
	}

	err = r.db.RecordEvent(c.Request.Context(),
		req.SiteID, req.Name, req.Category, req.Path, req.Value, req.SessionHash)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "failed to record event")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "tracked"})
}

// DurationRequest represents a page duration update.
type DurationRequest struct {
	SiteID      string `json:"site_id" binding:"required,max=64"`
	Path        string `json:"path" binding:"required,max=512"`
	Duration    int    `json:"duration" binding:"required,min=1,max=3600"`
	SessionHash string `json:"session_hash" binding:"max=64"`
	IsBounce    *bool  `json:"is_bounce"`
}

func (r *Router) trackDuration(c *gin.Context) {
	var req DurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	// Cap duration at a reasonable maximum (30 minutes = 1800 seconds)
	duration := req.Duration
	if duration > 1800 {
		duration = 1800
	}

	err := r.db.UpdatePageViewDurationAndBounce(c.Request.Context(), req.SiteID, req.Path, req.SessionHash, duration, req.IsBounce)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "failed to update duration")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (r *Router) analyticsSummary(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondUnauthorized(c, ErrUnauthorized, "unauthorized")
		return
	}
	since := time.Now().Add(-24 * time.Hour) // Last 24 hours

	totalViews, uniqueSessions, err := r.db.GetPageViewStats(c.Request.Context(), siteID, since)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	topPages, err := r.db.GetTopPages(c.Request.Context(), siteID, since, 10)
	if err != nil {
		slog.Error("failed to get top pages", "error", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"total_views":     totalViews,
		"unique_visitors": uniqueSessions,
		"top_pages":       topPages,
		"period":          "24h",
	})
}

func (r *Router) analyticsPages(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondUnauthorized(c, ErrUnauthorized, "unauthorized")
		return
	}
	since := time.Now().Add(-7 * 24 * time.Hour) // Last 7 days

	pages, err := r.db.GetTopPages(c.Request.Context(), siteID, since, 50)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	c.JSON(http.StatusOK, gin.H{"pages": pages, "period": "7d"})
}
