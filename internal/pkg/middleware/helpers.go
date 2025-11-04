package middleware

import "strings"

// isRoleAllowed checks if a user role is within the allowed roles list.
func IsRoleAllowed(userRole string, allowedRoles []string) bool {
	userRole = strings.ToLower(userRole)
	for _, allowed := range allowedRoles {
		if strings.ToLower(allowed) == userRole {
			return true
		}
	}
	return false
}
