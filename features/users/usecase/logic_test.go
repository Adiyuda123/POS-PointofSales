package usecase_test

import (
	"POS-PointofSales/features/users"
	"POS-PointofSales/features/users/mocks"
	"POS-PointofSales/features/users/usecase"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateProfile(t *testing.T) {
	repo := mocks.NewRepository(t)
	ul := usecase.New(repo)

	t.Run("Success Update Profile", func(t *testing.T) {
		file := &multipart.FileHeader{
			Filename: "example.jpg",
		}
		id := uint(1)
		name := "joko"
		email := "joko@mail.com"
		phone := "0812345678910"
		repo.On("UpdateProfile", id, name, email, phone, file).
			Return(nil).Once()

		err := ul.UpdateProfileLogic(id, name, email, phone, file)
		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})

}

func TestUserProfile(t *testing.T) {
	repo := mocks.NewRepository(t)
	ul := usecase.New(repo)

	t.Run("Success Update Profile", func(t *testing.T) {
		id := uint(1)
		resultRepo := users.Core{
			ID:       id,
			Name:     "Jarwo",
			Email:    "jarwo@mail.com",
			Phone:    "0812345678910",
			Pictures: "jarwo.jpg",
			Password: "123456",
		}

		repo.On("GetUserById", id).Return(resultRepo, nil).Once()
		resultRepo, err := ul.UserProfileLogic(id)
		assert.NoError(t, err)
		assert.Equal(t, resultRepo, resultRepo)
		repo.AssertExpectations(t)
	})

	t.Run("Error Update Profile", func(t *testing.T) {
		id := uint(1)
		resultRepo := users.Core{}

		repo.On("GetUserById", id).Return(resultRepo, errors.New("internal server error")).Once()
		result, err := ul.UserProfileLogic(id)
		assert.Error(t, err)
		assert.Equal(t, resultRepo, result)
		repo.AssertExpectations(t)
	})
}

func TestDeleteProfile(t *testing.T) {
	repo := mocks.NewRepository(t)
	ul := usecase.New(repo)

	t.Run("Success Delete Profile", func(t *testing.T) {
		id := uint(1)

		repo.On("DeleteUser", id).Return(nil).Once()
		err := ul.DeleteUserLogic(id)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Error Finding User", func(t *testing.T) {
		id := uint(1)

		repo.On("DeleteUser", id).Return(errors.New("finding user")).Once()
		err := ul.DeleteUserLogic(id)
		assert.Error(t, err)
		assert.Equal(t, errors.New("bad request, user not found"), err)
		repo.AssertExpectations(t)
	})
}
