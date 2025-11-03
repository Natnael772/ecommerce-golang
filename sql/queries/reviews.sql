-- name: CreateReview :one
INSERT INTO reviews (
    product_id, user_id, rating, comment
)
VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetReview :one
SELECT * FROM reviews
WHERE id = $1;

-- name: GetReviewsByUser :many
SELECT * FROM reviews
WHERE user_id = $1;

-- name: GetReviewByUserAndProduct :one
SELECT * FROM reviews
WHERE user_id = $1 AND product_id = $2
LIMIT 1;

-- name: ListReviewsByProduct :many
SELECT * FROM reviews
WHERE product_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountReviewsByProduct :one
SELECT COUNT(*) FROM reviews
WHERE product_id = $1;

-- name: UpdateReview :one
UPDATE reviews
SET rating = $2,
    comment = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1;
