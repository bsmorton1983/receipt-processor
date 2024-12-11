-- name: CreateReceiptItem :one
INSERT INTO receipt_items (
    receipt_id,
    short_description,
    price
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetReceiptItem :one
SELECT * FROM receipt_items
WHERE id = $1 LIMIT 1;

-- name: ListReceiptItems :many
SELECT * FROM receipt_items
WHERE receipt_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;