package handler

type InputUpdate struct {
	ID           uint   `json:"id"`
	Name         string `json:"product_name"`
	Descriptions string `json:"descriptions"`
	Pictures     string `json:"pictures"`
	Price        int    `json:"price"`
}

type InputRequest struct {
	Name         string `json:"product_name"`
	Descriptions string `json:"descriptions"`
	Pictures     string `json:"pictures"`
	Price        int    `json:"price"`
	Stock        int    `json:"stock"`
}

type InputRestock struct {
	ID    uint `json:"product_id"`
	Stock int  `json:"restock_quantity"`
}
