package product

import (
	"context"
	"ecommerce-app/internal/pkg/database/sqlc"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	// "ecommerce-app/internal/db/sql"
)

type Repository interface {
	Create(ctx context.Context, p Product) (Product, error)
	GetByID(ctx context.Context, id string) (Product, error)
	GetBySku(ctx context.Context, sku string) (Product, error)
	List(ctx context.Context, limit, offset int32) ([]Product, error)
	Count(ctx context.Context) (int32, error)
	UpdatePrice(ctx context.Context, id string, price int32) (Product, error)
	UpdateProduct(ctx context.Context,id string, req UpdateProductRequest) (Product, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repository{q}
}
func (r *repository) Create(ctx context.Context, p Product) (Product, error) {
	desc := pgtype.Text{String: p.Description, Valid: true}
	cat := pgtype.UUID{Bytes: p.CategoryID, Valid: true}
	mainImage := pgtype.Text{String: p.MainImageUrl, Valid: true}

	imagesBytes, err := json.Marshal(p.Images)
	if err != nil {
		return Product{}, err
	}

	attributes, err := json.Marshal(p.Attributes)
	if err != nil {
		// attributes = nil
		return Product{}, err
	}

	var discountValidUntil pgtype.Timestamptz
	if p.DiscountValidUntil != nil {
		discountValidUntil = pgtype.Timestamptz{
			Time:  *p.DiscountValidUntil,
			Valid: true,
		}
	} else {
		discountValidUntil = pgtype.Timestamptz{
			Valid: false,
		}
	}

	params := sqlc.CreateProductParams{
		Sku:            p.SKU,
		Name:           p.Name,
		Description:    desc,
		CategoryID:     cat,
		PriceCents:     p.PriceCents,
		Currency:       p.Currency,
		Attributes:     attributes,
		MainImageUrl:   mainImage,
		Images:         imagesBytes,
		DiscountPercent: pgtype.Int4{Int32: p.DiscountPercent, Valid: true},
		DiscountValidUntil: discountValidUntil,
	}

	row, err := r.q.CreateProduct(ctx, params)
	if err != nil {
		return Product{}, err
	}

	return mapProduct(row), nil
}

func (r *repository) GetByID(ctx context.Context, id string) (Product, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return Product{}, err
	}

	row, err := r.q.GetProductByID(ctx, uuid)
	if err != nil {
		return Product{}, err
	}
	return mapProduct(row), nil
}

func (r *repository) GetBySku(ctx context.Context, sku string) (Product, error) {
	row, err := r.q.GetProductBySKU(ctx, sku)
	if err != nil {
		return Product{}, err
	}
	
	return mapProduct(row), nil
}

func (r *repository) List(ctx context.Context, limit, offset int32) ([]Product, error) {
	rows, err := r.q.ListProducts(ctx, sqlc.ListProductsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]Product, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapProduct(r))
	}
	return out, nil
}

func (r *repository) Count(ctx context.Context) (int32, error) {
	count, err := r.q.CountProducts(ctx)
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

func (r *repository) UpdatePrice(ctx context.Context, id string, price int32) (Product, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return Product{}, err
	}

	row, err := r.q.UpdateProductPrice(ctx, sqlc.UpdateProductPriceParams{
		ID:         uuid,
		PriceCents: price,
	})
	if err != nil {
		return Product{}, err
	}
	return mapProduct(row), nil
}

func (r *repository) UpdateProduct(ctx context.Context,	id string, p UpdateProductRequest) (Product, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return Product{}, err
	}

	desc := pgtype.Text{String: *p.Description, Valid: true}
	cat := pgtype.UUID{Bytes: *p.CategoryID, Valid: true}
	mainImage := pgtype.Text{String: *p.Currency, Valid: true}

	imagesBytes, err := json.Marshal(p.Images)
	if err != nil {
		return Product{}, err
	}

	// check if attributes is nil
	attributes, err := json.Marshal(p.Attributes)
	if err != nil {
		return Product{}, err
	}
	

	params := sqlc.UpdateProductParams{
		ID:             uuid,
		Name:           *p.Name,
		Description:    desc,
		CategoryID:     cat,
		PriceCents:     *p.PriceCents,
		Currency:       *p.Currency,
		Attributes:     attributes,
		MainImageUrl:   mainImage,
		Images:         imagesBytes,
		DiscountPercent: pgtype.Int4{Int32: *p.DiscountPercent, Valid: true},
	}

	row, err := r.q.UpdateProduct(ctx, params)
	if err != nil {
		return Product{}, err
	}

	return mapProduct(row), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return err
	}
	return r.q.DeleteProduct(ctx, uuid)
}

func mapProduct(row sqlc.Product) Product {
	images:= []string{}
	if err := json.Unmarshal(row.Images, &images); err != nil {
		images = []string{}
	}

	var attributes map[string]interface{}
	err := json.Unmarshal(row.Attributes, &attributes)
	if err != nil {
		attributes=nil
	}

	var discountValidUntil *time.Time
	if row.DiscountValidUntil.Valid {
		t := row.DiscountValidUntil.Time.UTC()
		discountValidUntil = &t
	}


	return Product{
		ID:          row.ID.Bytes,
		SKU:         row.Sku,
		Name:        row.Name,
		Description: row.Description.String,
		CategoryID:  row.CategoryID.Bytes,
		PriceCents:  row.PriceCents,
		Currency:    row.Currency,
		Attributes:  attributes,
		MainImageUrl: row.MainImageUrl.String,
		Images: 	images,
		DiscountPercent: row.DiscountPercent.Int32,
		DiscountValidUntil: discountValidUntil,
		IsActive:    row.IsActive.Bool,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}
