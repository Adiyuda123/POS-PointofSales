package handler

import (
	"POS-PointofSales/features/products"
	"POS-PointofSales/features/transactions"
	"POS-PointofSales/features/users"
	"POS-PointofSales/helper"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type transactionController struct {
	service     transactions.UseCase
	productRepo products.Repository
	userRepo    users.Repository
}

func New(service transactions.UseCase, services products.Repository, servicess users.Repository) transactions.Handler {
	return &transactionController{
		service:     service,
		productRepo: services,
		userRepo:    servicess,
	}
}

// GetHistoryTransactionHandler implements transactions.Handler.
func (tc *transactionController) GetHistoryTransactionHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID == 0 {
			c.Logger().Error("decode token is blank")
			return c.JSON(http.StatusBadRequest, "jwt invalid")
		}

		search := c.QueryParam("search")

		limitStr := c.QueryParam("limit")
		offsetStr := c.QueryParam("offset")

		limit := 10
		if limitStr != "" {
			limitInt, err := strconv.Atoi(limitStr)
			if err != nil {
				c.Logger().Errorf("limit is not a number: %s", limitStr)
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid limit value", nil))
			}
			limit = limitInt
		}

		offset := 0
		if offsetStr != "" {
			offsetInt, err := strconv.Atoi(offsetStr)
			if err != nil {
				c.Logger().Errorf("offset is not a number: %s", offsetStr)
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid offset value", nil))
			}
			offset = offsetInt
		}

		fromDateStr := c.QueryParam("from_date")
		toDateStr := c.QueryParam("to_date")

		fromDate := time.Time{}
		if fromDateStr != "" {
			var err error
			fromDate, err = time.Parse("2006-01-02", fromDateStr)
			if err != nil {
				c.Logger().Errorf("from_date is not a valid date: %s", fromDateStr)
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid from_date value", nil))
			}
		}

		toDate := time.Now()
		if toDateStr != "" {
			var err error
			toDate, err = time.Parse("2006-01-02", toDateStr)
			if err != nil {
				c.Logger().Errorf("to_date is not a valid date: %s", toDateStr)
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid to_date value", nil))
			}
			toDate = toDate.AddDate(0, 0, 1)
		}

		data, totalData, err := tc.service.GetHistoryTransaction(userID, limit, offset, search, fromDate, toDate)
		if err != nil {
			c.Logger().Error("error occurs when calling GetHistoryTransaction Logic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Server Error", nil))
		}

		dataResponse := CoreToGetAllTransactionDetailsResponse(data)
		pagination := helper.Pagination(limit, offset, totalData)

		return c.JSON(helper.ReponseFormatWithMeta(http.StatusOK, "Successfully retrieved history transaction data", dataResponse, pagination))
	}
}

// AddPayments implements transactions.Handler.
func (tc *transactionController) AddPayments() echo.HandlerFunc {
	return func(c echo.Context) error {
		var inputRequest TransactionRequest
		if err := c.Bind(&inputRequest); err != nil {
			c.Logger().Error("error on bind payment input", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		items, err := tc.service.GetItemByOrderId(inputRequest.ReferenceID)
		if err != nil {
			c.Logger().Error("error get product data", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to get products", nil))
		}
		timestamp := time.Now().Unix()
		externalID := fmt.Sprintf("POSApp-Customer:%s-Timestamp:%d", items.Customer, timestamp)
		callbackURL := "https://www.pos-callback.com"
		currency := "IDR"
		types := "DYNAMIC"
		created := time.Now()
		expiresAt := created.Add(24 * time.Hour).Format("2006-01-02T15:04:05.999999Z")

		requestBody := TransactionRequest{
			ExternalID:  externalID,
			ReferenceID: inputRequest.ReferenceID,
			Type:        types,
			Currency:    currency,
			Amount:      float64(items.SubTotal),
			ExpiresAt:   expiresAt,
			CallbackURL: callbackURL,
		}

		qrCodeResponse, err := helper.SendXenditRequest(requestBody)
		if err != nil {
			c.Logger().Error("error generating QR code:", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to generate QR code", nil))
		}

		newTransaction := transactions.Core{
			ID:          qrCodeResponse.ID,
			ExternalID:  externalID,
			OrderID:     inputRequest.ReferenceID,
			Currency:    currency,
			Amount:      qrCodeResponse.Amount,
			ExpiresAt:   expiresAt,
			Created:     qrCodeResponse.Created,
			Updated:     qrCodeResponse.Updated,
			QRString:    qrCodeResponse.QRString,
			CallbackURL: callbackURL,
			Type:        types,
			Customer:    items.Customer,
			ItemID:      items.Id,
			Status:      qrCodeResponse.Status,
		}

		res, err := tc.service.AddPayments(newTransaction)
		if err != nil {
			c.Logger().Error("error inserting QR code payment:", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, err.Error(), nil))
		}

		dataResponse := CoreToTransactionResponse(res)

		return c.JSON(helper.ResponseFormat(http.StatusCreated, "Successfully add payment data", dataResponse))
	}
}

// AddTransactions implements transactions.Handler.
func (tc *transactionController) AddTransactions() echo.HandlerFunc {
	var orderSequence int
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID == 0 {
			c.Logger().Error("empty decode tokens")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "jwt invalid", nil))
		}

		users, err := tc.userRepo.GetUserById(userID)
		if err != nil {
			c.Logger().Error("failed to get user information: ", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "failed to get user information", nil))
		}

		detailRequest := new(TransactionDetailsRequest)
		if err := c.Bind(detailRequest); err != nil {
			c.Logger().Error("error on bind transaction input", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid request payload", nil))
		}

		productIDs := make([]uint, len(detailRequest.Details))
		for i, detail := range detailRequest.Details {
			productIDs[i] = detail.ProductID
		}

		products, err := tc.productRepo.GetProductByIds(productIDs)
		if err != nil {
			c.Logger().Error("error get product data", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to get products", nil))
		}

		itemDetails := make([]transactions.DetailCore, 0, len(detailRequest.Details))
		subTotal := 0
		for _, detail := range detailRequest.Details {
			product := findProductByID(products, detail.ProductID)
			if product == nil {
				c.Logger().Error("error find product detail", detail.ProductID)
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Product not found", nil))
			}

			itemDetail := transactions.DetailCore{
				ProductID: detail.ProductID,
				Quantity:  detail.Quantity,
				Price:     product.Price,
				Total:     product.Price * detail.Quantity,
			}

			err = tc.productRepo.UpdateProductStock(detail.ProductID, detail.Quantity)
			if err != nil {
				c.Logger().Error("error on update data", err.Error())
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to update product stock", nil))
			}

			subTotal += itemDetail.Total
			itemDetails = append(itemDetails, itemDetail)
		}

		orderSequence++
		timestamp := time.Now().Unix()
		orderID := fmt.Sprintf("ORDERID:%d-S:%d-C:%s-Timestamp:%d", orderSequence, userID, detailRequest.Customer, timestamp)
		newItem := transactions.ItemCore{
			SubTotal: subTotal,
			Customer: detailRequest.Customer,
			UserID:   userID,
			OrderID:  orderID,
			Status:   "Pending",
			Details:  itemDetails,
			UserName: users.Name,
		}

		createdDetailTransaction, err := tc.service.AddTransactions(userID, newItem)
		if err != nil {
			c.Logger().Error("error on add transaction", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, err.Error(), nil))
		}
		fmt.Printf("INI DATA USER NAME:%s", createdDetailTransaction.UserName)
		dataResponse := CoreToTransactionDetailsResponse(createdDetailTransaction)

		return c.JSON(helper.ReponseFormatWithMeta(http.StatusCreated, "Successfully add transaction detail data", dataResponse, nil))

	}
}

func findProductByID(products []products.Core, productID uint) *products.Core {
	for _, product := range products {
		if product.ID == productID {
			return &product
		}
	}
	return nil
}
