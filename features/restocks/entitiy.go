package restocks

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID        uint
	ProductID uint `validate:"required"`
	Quantity  int  `validate:"required,min=0"`
	CreatedAt time.Time
	UserID    uint
	UserName  string
}

type Handler interface {
	AddRestock() echo.HandlerFunc
	GetAllRestockHandler() echo.HandlerFunc
}

type UseCase interface {
	AddRestock(userId uint, stockInput Core) error
	GetAllRestock(userID uint, limit, offset int, search string, fromDate, toDate time.Time) ([]Core, int, error)
}

type Repository interface {
	AddRestock(userId uint, input Core) error
	SelectAllRestock(userID uint, limit, offset int, search string, fromDate, toDate time.Time) ([]Core, int, error)
}
