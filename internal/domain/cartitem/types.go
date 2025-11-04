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




// Dependency Injection Interface
type CartProvider interface {
    GetCartByUserID(ctx context.Context, userID string) (cart.Cart, *errs.AppError)
    CreateCart(ctx context.Context, req cart.CreateCartRequest) (cart.Cart, *errs.AppError)
}