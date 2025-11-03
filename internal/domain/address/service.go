package address

import (
	"context"
	"ecommerce-app/internal/pkg/errs"
	"errors"
)

type Service interface {
	CreateAddress(ctx context.Context, userID string, req CreateAddressRequest) (Address, *errs.AppError)
	GetAddressByID(ctx context.Context, id string) (Address, *errs.AppError)
	GetAddressesByUserID(ctx context.Context, userID string) ([]Address, *errs.AppError)
	UpdateAddress(ctx context.Context, id string, req UpdateAddressRequest) (Address, *errs.AppError)
	DeleteAddress(ctx context.Context, id string) *errs.AppError
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateAddress(ctx context.Context, userID string, req CreateAddressRequest) (Address, *errs.AppError) {
	address, err := s.repo.Create(ctx, userID, req)
	if err != nil {
		return Address{}, errs.ErrInternal.WithMessage("Failed to create address")
	}

	return address, nil
}

func (s *service) GetAddressByID(ctx context.Context, id string) (Address, *errs.AppError) {
	address, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Address{}, errs.ErrInternal.WithMessage("Failed to fetch address")
	}

	return address, nil
}

func (s *service) GetAddressesByUserID(ctx context.Context, userID string) ([]Address, *errs.AppError) {
	addresses, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
            return []Address{}, errs.ErrNotFound
        }

		return nil, errs.ErrInternal.WithMessage("Failed to fetch addresses")
	}

	return addresses, nil
}

func (s *service) UpdateAddress(ctx context.Context, id string, req UpdateAddressRequest) (Address, *errs.AppError) {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
            return Address{}, errs.ErrNotFound
        }

		return Address{}, errs.ErrInternal.WithMessage("Failed to fetch address")
	}

	address, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return Address{}, errs.ErrInternal.WithMessage("Failed to update address")
	}

	return address, nil
}

func (s *service) DeleteAddress(ctx context.Context, id string) *errs.AppError {
	addr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
            return errs.ErrNotFound.WithMessage("Resource not found with this id")
        }

		return errs.ErrInternal.WithMessage("Failed to fetch address")
	}
	
	if addr.IsDeleted{
		return errs.ErrNotFound.WithMessage("Resource not found with this id") 
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return errs.ErrInternal.WithMessage("Failed to delete address")
	}

	return nil
}