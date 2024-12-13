package api

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	db "github.com/bsmorton1983/receipt_processor/db/sqlc"
	"github.com/gin-gonic/gin"
)

type receiptItemRequest struct {
	ShortDescription string `json:"shortDescription" binding:"required`
	Price            string `json:"price" binding:"required`
}

type processReceiptRequest struct {
	Retailer     string               `json:"retailer" binding:"required`
	PurchaseDate string               `json:"purchaseDate" binding:"required`
	PurchaseTime string               `json:"purchaseTime" binding:"required`
	Items        []receiptItemRequest `json:"items" binding:"required`
	Total        string               `json:"total" binding:"required`
}

type processReceiptResponse struct {
	Id string
}

func (server *Server) processReceipt(ctx *gin.Context) {
	var req processReceiptRequest
	var res processReceiptResponse
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateReceiptParams{
		Retailer:     strings.TrimSpace(req.Retailer),
		PurchaseDate: strings.TrimSpace(req.PurchaseDate),
		PurchaseTime: strings.TrimSpace(req.PurchaseTime),
	}

	receipt, err := server.store.CreateReceipt(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, receipt_item := range req.Items {
		price, err := strconv.ParseFloat(receipt_item.Price, 64)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		arg := db.CreateReceiptItemParams{
			ReceiptID:        receipt.ID,
			ShortDescription: strings.TrimSpace(receipt_item.ShortDescription),
			Price:            price,
		}
		_, err = server.store.CreateReceiptItem(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	res.Id = strconv.FormatInt(receipt.ID, 10)
	ctx.JSON(http.StatusOK, res)
}

type getPointsRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type getPointsResponse struct {
	Points int64
}

func (server *Server) getPoints(ctx *gin.Context) {
	var req getPointsRequest
	var res getPointsResponse
	points := int64(0)
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	receipt, err := server.store.GetReceipt(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, c := range receipt.Retailer {
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			points += 1
		}
	}

	if p, err := calculateDatePoints(receipt.PurchaseDate); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	} else {
		points += p
	}

	if p, err := calculateTimePoints(receipt.PurchaseTime); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	} else {
		points += p
	}

	arg := db.ListReceiptItemsParams{
		ReceiptID: receipt.ID,
		Limit:     500,
		Offset:    0,
	}

	receipt_items, err := server.store.ListReceiptItems(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	points += calculateItemPoints(receipt_items)
	res.Points = points
	ctx.JSON(http.StatusOK, res)
}

func calculateDatePoints(purchaseDate string) (int64, error) {
	points := int64(0)
	date, err := strconv.ParseInt(purchaseDate[len(purchaseDate)-2:], 10, 64)
	if err != nil {
		return points, err
	}
	if date%2 != 0 {
		points = int64(6)
	}
	return points, err
}

func calculateTimePoints(purchaseTime string) (int64, error) {
	points := int64(0)
	hour, err := strconv.ParseInt(purchaseTime[0:2], 10, 64)
	if err != nil {
		return points, err
	}
	if hour >= 14 && hour < 16 {
		points = 10
	}
	return points, err
}

func calculateItemPoints(receipt_items []db.ReceiptItem) int64 {
	total := float64(0)
	points := int64(0)

	for i, item := range receipt_items {
		total += item.Price
		if (i+1)%2 == 0 {
			points += 5
		}
		if len(item.ShortDescription)%3 == 0 {
			points += int64(math.Ceil(item.Price * 0.2))
		}
	}

	total_str := fmt.Sprintf("%.2f", total)
	cents := total_str[len(total_str)-2:]

	if cents == "00" {
		points += 75
	} else if cents == "25" || cents == "50" || cents == "75" {
		points += 25
	}
	return points
}
