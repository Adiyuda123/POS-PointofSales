package handler

type InputRestock struct {
	ProductID uint `json:"product_id"`
	Stock     int  `json:"restock_quantity"`
}
