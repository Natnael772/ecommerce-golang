package order

import (
	"context"
	"ecommerce-app/internal/domain/product"
	"ecommerce-app/internal/pkg/errs"
	"ecommerce-app/internal/pkg/response"
	"time"

	"github.com/google/uuid"
)

// --- Domain Models ---
type Order struct {
	ID            uuid.UUID   `json:"id"`
	UserID        uuid.UUID   `json:"user_id"`
	OrderNumber   string      `json:"order_number"`
	SubtotalCents int64       `json:"subtotal_cents"`
	DiscountCents int64       `json:"discount_cents"`
	TaxCents      int64       `json:"tax_cents"`
	ShippingCents int64       `json:"shipping_cents"`
	TotalCents    int64       `json:"total_cents"`
	FinalCents    int64       `json:"final_cents"`
	Currency      string      `json:"currency"`
	Status        string      `json:"status"`
	ShippingInfo  interface{} `json:"shipping_info"`
	Notes         string      `json:"notes,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	Items         []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID         uuid.UUID `json:"id"`
	OrderID    uuid.UUID `json:"order_id"`
	ProductID  uuid.UUID `json:"product_id"`
	SKU        string    `json:"sku"`
	Name       string    `json:"name"`
	Qty   int       `json:"qty"`
	UnitPriceCents int64     `json:"unit_price_cents"`
	TotalPriceCents int64     `json:"total_cents"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

}

// --- Wrapper Types ---
type OrdersWithMeta struct {
	Orders []Order       `json:"orders"`
	Meta   response.Meta `json:"meta"`
}

type OrderWithClientSecret struct {
	Order        Order  `json:"order"`
	ClientSecret string `json:"client_secret"`
}

type OrderClientSecret string


// --- Dependency Injection Interface ---
type ProductProvider interface {
	GetProductByID(ctx context.Context, id string) (product.Product, *errs.AppError)
}

type PaymentProvider interface {
	// CreatePayment(ctx context.Context, req  ) (string, *errs.AppError)
}