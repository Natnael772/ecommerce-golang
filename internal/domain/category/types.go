package category

import (
	"ecommerce-app/internal/pkg/response"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoriesWithMeta struct {
	Categories []Category `json:"categories"`
	Meta       response.Meta       `json:"meta"`
}