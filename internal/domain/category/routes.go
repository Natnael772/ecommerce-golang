package category

import (
	"ecommerce-app/internal/pkg/validator"

	"github.com/go-chi/chi/v5"
)

func Routes(svc Service) chi.Router {
	h := NewHandler(svc)
	r := chi.NewRouter()

	r.With(validator.Validate[CreateCategoryRequest]()).Post("/", h.CreateCategory)

	r.Get("/", h.ListCategories)
	r.Get("/{id}", h.GetCategory)

	r.With(validator.Validate[UpdateCategoryRequest]()).Put("/{id}", h.UpdateCategory)
	r.Delete("/{id}", h.DeleteCategory)

	return r
}