package database

import (
	pRepo "POS-PointofSales/features/products/repository"
	rRepo "POS-PointofSales/features/restocks/repository"
	tRepo "POS-PointofSales/features/transactions/repository"
	uRepo "POS-PointofSales/features/users/repository"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&uRepo.User{})
	db.AutoMigrate(&pRepo.Product{})
	db.AutoMigrate(&rRepo.Restock{})
	db.AutoMigrate(&tRepo.Transaction{})
	db.AutoMigrate(&tRepo.Transactiondetail{})
}
