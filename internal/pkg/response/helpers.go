package response

import "net/http"

// Common response helpers

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message)
}

func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message)
}

func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message)
}


func NoContent(w http.ResponseWriter) {
	Error(w, http.StatusNoContent, "")
}

func OK[T any](w http.ResponseWriter, data T, messages ...string) {
	message := "Resource fetched successfully"

    if len(messages) > 0 {
        message = messages[0]
    }

	Success(w, http.StatusOK, data, message)
}

func OkWithMeta[T any](w http.ResponseWriter, data T, meta Meta) {
	SuccessWithMeta(w, http.StatusOK, data, meta)
}

func Created[T any](w http.ResponseWriter, data T, messages ...string) {
	message := "Resource created successfully"

	if len(messages) > 0 {
		message = messages[0]
	}

	Success(w, http.StatusCreated, data, message)
}

func Deleted(w http.ResponseWriter) {
	Success(w, http.StatusOK, struct{}{}, "Resource deleted successfully")
}
