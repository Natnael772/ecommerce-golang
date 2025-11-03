package review

import (
	"ecommerce-app/internal/pkg/validator"

	"github.com/go-chi/chi/v5"
)

func Routes(svc Service) chi.Router {
	h := NewHandler(svc)
	r := chi.NewRouter()

	r.With(validator.Validate[CreateReviewRequest]()).Post("/", h.CreateReview)

	r.Get("/product/{id}", h.GetReviewsByProduct)

	r.With(validator.Validate[UpdateReviewRequest]()).Put("/{id}", h.UpdateReview)
	r.Delete("/{id}", h.DeleteReview)
	return r
}