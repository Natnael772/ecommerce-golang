-- name: CreateInventory :one
INSERT INTO inventory (
    product_id,
    stock,
    reserved
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetInventoryByProductID :one
SELECT * FROM inventory
WHERE product_id = $1 LIMIT 1;

-- name: UpdateInventoryStock :one
UPDATE inventory
SET 
    stock = $2,
    updated_at = NOW()
WHERE product_id = $1
RETURNING *;

-- name: DeleteInventory :exec
DELETE FROM inventory
WHERE product_id = $1;