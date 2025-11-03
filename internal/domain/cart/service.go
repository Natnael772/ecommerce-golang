package cart

import (
	"context"
	"ecommerce-app/internal/pkg/errs"
	"ecommerce-app/internal/pkg/response"
	"ecommerce-app/pkg/pagination"
	"errors"

	"github.com/google/uuid"
)

type Service interface {
	CreateCart(ctx context.Context, req CreateCartRequest) (Cart, *errs.AppError)
	GetCartByID(ctx context.Context, id string) (Cart, *errs.AppError)
	GetCartByUserID(ctx context.Context, userID string) (Cart, *errs.AppError)
	GetCarts(ctx context.Context, page, perPage int) (CartsWithMeta, *errs.AppError)
	UpdateCart(ctx context.Context, id string, req UpdateCartRequest) (Cart, *errs.AppError)
	DeleteCart(ctx context.Context, id string) *errs.AppError
	DeleteExpiredCarts(ctx context.Context) *errs.AppError
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}


func (s *service) CreateCart(ctx context.Context, req CreateCartRequest) (Cart, *errs.AppError) {
	// Check if the user already has a cart
	existing, err := s.repo.GetByUserID(ctx, req.UserID)
	if err == nil && existing.ID != uuid.Nil {
		return existing, nil
	}

	// If it's not found, create a new cart
	if errors.Is(err, errs.ErrNotFound) {
		created, createErr := s.repo.Create(ctx, req)
		if createErr != nil {
			return Cart{}, errs.ErrInternal.WithMessage("Failed to create cart")
		}
		return created, nil
	}

	// Any other error (e.g., DB issue)
	if err != nil {
		return Cart{}, errs.ErrInternal.WithMessage("Failed to check existing cart")
	}

	return Cart{}, nil
}


// --- Get By ID ---
func (s *service) GetCartByID(ctx context.Context, id string) (Cart, *errs.AppError) {
	cart, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Cart{}, errs.ErrInternal.WithMessage("Failed to get cart by ID")
	}
	return cart, nil
}

// --- Get By User ID ---
func (s *service) GetCartByUserID(ctx context.Context, userID string) (Cart, *errs.AppError) {
	cart, err := s.repo.GetByUserID(ctx, userID)
	 if err != nil {
        if errors.Is(err, errs.ErrNotFound) {
            return Cart{}, errs.ErrNotFound
        }
		
        return Cart{}, errs.ErrInternal.WithMessage("Failed to get cart by user ID")
    }
	
	return cart, nil
}

// --- Paginated list of all carts ---
func (s *service) GetCarts(ctx context.Context, page, perPage int) (CartsWithMeta, *errs.AppError) {
	p := pagination.New(page, perPage)

	limit := int32(p.PerPage)
	offset := int32(p.Offset())

	carts, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return CartsWithMeta{}, errs.ErrInternal.WithMessage("Failed to get carts")
	}

	// In a real app, you might implement repo.CountCarts() for accurate total
	meta := response.Meta{
		Page:    p.Page,
		PerPage: p.PerPage,
		Total:   len(carts),
	}

	return CartsWithMeta{
		Carts: carts,
		Meta:  meta,
	}, nil
}

// --- Update (e.g., extend expiration) ---
func (s *service) UpdateCart(ctx context.Context, id string, req UpdateCartRequest) (Cart, *errs.AppError) {
	updated, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return Cart{}, errs.ErrInternal.WithMessage("Failed to update cart")
	}
	return updated, nil
}

// --- Delete ---
func (s *service) DeleteCart(ctx context.Context, id string) *errs.AppError {
	if err := s.repo.Delete(ctx, id); err != nil {
		return errs.ErrInternal.WithMessage("Failed to delete cart")
	}
	return nil
}

// --- Cleanup expired carts ---
func (s *service) DeleteExpiredCarts(ctx context.Context) *errs.AppError {
	if err := s.repo.DeleteExpired(ctx); err != nil {
		return errs.ErrInternal.WithMessage("Failed to delete expired carts")
	}
	return nil
}
