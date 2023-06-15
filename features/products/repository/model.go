package repository

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID           uint `gorm:"primarykey"`
	Name         string
	Descriptions string
	Price        int
	Pictures     string
	Stock        int
	UserID       uint
	UserName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
