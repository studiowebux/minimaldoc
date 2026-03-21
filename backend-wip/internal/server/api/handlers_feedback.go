package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FeedbackRequest represents a feedback submission.
type FeedbackRequest struct {
	SiteID      string `json:"site_id" binding:"required"`
	Path        string `json:"path" binding:"required"`
	Rating      int    `json:"rating" binding:"required,min=1,max=5"`
	Feedback    string `json:"feedback"`
	SessionHash string `json:"session_hash"`
}

func (r *Router) submitFeedback(c *gin.Context) {
	var req FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	err := r.db.RecordRating(c.Request.Context(), req.SiteID, req.Path, req.Rating, req.Feedback, req.SessionHash)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "failed to record")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "submitted"})
}

func (r *Router) feedbackStats(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondUnauthorized(c, ErrUnauthorized, "unauthorized")
		return
	}

	avgRating, totalRatings, err := r.db.GetRatingStats(c.Request.Context(), siteID)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"average_rating": avgRating,
		"total_ratings":  totalRatings,
	})
}

func (r *Router) feedbackList(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondUnauthorized(c, ErrUnauthorized, "unauthorized")
		return
	}

	ratings, err := r.db.ListRatings(c.Request.Context(), siteID, 50, 0)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	c.JSON(http.StatusOK, gin.H{"feedback": ratings})
}
