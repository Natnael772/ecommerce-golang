package payment

import (
	"context"
	"ecommerce-app/internal/domain/order"
	"ecommerce-app/internal/pkg/errs"
)

type Payment struct {
	ID             string
	OrderID        string
	Provider       string
	ProviderTxnID  string
	AmountCents    int64
	Currency       string
	PaymentMethod  string
	Status         string
	Details        interface{}
	FailureReason  string
	CreatedAt      string
}

type CreatePaymentRequest struct {
	OrderID       string      `json:"order_id" validate:"required,uuid4"`
	Provider      string      `json:"provider" validate:"required"`
	ProviderTxnID string      `json:"provider_txn_id,omitempty"`
	AmountCents   int64       `json:"amount_cents" validate:"required,gt=0"`
	Currency      string      `json:"currency" validate:"required,len=3"`
	PaymentMethod string      `json:"payment_method" validate:"required,oneof=CREDIT_CARD PAYPAL BANK_TRANSFER"`
	Details       interface{} `json:"details,omitempty"`
}

type UpdatePaymentStatusRequest struct {
	ID     string `json:"id" validate:"required,uuid4"`
	Status string `json:"status" validate:"required,oneof=INITIATED SUCCESS FAILED REFUNDED"`
}	

type PaymentResponse struct {
	ID             string      `json:"id"`
	OrderID        string      `json:"order_id"`
	Provider       string      `json:"provider"`
	ProviderTxnID  string      `json:"provider_txn_id"`
	AmountCents    int64       `json:"amount_cents"`
	Currency       string      `json:"currency"`
	PaymentMethod  string      `json:"payment_method"`
	Status         string      `json:"status"`
	Details        interface{} `json:"details"`
	FailureReason  string      `json:"failure_reason,omitempty"`
	CreatedAt      string      `json:"created_at"`
}

type ProviderResponse struct {
	Provider      string
	ProviderTxnID string
	Status        string
	Details       interface{}
}

type ProviderWebhookEvent struct {
	Provider       string
	ProviderTxnID  string
	OrderID        string
	Status         string
	FailureReason  string
	RawEvent       interface{}
}


// Dependency Injection Interfaces
type OrderProvider interface {
	UpdateOrderStatus(ctx context.Context, orderID string, status string) (order.Order, *errs.AppError)
}