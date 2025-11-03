package coupon

import (
	"context"
	"ecommerce-app/internal/pkg/database/sqlc"
	"ecommerce-app/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface{
	Create(ctx context.Context, c Coupon) (Coupon, error)
	GetByID(ctx context.Context, id string) (Coupon, error)
	GetByCode(ctx context.Context, code string) (Coupon, error)
	Get(ctx context.Context, limit, offset int32) ([]Coupon, error)
	Update(ctx context.Context, id string, req UpdateCouponRequest) (Coupon, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int32, error)
	IncrementUsage(ctx context.Context, code string) (Coupon, error)
}

type repository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repository{q}
}

func (r *repository) Create(ctx context.Context, c Coupon) (Coupon, error) {
	validFrom := pgtype.Timestamptz{Time: c.ValidFrom, Valid: true}
	validUntil := pgtype.Timestamptz{Time: c.ValidUntil, Valid: true}
	
	params := sqlc.CreateCouponParams{
		Code:            c.Code,
		Description:     pgtype.Text{String: c.Description, Valid: true},
		DiscountPercent: c.DiscountPercent,
		ValidFrom:       validFrom,
		ValidUntil:      validUntil,
		MaxUses:         pgtype.Int4{Int32: c.MaxUses, Valid: true},
		IsActive:        pgtype.Bool{Bool: c.IsActive, Valid: true},
	}

	row, err := r.q.CreateCoupon(ctx, params)
	if err != nil {
		logger.Error("Error creating coupon: %v", err)
		return Coupon{}, err
	}

	return mapCoupon(row), nil
}	

func (r *repository) Delete(ctx context.Context, id string) error {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return err
	}

	err := r.q.DeleteCoupon(ctx, uuid)
	if err != nil {
		logger.Error("Error deleting coupon: %v", err)
		return err
	}

	return nil
}

func (r *repository) GetByCode(ctx context.Context, code string) (Coupon, error) {
	row, err := r.q.GetCouponByCode(ctx, code)
	if err != nil {
		logger.Error("Error getting coupon by code: %v", err)
		return Coupon{}, err
	}
	return mapCoupon(row), nil
}

func (r *repository) GetByID(ctx context.Context, id string) (Coupon, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return Coupon{}, err
	}

	row, err := r.q.GetCouponById(ctx, uuid)
	if err != nil {
		logger.Error("Error getting coupon by ID: %v", err)
		return Coupon{}, err
	}
	return mapCoupon(row), nil
}
// func (r *repository) Update(ctx context.Context, id string, req UpdateCouponRequest) (Coupon, error) {
// 	var uuid pgtype.UUID
// 	if err := uuid.Scan(id); err != nil {
// 		return Coupon{}, err
// 	}

// 	// only set fields that are not nil in the request
// 	var code pgtype.Text
// 	if req.Code != nil {
// 		code = pgtype.Text{String: *req.Code, Valid: true}
// 	} else {
// 		code = pgtype.Text{Valid: false}
// 	}

// 	var description pgtype.Text
// 	if req.Description != nil {
// 		description = pgtype.Text{String: *req.Description, Valid: true}
// 	} else {
// 		description = pgtype.Text{Valid: false}
// 	}

// 	var discountPercent pgtype.Int4
// 	if req.DiscountPercent != nil {
// 		discountPercent = pgtype.Int4{Int32: int32(*req.DiscountPercent), Valid: true}
// 	} else {
// 		discountPercent = pgtype.Int4{Valid: false}
// 	}

// 	var validFrom pgtype.Timestamptz
// 	if req.ValidFrom != nil {
// 		validFrom = pgtype.Timestamptz{Time: *req.ValidFrom, Valid: true}
// 	} else {
// 		validFrom = pgtype.Timestamptz{Valid: false}
// 	}

// 	var validUntil pgtype.Timestamptz
// 	if req.ValidUntil != nil {
// 		validUntil = pgtype.Timestamptz{Time: *req.ValidUntil, Valid: true}
// 	} else {
// 		validUntil = pgtype.Timestamptz{Valid: false}
// 	}

// 	var maxUses pgtype.Int4
// 	if req.MaxUses != nil {
// 		maxUses = pgtype.Int4{Int32: *req.MaxUses, Valid: true}
// 	} else {
// 		maxUses = pgtype.Int4{Valid: false}
// 	}

// 	var isActive pgtype.Bool
// 	if req.IsActive != nil {
// 		isActive = pgtype.Bool{Bool: *req.IsActive, Valid: true}
// 	} else {
// 		isActive = pgtype.Bool{Valid: false}
// 	}

// 	params := sqlc.UpdateCouponParams{
// 		ID:              uuid,
// 		Code:            code.String,
// 		Description:     description,
// 		DiscountPercent: discountPercent.Int32,
// 		ValidFrom:       validFrom,
// 		ValidUntil:      validUntil,
// 		MaxUses:         maxUses,
// 		IsActive:        isActive,
// 	}

// 	row, err := r.q.UpdateCoupon(ctx, params)
// 	if err != nil {
// 		logger.Error("Error updating coupon: %v", err)
// 		return Coupon{}, err
// 	}

// 	return mapCoupon(row), nil
// }

func (r *repository) Update(ctx context.Context, id string, req UpdateCouponRequest) (Coupon, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return Coupon{}, err
	}

	var code pgtype.Text
	if req.Code != nil {
		code = pgtype.Text{String: *req.Code, Valid: true}
	} else {
		code = pgtype.Text{Valid: false}
	}

	var description pgtype.Text
	if req.Description != nil {
		description = pgtype.Text{String: *req.Description, Valid: true}
	} else {
		description = pgtype.Text{Valid: false}
	}

	var discountPercent pgtype.Int4
	if req.DiscountPercent != nil {
		discountPercent = pgtype.Int4{Int32: int32(*req.DiscountPercent), Valid: true}
	} else {
		discountPercent = pgtype.Int4{Valid: false}
	}

	var validFrom pgtype.Timestamptz
	if req.ValidFrom != nil {
		validFrom = pgtype.Timestamptz{Time: *req.ValidFrom, Valid: true}
	} else {
		validFrom = pgtype.Timestamptz{Valid: false}
	}

	var validUntil pgtype.Timestamptz
	if req.ValidUntil != nil {
		validUntil = pgtype.Timestamptz{Time: *req.ValidUntil, Valid: true}
	} else {
		validUntil = pgtype.Timestamptz{Valid: false}
	}

	var maxUses pgtype.Int4
	if req.MaxUses != nil {
		maxUses = pgtype.Int4{Int32: *req.MaxUses, Valid: true}
	} else {
		maxUses = pgtype.Int4{Valid: false}
	}

	var isActive pgtype.Bool
	if req.IsActive != nil {
		isActive = pgtype.Bool{Bool: *req.IsActive, Valid: true}
	} else {
		isActive = pgtype.Bool{Valid: false}
	}

	params := sqlc.UpdateCouponParams{
		ID:              uuid,
		Code:            code.String,
		Description:     description,
		DiscountPercent: discountPercent.Int32,
		ValidFrom:       validFrom,
		ValidUntil:      validUntil,
		MaxUses:         maxUses,
		IsActive:        isActive,
	}

	row, err := r.q.UpdateCoupon(ctx, params)
	if err != nil {
		logger.Error("Error updating coupon: %v", err)
		return Coupon{}, err
	}

	return mapCoupon(row), nil
}



func (r *repository) Get(ctx context.Context, limit, offset int32) ([]Coupon, error) {
	rows, err := r.q.GetCoupons(ctx, sqlc.GetCouponsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Error("Error listing coupons: %v", err)
		return nil, err
	}

	coupons := make([]Coupon, len(rows))
	for i, row := range rows {
		coupons[i] = mapCoupon(row)
	}

	return coupons, nil
}

func (r *repository) Count(ctx context.Context) (int32, error) {
	count, err := r.q.CountCoupons(ctx)
	if err != nil {
		logger.Error("Error counting coupons: %v", err)
		return 0, err
	}
	return int32(count), nil
}

func (r *repository) IncrementUsage(ctx context.Context, code string) (Coupon, error) {
	row, err := r.q.IncrementCouponUsage(ctx, code)
	if err != nil {
		logger.Error("Error incrementing coupon usage: %v", err)
		return Coupon{}, err
	}

	return mapCoupon(row), nil
}

func mapCoupon(row sqlc.Coupon) Coupon {

	return Coupon{
		ID:              uuid.UUID(row.ID.Bytes),
		Code:            row.Code,
		Description:     row.Description.String,
		DiscountPercent: row.DiscountPercent,
		ValidFrom:       row.ValidFrom.Time,
		ValidUntil:      row.ValidUntil.Time,
		MaxUses:         row.MaxUses.Int32,
		UsedCount:       row.UsedCount.Int32,
		IsActive:        row.IsActive.Bool,
		IsDeleted:       row.IsDeleted.Bool,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}
}