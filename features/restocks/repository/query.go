package repository

import (
	"POS-PointofSales/features/products/repository"
	"POS-PointofSales/features/restocks"
	"errors"

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

// AddRestock implements restocks.Repository.
func (rm *restockModel) AddRestock(userID uint, restockInput restocks.Core) error {
	productToRestock := &repository.Product{}
	if err := rm.db.First(productToRestock, "id = ? AND user_id = ?", restockInput.ProductID, userID).Error; err != nil {
		log.Error("Error finding product")
		return errors.New("product not found or not owned by the user")
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
