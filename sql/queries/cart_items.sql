-- name: AddCartItem :one
INSERT INTO cart_items (cart_id, product_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (cart_id, product_id)
DO UPDATE SET
    quantity = EXCLUDED.quantity,
    created_at = NOW()
RETURNING *;


-- name: AddCartItems :many
INSERT INTO cart_items (id, cart_id, product_id, quantity, created_at, updated_at)
SELECT
    gen_random_uuid() AS id,
    $1 AS cart_id,
    p.product_id,
    p.quantity,
    NOW() AS created_at,
    NOW() AS updated_at
FROM (
    SELECT unnest($2::uuid[]) AS product_id,
           unnest($3::int[]) AS quantity
) AS p(product_id, quantity)
ON CONFLICT (cart_id, product_id)
DO UPDATE SET
    quantity = EXCLUDED.quantity,  
    updated_at = NOW()
RETURNING *;



-- name: GetCartItem :one
SELECT *
FROM cart_items
WHERE id = $1
LIMIT 1;

-- name: GetCartItemByProduct :one
SELECT *
FROM cart_items
WHERE cart_id = $1
  AND product_id = $2
LIMIT 1;

-- name: GetCartItemsByUserID :many
SELECT ci.*
FROM cart_items ci
JOIN carts c ON ci.cart_id = c.id
WHERE c.user_id = $1
ORDER BY ci.created_at ASC;

-- name: ListCartItems :many
SELECT *
FROM cart_items
WHERE cart_id = $1
ORDER BY created_at ASC;

-- name: UpdateCartItemQuantity :one
UPDATE cart_items
SET 
    quantity = $2,
    created_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCartItem :exec
DELETE FROM cart_items
WHERE id = $1;

-- name: DeleteCartItemsByCart :exec
DELETE FROM cart_items
WHERE cart_id = $1;
