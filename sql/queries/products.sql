-- name: CreateProduct :one
INSERT INTO products (
    sku,
    name,
    description,
    category_id,
    price_cents,
    currency,
    attributes,
    main_image_url,
    images,
    discount_percent,
    discount_valid_until
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: GetProductBySKU :one
SELECT * FROM products
WHERE sku = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
WHERE is_deleted = FALSE
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: CountProducts :one
SELECT COUNT(*) FROM products WHERE is_deleted = FALSE;

-- name: UpdateProductPrice :one
UPDATE products
SET price_cents = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    category_id = COALESCE($4, category_id),
    price_cents = COALESCE($5, price_cents),
    currency = COALESCE($6, currency),
    attributes = COALESCE($7, attributes),
    main_image_url = COALESCE($8, main_image_url),
    images = COALESCE($9, images),
    discount_percent = COALESCE($10, discount_percent),
    discount_valid_until = COALESCE($11, discount_valid_until),
    updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: DeleteProduct :exec
UPDATE products
SET is_deleted = TRUE,
    updated_at = NOW()
WHERE id = $1;