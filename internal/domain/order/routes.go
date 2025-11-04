package order

import (
	"ecommerce-app/internal/pkg/middleware"
	"ecommerce-app/internal/pkg/validator"

	"github.com/go-chi/chi/v5"
)

func Routes(svc Service) chi.Router {
	h := NewHandler(svc)
	r := chi.NewRouter()

	r.With(validator.Validate[CreateOrderRequest]()).With(middleware.RoleMiddleware("customer")).Post("/", h.CreateOrder)
	
	r.With(middleware.RoleMiddleware("customer")).Get("/", h.GetOrdersByUser)
	r.Get("/{id}", h.GetOrderByID)
	r.With(middleware.RoleMiddleware("customer")).Get("/all", h.GetAllOrders)

	r.With(validator.Validate[UpdateOrderStatusRequest]()).With(middleware.RoleMiddleware("admin")).Put("/{id}", h.UpdateOrderStatus)
	r.With(middleware.RoleMiddleware("admin")).Delete("/{id}", h.DeleteOrder)

	return r
}
