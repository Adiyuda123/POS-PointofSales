package routes

import (
	"POS-PointofSales/features/products"
	"POS-PointofSales/helper"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ProductRoutes(e *echo.Echo, pc products.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/products", pc.AddHandler(), helper.JWTMiddleware())
	e.GET("/products", pc.GetAllHandler(), helper.JWTMiddleware())
	e.GET("/products/:id", pc.GetProductByIdHandler(), helper.JWTMiddleware())
	e.PUT("/products/:id", pc.UpdateHandler(), helper.JWTMiddleware())
	e.DELETE("/products/:id", pc.DeleteHandler(), helper.JWTMiddleware())
}
