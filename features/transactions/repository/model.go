package repository

import (
	"POS-PointofSales/features/products/repository"

	"gorm.io/gorm"
)

type Transaction struct {
	ID          string
	ExternalID  string
	Amount      int
	QRString    string
	CallbackURL string
	Type        string
	Status      string
	Created     string
	Updated     string
	Customer    string
	ItemID      uint
	UserID      uint
	OrderID     string
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
