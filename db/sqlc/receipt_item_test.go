package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bsmorton1983/receipt_processor/db/util"
	"github.com/stretchr/testify/require"
)

func createTestReceiptItem(t *testing.T, receipt Receipt, add_to_delete_queue bool) ReceiptItem {
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
	require.NotZero(t, receipt_item.CreationTime)

	if add_to_delete_queue {
		receipt_items_to_delete = append(receipt_items_to_delete, receipt_item.ID)
	}

	t.Cleanup(func() {
		clear_receipt_items()
	})

	return receipt_item
}

func TestCreateReceiptItem(t *testing.T) {
	createTestReceiptItem(t, createTestReceipt(t, true), true)
}

func TestGetReceiptItem(t *testing.T) {
	receipt_item1 := createTestReceiptItem(t, createTestReceipt(t, true), true)
	receipt_item2, err := testQueries.GetReceiptItem(context.Background(), receipt_item1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, receipt_item2)

	require.Equal(t, receipt_item1.ID, receipt_item2.ID)
	require.Equal(t, receipt_item1.ReceiptID, receipt_item2.ReceiptID)
	require.Equal(t, receipt_item1.ShortDescription, receipt_item2.ShortDescription)
	require.Equal(t, receipt_item1.Price, receipt_item2.Price)

	require.WithinDuration(t, receipt_item1.CreationTime, receipt_item2.CreationTime, time.Second)
}

func TestDeleteReceiptItem(t *testing.T) {
	receipt := createTestReceipt(t, true)
	receipt_item1 := createTestReceiptItem(t, receipt, false)
	err := testQueries.DeleteReceiptItem(context.Background(), receipt_item1.ID)
	require.NoError(t, err)

	receipt_item2, err := testQueries.GetReceiptItem(context.Background(), receipt_item1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, receipt_item2)
}

func TestListReceiptItems(t *testing.T) {
	receipt := createTestReceipt(t, true)
	for i := 0; i < 10; i++ {
		createTestReceiptItem(t, receipt, true)
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
