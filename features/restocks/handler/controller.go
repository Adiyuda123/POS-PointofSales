package handler

import (
	"POS-PointofSales/features/restocks"
	"POS-PointofSales/helper"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type restockController struct {
	service restocks.UseCase
}

func New(service restocks.UseCase) restocks.Handler {
	return &restockController{
		service: service,
	}
}

// AddRestock implements restocks.Handler.
func (pc *restockController) AddRestock() echo.HandlerFunc {
	return func(c echo.Context) error {
		var updateInput InputRestock

		userId := helper.DecodeToken(c)
		if userId == 0 {
			c.Logger().Error("empty decode tokens")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "jwt invalid", nil))
		}

		if err := c.Bind(&updateInput); err != nil {
			c.Logger().Error("failed to bind JSON data", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid request payload", nil))
		}

		if updateInput.ProductID == 0 {
			c.Logger().Error("Product ID is required")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Product ID is required", nil))
		}

		if updateInput.Stock <= 0 {
			c.Logger().Error("invalid stock value")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid stock value", nil))
		}

		productId := int(updateInput.ProductID)

		restockInput := restocks.Core{
			ProductID: uint(productId),
			Quantity:  updateInput.Stock,
		}

		if err := pc.service.AddRestock(userId, restockInput); err != nil {
			c.Logger().Error("failed to call restock log")
			if strings.Contains(err.Error(), "affected") {
				c.Logger().Error("no rows are affected on restock update")
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "the data has been updated", nil))
			}

			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "internal server error", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "product stock updated successfully", nil))
	}
}
