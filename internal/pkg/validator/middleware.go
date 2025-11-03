package validator

import (
	"context"
	"ecommerce-app/internal/pkg/response"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

type contextKey string

var CtxKeyValidatedBody = contextKey("validatedBody")

// Middleware to validate request body against a struct type
// func Validate[T any]() func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			var req T

// 			// Decode JSON body
// 			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 				response.BadRequest(w, "invalid JSON body")
// 				return
// 			}

// 			// Validate struct
// 			if err := validate.Struct(req); err != nil {
// 				fieldErrors := make(map[string]string)
// 				for _, fe := range err.(validator.ValidationErrors) {
// 					fieldErrors[fe.Field()] = fe.Tag()
// 					fieldErrors[fe.Field()] = msgForTag(fe)
// 					// jsonField := getJSONFieldName(req, fe.Field())
// 					// fieldErrors[jsonField] = msgForTag(fe)
// 				}

// 				errPayload := response.GenericResponse[any]{
// 					Status:  "error",
// 					Message: "validation failed",
// 					Errors:  fieldErrors,
// 				}

// 				response.JSON(w, http.StatusBadRequest, errPayload)
// 				return
// 			}

// 			// Store validated struct in context for handler
// 			ctx := r.Context()
// 			ctx = context.WithValue(ctx, CtxKeyValidatedBody, req)
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }

func Validate[T any]() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            var req T

            if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                response.BadRequest(w, "invalid JSON body")
                return
            }

            var err error
            // check if req is a slice
            rv := reflect.ValueOf(req)
            if rv.Kind() == reflect.Slice {
                for i := 0; i < rv.Len(); i++ {
                    if e := validate.Struct(rv.Index(i).Interface()); e != nil {
                        err = e
                        break
                    }
                }
            } else {
                err = validate.Struct(req)
            }

            if err != nil {
                fieldErrors := make(map[string]string)
                if ve, ok := err.(validator.ValidationErrors); ok {
                    for _, fe := range ve {
                        fieldErrors[fe.Field()] = msgForTag(fe)
                    }
                } else {
                    // fallback for InvalidValidationError
                    response.BadRequest(w, "validation failed")
                    return
                }

                response.JSON(w, http.StatusBadRequest, response.GenericResponse[any]{
                    Status:  "error",
                    Message: "validation failed",
                    Errors:  fieldErrors,
                })
                return
            }

            ctx := context.WithValue(r.Context(), CtxKeyValidatedBody, req)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}


