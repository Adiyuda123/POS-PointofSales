package repository

import (
	pRepo "POS-PointofSales/features/products/repository"
	rRepo "POS-PointofSales/features/restocks/repository"
	tRepo "POS-PointofSales/features/transactions/repository"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string
	Email        string
	Phone        string
	Pictures     string
	Password     string
	Products     []pRepo.Product     `gorm:"foreignKey:UserID"`
	Restocks     []rRepo.Restock     `gorm:"foreignKey:UserID"`
	Transactions []tRepo.Transaction `gorm:"foreignKey:UserID"`
}
