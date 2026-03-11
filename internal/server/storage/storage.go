// Package storage provides file storage abstractions for minimaldoc-server.
package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/studiowebux/minimaldoc/internal/server/config"
)

// Storage defines the interface for file storage operations.
type Storage interface {
	// Upload stores a file and returns its public URL.
	Upload(ctx context.Context, filename string, contentType string, reader io.Reader) (url string, path string, err error)

	// Delete removes a file by its storage path.
	Delete(ctx context.Context, path string) error

	// URL returns the public URL for a storage path.
	URL(path string) string
}

// File represents an uploaded file.
type File struct {
	ID          string
	Filename    string
	ContentType string
	Size        int64
	StoragePath string
	URL         string
	UploadedAt  time.Time
}

// New creates a storage provider based on configuration.
func New(cfg config.StorageConfig) (Storage, error) {
	switch cfg.Provider {
	case "local":
		return NewLocalStorage(cfg)
	case "s3":
		return NewS3Storage(cfg)
	default:
		return NewLocalStorage(cfg)
	}
}

// generatePath creates a unique storage path for a file.
func generatePath(filename string) string {
	ext := filepath.Ext(filename)
	id := uuid.New().String()
	// Organize by date for easier management
	now := time.Now()
	return fmt.Sprintf("%d/%02d/%s%s", now.Year(), now.Month(), id, ext)
}

// sanitizeFilename removes potentially dangerous characters from filenames.
func sanitizeFilename(filename string) string {
	// Get just the base name
	filename = filepath.Base(filename)
	// Remove any path separators that might have survived
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	// Remove null bytes
	filename = strings.ReplaceAll(filename, "\x00", "")
	return filename
}
