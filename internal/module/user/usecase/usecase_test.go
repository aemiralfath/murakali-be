package usecase

import (
	"context"
	"database/sql"
	"errors"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/user/delivery/body"
	"murakali/internal/module/user/mocks"
	"murakali/pkg/postgre"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserUC_CreateAddress(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		body        body.CreateAddressRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success Create Address",
			userID: "123456",
			body: body.CreateAddressRequest{
				IsDefault:     true,
				IsShopDefault: true,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultUserAddress", mock.Anything, mock.Anything).Return(&model.Address{IsDefault: true}, nil)
				r.On("UpdateDefaultAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetDefaultShopAddress", mock.Anything, mock.Anything).Return(&model.Address{IsShopDefault: true}, nil)
				r.On("UpdateDefaultShopAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("CreateAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "error user not found sql no rows",
			userID: "123456",
			body:   body.CreateAddressRequest{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New("User not exist."),
		},
		{
			name:   "error user sql",
			userID: "123456",
			body:   body.CreateAddressRequest{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "error address is default",
			userID: "123456",
			body:   body.CreateAddressRequest{IsDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultUserAddress", mock.Anything, mock.Anything).Return(nil, errors.New("test"))

			},
			expectedErr: nil,
		},
		{
			name:   "error update address isdefault",
			userID: "123456",
			body:   body.CreateAddressRequest{IsDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultUserAddress", mock.Anything, mock.Anything).Return(&model.Address{IsDefault: true}, nil)
				r.On("UpdateDefaultAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))

			},
			expectedErr: nil,
		},
		{
			name:   "error address is shop default",
			userID: "123456",
			body:   body.CreateAddressRequest{IsShopDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultShopAddress", mock.Anything, mock.Anything).Return(nil, errors.New("test"))

			},
			expectedErr: nil,
		},
		{
			name:   "error update address is shop default",
			userID: "123456",
			body:   body.CreateAddressRequest{IsShopDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultShopAddress", mock.Anything, mock.Anything).Return(&model.Address{IsShopDefault: true}, nil)
				r.On("UpdateDefaultShopAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))

			},
			expectedErr: nil,
		},
		{
			name:   "error update address is shop default",
			userID: "123456",
			body:   body.CreateAddressRequest{IsShopDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultShopAddress", mock.Anything, mock.Anything).Return(nil, nil)
				r.On("CreateAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))

			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.CreateAddress(context.Background(), tc.userID, tc.body)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_UpdateAddressByID(t *testing.T) {

	testCase := []struct {
		name        string
		userID      string
		addressID   string
		body        body.UpdateAddressRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:      "success Create Address",
			userID:    "123456",
			addressID: "123456",
			body: body.UpdateAddressRequest{
				IsDefault:     true,
				IsShopDefault: true,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetDefaultUserAddress", mock.Anything, mock.Anything).Return(&model.Address{IsDefault: true}, nil)
				r.On("UpdateDefaultAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetDefaultShopAddress", mock.Anything, mock.Anything).Return(&model.Address{IsShopDefault: true}, nil)
				r.On("UpdateDefaultShopAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "error user not found sql no rows",
			userID: "123456",
			body:   body.UpdateAddressRequest{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New("User not exist."),
		},
		{
			name:   "error user sql",
			userID: "123456",
			body:   body.UpdateAddressRequest{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "error address is default",
			userID: "123456",
			body:   body.UpdateAddressRequest{IsDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultUserAddress", mock.Anything, mock.Anything).Return(nil, errors.New("test"))

			},
			expectedErr: nil,
		},
		{
			name:   "error update address isdefault",
			userID: "123456",
			body:   body.UpdateAddressRequest{IsDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultUserAddress", mock.Anything, mock.Anything).Return(&model.Address{IsDefault: true}, nil)
				r.On("UpdateDefaultAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))

			},
			expectedErr: nil,
		},
		{
			name:   "error address is shop default",
			userID: "123456",
			body:   body.UpdateAddressRequest{IsShopDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultShopAddress", mock.Anything, mock.Anything).Return(nil, errors.New("test"))

			},
			expectedErr: nil,
		},
		{
			name:   "error update address is shop default",
			userID: "123456",
			body:   body.UpdateAddressRequest{IsShopDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultShopAddress", mock.Anything, mock.Anything).Return(&model.Address{IsShopDefault: true}, nil)
				r.On("UpdateDefaultShopAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))

			},
			expectedErr: nil,
		},
		{
			name:   "error update address is shop default",
			userID: "123456",
			body:   body.UpdateAddressRequest{IsShopDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultShopAddress", mock.Anything, mock.Anything).Return(nil, nil)
				r.On("CreateAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))

			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.UpdateAddressByID(context.Background(), tc.userID, tc.addressID, tc.body)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
