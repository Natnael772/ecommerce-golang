package inventory

import (
	"context"
	"ecommerce-app/internal/pkg/database/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)


type Repository interface {
	CreateInventory(ctx context.Context, productID string, stock int32, reserved int32) (Inventory, error)
	GetInventoryByProductID(ctx context.Context, productID string) (Inventory, error)
	UpdateInventoryStock(ctx context.Context, productID string, stock int32) (Inventory, error)
	DeleteInventory(ctx context.Context, productID string) error
}


type repository struct {
	queries *sqlc.Queries
}


func NewRepository(queries *sqlc.Queries) Repository {
	return &repository{queries: queries}
}

func (r *repository) CreateInventory(ctx context.Context, productID string, stock int32, reserved int32) (Inventory, error) {
	var productUUID pgtype.UUID
	if err := productUUID.Scan(productID); err != nil {
		return Inventory{}, err
	}

	params := sqlc.CreateInventoryParams{
		ProductID: productUUID,
		Stock:     stock,
		Reserved:  reserved,
	}

	row, err := r.queries.CreateInventory(ctx, params)
	if err != nil {
		return Inventory{}, err
	}

	return mapInventory(row), nil
}

func (r *repository) GetInventoryByProductID(ctx context.Context, productID string) (Inventory, error) {
	var productUUID pgtype.UUID
	if err := productUUID.Scan(productID); err != nil {
		return Inventory{}, err
	}

	row, err := r.queries.GetInventoryByProductID(ctx, productUUID)
	if err != nil {
		return Inventory{}, err
	}

	return mapInventory(row), nil
}

func (r *repository) UpdateInventoryStock(ctx context.Context, productID string, stock int32) (Inventory, error) {
	var productUUID pgtype.UUID
	if err := productUUID.Scan(productID); err != nil {
		return Inventory{}, err
	}

	params := sqlc.UpdateInventoryStockParams{
		ProductID: productUUID,
		Stock:     stock,
	}

	row, err := r.queries.UpdateInventoryStock(ctx, params)
	if err != nil {
		return Inventory{}, err
	}

	return mapInventory(row), nil
}

func (r *repository) DeleteInventory(ctx context.Context, productID string) error {
	var productUUID pgtype.UUID
	if err := productUUID.Scan(productID); err != nil {
		return err
	}

	return r.queries.DeleteInventory(ctx, productUUID)
}

func mapInventory(row sqlc.Inventory) Inventory {
	return Inventory{
		ProductID: row.ProductID.String(),
		Stock:     row.Stock,
		Reserved:  row.Reserved,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}


