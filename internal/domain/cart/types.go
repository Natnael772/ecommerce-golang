package cart

import (
	"ecommerce-app/internal/pkg/response"
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID     `json:"id"`
	UserID    uuid.UUID      `json:"user_id"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type CartItem struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Price     float64 `json:"price"`
}

// Pagination
type CartsWithMeta struct {
	Carts []Cart       `json:"carts"`
	Meta  response.Meta `json:"meta"`
}


