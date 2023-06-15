package handler

import (
	"POS-PointofSales/features/restocks"
	// "POS-PointofSales/features/users"
	"POS-PointofSales/helper"
	"net/http"
	"strconv"
	"strings"
	"time"

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

// GetAllRestockHandler implements restocks.Handler.
func (pc *restockController) GetAllRestockHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID == 0 {
			c.Logger().Error("decode token is blank")
			return c.JSON(http.StatusBadRequest, "jwt invalid")
		}

		search := c.QueryParam("search")

		limitStr := c.QueryParam("limit")
		offsetStr := c.QueryParam("offset")

		limit := 10
		if limitStr != "" {
			limitInt, err := strconv.Atoi(limitStr)
			if err != nil {
				c.Logger().Errorf("limit is not a number: %s", limitStr)
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid limit value", nil))
			}
			limit = limitInt
		}

		offset := 0
		if offsetStr != "" {
			offsetInt, err := strconv.Atoi(offsetStr)
			if err != nil {
				c.Logger().Errorf("offset is not a number: %s", offsetStr)
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid offset value", nil))
			}
			offset = offsetInt
		}

		fromDateStr := c.QueryParam("from_date")
		toDateStr := c.QueryParam("to_date")

		fromDate := time.Time{}
		if fromDateStr != "" {
			var err error
			fromDate, err = time.Parse("2006-01-02", fromDateStr)
			if err != nil {
				c.Logger().Errorf("from_date is not a valid date: %s", fromDateStr)
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid from_date value", nil))
			}
		}

		toDate := time.Now()
		if toDateStr != "" {
			var err error
			toDate, err = time.Parse("2006-01-02", toDateStr)
			if err != nil {
				c.Logger().Errorf("to_date is not a valid date: %s", toDateStr)
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid to_date value", nil))
			}
			toDate = toDate.AddDate(0, 0, 1)
		}

		data, totalData, err := pc.service.GetAllRestock(userID, limit, offset, search, fromDate, toDate)
		if err != nil {
			c.Logger().Error("error occurs when calling GetAll Logic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Server Error", nil))
		}

		dataResponse := CoreToGetAllRestockResponse(data)
		pagination := helper.Pagination(limit, offset, totalData)

		return c.JSON(helper.ReponseFormatWithMeta(http.StatusOK, "Successfully retrieved restock data", dataResponse, pagination))
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
