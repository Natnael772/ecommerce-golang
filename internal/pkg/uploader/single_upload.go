package uploader

import (
	"context"
	"fmt"
	"net/http"
)

// Middleware to handle single file upload
func UploadSingleFile(fieldName string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
            }

            r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
            if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
                http.Error(w, "File too large or malformed form", http.StatusBadRequest)
                return
            }

            file, header, err := r.FormFile(fieldName)
            if err != nil {
                http.Error(w, "Invalid or missing file", http.StatusBadRequest)
                return
            }
            defer file.Close()

            savedFilename, err := saveUploadedFile(file, header)
            if err != nil {
                http.Error(w, fmt.Sprintf("Failed to save file: %v", err), http.StatusInternalServerError)
                return
            }

            ctx := context.WithValue(r.Context(), UploadedFileKey, savedFilename)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUploadedFile(r *http.Request) (string, bool) {
    savedFile, ok := r.Context().Value(UploadedFileKey).(string)
    return savedFile, ok
}
