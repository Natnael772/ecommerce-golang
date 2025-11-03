package uploader

import (
	"context"
	"fmt"
	"net/http"
)

// Middleware to handle multiple file uploads
func UploadMultipleFile(fieldName string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
            }

            r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
            if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
                http.Error(w, "Files too large or malformed form", http.StatusBadRequest)
                return
            }

            files := r.MultipartForm.File[fieldName]
            if len(files) == 0 {
                http.Error(w, "No files uploaded", http.StatusBadRequest)
                return
            }

            var savedFiles []string
            for _, fh := range files {
                file, err := fh.Open()
                if err != nil {
                    http.Error(w, fmt.Sprintf("Failed to open uploaded file: %v", err), http.StatusInternalServerError)
                    return
                }

                savedFilename, err := saveUploadedFile(file, fh)
                file.Close()
                if err != nil {
                    http.Error(w, fmt.Sprintf("Failed to save file: %v", err), http.StatusInternalServerError)
                    return
                }
                savedFiles = append(savedFiles, savedFilename)
            }

            ctx := context.WithValue(r.Context(), UploadedFilesKey, savedFiles)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUploadedFiles(r *http.Request) ([]string, bool) {
    savedFiles, ok := r.Context().Value(UploadedFilesKey).([]string)
    return savedFiles, ok
}
