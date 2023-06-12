package restocks

import (
	"github.com/labstack/echo/v4"
)

type Core struct {
	ID        uint
	ProductID uint
	Quantity  int
	Date      string
	UserID    uint
}

type Handler interface {
	AddRestock() echo.HandlerFunc
}

type UseCase interface {
	AddRestock(userId uint, stockInput Core) error
}

type Repository interface {
	AddRestock(userId uint, input Core) error
}
