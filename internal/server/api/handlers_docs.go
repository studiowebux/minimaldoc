package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
)

// DocAccessRule is the API representation of a doc access rule.
type DocAccessRule struct {
	ID           string `json:"id"`
	PathPattern  string `json:"path_pattern"`
	RequiredRole string `json:"required_role"`
	Description  string `json:"description,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CheckAccessRequest is the request for checking doc access.
type CheckAccessRequest struct {
	Path string `json:"path" form:"path" binding:"required"`
}

// CheckAccessResponse is the response for access check.
type CheckAccessResponse struct {
	Path         string `json:"path"`
	IsProtected  bool   `json:"is_protected"`
	RequiredRole string `json:"required_role,omitempty"`
	HasAccess    bool   `json:"has_access"`
	Reason       string `json:"reason,omitempty"`
}

// CreateDocAccessRequest is the request for creating a doc access rule.
type CreateDocAccessRequest struct {
	PathPattern  string `json:"path_pattern" binding:"required"`
	RequiredRole string `json:"required_role" binding:"required,oneof=viewer author editor admin"`
	Description  string `json:"description"`
}

// UpdateDocAccessRequest is the request for updating a doc access rule.
type UpdateDocAccessRequest struct {
	PathPattern  string `json:"path_pattern" binding:"required"`
	RequiredRole string `json:"required_role" binding:"required,oneof=viewer author editor admin"`
	Description  string `json:"description"`
}

// checkDocAccess checks if a path requires authentication and if user has access.
// GET /api/docs/check?path=/docs/internal/api
func (r *Router) checkDocAccess(c *gin.Context) {
	var req CheckAccessRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondBadRequest(c, ErrPathRequired, "path parameter required")
		return
	}

	siteID, err := getSiteID(c)
	if err != nil {
		// Try to get from API key header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			respondBadRequest(c, ErrMissingSiteContext, "site context required")
			return
		}
		apiKeyHash := auth.HashAPIKey(apiKey, r.config.Auth.JWTSecret)
		site, err := r.db.GetSiteByAPIKey(c.Request.Context(), apiKeyHash)
		if err != nil || site == nil {
			respondBadRequest(c, ErrInvalidAPIKey, "invalid API key")
			return
		}
		siteID = site.ID
	}

	rule, err := r.db.CheckDocAccess(c.Request.Context(), siteID, req.Path)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "failed to check access")
		return
	}

	resp := CheckAccessResponse{
		Path:        req.Path,
		IsProtected: rule != nil,
	}

	if rule != nil {
		resp.RequiredRole = rule.RequiredRole

		// Check if user is authenticated and has required role
		userRole, exists := c.Get("user_role")
		if !exists {
			resp.HasAccess = false
			resp.Reason = "authentication_required"
		} else {
			resp.HasAccess = hasRequiredRole(userRole.(string), rule.RequiredRole)
			if !resp.HasAccess {
				resp.Reason = "insufficient_role"
			}
		}
	} else {
		resp.HasAccess = true
	}

	c.JSON(http.StatusOK, resp)
}

// getDocContent serves protected document content.
// GET /api/docs/content/*path
func (r *Router) getDocContent(c *gin.Context) {
	docPath := c.Param("path")
	if docPath == "" {
		respondBadRequest(c, ErrPathRequired, "path required")
		return
	}

	// Normalize path
	docPath = "/" + strings.TrimPrefix(docPath, "/")

	siteID, err := getSiteID(c)
	if err != nil {
		respondUnauthorized(c, ErrUnauthorized, "unauthorized")
		return
	}

	// Check access rule
	rule, err := r.db.CheckDocAccess(c.Request.Context(), siteID, docPath)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "failed to check access")
		return
	}

	if rule != nil {
		// Path is protected, verify role
		userRole, exists := c.Get("user_role")
		if !exists {
			respondUnauthorized(c, ErrAuthRequired, "authentication required")
			return
		}
		if !hasRequiredRole(userRole.(string), rule.RequiredRole) {
			respondError(c, http.StatusForbidden, ErrForbidden, "insufficient permissions")
			return
		}
	}

	// Serve the file from the configured docs directory
	// The docs directory should be configurable
	docsDir := r.config.Server.DocsDir
	if docsDir == "" {
		docsDir = "public" // default
	}

	// Get absolute path of docs directory for comparison
	absDocsDir, err := filepath.Abs(docsDir)
	if err != nil {
		respondInternalError(c, ErrInternalError, "failed to resolve docs directory")
		return
	}

	// Build full file path and clean it
	// filepath.Join already cleans the path, removing .. components
	fullPath := filepath.Join(absDocsDir, docPath)

	// Get absolute path of the requested file
	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		respondBadRequest(c, ErrInvalidPath, "invalid path")
		return
	}

	// Security: verify the resolved path is still inside the docs directory
	// This prevents path traversal attacks like /../../../etc/passwd
	if !strings.HasPrefix(absFullPath, absDocsDir+string(filepath.Separator)) && absFullPath != absDocsDir {
		respondError(c, http.StatusForbidden, ErrAccessDenied, "access denied")
		return
	}

	// Check if it's a directory, serve index.html
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Try with .html extension
			fullPath = fullPath + ".html"
			info, err = os.Stat(fullPath)
			if err != nil {
				respondNotFound(c, ErrDocNotFound, "document not found")
				return
			}
		} else {
			respondInternalError(c, ErrInternalError, "failed to access document")
			return
		}
	}

	if info.IsDir() {
		fullPath = filepath.Join(fullPath, "index.html")
		if _, err := os.Stat(fullPath); err != nil {
			respondNotFound(c, ErrDocNotFound, "document not found")
			return
		}
	}

	// Security: resolve symlinks and verify still inside docs directory
	// This prevents symlink attacks that escape the allowed directory
	realPath, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		respondNotFound(c, ErrDocNotFound, "document not found")
		return
	}
	if !strings.HasPrefix(realPath, absDocsDir+string(filepath.Separator)) && realPath != absDocsDir {
		respondError(c, http.StatusForbidden, ErrAccessDenied, "access denied")
		return
	}

	// Serve the file
	file, err := os.Open(realPath)
	if err != nil {
		respondInternalError(c, ErrInternalError, "failed to read document")
		return
	}
	defer file.Close()

	// Determine content type
	contentType := "text/html"
	ext := filepath.Ext(realPath)
	switch ext {
	case ".html":
		contentType = "text/html"
	case ".md":
		contentType = "text/markdown"
	case ".json":
		contentType = "application/json"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	}

	c.Header("Content-Type", contentType)
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, file)
}

// listDocAccessRules lists all access rules for the site.
// GET /api/docs/rules
func (r *Router) listDocAccessRules(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		respondUnauthorized(c, ErrUnauthorized, "unauthorized")
		return
	}

	rules, err := r.db.ListDocAccess(c.Request.Context(), siteID)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "failed to list rules")
		return
	}

	apiRules := make([]DocAccessRule, 0, len(rules))
	for _, rule := range rules {
		apiRules = append(apiRules, DocAccessRule{
			ID:           rule.ID,
			PathPattern:  rule.PathPattern,
			RequiredRole: rule.RequiredRole,
			Description:  rule.Description.String,
			CreatedAt:    rule.CreatedAt,
			UpdatedAt:    rule.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"rules": apiRules})
}

// createDocAccessRule creates a new access rule.
// POST /api/docs/rules
func (r *Router) createDocAccessRule(c *gin.Context) {
	var req CreateDocAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	siteID, err := getSiteID(c)
	if err != nil {
		respondUnauthorized(c, ErrUnauthorized, "unauthorized")
		return
	}
	id := uuid.New().String()

	rule, err := r.db.CreateDocAccess(c.Request.Context(), id, siteID, req.PathPattern, req.RequiredRole, req.Description)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			respondError(c, http.StatusConflict, ErrRuleExists, "rule for this path pattern already exists")
			return
		}
		respondInternalError(c, ErrDatabaseError, "failed to create rule")
		return
	}

	c.JSON(http.StatusCreated, DocAccessRule{
		ID:           rule.ID,
		PathPattern:  rule.PathPattern,
		RequiredRole: rule.RequiredRole,
		Description:  rule.Description.String,
		CreatedAt:    rule.CreatedAt,
		UpdatedAt:    rule.UpdatedAt,
	})
}

// updateDocAccessRule updates an existing access rule.
// PUT /api/docs/rules/:id
func (r *Router) updateDocAccessRule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondBadRequest(c, ErrRuleIDRequired, "rule ID required")
		return
	}

	var req UpdateDocAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, ErrBadRequest, err.Error())
		return
	}

	// Verify rule exists and belongs to site
	existing, err := r.db.GetDocAccessByID(c.Request.Context(), id)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "failed to check rule")
		return
	}
	if existing == nil {
		respondNotFound(c, ErrRuleNotFound, "rule not found")
		return
	}

	siteID, err := getSiteID(c)
	if err != nil {
		respondUnauthorized(c, ErrUnauthorized, "unauthorized")
		return
	}
	if existing.SiteID != siteID {
		respondError(c, http.StatusForbidden, ErrAccessDenied, "access denied")
		return
	}

	if err := r.db.UpdateDocAccess(c.Request.Context(), id, req.PathPattern, req.RequiredRole, req.Description); err != nil {
		respondInternalError(c, ErrDatabaseError, "failed to update rule")
		return
	}

	// Fetch updated rule
	updated, err := r.db.GetDocAccessByID(c.Request.Context(), id)
	if err != nil || updated == nil {
		respondInternalError(c, ErrDatabaseError, "failed to fetch updated rule")
		return
	}
	c.JSON(http.StatusOK, DocAccessRule{
		ID:           updated.ID,
		PathPattern:  updated.PathPattern,
		RequiredRole: updated.RequiredRole,
		Description:  updated.Description.String,
		CreatedAt:    updated.CreatedAt,
		UpdatedAt:    updated.UpdatedAt,
	})
}

// deleteDocAccessRule deletes an access rule.
// DELETE /api/docs/rules/:id
func (r *Router) deleteDocAccessRule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondBadRequest(c, ErrRuleIDRequired, "rule ID required")
		return
	}

	// Verify rule exists and belongs to site
	existing, err := r.db.GetDocAccessByID(c.Request.Context(), id)
	if err != nil {
		respondInternalError(c, ErrDatabaseError, "failed to check rule")
		return
	}
	if existing == nil {
		respondNotFound(c, ErrRuleNotFound, "rule not found")
		return
	}

	siteID, err := getSiteID(c)
	if err != nil {
		respondUnauthorized(c, ErrUnauthorized, "unauthorized")
		return
	}
	if existing.SiteID != siteID {
		respondError(c, http.StatusForbidden, ErrAccessDenied, "access denied")
		return
	}

	if err := r.db.DeleteDocAccess(c.Request.Context(), id); err != nil {
		respondInternalError(c, ErrDatabaseError, "failed to delete rule")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rule deleted"})
}

// hasRequiredRole checks if a user role meets the required role level.
// Role hierarchy: admin > editor > author > viewer
func hasRequiredRole(userRole, requiredRole string) bool {
	roleLevel := map[string]int{
		"viewer": 1,
		"author": 2,
		"editor": 3,
		"admin":  4,
	}

	userLevel, ok := roleLevel[userRole]
	if !ok {
		return false
	}

	requiredLevel, ok := roleLevel[requiredRole]
	if !ok {
		return false
	}

	return userLevel >= requiredLevel
}
