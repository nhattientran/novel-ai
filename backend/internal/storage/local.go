package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// LocalStorage handles file storage on local filesystem
type LocalStorage struct {
	uploadsDir string
	baseURL    string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(uploadsDir, baseURL string) *LocalStorage {
	return &LocalStorage{
		uploadsDir: uploadsDir,
		baseURL:    baseURL,
	}
}

// EnsureUploadsDir creates the uploads directory if it doesn't exist
func (s *LocalStorage) EnsureUploadsDir() error {
	return os.MkdirAll(s.uploadsDir, 0755)
}

// UploadImage handles image upload with validation
// Returns the public URL path for the uploaded file
func (s *LocalStorage) UploadImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate file size (max 5MB)
	const maxSize = 5 * 1024 * 1024 // 5MB
	if header.Size > maxSize {
		return "", fmt.Errorf("file size exceeds 5MB limit")
	}

	// Validate mime type
	contentType := header.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}

	if !allowedTypes[contentType] {
		return "", fmt.Errorf("invalid file type: %s (only jpeg, png, webp allowed)", contentType)
	}

	// Get file extension from content type
	ext := getExtensionFromMimeType(contentType)

	// Generate unique filename
	filename := uuid.New().String() + ext
	filepath := filepath.Join(s.uploadsDir, filename)

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return public URL path
	return "/uploads/" + filename, nil
}

// DeleteImage deletes an image file by its URL path
func (s *LocalStorage) DeleteImage(imageURL string) error {
	// Extract filename from URL
	filename := strings.TrimPrefix(imageURL, "/uploads/")
	if filename == imageURL || filename == "" {
		return fmt.Errorf("invalid image URL")
	}

	filepath := filepath.Join(s.uploadsDir, filename)

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}

	// Delete file
	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// getExtensionFromMimeType returns file extension for allowed mime types
func getExtensionFromMimeType(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
