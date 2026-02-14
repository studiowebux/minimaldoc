package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
)

// Site management handlers

func (r *Router) listSites(c *gin.Context) {
	sites, err := r.db.ListSites(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sites"})
		return
	}

	type siteResponse struct {
		ID        string  `json:"id"`
		Name      string  `json:"name"`
		Domain    *string `json:"domain,omitempty"`
		CreatedAt string  `json:"created_at"`
	}

	result := make([]siteResponse, len(sites))
	for i, s := range sites {
		result[i] = siteResponse{
			ID:        s.ID,
			Name:      s.Name,
			CreatedAt: s.CreatedAt,
		}
		if s.Domain.Valid {
			result[i].Domain = &s.Domain.String
		}
	}

	c.JSON(http.StatusOK, gin.H{"sites": result})
}

func (r *Router) createSite(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		Domain string `json:"domain"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	siteID := uuid.New().String()
	apiKey, err := auth.GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate API key"})
		return
	}
	apiKeyHash := auth.HashAPIKey(apiKey)

	site, err := r.db.CreateSite(c.Request.Context(), siteID, req.Name, req.Domain, apiKeyHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create site"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"site": gin.H{
			"id":   site.ID,
			"name": site.Name,
		},
		"api_key": apiKey,
	})
}

func (r *Router) getSite(c *gin.Context) {
	id := c.Param("id")

	site, err := r.db.GetSiteByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if site == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "site not found"})
		return
	}

	response := gin.H{
		"id":         site.ID,
		"name":       site.Name,
		"created_at": site.CreatedAt,
		"updated_at": site.UpdatedAt,
	}
	if site.Domain.Valid {
		response["domain"] = site.Domain.String
	}

	c.JSON(http.StatusOK, response)
}

func (r *Router) updateSite(c *gin.Context) {
	id := c.Param("id")

	site, err := r.db.GetSiteByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if site == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "site not found"})
		return
	}

	var req struct {
		Name   string `json:"name"`
		Domain string `json:"domain"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use existing values if not provided
	name := site.Name
	if req.Name != "" {
		name = req.Name
	}
	domain := site.Domain.String
	if req.Domain != "" {
		domain = req.Domain
	}

	if err := r.db.UpdateSite(c.Request.Context(), id, name, domain); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update site"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "site updated"})
}

func (r *Router) deleteSite(c *gin.Context) {
	id := c.Param("id")

	site, err := r.db.GetSiteByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if site == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "site not found"})
		return
	}

	if err := r.db.DeleteSite(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete site"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "site deleted"})
}

func (r *Router) regenerateAPIKey(c *gin.Context) {
	id := c.Param("id")

	site, err := r.db.GetSiteByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if site == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "site not found"})
		return
	}

	apiKey, err := auth.GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate API key"})
		return
	}
	apiKeyHash := auth.HashAPIKey(apiKey)

	if err := r.db.UpdateSiteAPIKey(c.Request.Context(), id, apiKeyHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"api_key": apiKey,
		"message": "API key regenerated",
	})
}
