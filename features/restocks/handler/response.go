package handler

import (
	"POS-PointofSales/features/restocks"
)

type RestockResponse struct {
	ID        uint   `json:"restock_id"`
	ProductID uint   `json:"product_id"`
	Quantity  int    `json:"restock_quantity"`
	Date      string `json:"restock_date"`
	UserID    uint   `json:"user_id"`
}

func CoreToRestockResponse(data restocks.Core) RestockResponse {
	return RestockResponse{
		ID:        data.ID,
		ProductID: data.ProductID,
		Quantity:  data.Quantity,
		UserID:    data.UserID,
	}
}

func CoreToGetAllProductResponse(data []restocks.Core) []RestockResponse {
	res := make([]RestockResponse, len(data))
	for i, val := range data {
		res[i] = CoreToRestockResponse(val)
	}
	return res
}
