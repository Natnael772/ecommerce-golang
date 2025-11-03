package user

import (
	"ecommerce-app/internal/pkg/response"
	"ecommerce-app/internal/pkg/validator"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	req := validator.GetValidatedBody[RegisterUserRequest](r)

	user, appErr := h.svc.Register(r.Context(), req)
	if appErr != nil {
	    response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Created(w, user) 
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

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	u, err := h.svc.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(u)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	page,_ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage,_ :=strconv.Atoi( r.URL.Query().Get("per_page"))

	result, appErr := h.svc.ListUsers(r.Context(), int32(page), int32(perPage))
	
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OkWithMeta(w, result.Users, result.Meta)
}
