package user

import (
	"context"
	"ecommerce-app/internal/pkg/errs"
	"ecommerce-app/internal/pkg/password"
	"ecommerce-app/internal/pkg/response"
	"ecommerce-app/pkg/pagination"
	"errors"
)

type Service interface {
	Register(ctx context.Context, req RegisterUserRequest) (User, *errs.AppError)
	GetUser(ctx context.Context, id string) (User, error)
	ListUsers(ctx context.Context, page, itemsPerPage int32) (UsersWithMeta, *errs.AppError)
	
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Register(ctx context.Context, req RegisterUserRequest) (User, *errs.AppError) {
	exists, err := s.repo.GetByEmail(ctx, req.Email)

	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return User{}, errs.ErrInternal.WithMessage("failed to check existing user")
	}


	if exists.ID != "" {
		return User{}, errs.ErrConflict.WithMessage("email already in use")
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		return User{}, errs.ErrInternal.WithMessage("failed to hash password")
	}
	
	userData := User{
		Email:    req.Email,
		Password: hash,
		FirstName: req.FirstName,
		LastName: req.LastName,
		Role: req.Role,
	}
	user, err := s.repo.Create(ctx, userData)
	if err != nil {
		return User{}, errs.ErrInternal.WithMessage("failed to create user")
	}

	return user, nil
}

func (s *service) GetUser(ctx context.Context, id string) (User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) ListUsers(ctx context.Context, page, perPage int32) (UsersWithMeta, *errs.AppError){
	p := pagination.New(int(page), int(perPage)) 
	limit := int32(p.PerPage)
	offset := int32(p.Offset())

	users,err:= s.repo.List(ctx, limit, offset)
	if err != nil {
		return UsersWithMeta{}, errs.ErrInternal.WithMessage("failed to list users")
	}

	total,err:= s.repo.Count(ctx)
	if err != nil {
		return UsersWithMeta{}, errs.ErrInternal.WithMessage("failed to count users")
	}

	meta := response.Meta{
		Page: 		 int(p.Page),
		PerPage:     int(p.PerPage),
		Total:       int(total),
	}

	result := UsersWithMeta{
		Users: users,
		Meta:  meta,
	}

	return result, nil
}
