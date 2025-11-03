package product

import (
	"ecommerce-app/internal/pkg/validator"

	"github.com/go-chi/chi/v5"
)

func Routes(svc Service) chi.Router {
	h := NewHandler(svc)
	r := chi.NewRouter()
	
	r.With(validator.Validate[CreateProductRequest]()).Post("/", h.CreateProduct)

	r.Get("/", h.GetProducts)
	r.Get("/{id}", h.GetProduct)
	
	r.With(validator.Validate[UpdatePriceRequest]()).Patch("/{id}/price", h.UpdatePrice)
	r.With(validator.Validate[UpdateProductRequest]()).Put("/{id}", h.UpdateProduct)

	return r
}