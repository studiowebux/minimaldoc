package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/server/config"
)

// LocalStorage implements Storage using the local filesystem.
type LocalStorage struct {
	basePath string
	baseURL  string
}

// NewLocalStorage creates a local filesystem storage provider.
func NewLocalStorage(cfg config.StorageConfig) (*LocalStorage, error) {
	// Ensure the upload directory exists (0750: owner rwx, group rx, others none)
	if err := os.MkdirAll(cfg.LocalPath, 0750); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	return &LocalStorage{
		basePath: cfg.LocalPath,
		baseURL:  "/uploads", // Served via static file handler
	}, nil
}

// Upload stores a file to the local filesystem.
func (s *LocalStorage) Upload(ctx context.Context, filename string, contentType string, reader io.Reader) (string, string, error) {
	// Generate unique path
	storagePath := generatePath(sanitizeFilename(filename))
	fullPath := filepath.Join(s.basePath, storagePath)

	// Ensure directory exists (0750: owner rwx, group rx, others none)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(fullPath) // #nosec G304 -- path from trusted user configuration
	if err != nil {
		return "", "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy content
	if _, err := io.Copy(file, reader); err != nil {
		// Close before removal to avoid file lock issues on Windows
		file.Close()
		_ = os.Remove(fullPath)
		return "", "", fmt.Errorf("failed to write file: %w", err)
	}

	url := s.URL(storagePath)
	return url, storagePath, nil
}

// Delete removes a file from the local filesystem.
func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)

	// Security: ensure the path is within basePath
	absBase, err := filepath.Abs(s.basePath)
	if err != nil {
		return fmt.Errorf("invalid base path")
	}
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return fmt.Errorf("invalid path")
	}

	// Use strings.HasPrefix on cleaned absolute paths
	// Add separator to ensure we're checking directory boundaries
	if !strings.HasPrefix(absPath, absBase+string(filepath.Separator)) && absPath != absBase {
		return fmt.Errorf("invalid path: outside storage directory")
	}

	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// URL returns the public URL for a storage path.
func (s *LocalStorage) URL(path string) string {
	return s.baseURL + "/" + path
}
