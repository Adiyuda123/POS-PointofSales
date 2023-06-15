package repository

import (
	"POS-PointofSales/features/restocks"
)

func CoreToModel(data restocks.Core) Restock {
	return Restock{
		ID:        data.ID,
		ProductID: data.ProductID,
		Quantity:  data.Quantity,
		UserID:    data.UserID,
		UserName:  data.UserName,
		CreatedAt: data.CreatedAt,
	}
}

func ModelToCore(data Restock) restocks.Core {

	result := restocks.Core{
		ID:        data.ID,
		ProductID: data.ProductID,
		Quantity:  data.Quantity,
		CreatedAt: data.CreatedAt,
		UserID:    data.UserID,
		UserName:  data.UserName,
	}

	return result
}

func ListrestockTorestockCore(restock []Restock) []restocks.Core {
	var data []restocks.Core
	for _, v := range restock {
		data = append(data, ModelToCore(v))
	}
	return data
}
