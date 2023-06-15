// handler/transaction.go
package handler

import (
	"POS-PointofSales/features/products"
	"POS-PointofSales/features/transactions"
	"POS-PointofSales/helper"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type transactionController struct {
	service     transactions.UseCase
	productRepo products.Repository
}

func New(service transactions.UseCase, services products.Repository) transactions.Handler {
	return &transactionController{
		service:     service,
		productRepo: services,
	}
}

// AddTransactions implements transactions.Handler.
func (tc *transactionController) AddTransactions() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID == 0 {
			c.Logger().Error("empty decode tokens")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "jwt invalid", nil))
		}
		detailRrequest := new(TransactionDetailsRequest)
		if err := c.Bind(detailRrequest); err != nil {
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid request payload", nil))
		}

		createdTransaction, err := tc.service.CreateTransactionsIfNotExists(userID, detailRrequest.TransactionID, detailRrequest.Customer)
		if err != nil {
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, err.Error(), nil))
		}

		product, err := tc.productRepo.GetProductById(detailRrequest.ProductID)
		if err != nil {
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to get product", nil))
		}

		total := product.Price * detailRrequest.Quantity

		newTransaction := transactions.DetailCore{
			ExternalID:    detailRrequest.ExternalID,
			ProductID:     detailRrequest.ProductID,
			Quantity:      detailRrequest.Quantity,
			Total:         total,
			Customer:      detailRrequest.Customer,
			UserID:        userID,
			TransactionID: createdTransaction.Id,
		}

		createdDetailTransaction, err := tc.service.AddTransactions(userID, newTransaction)
		if err != nil {
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, err.Error(), nil))
		}
		createdDetailTransaction.TransactionID = createdTransaction.Id

		dataResponse := CoreToTransactionDetailsResponse(createdDetailTransaction)
		responseString := fmt.Sprintf("%v", dataResponse)

		return c.JSON(helper.ResponseFormat(http.StatusCreated, responseString, nil))
	}
}

func (tc *transactionController) CreateTransactions() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := helper.DecodeToken(c)
		if userId == 0 {
			c.Logger().Error("empty decode tokens")
			return c.JSON(http.StatusBadRequest, "jwt invalid")
		}

		request := new(TransactionRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid request payload")
		}
		amount, err := tc.service.GetTotalAmount(request.ExternalID, request.Customer)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		invoiceURL, err := helper.CreateInvoice(amount, request.ExternalID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		newTransaction := transactions.Core{
			ExternalID: request.ExternalID,
			Status:     "Pending",
			InvoiceURL: invoiceURL,
			Amount:     amount,
			Customer:   request.Customer,
			UserID:     userId,
		}

		createdTransaction, err := tc.service.CreateTransactions(userId, newTransaction)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		response := CoreToTransactionResponse(createdTransaction)
		return c.JSON(http.StatusCreated, response)
	}
}
