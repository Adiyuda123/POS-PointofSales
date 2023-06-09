package routes

import (
	"POS-PointofSales/features/auth"
	"POS-PointofSales/helper"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AuthRoutes(e *echo.Echo, ac auth.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/register", ac.RegisterHandler(), helper.JWTMiddleWare())
	e.POST("/login", ac.LoginHandler(), helper.JWTMiddleWare())
	e.POST("/change_password", ac.ChangePasswordHandler(), helper.JWTMiddleWare())

}
