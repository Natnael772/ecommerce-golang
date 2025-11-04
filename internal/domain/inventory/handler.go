package inventory

import (
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

func (h *Handler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	req := validator.GetValidatedBody[CreateInventoryRequest](r)

	inv, appErr := h.svc.CreateInventory(r.Context(), req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Created(w, inv)
}

func (h *Handler) GetInventoryByProductID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	inv, appErr := h.svc.GetInventoryByProductID(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, inv, "Inventory fetched successfully")
}

func (h *Handler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req := validator.GetValidatedBody[UpdateInventoryRequest](r)

	inv, appErr := h.svc.UpdateInventory(r.Context(), id, *req.Stock)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, inv, "Inventory updated successfully")
}

func (h *Handler) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	appErr := h.svc.DeleteInventory(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.NoContent(w)
}