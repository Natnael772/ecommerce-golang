-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id,
    product_id,
    sku,
    name,
    qty,
    unit_price_cents,
    total_price_cents
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetOrderItemsByOrderID :many
SELECT * FROM order_items
WHERE order_id = $1;