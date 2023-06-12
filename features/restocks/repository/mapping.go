package repository

import (
	"POS-PointofSales/features/restocks"
)

func CoreToModel(data restocks.Core) Restock {
	return Restock{
		ProductID: data.ProductID,
		Quantity:  data.Quantity,
		UserID:    data.UserID,
	}
}

func ModelToCore(data Restock) restocks.Core {

	result := restocks.Core{
		ID:        data.ID,
		ProductID: data.ProductID,
		Quantity:  data.Quantity,
		UserID:    data.UserID,
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
