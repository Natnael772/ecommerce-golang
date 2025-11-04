package cartitem

import (
	"ecommerce-app/internal/pkg/middleware"
	"ecommerce-app/internal/pkg/response"
	"ecommerce-app/internal/pkg/validator"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	req := validator.GetValidatedBody[AddItemRequest](r)
	userID := r.Context().Value(middleware.UserIDKey).(string)
	

	item, appErr := h.svc.AddItem(r.Context(),userID, req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Created(w, item)
}

func (h *Handler) AddItems(w http.ResponseWriter, r *http.Request) {
	req := validator.GetValidatedBody[[]AddItemRequest](r)
	userID := r.Context().Value(middleware.UserIDKey).(string)

	items, appErr := h.svc.AddItems(r.Context(),userID, req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Created(w, items)
}

func (h *Handler) GetUserCartItems(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	items, appErr := h.svc.GetItemsByUserID(r.Context(), userID)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, items, "User cart items retrieved successfully")
}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	item, appErr := h.svc.GetItemByID(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, item, "Cart item retrieved successfully")
}

func (h *Handler) ListItemsByCart(w http.ResponseWriter, r *http.Request) {
	cartID := chi.URLParam(r, "cartID")

	items, appErr := h.svc.ListItemsByCart(r.Context(), cartID)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, items, "Cart items retrieved successfully")
}

func (h *Handler) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req := validator.GetValidatedBody[UpdateQuantityRequest](r)

	item, appErr := h.svc.UpdateItemQuantity(r.Context(), id, req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, item, "Cart item quantity updated successfully")
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	appErr := h.svc.DeleteItem(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.NoContent(w)
}

func (h *Handler) DeleteItemsByCart(w http.ResponseWriter, r *http.Request) {
	cartID := chi.URLParam(r, "cartID")

	appErr := h.svc.DeleteItem(r.Context(), cartID)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, interface{}(nil), "All cart items deleted successfully")
}


func (h *Handler) ClearCartItems(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	appErr := h.svc.ClearCartItems(r.Context(), userID)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.NoContent(w)
}