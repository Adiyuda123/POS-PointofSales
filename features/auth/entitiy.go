package auth

import (
	"POS-PointofSales/features/users"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	RegisterHandler() echo.HandlerFunc
	LoginHandler() echo.HandlerFunc
	ChangePasswordHandler() echo.HandlerFunc
}

type UseCase interface {
	RegisterUser(newUser users.Core) error
	LogInLogic(email string, password string) (users.Core, error)
	ChangePassword(id uint, oldPassword, newPassword, confirmPassword, hash string) error
}

type Repository interface {
	InsertUser(newUser users.Core) error
	Login(email string, password string) (users.Core, error)
	EditPassword(id uint, newPassword string) error
	GetUserByEmailOrId(email string, id uint) (users.Core, error)
}
