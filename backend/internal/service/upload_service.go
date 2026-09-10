package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"iot-rd-backend/internal/config"
)

// MaxUploadSize defines the maximum permitted file size (25 MB)
const MaxUploadSize = 25 * 1024 * 1024

// AllowedExtensions defines the permitted file types for IoT R&D documents and photos
var AllowedExtensions = map[string]string{
	// Images
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".webp": "image/webp",
	".gif":  "image/gif",
	".svg":  "image/svg+xml",

	// Documents & Datasheets
	".pdf":  "application/pdf",
	".csv":  "text/csv",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".xls":  "application/vnd.ms-excel",
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".txt":  "text/plain",
	".zip":  "application/zip",
}

// DisallowedExecutables defines dangerous extensions strictly prohibited
var DisallowedExecutables = []string{
	".exe", ".bat", ".cmd", ".sh", ".php", ".phtml", ".pl", ".py",
	".jsp", ".asp", ".aspx", ".cgi", ".dll", ".so", ".bin", ".msi",
	".vbs", ".js", ".html", ".htm", ".com", ".scr",
}

type UploadService struct {
	cfg *config.Config
}

func NewUploadService(cfg *config.Config) *UploadService {
	_ = os.MkdirAll(cfg.UploadPath, os.ModePerm)
	return &UploadService{cfg: cfg}
}

type UploadResult struct {
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
	FileSize string `json:"fileSize"`
	MimeType string `json:"mimeType"`
	URL      string `json:"url"`
}

// ValidateUpload performs strict validation on file size, extension whitelist, and binary magic bytes.
func ValidateUpload(file *multipart.FileHeader, src io.ReadSeeker) (string, error) {
	// 1. Check size limit
	if file.Size > MaxUploadSize {
		return "", fmt.Errorf("file size (%s) exceeds maximum allowed limit of %s", formatFileSize(file.Size), formatFileSize(MaxUploadSize))
	}
	if file.Size == 0 {
		return "", fmt.Errorf("empty files are not permitted")
	}

	// 2. Check extension against blacklist and whitelist
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		return "", fmt.Errorf("file must have a valid extension")
	}

	for _, bad := range DisallowedExecutables {
		if ext == bad {
			return "", fmt.Errorf("file extension '%s' is strictly prohibited for security reasons", ext)
		}
	}

	expectedMime, isAllowed := AllowedExtensions[ext]
	if !isAllowed {
		return "", fmt.Errorf("file type '%s' is not in the approved whitelist", ext)
	}

	// 3. Inspect Magic Bytes / Sniff Content
	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to inspect file content header: %w", err)
	}

	// Rewind file pointer after sniffing
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to reset file stream: %w", err)
	}

	sniffedMime := http.DetectContentType(buf[:n])

	// Strict check: Disallow executable binary signatures
	if strings.Contains(sniffedMime, "x-dosexec") || strings.Contains(sniffedMime, "x-executable") {
		return "", fmt.Errorf("upload rejected: file contains binary executable signatures")
	}

	// Validate content type matches extension family
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp", ".gif":
		if !strings.HasPrefix(sniffedMime, "image/") {
			return "", fmt.Errorf("content mismatch: file has %s extension but content is detected as %s", ext, sniffedMime)
		}
	case ".pdf":
		if sniffedMime != "application/pdf" {
			return "", fmt.Errorf("content mismatch: file has .pdf extension but content is detected as %s", sniffedMime)
		}
	case ".svg":
		// Check for potential XSS in SVG
		lowerContent := strings.ToLower(string(buf[:n]))
		if strings.Contains(lowerContent, "<script") || strings.Contains(lowerContent, "javascript:") {
			return "", fmt.Errorf("upload rejected: SVG contains embedded scripts or unsafe tags")
		}
	}

	finalMime := expectedMime
	if sniffedMime != "application/octet-stream" && sniffedMime != "text/plain" {
		finalMime = sniffedMime
	}

	return finalMime, nil
}

func (s *UploadService) SaveUploadedFile(file *multipart.FileHeader) (*UploadResult, error) {
	// Ensure directory exists
	if err := os.MkdirAll(s.cfg.UploadPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Perform Security Whitelist & Content Sniffing Validation
	validatedMime, err := ValidateUpload(file, src)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	baseName := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
	baseName = strings.ReplaceAll(baseName, " ", "_")

	uniqueName := fmt.Sprintf("%s_%d%s", baseName, time.Now().Unix(), ext)
	targetPath := filepath.Join(s.cfg.UploadPath, uniqueName)

	dst, err := os.Create(targetPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	sizeStr := formatFileSize(file.Size)

	return &UploadResult{
		FileName: file.Filename,
		FilePath: targetPath,
		FileSize: sizeStr,
		MimeType: validatedMime,
		URL:      fmt.Sprintf("/uploads/%s", uniqueName),
	}, nil
}

func formatFileSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
}
