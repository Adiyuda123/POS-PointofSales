package repository

import (
	"POS-PointofSales/features/transactions"
)

func CoreToModel(data transactions.DetailCore) Transactiondetail {
	return Transactiondetail{
		TransactionID: data.TransactionID,
		ExternalID:    data.ExternalID,
		ProductID:     data.ProductID,
		Quantity:      data.Quantity,
		Total:         data.Total,
		Customer:      data.Customer,
		// UserID:     data.UserID,
	}
}

func ModelToCore(data Transactiondetail) transactions.DetailCore {

	result := transactions.DetailCore{
		Id:         data.ID,
		ExternalID: data.ExternalID,
		ProductID:  data.ProductID,
		Quantity:   data.Quantity,
		Total:      data.Quantity,
		Customer:   data.Customer,
		// UserID:        data.UserID,
		// TransactionID: data.TransactionID,
	}

	return result
}

func ListtdetailsTotdetailsCore(detail []Transactiondetail) []transactions.DetailCore {
	var data []transactions.DetailCore
	for _, v := range detail {
		data = append(data, ModelToCore(v))
	}
	return data
}

func TransactionToModel(data transactions.Core) Transaction {
	return Transaction{
		ExternalID: data.ExternalID,
		Status:     data.Status,
		InvoiceURL: data.InvoiceURL,
		Amount:     data.Amount,
		Customer:   data.Customer,
		UserID:     data.UserID,
		// Transactiondetails: []Transactiondetail{},
	}
}

func ModelToTransaction(data Transaction) transactions.Core {

	result := transactions.Core{
		Id:         data.ID,
		ExternalID: data.ExternalID,
		Status:     data.Status,
		InvoiceURL: data.InvoiceURL,
		Amount:     data.Amount,
		Customer:   data.Customer,
		UserID:     data.UserID,
		Details:    []transactions.DetailCore{},
	}

	return result
}

func ListtsToTransactionCore(transaction []Transaction) []transactions.Core {
	var data []transactions.Core
	for _, v := range transaction {
		data = append(data, ModelToTransaction(v))
	}
	return data
}
