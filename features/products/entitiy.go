package products

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID           uint
	Name         string
	Descriptions string
	Price        int
	Pictures     string
	Stock        int
	UserID       uint
	UserName     string
}

type Handler interface {
	AddHandler() echo.HandlerFunc
	GetAllHandler() echo.HandlerFunc
	GetProductByIdHandler() echo.HandlerFunc
	UpdateHandler() echo.HandlerFunc
	DeleteHandler() echo.HandlerFunc
}

type UseCase interface {
	Add(newProduct Core, file *multipart.FileHeader) (Core, error)
	GetAll(limit, offset int, name string) ([]Core, int, error)
	GetProductById(id uint) (Core, error)
	Update(userId uint, id uint, updateProduct Core, file *multipart.FileHeader) error
	Delete(userId uint, id uint) error
}

type Repository interface {
	Insert(newProduct Core, file *multipart.FileHeader) (Core, error)
	SelectAll(limit, offset int, name string) ([]Core, int, error)
	GetProductById(id uint) (Core, error)
	Update(userId uint, id uint, input Core, file *multipart.FileHeader) error
	Delete(userId uint, id uint) error
}
