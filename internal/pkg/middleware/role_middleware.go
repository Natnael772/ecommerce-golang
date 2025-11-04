package middleware

import (
	"context"
	"net/http"
	"strings"

	"ecommerce-app/configs"
	"ecommerce-app/internal/pkg/jwt"
	"ecommerce-app/internal/pkg/logger"
	"ecommerce-app/internal/pkg/response"
)

type ctxKey string

const (
	UserIDKey   ctxKey = "user_id"
	UserRoleKey ctxKey = "user_role"
	UserEmailKey ctxKey = "user_email"
)

// RoleMiddleware validates JWT and enforces role-based access control.
// It attaches user ID, role, and email to the request context for downstream handlers.
func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "missing authorization header")
				return
			}

			tokenParts := strings.SplitN(authHeader, " ", 2)
			if len(tokenParts) != 2 || !strings.EqualFold(tokenParts[0], "bearer") {
				response.Unauthorized(w, "invalid authorization header format")
				return
			}

			cfg := configs.Load()
			if cfg.JWTSecret == "" {
				response.InternalServerError(w, "server misconfiguration: missing JWT_SECRET")
				return
			}

			claims, err := jwt.Verify(cfg.JWTSecret, tokenParts[1])
			if err != nil {
				if err == jwt.ErrExpiredToken {
					response.Unauthorized(w, "token has expired")
				} else {
					response.Unauthorized(w, "invalid token")
				}
				return
			}

			
			// Role-based authorization check (if roles specified)
			if len(allowedRoles) > 0 && !IsRoleAllowed(claims.Role, allowedRoles) {
				response.Forbidden(w, "insufficient privileges to access this resource")
				return
			}
			
			logger.Info(
				"RoleMiddleware: user_id=%s role=%s path=%s",
				claims.UserID, claims.Role, r.URL.Path,
			)

			// Attach user info to context early
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}


