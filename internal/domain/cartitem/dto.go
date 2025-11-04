package cartitem

import "github.com/google/uuid"

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