-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) AS total_count FROM users;


-- name: UpdateUser :one
UPDATE users
SET first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    profile_picture_url = COALESCE($4, profile_picture_url),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
