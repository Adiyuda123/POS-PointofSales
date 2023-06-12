package repository

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name         string
	Descriptions string
	Price        int
	Pictures     string
	Stock        int
	UserID       uint
	UserName     string
}
