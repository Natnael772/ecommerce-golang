package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"ecommerce-app/internal/pkg/response"
)

// BindJSON binds and validates JSON request body
func BindJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if r.Body == nil {
		response.BadRequest(w, "Request body is required")
		return fmt.Errorf("request body is required")
	}

	// Limit request body size to prevent abuse
	maxBytes := int64(1_048_576) // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case err.Error() == "http: request body too large":
			response.BadRequest(w, "Request body too large")
		case errors.As(err, &syntaxError):
			response.BadRequest(w, fmt.Sprintf("Invalid JSON syntax at position %d", syntaxError.Offset))
		case errors.As(err, &unmarshalTypeError):
			response.BadRequest(w, fmt.Sprintf("Invalid value for field '%s'", unmarshalTypeError.Field))
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			response.BadRequest(w, fmt.Sprintf("Unknown field %s", fieldName))
		default:
			response.BadRequest(w, "Invalid JSON payload")
		}
		return err
	}

	

	return nil
}