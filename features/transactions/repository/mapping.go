package repository

import (
	"POS-PointofSales/features/transactions"
)

func CoreToItem(data transactions.ItemCore) Item {
	result := Item{
		SubTotal: data.SubTotal,
		Customer: data.Customer,
		UserID:   data.UserID,
		UserName: data.UserName,
		OrderID:  data.OrderID,
		Status:   data.Status,
		Details:  make([]ItemDetail, len(data.Details)),
	}

	for i, detail := range data.Details {
		detailItem := ItemDetail{
			ItemID:    detail.ItemID,
			ProductID: detail.ProductID,
			Quantity:  detail.Quantity,
			Total:     detail.Total,
		}
		result.Details[i] = detailItem
	}

	return result
}

func ItemToCore(data Item) transactions.ItemCore {
	result := transactions.ItemCore{
		Id:       data.ID,
		SubTotal: data.SubTotal,
		Customer: data.Customer,
		OrderID:  data.OrderID,
		Status:   data.Status,
		UserID:   data.UserID,
		Details:  make([]transactions.DetailCore, len(data.Details)),
	}

	for i, detail := range data.Details {
		detailCore := transactions.DetailCore{
			Id:        detail.ID,
			ItemID:    detail.ItemID,
			ProductID: detail.ProductID,
			Quantity:  detail.Quantity,
			Total:     detail.Total,
			Price:     detail.Price,
		}
		result.Details[i] = detailCore
	}

	return result
}

func ListItemToCore(itemModel []Item) []transactions.ItemCore {
	var data []transactions.ItemCore
	for _, v := range itemModel {
		data = append(data, ItemToCore(v))
	}
	return data
}

// ###############################################################################

func CoreToTransaction(data transactions.Core) Transaction {
	return Transaction{
		ID:          data.ID,
		ExternalID:  data.ExternalID,
		OrderID:     data.OrderID,
		Currency:    data.Currency,
		Amount:      data.Amount,
		ExpiresAt:   data.ExpiresAt,
		Created:     data.Created,
		Updated:     data.Updated,
		QRString:    data.QRString,
		CallbackURL: data.CallbackURL,
		Type:        data.Type,
		Customer:    data.Customer,
		ItemID:      data.ItemID,
		UserID:      data.UserID,
		Status:      data.Status,
		Item:        Item{},
	}

}

func TransactionToCore(data Transaction) transactions.Core {
	return transactions.Core{
		ID:          data.ID,
		ExternalID:  data.ExternalID,
		OrderID:     data.OrderID,
		Currency:    data.Currency,
		Amount:      data.Amount,
		ExpiresAt:   data.ExpiresAt,
		Created:     data.Created,
		Updated:     data.Updated,
		QRString:    data.QRString,
		CallbackURL: data.CallbackURL,
		Type:        data.Type,
		Customer:    data.Customer,
		ItemID:      data.ItemID,
		UserID:      data.UserID,
		Status:      data.Status,
	}
}

func ListTransactionToCore(itemModel []Transaction) []transactions.Core {
	var data []transactions.Core
	for _, v := range itemModel {
		data = append(data, TransactionToCore(v))
	}
	return data
}
