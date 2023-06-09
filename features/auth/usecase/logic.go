package usecase

import (
	"POS-PointofSales/features/auth"
	"POS-PointofSales/features/users"
	"POS-PointofSales/helper"
	"errors"
	"strings"

	"github.com/labstack/gommon/log"
)

type authLogic struct {
	data auth.Repository
}

func New(r auth.Repository) auth.UseCase {
	return &authLogic{
		data: r,
	}
}

// ChangePassword implements auth.UseCase.
func (al *authLogic) ChangePassword(id uint, oldPassword string, newPassword string, confirmPassword string, hash string) error {
	if oldPassword == "" || newPassword == "" || confirmPassword == "" {
		return errors.New("old password,new password, and confirm password cannot be empty")
	}

	if newPassword != confirmPassword {
		return errors.New("new password and confirm password must be similarity")
	}

	user, _ := al.data.GetUserByEmailOrId(".", id)
	if !helper.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New("old password not match with exist password")
	}

	return al.data.EditPassword(id, hash)
}

// LogInLogic implements auth.UseCase.
func (al *authLogic) LogInLogic(email string, password string) (users.Core, error) {
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
func (al *authLogic) RegisterUser(newUser users.Core) error {
	if err := al.data.InsertUser(newUser); err != nil {
		log.Error("error on calling register insert user query", err.Error())
		if strings.Contains(err.Error(), "column") {
			return errors.New("server error")
		} else if strings.Contains(err.Error(), "value") {
			return errors.New("invalid value")
		} else if strings.Contains(err.Error(), "too short") {
			return errors.New("invalid password length")
		}
		return errors.New("server error")
	}
	return nil
}
