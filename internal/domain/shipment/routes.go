package shipment

import (
	"ecommerce-app/internal/pkg/validator"

	"github.com/go-chi/chi/v5"
)

func Routes(svc Service) chi.Router {
	h := NewHandler(svc)
	r := chi.NewRouter()

	r.With(validator.Validate[CreateShipmentRequest]()).Post("/", h.CreateShipment)

	r.Get("/order/{orderID}", h.GetShipmentsByOrderID)

	r.With(validator.Validate[UpdateShipmentStatusRequest]()).Patch("/{id}/status", h.UpdateShipmentStatus)

	return r
}