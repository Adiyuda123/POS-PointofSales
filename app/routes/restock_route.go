package routes

import (
	"POS-PointofSales/features/restocks"
	"POS-PointofSales/helper"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RestockRoutes(e *echo.Echo, pc restocks.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/restocks", pc.AddRestock(), helper.JWTMiddleware())
	e.GET("/restocks", pc.GetAllRestockHandler(), helper.JWTMiddleware())
}
