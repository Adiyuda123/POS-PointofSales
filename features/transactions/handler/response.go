package handler

import (
	"POS-PointofSales/features/transactions"
)

type TransactionResponse struct {
	ExternalID string                       `json:"external_id" form:"external_id"`
	Status     string                       `json:"status" form:"status"`
	InvoiceURL string                       `json:"invoice_url" form:"invoice_url"`
	Amount     int                          `json:"amount" form:"amount"`
	Customer   string                       `json:"customer" form:"customer"`
	UserID     uint                         `json:"user_id" form:"user_id"`
	Details    []TransactionDetailsResponse `json:"details" form:"details"`
}

func CoreToTransactionResponse(data transactions.Core) TransactionResponse {
	response := TransactionResponse{
		ExternalID: data.ExternalID,
		Status:     data.Status,
		InvoiceURL: data.InvoiceURL,
		Amount:     data.Amount,
		Customer:   data.Customer,
		UserID:     data.UserID,
		Details:    make([]TransactionDetailsResponse, len(data.Details)),
	}

	for i, detail := range data.Details {
		response.Details[i] = TransactionDetailsResponse{
			ExternalID: detail.ExternalID,
			ProductID:  detail.ProductID,
			Quantity:   detail.Quantity,
			Total:      detail.Total,
			Customer:   detail.Customer,
		}
	}

	return response
}

func CoreToGetAllTransactionsResponse(data []transactions.Core) []TransactionResponse {
	res := make([]TransactionResponse, len(data))
	for i, val := range data {
		res[i] = CoreToTransactionResponse(val)
	}
	return res
}

type TransactionDetailsResponse struct {
	ID         uint   `json:"id" form:"id"`
	ExternalID string `json:"external_id" form:"external_id"`
	ProductID  uint   `json:"product_id" form:"product_id"`
	Quantity   int    `json:"quantity" form:"quantity"`
	Total      int    `json:"total" form:"total"`
	Customer   string `json:"customer" form:"customer"`
	UserID     uint   `json:"user_id" form:"user_id"`
}

func CoreToTransactionDetailsResponse(data transactions.DetailCore) TransactionDetailsResponse {
	return TransactionDetailsResponse{
		ID:         data.Id,
		ExternalID: data.ExternalID,
		ProductID:  data.ProductID,
		Quantity:   data.Quantity,
		Total:      data.Total,
		Customer:   data.Customer,
		UserID:     data.UserID,
	}
}

func CoreToGetAllTransactionDetailsResponse(data []transactions.DetailCore) []TransactionDetailsResponse {
	res := make([]TransactionDetailsResponse, len(data))
	for i, val := range data {
		res[i] = CoreToTransactionDetailsResponse(val)
	}
	return res
}
