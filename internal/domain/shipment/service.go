package shipment

import (
	"context"
	"ecommerce-app/internal/pkg/errs"
)

type Service interface {
	CreateShipment(ctx context.Context, req CreateShipmentRequest) (Shipment, *errs.AppError)
	GetShipment(ctx context.Context, id string) (Shipment, *errs.AppError)
	GetShipmentsByOrderID(ctx context.Context, orderID string) ([]Shipment, *errs.AppError)
	UpdateShipmentStatus(ctx context.Context, id string, req UpdateShipmentStatusRequest) (Shipment, *errs.AppError)
	DeleteShipment(ctx context.Context, id string) *errs.AppError
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateShipment(ctx context.Context, req CreateShipmentRequest) (Shipment, *errs.AppError) {
	shipment, err := s.repo.CreateShipment(ctx, req.OrderID, req.Carrier, req.TrackingNumber,req.Status, req.ShippedAt, req.DeliveredAt)
	if err != nil {
		return Shipment{}, errs.ErrInternal.WithMessage("failed to create shipment")
	}

	return shipment, nil
}

func (s *service) GetShipment(ctx context.Context, id string) (Shipment, *errs.AppError) {
	shipment, err := s.repo.GetShipment(ctx, id)
	if err != nil {
		return Shipment{}, errs.ErrInternal.WithMessage("failed to get shipment")
	}

	return shipment, nil
}

func (s *service) GetShipmentsByOrderID(ctx context.Context, orderID string) ([]Shipment, *errs.AppError) {
	shipments, err := s.repo.ListShipmentsByOrder(ctx, orderID)
	if err != nil {
		return nil, errs.ErrInternal.WithMessage("failed to list shipments by order")
	}

	return shipments, nil
}

func (s *service) UpdateShipmentStatus(ctx context.Context, id string, req UpdateShipmentStatusRequest) (Shipment, *errs.AppError) {
	shipment, err := s.repo.UpdateShipmentStatus(ctx, id, req.Status, req.ShippedAt, req.DeliveredAt)
	if err != nil {
		return Shipment{}, errs.ErrInternal.WithMessage("failed to update shipment status")
	}

	return shipment, nil
}

func (s *service) DeleteShipment(ctx context.Context, id string) *errs.AppError {
	err := s.repo.DeleteShipment(ctx, id)
	if err != nil {
		return errs.ErrInternal.WithMessage("failed to delete shipment")
	}

	return nil
}