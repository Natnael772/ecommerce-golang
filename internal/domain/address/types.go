package address

import (
	"ecommerce-app/internal/pkg/response"
	"time"

	"github.com/google/uuid"
)


type Address struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Label      string `json:"label"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	IsDefault  bool   `json:"is_default"`
	IsDeleted  bool   `json:"is_deleted"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AddressesWithMeta struct {
	Addresses []Address     `json:"addresses"`
	Meta      response.Meta `json:"meta"`
}

// --- Request Dto ---
type CreateAddressRequest struct {
	// UserID     string `json:"user_id" validate:"required,uuid4"`
	Label      string `json:"label" validate:"required"`
	Line1      string `json:"line1" validate:"required"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city" validate:"required"`
	State      string `json:"state" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
	Country    string `json:"country" validate:"required"`
	IsDefault  bool   `json:"is_default,omitempty"`
}

type UpdateAddressRequest struct {
	Label      *string `json:"label,omitempty"`
	Line1      *string `json:"line1,omitempty"`
	Line2      *string `json:"line2,omitempty"`
	City       *string `json:"city,omitempty"`
	State      *string `json:"state,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
	Country    *string `json:"country,omitempty"`
	IsDefault  *bool   `json:"is_default,omitempty"`
}