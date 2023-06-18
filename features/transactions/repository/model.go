package repository

import (
	"POS-PointofSales/features/products/repository"

	"gorm.io/gorm"
)

type Transaction struct {
	ID          string
	ExternalID  string
	OrderID     string
	Currency    string
	Amount      int
	ExpiresAt   string
	Created     string
	Updated     string
	QRString    string
	CallbackURL string
	Type        string
	Customer    string
	ItemID      uint
	UserID      uint
	Status      string
	Item        Item `gorm:"foreignKey:ItemID"`
}

type Item struct {
	gorm.Model
	SubTotal int
	Customer string
	UserID   uint
	UserName string
	OrderID  string
	Status   string
	Details  []ItemDetail `gorm:"foreignKey:ItemID"`
}

type ItemDetail struct {
	gorm.Model
	ItemID    uint
	ProductID uint
	Quantity  int
	Total     int
	Price     int
	Product   repository.Product `gorm:"foreignKey:ProductID"`
}
