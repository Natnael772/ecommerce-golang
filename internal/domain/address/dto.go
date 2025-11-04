package address

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