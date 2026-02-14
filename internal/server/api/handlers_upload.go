package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/studiowebux/minimaldoc/internal/server/auth"
)

// uploadImage handles image uploads for blog posts.
func (r *Router) uploadImage(c *gin.Context) {
	if r.storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Storage not configured"})
		return
	}

	user := c.MustGet("user").(*auth.Claims)

	// Parse multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > r.config.Storage.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("File too large. Max size: %d bytes", r.config.Storage.MaxFileSize),
		})
		return
	}

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	if !isAllowedType(contentType, r.config.Storage.AllowedTypes) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "File type not allowed",
			"allowed_types": r.config.Storage.AllowedTypes,
		})
		return
	}

	// Upload to storage
	url, storagePath, err := r.storage.Upload(c.Request.Context(), header.Filename, contentType, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Save metadata to database
	id := uuid.New().String()
	upload, err := r.db.CreateUpload(c.Request.Context(), id, user.SiteID, user.ID, header.Filename, contentType, header.Size, storagePath, url)
	if err != nil {
		// Try to clean up the uploaded file
		r.storage.Delete(c.Request.Context(), storagePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save upload metadata"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"upload": upload,
		"url":    url,
	})
}

// deleteImage handles image deletion.
func (r *Router) deleteImage(c *gin.Context) {
	if r.storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Storage not configured"})
		return
	}

	uploadID := c.Param("id")
	user := c.MustGet("user").(*auth.Claims)

	// Get upload record
	upload, err := r.db.GetUpload(c.Request.Context(), uploadID)
	if err != nil || upload == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Upload not found"})
		return
	}

	// Check ownership (admins can delete any, others only their own)
	if user.Role != "admin" && upload.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete another user's upload"})
		return
	}

	// Delete from storage
	if err := r.storage.Delete(c.Request.Context(), upload.StoragePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	// Delete from database
	if err := r.db.DeleteUpload(c.Request.Context(), uploadID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete upload record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Upload deleted"})
}

// listUploads returns uploads for the current user (or all for admin).
func (r *Router) listUploads(c *gin.Context) {
	user := c.MustGet("user").(*auth.Claims)

	var uploads interface{}
	var err error

	if user.Role == "admin" {
		uploads, err = r.db.ListUploads(c.Request.Context(), user.SiteID, "", 100, 0)
	} else {
		uploads, err = r.db.ListUploads(c.Request.Context(), user.SiteID, user.ID, 100, 0)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list uploads"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"uploads": uploads})
}

// isAllowedType checks if a content type is in the allowed list.
func isAllowedType(contentType string, allowed []string) bool {
	for _, t := range allowed {
		if t == contentType {
			return true
		}
	}
	return false
}
