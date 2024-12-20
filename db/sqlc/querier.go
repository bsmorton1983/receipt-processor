// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateReceipt(ctx context.Context, arg CreateReceiptParams) (Receipt, error)
	CreateReceiptItem(ctx context.Context, arg CreateReceiptItemParams) (ReceiptItem, error)
	DeleteReceipt(ctx context.Context, id uuid.UUID) error
	DeleteReceiptItem(ctx context.Context, id uuid.UUID) error
	GetReceipt(ctx context.Context, id uuid.UUID) (Receipt, error)
	GetReceiptItem(ctx context.Context, id uuid.UUID) (ReceiptItem, error)
	ListReceiptItems(ctx context.Context, arg ListReceiptItemsParams) ([]ReceiptItem, error)
	ListReceipts(ctx context.Context, arg ListReceiptsParams) ([]Receipt, error)
}

var _ Querier = (*Queries)(nil)
