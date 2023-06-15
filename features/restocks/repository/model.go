package repository

import (
	"time"

	"gorm.io/gorm"
)

type Restock struct {
	ID        uint `gorm:"primarykey"`
	ProductID uint
	Quantity  int
	UserID    uint
	UserName  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
