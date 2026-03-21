package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

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
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	// Get site name for email
	site, err := r.db.GetSiteByID(c.Request.Context(), req.SiteID)
	if err != nil || site == nil {
		respondBadRequest(c, ErrSiteInvalid, "invalid site_id")
		return
	}

	// Check if already subscribed
	existing, err := r.db.GetSubscriberByEmail(c.Request.Context(), req.SiteID, req.Email)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if existing != nil {
		if existing.Verified {
			c.JSON(http.StatusOK, gin.H{"status": "already_subscribed"})
			return
		}
		// Generate fresh token and resend (can't retrieve hashed token)
		newToken, err := auth.GenerateVerificationToken()
		if err != nil {
			respondInternalError(c, ErrTokenGenerationFailed, "failed to generate token")
			return
		}
		tokenHash := auth.HashSessionToken(newToken)
		if err := r.db.UpdateSubscriberToken(c.Request.Context(), req.SiteID, req.Email, tokenHash); err != nil {
			respondInternalError(c, ErrSubscribeFailed, "failed to update token")
			return
		}
		r.sendVerificationEmail(site.Name, req.SiteID, req.Email, newToken)
		c.JSON(http.StatusOK, gin.H{"status": "verification_resent"})
		return
	}

	// Create subscriber with hashed verification token
	id, err := auth.GenerateSessionToken()
	if err != nil {
		respondInternalError(c, ErrIDGenFailed, "failed to generate ID")
		return
	}
	verifyToken, err := auth.GenerateVerificationToken()
	if err != nil {
		respondInternalError(c, ErrTokenGenerationFailed, "failed to generate token")
		return
	}
	tokenHash := auth.HashSessionToken(verifyToken)

	err = r.db.CreateSubscriber(c.Request.Context(), id, req.SiteID, req.Email, tokenHash)
	if err != nil {
		respondInternalError(c, ErrSubscribeFailed, "failed to subscribe")
		return
	}

	// Send verification email with raw token (DB stores hash)
	r.sendVerificationEmail(site.Name, req.SiteID, req.Email, verifyToken)

	c.JSON(http.StatusOK, gin.H{"status": "verification_sent"})
}

func (r *Router) sendVerificationEmail(siteName, siteID, emailAddr, token string) {
	templates := email.NewTemplates(siteName, r.config.Email.BaseURL)
	msg := templates.VerificationEmail(emailAddr, siteID, token)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := r.email.Send(ctx, msg); err != nil {
		slog.Error("failed to send verification email", "email", emailAddr, "error", err)
	}
}

func (r *Router) verifySubscription(c *gin.Context) {
	siteID := c.Query("site_id")
	token := c.Query("token")

	if siteID == "" || token == "" {
		respondBadRequest(c, ErrMissingParams, "missing parameters")
		return
	}

	// Hash the token to match the stored hash
	tokenHash := auth.HashSessionToken(token)

	// Get subscriber before verification to get email
	subscriber, err := r.db.GetSubscriberByToken(c.Request.Context(), siteID, tokenHash)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}
	if subscriber == nil {
		respondBadRequest(c, ErrInvalidVerifyToken, "invalid or expired token")
		return
	}

	err = r.db.VerifySubscriber(c.Request.Context(), siteID, tokenHash)
	if err != nil {
		respondBadRequest(c, ErrInvalidVerifyToken, "invalid or expired token")
		return
	}

	// Send welcome email
	site, err := r.db.GetSiteByID(c.Request.Context(), siteID)
	if err != nil {
		slog.Error("failed to get site for welcome email", "error", err)
	} else if site != nil {
		r.sendWelcomeEmail(site.Name, siteID, subscriber.Email)
	}

	c.JSON(http.StatusOK, gin.H{"status": "verified"})
}

func (r *Router) sendWelcomeEmail(siteName, siteID, emailAddr string) {
	templates := email.NewTemplates(siteName, r.config.Email.BaseURL)
	unsubToken := auth.SignUnsubscribeToken(siteID, emailAddr, r.config.Auth.JWTSecret)
	msg := templates.WelcomeEmail(emailAddr, unsubToken)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := r.email.Send(ctx, msg); err != nil {
		slog.Error("failed to send welcome email", "email", emailAddr, "error", err)
	}
}

// UnsubscribeRequest represents a token-based unsubscribe request.
type UnsubscribeRequest struct {
	Token string `json:"token" form:"token" binding:"required"`
}

func (r *Router) unsubscribe(c *gin.Context) {
	var req UnsubscribeRequest

	// Support both JSON body and query parameter (for email link clicks)
	if c.Request.Method == http.MethodGet {
		req.Token = c.Query("token")
	} else if err := c.ShouldBindJSON(&req); err != nil {
		req.Token = c.Query("token")
	}

	if req.Token == "" {
		respondBadRequest(c, ErrBadRequest, "missing unsubscribe token")
		return
	}

	siteID, email, err := auth.VerifyUnsubscribeToken(req.Token, r.config.Auth.JWTSecret)
	if err != nil {
		respondBadRequest(c, ErrBadRequest, "invalid or expired unsubscribe token")
		return
	}

	err = r.db.UnsubscribeByEmail(c.Request.Context(), siteID, email)
	if err != nil {
		respondInternalError(c, ErrUnsubscribeFailed, "failed to unsubscribe")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "unsubscribed"})
}

func (r *Router) listSubscribers(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondUnauthorized(c, ErrUnauthorized, "unauthorized")
		return
	}

	subscribers, err := r.db.ListSubscribers(c.Request.Context(), siteID, false)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "database error")
		return
	}

	count, err := r.db.CountSubscribers(c.Request.Context(), siteID, true)
	if err != nil {
		slog.Error("failed to count subscribers", "error", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"subscribers":    subscribers,
		"verified_count": count,
	})
}
