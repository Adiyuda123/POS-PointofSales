package handler

import (
	"POS-PointofSales/features/auth"
	"POS-PointofSales/features/users"
	"POS-PointofSales/helper"
	"net/http"
	"strconv"
	"strings"

	echo "github.com/labstack/echo/v4"
)

type authController struct {
	service auth.UseCase
}

func New(au auth.UseCase) auth.Handler {
	return &authController{
		service: au,
	}
}

// ChangePasswordHandler implements auth.Handler.
func (uc *authController) ChangePasswordHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userIdString := helper.DecodeToken(c)
		if userIdString == "" {
			c.Logger().Error("decode token is blank")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "jwt invalid", nil))
		}

		userId, err := strconv.Atoi(userIdString)
		if err != nil {
			c.Logger().Error("error converting userId to int", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid userId", nil))
		}

		r := ChangePasswordRequest{}
		if err := c.Bind(&r); err != nil {
			c.Logger().Error("error on bind login input", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		hash, _ := helper.HashPassword(r.NewPassword)
		if err := uc.service.ChangePassword(uint(userId), r.OldPassword, r.NewPassword, r.ConfirmPassword, hash); err != nil {
			c.Logger().Error("error on bind login input", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "change password Success", nil))
	}
}

// LoginHandler implements auth.Handler.
func (uc *authController) LoginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input LoginInput
		if err := c.Bind(&input); err != nil {
			c.Logger().Error("error on bind login input", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		res, err := uc.service.LogInLogic(input.Email, input.Password)
		if err != nil {
			c.Logger().Error("error on calling Login Logic", err.Error())

			if strings.Contains(err.Error(), "not exist") {
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "email does not exist, please sign up", nil))
			} else if strings.Contains(err.Error(), "wrong") {
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "password is wrong please try again", nil))
			} else if strings.Contains(err.Error(), "blank") {
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "email is blank please try again", nil))
			}
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "internal server error", nil))
		}
		token, err := helper.GenerateToken(strconv.Itoa(int(res.ID)))
		if err != nil {
			c.Logger().Error("error on generation token", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Internal server error", nil))
		}
		var data = new(LoginResponse)
		data.ID = res.ID
		data.Name = res.Name
		data.Token = token

		return c.JSON(helper.ResponseFormat(http.StatusOK, "success login!", data))
	}
}

// RegisterHandler implements auth.Handler.
func (uc *authController) RegisterHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user is not admin trying to access add position")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not an admin", nil))
		}

		name := c.FormValue("name")
		email := c.FormValue("email")
		phone := c.FormValue("phone_number")
		pictures := c.FormValue("pictures")

		err := uc.service.RegisterUser(users.Core{
			Name:     name,
			Email:    email,
			Phone:    phone,
			Pictures: pictures,
		})
		if err != nil {
			c.Logger().Error("error on calling RegisterUser", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, err.Error(), nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusCreated, "successfully created user", nil))
	}
}
