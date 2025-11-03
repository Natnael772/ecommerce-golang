package category

import (
	"context"
	"ecommerce-app/internal/pkg/database/sqlc"
	"ecommerce-app/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface{
	Create(ctx context.Context, c Category) (Category, error)
	GetByID(ctx context.Context, id string) (Category, error)
	GetBySlug(ctx context.Context, slug string) (Category, error)
	List(ctx context.Context, limit, offset int32) ([]Category, error)
	Update(ctx context.Context, c Category) (Category, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int32, error)
}

type repository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repository{q}
}

func (r *repository) Create(ctx context.Context, c Category) (Category, error) {
	var parentID pgtype.UUID
	if c.ParentID != nil {
		parentID = pgtype.UUID{Bytes: [16]byte(*c.ParentID), Valid: true}
	} else {
		parentID = pgtype.UUID{Valid: false}
	}

	params := sqlc.CreateCategoryParams{
		Name:     c.Name,
		Slug:     c.Slug,
		ParentID: parentID,
	}

	row, err := r.q.CreateCategory(ctx, params)
	if err != nil {
		logger.Error("Error creating category: %v", err)
		return Category{}, err
	}

	return mapCategory(row), nil
}

func (r *repository) GetByID(ctx context.Context, id string) (Category, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return Category{}, err
	}

	row, err := r.q.GetCategoryByID(ctx, uuid)
	if err != nil {
		return Category{}, err
	}
	return mapCategory(row), nil
}

func (r *repository) GetBySlug(ctx context.Context, slug string) (Category, error) {
	row, err := r.q.GetCategoryBySlug(ctx, slug)
	if err != nil {
		return Category{}, err
	}
	return mapCategory(row), nil
}

func (r *repository) List(ctx context.Context, limit, offset int32) ([]Category, error) {
	rows, err := r.q.ListCategories(ctx, sqlc.ListCategoriesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	categories := make([]Category, len(rows))
	for i, row := range rows {
		categories[i] = mapCategory(row)
	}

	return categories, nil
}

func (r *repository) Update(ctx context.Context, c Category) (Category, error) {
	var parentID pgtype.UUID
	if c.ParentID != nil {
		parentID = pgtype.UUID{Bytes: [16]byte(*c.ParentID), Valid: true}
	} else {
		parentID = pgtype.UUID{Valid: false}
	}

	params := sqlc.UpdateCategoryParams{
		ID:       pgtype.UUID{Bytes: c.ID, Valid: true},
		Name:     c.Name,
		Slug:     c.Slug,
		ParentID: parentID,
	}

	row, err := r.q.UpdateCategory(ctx, params)
	if err != nil {
		return Category{}, err
	}

	return mapCategory(row), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return err
	}
	return r.q.DeleteCategory(ctx, uuid)
}

func (r *repository) Count(ctx context.Context) (int32, error) {
	count, err := r.q.CountCategories(ctx)
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

func mapCategory(row sqlc.Category) Category {	
	var parentID *uuid.UUID
	if row.ParentID.Valid {
		id := uuid.UUID(row.ParentID.Bytes)
		parentID = &id
	}

	return Category{
		ID:        row.ID.Bytes,
		Name:      row.Name,
		Slug:      row.Slug,
		ParentID:  parentID,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
	