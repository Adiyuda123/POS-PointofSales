package repository

import (
	"POS-PointofSales/features/products"
	"POS-PointofSales/helper"
	"errors"
	"mime/multipart"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type productModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) products.Repository {
	return &productModel{
		db: db,
	}
}

// Delete implements products.Repository.
func (pm *productModel) Delete(userId uint, id uint) error {
	productToDelete := &Product{}
	if err := pm.db.First(productToDelete, id).Error; err != nil {
		log.Error("Error in finding product")
		return errors.New("product not found")
	}

	if productToDelete.UserID != userId && userId != 1 {
		log.Error("Unauthorized to delete product")
		return errors.New("unauthorized to delete product")
	}

	if err := pm.db.Delete(productToDelete).Error; err != nil {
		log.Error("Failed to delete product")
		return err
	}

	return nil
}

// GetProductById implements products.Repository.
func (pm *productModel) GetProductById(id uint) (products.Core, error) {
	res := Product{}
	if err := pm.db.Table("products").
		Select("products.id, products.name, products.descriptions, products.price, products.pictures, products.stock, users.name as user_name").
		Joins("JOIN users ON products.user_id = users.id").
		Where("products.id = ?", id).
		First(&res).Error; err != nil {
		log.Error("error occurs in finding product by id:", err.Error())
		return products.Core{}, err
	}

	return ModelToCore(res), nil
}

// Insert implements products.Repository.
func (pm *productModel) Insert(newProduct products.Core, file *multipart.FileHeader) (products.Core, error) {
	inputProduct := CoreToModel(newProduct)

	if file != nil {
		file, err := file.Open()
		if err != nil {
			log.Errorf("error occurred while opening picture: %v", err)
			return products.Core{}, errors.New("failed to open picture")
		}
		defer file.Close()

		uploadURL, err := helper.UploadFile(file, "/products")
		if err != nil {
			log.Errorf("error occurred while uploading file: %v", err)
			return products.Core{}, errors.New("failed to upload file")
		}
		inputProduct.Pictures = uploadURL[0]
	}

	if err := pm.db.Create(&inputProduct).Error; err != nil {
		log.Error("error occurred while add product:", err.Error())
		return products.Core{}, err
	}

	return ModelToCore(inputProduct), nil
}

// SelectAll implements products.Repository.
func (pm *productModel) SelectAll(limit int, offset int, name string) ([]products.Core, int, error) {
	nameSearch := "%" + name + "%"
	totalData := int64(-1)
	var productModel []Product

	query := pm.db.Table("products").
		Select("products.id, products.name, products.descriptions, products.price, products.pictures, products.stock, users.name AS user_name").
		Joins("JOIN users ON products.user_id = users.id").
		Group("products.id").Find(&productModel).
		Limit(limit).Offset(offset).
		Order("id DESC")

	if name != "" {
		if err := query.Where("products.name LIKE ? OR products.descriptions LIKE ? OR products.price LIKE ? OR products.stock LIKE ? OR users.name LIKE ?",
			nameSearch, nameSearch, nameSearch, nameSearch, nameSearch).Find(&productModel).Error; err != nil {
			log.Errorf("error on finding search: %w", err)
			return []products.Core{}, int(totalData), err
		}
		if err := pm.db.Table("products").
			Joins("JOIN users ON products.user_id = users.id").
			Where("products.name LIKE ? OR products.descriptions LIKE ? OR products.price LIKE ? OR products.stock LIKE ? OR users.name LIKE ?",
				nameSearch, nameSearch, nameSearch, nameSearch, nameSearch).
			Count(&totalData).Error; err != nil {
			log.Errorf("error on count filtered data: %w", err)
			return []products.Core{}, int(totalData), err
		}
	} else {
		if err := query.Find(&productModel).Error; err != nil {
			log.Errorf("error on finding data without search: %w", err)
			return []products.Core{}, int(totalData), err
		}
		if err := pm.db.Table("products").Count(&totalData).Error; err != nil {
			log.Errorf("error on counting data without search: %w", err)
			return []products.Core{}, int(totalData), err
		}
	}

	productCoreAll := ListproductToproductCore(productModel)
	return productCoreAll, int(totalData), nil
}

// Update implements products.Repository.
func (pm *productModel) Update(userId uint, id uint, input products.Core, file *multipart.FileHeader) error {
	data := CoreToModel(input)

	productToUpdate := &Product{}
	if err := pm.db.First(productToUpdate, "id = ? AND user_id = ?", id, userId).Error; err != nil {
		log.Error("Error in finding product")
		return errors.New("product not found or not owned by the user")
	}

	if file != nil {
		file, err := file.Open()
		if err != nil {
			log.Errorf("Error occurred while opening picture: %v", err)
			return errors.New("failed to open picture")
		}
		defer file.Close()

		uploadURL, err := helper.UploadFile(file, "/products")
		if err != nil {
			log.Errorf("Error occurred while uploading file: %v", err)
			return errors.New("failed to upload file")
		}
		data.Pictures = uploadURL[0]
	}

	tx := pm.db.Model(&Product{}).Where("id = ? AND user_id = ?", id, userId).Updates(&data)
	if tx.RowsAffected < 1 {
		log.Error("Failed to update product")
		return errors.New("product not updated")
	}
	if tx.Error != nil {
		log.Error("Product not found")
		return tx.Error
	}

	return nil
}
