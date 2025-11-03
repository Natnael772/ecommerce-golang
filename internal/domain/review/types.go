package review

import (
	"ecommerce-app/internal/pkg/response"
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	UserID    uuid.UUID `json:"user_id"`
	Rating    int32     `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReviewsWithMeta struct {
	Reviews     []Review `json:"reviews"`
	Meta       response.Meta       `json:"meta"`
}