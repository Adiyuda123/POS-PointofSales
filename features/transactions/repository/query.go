package repository

import (
	"POS-PointofSales/features/products/repository"
	"POS-PointofSales/features/transactions"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type transactionModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) transactions.Repository {
	return &transactionModel{
		db: db,
	}
}

// // GetTransactionByExternalID implements transactions.Repository.
// func (tm *transactionModel) GetTransactionByExternalID(externalID string) (transactions.Core, error) {
// 	res := Transaction{}
// 	if err := tm.db.Table("transactions").
// 		Where("external_id = ?", externalID).
// 		First(&res).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			// Return an empty Core if the transaction doesn't exist
// 			return transactions.Core{}, nil
// 		}
// 		log.Error("error occurs in finding transaction by external ID:", err.Error())
// 		return transactions.Core{}, err
// 	}

// 	return ModelToTransaction(res), nil
// }

// SelectTransactionById implements transactions.Repository.
func (tm *transactionModel) SelectTransactionById(id uint) (transactions.Core, error) {
	res := Transaction{}
	if err := tm.db.Table("transactions").
		Where("id = ?", id).
		First(&res).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return an empty Core if the transaction doesn't exist
			return transactions.Core{}, nil
		}
		log.Error("error occurs in finding transaction by id:", err.Error())
		return transactions.Core{}, err
	}

	return ModelToTransaction(res), nil
}

// GetTotalAmount implements transactions.Repository.
func (tm *transactionModel) GetTotalAmount(externalID string, customer string) (int, error) {
	var amount int
	if err := tm.db.Model(&Transactiondetail{}).
		Where("transactiondetails.external_id = ? AND transactiondetails.customer = ?", externalID, customer).
		Select("SUM(transactiondetails.total)").
		Scan(&amount).Error; err != nil {
		return 0, err
	}
	return amount, nil
}

// InsertDetailTransactions implements transactions.Repository.
func (tm *transactionModel) InsertDetailTransactions(userID uint, inputDetail transactions.DetailCore) (transactions.DetailCore, error) {
	productTransactionDetail := CoreToModel(inputDetail)

	var product repository.Product
	if err := tm.db.First(&product, productTransactionDetail.ProductID).Error; err != nil {
		return transactions.DetailCore{}, err
	}

	total := product.Price * productTransactionDetail.Quantity
	productTransactionDetail.Total = total

	if err := tm.db.Create(&productTransactionDetail).Error; err != nil {
		return transactions.DetailCore{}, err
	}

	tm.db.Model(&product).
		Where("id = ?", productTransactionDetail.ProductID).
		UpdateColumn("stock", gorm.Expr("stock - ?", productTransactionDetail.Quantity))

	return ModelToCore(productTransactionDetail), nil
}

// InsertTransactions implements transactions.Repository.
func (tm *transactionModel) InsertTransactions(userId uint, input transactions.Core) (transactions.Core, error) {
	transaction := TransactionToModel(input)
	if err := tm.db.Create(&transaction).Error; err != nil {
		return transactions.Core{}, err
	}
	// panic("QUERY PANIC SATU DUA")
	// Simpan detail transaksi
	// for _, detail := range transaction.Transactiondetails {
	// 	detail.ExternalID = transaction.ExternalID
	// 	if err := tm.db.Create(&detail).Error; err != nil {
	// 		return transactions.Core{}, err
	// 	}
	// }

	for _, detail := range transaction.Transactiondetails {
		detail.ExternalID = transaction.ExternalID
		if err := tm.db.Create(&detail).Error; err != nil {
			return transactions.Core{}, err
		}
		transaction.Transactiondetails = append(transaction.Transactiondetails, detail)
	}

	// var amount int
	// if err := tm.db.Model(&Transactiondetail{}).
	// 	Where("transactiondetails.external_id = ? AND transactiondetails.customer = ?",
	// 		input.ExternalID, input.Customer).
	// 	Select("SUM(transactiondetails.total)").
	// 	Scan(&amount).Error; err != nil {
	// 	return transactions.Core{}, err
	// }
	// // transaction.Amount = amount
	// fmt.Printf("amount: %d\n", amount)

	// if err := tm.db.Model(&transactions.Core{}).
	// 	Where("external_id = ? AND customer = ?", input.ExternalID, input.Customer).
	// 	Update("amount", amount).Error; err != nil {
	// 	return transactions.Core{}, err
	// }

	return input, nil
}
