package product

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

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	req := validator.GetValidatedBody[CreateProductRequest](r)

	product, appErr := h.svc.CreateProduct(r.Context(), req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.Created(w, product)
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page,_ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage,_ :=strconv.Atoi( r.URL.Query().Get("per_page"))

	result, appErr := h.svc.ListProducts(r.Context(), page, perPage)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OkWithMeta(w, result.Products,result.Meta)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	product, appErr := h.svc.GetProductByID(r.Context(), id)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, product, "Product retrieved successfully")
}

func (h *Handler) UpdatePrice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req := validator.GetValidatedBody[UpdatePriceRequest](r)

	updatedProduct, appErr := h.svc.UpdatePrice(r.Context(), id, req.PriceCents)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, updatedProduct, "Product price updated successfully")
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req := validator.GetValidatedBody[UpdateProductRequest](r)

	updatedProduct, appErr := h.svc.UpdateProduct(r.Context(), id, req)
	if appErr != nil {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}

	response.OK(w, updatedProduct, "Product updated successfully")
}
