package shipment

import "time"

type Shipment struct {
	ID            string
	OrderID       string
	Carrier       string
	TrackingNumber string
	Status        string
	ShippedAt     *time.Time
	DeliveredAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}


// Request DTOs

type CreateShipmentRequest struct {
	OrderID        string     `json:"order_id"`
	Carrier        string     `json:"carrier"`
	TrackingNumber string     `json:"tracking_number"`
	Status         string     `json:"status"`
	ShippedAt      *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
}


type UpdateShipmentStatusRequest struct {
	Status      string     `json:"status"`
	ShippedAt   *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
}