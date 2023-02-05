package usecase

import (
	"context"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/user/delivery/body"
	"murakali/internal/module/user/mocks"
	"murakali/pkg/postgre"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/DATA-DOG/go-sqlmock"
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
			name: "success Create Address",
			body: body.CreateAddressRequest{
				Name:          "test",
				ProvinceID:    1,
				CityID:        2,
				Province:      "Sumatera Selatan",
				City:          "Lahat",
				District:      "Lahat",
				SubDistrict:   "Pasar Lama",
				AddressDetail: "Jalan Mayor Ruslam",
				ZipCode:       "31413",
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
