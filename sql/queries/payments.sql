-- name: CreatePayment :one
INSERT INTO payments (
    order_id,
    provider,
    provider_txn_id,
    amount_cents,
    currency,
    payment_method,
    status,
    details
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetPaymentByOrderID :one
SELECT * FROM payments
WHERE order_id = $1 LIMIT 1;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET status = $2, created_at = NOW()
WHERE id = $1
RETURNING *;