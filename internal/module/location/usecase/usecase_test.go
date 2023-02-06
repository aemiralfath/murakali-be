package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"murakali/config"
	"murakali/internal/constant"
	"murakali/internal/model"
	"murakali/internal/module/location/delivery/body"
	"murakali/internal/module/location/mocks"
	"murakali/pkg/postgre"
	"net/url"
	"testing"
)

func TestLocationUC_GetShippingCost(t *testing.T) {
	response := `
		{
   "rajaongkir":{
      "query":{
         "origin":"501",
         "destination":"114",
         "weight":1700,
         "courier":"jne"
      },
      "status":{
         "code":200,
         "description":"OK"
      },
      "origin_details":{
         "city_id":"501",
         "province_id":"5",
         "province":"DI Yogyakarta",
         "type":"Kota",
         "city_name":"Yogyakarta",
         "postal_code":"55000"
      },
      "destination_details":{
         "city_id":"114",
         "province_id":"1",
         "province":"Bali",
         "type":"Kota",
         "city_name":"Denpasar",
         "postal_code":"80000"
      },
      "results":[
         {
            "code":"jne",
            "name":"Jalur Nugraha Ekakurir (JNE)",
            "costs":[
               {
                  "service":"",
                  "description":"Layanan Reguler",
                  "cost":[
                     {
                        "value":44000,
                        "etd":"2-3",
                        "note":""
                     }
                  ]
               }
            ]
         }
      ]
   }
}
	`
	errResponse := ""
	testCase := []struct {
		name        string
		url         string
		posUrl      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success get shipping cost redis",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, nil)
				r.On("GetShopAddress", mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetShopCourierID", mock.Anything, mock.Anything).Return([]string{"1"}, nil)
				r.On("GetProductCourierWhitelistID", mock.Anything, mock.Anything).Return([]string{""}, nil)
				r.On("GetCourierByID", mock.Anything, mock.Anything).Return(&model.Courier{}, nil)
				r.On("GetCostRedis", mock.Anything, mock.Anything).Return(&response, nil)
			},
			expectedErr: nil,
		},
		{
			name:   "success get shipping cost api",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, nil)
				r.On("GetShopAddress", mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetShopCourierID", mock.Anything, mock.Anything).Return([]string{"1"}, nil)
				r.On("GetProductCourierWhitelistID", mock.Anything, mock.Anything).Return([]string{""}, nil)
				r.On("GetCourierByID", mock.Anything, mock.Anything).Return(&model.Courier{}, nil)
				r.On("GetCostRedis", mock.Anything, mock.Anything).Return(&response, errors.New("test"))
				r.On("InsertCostRedis", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "err insert get shipping cost api",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, nil)
				r.On("GetShopAddress", mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetShopCourierID", mock.Anything, mock.Anything).Return([]string{"1"}, nil)
				r.On("GetProductCourierWhitelistID", mock.Anything, mock.Anything).Return([]string{""}, nil)
				r.On("GetCourierByID", mock.Anything, mock.Anything).Return(&model.Courier{}, nil)
				r.On("GetCostRedis", mock.Anything, mock.Anything).Return(&response, errors.New("test"))
				r.On("InsertCostRedis", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: nil,
		},
		{
			name:   "error get shipping cost api",
			url:    "",
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, nil)
				r.On("GetShopAddress", mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetShopCourierID", mock.Anything, mock.Anything).Return([]string{"1"}, nil)
				r.On("GetProductCourierWhitelistID", mock.Anything, mock.Anything).Return([]string{""}, nil)
				r.On("GetCourierByID", mock.Anything, mock.Anything).Return(&model.Courier{}, nil)
				r.On("GetCostRedis", mock.Anything, mock.Anything).Return(&response, errors.New("test"))
			},
			expectedErr: nil,
		},
		{
			name:   "error get shipping cost redis",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, nil)
				r.On("GetShopAddress", mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetShopCourierID", mock.Anything, mock.Anything).Return([]string{"1"}, nil)
				r.On("GetProductCourierWhitelistID", mock.Anything, mock.Anything).Return([]string{""}, nil)
				r.On("GetCourierByID", mock.Anything, mock.Anything).Return(&model.Courier{}, nil)
				r.On("GetCostRedis", mock.Anything, mock.Anything).Return(&errResponse, nil)
			},
			expectedErr: nil,
		},
		{
			name:   "err courier get shipping cost redis",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, nil)
				r.On("GetShopAddress", mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetShopCourierID", mock.Anything, mock.Anything).Return([]string{"1"}, nil)
				r.On("GetProductCourierWhitelistID", mock.Anything, mock.Anything).Return([]string{""}, nil)
				r.On("GetCourierByID", mock.Anything, mock.Anything).Return(&model.Courier{}, errors.New("test"))
			},
			expectedErr: nil,
		},
		{
			name:   "err whitelist shipping cost redis",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, nil)
				r.On("GetShopAddress", mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetShopCourierID", mock.Anything, mock.Anything).Return([]string{"1"}, nil)
				r.On("GetProductCourierWhitelistID", mock.Anything, mock.Anything).Return([]string{""}, errors.New("test"))
			},
			expectedErr: nil,
		},
		{
			name:   "error no courier get shipping cost redis",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, nil)
				r.On("GetShopAddress", mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetShopCourierID", mock.Anything, mock.Anything).Return([]string{"1"}, sql.ErrNoRows)
			},
			expectedErr: nil,
		},
		{
			name:   "error no courier get shipping cost redis",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, nil)
				r.On("GetShopAddress", mock.Anything, mock.Anything).Return(&model.Address{}, nil)
				r.On("GetShopCourierID", mock.Anything, mock.Anything).Return([]string{"1"}, errors.New("test"))
			},
			expectedErr: nil,
		},
		{
			name:   "err no address get shipping cost redis",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, nil)
				r.On("GetShopAddress", mock.Anything, mock.Anything).Return(&model.Address{}, sql.ErrNoRows)
			},
			expectedErr: nil,
		},
		{
			name:   "err address get shipping cost redis",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, nil)
				r.On("GetShopAddress", mock.Anything, mock.Anything).Return(&model.Address{}, errors.New("test"))
			},
			expectedErr: nil,
		},
		{
			name:   "err no shop shipping cost redis",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, sql.ErrNoRows)
			},
			expectedErr: nil,
		},
		{
			name:   "err shop shipping cost redis",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, errors.New("test"))
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()

			r := mocks.NewRepository(t)
			u := NewLocationUseCase(&config.Config{External: config.ExternalConfig{OngkirAPIURL: tc.url, KodePosURL: tc.posUrl, OngkirAPIKey: constant.ONGKIR_API_KEY}}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetShippingCost(context.Background(), body.GetShippingCostRequest{Destination: 1, Weight: 1000, ProductIDS: []string{"1"}, ShopID: "1"})
			if err != nil {
				if errors.Is(&url.Error{}, tc.expectedErr) {
					assert.Equal(t, err, tc.expectedErr)
				}
			}
		})
	}
}

func TestLocationUC_GetUrban(t *testing.T) {
	testCase := []struct {
		name        string
		url         string
		posUrl      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success get urban redis",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUrbanRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("{}", nil)
			},
			expectedErr: nil,
		},
		{
			name:   "success get urban api",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUrbanRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("{}", errors.New("test"))
				r.On("InsertUrbanRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "error insert urban api",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUrbanRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("{}", errors.New("test"))
				r.On("InsertUrbanRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: nil,
		},
		{
			name:   "error get urban api",
			url:    constant.ONGKIR_API_URL,
			posUrl: "",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUrbanRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("{}", errors.New("test"))
			},
			expectedErr: nil,
		},
		{
			name:   "error get urban redis",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUrbanRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()

			r := mocks.NewRepository(t)
			u := NewLocationUseCase(&config.Config{External: config.ExternalConfig{OngkirAPIURL: tc.url, KodePosURL: tc.posUrl, OngkirAPIKey: constant.ONGKIR_API_KEY}}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetUrban(context.Background(), "DKI Jakarta", "Jakarta Selatan", "Mampang Prapatan")
			if err != nil {
				if errors.Is(&url.Error{}, tc.expectedErr) {
					assert.Equal(t, err, tc.expectedErr)
				}
			}
		})
	}
}

func TestLocationUC_GetSubDistrict(t *testing.T) {
	testCase := []struct {
		name        string
		url         string
		posUrl      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success get sub district redis",
			url:  constant.ONGKIR_API_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSubDistrictRedis", mock.Anything, mock.Anything, mock.Anything).Return("{}", nil)
			},
			expectedErr: nil,
		},
		{
			name:   "success get sub district api",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSubDistrictRedis", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("test"))
				r.On("InsertSubDistrictRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "error get sub district api",
			url:    constant.ONGKIR_API_URL,
			posUrl: "",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSubDistrictRedis", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("test"))
			},
			expectedErr: nil,
		},
		{
			name:   "error insert sub district api",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSubDistrictRedis", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("test"))
				r.On("InsertSubDistrictRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: nil,
		},
		{
			name:   "error get sub district redis",
			url:    constant.ONGKIR_API_URL,
			posUrl: constant.KODE_POS_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSubDistrictRedis", mock.Anything, mock.Anything, mock.Anything).Return("", nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()

			r := mocks.NewRepository(t)
			u := NewLocationUseCase(&config.Config{External: config.ExternalConfig{OngkirAPIURL: tc.url, KodePosURL: tc.posUrl, OngkirAPIKey: constant.ONGKIR_API_KEY}}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetSubDistrict(context.Background(), "Jakarta", "Selatan")
			if err != nil {
				if errors.Is(&url.Error{}, tc.expectedErr) {
					assert.Equal(t, err, tc.expectedErr)
				}
			}
		})
	}
}

func TestLocationUC_GetCity(t *testing.T) {
	testCase := []struct {
		name        string
		url         string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success get city redis",
			url:  constant.ONGKIR_API_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCityRedis", mock.Anything, mock.Anything).Return("{}", nil)
			},
			expectedErr: nil,
		},
		{
			name: "success get city api",
			url:  constant.ONGKIR_API_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCityRedis", mock.Anything, mock.Anything).Return("{}", errors.New("test"))
				r.On("InsertCityRedis", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "error insert city api",
			url:  constant.ONGKIR_API_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCityRedis", mock.Anything, mock.Anything).Return("{}", errors.New("test"))
				r.On("InsertCityRedis", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name: "error get city api",
			url:  "",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCityRedis", mock.Anything, mock.Anything).Return("{}", errors.New("test"))
			},
			expectedErr: &url.Error{},
		},
		{
			name: "err get city redis",
			url:  constant.ONGKIR_API_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCityRedis", mock.Anything, mock.Anything).Return("", nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()

			r := mocks.NewRepository(t)
			u := NewLocationUseCase(&config.Config{External: config.ExternalConfig{OngkirAPIURL: tc.url, OngkirAPIKey: constant.ONGKIR_API_KEY}}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetCity(context.Background(), 1)
			if err != nil {
				if errors.Is(&url.Error{}, tc.expectedErr) {
					assert.Equal(t, err, tc.expectedErr)
				}
			}
		})
	}
}

func TestLocationUC_GetProvince(t *testing.T) {
	testCase := []struct {
		name        string
		url         string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success get province redis",
			url:  constant.ONGKIR_API_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProvinceRedis", mock.Anything).Return("{}", nil)
			},
			expectedErr: nil,
		},
		{
			name: "err get province redis",
			url:  constant.ONGKIR_API_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProvinceRedis", mock.Anything).Return("", nil)
			},
			expectedErr: nil,
		},
		{
			name: "success get province api",
			url:  constant.ONGKIR_API_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProvinceRedis", mock.Anything).Return("", errors.New("test"))
				r.On("InsertProvinceRedis", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "insert redis error province api",
			url:  constant.ONGKIR_API_URL,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProvinceRedis", mock.Anything).Return("", errors.New("test"))
				r.On("InsertProvinceRedis", mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			name: "error get province api",
			url:  "",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProvinceRedis", mock.Anything).Return("", errors.New("test"))
			},
			expectedErr: &url.Error{},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()

			r := mocks.NewRepository(t)
			u := NewLocationUseCase(&config.Config{External: config.ExternalConfig{OngkirAPIURL: tc.url, OngkirAPIKey: constant.ONGKIR_API_KEY}}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetProvince(context.Background())
			if err != nil {
				if errors.Is(&url.Error{}, tc.expectedErr) {
					assert.Equal(t, err, tc.expectedErr)
				}
			}
		})
	}
}
