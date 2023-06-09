package routes

import (
	"POS-PointofSales/features/users"
	"POS-PointofSales/helper"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func UserRoutes(e *echo.Echo, uc users.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/users", uc.UserProfileHandler(), helper.JWTMiddleWare())
	e.PUT("/users/:id", uc.UpdateProfileHandler(), helper.JWTMiddleWare())
	e.DELETE("/users/:id", uc.DeleteUserHandler(), helper.JWTMiddleWare())
}
