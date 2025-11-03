package product

import "github.com/google/uuid"

type CreateProductRequest struct {
	SKU         string    `json:"sku" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=200"`
	Description string    `json:"description" validate:"required,min=10"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	PriceCents  int32     `json:"price_cents" validate:"required,gt=0"`
	Currency    string    `json:"currency" validate:"required,len=3"`
	Attributes  map[string]interface{}     `json:"attributes,omitempty"`
	MainImageUrl string   `json:"main_image_url" validate:"required,url"`
	Images      []string  `json:"images,omitempty"`
	DiscountPercent int32  `json:"discount_percent,omitempty" validate:"omitempty,gte=0,lte=100"`
	DiscountValidUntil *string `json:"discount_valid_until,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type UpdateProductRequest struct {
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Description *string    `json:"description,omitempty" validate:"omitempty,min=10"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	PriceCents  *int32     `json:"price_cents,omitempty" validate:"omitempty,gt=0"`
	Currency    *string    `json:"currency,omitempty" validate:"omitempty,len=3"`
	Attributes  *map[string]interface{}    `json:"attributes,omitempty"`
	MainImageUrl *string   `json:"main_image_url,omitempty" validate:"omitempty,url"`
	Images      *[]string  `json:"images,omitempty"`
	DiscountPercent *int32  `json:"discount_percent,omitempty" validate:"omitempty,gte=0,lte=100"`
	DiscountValidUntil *string `json:"discount_valid_until,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	IsActive    *bool      `json:"is_active,omitempty"`
}

type UpdatePriceRequest struct {
	PriceCents int32 `json:"price_cents" validate:"required,gt=0"`
}