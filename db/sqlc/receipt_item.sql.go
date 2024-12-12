// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: receipt_item.sql

package db

import (
	"context"
)

const createReceiptItem = `-- name: CreateReceiptItem :one
INSERT INTO receipt_items (
    receipt_id,
    short_description,
    price
) VALUES (
    $1, $2, $3
) RETURNING id, receipt_id, short_description, price, creation_time
`

type CreateReceiptItemParams struct {
	ReceiptID        int64  `json:"receipt_id"`
	ShortDescription string `json:"short_description"`
	Price            int64  `json:"price"`
}

func (q *Queries) CreateReceiptItem(ctx context.Context, arg CreateReceiptItemParams) (ReceiptItem, error) {
	row := q.db.QueryRowContext(ctx, createReceiptItem, arg.ReceiptID, arg.ShortDescription, arg.Price)
	var i ReceiptItem
	err := row.Scan(
		&i.ID,
		&i.ReceiptID,
		&i.ShortDescription,
		&i.Price,
		&i.CreationTime,
	)
	return i, err
}

const deleteReceiptItem = `-- name: DeleteReceiptItem :exec
DELETE FROM receipt_items
WHERE id = $1
`

func (q *Queries) DeleteReceiptItem(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteReceiptItem, id)
	return err
}

const getReceiptItem = `-- name: GetReceiptItem :one
SELECT id, receipt_id, short_description, price, creation_time FROM receipt_items
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetReceiptItem(ctx context.Context, id int64) (ReceiptItem, error) {
	row := q.db.QueryRowContext(ctx, getReceiptItem, id)
	var i ReceiptItem
	err := row.Scan(
		&i.ID,
		&i.ReceiptID,
		&i.ShortDescription,
		&i.Price,
		&i.CreationTime,
	)
	return i, err
}

const listReceiptItems = `-- name: ListReceiptItems :many
SELECT id, receipt_id, short_description, price, creation_time FROM receipt_items
WHERE receipt_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListReceiptItemsParams struct {
	ReceiptID int64 `json:"receipt_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListReceiptItems(ctx context.Context, arg ListReceiptItemsParams) ([]ReceiptItem, error) {
	rows, err := q.db.QueryContext(ctx, listReceiptItems, arg.ReceiptID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ReceiptItem{}
	for rows.Next() {
		var i ReceiptItem
		if err := rows.Scan(
			&i.ID,
			&i.ReceiptID,
			&i.ShortDescription,
			&i.Price,
			&i.CreationTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
