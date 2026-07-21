package utilities

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SaveImage saves a multipart file into the specified directory
func OldSaveImage(file multipart.File, fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
	// Ensure the target directory exists
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("unable to create uploads directory (%s): %w", uploadDir, err)
	}

	ext := filepath.Ext(fileHeader.Filename) // keep original extension
	newName := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), rand.Intn(10000), ext)
	filePath := filepath.Join(uploadDir, newName)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("unable to save file: %w", err)
	}

	// Clean leading dots for relative web URL paths
	cleanPath := filepath.ToSlash(filepath.Clean(filePath))
	if len(cleanPath) > 0 && cleanPath[0] != '/' {
		cleanPath = "/" + cleanPath
	}

	// Return the relative URL/path for DB insertion
	return cleanPath, nil
}

func SaveImage(file multipart.File, fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExtensions[ext] {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("create upload directory: %w", err)
	}

	filename := fmt.Sprintf("%s%s", uuid.NewString(), ext)

	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("save file: %w", err)
	}

	// Convert filesystem path to URL path
	return "/" + filepath.ToSlash(filepath.Join(uploadDir, filename)), nil
}
