package address

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
	return &Handler{svc}
}

func (h *Handler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	req := validator.GetValidatedBody[CreateAddressRequest](r)

	address, appErr := h.svc.CreateAddress(r.Context(), userID, req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Created(w, address)
}

func (h *Handler) getAddressesByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	address, appErr := h.svc.GetAddressesByUserID(r.Context(), userID)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, address, "User Address fetched successfully")
}

func (h *Handler) GetAddressByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	address, appErr := h.svc.GetAddressByID(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, address, "Address fetched successfully")
}

func (h *Handler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req := validator.GetValidatedBody[UpdateAddressRequest](r)	

	address, appErr := h.svc.UpdateAddress(r.Context(), id, req)
	if appErr !=nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w,address, "Address updated successfully")
}

func (h *Handler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	appErr := h.svc.DeleteAddress(r.Context(), id)
	if appErr !=nil {
		response.Error(w, appErr.Code, appErr.Message)
		return	
	}

	response.Deleted(w)
}	