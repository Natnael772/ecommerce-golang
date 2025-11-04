package cart

import "time"

// --- Request Dto ---
type CreateCartRequest struct {
	UserID    string     `json:"user_id" validate:"required,uuid4"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type UpdateCartRequest struct {
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}