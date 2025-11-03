package review

import (
	"context"
	"ecommerce-app/internal/pkg/errs"
	"ecommerce-app/internal/pkg/logger"
	"ecommerce-app/internal/pkg/response"
	"ecommerce-app/pkg/pagination"
	"errors"

	"github.com/google/uuid"
)

type Service interface {
	CreateReview(ctx context.Context, userID string, req CreateReviewRequest) (Review, *errs.AppError)
	GetReview(ctx context.Context, id string) (Review, *errs.AppError)
	GetReviewsByProduct(ctx context.Context, productID string, page, perPage int) (ReviewsWithMeta, *errs.AppError)
	UpdateReview(ctx context.Context, id string, req UpdateReviewRequest) (Review, *errs.AppError)
	DeleteReview(ctx context.Context, id string) *errs.AppError
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateReview(ctx context.Context, userID string, req CreateReviewRequest) (Review, *errs.AppError) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return Review{}, errs.ErrBadRequest.WithMessage("Invalid user id")
	}

	prodExists, err := s.repo.CheckProductExists(ctx, req.ProductID.String())
	if err != nil {
		return Review{}, errs.ErrInternal.WithMessage("Failed to check product existence")
	}

	if !prodExists {
		return Review{}, errs.ErrBadRequest.WithMessage("Product does not exist")
	}

	reviewExists, err := s.repo.GetUserReviewForProduct(ctx, userID, req.ProductID.String())
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		logger.Error("Error checking existing review: %v", err)
		return Review{}, errs.ErrInternal.WithMessage("Failed to check existing review")
	}

	if reviewExists.ID != uuid.Nil {
		return Review{}, errs.ErrConflict.WithMessage("User has already reviewed this product")
	}

	review := Review{
		Rating:    req.Rating,
		Comment:   req.Comment,
		ProductID: req.ProductID,
		UserID:    parsedUserID,
	}

	createdRev, err := s.repo.Create(ctx, review)
	if err != nil {
		return Review{}, errs.ErrInternal.WithMessage("Failed to create review")
	}

	return createdRev, nil
}

func (s *service) GetReview(ctx context.Context, id string) (Review, *errs.AppError) {
	review, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Review{}, errs.ErrInternal.WithMessage("Failed to get review")
	}
	return review, nil
}

func (s *service) UpdateReview(ctx context.Context, id string, req UpdateReviewRequest) (Review, *errs.AppError) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return Review{}, errs.ErrBadRequest.WithMessage("Invalid review id")
	}

	review := Review{
		ID:       parsedID,
	}	

	if req.Rating != nil {
		review.Rating = *req.Rating
	}
	
	if req.Comment != nil {
		review.Comment = *req.Comment
	}

	updatedRev, err := s.repo.Update(ctx, review)
	if err != nil {
		return Review{}, errs.ErrInternal.WithMessage("Failed to update review")
	}

	return updatedRev, nil
}

func (s *service) GetReviewsByProduct(ctx context.Context, productID string, page, perPage int) (ReviewsWithMeta, *errs.AppError) {
	p:= pagination.New(page, perPage)

	limit:=int32(p.PerPage)
	offset:= int32(p.Offset())

	reviews,  err := s.repo.ListByProduct(ctx, productID, limit, offset)
	if err != nil {
		return ReviewsWithMeta{}, errs.ErrInternal.WithMessage("Failed to list Reviews")
	}

	total, err := s.repo.CountByProduct(ctx, productID)
	if err != nil {
		return ReviewsWithMeta{}, errs.ErrInternal.WithMessage("Failed to count Reviews")
	}

	meta:= response.Meta{
		Page:       int(p.Page),
		PerPage:    int(p.PerPage),
		Total:      int(total),
	}

	result := ReviewsWithMeta{
		Reviews: reviews,
		Meta:       meta,
	}

	return result, nil
}

func (s *service) DeleteReview(ctx context.Context, id string) *errs.AppError {

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errs.ErrNotFound.WithMessage("Review not found")
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return errs.ErrInternal.WithMessage("Failed to delete review")
	}
	return nil
}

