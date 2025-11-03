package category

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
	return &Handler{svc}
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	req := validator.GetValidatedBody[CreateCategoryRequest](r)

	category, appErr := h.svc.CreateCategory(r.Context(), req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Created(w, category)
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	page,_ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage,_ :=strconv.Atoi( r.URL.Query().Get("per_page"))

	result, appErr := h.svc.ListCategories(r.Context(), page, perPage)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OkWithMeta(w, result.Categories, result.Meta)
}

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	category, appErr := h.svc.GetCategory(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, category, "Category fetched successfully")
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req := validator.GetValidatedBody[UpdateCategoryRequest](r)

	category, appErr := h.svc.UpdateCategory(r.Context(), id, req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, category, "Category updated successfully")
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	appErr := h.svc.DeleteCategory(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Deleted(w)
}