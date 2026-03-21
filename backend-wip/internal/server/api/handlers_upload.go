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
)

// uploadUser holds the user info extracted safely from gin context.
type uploadUser struct {
	ID     string
	Email  string
	Role   string
	SiteID string
}

// getUploadUser extracts user info from context without panicking.
func (r *Router) getUploadUser(c *gin.Context) *uploadUser {
	u := &uploadUser{}
	if v, ok := c.Get("user_id"); ok {
		u.ID, _ = v.(string)
	}
	if v, ok := c.Get("user_email"); ok {
		u.Email, _ = v.(string)
	}
	if v, ok := c.Get("user_role"); ok {
		u.Role, _ = v.(string)
	}
	if v, ok := c.Get("site_id"); ok {
		u.SiteID, _ = v.(string)
	}
	return u
}

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
		respondError(c, http.StatusServiceUnavailable, ErrStorageNotConfig, "Storage not configured")
		return
	}

	user := r.getUploadUser(c)

	// Parse multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		respondBadRequest(c, ErrNoFileProvided, "No file provided")
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > r.config.Storage.MaxFileSize {
		respondBadRequest(c, ErrBadRequest, fmt.Sprintf("File too large. Max size: %d bytes", r.config.Storage.MaxFileSize))
		return
	}

	// Read file content for magic number detection (up to 512 bytes + full content)
	content, err := io.ReadAll(io.LimitReader(file, r.config.Storage.MaxFileSize+1))
	if err != nil {
		respondBadRequest(c, ErrBadRequest, "Failed to read file")
		return
	}
	if int64(len(content)) > r.config.Storage.MaxFileSize {
		respondBadRequest(c, ErrBadRequest, fmt.Sprintf("File too large. Max size: %d bytes", r.config.Storage.MaxFileSize))
		return
	}

	// Detect actual MIME type from file content (magic number)
	detected := mimetype.Detect(content)
	ext, ok := allowedMIMETypes[detected.String()]
	if !ok {
		respondBadRequest(c, ErrBadRequest, fmt.Sprintf("File type %q not allowed", detected.String()))
		return
	}

	// Generate safe filename: UUID + canonical extension (never use user-supplied name)
	safeFilename := uuid.New().String() + ext

	// Upload to storage with safe filename and detected MIME type
	url, storagePath, err := r.storage.Upload(c.Request.Context(), safeFilename, detected.String(), bytes.NewReader(content))
	if err != nil {
		respondInternalError(c, ErrUploadFailed, "Failed to upload file")
		return
	}

	// Save metadata to database (store original filename for display, safe filename for path)
	id := uuid.New().String()
	originalName := sanitizeFilename(header.Filename)
	upload, err := r.db.CreateUpload(c.Request.Context(), id, user.SiteID, user.ID, originalName, detected.String(), int64(len(content)), storagePath, url)
	if err != nil {
		_ = r.storage.Delete(c.Request.Context(), storagePath)
		respondInternalError(c, ErrUploadFailed, "Failed to save upload metadata")
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
		respondError(c, http.StatusServiceUnavailable, ErrStorageNotConfig, "Storage not configured")
		return
	}

	uploadID := c.Param("id")
	user := r.getUploadUser(c)

	// Get upload record
	upload, err := r.db.GetUpload(c.Request.Context(), uploadID)
	if err != nil || upload == nil {
		respondNotFound(c, ErrUploadNotFound, "Upload not found")
		return
	}

	// Verify site_id matches (prevent cross-site delete)
	if upload.SiteID != user.SiteID {
		respondNotFound(c, ErrUploadNotFound, "Upload not found")
		return
	}

	// Check ownership (admins can delete any, others only their own)
	if user.Role != "admin" && upload.UserID != user.ID {
		respondError(c, http.StatusForbidden, ErrNotUploadOwner, "Cannot delete another user's upload")
		return
	}

	// Delete from storage
	if err := r.storage.Delete(c.Request.Context(), upload.StoragePath); err != nil {
		respondInternalError(c, ErrInternalError, "Failed to delete file")
		return
	}

	// Delete from database
	if err := r.db.DeleteUpload(c.Request.Context(), uploadID); err != nil {
		respondInternalError(c, ErrDatabaseError, "Failed to delete upload record")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Upload deleted"})
}

// listUploads returns uploads for the current user (or all for admin).
func (r *Router) listUploads(c *gin.Context) {
	user := r.getUploadUser(c)

	var uploads interface{}
	var err error

	if user.Role == "admin" {
		uploads, err = r.db.ListUploads(c.Request.Context(), user.SiteID, "", 100, 0)
	} else {
		uploads, err = r.db.ListUploads(c.Request.Context(), user.SiteID, user.ID, 100, 0)
	}

	if err != nil {
		respondInternalError(c, ErrDatabaseError, "Failed to list uploads")
		return
	}

	c.JSON(http.StatusOK, gin.H{"uploads": uploads})
}

// serveUpload serves uploaded files behind auth, scoped to the user's site.
func (r *Router) serveUpload(c *gin.Context) {
	reqPath := c.Param("filepath")
	if reqPath == "" || reqPath == "/" {
		respondNotFound(c, ErrUploadNotFound, "Upload not found")
		return
	}

	// Strip leading slash for DB lookup
	storagePath := strings.TrimPrefix(reqPath, "/")

	// Look up the upload in the DB to validate it belongs to the user's site
	siteID, _ := c.Get("site_id")
	siteIDStr, _ := siteID.(string)

	upload, err := r.db.GetUploadByPath(c.Request.Context(), storagePath)
	if err != nil || upload == nil {
		respondNotFound(c, ErrUploadNotFound, "Upload not found")
		return
	}

	if upload.SiteID != siteIDStr {
		respondNotFound(c, ErrUploadNotFound, "Upload not found")
		return
	}

	// Serve the file from disk
	fullPath := filepath.Join(r.config.Storage.LocalPath, filepath.FromSlash(storagePath))

	// Prevent directory traversal
	absBase, err := filepath.Abs(r.config.Storage.LocalPath)
	if err != nil {
		respondNotFound(c, ErrUploadNotFound, "Upload not found")
		return
	}
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		respondNotFound(c, ErrUploadNotFound, "Upload not found")
		return
	}
	if !strings.HasPrefix(absPath, absBase+string(filepath.Separator)) {
		respondNotFound(c, ErrUploadNotFound, "Upload not found")
		return
	}

	c.File(fullPath)
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
