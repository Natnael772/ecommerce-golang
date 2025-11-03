package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"ecommerce-app/internal/pkg/database/sqlc"
	"ecommerce-app/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	Create(ctx context.Context,userID string, params CreateOrderRequestInput) (Order, error)
	GetByID(ctx context.Context, id string) (Order, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int32) ([]Order, error)
	GetAll(ctx context.Context, limit, offset int32) ([]Order, error)
	CountAll(ctx context.Context) (int32, error)
	CountByUserID(ctx context.Context, userID string) (int32, error)
	UpdateStatus(ctx context.Context, id, status string) (Order, error)
	Delete(ctx context.Context, id string) error
	CreateOrderPayment(ctx context.Context, params CreateOrderPaymentInput) error
}

// repository implements Repository
type repository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repository{q: q}
}

func (r *repository) Create(ctx context.Context, userID string, req CreateOrderRequestInput) (Order, error) {
	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return Order{}, err
	}

	shippingJSON, err := json.Marshal(req.ShippingInfo)
	if err != nil {
		return Order{}, err
	}

	params := sqlc.CreateOrderParams{
		UserID:       userUUID,
		OrderNumber:  req.OrderNumber,
		SubtotalCents: req.SubtotalCents,
		DiscountCents: pgtype.Int8{Int64: req.DiscountCents, Valid: true},
		TaxCents:      pgtype.Int8{Int64: req.TaxCents, Valid: true},
		ShippingCents: pgtype.Int8{Int64: req.ShippingCents, Valid: true},
		TotalCents:   req.TotalCents,
		FinalCents:   req.FinalCents,
		ShippingInfo: shippingJSON,
		Notes:        pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	}

	row, err := r.q.CreateOrder(ctx, params)
	if err != nil {
		return Order{}, err
	}

	return mapOrder(row), nil
}

func (r *repository) GetByID(ctx context.Context, id string) (Order, error) {
	var uuidID pgtype.UUID
	if err := uuidID.Scan(id); err != nil {
		return Order{}, err
	}

	row, err := r.q.GetOrderWithItemsByID(ctx, uuidID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Order{}, errs.ErrNotFound
		}
		return Order{}, err
	}

	items, err := mapOrderItems(row.Items)
	if err != nil {
		return Order{}, err
	}

	var shipping map[string]interface{}
	if err := json.Unmarshal(row.ShippingInfo, &shipping); err != nil {
		shipping = map[string]interface{}{}
	}

	return Order{
		ID:            uuid.UUID(row.ID.Bytes),
		UserID:        uuid.UUID(row.UserID.Bytes),
		OrderNumber:   row.OrderNumber,
		SubtotalCents: row.SubtotalCents,
		DiscountCents: int64(row.DiscountCents.Int64),
		TaxCents:      int64(row.TaxCents.Int64),
		ShippingCents: int64(row.ShippingCents.Int64),
		TotalCents:    row.TotalCents,
		FinalCents:    row.FinalCents,
		Currency:      row.Currency,
		Status:        row.Status,
		ShippingInfo:  shipping,
		Notes:         row.Notes.String,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
		Items:         items,
	}, nil
}

func (r *repository) GetByUserID(ctx context.Context, userID string, limit, offset int32) ([]Order, error) {
	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return nil, err
	}

	params := sqlc.GetOrdersWithItemsByUserIDParams{
		UserID: userUUID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	rows, err := r.q.GetOrdersWithItemsByUserID(ctx, params)
	if err != nil {
		return nil, err
	}

	orders := make([]Order, len(rows))
	for i, row := range rows {
		items := []OrderItem{}

		var shipping map[string]interface{}
		_ = json.Unmarshal(row.ShippingInfo, &shipping)

		orders[i] = Order{
			ID:            uuid.UUID(row.ID.Bytes),
			UserID:        uuid.UUID(row.UserID.Bytes),
			OrderNumber:   row.OrderNumber,
			SubtotalCents: row.SubtotalCents,
			DiscountCents: int64(row.DiscountCents.Int64),
			TaxCents:      int64(row.TaxCents.Int64),
			ShippingCents: int64(row.ShippingCents.Int64),
			TotalCents:    row.TotalCents,
			FinalCents:    row.FinalCents,
			Currency:      row.Currency,
			Status:        row.Status,
			ShippingInfo:  shipping,
			Notes:         row.Notes.String,
			CreatedAt:     row.CreatedAt.Time,
			UpdatedAt:     row.UpdatedAt.Time,
			Items:         items,
		}
	}

	return orders, nil
}

func (r *repository) GetAll(ctx context.Context, limit, offset int32) ([]Order, error) {
	params := sqlc.GetOrdersWithItemsParams{
		Limit:  limit,
		Offset: offset,
	}

	rows, err := r.q.GetOrdersWithItems(ctx, params)
	if err != nil {
		return nil, err
	}

	orders := make([]Order, len(rows))
	for i, row := range rows {
		items, _ := mapOrderItems(row.Items)

		var shipping map[string]interface{}
		_ = json.Unmarshal(row.ShippingInfo, &shipping)

		orders[i] = Order{
			ID:            uuid.UUID(row.ID.Bytes),
			UserID:        uuid.UUID(row.UserID.Bytes),
			OrderNumber:   row.OrderNumber,
			SubtotalCents: row.SubtotalCents,
			DiscountCents: int64(row.DiscountCents.Int64),
			TaxCents:      int64(row.TaxCents.Int64),
			ShippingCents: int64(row.ShippingCents.Int64),
			TotalCents:    row.TotalCents,
			FinalCents:    row.FinalCents,
			Currency:      row.Currency,
			Status:        row.Status,
			ShippingInfo:  shipping,
			Notes:         row.Notes.String,
			CreatedAt:     row.CreatedAt.Time,
			UpdatedAt:     row.UpdatedAt.Time,
			Items:         items,
		}
	}

	return orders, nil
}

func (r *repository) CountAll(ctx context.Context) (int32, error) {
	count, err := r.q.CountOrders(ctx)
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

func (r *repository) CountByUserID(ctx context.Context, userID string) (int32, error) {
	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return 0, err
	}

	count, err := r.q.CountOrdersByUser(ctx, userUUID)
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

func (r *repository) UpdateStatus(ctx context.Context, id, status string) (Order, error) {
	var uuidID pgtype.UUID
	if err := uuidID.Scan(id); err != nil {
		return Order{}, err
	}

	params := sqlc.UpdateOrderStatusParams{
		ID:     uuidID,
		Status: status,
	}

	row, err := r.q.UpdateOrderStatus(ctx, params)
	if err != nil {
		return Order{}, err
	}

	items := []OrderItem{}

	var shipping map[string]interface{}
	_ = json.Unmarshal(row.ShippingInfo, &shipping)

	return Order{
		ID:            uuid.UUID(row.ID.Bytes),
		UserID:        uuid.UUID(row.UserID.Bytes),
		OrderNumber:   row.OrderNumber,
		SubtotalCents: row.SubtotalCents,
		DiscountCents: int64(row.DiscountCents.Int64),
		TaxCents:      int64(row.TaxCents.Int64),
		ShippingCents: int64(row.ShippingCents.Int64),
		TotalCents:    row.TotalCents,
		FinalCents:    row.FinalCents,
		Currency:      row.Currency,
		Status:        row.Status,
		ShippingInfo:  shipping,
		Notes:         row.Notes.String,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
		Items:         items,
	}, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	var uuidID pgtype.UUID
	if err := uuidID.Scan(id); err != nil {
		return err
	}

	return r.q.DeleteOrder(ctx, uuidID)
}

func (r *repository) CreateOrderPayment(ctx context.Context, req CreateOrderPaymentInput) error {
	var orderUUID pgtype.UUID
	if err := orderUUID.Scan(req.OrderID); err != nil {
		return err
	}

	params := sqlc.CreatePaymentParams{
		OrderID:       orderUUID,
		Provider:      req.Provider,
		ProviderTxnID: pgtype.Text{String: req.ProviderTxnID, Valid: req.ProviderTxnID != ""},
		PaymentMethod: req.PaymentMethod,		
		AmountCents:   req.AmountCents,
		Currency:      req.Currency,
		Status:        req.Status,
	}

	_, err := r.q.CreatePayment(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

// --- Helper to map JSON items ---
func mapOrderItems(itemsJSON interface{}) ([]OrderItem, error) {
	bytes, err := json.Marshal(itemsJSON)
	if err != nil {
		return nil, err
	}
	var items []OrderItem
	if err := json.Unmarshal(bytes, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func mapOrder(row sqlc.Order) Order {
	return Order{
		ID:            uuid.UUID(row.ID.Bytes),
		UserID:        uuid.UUID(row.UserID.Bytes),
		OrderNumber:   row.OrderNumber,
		SubtotalCents: row.SubtotalCents,
		DiscountCents: int64(row.DiscountCents.Int64),
		TaxCents:      int64(row.TaxCents.Int64),
		ShippingCents: int64(row.ShippingCents.Int64),
		TotalCents:    row.TotalCents,
		FinalCents:    row.FinalCents,
		Currency:      row.Currency,
		Status:        row.Status,
		Notes:         row.Notes.String,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
		Items: 	   []OrderItem{},
}}