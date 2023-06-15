package handler

import (
	"POS-PointofSales/features/restocks"
	"time"
)

type RestockResponse struct {
	ID        uint   `json:"restock_id"`
	ProductID uint   `json:"product_id"`
	Quantity  int    `json:"restock_quantity"`
	Date      string `json:"restock_date"`
	UserID    uint   `json:"user_id"`
	UserName  string `json:"user_name"`
}

func CoreToRestockResponse(data restocks.Core) RestockResponse {
	return RestockResponse{
		ID:        data.ID,
		ProductID: data.ProductID,
		Quantity:  data.Quantity,
		Date:      data.CreatedAt.Add(7 * time.Hour).Format("2006-01-02 15:04"),
		UserID:    data.UserID,
		UserName:  data.UserName,
	}
}

func CoreToGetAllRestockResponse(data []restocks.Core) []RestockResponse {
	res := make([]RestockResponse, len(data))
	for i, val := range data {
		res[i] = CoreToRestockResponse(val)
	}
	return res
}
