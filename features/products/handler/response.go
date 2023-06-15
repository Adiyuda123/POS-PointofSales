package handler

import (
	"POS-PointofSales/features/products"
	"time"
)

type ProductResponse struct {
	ID           uint   `json:"id"`
	UserName     string `json:"user_name"`
	Name         string `json:"product_name"`
	Descriptions string `json:"descriptions"`
	Price        int    `json:"price"`
	Pictures     string `json:"pictures"`
	Stock        int    `json:"stock"`
	CreatedAt    string `json:"created_at"`
}

func CoreToProductResponse(data products.Core) ProductResponse {
	return ProductResponse{
		ID:           data.ID,
		UserName:     data.UserName,
		Name:         data.Name,
		Descriptions: data.Descriptions,
		Price:        data.Price,
		Pictures:     data.Pictures,
		Stock:        data.Stock,
		CreatedAt:    data.CreatedAt.Add(7 * time.Hour).Format("2006-01-02 15:04"),
	}
}

func CoreToGetAllProductResponse(data []products.Core) []ProductResponse {
	res := make([]ProductResponse, len(data))
	for i, val := range data {
		res[i] = CoreToProductResponse(val)
	}
	return res
}
