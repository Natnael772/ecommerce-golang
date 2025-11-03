-- name: CreateCoupon :one
INSERT INTO coupons (
    code, description, discount_percent, valid_from, valid_until, max_uses, is_active
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetCouponByCode :one
SELECT * FROM coupons
WHERE code = $1;

-- name: GetCouponById :one
SELECT * FROM coupons
WHERE id = $1 LIMIT 1;

-- name: GetCoupons :many
SELECT * FROM coupons
LIMIT $1 OFFSET $2;

-- name: CountCoupons :one
SELECT COUNT(*) FROM coupons;

-- name: IncrementCouponUsage :one
UPDATE coupons
SET used_count = used_count + 1,
    updated_at = NOW()
WHERE code = $1
RETURNING *;

-- name: UpdateCoupon :one
UPDATE coupons
SET code = COALESCE($2, code),
    description = COALESCE($3, description),
    discount_percent = COALESCE($4, discount_percent),
    valid_from = COALESCE($5, valid_from),
    valid_until = COALESCE($6, valid_until),
    max_uses = COALESCE($7, max_uses),
    is_active = COALESCE($8, is_active),
    updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: DeleteCoupon :exec
DELETE FROM coupons
WHERE id = $1;
