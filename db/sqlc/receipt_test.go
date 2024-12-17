package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bsmorton1983/receipt_processor/util"
	"github.com/stretchr/testify/require"
)

// TestCreateReceipt test receipt db create
func TestCreateReceipt(t *testing.T) {
	createTestReceipt(t, true)
}

// TestGetReceipt test receipt db get
func TestGetReceipt(t *testing.T) {
	receipt1 := createTestReceipt(t, true)
	receipt2, err := testQueries.GetReceipt(context.Background(), receipt1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, receipt2)

	require.Equal(t, receipt1.ID, receipt2.ID)
	require.Equal(t, receipt1.Retailer, receipt2.Retailer)
	require.Equal(t, receipt1.PurchaseDate, receipt2.PurchaseDate)
	require.Equal(t, receipt1.PurchaseTime, receipt2.PurchaseTime)

	require.WithinDuration(t, receipt1.CreationTime, receipt2.CreationTime, time.Second)
}

// TestDeleteReceipt test receipt db delete
func TestDeleteReceipt(t *testing.T) {
	receipt1 := createTestReceipt(t, false)
	err := testQueries.DeleteReceipt(context.Background(), receipt1.ID)
	require.NoError(t, err)

	receipt2, err := testQueries.GetReceipt(context.Background(), receipt1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, receipt2)
}

// TestListReceipts test receipt db list
func TestListReceipts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestReceipt(t, true)
	}

	arg := ListReceiptsParams{
		Limit:  5,
		Offset: 5,
	}

	receipts, err := testQueries.ListReceipts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, receipts, 5)

	for _, receipt := range receipts {
		require.NotEmpty(t, receipt)
	}
}

// createTestReceipt create receipt in db and verify results
func createTestReceipt(t *testing.T, add_to_delete_queue bool) Receipt {
	arg := CreateReceiptParams{
		Retailer:     util.RandomRetailer(),
		PurchaseDate: util.CurrentDate(),
		PurchaseTime: util.CurrentTime(),
	}

	receipt, err := testQueries.CreateReceipt(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, receipt)

	require.Equal(t, arg.Retailer, receipt.Retailer)
	require.Equal(t, arg.PurchaseDate, receipt.PurchaseDate)
	require.Equal(t, arg.PurchaseTime, receipt.PurchaseTime)

	require.NotZero(t, receipt.ID)
	require.NotZero(t, receipt.CreationTime)

	if add_to_delete_queue {
		receipts_to_delete = append(receipt_items_to_delete, receipt.ID)
	}

	t.Cleanup(func() {
		clear_receipts()
	})

	return receipt
}
