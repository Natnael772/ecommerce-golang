package cart

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
	return &Handler{svc: svc}
}

// --- POST /carts ---
func (h *Handler) CreateCart(w http.ResponseWriter, r *http.Request) {
	req := validator.GetValidatedBody[CreateCartRequest](r)

	cart, appErr := h.svc.CreateCart(r.Context(), req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Created(w, cart)
}

// --- GET /carts ---
func (h *Handler) GetCarts(w http.ResponseWriter, r *http.Request) {
	page, perPage := pagination.GetPaginationParams(r)

	result, appErr := h.svc.GetCarts(r.Context(), page, perPage)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OkWithMeta(w, result.Carts, result.Meta)
}

// --- GET /carts/{id} ---
func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cart, appErr := h.svc.GetCartByID(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, cart, "Cart retrieved successfully")
}

// --- GET /users/{userID}/cart ---
func (h *Handler) GetCartByUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	cart, appErr := h.svc.GetCartByUserID(r.Context(), userID)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, cart, "User cart retrieved successfully")
}

// --- PUT /carts/{id} ---
func (h *Handler) UpdateCart(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req := validator.GetValidatedBody[UpdateCartRequest](r)

	updatedCart, appErr := h.svc.UpdateCart(r.Context(), id, req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, updatedCart, "Cart updated successfully")
}

// --- DELETE /carts/{id} ---
func (h *Handler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	appErr := h.svc.DeleteCart(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.NoContent(w)
}

// --- DELETE /carts/expired ---
func (h *Handler) DeleteExpiredCarts(w http.ResponseWriter, r *http.Request) {
	appErr := h.svc.DeleteExpiredCarts(r.Context())
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, interface{}(nil), "Expired carts deleted successfully")
}
