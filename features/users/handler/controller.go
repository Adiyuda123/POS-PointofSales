package handler

import (
	"POS-PointofSales/features/users"
	"POS-PointofSales/helper"
	"log"
	"net/http"
	"strconv"
	"strings"

	echo "github.com/labstack/echo/v4"
)

type userController struct {
	service users.UseCase
}

func New(us users.UseCase) users.Handler {
	return &userController{
		service: us,
	}
}

// DeleteUserHandler implements users.Handler.
func (uc *userController) DeleteUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := helper.DecodeToken(c)
		if userId == "" {
			c.Logger().Error("decode token is blank")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "jwt invalid", nil))
		}

		userPath, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error("cannot use path param", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusNotFound, "path invalid", nil))
		}

		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			c.Logger().Error("failed to convert userId to integer")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "jwt invalid", nil))
		}

		if uint(userIdInt) != uint(userPath) {
			c.Logger().Error("userpath is not equal with userId")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "user are not authorized to delete other user account", nil))
		}

		if err = uc.service.DeleteUserLogic(uint(userPath)); err != nil {
			c.Logger().Error("error in calling DeletUserLogic")
			if strings.Contains(err.Error(), "user not found") {
				c.Logger().Error("error in calling DeletUserLogic, user not found")
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "user not found", nil))

			} else if strings.Contains(err.Error(), "cannot delete") {
				c.Logger().Error("error in calling DeletUserLogic, cannot delete")
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error in delete user", nil))
			}

			c.Logger().Error("error in calling DeletUserLogic, cannot delete")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error in delete user", nil))

		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "success to delete user", nil))
	}
}

// UpdateProfileHandler implements users.Handler.
func (uc *userController) UpdateProfileHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var updateInput InputUpdateProfile
		userId := helper.DecodeToken(c)
		if userId == "" {
			c.Logger().Error("decode token is blank")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "jwt invalid", nil))
		}
		userPath, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error("cannot use path param", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusNotFound, "path invalid", nil))
		}

		if userId != strconv.Itoa(userPath) {
			c.Logger().Error("userpath is not equal with userId")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "user are not authorized to delete other user account", nil))
		}
		updateInput.ID = uint(userPath)
		updateInput.Name = c.FormValue("name")
		updateInput.Email = c.FormValue("email")
		updateInput.Phone = c.FormValue("phone_number")
		updateInput.Pictures, err = c.FormFile("pictures")
		if err != nil {
			log.Println("error occurs on reading form image")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "error from reading picture file", nil))
		}

		if err := uc.service.UpdateProfileLogic(updateInput.ID, updateInput.Name, updateInput.Email, updateInput.Phone, updateInput.Pictures); err != nil {
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

		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to update user data", nil))
	}
}

// UserProfileHandler implements users.Handler.
func (uc *userController) UserProfileHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var data = new(GetUserByIdResponse)
		userId := helper.DecodeToken(c)
		if userId == "" {
			c.Logger().Error("decode token is blank")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "jwt invalid", nil))
		}

		userPath, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error("cannot use path param", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusNotFound, "path invalid", nil))
		}

		result, err := uc.service.UserProfileLogic(uint(userPath))
		if err != nil {
			c.Logger().Error("error on calling userpofilelogic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error", nil))
		}

		data.ID = result.ID
		data.Email = result.Email
		data.Name = result.Name
		data.Pictures = result.Pictures

		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to check user profile", data))
	}
}
