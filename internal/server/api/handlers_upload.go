package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/studiowebux/minimaldoc/internal/server/auth"
)

// allowedMIMETypes maps MIME types to their canonical file extensions.
var allowedMIMETypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

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

	// Read file content for magic number detection (up to 512 bytes + full content)
	content, err := io.ReadAll(io.LimitReader(file, r.config.Storage.MaxFileSize+1))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file"})
		return
	}
	if int64(len(content)) > r.config.Storage.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("File too large. Max size: %d bytes", r.config.Storage.MaxFileSize),
		})
		return
	}

	// Detect actual MIME type from file content (magic number)
	detected := mimetype.Detect(content)
	ext, ok := allowedMIMETypes[detected.String()]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         fmt.Sprintf("File type %q not allowed", detected.String()),
			"allowed_types": []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
		})
		return
	}

	// Generate safe filename: UUID + canonical extension (never use user-supplied name)
	safeFilename := uuid.New().String() + ext

	// Upload to storage with safe filename and detected MIME type
	url, storagePath, err := r.storage.Upload(c.Request.Context(), safeFilename, detected.String(), bytes.NewReader(content))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Save metadata to database (store original filename for display, safe filename for path)
	id := uuid.New().String()
	originalName := sanitizeFilename(header.Filename)
	upload, err := r.db.CreateUpload(c.Request.Context(), id, user.SiteID, user.ID, originalName, detected.String(), int64(len(content)), storagePath, url)
	if err != nil {
		_ = r.storage.Delete(c.Request.Context(), storagePath)
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

	// Verify site_id matches (prevent cross-site delete)
	if upload.SiteID != user.SiteID {
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

// sanitizeFilename removes path components and limits length for safe display.
func sanitizeFilename(name string) string {
	// Strip directory components
	name = filepath.Base(name)
	// Remove any path separators that might survive
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	// Limit length
	if len(name) > 255 {
		name = name[:255]
	}
	return name
}
