package uploader

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func saveUploadedFile(file multipart.File, header *multipart.FileHeader) (string, error) {
    filename := sanitizeFilename(header.Filename)
    timestamp := time.Now().Format("20060102-150405")
    newFilename := fmt.Sprintf("%s-%s", timestamp, filename)
    fullPath := filepath.Join(UploadPath, newFilename)

    dst, err := os.Create(fullPath)
    if err != nil {
        return "", fmt.Errorf("could not create file: %w", err)
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        return "", fmt.Errorf("could not copy file: %w", err)
    }

    return newFilename, nil
}

func sanitizeFilename(name string) string {
    name = filepath.Base(name)
    name = strings.ReplaceAll(name, " ", "_")
    return name
}
