package review

import (
	"context"
	"database/sql"
	"ecommerce-app/internal/pkg/database/sqlc"
	"ecommerce-app/internal/pkg/errs"
	"ecommerce-app/internal/pkg/logger"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface{
	Create(ctx context.Context, r Review) (Review, error)
	GetByID(ctx context.Context, id string) (Review, error)
	GetUserReviews(ctx context.Context, userID string) ([]Review, error)
	GetUserReviewForProduct(ctx context.Context, userID, productID string) (Review, error)
	ListByProduct(ctx context.Context, productID string, limit, offset int32) ([]Review, error)
	Update(ctx context.Context, r Review) (Review, error)
	Delete(ctx context.Context, id string) error
	CountByProduct(ctx context.Context, productID string) (int32, error)
	CheckProductExists(ctx context.Context, productID string) (bool, error)
}

type repository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repository{q}
}

func (r *repository) Create(ctx context.Context, rev Review) (Review, error) {
	params:= sqlc.CreateReviewParams{
		Rating: rev.Rating,
		Comment: pgtype.Text{String: rev.Comment, Valid: true},
		ProductID: pgtype.UUID{Bytes: rev.ProductID, Valid: true},
		UserID: pgtype.UUID{Bytes: rev.UserID, Valid: true},
	}

	row, err := r.q.CreateReview(ctx, params)
	if err != nil {
		logger.Error("Error creating review: %v", err)
		return Review{}, err
	}

	return mapReview(row), nil
}

func (r *repository) GetByID(ctx context.Context, id string) (Review, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return Review{}, err
	}

	row, err := r.q.GetReview(ctx, uuid)
	if err != nil {
		logger.Error("Error getting review by ID: %v", err)
		return Review{}, err
	}

	return mapReview(row), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return err
	}

	err := r.q.DeleteReview(ctx, uuid)
	if err != nil {
		logger.Error("Error deleting review: %v", err)
		return err
	}

	return nil
}

func (r *repository) GetUserReviews(ctx context.Context, userID string) ([]Review, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(userID); err != nil {
		return nil, err
	}

	rows, err := r.q.GetReviewsByUser(ctx, uuid)
	if err != nil {
		logger.Error("Error getting reviews by user ID: %v", err)
		return nil, err
	}

	out := make([]Review, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapReview(r))
	}
	return out, nil
}	
		

func (r *repository) GetUserReviewForProduct(ctx context.Context, userID, productID string) (Review, error) {
	var userUUID, productUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return Review{}, err
	}
	if err := productUUID.Scan(productID); err != nil {
		return Review{}, err
	}

	row, err := r.q.GetReviewByUserAndProduct(ctx, sqlc.GetReviewByUserAndProductParams{
		UserID: userUUID,
		ProductID: productUUID,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Review{}, errs.ErrNotFound
		}
		return Review{}, err
	}

	if !row.ID.Valid {
		return Review{}, nil
	}

	return mapReview(row), nil
}


func (r *repository) Count(ctx context.Context, productID string) (int32, error) {

	var prodId pgtype.UUID
	if err := prodId.Scan(productID); err != nil {
		return 0,  err
	}

	count, err := r.q.CountReviewsByProduct(ctx, prodId)
	if err != nil {
		logger.Error("Error counting reviews: %v", err)
		return 0, err
	}
	return int32(count), nil
}

func (r *repository) ListByProduct(ctx context.Context, productID string, limit, offset int32) ([]Review, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(productID); err != nil {
		return nil,  err
	}

	params:= sqlc.ListReviewsByProductParams{
		ProductID: uuid,
		Limit: limit,
		Offset: offset,
	}

	rows, err := r.q.ListReviewsByProduct(ctx, params)
	if err != nil {
		logger.Error("Error listing reviews by product: %v", err)
		return nil, err
	}

	out := make([]Review, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapReview(r))
	}
	return out, nil
}

	

func (r *repository) CountByProduct(ctx context.Context, productID string) (int32, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(productID); err != nil {
		return 0,  err
	}

	count, err := r.q.CountReviewsByProduct(ctx, uuid)
	if err != nil {
		logger.Error("Error counting reviews by product: %v", err)
		return 0, err
	}
	return int32(count), nil
}

func (r *repository) Update(ctx context.Context, rev Review) (Review, error) {
	params:= sqlc.UpdateReviewParams{
		ID: pgtype.UUID{Bytes: rev.ID, Valid: true},
		Rating: rev.Rating,
		Comment: pgtype.Text{String: rev.Comment, Valid: true},
	}

	row, err := r.q.UpdateReview(ctx, params)
	if err != nil {
		logger.Error("Error updating review: %v", err)
		return Review{}, err
	}

	return mapReview(row), nil
}

func (r *repository) CheckProductExists(ctx context.Context, productID string) (bool, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(productID); err != nil {
		return false,  err
	}

	prod, err := r.q.GetProductByID(ctx, uuid)
	if err != nil {
		logger.Error("Error checking product existence: %v", err)
		return false, err
	}		

	return prod.ID.Valid, nil
}



func mapReview(row sqlc.Review) Review {
	return Review{
		ID:        uuid.UUID(row.ID.Bytes),
		ProductID: uuid.UUID(row.ProductID.Bytes),
		UserID:    uuid.UUID(row.UserID.Bytes),
		Rating:    row.Rating,
		Comment:   row.Comment.String,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}