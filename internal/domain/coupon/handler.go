package coupon

import (
	"ecommerce-app/internal/pkg/response"
	"ecommerce-app/internal/pkg/validator"
	"ecommerce-app/pkg/pagination"
	"net/http"

	"github.com/go-chi/chi/v5"
)



type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	req := validator.GetValidatedBody[CreateCouponRequest](r)

	coupon, appErr := h.svc.CreateCoupon(r.Context(), req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Created(w, coupon)
}

func (h *Handler) GetCoupons(w http.ResponseWriter, r *http.Request) {
	page, perPage := pagination.GetPaginationParams(r)

	result, appErr := h.svc.GetCoupons(r.Context(), page, perPage)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OkWithMeta(w, result.Coupons, result.Meta)
}

func (h *Handler) GetCoupon(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	coupon, appErr := h.svc.GetCouponByID(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, coupon, "Coupon retrieved successfully")
}

func (h *Handler) IncrementCouponUsage(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	updatedCoupon, appErr := h.svc.IncrementCouponUsage(r.Context(), code)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, updatedCoupon, "Coupon usage incremented successfully")
}

func (h *Handler) UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req := validator.GetValidatedBody[UpdateCouponRequest](r)

	updatedCoupon, appErr := h.svc.UpdateCoupon(r.Context(), id, req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, updatedCoupon, "Coupon updated successfully")
}

func (h *Handler) DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	appErr := h.svc.DeleteCoupon(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.NoContent(w)
}