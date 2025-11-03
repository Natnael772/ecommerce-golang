package cart

import (
	"context"
	"database/sql"
	"ecommerce-app/internal/pkg/database/sqlc"
	"ecommerce-app/internal/pkg/errs"
	"ecommerce-app/internal/pkg/logger"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	Create(ctx context.Context, c CreateCartRequest) (Cart, error)
	GetByID(ctx context.Context, id string) (Cart, error)
	GetByUserID(ctx context.Context, userID string) (Cart, error)
	List(ctx context.Context, limit, offset int32) ([]Cart, error)
	Update(ctx context.Context, id string, req UpdateCartRequest) (Cart, error)
	Delete(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}

type repository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repository{q}
}

// --- Create ---
func (r *repository) Create(ctx context.Context, c CreateCartRequest) (Cart, error) {
	expiresAt := pgtype.Timestamptz{Valid: false}
	if !c.ExpiresAt.IsZero() {
		expiresAt = pgtype.Timestamptz{Time: *c.ExpiresAt, Valid: true}
	}

	var uuidUser pgtype.UUID
	if err := uuidUser.Scan(c.UserID); err != nil {
		logger.Error("Invalid user ID: %v", err)
		return Cart{}, err
	}

	params := sqlc.CreateCartParams{
		UserID:    uuidUser,
		ExpiresAt: expiresAt,
	}

	row, err := r.q.CreateCart(ctx, params)
	if err != nil {
		logger.Error("Error creating cart: %v", err)
		return Cart{}, err
	}

	return MapCart(row), nil
}

// --- Get By ID ---
func (r *repository) GetByID(ctx context.Context, id string) (Cart, error) {
	var uuidID pgtype.UUID
	if err := uuidID.Scan(id); err != nil {
		return Cart{}, err
	}

	row, err := r.q.GetCartByID(ctx, uuidID)
	if err != nil {
		logger.Error("Error getting cart by ID: %v", err)
		return Cart{}, err
	}

	return MapCart(row), nil
}

// --- Get By User ID ---
func (r *repository) GetByUserID(ctx context.Context, userID string) (Cart, error) {
	var uuidUser pgtype.UUID
	if err := uuidUser.Scan(userID); err != nil {
		return Cart{}, err
	}

	row, err := r.q.GetCartByUserID(ctx, uuidUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Cart{}, errs.ErrNotFound
		}

		return Cart{}, err
	}

	return MapCart(row), nil
}

// --- List ---
func (r *repository) List(ctx context.Context, limit, offset int32) ([]Cart, error) {
	rows, err := r.q.ListCarts(ctx, sqlc.ListCartsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Error("Error listing carts: %v", err)
		return nil, err
	}

	carts := make([]Cart, len(rows))
	for i, row := range rows {
		carts[i] = MapCart(row)
	}

	return carts, nil
}

// --- Update ---
func (r *repository) Update(ctx context.Context, id string, req UpdateCartRequest) (Cart, error) {
	var uuidID pgtype.UUID
	if err := uuidID.Scan(id); err != nil {
		return Cart{}, err
	}

	var expiresAt pgtype.Timestamptz
	if req.ExpiresAt != nil {
		expiresAt = pgtype.Timestamptz{Time: *req.ExpiresAt, Valid: true}
	} else {
		expiresAt = pgtype.Timestamptz{Valid: false}
	}

	params := sqlc.UpdateCartParams{
		ID:        uuidID,
		ExpiresAt: expiresAt,
	}

	row, err := r.q.UpdateCart(ctx, params)
	if err != nil {
		logger.Error("Error updating cart: %v", err)
		return Cart{}, err
	}

	return MapCart(row), nil
}

// --- Delete ---
func (r *repository) Delete(ctx context.Context, id string) error {
	var uuidID pgtype.UUID
	if err := uuidID.Scan(id); err != nil {
		return err
	}

	err := r.q.DeleteCart(ctx, uuidID)
	if err != nil {
		logger.Error("Error deleting cart: %v", err)
		return err
	}
	return nil
}

// --- Delete Expired ---
func (r *repository) DeleteExpired(ctx context.Context) error {
	err := r.q.DeleteExpiredCarts(ctx)
	if err != nil {
		logger.Error("Error deleting expired carts: %v", err)
		return err
	}
	return nil
}

// --- Mapper ---
func MapCart(row sqlc.Cart) Cart {
	return Cart{
		ID:        row.ID.Bytes,
		UserID:    row.UserID.Bytes,
		ExpiresAt: row.ExpiresAt.Time,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
