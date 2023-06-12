package repository

import (
	"POS-PointofSales/features/users"
)

func CoreToModel(data users.Core) User {
	return User{
		Name:     data.Name,
		Email:    data.Email,
		Phone:    data.Phone,
		Pictures: data.Pictures,
		Password: data.Password,
	}
}

func ModelToCore(data User) users.Core {

	result := users.Core{
		ID:       data.ID,
		Name:     data.Name,
		Email:    data.Email,
		Phone:    data.Phone,
		Pictures: data.Pictures,
		Password: data.Password,
	}

	return result
}

func ListUserToUserCore(user []User) []users.Core {
	var data []users.Core
	for _, v := range user {
		data = append(data, ModelToCore(v))
	}
	return data
}
