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

// --- DTOs ---
type CreateOrderRequest struct {
	Items        []CreateOrderItem `json:"items" validate:"required,min=1,dive"`
	ShippingInfo interface{}       `json:"shipping_info" validate:"required"`
	Notes        string            `json:"notes,omitempty"`
}

type CreateOrderItem struct {
	ProductID string  `json:"product_id" validate:"required,uuid4"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=PENDING PAID SHIPPED CANCELLED"`
}


// --- DB (Repository) DTOs ---
type CreateOrderItemInput struct {
	ProductID string `json:"product_id" validate:"required,uuid4"`
	SKU       string `json:"sku,omitempty"`
	Name      string `json:"name" validate:"required"`
	Qty       int    `json:"qty" validate:"required,min=1"`
	PriceCents int   `json:"price_cents" validate:"required,min=0"`
}

type CreateOrderRequestInput struct {
	UserID       string                  `json:"user_id" validate:"required,uuid4"`
	OrderNumber  string                  `json:"order_number" validate:"required"`
	SubtotalCents int64                  `json:"subtotal_cents" validate:"required,min=0"`
	DiscountCents int64                  `json:"discount_cents" validate:"min=0"`
	TaxCents      int64                  `json:"tax_cents" validate:"min=0"`
	ShippingCents int64                  `json:"shipping_cents" validate:"min=0"`
	TotalCents    int64                  `json:"total_cents" validate:"required,min=0"`
	FinalCents	int64                  `json:"final_cents" validate:"required,min=0"`
	Currency      string                 `json:"currency" validate:"required,len=3"`
	Status        string                 `json:"status" validate:"required,oneof=CREATED PAID SHIPPED CANCELLED"`
	ShippingInfo interface{}             `json:"shipping_info" validate:"required"`
	Notes        string                  `json:"notes,omitempty"`
	Items        []CreateOrderItem     `json:"items" validate:"required,min=1,dive"`
}

type CreateOrderPaymentInput struct {
	OrderID       string `json:"order_id" validate:"required,uuid4"`
	Provider      string `json:"provider" validate:"required"`
	ProviderTxnID string `json:"provider_txn_id,omitempty"`
	AmountCents   int64  `json:"amount_cents" validate:"required,min=0"`
	Currency      string `json:"currency" validate:"required,len=3"`
	PaymentMethod string `json:"payment_method" validate:"required,oneof=CREDIT_CARD PAYPAL BANK_TRANSFER STRIPE APPLE_PAY GOOGLE_PAY"`
	Status        string `json:"status" validate:"required,oneof=INITIATED COMPLETED FAILED"`
}


// --- Dependency Injection Interface ---
type ProductProvider interface {
	GetProductByID(ctx context.Context, id string) (product.Product, *errs.AppError)
}

type PaymentProvider interface {
	// CreatePayment(ctx context.Context, req  ) (string, *errs.AppError)
}