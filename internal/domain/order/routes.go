package order

import (
	"ecommerce-app/internal/pkg/validator"

	"github.com/go-chi/chi/v5"
)

func Routes(svc Service) chi.Router {
	h := NewHandler(svc)
	r := chi.NewRouter()

	r.With(validator.Validate[CreateOrderRequest]()).Post("/", h.CreateOrder)
	
	r.Get("/", h.GetOrdersByUser)
	r.Get("/{id}", h.GetOrderByID)
	r.Get("/all", h.GetAllOrders)

	r.With(validator.Validate[UpdateOrderStatusRequest]()).Put("/{id}", h.UpdateOrderStatus)
	r.Delete("/{id}", h.DeleteOrder)

	return r
}
