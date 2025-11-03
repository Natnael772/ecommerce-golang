package coupon

import "time"

type CreateCouponRequest struct {
	Code            string  `json:"code" validate:"required,alphanum,min=3,max=20"`
	Description     string  `json:"description" validate:"required,min=5,max=255"`
	DiscountPercent float64 `json:"discount_percent" validate:"required,gt=0,lt=100"`
	ValidFrom       time.Time `json:"valid_from" validate:"required"`
	ValidUntil      time.Time `json:"valid_until" validate:"required,gtfield=ValidFrom"`
	MaxUses         int32   `json:"max_uses" validate:"required,gt=0"`
	IsActive        bool    `json:"is_active" validate:"required"`
}

type UpdateCouponRequest struct {
	Code            *string  `json:"code,omitempty" validate:"omitempty,alphanum,min=3,max=20"`
	Description     *string  `json:"description,omitempty" validate:"omitempty,min=5,max=255"`
	DiscountPercent *float64 `json:"discount_percent,omitempty" validate:"omitempty,gt=0,lt=100"`
	ValidFrom       *time.Time `json:"valid_from,omitempty"`
	ValidUntil      *time.Time `json:"valid_until,omitempty" validate:"omitempty,gtfield=ValidFrom"`
	MaxUses         *int32   `json:"max_uses,omitempty" validate:"omitempty,gt=0"`
	IsActive        *bool    `json:"is_active,omitempty"`
}