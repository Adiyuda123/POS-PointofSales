package handler

import (
	"POS-PointofSales/features/products"
	"POS-PointofSales/features/users"
	"POS-PointofSales/helper"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	echo "github.com/labstack/echo/v4"
)

type productController struct {
	service  products.UseCase
	userRepo users.Repository
}

func New(service products.UseCase, services users.Repository) products.Handler {
	return &productController{
		service:  service,
		userRepo: services,
	}
}

// AddHandler implements products.Handler.
func (pc *productController) AddHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input InputRequest
		userId := helper.DecodeToken(c)
		if userId == 0 {
			c.Logger().Error("decode token is blank")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "jwt invalid", nil))
		}

		input.Name = c.FormValue("product_name")
		input.Descriptions = c.FormValue("descriptions")
		priceStr := c.FormValue("price")
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			c.Logger().Error("cannot convert price to int", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid price value", nil))
		}
		input.Price = price
		stockStr := c.FormValue("stock")
		stock, err := strconv.Atoi(stockStr)
		if err != nil {
			c.Logger().Error("cannot convert stock to int", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid stock value", nil))
		}
		input.Stock = stock
		picture, err := c.FormFile("pictures")
		if err != nil {
			if err == http.ErrMissingFile {
			} else {
				c.Logger().Error("error on retrieving uploaded file:", err.Error())
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "failed to retrieve uploaded file", nil))
			}
		}

		user, err := pc.userRepo.GetUserById(userId)
		if err != nil {
			c.Logger().Error("failed to get user information: ", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "failed to get user information", nil))
		}

		newProduct := products.Core{
			Name:         input.Name,
			Descriptions: input.Descriptions,
			Price:        input.Price,
			Pictures:     input.Pictures,
			Stock:        input.Stock,
			UserID:       userId,
			UserName:     user.Name,
		}

		copier.Copy(&newProduct, &input)

		data, err := pc.service.Add(newProduct, picture)
		if err != nil {
			c.Logger().Error("failed on calling add product log")
			if strings.Contains(err.Error(), "open") {
				c.Logger().Error("errors occurs on opening picture file")
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "error on opening picture", nil))
			} else if strings.Contains(err.Error(), "upload file in path") {
				c.Logger().Error("upload file in path are error")
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "error upload image", nil))
			} else if strings.Contains(err.Error(), "affected") {
				c.Logger().Error("no rows affected on add product")
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "data is up to date", nil))
			}

			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "internal server error", nil))
		}
		dataResponse := CoreToProductResponse(data)
		return c.JSON(helper.ResponseFormat(http.StatusOK, "successfully add product", dataResponse))
	}
}

// DeleteHandler implements products.Handler.
func (pc *productController) DeleteHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := helper.DecodeToken(c)
		if userId == 0 {
			c.Logger().Error("decode token is blank")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "jwt invalid", nil))
		}

		productId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error("cannot use path param", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusNotFound, "path invalid", nil))
		}

		if err = pc.service.Delete(userId, uint(productId)); err != nil {
			c.Logger().Error("error in calling Delete Logic")
			if strings.Contains(err.Error(), "product not found") {
				c.Logger().Error("error in calling Delete Logic, product not found")
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "product not found", nil))
			} else if strings.Contains(err.Error(), "cannot delete") {
				c.Logger().Error("error in calling Delete Logic, cannot delete")
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error in delete product", nil))
			}

			c.Logger().Error("error in calling Delete Logic, cannot delete")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error in delete product", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "successfully deleted product", nil))
	}
}

// GetAllHandler implements products.Handler.
func (pc *productController) GetAllHandler() echo.HandlerFunc {
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

		data, totaldata, err := pc.service.GetAll(limit, offset, search)
		if err != nil {
			c.Logger().Error("error occurs when calling GetAll Logic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Server Error", nil))
		}

		dataResponse := CoreToGetAllProductResponse(data)
		pagination := helper.Pagination(limit, offset, totaldata)

		return c.JSON(helper.ReponseFormatWithMeta(http.StatusOK, "Successfully retrieved product data", dataResponse, pagination))
	}
}

// GetProductByIdHandler implements products.Handler.
func (pc *productController) GetProductByIdHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID == 0 {
			c.Logger().Error("decode token is blank")
			return c.JSON(http.StatusBadRequest, "jwt invalid")
		}

		productId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error("cannot use path param", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusNotFound, "path invalid", nil))
		}

		data, err := pc.service.GetProductById(uint(productId))
		if err != nil {
			c.Logger().Error("error on calling product by id logic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error", nil))
		}
		dataResponse := CoreToProductResponse(data)
		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to get product by id", dataResponse))
	}
}

// UpdateHandler implements products.Handler.
func (pc *productController) UpdateHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var updateInput InputUpdate
		userId := helper.DecodeToken(c)
		if userId == 0 {
			c.Logger().Error("decode token is blank")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "jwt invalid", nil))
		}
		productId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error("cannot use path param", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusNotFound, "path invalid", nil))
		}

		updateInput.ID = uint(productId)
		updateInput.Name = c.FormValue("product_name")
		updateInput.Descriptions = c.FormValue("descriptions")
		priceStr := c.FormValue("price")
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			c.Logger().Error("cannot convert price to int", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid price value", nil))
		}
		updateInput.Price = price
		picture, err := c.FormFile("pictures")
		if err != nil {
			if err == http.ErrMissingFile {
			} else {
				c.Logger().Error("error on retrieving uploaded file:", err.Error())
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "failed to retrieve uploaded file", nil))
			}
		}

		updateProduct := products.Core{}
		copier.Copy(&updateProduct, &updateInput)

		if err := pc.service.Update(userId, uint(productId), updateProduct, picture); err != nil {
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

		return c.JSON(helper.ResponseFormat(http.StatusOK, "successfully updated user data", nil))
	}
}
