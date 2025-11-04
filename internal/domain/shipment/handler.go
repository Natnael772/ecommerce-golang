package shipment

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
	return &Handler{svc}
}

func (h *Handler) CreateShipment(w http.ResponseWriter, r *http.Request) {
	req := validator.GetValidatedBody[CreateShipmentRequest](r)

	shipment, appErr := h.svc.CreateShipment(r.Context(), req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Created(w, shipment)
}

func (h *Handler) GetShipmentsByOrderID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	shipment, appErr := h.svc.GetShipmentsByOrderID(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, shipment, "Shipment retrieved successfully")
}

func (h *Handler) GetShipment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	shipment, appErr := h.svc.GetShipment(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, shipment)
}

func (h *Handler) UpdateShipmentStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req := validator.GetValidatedBody[UpdateShipmentStatusRequest](r)

	updatedShipment, appErr := h.svc.UpdateShipmentStatus(r.Context(), id, req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, updatedShipment, "Shipment status updated successfully")
}

func (h *Handler) DeleteShipment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	appErr := h.svc.DeleteShipment(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.NoContent(w)
}

