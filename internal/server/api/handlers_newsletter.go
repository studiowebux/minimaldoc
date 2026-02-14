package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
	"github.com/studiowebux/minimaldoc/internal/server/email"
)

// SubscribeRequest represents a newsletter subscription request.
type SubscribeRequest struct {
	SiteID string `json:"site_id" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
}

func (r *Router) subscribe(c *gin.Context) {
	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get site name for email
	site, err := r.db.GetSiteByID(c.Request.Context(), req.SiteID)
	if err != nil || site == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid site_id"})
		return
	}

	// Check if already subscribed
	existing, err := r.db.GetSubscriberByEmail(c.Request.Context(), req.SiteID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if existing != nil {
		if existing.Verified {
			c.JSON(http.StatusOK, gin.H{"status": "already_subscribed"})
			return
		}
		// Resend verification email
		r.sendVerificationEmail(site.Name, req.SiteID, req.Email, existing.VerifyToken.String)
		c.JSON(http.StatusOK, gin.H{"status": "verification_resent"})
		return
	}

	// Create subscriber with verification token
	id, err := auth.GenerateSessionToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate ID"})
		return
	}
	verifyToken, err := auth.GenerateVerificationToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	err = r.db.CreateSubscriber(c.Request.Context(), id, req.SiteID, req.Email, verifyToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to subscribe"})
		return
	}

	// Send verification email
	r.sendVerificationEmail(site.Name, req.SiteID, req.Email, verifyToken)

	c.JSON(http.StatusOK, gin.H{"status": "verification_sent"})
}

func (r *Router) sendVerificationEmail(siteName, siteID, emailAddr, token string) {
	templates := email.NewTemplates(siteName, r.config.Email.BaseURL)
	msg := templates.VerificationEmail(emailAddr, siteID, token)

	if err := r.email.Send(context.Background(), msg); err != nil {
		log.Printf("Failed to send verification email to %s: %v", emailAddr, err)
	}
}

func (r *Router) verifySubscription(c *gin.Context) {
	siteID := c.Query("site_id")
	token := c.Query("token")

	if siteID == "" || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing parameters"})
		return
	}

	// Get subscriber before verification to get email
	subscriber, err := r.db.GetSubscriberByToken(c.Request.Context(), siteID, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if subscriber == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	err = r.db.VerifySubscriber(c.Request.Context(), siteID, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	// Send welcome email
	site, err := r.db.GetSiteByID(c.Request.Context(), siteID)
	if err != nil {
		log.Printf("Failed to get site for welcome email: %v", err)
	} else if site != nil {
		r.sendWelcomeEmail(site.Name, siteID, subscriber.Email)
	}

	c.JSON(http.StatusOK, gin.H{"status": "verified"})
}

func (r *Router) sendWelcomeEmail(siteName, siteID, emailAddr string) {
	templates := email.NewTemplates(siteName, r.config.Email.BaseURL)
	msg := templates.WelcomeEmail(emailAddr, siteID)

	if err := r.email.Send(context.Background(), msg); err != nil {
		log.Printf("Failed to send welcome email to %s: %v", emailAddr, err)
	}
}

func (r *Router) unsubscribe(c *gin.Context) {
	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.db.UnsubscribeByEmail(c.Request.Context(), req.SiteID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "unsubscribed"})
}

func (r *Router) listSubscribers(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subscribers, err := r.db.ListSubscribers(c.Request.Context(), siteID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	count, err := r.db.CountSubscribers(c.Request.Context(), siteID, true)
	if err != nil {
		log.Printf("Failed to count subscribers: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"subscribers":    subscribers,
		"verified_count": count,
	})
}
