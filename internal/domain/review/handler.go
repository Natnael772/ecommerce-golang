package review

import (
	"ecommerce-app/internal/pkg/middleware"
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

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	req := validator.GetValidatedBody[CreateReviewRequest](r)

	category, appErr := h.svc.CreateReview(r.Context(), userID, req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Created(w, category)
}

func (h *Handler) GetReviewsByProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	page, perPage := pagination.GetPaginationParams(r)

	result, appErr := h.svc.GetReviewsByProduct(r.Context(),id, page, perPage)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OkWithMeta(w, result.Reviews, result.Meta)
}


func (h *Handler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req := validator.GetValidatedBody[UpdateReviewRequest](r)

	category, appErr := h.svc.UpdateReview(r.Context(), id, req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, category, "Review updated successfully")
}

func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	appErr := h.svc.DeleteReview(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Deleted(w)
}