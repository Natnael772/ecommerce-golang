package validator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func GetValidatedBody[T any](r *http.Request) T {
	if val := r.Context().Value(CtxKeyValidatedBody); val != nil {
		// Convert any to JSON and back to T
		var req T
		b, _ := json.Marshal(val)
		json.Unmarshal(b, &req)
		return req
	}
	var zero T
	return zero
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return "invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
	case "password":
		return "password must include upper, lower, number, and special character"
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}


func getJSONFieldName(structType any, fieldName string) string {
	t := reflect.TypeOf(structType)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if f, ok := t.FieldByName(fieldName); ok {
		tag := f.Tag.Get("json")
		if tag != "" && tag != "-" {
			return strings.Split(tag, ",")[0]
		}
	}

	return strings.ToLower(fieldName) // fallback
}
