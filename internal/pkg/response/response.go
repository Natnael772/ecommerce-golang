package response

import (
	"ecommerce-app/internal/pkg/errs"
	"encoding/json"
	"net/http"
	"time"
)

// GenericResponse is the standard envelope for API responses.
type GenericResponse[T any] struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    *T          `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta holds metadata like pagination etc.
type Meta struct {
	RequestID string    `json:"requestId,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Page      int       `json:"page,omitempty"`
	PerPage   int       `json:"per_page,omitempty"`
	Total     int       `json:"total,omitempty"`
	Next      string    `json:"next,omitempty"`
	Prev      string    `json:"prev,omitempty"`
}

func Error(w http.ResponseWriter, code int, message string) {
	resp := GenericResponse[any]{
		Status:  "error",
		Message: message,
		// Meta: &Meta{
		// 	Timestamp: time.Now().UTC(),
		// },
	}
	JSON(w, code, resp)
}

func JSON[T any](w http.ResponseWriter, code int, payload GenericResponse[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, `{"status":"error","message":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}

// Success helper for common success cases
func Success[T any](w http.ResponseWriter, code int, data T, message string) {
	resp := GenericResponse[T]{
		Status:  "success",
		Message: message,
		Data:    &data,
	}
	JSON(w, code, resp)
}

// SuccessWithMeta for responses with pagination etc.
func SuccessWithMeta[T any](w http.ResponseWriter, code int, data T, meta Meta) {
	resp := GenericResponse[T]{
		Status: "success",
		Data:   &data,
		Meta:   &meta,
	}
	JSON(w, code, resp)
}


// FromError maps an *errs.AppError to a proper JSON response.
func FromError(w http.ResponseWriter, appErr *errs.AppError) {
	if appErr == nil {
		Error(w, http.StatusInternalServerError, "Unknown error occurred")
		return
	}

	resp := GenericResponse[any]{
		Status:  "error",
		Message: appErr.Message,
		Errors:  map[string]string{"field": appErr.Field},
		Meta: &Meta{
			Timestamp: time.Now().UTC(),
		},
	}
	JSON(w, appErr.Code, resp)
}

// AppError is a convenience wrapper to handle any error gracefully.
func AppError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	appErr := errs.EnsureAppError(err)
	FromError(w, appErr)
}