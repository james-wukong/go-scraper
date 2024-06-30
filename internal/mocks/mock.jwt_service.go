// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	jwt "github.com/james-wukong/go-app/pkg/jwt"
	mock "github.com/stretchr/testify/mock"
)

// JWTService is an autogenerated mock type for the JWTService type
type JWTService struct {
	mock.Mock
}

// GenerateToken provides a mock function with given fields: userId, isAdmin, email, password
func (_m *JWTService) GenerateToken(userId string, isAdmin bool, email string, password string) (string, error) {
	ret := _m.Called(userId, isAdmin, email, password)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, bool, string, string) string); ok {
		r0 = rf(userId, isAdmin, email, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, bool, string, string) error); ok {
		r1 = rf(userId, isAdmin, email, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseToken provides a mock function with given fields: tokenString
func (_m *JWTService) ParseToken(tokenString string) (jwt.JwtCustomClaim, error) {
	ret := _m.Called(tokenString)

	var r0 jwt.JwtCustomClaim
	if rf, ok := ret.Get(0).(func(string) jwt.JwtCustomClaim); ok {
		r0 = rf(tokenString)
	} else {
		r0 = ret.Get(0).(jwt.JwtCustomClaim)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tokenString)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewJWTService interface {
	mock.TestingT
	Cleanup(func())
}

// NewJWTService creates a new instance of JWTService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewJWTService(t mockConstructorTestingTNewJWTService) *JWTService {
	mock := &JWTService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
