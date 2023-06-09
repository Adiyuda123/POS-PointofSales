package database

import (
	pRepo "POS-PointofSales/features/products/repository"
	uRepo "POS-PointofSales/features/users/repository"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(uRepo.User{})
	db.AutoMigrate(pRepo.Product{})
}
