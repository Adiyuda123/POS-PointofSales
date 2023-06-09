package repository

import (
	"POS-PointofSales/features/products"

	"gorm.io/gorm"
)

type productModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) products.Repository {
	return &productModel{
		db: db,
	}
}

// Delete implements products.Repository.
func (*productModel) Delete(userId uint, id uint) error {
	panic("unimplemented")
}

// GetProductById implements products.Repository.
func (*productModel) GetProductById(id uint) (products.Core, error) {
	panic("unimplemented")
}

// Insert implements products.Repository.
func (*productModel) Insert(input products.Core) error {
	panic("unimplemented")
}

// SelectAll implements products.Repository.
func (*productModel) SelectAll(limit int, offset int, name string) ([]products.Core, error) {
	panic("unimplemented")
}

// Update implements products.Repository.
func (*productModel) Update(userId uint, id uint, input products.Core) error {
	panic("unimplemented")
}
