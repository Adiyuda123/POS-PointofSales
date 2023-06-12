package repository

import (
	"gorm.io/gorm"
)

type Restock struct {
	gorm.Model
	ProductID uint
	Quantity  int
	Date      string
	UserID    uint
}
