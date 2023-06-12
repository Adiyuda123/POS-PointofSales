package handler

import (
	"POS-PointofSales/features/auth"
	"POS-PointofSales/features/users"
	"POS-PointofSales/helper"
	"net/http"
	"strings"

	"github.com/jinzhu/copier"
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
		userId := helper.DecodeToken(c)
		if userId == 0 {
			c.Logger().Error("decode token is blank")
			return c.JSON(http.StatusBadRequest, "jwt invalid")
		}

		req := new(ChangePasswordRequest)
		if err := c.Bind(req); err != nil {
			c.Logger().Error("error on bind change password request", err.Error())
			return c.JSON(http.StatusBadRequest, "invalid request")
		}

		oldPasswordHashed, err := helper.HashPassword(req.OldPassword)
		if err != nil {
			return err
		}

		if err := uc.service.ChangePassword(userId, oldPasswordHashed, req.NewPassword, req.ConfirmPassword); err != nil {
			c.Logger().Error("error on changing password", err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, "change password success")
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
		token, err := helper.GenerateToken(res.ID)
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
		var registerInput RegisterInput
		userID := helper.DecodeToken(c)
		if userID == 0 {
			c.Logger().Error("gagal mendekode userID dari token")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "token tidak valid", nil))
		}

		if userID != 1 {
			c.Logger().Error("pengguna bukan admin yang mencoba mengakses penambahan posisi")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "Anda bukan admin", nil))
		}

		if err := c.Bind(&registerInput); err != nil {
			c.Logger().Error("error on bind login input", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		registerInput.Name = c.FormValue("name")
		registerInput.Email = c.FormValue("email")
		registerInput.Phone = c.FormValue("phone_number")
		file, err := c.FormFile("pictures")
		if err != nil {
			if err == http.ErrMissingFile {
			} else {
				c.Logger().Error("error on retrieving uploaded file:", err.Error())
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "failed to retrieve uploaded file", nil))
			}
		}

		registerUser := users.Core{}
		copier.Copy(&registerUser, &registerInput)

		data, err := uc.service.RegisterUser(registerUser, file)
		if err != nil {
			c.Logger().Error("failed on calling updateprofile log")
			if strings.Contains(err.Error(), "open") {
				c.Logger().Error("errors occurs on opening picture file")
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "error on opening picture", nil))
			} else if strings.Contains(err.Error(), "upload file in path") {
				c.Logger().Error("upload file in path are error")
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "error upload image", nil))
			} else if strings.Contains(err.Error(), "affected") {
				c.Logger().Error("no rows affected on update user")
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "data is up to date", nil))
			}

			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "internal server error", nil))
		}
		res := RegisterResponse{}
		copier.Copy(&res, &data)

		return c.JSON(helper.ResponseFormat(http.StatusCreated, "pengguna berhasil dibuat", res))
	}
}
