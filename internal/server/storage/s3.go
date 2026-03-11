package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/studiowebux/minimaldoc/internal/server/config"
)

// S3Storage implements Storage using S3-compatible object storage.
type S3Storage struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

// NewS3Storage creates an S3 storage provider.
func NewS3Storage(cfg config.StorageConfig) (Storage, error) {
	if cfg.S3Bucket == "" {
		return nil, fmt.Errorf("S3 bucket is required")
	}

	// Build S3 client options
	opts := []func(*s3.Options){
		func(o *s3.Options) {
			o.Region = cfg.S3Region
			o.Credentials = credentials.NewStaticCredentialsProvider(
				cfg.S3AccessKey,
				cfg.S3SecretKey,
				"",
			)
		},
	}

	// Custom endpoint for MinIO, R2, etc.
	if cfg.S3Endpoint != "" {
		opts = append(opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.S3Endpoint)
			o.UsePathStyle = true
		})
	}

	client := s3.New(s3.Options{}, opts...)

	// Determine public URL
	publicURL := cfg.S3PublicURL
	if publicURL == "" {
		if cfg.S3Endpoint != "" {
			publicURL = fmt.Sprintf("%s/%s", cfg.S3Endpoint, cfg.S3Bucket)
		} else {
			publicURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", cfg.S3Bucket, cfg.S3Region)
		}
	}

	return &S3Storage{
		client:    client,
		bucket:    cfg.S3Bucket,
		publicURL: publicURL,
	}, nil
}

// Upload stores a file to S3.
func (s *S3Storage) Upload(ctx context.Context, filename string, contentType string, reader io.Reader) (string, string, error) {
	storagePath := generatePath(sanitizeFilename(filename))

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(storagePath),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	url := s.URL(storagePath)
	return url, storagePath, nil
}

// Delete removes a file from S3.
func (s *S3Storage) Delete(ctx context.Context, path string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}
	return nil
}

// URL returns the public URL for a storage path.
func (s *S3Storage) URL(path string) string {
	return s.publicURL + "/" + path
}
