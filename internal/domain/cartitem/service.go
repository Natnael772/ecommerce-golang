package cartitem

import (
	"context"
	"ecommerce-app/internal/domain/cart"
	"ecommerce-app/internal/pkg/errs"
	"ecommerce-app/internal/pkg/logger"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	AddItem(ctx context.Context,userId string, req AddItemRequest) (CartItem, *errs.AppError)
	AddItems (ctx context.Context,userId string, req []AddItemRequest) ([]CartItem, *errs.AppError)
	GetItemByID(ctx context.Context, id string) (CartItem, *errs.AppError)
	GetItemByCartAndProduct(ctx context.Context, cartID, productID string) (CartItem, *errs.AppError)
	GetItemsByUserID(ctx context.Context, userID string) ([]CartItem, *errs.AppError)
	ListItemsByCart(ctx context.Context, cartID string) ([]CartItem, *errs.AppError)
	UpdateItemQuantity(ctx context.Context, id string, req UpdateQuantityRequest) (CartItem, *errs.AppError)
	DeleteItem(ctx context.Context, id string) *errs.AppError
	ClearCartItems(ctx context.Context, cartID string) *errs.AppError
}

type service struct {
	repo Repository
	cartSvc CartProvider
}

func NewService(repo Repository, cartSvc CartProvider) Service {
	return &service{repo: repo, cartSvc: cartSvc}
}

// --- Add or Upsert Cart Item ---
func (s *service) AddItem(ctx context.Context, userId string, req AddItemRequest) (CartItem, *errs.AppError) {
	if _, err := uuid.Parse(req.ProductID); err != nil {
		return CartItem{}, errs.ErrBadRequest.WithMessage("Invalid product ID")
	}

	// Fetch existing cart or create a new one if necessary
	var c cart.Cart
	c, appErr := s.cartSvc.GetCartByUserID(ctx, userId)
	if appErr != nil {
		if errors.Is(appErr, errs.ErrNotFound) {
			logger.Info("No existing cart for user %s, creating a new one", userId)

			expiresAt := time.Now().Add(14 * 24 * time.Hour)

			newCart, createErr := s.cartSvc.CreateCart(ctx, cart.CreateCartRequest{UserID: userId, ExpiresAt: &expiresAt})
			if createErr != nil {
				logger.Error("Failed to create cart for user %s: %v", userId, createErr)
				return CartItem{}, errs.ErrInternal.WithMessage("Failed to create cart for user")
			}
			c = newCart
			logger.Info("Created new cart for user %s: %+v", userId, c)			

		} else {
			return CartItem{}, appErr
		}
	}

	item := CartItem{
		CartID:    c.ID,
		ProductID: uuid.MustParse(req.ProductID),
		Quantity:  req.Quantity,
	}

	created, err := s.repo.Add(ctx, item)
	if err != nil {
		logger.Error("Error adding or updating cart item: %v", err)
		return CartItem{}, errs.ErrInternal.WithMessage("Failed to add or update cart item")
	}

	return created, nil
}

// --- Add or Upsert Multiple Cart Items ---
func (s *service) AddItems(ctx context.Context, userId string, req []AddItemRequest) ([]CartItem, *errs.AppError) {
	// Fetch existing cart or create a new one if necessary
	var c cart.Cart
	c, appErr := s.cartSvc.GetCartByUserID(ctx, userId)
	if appErr != nil {
		if errors.Is(appErr, errs.ErrNotFound) {			
			expiresAt := time.Now().Add(14 * 24 * time.Hour)

			newCart, createErr := s.cartSvc.CreateCart(ctx, cart.CreateCartRequest{UserID: userId, ExpiresAt: &expiresAt})
			if createErr != nil {
				return nil, errs.ErrInternal.WithMessage("Failed to create cart for user")
			}
			c = newCart			

		} else {
			return nil, appErr
		}
	}

	cartItemsReq := AddItemsRequest{
		CartID: c.ID,
		Items: req,
	}

	createdItems, err := s.repo.AddItems(ctx, cartItemsReq)
	if err != nil {
		logger.Error("Error adding or updating cart items: %v", err)
		return nil, errs.ErrInternal.WithMessage("Failed to add or update cart items")
	}

	return createdItems, nil
}

// --- Get items by UserID ---
func (s *service) GetItemsByUserID(ctx context.Context, userID string) ([]CartItem, *errs.AppError) {
	items, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errs.ErrInternal.WithMessage("Failed to get cart items by user ID")
	}
	return items, nil
}

// --- Get item by ID ---
func (s *service) GetItemByID(ctx context.Context, id string) (CartItem, *errs.AppError) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return CartItem{}, errs.ErrInternal.WithMessage("Failed to get cart item by ID")
	}
	return item, nil
}

// --- Get item by CartID + ProductID ---
func (s *service) GetItemByCartAndProduct(ctx context.Context, cartID, productID string) (CartItem, *errs.AppError) {
	item, err := s.repo.GetByCartAndProduct(ctx, cartID, productID)
	if err != nil {
		return CartItem{}, errs.ErrInternal.WithMessage("Failed to get cart item by cart and product")
	}
	return item, nil
}

// --- List all items in a cart ---
func (s *service) ListItemsByCart(ctx context.Context, cartID string) ([]CartItem, *errs.AppError) {
	items, err := s.repo.ListByCart(ctx, cartID)
	if err != nil {
		return nil, errs.ErrInternal.WithMessage("Failed to list cart items")
	}
	return items, nil
}

// --- Update quantity ---
func (s *service) UpdateItemQuantity(ctx context.Context, id string, req UpdateQuantityRequest) (CartItem, *errs.AppError) {
	if req.Quantity <= 0 {
		return CartItem{}, errs.ErrBadRequest.WithMessage("Quantity must be greater than zero")
	}

	item, err := s.repo.UpdateQuantity(ctx, id, req.Quantity)
	if err != nil {
		return CartItem{}, errs.ErrInternal.WithMessage("Failed to update cart item quantity")
	}
	return item, nil
}

// --- Delete single item ---
func (s *service) DeleteItem(ctx context.Context, id string) *errs.AppError {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return errs.ErrInternal.WithMessage("Failed to delete cart item")
	}
	return nil
}

// --- Clear all items for a cart ---
func (s *service) ClearCartItems(ctx context.Context, cartID string) *errs.AppError {
	err := s.repo.DeleteByCart(ctx, cartID)
	if err != nil {
		return errs.ErrInternal.WithMessage("Failed to clear cart items")
	}
	return nil
}
