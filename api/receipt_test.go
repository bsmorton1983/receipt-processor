package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	mockdb "github.com/bsmorton1983/receipt_processor/db/mock"
	db "github.com/bsmorton1983/receipt_processor/db/sqlc"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// mock process receipt request for testing
var testProcessReceiptRequest = processReceiptRequest{
	Retailer:     "Target",
	PurchaseDate: "2022-01-01",
	PurchaseTime: "13:01",
	Items: []receiptItemRequest{
		{
			ShortDescription: "Mountain Dew 12PK",
			Price:            "6.49",
		},
		{
			ShortDescription: "Emils Cheese Pizza",
			Price:            "12.25",
		},
		{
			ShortDescription: "Knorr Creamy Chicken",
			Price:            "1.26",
		},
		{
			ShortDescription: "Doritos Nacho Cheese",
			Price:            "3.35",
		},
		{
			ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
			Price:            "12.00",
		},
	},
	Total: "35.35",
}

type invalidIdResponse struct {
	Error string
}

type testCase struct {
	name          string
	receiptID     string
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}

// testReceipt returns a mock receipt
func testReceipt() db.Receipt {
	return db.Receipt{
		ID:           uuid.New(),
		Retailer:     testProcessReceiptRequest.Retailer,
		PurchaseDate: testProcessReceiptRequest.PurchaseDate,
		PurchaseTime: testProcessReceiptRequest.PurchaseTime,
	}
}

// testReceiptItem returns a mock receipt item
func testReceiptItem(req receiptItemRequest, receipt db.Receipt, price float64) db.ReceiptItem {
	return db.ReceiptItem{
		ID:               uuid.New(),
		ReceiptID:        receipt.ID,
		ShortDescription: strings.TrimSpace(req.ShortDescription),
		Price:            price,
	}
}

// TestProcessReceiptAPI tests the /receipts/process api call
func TestProcessReceiptAPI(t *testing.T) {
	receipt := testReceipt()

	testCases := []testCase{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateReceiptParams{
					Retailer:     strings.TrimSpace(testProcessReceiptRequest.Retailer),
					PurchaseDate: strings.TrimSpace(testProcessReceiptRequest.PurchaseDate),
					PurchaseTime: strings.TrimSpace(testProcessReceiptRequest.PurchaseTime),
				}
				store.EXPECT().CreateReceipt(gomock.Any(), gomock.Eq(arg)).Times(1).Return(receipt, nil)
				for _, req := range testProcessReceiptRequest.Items {
					price, err := strconv.ParseFloat(req.Price, 64)
					require.NoError(t, err)
					receipt_item := testReceiptItem(req, receipt, price)
					arg := db.CreateReceiptItemParams{
						ReceiptID:        receipt.ID,
						ShortDescription: strings.TrimSpace(req.ShortDescription),
						Price:            price,
					}
					store.EXPECT().CreateReceiptItem(gomock.Any(), gomock.Eq(arg)).Times(1).Return(receipt_item, nil)
				}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchProcessReceiptResponce(t, recorder.Body, processReceiptResponse{
					ID: receipt.ID,
				})
			},
		},
		{
			name: "InternalError(Receipt DB)",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateReceiptParams{
					Retailer:     strings.TrimSpace(testProcessReceiptRequest.Retailer),
					PurchaseDate: strings.TrimSpace(testProcessReceiptRequest.PurchaseDate),
					PurchaseTime: strings.TrimSpace(testProcessReceiptRequest.PurchaseTime),
				}
				store.EXPECT().CreateReceipt(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Receipt{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalError(Receipt Item DB)",
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := db.CreateReceiptParams{
					Retailer:     strings.TrimSpace(testProcessReceiptRequest.Retailer),
					PurchaseDate: strings.TrimSpace(testProcessReceiptRequest.PurchaseDate),
					PurchaseTime: strings.TrimSpace(testProcessReceiptRequest.PurchaseTime),
				}
				store.EXPECT().CreateReceipt(gomock.Any(), gomock.Eq(arg1)).Times(1).Return(receipt, nil)
				req := testProcessReceiptRequest.Items[0]
				price, err := strconv.ParseFloat(req.Price, 64)
				require.NoError(t, err)
				arg2 := db.CreateReceiptItemParams{
					ReceiptID:        receipt.ID,
					ShortDescription: strings.TrimSpace(req.ShortDescription),
					Price:            price,
				}
				store.EXPECT().CreateReceiptItem(gomock.Any(), gomock.Eq(arg2)).Times(1).Return(db.ReceiptItem{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	request_json, err := json.Marshal(testProcessReceiptRequest)
	require.NoError(t, err)
	excecuteTestCases(t, testCases, request_json, "/receipts/process", http.MethodPost)
}

// TestProcessReceiptAPI tests the /receipts/:id/points api call
func TestGetPointsAPI(t *testing.T) {
	receipt := testReceipt()

	testCases := []testCase{
		{
			name:      "OK",
			receiptID: receipt.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetReceipt(gomock.Any(), gomock.Eq(receipt.ID)).Times(1).Return(receipt, nil)
				receipt_items := []db.ReceiptItem{}
				for _, req := range testProcessReceiptRequest.Items {
					price, err := strconv.ParseFloat(req.Price, 64)
					require.NoError(t, err)
					receipt_items = append(receipt_items, testReceiptItem(req, receipt, price))
				}
				store.EXPECT().ListReceiptItems(gomock.Any(), gomock.Any()).Times(1).Return(receipt_items, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGetPointsResponce(t, recorder.Body, getPointsResponse{
					Points: 28,
				})
			},
		},
		{
			name:      "InternalError(Receipt DB)",
			receiptID: receipt.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetReceipt(gomock.Any(), gomock.Eq(receipt.ID)).Times(1).Return(db.Receipt{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InternalError(Receipt Item DB)",
			receiptID: receipt.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetReceipt(gomock.Any(), gomock.Eq(receipt.ID)).Times(1).Return(receipt, nil)
				store.EXPECT().ListReceiptItems(gomock.Any(), gomock.Any()).Times(1).Return([]db.ReceiptItem{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			receiptID: receipt.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetReceipt(gomock.Any(), gomock.Eq(receipt.ID)).Times(1).Return(db.Receipt{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			receiptID: "BAD-UUID",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetReceipt(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				fmt.Println(recorder.Body)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				requireBodyMatchInvalidIdResponce(t, recorder.Body, invalidIdResponse{
					Error: "invalid UUID length: 8",
				})
			},
		},
	}

	excecuteTestCases(t, testCases, nil, "/receipts/id/points", http.MethodGet)
}

// helper function to check http responce for process receipt
func requireBodyMatchProcessReceiptResponce(t *testing.T, body *bytes.Buffer, response processReceiptResponse) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotResponse processReceiptResponse
	err = json.Unmarshal(data, &gotResponse)
	require.NoError(t, err)
	require.Equal(t, response, gotResponse)
}

// helper function to check http responce for get points
func requireBodyMatchGetPointsResponce(t *testing.T, body *bytes.Buffer, response getPointsResponse) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotResponse getPointsResponse
	err = json.Unmarshal(data, &gotResponse)
	require.NoError(t, err)
	require.Equal(t, response, gotResponse)
}

// helper function to check http responce for invalid id
func requireBodyMatchInvalidIdResponce(t *testing.T, body *bytes.Buffer, response invalidIdResponse) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotResponse invalidIdResponse
	err = json.Unmarshal(data, &gotResponse)
	require.NoError(t, err)
	require.Equal(t, response, gotResponse)
}

// helper function to execute batch receipt tests
func excecuteTestCases(t *testing.T, testCases []testCase, request_json []byte, url, method string) {
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			new_url := url
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			if tc.receiptID != "" {
				new_url = strings.Replace(new_url, "id", tc.receiptID, 1)
			}

			var request *http.Request
			var err error
			if request_json != nil {
				request, err = http.NewRequest(method, new_url, bytes.NewBuffer(request_json))
			} else {
				request, err = http.NewRequest(method, new_url, nil)
			}
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
