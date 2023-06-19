// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	multipart "mime/multipart"

	mock "github.com/stretchr/testify/mock"

	users "POS-PointofSales/features/users"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// DeleteUser provides a mock function with given fields: id
func (_m *Repository) DeleteUser(id uint) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserById provides a mock function with given fields: id
func (_m *Repository) GetUserById(id uint) (users.Core, error) {
	ret := _m.Called(id)

	var r0 users.Core
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) (users.Core, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint) users.Core); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(users.Core)
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProfile provides a mock function with given fields: id, name, email, phone, picture
func (_m *Repository) UpdateProfile(id uint, name string, email string, phone string, picture *multipart.FileHeader) error {
	ret := _m.Called(id, name, email, phone, picture)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, string, string, string, *multipart.FileHeader) error); ok {
		r0 = rf(id, name, email, phone, picture)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}