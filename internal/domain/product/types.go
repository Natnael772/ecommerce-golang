package product

import (
	"ecommerce-app/internal/pkg/response"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID  `json:"id"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryID  uuid.UUID `json:"category_id"`
	PriceCents  int32  `json:"price_cents"`
	Currency    string `json:"currency"`
	Attributes  map[string]interface{} `json:"attributes"`
	MainImageUrl string `json:"main_image_url"`
	Images      []string  `json:"images"`
	DiscountPercent int32 `json:"discount_percent"`
	DiscountValidUntil *time.Time `json:"discount_valid_until,omitempty"`
	IsActive    bool      `json:"is_active"`
	IsDeleted   bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type ProductsWithMeta struct {
	Products []Product    `json:"products"`
	Meta     response.Meta `json:"meta"`
}