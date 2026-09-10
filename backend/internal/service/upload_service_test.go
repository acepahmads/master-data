package service

import (
	"bytes"
	"io"
	"mime/multipart"
	"testing"
)

type mockSeeker struct {
	*bytes.Reader
}

func (m *mockSeeker) Close() error {
	return nil
}

func newMockFile(filename string, content []byte) (*multipart.FileHeader, io.ReadSeeker) {
	header := &multipart.FileHeader{
		Filename: filename,
		Size:     int64(len(content)),
	}
	seeker := &mockSeeker{Reader: bytes.NewReader(content)}
	return header, seeker
}

func TestValidateUpload(t *testing.T) {
	// 1. Test Valid PNG image (with standard PNG magic bytes)
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00}
	h1, r1 := newMockFile("camera_preview.png", pngHeader)
	mime1, err := ValidateUpload(h1, r1)
	if err != nil {
		t.Fatalf("Expected valid PNG to pass, got error: %v", err)
	}
	if mime1 != "image/png" {
		t.Errorf("Expected image/png MIME, got %s", mime1)
	}

	// 2. Test Valid PDF document (with %PDF header)
	pdfHeader := []byte("%PDF-1.7 standard technical datasheet specifications")
	h2, r2 := newMockFile("sensor_datasheet.pdf", pdfHeader)
	mime2, err := ValidateUpload(h2, r2)
	if err != nil {
		t.Fatalf("Expected valid PDF to pass, got error: %v", err)
	}
	if mime2 != "application/pdf" {
		t.Errorf("Expected application/pdf MIME, got %s", mime2)
	}

	// 3. Test Disallowed Extension (.php script)
	phpContent := []byte("<?php echo 'malicious script'; ?>")
	h3, r3 := newMockFile("backdoor.php", phpContent)
	_, err = ValidateUpload(h3, r3)
	if err == nil {
		t.Fatalf("Expected .php upload to be rejected, but it was allowed")
	}

	// 4. Test Disallowed Extension (.exe executable)
	exeHeader := []byte{'M', 'Z', 0x90, 0x00, 0x03}
	h4, r4 := newMockFile("patch.exe", exeHeader)
	_, err = ValidateUpload(h4, r4)
	if err == nil {
		t.Fatalf("Expected .exe upload to be rejected, but it was allowed")
	}

	// 5. Test File Sniffing / Magic Bytes mismatch (Executable disguised as PNG)
	h5, r5 := newMockFile("photo.png", exeHeader) // Named .png but has MZ executable magic bytes
	_, err = ValidateUpload(h5, r5)
	if err == nil {
		t.Fatalf("Expected disguised executable (.png with MZ bytes) to be rejected, but it was allowed")
	}

	// 6. Test Unlisted extension (.tar.gz)
	tarContent := []byte("tar archive content placeholder")
	h6, r6 := newMockFile("archive.tar.gz", tarContent)
	_, err = ValidateUpload(h6, r6)
	if err == nil {
		t.Fatalf("Expected unlisted extension to be rejected, but it was allowed")
	}

	// 7. Test Oversized file
	largeHeader := &multipart.FileHeader{
		Filename: "huge.png",
		Size:     MaxUploadSize + 1024,
	}
	_, err = ValidateUpload(largeHeader, r1)
	if err == nil {
		t.Fatalf("Expected oversized file to be rejected, but it was allowed")
	}
}
