package usecase

import (
	"POS-PointofSales/features/users"
	"errors"
	"mime/multipart"
	"strings"

	"github.com/labstack/gommon/log"
)

type userLogic struct {
	u users.Repository
}

func New(r users.Repository) users.UseCase {
	return &userLogic{
		u: r,
	}
}

// DeleteUserLogic implements users.UseCase.
func (ul *userLogic) DeleteUserLogic(id uint) error {
	err := ul.u.DeleteUser(id)
	if err != nil {
		log.Error("failed on calling deleteuser query")
		if strings.Contains(err.Error(), "finding user") {
			log.Error("error on finding user (not found)")
			return errors.New("bad request, user not found")
		} else if strings.Contains(err.Error(), "cannot delete") {
			log.Error("error on delete user")
			return errors.New("internal server error, cannot delete user")
		}
		log.Error("error in delete user (else)")
		return err
	}
	return nil
}

// UpdateProfileLogic implements users.UseCase.
func (ul *userLogic) UpdateProfileLogic(id uint, name string, email string, phone string, picture *multipart.FileHeader) error {
	if err := ul.u.UpdateProfile(id, name, email, phone, picture); err != nil {
		log.Error("failed on calling updateprofile query")
		if strings.Contains(err.Error(), "open") {
			log.Error("errors occurs on opening picture file")
			return errors.New("user photo are not allowed")
		} else if strings.Contains(err.Error(), "upload file in path") {
			log.Error("upload file in path are error")
			return errors.New("cannot upload file in path")
		} else if strings.Contains(err.Error(), "hashing password") {
			log.Error("hashing password error")
			return errors.New("is invalid")
		} else if strings.Contains(err.Error(), "affected") {
			log.Error("no rows affected on update user")
			return errors.New("data is up to date")
		}
		return err
	}
	return nil
}

// UserProfileLogic implements users.UseCase.
func (ul *userLogic) UserProfileLogic(id uint) (users.Core, error) {
	result, err := ul.u.GetUserById(id)
	if err != nil {
		log.Error("failed to find user", err.Error())
		return users.Core{}, errors.New("internal server error")
	}

	return result, nil
}
