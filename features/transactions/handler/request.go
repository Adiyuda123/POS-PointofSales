package handler

type TransactionRequest struct {
	ExternalID string `json:"external_id" form:"external_id"`
	Status     string `json:"status" form:"status"`
	InvoiceUrl string `json:"invoice_url" form:"invoice_url"`
	Amount     int    `json:"amount" form:"amount"`
	Customer   string `json:"customer" form:"customer"`
	UserID     uint   `json:"user_id" form:"user_id"`
	// Details    []transactiondetails.Core
}

type TransactionDetailsRequest struct {
	ExternalID    string `json:"external_id" form:"external_id"`
	TransactionID uint   `json:"transaction_id" form:"transaction_id"`
	ProductID     uint   `json:"product_id" form:"product_id"`
	Quantity      int    `json:"quantity" form:"quantity"`
	Customer      string `json:"customer" form:"customer"`
	// Total         int  `json:"total" form:"total"`
}
