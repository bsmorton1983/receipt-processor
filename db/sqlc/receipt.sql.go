// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: receipt.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createReceipt = `-- name: CreateReceipt :one
INSERT INTO receipts (
    retailer,
    purchase_date,
    purchase_time
) VALUES (
    $1, $2, $3
) RETURNING id, retailer, purchase_date, purchase_time, creation_time
`

type CreateReceiptParams struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchase_date"`
	PurchaseTime string `json:"purchase_time"`
}

func (q *Queries) CreateReceipt(ctx context.Context, arg CreateReceiptParams) (Receipt, error) {
	row := q.db.QueryRowContext(ctx, createReceipt, arg.Retailer, arg.PurchaseDate, arg.PurchaseTime)
	var i Receipt
	err := row.Scan(
		&i.ID,
		&i.Retailer,
		&i.PurchaseDate,
		&i.PurchaseTime,
		&i.CreationTime,
	)
	return i, err
}

const deleteReceipt = `-- name: DeleteReceipt :exec
DELETE FROM receipts 
WHERE id = $1
`

func (q *Queries) DeleteReceipt(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteReceipt, id)
	return err
}

const getReceipt = `-- name: GetReceipt :one
SELECT id, retailer, purchase_date, purchase_time, creation_time FROM receipts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetReceipt(ctx context.Context, id uuid.UUID) (Receipt, error) {
	row := q.db.QueryRowContext(ctx, getReceipt, id)
	var i Receipt
	err := row.Scan(
		&i.ID,
		&i.Retailer,
		&i.PurchaseDate,
		&i.PurchaseTime,
		&i.CreationTime,
	)
	return i, err
}

const listReceipts = `-- name: ListReceipts :many
SELECT id, retailer, purchase_date, purchase_time, creation_time FROM receipts
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListReceiptsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListReceipts(ctx context.Context, arg ListReceiptsParams) ([]Receipt, error) {
	rows, err := q.db.QueryContext(ctx, listReceipts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Receipt{}
	for rows.Next() {
		var i Receipt
		if err := rows.Scan(
			&i.ID,
			&i.Retailer,
			&i.PurchaseDate,
			&i.PurchaseTime,
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
