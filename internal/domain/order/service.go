// service.go
package order

import (
	"context"
	"ecommerce-app/internal/domain/stripe"
	"ecommerce-app/internal/pkg/errs"
	"ecommerce-app/internal/pkg/logger"
	"ecommerce-app/internal/pkg/response"
	"ecommerce-app/pkg/idgen"
	"ecommerce-app/pkg/pagination"
)

type Service interface {
	CreateOrder(ctx context.Context, userID string, req CreateOrderRequest) (OrderWithClientSecret, *errs.AppError)
	GetOrderByID(ctx context.Context, id string) (Order, *errs.AppError)
	GetOrdersByUserID(ctx context.Context, userID string, page, perPage int) (OrdersWithMeta, *errs.AppError)
	GetAllOrders(ctx context.Context, page, perPage int) (OrdersWithMeta, *errs.AppError)
	UpdateOrderStatus(ctx context.Context, id string, status string) (Order, *errs.AppError)
	DeleteOrder(ctx context.Context, id string) *errs.AppError
	CreateOrderPayment(ctx context.Context, order Order, providerName, providerTxnID, status string) *errs.AppError
}

type service struct {
	repo Repository
	productSvc ProductProvider
}

func NewService(repo Repository, productSvc ProductProvider) Service {
	return &service{repo: repo, productSvc: productSvc}
}

func (s *service) CreateOrder(ctx context.Context, userID string, req CreateOrderRequest) (OrderWithClientSecret, *errs.AppError) {

	// Calculate order total, generate order number, etc.
	subTotalCents := int64(0)
	for _, item := range req.Items {
		// Fetch product price
		prod, appErr := s.productSvc.GetProductByID(ctx, item.ProductID)
		if appErr != nil {
			return OrderWithClientSecret{}, appErr
		}

		itemPriceCents := prod.PriceCents

		subTotalCents += int64(item.Quantity) * int64(itemPriceCents)
	}
	
	orderNumber := idgen.GenerateReadableID("ORD")
	dbReq := CreateOrderRequestInput{
		UserID:       userID,
		Items:        req.Items,
		ShippingInfo: req.ShippingInfo,
		Notes:        req.Notes,
		OrderNumber:  orderNumber,
		SubtotalCents: subTotalCents,
		TotalCents:   subTotalCents,
		FinalCents:    subTotalCents,
	}

	// Create order in DB
	order, err := s.repo.Create(ctx, userID, dbReq)
	if err != nil {
		return OrderWithClientSecret{}, errs.ErrInternal.WithMessage("Failed to create order")
	}

	// Create stripe payment intent
	stripeClient := stripe.NewStripeProvider()
	meta := map[string]string{"user_id": userID, "order_id": order.ID.String()}

	paymentIntent, err := stripeClient.CreatePaymentIntent(ctx, subTotalCents, "usd", meta)

	if err != nil {
		logger.Error("Failed to create payment intent for order %s: %v", order.ID.String(), err)
		return OrderWithClientSecret{}, errs.ErrInternal.WithMessage("Failed to create payment intent")
	}

	logger.Info("Created Stripe PaymentIntent %s for Order %s", paymentIntent.ID, order.ID.String())

	// Create payment record with INITIATED status in DB
	err = s.repo.CreateOrderPayment(ctx, CreateOrderPaymentInput{
		OrderID:       order.ID.String(),
		Provider:      "STRIPE",
		ProviderTxnID: paymentIntent.ID,
		PaymentMethod: "CREDIT_CARD",
		AmountCents:   order.FinalCents,
		Currency:      order.Currency,
		Status:        "INITIATED",
	})
	if err != nil {
		logger.Error("Failed to create payment record for order %s: %v", order.ID.String(), err)
		return OrderWithClientSecret{}, errs.ErrInternal.WithMessage("Failed to create payment record")
	}

	res:= OrderWithClientSecret{
		Order: order,
		ClientSecret: paymentIntent.ClientSecret,
	}

	return res, nil
}

func (s *service) GetOrderByID(ctx context.Context, id string) (Order, *errs.AppError) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Order{}, errs.ErrInternal.WithMessage("Failed to get order")
	}

	return order, nil
}

func (s *service) GetOrdersByUserID(ctx context.Context, userID string, page, perPage int) (OrdersWithMeta, *errs.AppError) {
	p := pagination.New(page, perPage)
	limit := int32(p.PerPage)
	offset := int32(p.Offset())

	orders, err := s.repo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return OrdersWithMeta{}, errs.ErrInternal.WithMessage("Failed to get orders for user")
	}

	total, err := s.repo.CountByUserID(ctx, userID)
	if err != nil {
		return OrdersWithMeta{}, errs.ErrInternal.WithMessage("Failed to count orders for user")
	}

	meta:= response.Meta{
		Page: p.Page,
		PerPage: p.PerPage,
		Total:   int(total),
	}

	result := OrdersWithMeta{
		Orders: orders,
		Meta:  meta,
	}

	return result, nil
}

func (s *service) GetAllOrders(ctx context.Context, page, perPage int) (OrdersWithMeta, *errs.AppError) {
    p := pagination.New(page, perPage)
	limit := int32(p.PerPage)
	offset := int32(p.Offset())

	orders, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return OrdersWithMeta{}, errs.ErrInternal.WithMessage("Failed to get all orders")
	}

	total, err := s.repo.CountAll(ctx)
	if err != nil {

		return OrdersWithMeta{}, errs.ErrInternal.WithMessage("Failed to count orders")
	}

	meta := response.Meta{
		Page:        p.Page,
		PerPage: p.PerPage,
		Total:       int(total),
	}

	result := OrdersWithMeta{
		Orders: orders,
		Meta:  meta,
	}

	return result, nil
}

func (s *service) UpdateOrderStatus(ctx context.Context, id string, status string) (Order, *errs.AppError) {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Order{}, errs.ErrInternal.WithMessage("Failed to get order for update")
	}

	updateOrder, err:= s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		return Order{}, errs.ErrInternal.WithMessage("Failed to update order status")
	}

	return updateOrder, nil
}

func (s *service) DeleteOrder(ctx context.Context, id string) *errs.AppError {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return errs.ErrInternal.WithMessage("Failed to delete order")
	}

	return nil
}

func (s *service) CreateOrderPayment(ctx context.Context, order Order, providerName, providerTxnID, status string) ( *errs.AppError) {
	 err := s.repo.CreateOrderPayment(ctx, CreateOrderPaymentInput{
		OrderID:       order.ID.String(),
		Provider:      providerName,
		ProviderTxnID: providerTxnID,
		AmountCents:   order.FinalCents,
		Currency:      order.Currency,
		Status:        status,
	})
	if err != nil {
		return errs.ErrInternal.WithMessage("Failed to create order payment")
	}

	return nil
}