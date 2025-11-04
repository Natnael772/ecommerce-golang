package auth

import (
	"context"
	"ecommerce-app/configs"
	"ecommerce-app/internal/pkg/errs"
	"ecommerce-app/internal/pkg/jwt"
	"ecommerce-app/internal/pkg/password"
	"time"
)

type Service interface {
	Login(ctx context.Context, req LoginRequest) (UserWithToken, *errs.AppError)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}


func (s *service) Login(ctx context.Context, req LoginRequest) (UserWithToken, *errs.AppError) {
	u, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return UserWithToken{}, errs.ErrUnauthorized.WithMessage("invalid email or password")
	}

	isPasswordValid := password.Check(req.Password, u.Password)
	if !isPasswordValid {
		return UserWithToken{}, errs.ErrUnauthorized.WithMessage("invalid email or password")
	}

	cfg:= configs.Load()
	jwtSecret := cfg.JWTSecret
	jwtDuration := cfg.JWTDuration

	token, err := jwt.GenerateToken(jwtSecret, u.ID, u.Email, u.Role, time.Duration(jwtDuration)*time.Second)
	if err != nil {
		return UserWithToken{}, errs.ErrInternal.WithMessage("failed to generate auth token")
	}

	data := UserWithToken{
		User:  u,
		Token: token,
	}
	
	return data, nil
}


