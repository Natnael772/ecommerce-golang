package payment

import (
	"context"
	db "ecommerce-app/internal/pkg/database/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, arg db.CreatePaymentParams) (db.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID pgtype.UUID) (db.Payment, error)
	UpdatePaymentStatus(ctx context.Context, arg db.UpdatePaymentStatusParams) (db.Payment, error)
}

type paymentRepository struct {
	q *db.Queries
}

func NewPaymentRepository(q *db.Queries) PaymentRepository {
	return &paymentRepository{q}
}

func (r *paymentRepository) CreatePayment(ctx context.Context, arg db.CreatePaymentParams) (db.Payment, error) {
	return r.q.CreatePayment(ctx, arg)
}

func (r *paymentRepository) GetPaymentByOrderID(ctx context.Context, orderID pgtype.UUID) (db.Payment, error) {
	return r.q.GetPaymentByOrderID(ctx, orderID)
}

func (r *paymentRepository) UpdatePaymentStatus(ctx context.Context, arg db.UpdatePaymentStatusParams) (db.Payment, error) {
	return r.q.UpdatePaymentStatus(ctx, arg)
}
