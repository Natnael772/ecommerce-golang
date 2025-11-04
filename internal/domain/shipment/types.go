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
