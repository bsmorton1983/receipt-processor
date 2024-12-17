// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type Receipt struct {
	ID           uuid.UUID `json:"id"`
	Retailer     string    `json:"retailer"`
	PurchaseDate string    `json:"purchase_date"`
	PurchaseTime string    `json:"purchase_time"`
	CreationTime time.Time `json:"creation_time"`
}

type ReceiptItem struct {
	ID               uuid.UUID `json:"id"`
	ReceiptID        uuid.UUID `json:"receipt_id"`
	ShortDescription string    `json:"short_description"`
	Price            float64   `json:"price"`
	CreationTime     time.Time `json:"creation_time"`
}
