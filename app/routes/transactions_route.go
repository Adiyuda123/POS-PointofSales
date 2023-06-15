package routes

import (
	transactions "POS-PointofSales/features/transactions"
	"POS-PointofSales/helper"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func TransactionRoutes(e *echo.Echo, pc transactions.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/payments", pc.CreateTransactions(), helper.JWTMiddleware())
	e.POST("/transactions", pc.AddTransactions(), helper.JWTMiddleware())
	// e.GET("/restocks", pc.GetAllHandler(), helper.JWTMiddleware())
}
