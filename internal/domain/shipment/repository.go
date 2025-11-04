package shipment

import (
	"context"
	"ecommerce-app/internal/pkg/database/sqlc"

	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	CreateShipment(ctx context.Context, orderID, carrier, trackingNumber, status string, shippedAt, deliveredAt *time.Time) (Shipment, error)
	GetShipment(ctx context.Context, id string) (Shipment, error)
	ListShipmentsByOrder(ctx context.Context, orderID string) ([]Shipment, error)
	UpdateShipmentStatus(ctx context.Context, id, status string, shippedAt, deliveredAt *time.Time) (Shipment, error)
	DeleteShipment(ctx context.Context, id string) error
}

type repository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repository{q: q}
}

func (r *repository) CreateShipment(ctx context.Context, orderID, carrier, trackingNumber, status string, shippedAt, deliveredAt *time.Time) (Shipment, error) {
	var orderUUID pgtype.UUID
	if err := orderUUID.Scan(orderID); err != nil {
		return Shipment{}, err
	}

	var shippedAtPg pgtype.Timestamptz
	if shippedAt != nil {
		shippedAtPg = pgtype.Timestamptz{Time: *shippedAt, Valid: true}
	} else {
		shippedAtPg = pgtype.Timestamptz{Valid: false}
	}

	var deliveredAtPg pgtype.Timestamptz
	if deliveredAt != nil {
		deliveredAtPg = pgtype.Timestamptz{Time: *deliveredAt, Valid: true}
	} else {
		deliveredAtPg = pgtype.Timestamptz{Valid: false}
	}

	params := sqlc.CreateShipmentParams{
		OrderID:        orderUUID,
		Carrier:        carrier,
		TrackingNumber: pgtype.Text{String: trackingNumber, Valid: true},
		ShippedAt:      shippedAtPg,
		DeliveredAt:    deliveredAtPg,
	}

	row, err := r.q.CreateShipment(ctx, params)
	if err != nil {
		return Shipment{}, err
	}

	return mapShipment(row), nil
}

func (r *repository) GetShipment(ctx context.Context, id string) (Shipment, error) {
	var shipmentUUID pgtype.UUID
	if err := shipmentUUID.Scan(id); err != nil {
		return Shipment{}, err
	}

	row, err := r.q.GetShipment(ctx, shipmentUUID)
	if err != nil {
		return Shipment{}, err
	}

	return mapShipment(row), nil
}

func (r *repository) ListShipmentsByOrder(ctx context.Context, orderID string) ([]Shipment, error) {
	var orderUUID pgtype.UUID
	if err := orderUUID.Scan(orderID); err != nil {
		return nil, err
	}

	rows, err := r.q.ListShipmentsByOrder(ctx, orderUUID)
	if err != nil {
		return nil, err
	}

	shipments := make([]Shipment, len(rows))
	for i, row := range rows {
		shipments[i] = mapShipment(row)
	}

	return shipments, nil
}

func (r *repository) UpdateShipmentStatus(ctx context.Context, id, status string, shippedAt, deliveredAt *time.Time) (Shipment, error) {
	var shipmentUUID pgtype.UUID
	if err := shipmentUUID.Scan(id); err != nil {
		return Shipment{}, err
	}

	var shippedAtPg pgtype.Timestamptz
	if shippedAt != nil {
		shippedAtPg = pgtype.Timestamptz{Time: *shippedAt, Valid: true}
	} else {
		shippedAtPg = pgtype.Timestamptz{Valid: false}
	}

	var deliveredAtPg pgtype.Timestamptz
	if deliveredAt != nil {
		deliveredAtPg = pgtype.Timestamptz{Time: *deliveredAt, Valid: true}
	} else {
		deliveredAtPg = pgtype.Timestamptz{Valid: false}
	}

	params := sqlc.UpdateShipmentStatusParams{
		ID:          shipmentUUID,
		Status:      status,
		ShippedAt:   shippedAtPg,
		DeliveredAt: deliveredAtPg,
	}

	row, err := r.q.UpdateShipmentStatus(ctx, params)
	if err != nil {
		return Shipment{}, err
	}

	return mapShipment(row), nil
}

func (r *repository) DeleteShipment(ctx context.Context, id string) error {
	var shipmentUUID pgtype.UUID
	if err := shipmentUUID.Scan(id); err != nil {
		return err
	}

	return r.q.DeleteShipment(ctx, shipmentUUID)
}

func mapShipment(row sqlc.Shipment) Shipment {
	var shippedAt *time.Time
	if row.ShippedAt.Valid {
		shippedAt = &row.ShippedAt.Time
	}

	var deliveredAt *time.Time
	if row.DeliveredAt.Valid {
		deliveredAt = &row.DeliveredAt.Time
	}

	return Shipment{
		ID:            row.ID.String(),
		OrderID:       row.OrderID.String(),
		Carrier:       row.Carrier,
		TrackingNumber: row.TrackingNumber.String,
		Status:        row.Status,
		ShippedAt:     shippedAt,
		DeliveredAt:   deliveredAt,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
}