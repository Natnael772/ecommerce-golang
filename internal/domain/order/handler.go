// handler.go
package order

import (
	"ecommerce-app/internal/pkg/response"
	"ecommerce-app/internal/pkg/validator"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// userID := r.Context().Value("user_id").(string)
	userID := "0e451c10-a776-4bf1-be33-e5684d954dc3"
	req := validator.GetValidatedBody[CreateOrderRequest](r)

	res, appErr := h.svc.CreateOrder(r.Context(), userID, req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Created(w, res, "Order created successfully")
}

func (h *Handler) GetOrdersByUser(w http.ResponseWriter, r *http.Request) {
	// userID := r.Context().Value("user_id").(string)
	userID := "0e451c10-a776-4bf1-be33-e5684d954dc3"
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

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
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	
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
