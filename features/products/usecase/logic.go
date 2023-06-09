package usecase

import (
	"POS-PointofSales/features/products"
	"mime/multipart"
)

type productLogic struct {
	u products.Repository
}

func New(r products.Repository) products.UseCase {
	return &productLogic{
		u: r,
	}
}

// Add implements products.UseCase.
func (*productLogic) Add(newProduct products.Core, file *multipart.FileHeader) error {
	panic("unimplemented")
}

// Delete implements products.UseCase.
func (*productLogic) Delete(userId uint, id uint) error {
	panic("unimplemented")
}

// GetAll implements products.UseCase.
func (*productLogic) GetAll(page int, name string) ([]products.Core, error) {
	panic("unimplemented")
}

// GetProductById implements products.UseCase.
func (*productLogic) GetProductById(id uint) (products.Core, error) {
	panic("unimplemented")
}

// Update implements products.UseCase.
func (*productLogic) Update(userId uint, id uint, updateProduct products.Core, file *multipart.FileHeader) error {
	panic("unimplemented")
}
