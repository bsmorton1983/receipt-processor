package db

import (
	"context"
	"testing"

	"github.com/bsmorton1983/receipt_processor/db/util"
	"github.com/stretchr/testify/require"
)

func createTestReceiptItem(t *testing.T, receipt Receipt) ReceiptItem {
	arg := CreateReceiptItemParams{
		ReceiptID:        receipt.ID,
		ShortDescription: util.RandomDescription(),
		Price:            util.RandomPrice(),
	}

	receipt_item, err := testQueries.CreateReceiptItem(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, receipt_item)

	require.Equal(t, arg.ReceiptID, receipt_item.ReceiptID)
	require.Equal(t, arg.ShortDescription, receipt_item.ShortDescription)
	require.Equal(t, arg.Price, receipt_item.Price)

	require.NotZero(t, receipt_item.ID)

	return receipt_item
}

func TestCreateReceiptItem(t *testing.T) {
	createTestReceiptItem(t, createTestReceipt(t))
}

func TestGetReceiptItem(t *testing.T) {
	receipt_item1 := createTestReceiptItem(t, createTestReceipt(t))
	receipt_item2, err := testQueries.GetReceiptItem(context.Background(), receipt_item1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, receipt_item2)

	require.Equal(t, receipt_item1.ID, receipt_item2.ID)
	require.Equal(t, receipt_item1.ReceiptID, receipt_item2.ReceiptID)
	require.Equal(t, receipt_item1.ShortDescription, receipt_item2.ShortDescription)
	require.Equal(t, receipt_item1.Price, receipt_item2.Price)
}

func TestListReceiptItems(t *testing.T) {
	receipt := createTestReceipt(t)
	for i := 0; i < 10; i++ {
		createTestReceiptItem(t, receipt)
	}

	arg := ListReceiptItemsParams{
		ReceiptID: receipt.ID,
		Limit:     5,
		Offset:    5,
	}

	receipt_items, err := testQueries.ListReceiptItems(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, receipt_items, 5)

	for _, receipt_item := range receipt_items {
		require.NotEmpty(t, receipt_item)
		require.Equal(t, arg.ReceiptID, receipt_item.ReceiptID)
	}
}
