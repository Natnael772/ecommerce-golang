package inventory

import (
	"context"
	"ecommerce-app/internal/pkg/errs"
)

type Service interface {
	CreateInventory(ctx context.Context, req CreateInventoryRequest) (Inventory, *errs.AppError)
	GetInventoryByProductID(ctx context.Context, id string) (Inventory, *errs.AppError)
	UpdateInventory(ctx context.Context, id string, stock int32) (Inventory, *errs.AppError)
	DeleteInventory(ctx context.Context, id string) *errs.AppError
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateInventory(ctx context.Context, req CreateInventoryRequest) (Inventory, *errs.AppError) {

	inv, err := s.repo.CreateInventory(ctx, req.ProductID, req.Stock, req.Reserved)
	if err != nil {
		return Inventory{}, errs.ErrInternal.WithMessage("failed to create inventory")
	}

	return inv, nil
}

func (s *service) GetInventoryByProductID(ctx context.Context, id string) (Inventory, *errs.AppError) {
	res,err:= s.repo.GetInventoryByProductID(ctx, id)
	if err != nil {
		return Inventory{}, errs.ErrInternal.WithMessage("failed to get inventory")
	}

	return res, nil
}

func (s *service) UpdateInventory(ctx context.Context, productId string, stock int32) (Inventory, *errs.AppError) {
   res, err:= s.repo.UpdateInventoryStock(ctx, productId, stock)
   if err != nil {
	   return Inventory{}, errs.ErrInternal.WithMessage("failed to update inventory")
   }

   return res, nil
}

func (s *service) DeleteInventory(ctx context.Context, id string) *errs.AppError {
	err:= s.repo.DeleteInventory(ctx, id)
	if err != nil {
		return errs.ErrInternal.WithMessage("failed to delete inventory")
	}

	return nil
}