package coupon

import (
	"context"
	"ecommerce-app/internal/pkg/errs"
	"ecommerce-app/internal/pkg/response"
	"ecommerce-app/pkg/pagination"

	"github.com/google/uuid"
)

type Service interface {
	CreateCoupon(ctx context.Context, req CreateCouponRequest) (Coupon, *errs.AppError)
	GetCouponByID(ctx context.Context, id string) (Coupon, *errs.AppError)
	GetCouponByCode(ctx context.Context, code string) (Coupon, *errs.AppError)
	GetCoupons(ctx context.Context, page, perPage int) (CouponsWithMeta, *errs.AppError)
	IncrementCouponUsage(ctx context.Context, code string) (Coupon, *errs.AppError)
	UpdateCoupon(ctx context.Context, id string, req UpdateCouponRequest) (Coupon, *errs.AppError)
	DeleteCoupon(ctx context.Context, id string) *errs.AppError
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateCoupon(ctx context.Context, req CreateCouponRequest) (Coupon, *errs.AppError) {
	existingCoupon, err := s.repo.GetByCode(ctx, req.Code)
	if err == nil && existingCoupon.ID != uuid.Nil {
		return Coupon{}, errs.ErrConflict.WithMessage("Coupon with the same code already exists")
	}

	coupon := Coupon{
		Code:            req.Code,
		Description:     req.Description,
		DiscountPercent: int32(req.DiscountPercent),
		ValidFrom:       req.ValidFrom,
		ValidUntil:      req.ValidUntil,
		MaxUses:         req.MaxUses,
		IsActive:        req.IsActive,
	}

	createdCoupon, err := s.repo.Create(ctx, coupon)
	if err != nil {
		return Coupon{}, errs.ErrInternal.WithMessage("Failed to create coupon")
	}

	return createdCoupon, nil
}


func (s *service) GetCouponByID(ctx context.Context, id string) (Coupon, *errs.AppError) {
	coupon, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Coupon{}, errs.ErrInternal.WithMessage("Failed to get coupon by ID")
	}
	return coupon, nil
}

func (s *service) GetCoupons(ctx context.Context, page, perPage int) (CouponsWithMeta, *errs.AppError) {
	p:= pagination.New(page, perPage)

	limit := int32(p.PerPage)
	offset := int32(p.Offset())

	coupons, err := s.repo.Get(ctx, limit, offset)
	if err != nil {
		return CouponsWithMeta{}, errs.ErrInternal.WithMessage("Failed to get coupons")
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return CouponsWithMeta{}, errs.ErrInternal.WithMessage("Failed to count coupons")
	}

	meta := response.Meta{
		Page:         p.Page,
		PerPage:      p.PerPage,
		Total:   	(int(count)),
	}

	result := CouponsWithMeta{
		Coupons: coupons,
		Meta:    meta,
	}

	return result, nil
}

func (s *service) GetCouponByCode(ctx context.Context, code string) (Coupon, *errs.AppError) {
	coupon, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return Coupon{}, errs.ErrInternal.WithMessage("Failed to get coupon by code")
	}
	return coupon, nil
}


func (s *service) IncrementCouponUsage(ctx context.Context, code string) (Coupon, *errs.AppError) {
	updatedCoupon, err := s.repo.IncrementUsage(ctx, code)
	if err != nil {
		return Coupon{}, errs.ErrInternal.WithMessage("Failed to increment coupon usage")
	}
	return updatedCoupon, nil
}

func (s *service) UpdateCoupon(ctx context.Context, id string, req UpdateCouponRequest) (Coupon, *errs.AppError) {
	updatedCoupon, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return Coupon{}, errs.ErrInternal.WithMessage("Failed to update coupon")
	}
	return updatedCoupon, nil
}


func (s *service) DeleteCoupon(ctx context.Context, id string) *errs.AppError {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return errs.ErrInternal.WithMessage("Failed to delete coupon")
	}
	return nil
}