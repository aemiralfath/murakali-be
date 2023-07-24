package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"murakali/config"
	"murakali/internal/constant"
	"murakali/internal/model"
	body2 "murakali/internal/module/location/delivery/body"
	"murakali/internal/module/user/delivery/body"
	"murakali/internal/module/user/mocks"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
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
			name:      "error user not found sql no rows",
			userID:    "123456",
			addressID: "",
			body:      body.UpdateAddressRequest{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New("User not exist."),
		},
		{
			name:      "error user sql",
			userID:    "123456",
			addressID: "",
			body:      body.UpdateAddressRequest{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:      "error address sql no rows",
			userID:    "123456",
			addressID: "123456",
			body:      body.UpdateAddressRequest{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New("Address not exist."),
		},
		{
			name:      "error address sql",
			userID:    "123456",
			addressID: "123456",
			body:      body.UpdateAddressRequest{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:      "error address is default",
			userID:    "123456",
			addressID: "123456",
			body:      body.UpdateAddressRequest{IsDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetDefaultUserAddress", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:      "error update address isdefault",
			userID:    "123456",
			addressID: "123456",
			body:      body.UpdateAddressRequest{IsDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetDefaultUserAddress", mock.Anything, mock.Anything).Return(&model.Address{IsDefault: true}, nil)
				r.On("UpdateDefaultAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))

			},
			expectedErr: errors.New("test"),
		},
		{
			name:      "error address is shop default",
			userID:    "123456",
			addressID: "123456",
			body:      body.UpdateAddressRequest{IsShopDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetDefaultShopAddress", mock.Anything, mock.Anything).Return(nil, errors.New("test"))

			},
			expectedErr: errors.New("test"),
		},
		{
			name:      "error update address is shop default",
			userID:    "123456",
			addressID: "123456",
			body:      body.UpdateAddressRequest{IsShopDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetDefaultShopAddress", mock.Anything, mock.Anything).Return(&model.Address{IsShopDefault: true}, nil)
				r.On("UpdateDefaultShopAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))

			},
			expectedErr: errors.New("test"),
		},
		{
			name:      "error update address is shop default",
			addressID: "123456",
			userID:    "123456",
			body:      body.UpdateAddressRequest{IsShopDefault: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetDefaultShopAddress", mock.Anything, mock.Anything).Return(nil, nil)
				r.On("UpdateAddress", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))

			},
			expectedErr: errors.New("test"),
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

func Test_userUC_GetAddress(t *testing.T) {
	testCase := []struct {
		name         string
		userID       string
		pgn          *pagination.Pagination
		queryRequest *body.GetAddressQueryRequest
		mock         func(t *testing.T, r *mocks.Repository)
		expectedErr  error
	}{
		{
			name:         "success Get Address Default:ShopDefault - false:false",
			userID:       "123456",
			pgn:          &pagination.Pagination{},
			queryRequest: &body.GetAddressQueryRequest{IsDefaultBool: false, IsShopDefaultBool: false},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(10)
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetTotalAddress", mock.Anything, mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("GetAllAddresses", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*model.Address{}, nil)
			},
			expectedErr: nil,
		},
		{
			name:         "Error user sql no rows",
			userID:       "123456",
			pgn:          &pagination.Pagination{},
			queryRequest: &body.GetAddressQueryRequest{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New("Invalid Credentials."),
		},
		{
			name:         "Error user",
			userID:       "123456",
			pgn:          &pagination.Pagination{},
			queryRequest: &body.GetAddressQueryRequest{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:         "Errror Totalrows",
			userID:       "123456",
			pgn:          &pagination.Pagination{},
			queryRequest: &body.GetAddressQueryRequest{IsDefaultBool: false, IsShopDefaultBool: false},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(0)
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetTotalAddress", mock.Anything, mock.Anything, mock.Anything).Return(tempInt64, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:         "Errror GetAllAddresses",
			userID:       "123456",
			pgn:          &pagination.Pagination{},
			queryRequest: &body.GetAddressQueryRequest{IsDefaultBool: false, IsShopDefaultBool: false},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(10)
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetTotalAddress", mock.Anything, mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("GetAllAddresses", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:         "success Get Address Default:ShopDefault - true:false",
			userID:       "123456",
			pgn:          &pagination.Pagination{},
			queryRequest: &body.GetAddressQueryRequest{IsDefaultBool: true, IsShopDefaultBool: false},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultUserAddress", mock.Anything, mock.Anything).Return(&model.Address{}, nil)
			},
			expectedErr: nil,
		},
		{
			name:         "Error GetDefaultUserAddress no sql row",
			userID:       "123456",
			pgn:          &pagination.Pagination{},
			queryRequest: &body.GetAddressQueryRequest{IsDefaultBool: true, IsShopDefaultBool: false},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultUserAddress", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New("Default address not found."),
		},
		{
			name:         "Error GetDefaultUserAddress",
			userID:       "123456",
			pgn:          &pagination.Pagination{},
			queryRequest: &body.GetAddressQueryRequest{IsDefaultBool: true, IsShopDefaultBool: false},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultUserAddress", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:         "success Get Address Default:ShopDefault - false:true",
			userID:       "123456",
			pgn:          &pagination.Pagination{},
			queryRequest: &body.GetAddressQueryRequest{IsDefaultBool: false, IsShopDefaultBool: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultShopAddress", mock.Anything, mock.Anything).Return(&model.Address{}, nil)
			},
			expectedErr: nil,
		},
		{
			name:         "Error GetDefaultShopAddress no sql row",
			userID:       "123456",
			pgn:          &pagination.Pagination{},
			queryRequest: &body.GetAddressQueryRequest{IsDefaultBool: false, IsShopDefaultBool: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultShopAddress", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New("Shop address not found."),
		},
		{
			name:         "Error GetDefaultShopAddress",
			userID:       "123456",
			pgn:          &pagination.Pagination{},
			queryRequest: &body.GetAddressQueryRequest{IsDefaultBool: false, IsShopDefaultBool: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetDefaultShopAddress", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetAddress(context.Background(), tc.userID, tc.pgn, tc.queryRequest)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_GetOrder(t *testing.T) {
	testCase := []struct {
		name          string
		userID        string
		orderStatusID string
		pgn           *pagination.Pagination
		mock          func(t *testing.T, r *mocks.Repository)
		expectedErr   error
	}{
		{
			name:          "success Get order",
			userID:        "123456",
			orderStatusID: "1",
			pgn:           &pagination.Pagination{},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(10)
				r.On("GetTotalOrder", mock.Anything, mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("GetOrders", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*model.Order{}, nil)
			},
			expectedErr: nil,
		},
		{
			name:          "Error Total Errows",
			userID:        "123456",
			orderStatusID: "1",
			pgn:           &pagination.Pagination{},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(-1)
				r.On("GetTotalOrder", mock.Anything, mock.Anything, mock.Anything).Return(tempInt64, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:          "Error Get orders",
			userID:        "123456",
			orderStatusID: "1",
			pgn:           &pagination.Pagination{},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(10)
				r.On("GetTotalOrder", mock.Anything, mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("GetOrders", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetOrder(context.Background(), tc.userID, tc.orderStatusID, tc.pgn)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_GetOrderByOrderID(t *testing.T) {
	testCase := []struct {
		name        string
		orderID     string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:    "Success Get Order By Order ID",
			orderID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempString := "{\"rajaongkir\":{\"query\":{\"origin\":\"501\",\"destination\":\"114\",\"weight\":1700,\"courier\":\"jne\"},\"status\":{\"code\":200,\"description\":\"OK\"},\"origin_details\":{\"city_id\":\"501\",\"province_id\":\"5\",\"province\":\"DI Yogyakarta\",\"type\":\"Kota\",\"city_name\":\"Yogyakarta\",\"postal_code\":\"55000\"},\"destination_details\":{\"city_id\":\"114\",\"province_id\":\"1\",\"province\":\"Bali\",\"type\":\"Kota\",\"city_name\":\"Denpasar\",\"postal_code\":\"80000\"},\"results\":[{\"code\":\"jne\",\"name\":\"Jalur Nugraha Ekakurir (JNE)\",\"costs\":[{\"service\":\"OKE\",\"description\":\"Ongkos Kirim Ekonomis\",\"cost\":[{\"value\":38000,\"etd\":\"4-5\",\"note\":\"\"}]},{\"service\":\"REG\",\"description\":\"Layanan Reguler\",\"cost\":[{\"value\":44000,\"etd\":\"2-3\",\"note\":\"\"}]},{\"service\":\"SPS\",\"description\":\"Super Speed\",\"cost\":[{\"value\":349000,\"etd\":\"\",\"note\":\"\"}]},{\"service\":\"YES\",\"description\":\"Yakin Esok Sampai\",\"cost\":[{\"value\":98000,\"etd\":\"1-1\",\"note\":\"\"}]}]}]}}"
				tempRajaOngkir := body2.RajaOngkirCostResponse{
					Rajaongkir: struct {
						Query struct {
							Origin      string "json:\"origin,omitempty\""
							Destination string "json:\"destination,omitempty\""
							Weight      int    "json:\"weight,omitempty\""
							Courier     string "json:\"courier,omitempty\""
						} "json:\"query,omitempty\""
						Status struct {
							Code        int    "json:\"code,omitempty\""
							Description string "json:\"description,omitempty\""
						} "json:\"status,omitempty\""
						OriginDetails struct {
							CityID     string "json:\"city_id,omitempty\""
							ProvinceID string "json:\"province_id,omitempty\""
							Province   string "json:\"province,omitempty\""
							Type       string "json:\"type,omitempty\""
							CityName   string "json:\"city_name,omitempty\""
							PostalCode string "json:\"postal_code,omitempty\""
						} "json:\"origin_details,omitempty\""
						DestinationDetails struct {
							CityID     string "json:\"city_id,omitempty\""
							ProvinceID string "json:\"province_id,omitempty\""
							Province   string "json:\"province,omitempty\""
							Type       string "json:\"type,omitempty\""
							CityName   string "json:\"city_name,omitempty\""
							PostalCode string "json:\"postal_code,omitempty\""
						} "json:\"destination_details,omitempty\""
						Results []struct {
							Code  string "json:\"code,omitempty\""
							Name  string "json:\"name,omitempty\""
							Costs []struct {
								Service     string "json:\"service,omitempty\""
								Description string "json:\"description,omitempty\""
								Cost        []struct {
									Value int    "json:\"value,omitempty\""
									Etd   string "json:\"etd,omitempty\""
									Note  string "json:\"note,omitempty\""
								} "json:\"cost,omitempty\""
							} "json:\"costs,omitempty\""
						} "json:\"results,omitempty\""
					}{
						Query: struct {
							Origin      string "json:\"origin,omitempty\""
							Destination string "json:\"destination,omitempty\""
							Weight      int    "json:\"weight,omitempty\""
							Courier     string "json:\"courier,omitempty\""
						}{
							Origin:      "origin",
							Destination: "des",
							Weight:      11,
							Courier:     "jne",
						},
						Status: struct {
							Code        int    "json:\"code,omitempty\""
							Description string "json:\"description,omitempty\""
						}{
							Code:        1,
							Description: "descrpition",
						},
						OriginDetails: struct {
							CityID     string "json:\"city_id,omitempty\""
							ProvinceID string "json:\"province_id,omitempty\""
							Province   string "json:\"province,omitempty\""
							Type       string "json:\"type,omitempty\""
							CityName   string "json:\"city_name,omitempty\""
							PostalCode string "json:\"postal_code,omitempty\""
						}{
							CityID:     "1",
							ProvinceID: "1",
							Province:   "palembang",
							Type:       "type",
							CityName:   "palembang",
							PostalCode: "12212",
						},
						DestinationDetails: struct {
							CityID     string "json:\"city_id,omitempty\""
							ProvinceID string "json:\"province_id,omitempty\""
							Province   string "json:\"province,omitempty\""
							Type       string "json:\"type,omitempty\""
							CityName   string "json:\"city_name,omitempty\""
							PostalCode string "json:\"postal_code,omitempty\""
						}{
							CityID:     "test",
							ProvinceID: "test",
							Province:   "test",
							Type:       "test",
							CityName:   "test",
							PostalCode: "test",
						},
						Results: make([]struct {
							Code  string "json:\"code,omitempty\""
							Name  string "json:\"name,omitempty\""
							Costs []struct {
								Service     string "json:\"service,omitempty\""
								Description string "json:\"description,omitempty\""
								Cost        []struct {
									Value int    "json:\"value,omitempty\""
									Etd   string "json:\"etd,omitempty\""
									Note  string "json:\"note,omitempty\""
								} "json:\"cost,omitempty\""
							} "json:\"costs,omitempty\""
						}, 1),
					}}
				for _, res := range tempRajaOngkir.Rajaongkir.Results {
					fmt.Println("result index")
					res.Code = "code"
					res.Name = "name"
					res.Costs = make([]struct {
						Service     string "json:\"service,omitempty\""
						Description string "json:\"description,omitempty\""
						Cost        []struct {
							Value int    "json:\"value,omitempty\""
							Etd   string "json:\"etd,omitempty\""
							Note  string "json:\"note,omitempty\""
						} "json:\"cost,omitempty\""
					}, 1)
					fmt.Println("tempRajaOngkir.Rajaongkir.Results", res)
					for _, c := range res.Costs {
						fmt.Println("cost masuk")
						c.Service = "service"
						c.Description = "description"
						c.Cost = make([]struct {
							Value int    "json:\"value,omitempty\""
							Etd   string "json:\"etd,omitempty\""
							Note  string "json:\"note,omitempty\""
						}, 1)
						fmt.Println("res.Costs", c)
						for _, cc := range c.Cost {
							fmt.Println("ccccc")
							cc.Value = 111
							cc.Etd = "etd"
							cc.Note = "note"
							fmt.Println("cc", cc)
						}
						fmt.Println("after res.Costs", c)
					}
				}
				fmt.Println("test - rajaongkir", tempRajaOngkir)
				r.On("GetOrderByOrderID", mock.Anything, mock.Anything).Return(&model.Order{
					SellerAddress: &model.Address{CityID: 501},
					BuyerAddress:  &model.Address{CityID: 114},
					Detail: []*model.OrderDetail{
						{
							ProductWeight: 100,
							OrderQuantity: 1,
						},
					},
					CourierCode:    "jne",
					CourierService: "OKE",
				}, nil)
				r.On("GetBuyerIDByOrderID", mock.Anything, mock.Anything).Return("buyer", nil)
				r.On("GetSellerIDByOrderID", mock.Anything, mock.Anything).Return("seller", nil)
				r.On("GetCostRedis", mock.Anything, mock.Anything).Return(&tempString, nil)
				// r.On("GetCostRajaOngkir", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&tempRajaOngkir, nil)
				// r.On("InsertCostRedis", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetOrderByOrderID(context.Background(), tc.orderID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_ChangeOrderStatus(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		requestBody body.ChangeOrderStatusRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success Change Order Status",
			userID: "123456",
			requestBody: body.ChangeOrderStatusRequest{
				OrderID:       "123",
				OrderStatusID: 7,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempBuyerID := "123456"
				tempProductSold := make([]*body.ProductUnitSoldOrderQty, 0)
				tempBody := &body.ProductUnitSoldOrderQty{
					Quantity: 1, UnitSold: 1,
				}
				tempProductSold = append(tempProductSold, tempBody)
				r.On("GetBuyerIDByOrderID", mock.Anything, mock.Anything).Return(tempBuyerID, nil)
				r.On("ChangeOrderStatus", mock.Anything, mock.Anything).Return(nil)
				r.On("GetProductUnitSoldByOrderID", mock.Anything, mock.Anything, mock.Anything).Return(tempProductSold, nil)
				r.On("UpdateProductUnitSold", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:        "Error userID != order buyer ID",
			userID:      "123456",
			requestBody: body.ChangeOrderStatusRequest{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetBuyerIDByOrderID", mock.Anything, mock.Anything).Return("", errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:        "Error userID != order buyer ID",
			userID:      "123456",
			requestBody: body.ChangeOrderStatusRequest{},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempBuyerID := "111111"
				r.On("GetBuyerIDByOrderID", mock.Anything, mock.Anything).Return(tempBuyerID, nil)
			},
			expectedErr: errors.New("Invalid Credentials."),
		},
		{
			name:   "Error Order Status ID is not 7 and 6",
			userID: "123456",
			requestBody: body.ChangeOrderStatusRequest{
				OrderID:       "123",
				OrderStatusID: 1,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempBuyerID := "123456"
				r.On("GetBuyerIDByOrderID", mock.Anything, mock.Anything).Return(tempBuyerID, nil)
			},
			expectedErr: errors.New("Invalid request."),
		},
		{
			name:   "Error ChangeOrderStatus",
			userID: "123456",
			requestBody: body.ChangeOrderStatusRequest{
				OrderID:       "123",
				OrderStatusID: 7,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempBuyerID := "123456"
				r.On("GetBuyerIDByOrderID", mock.Anything, mock.Anything).Return(tempBuyerID, nil)
				r.On("ChangeOrderStatus", mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "Error GetProductUnitSoldByOrderID",
			userID: "123456",
			requestBody: body.ChangeOrderStatusRequest{
				OrderID:       "123",
				OrderStatusID: 7,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempBuyerID := "123456"
				r.On("GetBuyerIDByOrderID", mock.Anything, mock.Anything).Return(tempBuyerID, nil)
				r.On("ChangeOrderStatus", mock.Anything, mock.Anything).Return(nil)
				r.On("GetProductUnitSoldByOrderID", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "Error UpdateProductUnitSold",
			userID: "123456",
			requestBody: body.ChangeOrderStatusRequest{
				OrderID:       "123",
				OrderStatusID: 7,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempBuyerID := "123456"
				tempProductSold := make([]*body.ProductUnitSoldOrderQty, 0)
				tempBody := &body.ProductUnitSoldOrderQty{
					Quantity: 1, UnitSold: 1,
				}
				tempProductSold = append(tempProductSold, tempBody)
				r.On("GetBuyerIDByOrderID", mock.Anything, mock.Anything).Return(tempBuyerID, nil)
				r.On("ChangeOrderStatus", mock.Anything, mock.Anything).Return(nil)
				r.On("GetProductUnitSoldByOrderID", mock.Anything, mock.Anything, mock.Anything).Return(tempProductSold, nil)
				r.On("UpdateProductUnitSold", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.ChangeOrderStatus(context.Background(), tc.userID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_GetTransactionDetailByID(t *testing.T) {
	testCase := []struct {
		name          string
		transactionID string
		userID        string
		mock          func(t *testing.T, r *mocks.Repository)
		expectedErr   error
	}{
		{
			name:          "success Get Transaction Detail By ID",
			transactionID: "123456",
			userID:        "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&model.Transaction{}, nil)
				r.On("GetOrdersByTransactionID", mock.Anything, mock.Anything, mock.Anything).Return([]*model.Order{}, nil)
			},
			expectedErr: nil,
		},
		{
			name:          "Error repo GetTransactionByID no sql rows",
			transactionID: "123456",
			userID:        "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New("Transaction not found."),
		},
		{
			name:          "Error repo GetTransactionByID",
			transactionID: "123456",
			userID:        "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:          "Error repo GetOrdersByTransactionID",
			transactionID: "123456",
			userID:        "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&model.Transaction{}, nil)
				r.On("GetOrdersByTransactionID", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetTransactionDetailByID(context.Background(), tc.transactionID, tc.userID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_GetAddressByID(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		addressID   string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:      "success Get Address By ID",
			userID:    "123456",
			addressID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{}, nil)
			},
			expectedErr: nil,
		},
		{
			name:      "error repo GetUserByID no sql rows",
			userID:    "123456",
			addressID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, sql.ErrNoRows)
			},
			expectedErr: errors.New("Invalid Credentials."),
		},
		{
			name:      "error repo GetUserByID",
			userID:    "123456",
			addressID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:      "error repo GetAddressByID no sql rows",
			userID:    "123456",
			addressID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{}, sql.ErrNoRows)
			},
			expectedErr: errors.New("Address not exist."),
		},
		{
			name:      "error repo GetAddressByID",
			userID:    "123456",
			addressID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{}, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetAddressByID(context.Background(), tc.userID, tc.addressID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_DeleteAddressByID(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		addressID   string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:      "success Delete Address By ID",
			userID:    "123456",
			addressID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("DeleteAddress", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:      "error repo GetUserByID no sql rows",
			userID:    "123456",
			addressID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, sql.ErrNoRows)
			},
			expectedErr: errors.New("Invalid Credentials."),
		},
		{
			name:      "error repo GetUserByID",
			userID:    "123456",
			addressID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:      "error repo GetAddressByID no sql rows",
			userID:    "123456",
			addressID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{}, sql.ErrNoRows)
			},
			expectedErr: errors.New("Address not exist."),
		},
		{
			name:      "error repo GetAddressByID",
			userID:    "123456",
			addressID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{}, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:      "error delete address",
			userID:    "123456",
			addressID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{IsDefault: true}, nil)
			},
			expectedErr: errors.New("Address is default."),
		},
		{
			name:      "error delete address",
			userID:    "123456",
			addressID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("DeleteAddress", mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			err := u.DeleteAddressByID(context.Background(), tc.userID, tc.addressID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_CompletedRejectedRefund(t *testing.T) {
	testCase := []struct {
		name        string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success Completed Rejected Refund",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempBuyerID := uuid.Nil.String()
				tempProductSold := make([]*body.ProductUnitSoldOrderQty, 0)
				tempBody := &body.ProductUnitSoldOrderQty{
					Quantity: 1, UnitSold: 1,
				}
				tempProductSold = append(tempProductSold, tempBody)

				tempRefunds := make([]*model.RefundOrder, 0)
				refund := &model.RefundOrder{Order: &model.OrderModel{UserID: uuid.Nil, ID: uuid.Nil}}
				tempRefunds = append(tempRefunds, refund)
				r.On("GetRejectedRefund", mock.Anything).Return(tempRefunds, nil)
				r.On("GetBuyerIDByOrderID", mock.Anything, mock.Anything).Return(tempBuyerID, nil)
				r.On("ChangeOrderStatus", mock.Anything, mock.Anything).Return(nil)
				r.On("GetProductUnitSoldByOrderID", mock.Anything, mock.Anything, mock.Anything).Return(tempProductSold, nil)
				r.On("UpdateProductUnitSold", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Error Change Order Status ",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempRefunds := make([]*model.RefundOrder, 0)
				refund := &model.RefundOrder{Order: &model.OrderModel{UserID: uuid.Nil, ID: uuid.Nil}}
				tempRefunds = append(tempRefunds, refund)
				r.On("GetRejectedRefund", mock.Anything).Return(tempRefunds, nil)
				r.On("GetBuyerIDByOrderID", mock.Anything, mock.Anything).Return("", errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name: "Error Repo GetRejectedRefund",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRejectedRefund", mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()

			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.CompletedRejectedRefund(context.Background())
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_EditUser(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		requestBody body.EditUserRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success Edit User",
			userID: "123456",
			requestBody: body.EditUserRequest{
				Username:  "username",
				PhoneNo:   "19283",
				BirthDate: "02-01-2006",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{IsVerify: true}, nil)
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				r.On("GetUserByPhoneNo", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				r.On("UpdateUserField", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "Error repo UpdateUserField",
			userID: "123456",
			requestBody: body.EditUserRequest{
				Username:  "username",
				PhoneNo:   "19283",
				BirthDate: "02-01-2006",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{IsVerify: true}, nil)
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				r.On("GetUserByPhoneNo", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				r.On("UpdateUserField", mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "Error user is verify false",
			userID: "123456",
			requestBody: body.EditUserRequest{
				Username:  "username",
				PhoneNo:   "19283",
				BirthDate: "02-01-2006",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{IsVerify: false}, nil)
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				r.On("GetUserByPhoneNo", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(response.UserNotVerifyMessage),
		},
		{
			name:   "Error user Phone No is different",
			userID: "123456",
			requestBody: body.EditUserRequest{
				Username:  "username",
				PhoneNo:   "19283",
				BirthDate: "02-01-2006",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempPhoneUser := "1111"
				tempPhone := "123456"
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{IsVerify: true, PhoneNo: &tempPhoneUser}, nil)
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				r.On("GetUserByPhoneNo", mock.Anything, mock.Anything).Return(&model.User{PhoneNo: &tempPhone}, nil)
			},
			expectedErr: errors.New(response.PhoneNoAlreadyExistMessage),
		},
		{
			name:   "Error repo GetUserByPhoneNo",
			userID: "123456",
			requestBody: body.EditUserRequest{
				Username:  "username",
				PhoneNo:   "19283",
				BirthDate: "02-01-2006",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{IsVerify: true}, nil)
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				r.On("GetUserByPhoneNo", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "Error Username is different",
			userID: "123456",
			requestBody: body.EditUserRequest{
				Username:  "username",
				PhoneNo:   "19283",
				BirthDate: "02-01-2006",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempUsername := "username"
				tempUsername1 := "username1"
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{IsVerify: true, Username: &tempUsername}, nil)
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(&model.User{Username: &tempUsername1}, nil)
			},
			expectedErr: errors.New(response.UserNameAlreadyExistMessage),
		},
		{
			name:   "Error Username is different",
			userID: "123456",
			requestBody: body.EditUserRequest{
				Username:  "username",
				PhoneNo:   "19283",
				BirthDate: "02-01-2006",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{IsVerify: true}, nil)
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "Error Username is different",
			userID: "123456",
			requestBody: body.EditUserRequest{
				Username:  "username",
				PhoneNo:   "19283",
				BirthDate: "02-01-2006",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(response.UnauthorizedMessage),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.EditUser(context.Background(), tc.userID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_EditEmail(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		requestBody body.EditEmailRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success Edit User",
			userID: "123456",
			requestBody: body.EditEmailRequest{
				Email: "email@gmail.com",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				r.On("InsertNewOTPKey", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "Error repo InsertNewOTPKey",
			userID: "123456",
			requestBody: body.EditEmailRequest{
				Email: "email@gmail.com",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				r.On("InsertNewOTPKey", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "Error Email History found same as UserModel",
			userID: "123456",
			requestBody: body.EditEmailRequest{
				Email: "email@gmail.com",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email@gmail.com"}, nil)
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(&model.EmailHistory{Email: "email@gmail.com"}, nil)
			},
			expectedErr: errors.New(response.EmailSamePreviousEmailMessage),
		},
		{
			name:   "Error Email History found",
			userID: "123456",
			requestBody: body.EditEmailRequest{
				Email: "email@gmail.com",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email1@gmail.com"}, nil)
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(&model.EmailHistory{Email: "email@gmail.com"}, nil)
			},
			expectedErr: errors.New(response.EmailAlreadyExistMessage),
		},
		{
			name:   "Error Repo CheckEmailHistory",
			userID: "123456",
			requestBody: body.EditEmailRequest{
				Email: "email@gmail.com",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email1@gmail.com"}, nil)
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "Error repo GetUserByID sql no rows",
			userID: "123456",
			requestBody: body.EditEmailRequest{
				Email: "email@gmail.com",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(response.UnauthorizedMessage),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.EditEmail(context.Background(), tc.userID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_EditEmailUser(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		requestBody body.EditEmailUserRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success Edit User",
			userID: "123456",
			requestBody: body.EditEmailUserRequest{
				Email: "email@gmail.com",
				Code:  "5694d08a2e53ffcae0c3103e5ad6f6076abd960eb1f8a56577040bc1028f702b",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempnt64 := int64(1)
				r.On("GetOTPValue", mock.Anything, mock.Anything).Return("code", nil)
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("UpdateUserEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("CreateEmailHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("DeleteOTPValue", mock.Anything, mock.Anything).Return(tempnt64, nil)
			},
			expectedErr: nil,
		},
		{
			name:   "Error repo DeleteOTPValue",
			userID: "123456",
			requestBody: body.EditEmailUserRequest{
				Email: "email@gmail.com",
				Code:  "5694d08a2e53ffcae0c3103e5ad6f6076abd960eb1f8a56577040bc1028f702b",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempnt64 := int64(1)
				r.On("GetOTPValue", mock.Anything, mock.Anything).Return("code", nil)
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("UpdateUserEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("CreateEmailHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("DeleteOTPValue", mock.Anything, mock.Anything).Return(tempnt64, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "Error repo DeleteOTPValue",
			userID: "123456",
			requestBody: body.EditEmailUserRequest{
				Email: "email@gmail.com",
				Code:  "5694d08a2e53ffcae0c3103e5ad6f6076abd960eb1f8a56577040bc1028f702b",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetOTPValue", mock.Anything, mock.Anything).Return("code", nil)
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("UpdateUserEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("CreateEmailHistory", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "Error repo UpdateUserEmail",
			userID: "123456",
			requestBody: body.EditEmailUserRequest{
				Email: "email@gmail.com",
				Code:  "5694d08a2e53ffcae0c3103e5ad6f6076abd960eb1f8a56577040bc1028f702b",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetOTPValue", mock.Anything, mock.Anything).Return("code", nil)
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("UpdateUserEmail", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "Error repo GetUserByID no sql rows",
			userID: "123456",
			requestBody: body.EditEmailUserRequest{
				Email: "email@gmail.com",
				Code:  "5694d08a2e53ffcae0c3103e5ad6f6076abd960eb1f8a56577040bc1028f702b",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetOTPValue", mock.Anything, mock.Anything).Return("code", nil)
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(response.UnauthorizedMessage),
		},
		{
			name:   "Error OTP Code Hashed",
			userID: "123456",
			requestBody: body.EditEmailUserRequest{
				Email: "email@gmail.com",
				Code:  "5694d08a2e53ffcae0c3103e5ad6f6076abd960eb1f8a56577040bc1028f702b",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetOTPValue", mock.Anything, mock.Anything).Return("code1", nil)
			},
			expectedErr: errors.New(response.OTPIsNotValidMessage),
		},
		{
			name:   "Error repo GetOTPValue",
			userID: "123456",
			requestBody: body.EditEmailUserRequest{
				Email: "email@gmail.com",
				Code:  "5694d08a2e53ffcae0c3103e5ad6f6076abd960eb1f8a56577040bc1028f702b",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetOTPValue", mock.Anything, mock.Anything).Return("", errors.New("test"))
			},
			expectedErr: errors.New(response.OTPAlreadyExpiredMessage),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()

			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.EditEmailUser(context.Background(), tc.userID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_GetSealabsPay(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success Edit User",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSealabsPay", mock.Anything, mock.Anything).Return([]*model.SealabsPay{}, nil)
			},
			expectedErr: nil,
		},
		{
			name:   "error repo GetSealabsPay",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSealabsPay", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetSealabsPay(context.Background(), tc.userID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_AddSealabsPay(t *testing.T) {
	testCase := []struct {
		name        string
		request     body.AddSealabsPayRequest
		userid      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:    "success Edit User",
			request: body.AddSealabsPayRequest{},
			userid:  "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckUserSealabsPay", mock.Anything, mock.Anything).Return(0, nil)
				r.On("CheckDeletedSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(0, nil)
				r.On("AddSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "Error slp found",
			request: body.AddSealabsPayRequest{CardNumber: "123456"},
			userid:  "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempCardNumber := "222222"
				r.On("CheckUserSealabsPay", mock.Anything, mock.Anything).Return(1, nil)
				r.On("CheckDeletedSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(0, nil)
				r.On("CheckDefaultSealabsPay", mock.Anything, mock.Anything).Return(&tempCardNumber, nil)
				r.On("SetDefaultSealabsPayTrans", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("AddSealabsPayTrans", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "Error slp found and failed to save",
			request: body.AddSealabsPayRequest{CardNumber: "123456"},
			userid:  "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempCardNumber := "222222"
				r.On("CheckUserSealabsPay", mock.Anything, mock.Anything).Return(1, nil)
				r.On("CheckDeletedSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(0, nil)
				r.On("CheckDefaultSealabsPay", mock.Anything, mock.Anything).Return(&tempCardNumber, nil)
				r.On("SetDefaultSealabsPayTrans", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("AddSealabsPayTrans", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New(response.SealabsCardAlreadyExist),
		},
		{
			name:    "Error slp found and failed to update",
			request: body.AddSealabsPayRequest{CardNumber: "123456"},
			userid:  "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempCardNumber := "222222"
				r.On("CheckUserSealabsPay", mock.Anything, mock.Anything).Return(1, nil)
				r.On("CheckDeletedSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
				r.On("CheckDefaultSealabsPay", mock.Anything, mock.Anything).Return(&tempCardNumber, nil)
				r.On("SetDefaultSealabsPayTrans", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateUserSealabsPayTrans", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New(response.SealabsCardAlreadyExist),
		},
		{
			name:    "Error repo SetDefaultSealabsPayTrans",
			request: body.AddSealabsPayRequest{CardNumber: "123456"},
			userid:  "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempCardNumber := "222222"
				r.On("CheckUserSealabsPay", mock.Anything, mock.Anything).Return(1, nil)
				r.On("CheckDeletedSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
				r.On("CheckDefaultSealabsPay", mock.Anything, mock.Anything).Return(&tempCardNumber, nil)
				r.On("SetDefaultSealabsPayTrans", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New(response.SealabsCardAlreadyExist),
		},
		{
			name:    "Error Card Number ",
			request: body.AddSealabsPayRequest{CardNumber: "123456"},
			userid:  "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempCardNumber := "123456"
				r.On("CheckUserSealabsPay", mock.Anything, mock.Anything).Return(1, nil)
				r.On("CheckDeletedSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
				r.On("CheckDefaultSealabsPay", mock.Anything, mock.Anything).Return(&tempCardNumber, nil)
			},
			expectedErr: errors.New(response.SealabsCardAlreadyExist),
		},
		{
			name:    "Error repo CheckDefaultSealabsPay",
			request: body.AddSealabsPayRequest{CardNumber: "123456"},
			userid:  "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckUserSealabsPay", mock.Anything, mock.Anything).Return(1, nil)
				r.On("CheckDeletedSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
				r.On("CheckDefaultSealabsPay", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:    "Error repo AddSealabsPay ",
			request: body.AddSealabsPayRequest{CardNumber: "123456"},
			userid:  "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckUserSealabsPay", mock.Anything, mock.Anything).Return(0, nil)
				r.On("CheckDeletedSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(0, nil)
				r.On("AddSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New(response.SealabsCardAlreadyExist),
		},
		{
			name:    "Error repo UpdateUserSealabsPay ",
			request: body.AddSealabsPayRequest{CardNumber: "123456"},
			userid:  "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckUserSealabsPay", mock.Anything, mock.Anything).Return(0, nil)
				r.On("CheckDeletedSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
				r.On("UpdateUserSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New(response.SealabsCardAlreadyExist),
		},
		{
			name:    "Error repo CheckDeletedSealabsPay ",
			request: body.AddSealabsPayRequest{CardNumber: "123456"},
			userid:  "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckUserSealabsPay", mock.Anything, mock.Anything).Return(0, nil)
				r.On("CheckDeletedSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(-1, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},

		{
			name:    "Error repo CheckUserSealabsPay ",
			request: body.AddSealabsPayRequest{CardNumber: "123456"},
			userid:  "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckUserSealabsPay", mock.Anything, mock.Anything).Return(-1, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()

			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.AddSealabsPay(context.Background(), tc.request, tc.name)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_PatchSealabsPay(t *testing.T) {
	testCase := []struct {
		name        string
		cardNumber  string
		userid      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:       "success Patch Sealabs Pay",
			cardNumber: "123456",
			userid:     "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("PatchSealabsPay", mock.Anything, mock.Anything).Return(nil)
				r.On("SetDefaultSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:       "error repo SetDefaultSealabsPay",
			cardNumber: "123456",
			userid:     "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("PatchSealabsPay", mock.Anything, mock.Anything).Return(nil)
				r.On("SetDefaultSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:       "error repo PatchSealabsPay",
			cardNumber: "123456",
			userid:     "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("PatchSealabsPay", mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			err := u.PatchSealabsPay(context.Background(), tc.cardNumber, tc.userid)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_DeleteSealabsPay(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		cardNumber  string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:       "success Delete Sealabs Pay",
			userID:     "123456",
			cardNumber: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSealabsPayUser", mock.Anything, mock.Anything, mock.Anything).Return(&model.SealabsPay{}, nil)
				r.On("DeleteSealabsPay", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:       "error repo DeleteSealabsPay",
			userID:     "123456",
			cardNumber: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSealabsPayUser", mock.Anything, mock.Anything, mock.Anything).Return(&model.SealabsPay{}, nil)
				r.On("DeleteSealabsPay", mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:       "error slp is defult true",
			userID:     "123456",
			cardNumber: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSealabsPayUser", mock.Anything, mock.Anything, mock.Anything).Return(&model.SealabsPay{IsDefault: true}, nil)
			},
			expectedErr: errors.New(response.SealabsCardIsDefault),
		},
		{
			name:       "error repo GetSealabsPayUser",
			userID:     "123456",
			cardNumber: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSealabsPayUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:       "error repo GetSealabsPayUser no sql row",
			userID:     "123456",
			cardNumber: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSealabsPayUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(response.SealabsCardNotFound),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			err := u.DeleteSealabsPay(context.Background(), tc.cardNumber, tc.userID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_ActivateWallet(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		pin         string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success Activate Wallet",
			userID: "123456",
			pin:    "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				r.On("CreateWallet", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "error repo CreateWallet",
			userID: "123456",
			pin:    "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				r.On("CreateWallet", mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "error wallet found",
			userID: "123456",
			pin:    "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(&model.Wallet{}, nil)
			},
			expectedErr: errors.New(response.WalletAlreadyActivated),
		},
		{
			name:   "error repo GetWalletByUserID",
			userID: "123456",
			pin:    "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "error repo GetUserByID",
			userID: "123456",
			pin:    "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			err := u.ActivateWallet(context.Background(), tc.userID, tc.pin)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_RegisterMerchant(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		shopName    string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:     "success Register Merchant",
			userID:   "123456",
			shopName: "shop name",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(0)
				r.On("CheckShopByID", mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("CheckShopUnique", mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, nil)
				r.On("AddShop", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateRole", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:     "error repo UpdateRole",
			userID:   "123456",
			shopName: "shop name",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(0)
				r.On("CheckShopByID", mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("CheckShopUnique", mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, nil)
				r.On("AddShop", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateRole", mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:     "error repo AddShop",
			userID:   "123456",
			shopName: "shop name",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(0)
				r.On("CheckShopByID", mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("CheckShopUnique", mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, nil)
				r.On("AddShop", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:     "error repo GetWalletByUserID",
			userID:   "123456",
			shopName: "shop name",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(0)
				r.On("CheckShopByID", mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("CheckShopUnique", mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:     "error repo GetWalletByUserID",
			userID:   "123456",
			shopName: "shop name",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(0)
				r.On("CheckShopByID", mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("CheckShopUnique", mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(response.WalletIsNotActivated),
		},
		{
			name:     "error shop unique is found",
			userID:   "123456",
			shopName: "shop name",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(0)
				r.On("CheckShopByID", mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("CheckShopUnique", mock.Anything, mock.Anything).Return(int64(1), nil)
			},
			expectedErr: errors.New(response.ShopAlreadyExists),
		},
		{
			name:     "error repo CheckShopUnique",
			userID:   "123456",
			shopName: "shop name",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(0)
				r.On("CheckShopByID", mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("CheckShopUnique", mock.Anything, mock.Anything).Return(tempInt64, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:     "error shop CheckShopByID is found",
			userID:   "123456",
			shopName: "shop name",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckShopByID", mock.Anything, mock.Anything).Return(int64(1), nil)
			},
			expectedErr: errors.New(response.UserAlreadyHaveShop),
		},
		{
			name:     "error repo CheckShopByID",
			userID:   "123456",
			shopName: "shop name",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(0)
				r.On("CheckShopByID", mock.Anything, mock.Anything).Return(tempInt64, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			err := u.RegisterMerchant(context.Background(), tc.userID, tc.shopName)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_GetUserProfile(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success Get User Profile",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
			},
			expectedErr: nil,
		},
		{
			name:   "error repo GetUserProfile",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "error repo GetUserProfile no sql row",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(response.UserNotExistMessage),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetUserProfile(context.Background(), tc.userID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_UploadProfilePicture(t *testing.T) {
	testCase := []struct {
		name        string
		imgURL      string
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success Upload Profile Picture",
			imgURL: "1233456",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("UpdateProfileImage", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "error repo UpdateProfileImage",
			imgURL: "1233456",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("UpdateProfileImage", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			err := u.UploadProfilePicture(context.Background(), tc.imgURL, tc.name)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_VerifyPasswordChange(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success Verify Password Change",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email@gmail.com"}, nil)
				r.On("InsertNewOTPKey", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "error repo InsertNewOTPKey",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email@gmail.com"}, nil)
				r.On("InsertNewOTPKey", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "error repo GetUserByID",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			err := u.VerifyPasswordChange(context.Background(), tc.userID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_VerifyOTP(t *testing.T) {
	testCase := []struct {
		name        string
		requestBody body.VerifyOTPRequest
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success VerifyOTP",
			requestBody: body.VerifyOTPRequest{
				OTP: "123456",
			},
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email@gmail.com"}, nil)
				r.On("GetOTPValue", mock.Anything, mock.Anything).Return("123456", nil)
				r.On("DeleteOTPValue", mock.Anything, mock.Anything).Return(int64(1), nil)
			},
			expectedErr: nil,
		},
		{
			name: "error repo DeleteOTPValue",
			requestBody: body.VerifyOTPRequest{
				OTP: "123456",
			},
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email@gmail.com"}, nil)
				r.On("GetOTPValue", mock.Anything, mock.Anything).Return("123456", nil)
				r.On("DeleteOTPValue", mock.Anything, mock.Anything).Return(int64(0), errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name: "error repo DeleteOTPValue",
			requestBody: body.VerifyOTPRequest{
				OTP: "123456",
			},
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email@gmail.com"}, nil)
				r.On("GetOTPValue", mock.Anything, mock.Anything).Return("123451", nil)
			},
			expectedErr: errors.New(response.OTPIsNotValidMessage),
		},
		{
			name: "error repo DeleteOTPValue",
			requestBody: body.VerifyOTPRequest{
				OTP: "123456",
			},
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email@gmail.com"}, nil)
				r.On("GetOTPValue", mock.Anything, mock.Anything).Return("", errors.New("test"))
			},
			expectedErr: errors.New(response.OTPAlreadyExpiredMessage),
		},
		{
			name: "error repo DeleteOTPValue",
			requestBody: body.VerifyOTPRequest{
				OTP: "123456",
			},
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.VerifyOTP(context.Background(), tc.requestBody, tc.userID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_ChangePassword(t *testing.T) {
	passwordHash := "$2a$10$WKul/6gjYoYjOXuNVX4XGen1ZkWYb1PKFiI5vlZp5TFerZh6nTujG"
	testCase := []struct {
		name        string
		userID      string
		newPassword string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:        "success ChangePassword",
			userID:      "123456",
			newPassword: "Tested7*",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempUsername := "juww"
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email@gmail.com", Username: &tempUsername}, nil)
				r.On("GetPasswordByID", mock.Anything, mock.Anything).Return(passwordHash, nil)
				r.On("UpdatePasswordByID", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetSessionKeyRedis", mock.Anything, mock.Anything, mock.Anything).Return([]string{"asd"}, nil)
				r.On("InsertSessionRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:        "error repo InsertSessionRedis",
			userID:      "123456",
			newPassword: "Tested7*",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempUsername := "juww"
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email@gmail.com", Username: &tempUsername}, nil)
				r.On("GetPasswordByID", mock.Anything, mock.Anything).Return(passwordHash, nil)
				r.On("UpdatePasswordByID", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetSessionKeyRedis", mock.Anything, mock.Anything, mock.Anything).Return([]string{"asd"}, nil)
				r.On("InsertSessionRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:        "error repo GetSessionKeyRedis",
			userID:      "123456",
			newPassword: "Tested7*",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempUsername := "juww"
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email@gmail.com", Username: &tempUsername}, nil)
				r.On("GetPasswordByID", mock.Anything, mock.Anything).Return(passwordHash, nil)
				r.On("UpdatePasswordByID", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetSessionKeyRedis", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:        "error repo UpdatePasswordByID",
			userID:      "123456",
			newPassword: "Tested7*",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempUsername := "juww"
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email@gmail.com", Username: &tempUsername}, nil)
				r.On("GetPasswordByID", mock.Anything, mock.Anything).Return(passwordHash, nil)
				r.On("UpdatePasswordByID", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:        "error password contain",
			userID:      "123456",
			newPassword: "Tested7juww",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempUsername := "juww"
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email@gmail.com", Username: &tempUsername}, nil)
				r.On("GetPasswordByID", mock.Anything, mock.Anything).Return(passwordHash, nil)
			},
			expectedErr: errors.New(response.PasswordContainUsernameMessage),
		},
		{
			name:        "error password same old password",
			userID:      "123456",
			newPassword: "Tested8*",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempUsername := "juww"
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email@gmail.com", Username: &tempUsername}, nil)
				r.On("GetPasswordByID", mock.Anything, mock.Anything).Return(passwordHash, nil)
			},
			expectedErr: errors.New(response.PasswordSameOldPasswordMessage),
		},
		{
			name:        "error repo GetPasswordByID",
			userID:      "123456",
			newPassword: "Tested8*",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempUsername := "juww"
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "email@gmail.com", Username: &tempUsername}, nil)
				r.On("GetPasswordByID", mock.Anything, mock.Anything).Return("", errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:        "error repo GetUserByID",
			userID:      "123456",
			newPassword: "Tested8*",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			err := u.ChangePassword(context.Background(), tc.userID, tc.newPassword)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_TopUpWallet(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		requestBody body.TopUpWalletRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success TopUpWallet",
			userID: "123456",
			requestBody: body.TopUpWalletRequest{
				Amount: 123456,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(&model.Wallet{ID: uuid.Nil}, nil)
				r.On("GetSealabsPayUser", mock.Anything, mock.Anything, mock.Anything).Return(&model.SealabsPay{CardNumber: "123456"}, nil)
				r.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(&uuid.Nil, nil)
			},
			expectedErr: nil,
		},
		// {
		// 	name:   "error CreateTransaction",
		// 	userID: "123456",
		// 	requestBody: body.TopUpWalletRequest{
		// 		Amount: 123456,
		// 	},
		// 	mock: func(t *testing.T, r *mocks.Repository) {
		// 		r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(&model.Wallet{ID: uuid.Nil}, nil)
		// 		r.On("GetSealabsPayUser", mock.Anything, mock.Anything, mock.Anything).Return(&model.SealabsPay{CardNumber: "123456"}, nil)
		// 		r.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New(response.ProductQuantityNotAvailable))
		// 	},
		// 	expectedErr: errors.New("test"),
		// },

		{
			name:   "error GetSealabsPayUser",
			userID: "123456",
			requestBody: body.TopUpWalletRequest{
				Amount: 123456,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(&model.Wallet{ID: uuid.Nil}, nil)
				r.On("GetSealabsPayUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "error GetSealabsPayUser no sql row",
			userID: "123456",
			requestBody: body.TopUpWalletRequest{
				Amount: 123456,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(&model.Wallet{ID: uuid.Nil}, nil)
				r.On("GetSealabsPayUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(response.SealabsCardNotFound),
		},
		{
			name:   "error GetWalletByUserID",
			userID: "123456",
			requestBody: body.TopUpWalletRequest{
				Amount: 123456,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "error GetWalletByUserID sql no row",
			userID: "123456",
			requestBody: body.TopUpWalletRequest{
				Amount: 123456,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(response.WalletIsNotActivated),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()

			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.TopUpWallet(context.Background(), tc.userID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_CreateSLPPayment(t *testing.T) {
	testCase := []struct {
		name          string
		transactionID string
		mock          func(t *testing.T, r *mocks.Repository)
		expectedErr   error
	}{
		// {
		// 	name:          "success CreateSLPPayment",
		// 	transactionID: "123456",
		// 	mock: func(t *testing.T, r *mocks.Repository) {
		// 		tempCardNumber := "123456"
		// 		r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&model.Transaction{
		// 			ExpiredAt: sql.NullTime{
		// 				Valid: true,
		//				Time:  time.Now().Add(time.Hour),
		// 			},
		// 			CardNumber: &tempCardNumber,
		// 		}, nil)
		// 		r.On("GetSealabsPayUser", mock.Anything, mock.Anything, mock.Anything).Return(&model.SealabsPay{CardNumber: "123456"}, nil)
		// 		r.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(&uuid.Nil, nil)
		// 	},
		// 	expectedErr: nil,
		// },
		{
			name:          "error Card number nil",
			transactionID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&model.Transaction{
					ExpiredAt: sql.NullTime{
						Valid: true,
						Time:  time.Now().Add(time.Hour),
					},
					CardNumber: nil,
				}, nil)
			},
			expectedErr: errors.New(response.InvalidPaymentMethod),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()

			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.CreateSLPPayment(context.Background(), tc.transactionID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_CreateWalletPayment(t *testing.T) {
	testCase := []struct {
		name          string
		transactionID string
		mock          func(t *testing.T, r *mocks.Repository)
		expectedErr   error
	}{
		{
			name:          "success CreateSLPPayment",
			transactionID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempCardNumber := "123456"
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&model.Transaction{
					ID: uuid.Nil,
					ExpiredAt: sql.NullTime{
						Valid: true,
						Time:  time.Now().Add(time.Hour),
					},
					CardNumber: &tempCardNumber,
					TotalPrice: 100,
					WalletID:   &uuid.Nil,
				}, nil)
				r.On("GetWalletUser", mock.Anything, mock.Anything).Return(&model.Wallet{Balance: 1000}, nil)
				r.On("GetOrderByTransactionID", mock.Anything, mock.Anything).Return([]*model.OrderModel{{OrderStatusID: 1}}, nil)
				r.On("UpdateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateOrder", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("InsertWalletHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateWalletBalance", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(&model.Wallet{
					ID: uuid.Nil,
				}, nil)
				r.On("InsertWalletHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateWalletBalance", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			},
			expectedErr: nil,
		},
		{
			name:          "error GetTransactionByID",
			transactionID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:          "error GetTransactionByID sql no rows",
			transactionID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(response.TransactionIDNotExist),
		},
		{
			name:          "error transaction expired",
			transactionID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&model.Transaction{
					ExpiredAt: sql.NullTime{
						Valid: true,
						Time:  time.Now(),
					},
					CardNumber: nil,
				}, nil)
			},
			expectedErr: errors.New(response.TransactionAlreadyExpired),
		},
		{
			name:          "error transaction has paid",
			transactionID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&model.Transaction{
					ExpiredAt: sql.NullTime{
						Valid: true,
						Time:  time.Now().Add(time.Hour),
					},
					PaidAt: sql.NullTime{
						Valid: true,
					},
					CardNumber: nil,
				}, nil)
			},
			expectedErr: errors.New(response.TransactionAlreadyFinished),
		},
		{
			name:          "error Card number nil",
			transactionID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&model.Transaction{
					ExpiredAt: sql.NullTime{
						Valid: true,
						Time:  time.Now().Add(time.Hour),
					},
					CardNumber: nil,
				}, nil)
			},
			expectedErr: errors.New(response.InvalidPaymentMethod),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()

			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.CreateWalletPayment(context.Background(), tc.transactionID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_GetTransactionByUserID(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		status      int
		pgn         *pagination.Pagination
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success GetTransactionByUserID",
			userID: "123456",
			status: 1,
			pgn:    &pagination.Pagination{},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempTotal := int64(10)
				r.On("GetTotalTransactionByUserID", mock.Anything, mock.Anything, mock.Anything).Return(tempTotal, nil)
				r.On("GetTransactionByUserID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*model.Transaction{{}}, nil)
				r.On("GetOrderDetailByTransactionID", mock.Anything, mock.Anything).Return([]*model.Order{}, nil)
				r.On("GetVoucherMarketplaceByID", mock.Anything, mock.Anything).Return(&model.Voucher{}, nil)
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
			_, err := u.GetTransactionByUserID(context.Background(), tc.userID, tc.status, tc.pgn)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_GetTransactionByID(t *testing.T) {
	testCase := []struct {
		name          string
		transactionID string
		mock          func(t *testing.T, r *mocks.Repository)
		expectedErr   error
	}{
		{
			name:          "success GetTransactionByID",
			transactionID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&model.Transaction{}, nil)
				r.On("GetOrderDetailByTransactionID", mock.Anything, mock.Anything).Return([]*model.Order{}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetTransactionByID(context.Background(), tc.transactionID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_UpdateTransaction(t *testing.T) {
	testCase := []struct {
		name          string
		transactionID string
		requestBody   body.SLPCallbackRequest
		mock          func(t *testing.T, r *mocks.Repository)
		expectedErr   error
	}{
		{
			name:          "success UpdateTransaction",
			transactionID: "123456",
			requestBody: body.SLPCallbackRequest{
				Status:  constant.SLPStatusPaid,
				Message: constant.SlPMessagePaid,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempCardNumber := "123456"
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&model.Transaction{
					ID: uuid.Nil,
					ExpiredAt: sql.NullTime{
						Valid: true,
						Time:  time.Now().Add(time.Hour),
					},
					CardNumber: &tempCardNumber,
					TotalPrice: 100,
					WalletID:   &uuid.Nil,
				}, nil)
				r.On("GetOrderByTransactionID", mock.Anything, mock.Anything).Return([]*model.OrderModel{{}}, nil)
				r.On("UpdateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateOrder", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(&model.Wallet{}, nil)
				r.On("InsertWalletHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateWalletBalance", mock.Anything, mock.Anything, mock.Anything).Return(nil)
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
			err := u.UpdateTransaction(context.Background(), tc.transactionID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_UpdateTransactionPaymentMethod(t *testing.T) {
	testCase := []struct {
		name          string
		transactionID string
		cardNumber    string
		mock          func(t *testing.T, r *mocks.Repository)
		expectedErr   error
	}{
		{
			name:          "success UpdateTransactionPaymentMethod",
			transactionID: "123456",
			cardNumber:    "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempCardNumber := "123456"
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&model.Transaction{
					ID: uuid.Nil,
					ExpiredAt: sql.NullTime{
						Valid: true,
						Time:  time.Now().Add(time.Hour),
					},
					CardNumber: &tempCardNumber,
					TotalPrice: 100,
					WalletID:   &uuid.Nil,
				}, nil)
				r.On("UpdateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil)
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
			err := u.UpdateTransactionPaymentMethod(context.Background(), tc.transactionID, tc.cardNumber)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_UpdateWalletTransaction(t *testing.T) {
	testCase := []struct {
		name          string
		transactionID string
		requestBody   body.SLPCallbackRequest
		mock          func(t *testing.T, r *mocks.Repository)
		expectedErr   error
	}{
		{
			name:          "success UpdateWalletTransaction cancel",
			transactionID: "123456",
			requestBody: body.SLPCallbackRequest{
				Status:  constant.SLPStatusCanceled,
				Message: constant.SLPMessageCanceled,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempCardNumber := "123456"
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&model.Transaction{
					ID: uuid.Nil,
					ExpiredAt: sql.NullTime{
						Valid: true,
						Time:  time.Now().Add(time.Hour),
					},
					CardNumber: &tempCardNumber,
					TotalPrice: 100,
					WalletID:   &uuid.Nil,
				}, nil)
				r.On("GetWalletUser", mock.Anything, mock.Anything).Return(&model.Wallet{}, nil)
				r.On("UpdateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:          "success UpdateWalletTransaction paid",
			transactionID: "123456",
			requestBody: body.SLPCallbackRequest{
				Status:  constant.SLPStatusPaid,
				Message: constant.SlPMessagePaid,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempCardNumber := "123456"
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&model.Transaction{
					ID: uuid.Nil,
					ExpiredAt: sql.NullTime{
						Valid: true,
						Time:  time.Now().Add(time.Hour),
					},
					CardNumber: &tempCardNumber,
					TotalPrice: 100,
					WalletID:   &uuid.Nil,
				}, nil)
				r.On("GetWalletUser", mock.Anything, mock.Anything).Return(&model.Wallet{}, nil)
				r.On("UpdateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("InsertWalletHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateWalletBalance", mock.Anything, mock.Anything, mock.Anything).Return(nil)
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
			err := u.UpdateWalletTransaction(context.Background(), tc.transactionID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_GetWallet(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success GetWallet",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(&model.Wallet{}, nil)
			},
			expectedErr: nil,
		},
		{
			name:   "error GetWallet",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "error GetWallet no sql row",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(response.WalletIsNotActivated),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetWallet(context.Background(), tc.userID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_GetWalletHistory(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		pgn         *pagination.Pagination
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success GetWalletHistory",
			userID: "123456",
			pgn:    &pagination.Pagination{},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempInt64 := int64(1)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(&model.Wallet{
					ID: uuid.Nil,
				}, nil)
				r.On("GetTotalWalletHistoryByWalletID", mock.Anything, mock.Anything).Return(tempInt64, nil)
				r.On("GetWalletHistoryByWalletID", mock.Anything, mock.Anything, mock.Anything).Return([]*body.HistoryWalletResponse{}, nil)
			},
			expectedErr: nil,
		},
		{
			name:   "error GetWallet",
			userID: "123456",
			pgn:    &pagination.Pagination{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "error GetWallet no sql row",
			userID: "123456",
			pgn:    &pagination.Pagination{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetWalletHistory(context.Background(), tc.userID, tc.pgn)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_GetDetailWalletHistory(t *testing.T) {
	testCase := []struct {
		name            string
		walletHistoryID string
		userID          string
		mock            func(t *testing.T, r *mocks.Repository)
		expectedErr     error
	}{
		{
			name:            "success GetDetailWalletHistory",
			userID:          "123456",
			walletHistoryID: "ab80c496-387b-4989-bf3b-a6f68a05940d",
			mock: func(t *testing.T, r *mocks.Repository) {
				tempUUID, _ := uuid.Parse("ab80c496-387b-4989-bf3b-a6f68a05940d")
				r.On("GetWalletHistoryByID", mock.Anything, mock.Anything).Return(&model.WalletHistory{
					ID:            tempUUID,
					WalletID:      tempUUID,
					TransactionID: uuid.Nil,
					From:          "ab80c496-387b-4989-bf3b-a6f68a05940d",
				}, nil)
				r.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&model.Transaction{}, nil)
				r.On("GetOrdersByTransactionID", mock.Anything, mock.Anything, mock.Anything).Return([]*model.Order{}, nil)
			},
			expectedErr: nil,
		},
		{
			name:            "error GetWallet",
			userID:          "123456",
			walletHistoryID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletHistoryByID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:            "error GetWallet no sql row",
			userID:          "123456",
			walletHistoryID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletHistoryByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(response.WalletHistoryNotFound),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetDetailWalletHistory(context.Background(), tc.walletHistoryID, tc.userID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_WalletStepUp(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		requestBody body.WalletStepUpRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "error WalletStepUp wrong error",
			userID: "123456",
			requestBody: body.WalletStepUpRequest{
				Amount: 1000,
				Pin:    "123456",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(&model.Wallet{
					Balance: 100000,
					PIN:     "123456",
				}, nil)
				r.On("UpdateWallet", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: errors.New(response.WalletPinIsInvalid),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.WalletStepUp(context.Background(), tc.userID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_ChangeWalletPinStepUp(t *testing.T) {
	passwordHash := "$2a$10$WKul/6gjYoYjOXuNVX4XGen1ZkWYb1PKFiI5vlZp5TFerZh6nTujG"
	testCase := []struct {
		name        string
		userID      string
		requestBody body.ChangeWalletPinStepUpRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "error ChangeWalletPinStepUp",
			userID: "123456",
			requestBody: body.ChangeWalletPinStepUpRequest{
				Password: "Tested8*",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(&model.Wallet{
					ID:      uuid.Nil,
					Balance: 100000,
					PIN:     "123456",
				}, nil)
				r.On("GetUserPasswordByID", mock.Anything, mock.Anything).Return(&model.User{
					Password: &passwordHash,
				}, nil)
			},
			expectedErr: nil,
		},
		{
			name:   "error ChangeWalletPinStepUp",
			userID: "123456",
			requestBody: body.ChangeWalletPinStepUpRequest{
				Password: "Tested8*",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "error ChangeWalletPinStepUp",
			userID: "123456",
			requestBody: body.ChangeWalletPinStepUpRequest{
				Password: "Tested8*",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(response.WalletIsNotActivated),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.ChangeWalletPinStepUp(context.Background(), tc.userID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_ChangeWalletPin(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		pin         string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success ChangeWalletPin",
			userID: "123456",
			pin:    "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(&model.Wallet{
					ID:      uuid.Nil,
					Balance: 100000,
					PIN:     "123456",
				}, nil)
				r.On("UpdateWalletPin", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "error ChangeWalletPin",
			userID: "123456",
			pin:    "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name:   "error ChangeWalletPin sql no rows",
			userID: "123456",
			pin:    "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetWalletByUserID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(response.WalletIsNotActivated),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			err := u.ChangeWalletPin(context.Background(), tc.userID, tc.pin)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_CreateTransaction(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		requestBody body.CreateTransactionRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success CreateTransaction",
			userID: "ab80c496-387b-4989-bf3b-a6f68a05940d",
			requestBody: body.CreateTransactionRequest{
				WalletID:             "c737a0f0-00e0-4dd5-89eb-dc44a9ea3413",
				CardNumber:           "",
				VoucherMarketplaceID: "f483bf6b-6293-428b-b87a-5892aacb4efa",
				CartItems: []body.CartItem{
					{
						ShopID:        "33ee7825-461b-40ca-8d6e-09ce7f2851fb",
						VoucherShopID: "b558e7e0-a39b-420e-ada4-ce4b180e7e9a",
						CourierID:     "1",
						CourierFee:    100,
						ProductDetails: []body.ProductDetail{
							{
								ID:       "c62f09d8-290d-496c-8413-e7d40ceaed05",
								Quantity: 1,
								SubPrice: 100,
								Note:     "noted",
							},
						},
					},
				},
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				tempUserID, _ := uuid.Parse("ab80c496-387b-4989-bf3b-a6f68a05940d")
				tempShopID, _ := uuid.Parse("33ee7825-461b-40ca-8d6e-09ce7f2851fb")
				tempPromoID, _ := uuid.Parse("7975f81f-5c51-46b8-8ea0-8362cc419c9d")
				tempProductID, _ := uuid.Parse("bd0b620f-1b07-46c7-aede-4de60d493450")
				tempProductDetailID, _ := uuid.Parse("e8590820-a776-470f-88e2-65961f0bd80e")
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{
					ID: tempUserID,
				}, nil)
				r.On("GetWalletUser", mock.Anything, mock.Anything).Return(&model.Wallet{
					UserID: tempUserID,
				}, nil)

				tempVoucherID, _ := uuid.Parse("f483bf6b-6293-428b-b87a-5892aacb4efa")
				tempVoucherID1, _ := uuid.Parse("f483bf6b-6293-428b-b87a-5892aacb4efa")
				tempVDiscountPercentage := float64(0)
				tempVDiscountFixPrice := float64(10)
				tempVMinProductPrice := float64(1)
				tempVMaxDiscountPrice := float64(100000)
				r.On("GetVoucherMarketplaceByID", mock.Anything, mock.Anything).Return(&model.Voucher{
					ID:                 tempVoucherID,
					ShopID:             uuid.Nil,
					Code:               "ABC",
					Quota:              10,
					ActivedDate:        time.Now(),
					ExpiredDate:        time.Now().Add(time.Hour),
					DiscountPercentage: &tempVDiscountPercentage,
					DiscountFixPrice:   &tempVDiscountFixPrice,
					MinProductPrice:    &tempVMinProductPrice,
					MaxDiscountPrice:   &tempVMaxDiscountPrice,
					CreatedAt:          time.Now(),
					UpdatedAt:          sql.NullTime{},
					DeletedAt:          sql.NullTime{},
				}, nil)
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{
					ID:     tempShopID,
					UserID: tempShopID,
				}, nil)
				r.On("GetVoucherShopByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Voucher{
					ID:                 tempVoucherID1,
					ShopID:             tempShopID,
					Code:               "ABCD",
					Quota:              2,
					ActivedDate:        time.Now(),
					ExpiredDate:        time.Now().Add(time.Hour),
					DiscountPercentage: &tempVDiscountPercentage,
					DiscountFixPrice:   &tempVDiscountFixPrice,
					MinProductPrice:    &tempVMinProductPrice,
					MaxDiscountPrice:   &tempVMaxDiscountPrice,
					CreatedAt:          time.Now(),
					UpdatedAt:          sql.NullTime{},
					DeletedAt:          sql.NullTime{},
				}, nil)
				r.On("GetCourierShopByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Courier{
					ID:          uuid.Nil,
					Name:        "jne",
					Code:        "JNE",
					Service:     "sell",
					Description: "asdkajsd",
					CreatedAt:   time.Now(),
					UpdatedAt:   sql.NullTime{},
				}, nil)

				r.On("GetProductDetailByID", mock.Anything, mock.Anything, mock.Anything).Once().Return(&model.ProductDetail{
					ID: tempProductDetailID,
				}, nil)
				r.On("GetCartItemUser", mock.Anything, mock.Anything, mock.Anything).Once().Return(&model.CartItem{}, nil)

				r.On("GetProductDetailByID", mock.Anything, mock.Anything, mock.Anything).Once().Return(&model.ProductDetail{
					ID:        tempProductDetailID,
					ProductID: tempProductID,
					Price:     10000,
					Stock:     10,
					Weight:    100,
					Size:      100,
					Hazardous: false,
					Condition: "good",
					BulkPrice: false,
					CreatedAt: time.Now(),
					UpdatedAt: sql.NullTime{},
					DeletedAt: sql.NullTime{},
				}, nil)
				r.On("GetCartItemUser", mock.Anything, mock.Anything, mock.Anything).Once().Return(&model.CartItem{}, nil)
				r.On("GetProductPromotionByProductID", mock.Anything, mock.Anything).Once().Return(&model.Promotion{
					ID:                 tempPromoID,
					Name:               "asd",
					ProductID:          tempProductID,
					DiscountPercentage: &tempVDiscountPercentage,
					DiscountFixPrice:   &tempVDiscountFixPrice,
					MinProductPrice:    &tempVMinProductPrice,
					MaxDiscountPrice:   &tempVMaxDiscountPrice,
					Quota:              3,
					MaxQuantity:        1,
					ActivedDate:        time.Now(),
					ExpiredDate:        time.Now(),
					CreatedAt:          time.Now(),
					UpdatedAt:          sql.NullTime{},
					DeletedAt:          sql.NullTime{},
				}, nil)

				tempTransactionID, _ := uuid.Parse("b7938be2-0d48-4ba8-af6b-465b79eb0891")
				tempOrderID, _ := uuid.Parse("ccbcbe3e-3cb1-4aae-abb8-e42d6bc587c0")
				r.On("GetAddressByBuyerID", mock.Anything, mock.Anything).Once().Return(&model.Address{}, nil)
				r.On("GetAddressBySellerID", mock.Anything, mock.Anything).Once().Return(&model.Address{}, nil)
				r.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Once().Return(&tempTransactionID, nil)
				r.On("UpdateVoucherQuota", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
				r.On("UpdateVoucherQuota", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
				r.On("UpdatePromotionQuota", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
				r.On("CreateOrder", mock.Anything, mock.Anything, mock.Anything).Once().Return(&tempOrderID, nil)
				r.On("CreateOrderItem", mock.Anything, mock.Anything, mock.Anything).Once().Return(&tempProductDetailID, nil)
				r.On("UpdateProductDetailStock", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
				r.On("DeleteCartItemByID", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
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
			_, err := u.CreateTransaction(context.Background(), tc.userID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_userUC_GetRefundOrder(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		orderID     string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr bool
	}{
		{name: "error GetOrderModelByID error",
			userID:  "123456",
			orderID: "abc",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetOrderModelByID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expectedErr: true,
		}, {name: "error orderData string not equal to userID",
			userID:  "123456",
			orderID: "abc",
			mock: func(t *testing.T, r *mocks.Repository) {
				testUUID, _ := uuid.Parse("123")
				r.On("GetOrderModelByID", mock.Anything, mock.Anything).Return(&model.OrderModel{
					UserID: testUUID,
				}, nil)
			},
			expectedErr: true,
		}, {name: "GetRefundByOrderID error",
			userID:  "ab80c496-387b-4989-bf3b-a6f68a05940d",
			orderID: "abc",
			mock: func(t *testing.T, r *mocks.Repository) {
				testUUID, _ := uuid.Parse("ab80c496-387b-4989-bf3b-a6f68a05940d")
				r.On("GetOrderModelByID", mock.Anything, mock.Anything).Return(&model.OrderModel{
					UserID: testUUID,
				}, nil)
				r.On("GetRefundOrderByOrderID", mock.Anything, mock.Anything).Return(nil, errors.New("hello"))
			},
			expectedErr: true,
		}, {name: "GetUserByID error",
			userID:  "ab80c496-387b-4989-bf3b-a6f68a05940d",
			orderID: "abc",
			mock: func(t *testing.T, r *mocks.Repository) {
				testUUID, _ := uuid.Parse("ab80c496-387b-4989-bf3b-a6f68a05940d")
				mockBool := true
				r.On("GetOrderModelByID", mock.Anything, mock.Anything).Return(&model.OrderModel{
					UserID: testUUID,
				}, nil)
				r.On("GetRefundOrderByOrderID", mock.Anything, mock.Anything).Return(&model.Refund{
					IsBuyerRefund: &mockBool,
				}, sql.ErrNoRows)
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, errors.New("halo"))
			},
			expectedErr: true,
		}, {name: "GetShopByID error",
			userID:  "ab80c496-387b-4989-bf3b-a6f68a05940d",
			orderID: "abc",
			mock: func(t *testing.T, r *mocks.Repository) {
				testUUID, _ := uuid.Parse("ab80c496-387b-4989-bf3b-a6f68a05940d")
				mockBool := true
				mockBoolFalse := false
				r.On("GetOrderModelByID", mock.Anything, mock.Anything).Return(&model.OrderModel{
					UserID: testUUID,
				}, nil)
				r.On("GetRefundOrderByOrderID", mock.Anything, mock.Anything).Return(&model.Refund{
					IsSellerRefund: &mockBool,
					IsBuyerRefund:  &mockBoolFalse,
				}, sql.ErrNoRows)
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(nil, errors.New("halo"))
			},
			expectedErr: true,
		}, {name: "GetUserByID error",
			userID:  "ab80c496-387b-4989-bf3b-a6f68a05940d",
			orderID: "abc",
			mock: func(t *testing.T, r *mocks.Repository) {
				testUUID, _ := uuid.Parse("ab80c496-387b-4989-bf3b-a6f68a05940d")
				mockBool := true
				mockBoolFalse := false
				r.On("GetOrderModelByID", mock.Anything, mock.Anything).Return(&model.OrderModel{
					UserID: testUUID,
				}, nil)
				r.On("GetRefundOrderByOrderID", mock.Anything, mock.Anything).Return(&model.Refund{
					IsSellerRefund: &mockBool,
					IsBuyerRefund:  &mockBoolFalse,
				}, sql.ErrNoRows)
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, nil)
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, errors.New("abc"))
			},
			expectedErr: true,
		}, {name: "GetRefundThreadByRefundID error",
			userID:  "ab80c496-387b-4989-bf3b-a6f68a05940d",
			orderID: "abc",
			mock: func(t *testing.T, r *mocks.Repository) {
				testUUID, _ := uuid.Parse("ab80c496-387b-4989-bf3b-a6f68a05940d")
				mockBool := true
				mockBoolFalse := false
				mockString := "abc"
				r.On("GetOrderModelByID", mock.Anything, mock.Anything).Return(&model.OrderModel{
					UserID: testUUID,
				}, nil)
				r.On("GetRefundOrderByOrderID", mock.Anything, mock.Anything).Return(&model.Refund{
					IsSellerRefund: &mockBool,
					IsBuyerRefund:  &mockBoolFalse,
				}, sql.ErrNoRows)
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, nil)
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{
					Username: &mockString,
					PhotoURL: &mockString,
				}, nil)
				r.On("GetRefundThreadByRefundID", mock.Anything, mock.Anything).Return(nil, errors.New("abc"))
			},
			expectedErr: true,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewUserUseCase(&config.Config{}, &postgre.TxRepo{}, r)
			tc.mock(t, r)
			_, err := u.GetRefundOrder(context.Background(), tc.userID, tc.orderID)
			if tc.expectedErr {
				assert.NotNil(t, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
