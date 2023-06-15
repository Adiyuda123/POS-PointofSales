package usecase

import (
	"POS-PointofSales/features/auth"
	"POS-PointofSales/features/users"
	"POS-PointofSales/helper"
	"errors"
	"mime/multipart"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
)

type authLogic struct {
	data      auth.Repository
	validator *validator.Validate
}

func New(r auth.Repository) auth.UseCase {
	validate := validator.New()
	validate.RegisterValidation("validImageFormat", helper.ValidImageFormat)
	return &authLogic{
		data:      r,
		validator: validate,
	}
}

// ChangePassword implements auth.UseCase.
func (al *authLogic) ChangePassword(id uint, oldPassword string, newPassword string, confirmPassword string) error {
	oldPasswordHashed, err := helper.HashPassword(oldPassword)
	if err != nil {
		return err
	}
	user, err := al.data.GetUserByEmailOrId(".", id)
	if err != nil {
		return err
	}

	if match := helper.CheckPasswordHash(oldPasswordHashed, user.Password); !match {
		return errors.New("old password does not match with the existing password")
	}

	if newPassword != confirmPassword {
		return errors.New("new password and confirm password must be similar")
	}

	newPasswordHashed, err := helper.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = newPasswordHashed

	if err := al.data.EditPassword(id, oldPassword, newPassword, confirmPassword); err != nil {
		log.Error("error on loginlogic, internal server error", err.Error())
		return errors.New("internal server error")
	}

	return nil
}

// LogInLogic implements auth.UseCase.
func (al *authLogic) LogInLogic(email string, password string) (users.Core, error) {
	// loginData := users.Core{
	// 	Email:    email,
	// 	Password: password,
	// }

	// err := al.validator.Struct(loginData)
	// if err != nil {
	// 	log.Error("validation error:", err.Error())
	// 	return users.Core{}, err
	// }

	res, err := al.data.Login(email, password)
	if err != nil {

		if strings.Contains(err.Error(), "not exist") {
			return users.Core{}, errors.New("email cannot be blank")

		} else if strings.Contains(err.Error(), "wrong") {
			return users.Core{}, errors.New("password is wrong")

		}
		log.Error("error on loginlogic, internal server error", err.Error())
		return users.Core{}, errors.New("internal server error")

	}

	return res, nil
}

// RegisterUser implements auth.UseCase.
func (al *authLogic) RegisterUser(newUser users.Core, picture *multipart.FileHeader) (users.Core, error) {
	registerData := users.Core{
		Name:     newUser.Name,
		Email:    newUser.Email,
		Phone:    newUser.Phone,
		Pictures: picture.Filename,
	}

	err := al.validator.Struct(registerData)
	if err != nil {
		log.Error("validation error:", err.Error())
		return users.Core{}, err
	}

	res, err := al.data.InsertUser(newUser, picture)
	if err != nil {
		log.Error("error on calling register insert user query", err.Error())
		if strings.Contains(err.Error(), "column") {
			return users.Core{}, errors.New("server error")
		} else if strings.Contains(err.Error(), "value") {
			return users.Core{}, errors.New("invalid value")
		} else if strings.Contains(err.Error(), "too short") {
			return users.Core{}, errors.New("invalid password length")
		}
		return users.Core{}, errors.New("server error")
	}
	return res, nil
}
