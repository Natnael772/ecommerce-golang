package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// Middleware adds logging to HTTP requests
func Middleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Create a response wrapper to capture status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			
			// Process request
			next.ServeHTTP(ww, r)
			
			// Log request details
			duration := time.Since(start)
			
			logger.Info("HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration", duration.String(),
				"bytes", ww.BytesWritten(),
				"user_agent", r.UserAgent(),
				"ip", r.RemoteAddr,
			)
		})
	}
}