package auth

import (
	"context"
	"database/sql"
	"ecommerce-app/internal/pkg/database/sqlc"
	"ecommerce-app/internal/pkg/errs"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id string) (User, error)
}

type repository struct {
	queries *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repository{q}
}



func (r *repository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, errs.ErrNotFound
		}
		return User{}, err
	}

	return mapUser(row), nil
}

func (r *repository) GetUserByID(ctx context.Context, id string) (User, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return User{}, err
	}
	row, err := r.queries.GetUserByID(ctx, uuid)
	if err != nil {
		return User{}, err
	}
	return mapUser(row), nil
}


func mapUser(row sqlc.User) User {
	return User{
		ID:        row.ID.String(),
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Email:     row.Email,
		Password:  row.Password,
		Role:      string(row.Role),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
