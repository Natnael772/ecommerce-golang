package errs

import (
	"errors"
	"net/http"
)

// Error implements the error interface for AppError.
func (e *AppError) Error() string {
    if e == nil {
        return ""
    }
    return e.Message
}

// Copy returns a shallow copy of the AppError.
func (e *AppError) Copy() *AppError {
	newErr := *e
	return &newErr
}

// WithMessage overrides the message.
func (e *AppError) WithMessage(msg string) *AppError {
	err := e.Copy()
	err.Message = msg
	return err
}

// WithField adds field info (useful for validation / conflict fields).
func (e *AppError) WithField(field string) *AppError {
	err := e.Copy()
	err.Field = field
	return err
}

// EnsureAppError converts any standard error into an *AppError.
func EnsureAppError(err error) *AppError {
    if err == nil {
        return nil
    }

    var appErr *AppError
    if errors.As(err, &appErr) {
        return appErr
    }

    // Fallback for unhandled errors
    return &AppError{
        Code:    http.StatusInternalServerError,
        Message: err.Error(),
		Err: 	 err,
    }
}
