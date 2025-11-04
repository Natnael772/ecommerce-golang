package inventory

import "time"

type Inventory struct {
	ProductID string `json:"product_id"`
	Stock     int32  `json:"stock"`
	Reserved  int32  `json:"reserved"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

