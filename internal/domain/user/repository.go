package user

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
	Create(ctx context.Context, u User) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id string) (User, error)
	List(ctx context.Context, limit, offset int32) ([]User, error)
	Count(ctx context.Context) (int32, error)
}

type repository struct {
	queries *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repository{q}
}

func (r *repository) Create(ctx context.Context, u User) (User, error) {
	params := sqlc.CreateUserParams{
		Email:     u.Email,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      sqlc.UserRole(u.Role),
	}
	row, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		logger.Error("Error creating user: %v", err)
		return User{}, err
	}
	return mapUser(row), nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, errs.ErrNotFound
		}
		return User{}, err
	}

	return mapUser(row), nil
}

func (r *repository) GetByID(ctx context.Context, id string) (User, error) {
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

func (r *repository) List(ctx context.Context, limit, offset int32) ([]User, error) {
	rows, err := r.queries.ListUsers(ctx, sqlc.ListUsersParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, len(rows))
	for _, row := range rows {
		users = append(users, mapUser(row))
	}
	return users, nil
}


func (r *repository) Count(ctx context.Context) (int32, error) {
	count, err := r.queries.CountUsers(ctx)
	if err != nil {
		return 0, err
	}
	return int32(count), nil
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
