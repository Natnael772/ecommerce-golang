package errs

import "net/http"

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Err     error  `json:"-"`
}

// Standard domain errors
var (
	ErrNotFound       = &AppError{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrConflict        = &AppError{Code: http.StatusConflict, Message: "Resource already exists"}
	ErrBadRequest      = &AppError{Code: http.StatusBadRequest, Message: "Invalid request"}
	ErrInternal        = &AppError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ErrUnauthorized    = &AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
)


