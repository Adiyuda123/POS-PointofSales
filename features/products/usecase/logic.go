package usecase

import (
	"POS-PointofSales/features/products"
	"POS-PointofSales/helper"
	"errors"
	"mime/multipart"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
)

type productLogic struct {
	u         products.Repository
	validator *validator.Validate
}

func New(r products.Repository) products.UseCase {
	validate := validator.New()
	validate.RegisterValidation("validImageFormat", helper.ValidImageFormat)
	return &productLogic{
		u:         r,
		validator: validate,
	}
}

// Add implements products.UseCase.
func (pl *productLogic) Add(newProduct products.Core, file *multipart.FileHeader) (products.Core, error) {
	productData := products.Core{
		ID:           newProduct.ID,
		Name:         newProduct.Name,
		Descriptions: newProduct.Descriptions,
		Price:        newProduct.Price,
		Pictures:     file.Filename,
		Stock:        newProduct.Stock,
		UserID:       newProduct.UserID,
		UserName:     newProduct.Pictures,
	}

	err := pl.validator.Struct(productData)
	if err != nil {
		log.Error("validation error:", err.Error())
		return products.Core{}, err
	}

	res, err := pl.u.Insert(newProduct, file)
	if err != nil {
		log.Error("failed on calling add product query")
		if strings.Contains(err.Error(), "open") {
			log.Error("errors occurs on opening picture file")
			return products.Core{}, errors.New("product photo is not allowed")
		} else if strings.Contains(err.Error(), "upload file in path") {
			log.Error("upload file in path error")
			return products.Core{}, errors.New("cannot upload file in path")
		} else if strings.Contains(err.Error(), "affected") {
			log.Error("no rows affected on add product")
			return products.Core{}, errors.New("data is up to date")
		}
		return products.Core{}, err
	}
	return res, nil
}

// Delete implements products.UseCase.
func (pl *productLogic) Delete(userId uint, id uint) error {
	err := pl.u.Delete(userId, id)
	if err != nil {
		log.Error("failed on calling delete product query")
		if strings.Contains(err.Error(), "finding product") {
			log.Error("error on finding product (not found)")
			return errors.New("bad request, product not found")
		} else if strings.Contains(err.Error(), "cannot delete") {
			log.Error("error on delete product")
			return errors.New("internal server error, cannot delete product")
		}
		log.Error("error in delete product (else)")
		return err
	}
	return nil
}

// GetAll implements products.UseCase.
func (pl *productLogic) GetAll(userID uint, limit, offset int, search string, fromDate, toDate time.Time) ([]products.Core, int, error) {
	result, totaldata, err := pl.u.SelectAll(userID, limit, offset, search, fromDate, toDate)
	if err != nil {
		log.Error("failed to find all product", err.Error())
		return []products.Core{}, totaldata, errors.New("internal server error")
	}

	return result, totaldata, nil
}

// GetProductById implements products.UseCase.
func (pl *productLogic) GetProductById(id uint) (products.Core, error) {
	result, err := pl.u.GetProductById(id)
	if err != nil {
		log.Error("failed to find product", err.Error())
		return products.Core{}, errors.New("internal server error")
	}

	return result, nil
}

// Update implements products.UseCase.
func (pl *productLogic) Update(userId uint, id uint, updateProduct products.Core, file *multipart.FileHeader) error {
	productData := products.Core{
		ID:           id,
		Name:         updateProduct.Name,
		Descriptions: updateProduct.Descriptions,
		Price:        updateProduct.Price,
		Pictures:     file.Filename,
		Stock:        updateProduct.Stock,
	}

	err := pl.validator.Struct(productData)
	if err != nil {
		log.Error("validation error:", err.Error())
		return err
	}

	if err := pl.u.Update(userId, id, updateProduct, file); err != nil {
		log.Error("failed on calling update product query")
		if strings.Contains(err.Error(), "open") {
			log.Error("errors occurs on opening picture file")
			return errors.New("product photo are not allowed")
		} else if strings.Contains(err.Error(), "upload file in path") {
			log.Error("upload file in path are error")
			return errors.New("cannot upload file in path")
		} else if strings.Contains(err.Error(), "affected") {
			log.Error("no rows affected on update product")
			return errors.New("data is up to date")
		}
		return err
	}
	return nil
}
