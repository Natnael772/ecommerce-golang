package cartitem

import (
	"context"
	"ecommerce-app/internal/domain/cart"
	"ecommerce-app/internal/pkg/database/sqlc"
	"ecommerce-app/internal/pkg/logger"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	Add(ctx context.Context, item CartItem) (CartItem, error)
	AddItems(ctx context.Context, req AddItemsRequest) ([]CartItem, error)
	GetByID(ctx context.Context, id string) (CartItem, error)
	GetByUserID(ctx context.Context, userID string) ([]CartItem, error)
	GetByCartAndProduct(ctx context.Context, cartID, productID string) (CartItem, error)
	ListByCart(ctx context.Context, cartID string) ([]CartItem, error)
	UpdateQuantity(ctx context.Context, id string, quantity int32) (CartItem, error)
	Delete(ctx context.Context, id string) error
	DeleteByCart(ctx context.Context, cartID string) error
	GetCartByUserID(ctx context.Context, userID string) (cart.Cart, error)
}

type repository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repository{q: q}
}

// Add or Upsert Item
func (r *repository) Add(ctx context.Context, item CartItem) (CartItem, error) {
	var cartUUID, productUUID pgtype.UUID
	if err := cartUUID.Scan(item.CartID.String()); err != nil {
		return CartItem{}, err
	}
	if err := productUUID.Scan(item.ProductID.String()); err != nil {
		return CartItem{}, err
	}
	params := sqlc.AddCartItemParams{
		CartID:    cartUUID,
		ProductID: productUUID,
		Quantity:  item.Quantity,
	}

	row, err := r.q.AddCartItem(ctx, params)
	if err != nil {
		logger.Error("Error adding cart item: %v", err)
		return CartItem{}, err
	}

	return mapCartItem(row), nil
}

// Add or Upsert Multiple Items
func (r *repository) AddItems(ctx context.Context, req AddItemsRequest) ([]CartItem, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("no items to add")
	}

	// Aggregate duplicate products
	aggMap := make(map[uuid.UUID]int32)
	for _, item := range req.Items {
		pid, err := uuid.Parse(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID %s: %w", item.ProductID, err)
		}
		aggMap[pid] += item.Quantity
	}

	// Convert map to slices for SQL parameters and to pgtype.UUID
	productIDsPg := make([]pgtype.UUID, 0, len(aggMap))
	quantities := make([]int32, 0, len(aggMap))
	for id, qty := range aggMap {
		var pgID pgtype.UUID
		if err := pgID.Scan(id.String()); err != nil {
			return nil, fmt.Errorf("failed to convert product id %s: %w", id.String(), err)
		}
		productIDsPg = append(productIDsPg, pgID)
		quantities = append(quantities, qty)
	}

	// Convert cart id to pgtype.UUID
	var cartUUID pgtype.UUID
	if err := cartUUID.Scan(req.CartID.String()); err != nil {
		return nil, fmt.Errorf("invalid cart ID %s: %w", req.CartID.String(), err)
	}

	// Prepare params for sqlc query
	params := sqlc.AddCartItemsParams{
		CartID:  cartUUID,
		Column2: productIDsPg,
		Column3: quantities,
	}

	// Execute query
	rows, err := r.q.AddCartItems(ctx, params)
	if err != nil {
		logger.Error("Error adding or updating cart items: %v", err)
		return nil, fmt.Errorf("failed to add cart items: %w", err)
	}

	// Map results to CartItem
	items := make([]CartItem, len(rows))
	for i, row := range rows {
		items[i] = mapCartItem(row)
	}

	return items, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (CartItem, error) {
	var uuidID pgtype.UUID
	if err := uuidID.Scan(id); err != nil {
		return CartItem{}, err
	}

	row, err := r.q.GetCartItem(ctx, uuidID)
	if err != nil {
		logger.Error("Error getting cart item by ID: %v", err)
		return CartItem{}, err
	}

	return mapCartItem(row), nil
}

func (r *repository) GetByUserID(ctx context.Context, userID string) ([]CartItem, error) {
	var uuidUser pgtype.UUID
	if err := uuidUser.Scan(userID); err != nil {
		return nil, err
	}

	rows, err := r.q.GetCartItemsByUserID(ctx, uuidUser)
	if err != nil {
		logger.Error("Error getting cart items by user ID: %v", err)
		return nil, err
	}

	items := make([]CartItem, len(rows))
	for i, row := range rows {
		items[i] = mapCartItem(row)
	}

	return items, nil
}

func (r *repository) GetByCartAndProduct(ctx context.Context, cartID, productID string) (CartItem, error) {
	var cartUUID, productUUID pgtype.UUID
	if err := cartUUID.Scan(cartID); err != nil {
		return CartItem{}, err
	}
	if err := productUUID.Scan(productID); err != nil {
		return CartItem{}, err
	}

	row, err := r.q.GetCartItemByProduct(ctx, sqlc.GetCartItemByProductParams{
		CartID:    cartUUID,
		ProductID: productUUID,
	})
	if err != nil {
		logger.Error("Error getting cart item by cart & product: %v", err)
		return CartItem{}, err
	}

	return mapCartItem(row), nil
}

func (r *repository) ListByCart(ctx context.Context, cartID string) ([]CartItem, error) {
	var uuidCart pgtype.UUID
	if err := uuidCart.Scan(cartID); err != nil {
		return nil, err
	}

	rows, err := r.q.ListCartItems(ctx, uuidCart)
	if err != nil {
		logger.Error("Error listing cart items: %v", err)
		return nil, err
	}

	items := make([]CartItem, len(rows))
	for i, row := range rows {
		items[i] = mapCartItem(row)
	}

	return items, nil
}

func (r *repository) UpdateQuantity(ctx context.Context, id string, quantity int32) (CartItem, error) {
	var uuidID pgtype.UUID
	if err := uuidID.Scan(id); err != nil {
		return CartItem{}, err
	}

	params := sqlc.UpdateCartItemQuantityParams{
		ID:       uuidID,
		Quantity: quantity,
	}

	row, err := r.q.UpdateCartItemQuantity(ctx, params)
	if err != nil {
		logger.Error("Error updating cart item quantity: %v", err)
		return CartItem{}, err
	}

	return mapCartItem(row), nil
}

// Delete single item
func (r *repository) Delete(ctx context.Context, id string) error {
	var uuidID pgtype.UUID
	if err := uuidID.Scan(id); err != nil {
		return err
	}

	if err := r.q.DeleteCartItem(ctx, uuidID); err != nil {
		logger.Error("Error deleting cart item: %v", err)
		return err
	}
	return nil
}

// Delete all items for a cart
func (r *repository) DeleteByCart(ctx context.Context, cartID string) error {
	var uuidCart pgtype.UUID
	if err := uuidCart.Scan(cartID); err != nil {
		return err
	}

	if err := r.q.DeleteCartItemsByCart(ctx, uuidCart); err != nil {
		logger.Error("Error deleting cart items by cart: %v", err)
		return err
	}
	return nil
}

func (r *repository) GetCartByUserID(ctx context.Context, userID string) (cart.Cart, error) {
	var uuidUser pgtype.UUID
	if err := uuidUser.Scan(userID); err != nil {
		return cart.Cart{}, err
	}

	row, err := r.q.GetCartByUserID(ctx, uuidUser)
	if err != nil {
		logger.Error("Error getting cart by user ID: %v", err)
		return cart.Cart{}, err
	}

	return cart.MapCart(row), nil
}

func mapCartItem(row sqlc.CartItem) CartItem {
	return CartItem{
		ID:        uuid.UUID(row.ID.Bytes),
		CartID:    uuid.UUID(row.CartID.Bytes),
		ProductID: uuid.UUID(row.ProductID.Bytes),
		Quantity:  row.Quantity,
		CreatedAt: row.CreatedAt.Time,
	}
}
