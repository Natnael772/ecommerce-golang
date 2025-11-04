package inventory

import "time"

type Inventory struct {
	ProductID string `json:"product_id"`
	Stock     int32  `json:"stock"`
	Reserved  int32  `json:"reserved"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- Request Dto ---
type CreateInventoryRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid4"`
	Stock     int32  `json:"stock" validate:"required,min=0"`
	Reserved  int32  `json:"reserved" validate:"min=0"`
}

type UpdateInventoryRequest struct {
	Stock    *int32 `json:"stock,omitempty" validate:"min=0"`
	// Reserved *int32 `json:"reserved,omitempty" validate:"min=0"`
}