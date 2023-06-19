package repository

import (
	"POS-PointofSales/features/products/repository"
	"POS-PointofSales/features/restocks"
	"errors"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type restockModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) restocks.Repository {
	return &restockModel{
		db: db,
	}
}

// SelectAllRestock implements restocks.Repository.
func (rm *restockModel) SelectAllRestock(userID uint, limit, offset int, search string, fromDate, toDate time.Time) ([]restocks.Core, int, error) {
	nameSearch := "%" + search + "%"
	totalData := int64(-1)
	var restockModel []Restock

	query := rm.db.Table("restocks").
		Select("restocks.id, restocks.quantity,  restocks.product_id,  restocks.user_id, users.name AS user_name").
		Joins("JOIN users ON restocks.user_id = users.id").
		Joins("JOIN products ON restocks.product_id = products.id").
		Group("restocks.id").
		Limit(limit).Offset(offset).
		Order("restocks.id DESC")

	if userID != 1 {
		query = query.Where("restocks.user_id = ?", userID)
	}

	if !fromDate.IsZero() {
		query = query.Where("restocks.created_at >= ?", fromDate)
	}

	if !toDate.IsZero() {
		query = query.Where("restocks.created_at < ?", toDate)
	}

	if search != "" {
		if err := query.Where("restocks.id LIKE ? OR restocks.quantity LIKE ? OR restocks.product_id LIKE ? OR restocks.user_id LIKE ? OR users.name LIKE ?",
			nameSearch, nameSearch, nameSearch, nameSearch, nameSearch).Find(&restockModel).Error; err != nil {
			log.Errorf("error on finding search: %w", err)
			return []restocks.Core{}, int(totalData), err
		}
		if err := rm.db.Table("restocks").
			Joins("JOIN users ON restocks.user_id = users.id").
			Joins("JOIN products ON restocks.product_id = products.id").
			Where("restocks.id LIKE ? OR restocks.quantity LIKE ? OR restocks.product_id LIKE ? OR restocks.user_id LIKE ? OR users.name LIKE ?",
				nameSearch, nameSearch, nameSearch, nameSearch, nameSearch).
			Count(&totalData).Error; err != nil {
			log.Errorf("error on count filtered data: %w", err)
			return []restocks.Core{}, int(totalData), err
		}
	} else {
		if err := query.Find(&restockModel).Error; err != nil {
			log.Errorf("error on finding data without search: %w", err)
			return []restocks.Core{}, int(totalData), err
		}
		if err := rm.db.Table("restocks").Count(&totalData).Error; err != nil {
			log.Errorf("error on counting data without search: %w", err)
			return []restocks.Core{}, int(totalData), err
		}
	}

	restockCoreAll := ListrestockTorestockCore(restockModel)
	return restockCoreAll, int(totalData), nil
}

// AddRestock implements restocks.Repository.
func (rm *restockModel) AddRestock(userID uint, restockInput restocks.Core) error {
	productToRestock := &repository.Product{}
	if err := rm.db.First(productToRestock, "id = ?", restockInput.ProductID).Error; err != nil {
		log.Error("Error finding product")
		return errors.New("product not found")
	}

	// Check if userID in product is different from the provided userID
	if userID != 1 && productToRestock.UserID != userID {
		log.Error("User is not authorized to restock this product")
		return errors.New("user is not authorized to restock this product")
	}

	productToRestock.Stock += restockInput.Quantity

	if err := rm.db.Save(productToRestock).Error; err != nil {
		log.Error("Failed to update product")
		return errors.New("product not updated")
	}

	newRestock := &Restock{
		ProductID: restockInput.ProductID,
		Quantity:  restockInput.Quantity,
		UserID:    userID,
	}
	if err := rm.db.Create(newRestock).Error; err != nil {
		log.Error("Failed to create restock")
		return errors.New("failed to create restock")
	}

	return nil
}
