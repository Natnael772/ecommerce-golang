package cart

import (
	"github.com/go-chi/chi/v5"
)

func Routes (svc Service) chi.Router {
	h := NewHandler(svc)
	r := chi.NewRouter()

	// r.With(validator.Validate[CreateCartRequest]()).Post("/", h.CreateCart)

	r.Get("/", h.GetCarts)
	r.Get("/{id}", h.GetCart)

	r.Get("/users/{userID}/cart", h.GetCartByUser)

	return r
}