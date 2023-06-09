package repository

import (
	pRepo "POS-PointofSales/features/products/repository"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(50);not null"`
	Email    string `gorm:"type:varchar(50);unique;not null"`
	Phone    string `gorm:"type:varchar(50);unique;not null"`
	Pictures string `gorm:"type:text"`
	Password string `gorm:"type:varchar(255);not null"`
	Products []pRepo.Product
}
