package repository

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ExternalID         string
	Status             string
	InvoiceURL         string
	Amount             int
	Customer           string
	UserID             uint
	Transactiondetails []Transactiondetail `gorm:"foreignKey:TransactionID"`
}

type Transactiondetail struct {
	gorm.Model
	TransactionID uint
	ExternalID    string
	ProductID     uint
	Quantity      int
	Total         int
	Customer      string
	UserID        uint
	Transaction   Transaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:TransactionID"`
}
