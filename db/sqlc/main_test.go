package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/bsmorton1983/receipt_processor/db/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

var receipts_to_delete = []int64{}
var receipt_items_to_delete = []int64{}

func clear_receipts() {
	for _, id := range receipts_to_delete {
		err := testQueries.DeleteReceipt(context.Background(), id)
		if err != nil {
			fmt.Println("Error deleting receipt:", err)
		}
	}
}

func clear_receipt_items() {
	for _, id := range receipt_items_to_delete {
		err := testQueries.DeleteReceiptItem(context.Background(), id)
		if err != nil {
			fmt.Println("Error deleting receipt item:", err)
		}
	}
}

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
