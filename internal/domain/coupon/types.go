package coupon

import (
	"ecommerce-app/internal/pkg/response"
	"time"

	"github.com/google/uuid"
)

type Coupon struct {
	ID              uuid.UUID       `json:"id"`
	Code            string          `json:"code"`
	Description     string    `json:"description"`
	DiscountPercent int32           `json:"discount_percent"`
	ValidFrom       time.Time   `json:"valid_from"`
	ValidUntil      time.Time   `json:"valid_until"`
	MaxUses         int32     `json:"max_uses"`
	UsedCount       int32           `json:"used_count"`
	IsActive        bool            `json:"is_active"`
	IsDeleted       bool            `json:"is_deleted"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type CouponsWithMeta struct {
	Coupons []Coupon     `json:"coupons"`
	Meta    response.Meta `json:"meta"`
}

