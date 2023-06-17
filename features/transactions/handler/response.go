package handler

import (
	"POS-PointofSales/features/transactions"
)

type TransactionDetailsResponse struct {
	ID       uint                 `json:"id" form:"id"`
	SubTotal int                  `json:"sub_total" form:"sub_total"`
	Customer string               `json:"customer" form:"customer"`
	Status   string               `json:"status" form:"status"`
	UserID   uint                 `json:"user_id" form:"user_id"`
	UserName string               `json:"user_name" form:"user_name"`
	OrderID  string               `json:"order_id" form:"order_id"`
	Details  []ItemDetailResponse `json:"details" form:"details"`
}
type ItemDetailResponse struct {
	ProductID uint `json:"product_id" form:"product_id"`
	Quantity  int  `json:"quantity" form:"quantity"`
}

func CoreToItemDetailResponse(details []transactions.DetailCore) []ItemDetailResponse {
	res := make([]ItemDetailResponse, len(details))
	for i, detail := range details {
		res[i] = ItemDetailResponse{
			ProductID: detail.ProductID,
			Quantity:  detail.Quantity,
		}
	}
	return res
}

func CoreToTransactionDetailsResponse(data transactions.ItemCore) TransactionDetailsResponse {
	response := TransactionDetailsResponse{
		ID:       data.Id,
		SubTotal: data.SubTotal,
		Customer: data.Customer,
		Status:   data.Status,
		UserID:   data.UserID,
		UserName: data.UserName,
		OrderID:  data.OrderID,
		Details:  make([]ItemDetailResponse, len(data.Details)),
	}

	for i, detail := range data.Details {
		detailResponse := ItemDetailResponse{
			ProductID: detail.ProductID,
			Quantity:  detail.Quantity,
		}
		response.Details[i] = detailResponse
	}

	return response
}

func CoreToItemDetailResponse2(details []transactions.DetailCore) []ItemDetailResponse {
	res := make([]ItemDetailResponse, len(details))
	for i, detail := range details {
		res[i] = ItemDetailResponse{
			ProductID: detail.ProductID,
			Quantity:  detail.Quantity,
		}
	}
	return res
}

func CoreToGetAllTransactionDetailsResponse(data []transactions.ItemCore) []TransactionDetailsResponse {
	res := make([]TransactionDetailsResponse, len(data))
	for i, val := range data {
		res[i] = CoreToTransactionDetailsResponse(val)
	}
	return res
}

//##################################################################################################

type TransactionResponse struct {
	ID          string `json:"id"`
	ExternalID  string `json:"external_id"`
	OrderID     string `json:"order_id"`
	Amount      int    `json:"amount"`
	QRString    string `json:"qr_string"`
	CallbackURL string `json:"callback_url"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	Customer    string `json:"customer" form:"customer"`
	ItemID      uint   `json:"item_id" form:"item_id"`
	UserID      uint   `json:"user_id" form:"user_id"`
}

func CoreToTransactionResponse(data transactions.Core) TransactionResponse {
	return TransactionResponse{
		ID:          data.ID,
		ExternalID:  data.ExternalID,
		OrderID:     data.OrderID,
		Amount:      data.Amount,
		QRString:    data.QRString,
		CallbackURL: data.CallbackURL,
		Type:        data.Type,
		Status:      data.Status,
		Created:     data.Created,
		Updated:     data.Updated,
		Customer:    data.Customer,
		ItemID:      data.ItemID,
		UserID:      data.UserID,
	}
}

func CoreToGetAllTransactionResponse(data []transactions.Core) []TransactionResponse {
	res := make([]TransactionResponse, len(data))
	for i, val := range data {
		res[i] = CoreToTransactionResponse(val)
	}
	return res
}
