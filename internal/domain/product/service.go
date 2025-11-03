package product

import (
	"context"
	"database/sql"
	"ecommerce-app/internal/pkg/errs"
	"ecommerce-app/internal/pkg/logger"
	"ecommerce-app/internal/pkg/response"
	"ecommerce-app/pkg/pagination"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (Product, *errs.AppError)
	GetProductByID(ctx context.Context, id string) (Product, *errs.AppError)
	ListProducts(ctx context.Context, page, perPage int) (ProductsWithMeta, *errs.AppError)
	UpdatePrice(ctx context.Context, id string, price int32) (Product, *errs.AppError)
	UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (Product, *errs.AppError)
	DeleteProduct(ctx context.Context, id string) *errs.AppError
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateProduct(ctx context.Context,req CreateProductRequest) (Product, *errs.AppError) {
	existingProduct, err := s.repo.GetBySku(ctx, req.SKU)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return Product{}, errs.ErrInternal.WithMessage("Failed to check existing product")
	}


	if err == nil && existingProduct.ID != uuid.Nil {
		return Product{}, errs.ErrConflict.WithMessage("Product with the same SKU already exists")
	}

	var discountValidUntil *time.Time
	if req.DiscountValidUntil != nil {
		parsed, perr := time.Parse(time.RFC3339, *req.DiscountValidUntil)
		if perr != nil {
			return Product{}, errs.ErrBadRequest.WithMessage("Invalid discountValidUntil format, expected RFC3339")
		}
		discountValidUntil = &parsed
	}	

	product := Product{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		PriceCents:       req.PriceCents,
		SKU:         req.SKU,
		Currency:       req.Currency,
		Attributes: req.Attributes,
		MainImageUrl: req.MainImageUrl,
		Images: 	req.Images,
		DiscountPercent: req.DiscountPercent,
		DiscountValidUntil: discountValidUntil,
	}

	createdProduct, err := s.repo.Create(ctx, product)
	
	if err != nil {
		logger.Info("Error creating product: %v", err)
		return Product{}, errs.ErrInternal.WithMessage("Failed to create product")
	}
	return createdProduct, nil
}

func (s *service) GetProductByID(ctx context.Context, id string) (Product, *errs.AppError) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if product.ID == uuid.Nil {
			return Product{}, errs.ErrNotFound.WithMessage("Product not found")
		}

		return Product{}, errs.ErrInternal.WithMessage("Failed to get product")
	}

	return product, nil
}


func (s *service) ListProducts(ctx context.Context, page, perPage int) (ProductsWithMeta, *errs.AppError) {
	p:= pagination.New(page, perPage)

	limit:=int32(p.PerPage)
	offset:= int32(p.Offset())

	products,err:= s.repo.List(ctx, limit, offset)
	if err != nil {
		return ProductsWithMeta{}, errs.ErrInternal.WithMessage("Failed to list products")
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return ProductsWithMeta{}, errs.ErrInternal.WithMessage("Failed to count products")
	}

	meta := response.Meta{
		Page:        p.Page,
		PerPage: p.PerPage,
		Total:       int(total),
	}

	result := ProductsWithMeta{
		Products: products,
		Meta:  meta,
	}

	return result, nil
}


func (s *service) UpdatePrice(ctx context.Context, id string, price int32) (Product, *errs.AppError) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Product{}, errs.ErrInternal.WithMessage("Failed to get product")
	}

	if product.ID == uuid.Nil {
		return Product{}, errs.ErrNotFound.WithMessage("Product not found")
	}

	updatedProduct, err := s.repo.UpdatePrice(ctx, id, price)
	if err != nil {
		return Product{}, errs.ErrInternal.WithMessage("Failed to update product price")
	}
	return updatedProduct, nil
}

func (s *service) UpdateProduct(ctx context.Context,id string, req UpdateProductRequest) (Product, *errs.AppError) {
	existingProduct, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Product{}, errs.ErrInternal.WithMessage("Failed to get product")
	}

	if existingProduct.ID == uuid.Nil {
		return Product{}, errs.ErrNotFound.WithMessage("Product not found")
	}

	updatedProduct, err := s.repo.UpdateProduct(ctx, id, req)
	if err != nil {
		return Product{}, errs.ErrInternal.WithMessage("Failed to update product")
	}
	return updatedProduct, nil
}

func (s *service) DeleteProduct(ctx context.Context, id string) *errs.AppError {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errs.ErrInternal.WithMessage("Failed to get product")
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return errs.ErrInternal.WithMessage("Failed to delete product")
	}
	return nil
}