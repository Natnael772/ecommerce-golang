package category

import (
	"context"
	"ecommerce-app/internal/pkg/errs"
	"ecommerce-app/internal/pkg/response"
	"ecommerce-app/pkg/pagination"
	"ecommerce-app/pkg/slug"

	"github.com/google/uuid"
)

type Service interface {
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (Category, *errs.AppError)
	GetCategory(ctx context.Context, id string) (Category, *errs.AppError)
	ListCategories(ctx context.Context, page, perPage int) (CategoriesWithMeta, *errs.AppError)
	UpdateCategory(ctx context.Context, id string, req UpdateCategoryRequest) (Category, *errs.AppError)
	DeleteCategory(ctx context.Context, id string) *errs.AppError
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateCategory(ctx context.Context, req CreateCategoryRequest) (Category, *errs.AppError) {
	slugStr, err := slug.GenerateSlug(req.Name)
	if err != nil {
		return Category{}, errs.ErrInternal.WithMessage("Failed to generate slug")
	}

	catExists,err:= s.repo.GetBySlug(ctx, slugStr)
	if err == nil && catExists.ID != uuid.Nil {
		return Category{}, errs.ErrConflict.WithMessage("Category with the same name or slug already exists")
	}

	category := Category{
		Name:     req.Name,
		ParentID: req.ParentID,
		Slug:     slugStr,
	}

	createdCat, err := s.repo.Create(ctx, category)
	if err != nil {
		return Category{}, errs.ErrInternal.WithMessage("Failed to create category")
	}

	return createdCat, nil
}

func (s *service) GetCategory(ctx context.Context, id string) (Category, *errs.AppError) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Category{}, errs.ErrInternal.WithMessage("Failed to get category")
	}
	return category, nil
}

func (s *service) ListCategories(ctx context.Context, page, perPage int) (CategoriesWithMeta, *errs.AppError) {
	p:= pagination.New(page, perPage)

	limit:=int32(p.PerPage)
	offset:= int32(p.Offset())

	cats,  err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return CategoriesWithMeta{}, errs.ErrInternal.WithMessage("Failed to list categories")
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return CategoriesWithMeta{}, errs.ErrInternal.WithMessage("Failed to count categories")
	}

	meta:= response.Meta{
		Page:       int(p.Page),
		PerPage:    int(p.PerPage),
		Total:      int(total),
	}

	result := CategoriesWithMeta{
		Categories: cats,
		Meta:       meta,
	}

	return result, nil
}

func (s *service) UpdateCategory(ctx context.Context, id string, req UpdateCategoryRequest) (Category, *errs.AppError) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return Category{}, errs.ErrBadRequest.WithMessage("Invalid category id")
	}

	category := Category{
		ID:       parsedID,
		ParentID: req.ParentID,
	}	

	if req.Name != nil {
		category.Name = *req.Name
	}
	
	if req.Name != nil {
		slugStr, err := slug.GenerateSlug(*req.Name)
		if err != nil {
			return Category{}, errs.ErrInternal.WithMessage("Failed to generate slug")
		}
		category.Slug = slugStr
	}

	updatedCat, err := s.repo.Update(ctx, category)
	if err != nil {
		return Category{}, errs.ErrInternal.WithMessage("Failed to update category")
	}

	return updatedCat, nil
}

func (s *service) DeleteCategory(ctx context.Context, id string) *errs.AppError {

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errs.ErrNotFound.WithMessage("Category not found")
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return errs.ErrInternal.WithMessage("Failed to delete category")
	}
	return nil
}
