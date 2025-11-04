// handler.go
package order

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
	return &Handler{svc: svc}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	req := validator.GetValidatedBody[CreateOrderRequest](r)
	userID := r.Context().Value(middleware.UserIDKey).(string)

	res, appErr := h.svc.CreateOrder(r.Context(), userID, req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Created(w, res, "Order created successfully")
}

func (h *Handler) GetOrdersByUser(w http.ResponseWriter, r *http.Request) {
	page, perPage := pagination.GetPaginationParams(r)
	userID := r.Context().Value(middleware.UserIDKey).(string)

	result, appErr := h.svc.GetOrdersByUserID(r.Context(), userID, page, perPage)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OkWithMeta(w, result.Orders, result.Meta)
}

func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	order, appErr := h.svc.GetOrderByID(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, order, "Order fetched successfully")
}

func (h *Handler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	page, perPage := pagination.GetPaginationParams(r)
	
	result, appErr := h.svc.GetAllOrders(r.Context(), page, perPage)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OkWithMeta(w, result.Orders, result.Meta)
}

func (h *Handler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req := validator.GetValidatedBody[UpdateOrderStatusRequest](r)

	order, appErr := h.svc.UpdateOrderStatus(r.Context(), id, req.Status)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, order, "Order status updated successfully")
}

func (h *Handler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	appErr := h.svc.DeleteOrder(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.NoContent(w)
}
