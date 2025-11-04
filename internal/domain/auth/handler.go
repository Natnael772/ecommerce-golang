package auth

import (
	"ecommerce-app/internal/pkg/response"
	"ecommerce-app/internal/pkg/validator"
	"net/http"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	req := validator.GetValidatedBody[LoginRequest](r)
	
	resp, appErr := h.svc.Login(r.Context(), req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, resp, "Login successful")
}
