package repository

import (
	"POS-PointofSales/features/users"
	"POS-PointofSales/helper"
	"errors"
	"mime/multipart"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type userModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) users.Repository {
	return &userModel{
		db: db,
	}
}

// DeleteUser implements users.Repository.
func (um *userModel) DeleteUser(id uint) error {
	userToDelete := &User{}
	if err := um.db.First(userToDelete, id).Error; err != nil {
		log.Error("Error in finding user id")
		return errors.New("error in finding user")
	}

	if err := um.db.Delete(userToDelete).Error; err != nil {
		log.Error("cannot delete user")
		return err
	}

	return nil
}

// GetUserById implements users.Repository.
func (um *userModel) GetUserById(id uint) (users.Core, error) {
	var res users.Core
	if err := um.db.Table("users").Select("id, name, email, pictures").Where("id = ?", id).First(&res).Error; err != nil {
		log.Error("error occurs in finding user profile", err.Error())
		return users.Core{}, err
	}

	return res, nil
}

// UpdateProfile implements users.Repository.
func (um *userModel) UpdateProfile(id uint, name string, email string, phone string, picture *multipart.FileHeader) error {
	var UpdateUser User
	if picture != nil {
		file, err := picture.Open()
		if err != nil {
			log.Errorf("error occurred while opening picture: %v", err)
			return errors.New("failed to open picture")
		}
		defer file.Close()

		uploadURL, err := helper.UploadFile(file, "/users")
		if err != nil {
			log.Errorf("error occurred while uploading file: %v", err)
			return errors.New("failed to upload file")
		}
		UpdateUser.Pictures = uploadURL[0]
	}

	UpdateUser.ID = id
	UpdateUser.Name = name
	UpdateUser.Email = email
	UpdateUser.Phone = phone

	tx := um.db.Model(&User{}).Where("id = ?", id).Updates(&UpdateUser)
	if tx.RowsAffected < 1 {
		log.Error("there is no column to change on update user")
		return errors.New("no data affected")
	}
	if tx.Error != nil {
		log.Error("error on update user")
		return errors.New("no data updated")
	}

	return nil
}
