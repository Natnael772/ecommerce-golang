package cartitem

import (
	"context"
	"ecommerce-app/internal/domain/cart"
	"ecommerce-app/internal/pkg/errs"
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID        uuid.UUID `json:"id"`
	CartID    uuid.UUID `json:"cart_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}


// Request DTOs

type AddItemRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid4"`
	Quantity  int32  `json:"quantity" validate:"required,gte=1"`
}


type AddItemsRequest struct {
	CartID     uuid.UUID   `json:"cart_id" validate:"uuid4"`
	Items      []AddItemRequest `json:"items" validate:"required,dive"`
}

type UpdateItemRequest struct {
	Quantity *int32 `json:"quantity,omitempty" validate:"omitempty,gte=1"`
}

type UpdateQuantityRequest struct {
	CartID     string `json:"cart_id" validate:"required,uuid4"`
	CartItemID string `json:"cart_item_id" validate:"required,uuid4"`
	Quantity   int32  `json:"quantity" validate:"required,gte=1"`
}

type RemoveCartItemRequest struct {
	CartID     string `json:"cart_id" validate:"required,uuid4"`
	CartItemID string `json:"cart_item_id" validate:"required,uuid4"`
}

// Dependency Injection Interface
type CartProvider interface {
    GetCartByUserID(ctx context.Context, userID string) (cart.Cart, *errs.AppError)
    CreateCart(ctx context.Context, req cart.CreateCartRequest) (cart.Cart, *errs.AppError)
}