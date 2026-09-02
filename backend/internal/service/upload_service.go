package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"iot-rd-backend/internal/config"
)

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

func (s *UploadService) SaveUploadedFile(file *multipart.FileHeader) (*UploadResult, error) {
	// Ensure directory exists
	if err := os.MkdirAll(s.cfg.UploadPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	baseName := strings.TrimSuffix(file.Filename, ext)
	baseName = strings.ReplaceAll(baseName, " ", "_")

	uniqueName := fmt.Sprintf("%s_%d%s", baseName, time.Now().Unix(), ext)
	targetPath := filepath.Join(s.cfg.UploadPath, uniqueName)

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(targetPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	sizeStr := formatFileSize(file.Size)
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return &UploadResult{
		FileName: file.Filename,
		FilePath: targetPath,
		FileSize: sizeStr,
		MimeType: mimeType,
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
