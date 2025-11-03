-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    order_number,
    subtotal_cents,
    discount_cents,
    tax_cents,
    shipping_cents,
    total_cents,
    final_cents,
    -- currency,
    shipping_info,
    notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
 RETURNING *;

-- -- name: GetOrderByID :one
-- SELECT * FROM orders
-- WHERE id = $1 LIMIT 1;

-- -- name: GetOrdersByUserID :many
-- SELECT * FROM orders
-- WHERE user_id = $1
-- ORDER BY created_at DESC;

-- -- name: GetOrders :many
-- SELECT * FROM orders
-- ORDER BY created_at DESC
-- LIMIT $1 OFFSET $2;

-- name: GetOrderWithItemsByID :one
SELECT 
    o.*,
    COALESCE(
        json_agg(to_jsonb(oi)) FILTER (WHERE oi.id IS NOT NULL), '[]'
    ) AS items
FROM orders o
LEFT JOIN order_items oi ON o.id = oi.order_id
WHERE o.id = $1
GROUP BY o.id;

-- name: GetOrdersWithItemsByUserID :many
SELECT 
    o.*,
    COALESCE(
        json_agg(to_jsonb(oi)) FILTER (WHERE oi.id IS NOT NULL), '[]'
    ) AS items
FROM orders o
LEFT JOIN order_items oi ON o.id = oi.order_id
WHERE o.user_id = $1
GROUP BY o.id
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetOrdersWithItems :many
SELECT 
    o.*,
    COALESCE(
        json_agg(to_jsonb(oi)) FILTER (WHERE oi.id IS NOT NULL), '[]'
    ) AS items
FROM orders o
LEFT JOIN order_items oi ON o.id = oi.order_id
GROUP BY o.id
ORDER BY o.created_at DESC
LIMIT $1 OFFSET $2;


-- name: CountOrders :one
SELECT COUNT(*) AS total_count FROM orders;

-- name: CountOrdersByUser :one
SELECT COUNT(*) AS total_count FROM orders
WHERE user_id = $1;

-- name: UpdateOrderStatus :one
UPDATE orders
SET 
    status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;