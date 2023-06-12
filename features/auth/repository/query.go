package repository

import (
	"POS-PointofSales/features/auth"
	"POS-PointofSales/features/users"
	"POS-PointofSales/features/users/repository"
	"POS-PointofSales/helper"
	"errors"
	"mime/multipart"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type authModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) auth.Repository {
	return &authModel{
		db: db,
	}
}

// GetUserByEmailOrId implements auth.Repository.
func (am *authModel) GetUserByEmailOrId(email string, id uint) (users.Core, error) {
	var user repository.User
	err := am.db.Where("email = ? OR id = ?", email, id).First(&user).Error
	if err != nil {
		return users.Core{}, err
	}
	return users.Core{}, nil
}

// EditPassword implements auth.Repository.
func (am *authModel) EditPassword(id uint, oldPassword, newPassword, confirmPassword string) error {
	var user repository.User
	if err := am.db.First(&user, id).Error; err != nil {
		return err
	}

	passwordMatch := helper.CheckPasswordHash(oldPassword, user.Password)
	if !passwordMatch {
		return errors.New("old password does not match with the existing password")
	}

	if newPassword != confirmPassword {
		return errors.New("new password and confirm password must be similar")
	}

	newPasswordHashed, err := helper.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = newPasswordHashed

	if err := am.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// InsertUser implements auth.Repository.
func (am *authModel) InsertUser(newUser users.Core, picture *multipart.FileHeader) (users.Core, error) {
	inputUser := repository.User{}
	samePassword := "alta123"
	hashedPassword, err := helper.HashPassword(samePassword)
	if err != nil {
		log.Error("error occurred while hashing password:", err.Error())
		return users.Core{}, err
	}

	if picture != nil {
		file, err := picture.Open()
		if err != nil {
			log.Errorf("error occurred while opening picture: %v", err)
			return users.Core{}, errors.New("failed to open picture")
		}
		defer file.Close()

		uploadURL, err := helper.UploadFile(file, "/users")
		if err != nil {
			log.Errorf("error occurred while uploading file: %v", err)
			return users.Core{}, errors.New("failed to upload file")
		}
		inputUser.Pictures = uploadURL[0]
	}

	inputUser.Name = newUser.Name
	inputUser.Email = newUser.Email
	inputUser.Phone = newUser.Phone
	inputUser.Password = hashedPassword

	existingUser := repository.User{}
	err = am.db.Where("email = ?", newUser.Email).First(&existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("error occurred while checking duplicate email:", err.Error())
		return users.Core{}, err
	}
	if err == nil {
		return users.Core{}, errors.New("email already exists")
	}

	err = am.db.Where("phone = ?", newUser.Phone).First(&existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("error occurred while checking duplicate phone number:", err.Error())
		return users.Core{}, err
	}
	if err == nil {
		return users.Core{}, errors.New("phone number already exists")
	}

	if err := am.db.Create(&inputUser).Error; err != nil {
		log.Error("error occurred while creating user:", err.Error())
		return users.Core{}, err
	}

	return repository.ModelToCore(inputUser), nil
}

// Login implements auth.Repository.
func (am *authModel) Login(email string, password string) (users.Core, error) {
	inputUser := repository.User{}

	if email == "" {
		log.Error("email login is blank")
		return users.Core{}, errors.New("data does not exist")
	}

	if err := am.db.Where("email = ?", email).First(&inputUser).Error; err != nil {
		log.Error("error occurs on select users login", err.Error())
		return users.Core{}, err
	}

	if inputUser.ID == 1 {
		if inputUser.Password != password {
			log.Error("user input for password is wrong")
			return users.Core{}, errors.New("wrong password")
		}
	} else {
		if !helper.CheckPasswordHash(inputUser.Password, password) {
			log.Error("user input for password is wrong")
			return users.Core{}, errors.New("wrong password")
		}
	}

	return users.Core{
		ID:    inputUser.ID,
		Name:  inputUser.Name,
		Email: inputUser.Email,
		Phone: inputUser.Phone,
	}, nil
}
