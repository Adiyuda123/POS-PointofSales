package repository

import (
	"POS-PointofSales/features/auth"
	"POS-PointofSales/features/users"
	"POS-PointofSales/features/users/repository"
	"POS-PointofSales/helper"
	"errors"

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
	if err := am.db.Where("email = ? or id = ?", email, id).First(&user); err.Error != nil {
		return users.Core{}, err.Error
	}

	return repository.ModelToCore(user), nil
}

// EditPassword implements auth.Repository.
func (am *authModel) EditPassword(id uint, newPassword string) error {
	var user repository.User
	if err := am.db.Model(&user).Where("id", id).Update("password", newPassword); err.Error != nil {
		return err.Error
	}
	return nil
}

// InsertUser implements auth.Repository.
func (am *authModel) InsertUser(newUser users.Core) error {
	inputUser := repository.User{}
	samePassword := "alta123"
	hashedPassword, err := helper.HashPassword(samePassword)
	if err != nil {
		log.Error("error occurs on hashing password", err.Error())
		return err
	}

	inputUser.Name = newUser.Name
	inputUser.Email = newUser.Email
	inputUser.Phone = newUser.Phone
	inputUser.Pictures = newUser.Pictures
	inputUser.Password = hashedPassword

	existingUser := repository.User{}
	if err := am.db.Where("email = ?", newUser.Email).First(&existingUser).Error; err == nil {
		return errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("error occurred while checking duplicate email", err.Error())
		return err
	}

	if err := am.db.Where("phone = ?", newUser.Phone).First(&existingUser).Error; err == nil {
		return errors.New("phone number already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("error occurred while checking duplicate phone number", err.Error())
		return err
	}

	if err := am.db.Table("users").Create(&inputUser).Error; err != nil {
		log.Error("error on create table users", err.Error())
		return err
	}

	return nil
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

	if !helper.CheckPasswordHash(inputUser.Password, password) {
		log.Error("user input for password is wrong")
		return users.Core{}, errors.New("wrong password")
	}

	return users.Core{
		ID:    inputUser.Model.ID,
		Name:  inputUser.Name,
		Email: inputUser.Email,
		Phone: inputUser.Phone,
	}, nil
}
