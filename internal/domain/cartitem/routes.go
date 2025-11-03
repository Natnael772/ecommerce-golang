package cartitem

import (
	"ecommerce-app/internal/pkg/validator"

	"github.com/go-chi/chi/v5"
)

func Routes(svc Service) chi.Router {
	h := NewHandler(svc)
	r := chi.NewRouter()

	r.With(validator.Validate[AddItemRequest]()).Post("/", h.AddItem)

	r.With(validator.Validate[[]AddItemRequest]()).Post("/batch", h.AddItems)

	r.Get("/user/items", h.GetUserCartItems)

	r.Get("/{id}", h.GetItem)
	r.Get("/carts/{cartID}/items", h.ListItemsByCart)

	r.Delete("/{id}", h.DeleteItem)

	r.Delete("/clear", h.ClearCartItems)

	return r
}