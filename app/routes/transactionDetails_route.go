package routes

// import (
// 	transactiondetails "POS-PointofSales/features/transactionDetails"
// 	"POS-PointofSales/helper"

// 	"github.com/labstack/echo/v4"
// 	"github.com/labstack/echo/v4/middleware"
// )

// func TransactionDetailRoutes(e *echo.Echo, pc transactiondetails.Handler) {
// 	e.Pre(middleware.RemoveTrailingSlash())

// 	e.Use(middleware.CORS())
// 	e.Use(middleware.Logger())

// 	e.POST("/transactions", pc.AddTransactions(), helper.JWTMiddleware())
// 	// e.GET("/restocks", pc.GetAllHandler(), helper.JWTMiddleware())
// }
