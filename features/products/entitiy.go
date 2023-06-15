package products

import (
	"mime/multipart"
	"time"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID           uint
	Name         string `validate:"omitempty"`
	Descriptions string `validate:"omitempty"`
	Price        int    `validate:"omitempty"`
	Pictures     string `validate:"validImageFormat"`
	Stock        int    `validate:"omitempty"`
	UserID       uint
	UserName     string
	CreatedAt    time.Time
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
	GetAll(userID uint, limit, offset int, search string, fromDate, toDate time.Time) ([]Core, int, error)
	GetProductById(id uint) (Core, error)
	Update(userId uint, id uint, updateProduct Core, file *multipart.FileHeader) error
	Delete(userId uint, id uint) error
}

type Repository interface {
	Insert(newProduct Core, file *multipart.FileHeader) (Core, error)
	SelectAll(userID uint, limit, offset int, search string, fromDate, toDate time.Time) ([]Core, int, error)
	GetProductById(id uint) (Core, error)
	Update(userId uint, id uint, input Core, file *multipart.FileHeader) error
	Delete(userId uint, id uint) error
}
