-- name: CreateAddress :one
INSERT INTO addresses (
    user_id, label, line1, line2, city, state, postal_code, country, is_default
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetAddress :one
SELECT * FROM addresses
WHERE id = $1 AND is_deleted = false;

-- name: GetAddressesByUser :many
SELECT * FROM addresses
WHERE user_id = $1 AND is_deleted = false
ORDER BY is_default DESC, created_at DESC;

-- name: UpdateAddress :one
UPDATE addresses
SET label = COALESCE($2, label),
    line1 = COALESCE($3, line1),
    line2 = COALESCE($4, line2),
    city = COALESCE($5, city),
    state = COALESCE($6, state),
    postal_code = COALESCE($7, postal_code),
    country = COALESCE($8, country),
    is_default = COALESCE($9, is_default),
    updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: DeleteAddress :exec
UPDATE addresses
SET is_deleted=true,
    updated_at = NOW()
WHERE id = $1;
