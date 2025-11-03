package coupon

import (
	"ecommerce-app/internal/pkg/validator"

	"github.com/go-chi/chi/v5"
)

func Routes (svc Service) chi.Router {
	h := NewHandler(svc)
	r := chi.NewRouter()

	r.With(validator.Validate[CreateCouponRequest]()).Post("/", h.CreateCoupon)

	r.Get("/", h.GetCoupons)
	r.Get("/{id}", h.GetCoupon)
	r.Post("/{code}/increment-usage", h.IncrementCouponUsage)

	r.With(validator.Validate[UpdateCouponRequest]()).Put("/{id}", h.UpdateCoupon)

	r.Delete("/{id}", h.DeleteCoupon)

	return r
}