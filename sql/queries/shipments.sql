-- name: CreateShipment :one
INSERT INTO shipments (
    order_id, carrier, tracking_number, status, shipped_at, delivered_at
)
VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetShipment :one
SELECT * FROM shipments
WHERE id = $1;

-- name: ListShipmentsByOrder :many
SELECT * FROM shipments
WHERE order_id = $1
ORDER BY created_at DESC;

-- name: UpdateShipmentStatus :one
UPDATE shipments
SET status = $2,
    shipped_at = $3,
    delivered_at = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteShipment :exec
DELETE FROM shipments
WHERE id = $1;
