-- name: CreateCart :one
INSERT INTO carts (user_id, expires_at)
VALUES ($1, $2)
RETURNING *;

-- name: GetCartByID :one
SELECT *
FROM carts
WHERE id = $1
  AND (expires_at IS NULL OR expires_at > NOW())
LIMIT 1;

-- name: GetCartByUserID :one
SELECT *
FROM carts
WHERE user_id = $1
  AND (expires_at IS NULL OR expires_at > NOW())
ORDER BY created_at DESC
LIMIT 1;

-- name: ListCarts :many
SELECT *
FROM carts
WHERE expires_at IS NULL OR expires_at > NOW()
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateCart :one
UPDATE carts
SET 
    expires_at = COALESCE($2, expires_at),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCart :exec
DELETE FROM carts
WHERE id = $1;

-- name: DeleteExpiredCarts :exec
DELETE FROM carts
WHERE expires_at IS NOT NULL
  AND expires_at < NOW();
