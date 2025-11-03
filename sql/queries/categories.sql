-- name: CreateCategory :one
INSERT INTO categories (
    name,
    slug,
    parent_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: GetCategoryBySlug :one
SELECT * FROM categories
WHERE slug = $1 LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: CountCategories :one
SELECT COUNT(*) FROM categories;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2,
    slug = $3,
    parent_id = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;