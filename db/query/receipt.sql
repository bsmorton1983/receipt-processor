-- name: CreateReceipt :one
INSERT INTO receipts (
    retailer,
    purchase_date,
    purchase_time
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetReceipt :one
SELECT * FROM receipts
WHERE id = $1 LIMIT 1;

-- name: ListReceipts :many
SELECT * FROM receipts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteReceipt :exec
DELETE FROM receipts 
WHERE id = $1;