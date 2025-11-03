-- name: CreateUserOTP :one
INSERT INTO user_otps (user_id, code, type, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserOTP :one
SELECT * FROM user_otps
WHERE user_id = $1 AND type = $2
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkUserOTPUsed :exec
UPDATE user_otps SET used = TRUE WHERE id = $1;

-- name: DeleteExpiredOTPs :exec
DELETE FROM user_otps WHERE expires_at < NOW();
