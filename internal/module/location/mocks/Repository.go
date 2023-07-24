// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "murakali/internal/model"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// GetCityRedis provides a mock function with given fields: ctx, provinceID
func (_m *Repository) GetCityRedis(ctx context.Context, provinceID int) (string, error) {
	ret := _m.Called(ctx, provinceID)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, int) string); ok {
		r0 = rf(ctx, provinceID)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, provinceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCostRedis provides a mock function with given fields: ctx, key
func (_m *Repository) GetCostRedis(ctx context.Context, key string) (*string, error) {
	ret := _m.Called(ctx, key)

	var r0 *string
	if rf, ok := ret.Get(0).(func(context.Context, string) *string); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCourierByID provides a mock function with given fields: ctx, courierID
func (_m *Repository) GetCourierByID(ctx context.Context, courierID string) (*model.Courier, error) {
	ret := _m.Called(ctx, courierID)

	var r0 *model.Courier
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Courier); ok {
		r0 = rf(ctx, courierID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Courier)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, courierID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProductCourierWhitelistID provides a mock function with given fields: ctx, productID
func (_m *Repository) GetProductCourierWhitelistID(ctx context.Context, productID string) ([]string, error) {
	ret := _m.Called(ctx, productID)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, string) []string); ok {
		r0 = rf(ctx, productID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, productID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProvinceRedis provides a mock function with given fields: ctx
func (_m *Repository) GetProvinceRedis(ctx context.Context) (string, error) {
	ret := _m.Called(ctx)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetShopAddress provides a mock function with given fields: ctx, userID
func (_m *Repository) GetShopAddress(ctx context.Context, userID string) (*model.Address, error) {
	ret := _m.Called(ctx, userID)

	var r0 *model.Address
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Address); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Address)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetShopByID provides a mock function with given fields: ctx, shopID
func (_m *Repository) GetShopByID(ctx context.Context, shopID string) (*model.Shop, error) {
	ret := _m.Called(ctx, shopID)

	var r0 *model.Shop
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Shop); ok {
		r0 = rf(ctx, shopID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Shop)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, shopID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetShopCourierID provides a mock function with given fields: ctx, shopID
func (_m *Repository) GetShopCourierID(ctx context.Context, shopID string) ([]string, error) {
	ret := _m.Called(ctx, shopID)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, string) []string); ok {
		r0 = rf(ctx, shopID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, shopID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubDistrictRedis provides a mock function with given fields: ctx, province, city
func (_m *Repository) GetSubDistrictRedis(ctx context.Context, province string, city string) (string, error) {
	ret := _m.Called(ctx, province, city)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, province, city)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, province, city)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUrbanRedis provides a mock function with given fields: ctx, province, city, subDistrict
func (_m *Repository) GetUrbanRedis(ctx context.Context, province string, city string, subDistrict string) (string, error) {
	ret := _m.Called(ctx, province, city, subDistrict)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) string); ok {
		r0 = rf(ctx, province, city, subDistrict)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, province, city, subDistrict)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertCityRedis provides a mock function with given fields: ctx, provinceID, value
func (_m *Repository) InsertCityRedis(ctx context.Context, provinceID int, value string) error {
	ret := _m.Called(ctx, provinceID, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) error); ok {
		r0 = rf(ctx, provinceID, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertCostRedis provides a mock function with given fields: ctx, key, value
func (_m *Repository) InsertCostRedis(ctx context.Context, key string, value string) error {
	ret := _m.Called(ctx, key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertProvinceRedis provides a mock function with given fields: ctx, value
func (_m *Repository) InsertProvinceRedis(ctx context.Context, value string) error {
	ret := _m.Called(ctx, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertSubDistrictRedis provides a mock function with given fields: ctx, province, city, value
func (_m *Repository) InsertSubDistrictRedis(ctx context.Context, province string, city string, value string) error {
	ret := _m.Called(ctx, province, city, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, province, city, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertUrbanRedis provides a mock function with given fields: ctx, province, city, subDistrict, value
func (_m *Repository) InsertUrbanRedis(ctx context.Context, province string, city string, subDistrict string, value string) error {
	ret := _m.Called(ctx, province, city, subDistrict, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) error); ok {
		r0 = rf(ctx, province, city, subDistrict, value)
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
