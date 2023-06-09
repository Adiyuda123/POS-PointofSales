package users

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID       uint
	Name     string
	Email    string
	Phone    string
	Pictures string
	Password string
}

type Handler interface {
	UserProfileHandler() echo.HandlerFunc
	UpdateProfileHandler() echo.HandlerFunc
	DeleteUserHandler() echo.HandlerFunc
}

type UseCase interface {
	UserProfileLogic(id uint) (Core, error)
	UpdateProfileLogic(id uint, name string, email string, phone string, picture *multipart.FileHeader) error
	DeleteUserLogic(id uint) error
}

type Repository interface {
	GetUserById(id uint) (Core, error)
	UpdateProfile(id uint, name string, email string, phone string, picture *multipart.FileHeader) error
	DeleteUser(id uint) error
}
