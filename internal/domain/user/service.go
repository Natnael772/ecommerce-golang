package user

import (
	"context"
	"ecommerce-app/configs"
	"ecommerce-app/internal/pkg/errs"
	"ecommerce-app/internal/pkg/jwt"
	"ecommerce-app/internal/pkg/password"
	"ecommerce-app/internal/pkg/response"
	"ecommerce-app/pkg/pagination"
	"errors"
	"time"
)

type Service interface {
	Register(ctx context.Context, req RegisterUserRequest) (User, *errs.AppError)
	Login(ctx context.Context, req LoginRequest) (UserWithToken, *errs.AppError)
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
	user := User{
		Email:    req.Email,
		Password: hash,
		FirstName: req.FirstName,
		LastName: req.LastName,
		Role: req.Role,
	}
	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return User{}, errs.ErrInternal.WithMessage("failed to create user")
	}
	return createdUser, nil
}

// Login checks user credentials
func (s *service) Login(ctx context.Context, req LoginRequest) (UserWithToken, *errs.AppError) {
	u, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return UserWithToken{}, errs.ErrUnauthorized.WithMessage("invalid email or password")
	}

	isPasswordValid := password.Check(req.Password, u.Password)
	if !isPasswordValid {
		return UserWithToken{}, errs.ErrUnauthorized.WithMessage("invalid email or password")
	}

	jwtSecret := configs.Load().JWTSecret
	jwtDuration := configs.Load().JWTDuration
	
	token, err := jwt.GenerateToken(jwtSecret, u.ID, u.Email, time.Duration(jwtDuration))
	if err != nil {
		return UserWithToken{}, errs.ErrInternal.WithMessage("failed to generate auth token")
	}

	data := UserWithToken{
		User:  u,
		Token: token,
	}
	
	return data, nil
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
