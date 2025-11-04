package category

import (
	"ecommerce-app/internal/pkg/middleware"
	"ecommerce-app/internal/pkg/validator"

	"github.com/go-chi/chi/v5"
)

func Routes(svc Service) chi.Router {
	h := NewHandler(svc)
	r := chi.NewRouter()

	r.With(validator.Validate[CreateCategoryRequest]()).With(middleware.RoleMiddleware("admin")).Post("/", h.CreateCategory)

	r.Get("/", h.ListCategories)
	r.Get("/{id}", h.GetCategory)

	r.With(validator.Validate[UpdateCategoryRequest]()).With(middleware.RoleMiddleware("admin")).Put("/{id}", h.UpdateCategory)
	r.With(middleware.RoleMiddleware("admin")).Delete("/{id}", h.DeleteCategory)

	return r
}