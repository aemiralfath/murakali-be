// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	body "murakali/internal/module/location/delivery/body"

	mock "github.com/stretchr/testify/mock"
)

// UseCase is an autogenerated mock type for the UseCase type
type UseCase struct {
	mock.Mock
}

// GetCity provides a mock function with given fields: ctx, provinceID
func (_m *UseCase) GetCity(ctx context.Context, provinceID int) (*body.CityResponse, error) {
	ret := _m.Called(ctx, provinceID)

	var r0 *body.CityResponse
	if rf, ok := ret.Get(0).(func(context.Context, int) *body.CityResponse); ok {
		r0 = rf(ctx, provinceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*body.CityResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, provinceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProvince provides a mock function with given fields: ctx
func (_m *UseCase) GetProvince(ctx context.Context) (*body.ProvinceResponse, error) {
	ret := _m.Called(ctx)

	var r0 *body.ProvinceResponse
	if rf, ok := ret.Get(0).(func(context.Context) *body.ProvinceResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*body.ProvinceResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetShippingCost provides a mock function with given fields: ctx, requestBody
func (_m *UseCase) GetShippingCost(ctx context.Context, requestBody body.GetShippingCostRequest) (*body.GetShippingCostResponse, error) {
	ret := _m.Called(ctx, requestBody)

	var r0 *body.GetShippingCostResponse
	if rf, ok := ret.Get(0).(func(context.Context, body.GetShippingCostRequest) *body.GetShippingCostResponse); ok {
		r0 = rf(ctx, requestBody)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*body.GetShippingCostResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, body.GetShippingCostRequest) error); ok {
		r1 = rf(ctx, requestBody)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubDistrict provides a mock function with given fields: ctx, province, city
func (_m *UseCase) GetSubDistrict(ctx context.Context, province string, city string) (*body.SubDistrictResponse, error) {
	ret := _m.Called(ctx, province, city)

	var r0 *body.SubDistrictResponse
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *body.SubDistrictResponse); ok {
		r0 = rf(ctx, province, city)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*body.SubDistrictResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, province, city)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUrban provides a mock function with given fields: ctx, province, city, subdistrict
func (_m *UseCase) GetUrban(ctx context.Context, province string, city string, subdistrict string) (*body.UrbanResponse, error) {
	ret := _m.Called(ctx, province, city, subdistrict)

	var r0 *body.UrbanResponse
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *body.UrbanResponse); ok {
		r0 = rf(ctx, province, city, subdistrict)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*body.UrbanResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, province, city, subdistrict)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUseCase interface {
	mock.TestingT
	Cleanup(func())
}

// NewUseCase creates a new instance of UseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUseCase(t mockConstructorTestingTNewUseCase) *UseCase {
	mock := &UseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}