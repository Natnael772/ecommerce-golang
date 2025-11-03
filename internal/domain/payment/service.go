package payment

import (
	"context"
	"ecommerce-app/internal/domain/stripe"
	db "ecommerce-app/internal/pkg/database/sqlc"
	"ecommerce-app/internal/pkg/logger"

	"github.com/jackc/pgx/v5/pgtype"
)

type PaymentService interface {
	HandleWebhook(ctx context.Context, providerName string, payload []byte, sigHeader string) error
}

type paymentService struct {
	repo      PaymentRepository
	orderSvc  OrderProvider
}

func NewPaymentService(repo PaymentRepository, orderSvc OrderProvider) PaymentService {
	return &paymentService{repo: repo, orderSvc: orderSvc}
}


func (s *paymentService) HandleWebhook(ctx context.Context, providerName string, payload []byte, sigHeader string) error {
	stripeClient := stripe.NewStripeProvider()
	event, err := stripeClient.HandleWebhook(payload, sigHeader)
	if err != nil {
		return err
	}

	if event.Status == "INITIATED" {
		// No action needed for initiated payments
		return nil
	}

	// Map order_id to UUID
	orderUUID := pgtype.UUID{}
	_ = orderUUID.Scan(event.OrderID)

	// Get payment record for that order
	payment, err := s.repo.GetPaymentByOrderID(ctx, orderUUID)
	if err != nil {
		logger.Error("Failed to get payment by order ID: %v", err)
		return err
	}

	// Update payment status based on event
	update := db.UpdatePaymentStatusParams{
		ID:     payment.ID,
		Status: event.Status,
	}

	_, err = s.repo.UpdatePaymentStatus(ctx, update)
	if err != nil {
		return err
	}

	orderStatus := event.Status
	if event.Status == "SUCCESS" {
		orderStatus = "PAID"
	}

	_, appErr := s.orderSvc.UpdateOrderStatus(ctx, payment.OrderID.String(), orderStatus)

	if appErr != nil {
		return appErr
	}

	return nil
}
