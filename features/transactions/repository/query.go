package repository

import (
	"POS-PointofSales/features/products/repository"
	"POS-PointofSales/features/transactions"
	"time"

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

// SelectHistoryTransaction implements transactions.Repository.
func (tm *transactionModel) SelectHistoryTransaction(userID uint, limit int, offset int, search string, fromDate time.Time, toDate time.Time) ([]transactions.ItemCore, int, error) {
	nameSearch := "%" + search + "%"
	totalData := int64(-1)
	var ItemModel []Item

	query := tm.db.Table("items").
		Select("items.id, items.sub_total, items.customer, items.order_id, items.status, items.user_name").
		Joins("JOIN users ON items.user_id = users.id").
		Joins("JOIN item_details ON items.id = item_details.item_id").
		Joins("JOIN products ON item_details.product_id = products.id").
		Group("items.id").
		Preload("Details").
		Preload("Details.Product").
		Limit(limit).Offset(offset).
		Order("items.id DESC")

	if userID != 1 {
		query = query.Where("items.user_id = ?", userID)
	}

	if !fromDate.IsZero() {
		query = query.Where("items.created_at >= ?", fromDate)
	}

	if !toDate.IsZero() {
		query = query.Where("items.created_at < ?", toDate)
	}

	if search != "" {
		if err := query.Where("items.id LIKE ? OR item_details.quantity LIKE ? OR item_details.product_id LIKE ? OR items.user_id LIKE ? OR items.user_name LIKE ?",
			nameSearch, nameSearch, nameSearch, nameSearch, nameSearch).Find(&ItemModel).Error; err != nil {
			log.Errorf("error on finding search: %w", err)
			return []transactions.ItemCore{}, int(totalData), err
		}
		if err := tm.db.Table("items").
			Joins("JOIN users ON items.user_id = users.id").
			Joins("JOIN item_details ON items.id = item_details.item_id").
			Joins("JOIN products ON item_details.product_id = products.id").
			Where("items.id LIKE ? OR item_details.quantity LIKE ? OR item_details.product_id LIKE ? OR items.user_id LIKE ? OR items.user_name LIKE ?",
				nameSearch, nameSearch, nameSearch, nameSearch, nameSearch).
			Count(&totalData).Error; err != nil {
			log.Errorf("error on count filtered data: %w", err)
			return []transactions.ItemCore{}, int(totalData), err
		}
	} else {
		if err := query.Find(&ItemModel).Error; err != nil {
			log.Errorf("error on finding data without search: %w", err)
			return []transactions.ItemCore{}, int(totalData), err
		}
		if err := tm.db.Table("items").Count(&totalData).Error; err != nil {
			log.Errorf("error on counting data without search: %w", err)
			return []transactions.ItemCore{}, int(totalData), err
		}
	}

	restockCoreAll := ListItemToCore(ItemModel)
	return restockCoreAll, int(totalData), nil
}

// GetItemById implements transactions.Repository.
func (tm *transactionModel) SelectItemByOrderId(orderID string) (transactions.ItemCore, error) {
	res := Item{}
	if err := tm.db.Table("items").
		Select("items.id, items.sub_total, items.customer, items.order_id, items.status").
		Joins("JOIN users ON items.user_id = users.id").
		Where("items.order_id = ?", orderID).
		First(&res).Error; err != nil {
		log.Error("error occurs in finding item by id:", err.Error())
		return transactions.ItemCore{}, err
	}

	return ItemToCore(res), nil
}

// InsertPayments implements transactions.Repository.
func (tm *transactionModel) InsertPayments(newTransaction transactions.Core) (transactions.Core, error) {
	payment := CoreToTransaction(newTransaction)
	qrCodePayment := Transaction{
		ID:          payment.ID,
		ExternalID:  payment.ExternalID,
		OrderID:     payment.OrderID,
		Currency:    payment.Currency,
		Amount:      payment.Amount,
		ExpiresAt:   payment.ExpiresAt,
		Created:     payment.Created,
		Updated:     payment.Updated,
		QRString:    payment.QRString,
		CallbackURL: payment.CallbackURL,
		Type:        payment.Type,
		Customer:    payment.Customer,
		ItemID:      payment.ItemID,
		UserID:      payment.UserID,
		Status:      payment.Status,
		Item:        Item{},
	}

	if err := tm.db.Create(&qrCodePayment).Error; err != nil {
		return transactions.Core{}, err
	}

	return TransactionToCore(qrCodePayment), nil
}

// InsertDetailTransactions implements transactions.Repository.
func (tm *transactionModel) InsertDetailTransactions(userID uint, inputDetail transactions.ItemCore) (transactions.ItemCore, error) {
	productTransactionDetail := CoreToItem(inputDetail)

	var product repository.Product
	if err := tm.db.First(&product, productTransactionDetail.Details[0].ProductID).Error; err != nil {
		log.Error("ERROR FIRST", err.Error())
		return transactions.ItemCore{}, err
	}

	if err := tm.db.Create(&productTransactionDetail); err.Error != nil {
		return transactions.ItemCore{}, err.Error
	}

	totalQuantity := 0
	for i := 0; i < len(productTransactionDetail.Details); i++ {
		totalQuantity /= productTransactionDetail.Details[i].Quantity
	}
	for i := 0; i < len(productTransactionDetail.Details); i++ {
		productID := productTransactionDetail.Details[i].ProductID
		if err := tm.db.Model(&repository.Product{}).Where("id = ?", productID).Update("stock", gorm.Expr("stock - ?", totalQuantity)).Error; err != nil {
			log.Error("error occurred in updating stock after transaction")
			return transactions.ItemCore{}, err
		}
	}

	return ItemToCore(productTransactionDetail), nil
}
