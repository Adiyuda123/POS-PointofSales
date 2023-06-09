package handler

import (
	"POS-PointofSales/features/products"

	echo "github.com/labstack/echo/v4"
)

type productController struct {
	service products.UseCase
}

func New(us products.UseCase) products.Handler {
	return &productController{
		service: us,
	}
}

// AddHandler implements products.Handler.
func (*productController) AddHandler() echo.HandlerFunc {
	panic("unimplemented")
}

// DeleteHandler implements products.Handler.
func (*productController) DeleteHandler() echo.HandlerFunc {
	panic("unimplemented")
}

// GetAllHandler implements products.Handler.
func (*productController) GetAllHandler() echo.HandlerFunc {
	panic("unimplemented")
}

// GetProductByIdHandler implements products.Handler.
func (*productController) GetProductByIdHandler() echo.HandlerFunc {
	panic("unimplemented")
}

// UpdateHandler implements products.Handler.
func (*productController) UpdateHandler() echo.HandlerFunc {
	panic("unimplemented")
}
