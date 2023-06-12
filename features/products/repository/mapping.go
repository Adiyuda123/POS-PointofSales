package repository

import (
	"POS-PointofSales/features/products"
)

func CoreToModel(data products.Core) Product {
	return Product{
		Name:         data.Name,
		Descriptions: data.Descriptions,
		Price:        data.Price,
		Pictures:     data.Pictures,
		Stock:        data.Stock,
		UserID:       data.UserID,
		UserName:     data.UserName,
	}
}

func ModelToCore(data Product) products.Core {

	result := products.Core{
		ID:           data.ID,
		Name:         data.Name,
		Descriptions: data.Descriptions,
		Price:        data.Price,
		Pictures:     data.Pictures,
		Stock:        data.Stock,
		UserID:       data.UserID,
		UserName:     data.UserName,
	}

	return result
}

func ListproductToproductCore(product []Product) []products.Core {
	var data []products.Core
	for _, v := range product {
		data = append(data, ModelToCore(v))
	}
	return data
}
