package inventory

import (
	"ecommerce-app/internal/pkg/middleware"
	"ecommerce-app/internal/pkg/validator"

	"github.com/go-chi/chi/v5"
)

func Routes(svc Service) chi.Router {
	h := NewHandler(svc)
	r := chi.NewRouter()

	r.With(validator.Validate[CreateInventoryRequest]()).With(middleware.RoleMiddleware("admin")).Post("/", h.CreateInventory)

	r.With(middleware.RoleMiddleware("admin")).Get("/{id}", h.GetInventoryByProductID)

	r.With(validator.Validate[UpdateInventoryRequest]()).With(middleware.RoleMiddleware("admin")).Put("/{id}", h.UpdateInventory)

	r.With(middleware.RoleMiddleware("admin")).Delete("/{id}", h.DeleteInventory)

	return r
}