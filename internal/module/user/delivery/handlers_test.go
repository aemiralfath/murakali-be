package delivery

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"murakali/config"
	"murakali/internal/constant"
	"murakali/internal/model"
	"murakali/internal/module/user/delivery/body"
	"murakali/internal/module/user/mocks"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"murakali/pkg/pagination"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func MockJsonPost(c *gin.Context, content interface{}) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
}

func MockJsonPatch(c *gin.Context, content interface{}) {
	c.Request.Method = "PATCH"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
}

func MockJsonPut(c *gin.Context, content interface{}) {
	c.Request.Method = "PUT"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
}

func TestUserHandlers_RegisterMerchant(t *testing.T) {
	invalidRequestBody := struct {
		ShopName int `json:"shop_name"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Register Merchant",
			body: body.RegisterMerchant{
				ShopName: "nama shop",
			},
			mock: func(s *mocks.UseCase) {
				s.On("RegisterMerchant", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Invalid Request Body shouldBind",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Invalid Request Body Field",
			body: body.RegisterMerchant{
				ShopName: "",
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Register Merchant Internal Error",
			body: body.RegisterMerchant{
				ShopName: "nama shop",
			},
			mock: func(s *mocks.UseCase) {
				s.On("RegisterMerchant", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Register Merchant Error Custom",
			body: body.RegisterMerchant{
				ShopName: "nama shop",
			},
			mock: func(s *mocks.UseCase) {
				s.On("RegisterMerchant", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/register-merchant", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)
			if tc.authorized {
				c.Set("userID", "123456")
			}

			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.RegisterMerchant(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_GetWallet(t *testing.T) {
	testCase := []struct {
		name       string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Get Wallet",
			mock: func(s *mocks.UseCase) {
				s.On("GetWallet", mock.Anything, mock.Anything).Return(&model.Wallet{
					ID:           uuid.Nil,
					UserID:       uuid.Nil,
					Balance:      0,
					PIN:          "",
					AttemptCount: 0,
					AttemptAt:    sql.NullTime{},
					UnlockedAt:   sql.NullTime{},
					ActiveDate:   sql.NullTime{},
					UpdatedAt:    sql.NullTime{},
				}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: " Get Wallet Internal Error",
			mock: func(s *mocks.UseCase) {
				s.On("GetWallet", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: " Get Wallet Custom  error",
			mock: func(s *mocks.UseCase) {
				s.On("GetWallet", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/user/wallet", nil)
			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetWallet(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_GetWalletHistory(t *testing.T) {
	testCase := []struct {
		name       string
		queries    map[string]string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:    "Success Get Wallet History",
			queries: map[string]string{},
			mock: func(s *mocks.UseCase) {
				s.On("GetWalletHistory", mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{
					Limit:      18,
					Page:       1,
					Sort:       "created_at desc",
					TotalRows:  1,
					TotalPages: 1,
					Rows:       nil,
				}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			queries:    nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:    "Get Wallet History Internal Error",
			queries: map[string]string{},
			mock: func(s *mocks.UseCase) {
				s.On("GetWalletHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Get Wallet History Custom  error",
			queries: map[string]string{
				"sort": "asc",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetWalletHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/user/address/", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			if tc.authorized {
				c.Set("userID", "123456")
			}

			if tc.queries != nil && len(tc.queries) > 0 {
				u := url.Values{}
				for key, value := range tc.queries {
					u.Set(key, value)
				}
				c.Request.URL.RawQuery = u.Encode()
			}

			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetWalletHistory(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_GetWalletHistoryByID(t *testing.T) {
	testCase := []struct {
		name       string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Get Wallet History By Wallet ID",
			mock: func(s *mocks.UseCase) {
				s.On("GetDetailWalletHistory", mock.Anything, mock.Anything, mock.Anything).Return(&body.DetailHistoryWalletResponse{
					ID:          "",
					Transaction: nil,
					From:        "",
					To:          "",
					Amount:      0,
					Description: "",
					CreatedAt:   time.Now(),
				}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},

		{
			name: "Get Wallet History By Wallet ID Internal Error",
			mock: func(s *mocks.UseCase) {
				s.On("GetDetailWalletHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Get Wallet History By Wallet ID Custom  error",
			mock: func(s *mocks.UseCase) {
				s.On("GetDetailWalletHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/user/wallet/history/:wallet_history_id", nil)
			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetWalletHistoryByID(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_TopUpWallet(t *testing.T) {
	invalidRequestBody := struct {
		CardNumber int    `json:"card_number"`
		Amount     string `json:"amount"`
	}{123, "123"}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success TopUp Wallet",
			body: body.TopUpWalletRequest{
				CardNumber: "123456",
				Amount:     10000,
			},
			mock: func(s *mocks.UseCase) {
				s.On("TopUpWallet", mock.Anything, mock.Anything, mock.Anything).Return("test", nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "TopUp Wallet error Shouldbind",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Invalid Request Body Validate TopUp Wallet ",
			body: body.TopUpWalletRequest{
				CardNumber: "",
				Amount:     0,
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "TopUp Wallet Internal Error",
			body: body.TopUpWalletRequest{
				CardNumber: "123456",
				Amount:     10000,
			},
			mock: func(s *mocks.UseCase) {
				s.On("TopUpWallet", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "TopUp Wallet Custom  error",
			body: body.TopUpWalletRequest{
				CardNumber: "123456",
				Amount:     10000,
			},
			mock: func(s *mocks.UseCase) {
				s.On("TopUpWallet", mock.Anything, mock.Anything, mock.Anything).Return("", httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPatch, "/api/v1/user/wallet", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPatch(c, tc.body)
			if tc.authorized {
				c.Set("userID", "123456")
			}

			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.TopUpWallet(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_ActivateWallet(t *testing.T) {
	invalidRequestBody := struct {
		Pin int `json:"pin"`
	}{123456}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Activate Wallet",
			body: body.ActivateWalletRequest{
				Pin: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ActivateWallet", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Activate Wallet error Shouldbind",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Invalid Request Body Validate Activate Wallet",
			body: body.ActivateWalletRequest{
				Pin: "",
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Activate Wallet Internal Error",
			body: body.ActivateWalletRequest{
				Pin: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ActivateWallet", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Activate Wallet Custom  error",
			body: body.ActivateWalletRequest{
				Pin: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ActivateWallet", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/wallet", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)
			if tc.authorized {
				c.Set("userID", "123456")
			}

			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.ActivateWallet(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_DeleteAddressByID(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Delete Address By ID",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			param:      "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "param id Parse Error",
			param:      "test",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Delete Address By ID Internal Error",
			param: uuid.Nil.String(),
			mock: func(s *mocks.UseCase) {
				s.On("DeleteAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name:  "Delete Address By ID Custom error",
			param: uuid.Nil.String(),
			mock: func(s *mocks.UseCase) {
				s.On("DeleteAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/user/address/%s", tc.param), nil)

			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "id",
						Value: tc.param,
					},
				}
			}

			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.DeleteAddressByID(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_GetAddressByID(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Get Address By ID",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Address{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			param:      "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "param id Parse Error",
			param:      "test",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Get Address By ID Internal Error",
			param: uuid.Nil.String(),
			mock: func(s *mocks.UseCase) {
				s.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name:  "Get Address By ID Custom error",
			param: uuid.Nil.String(),
			mock: func(s *mocks.UseCase) {
				s.On("GetAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/user/address/%s", tc.param), nil)

			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "id",
						Value: tc.param,
					},
				}
			}

			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetAddressByID(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_CreateAddress(t *testing.T) {
	invalidRequestBody := struct {
		Name int `json:"name"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Create Address",
			body: body.CreateAddressRequest{
				Name:          "test",
				ProvinceID:    1,
				CityID:        2,
				Province:      "Sumatera Selatan",
				City:          "Lahat",
				District:      "Lahat",
				SubDistrict:   "Pasar Lama",
				AddressDetail: "jalan Mayor Ruslam",
				ZipCode:       "31413",
				IsDefault:     true,
				IsShopDefault: false,
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateAddress", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Invalid Request Body shouldBind",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Invalid Request Body Field",
			body:       body.CreateAddressRequest{},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Create Address Internal Error",
			body: body.CreateAddressRequest{
				Name:          "test",
				ProvinceID:    1,
				CityID:        2,
				Province:      "Sumatera Selatan",
				City:          "Lahat",
				District:      "Lahat",
				SubDistrict:   "Pasar Lama",
				AddressDetail: "jalan Mayor Ruslam",
				ZipCode:       "31413",
				IsDefault:     true,
				IsShopDefault: false,
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateAddress", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Create Address Error Custom",
			body: body.CreateAddressRequest{
				Name:          "test",
				ProvinceID:    1,
				CityID:        2,
				Province:      "Sumatera Selatan",
				City:          "Lahat",
				District:      "Lahat",
				SubDistrict:   "Pasar Lama",
				AddressDetail: "jalan Mayor Ruslam",
				ZipCode:       "31413",
				IsDefault:     true,
				IsShopDefault: false,
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateAddress", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/address", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)
			if tc.authorized {
				c.Set("userID", "123456")
			}

			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CreateAddress(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_UpdateAddressByID(t *testing.T) {
	invalidRequestBody := struct {
		Name int `json:"name"`
	}{123}

	testCase := []struct {
		name       string
		param      string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success UpdateAddress By ID",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			body: body.UpdateAddressRequest{
				Name:          "test",
				ProvinceID:    1,
				CityID:        2,
				Province:      "Sumatera Selatan",
				City:          "Lahat",
				District:      "Lahat",
				SubDistrict:   "Pasar Lama",
				AddressDetail: "jalan Mayor Ruslam",
				ZipCode:       "31413",
				IsDefault:     true,
				IsShopDefault: false,
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateAddressByID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			param:      "",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Invalid Param ID uuid parse",
			param:      "test",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Invalid Request Body shouldBind",
			param:      "8302755e-25c5-4523-8498-7dc8b9e3a098",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Invalid Request Body Field",
			param:      "8302755e-25c5-4523-8498-7dc8b9e3a098",
			body:       body.UpdateAddressRequest{},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name:  "Update Address By ID Internal Error",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			body: body.UpdateAddressRequest{
				Name:          "test",
				ProvinceID:    1,
				CityID:        2,
				Province:      "Sumatera Selatan",
				City:          "Lahat",
				District:      "Lahat",
				SubDistrict:   "Pasar Lama",
				AddressDetail: "jalan Mayor Ruslam",
				ZipCode:       "31413",
				IsDefault:     true,
				IsShopDefault: false,
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateAddressByID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name:  "Update Address By ID Error Custom",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			body: body.UpdateAddressRequest{
				Name:          "test",
				ProvinceID:    1,
				CityID:        2,
				Province:      "Sumatera Selatan",
				City:          "Lahat",
				District:      "Lahat",
				SubDistrict:   "Pasar Lama",
				AddressDetail: "jalan Mayor Ruslam",
				ZipCode:       "31413",
				IsDefault:     true,
				IsShopDefault: false,
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateAddressByID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/user/address/%s", tc.param), bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPut(c, tc.body)
			if tc.authorized {
				c.Set("userID", "123456")
			}

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "id",
						Value: tc.param,
					},
				}
			}

			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.UpdateAddressByID(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_GetAddress(t *testing.T) {
	testCase := []struct {
		name       string
		queries    map[string]string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:    "Success Get Address",
			queries: map[string]string{},
			mock: func(s *mocks.UseCase) {
				s.On("GetAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			queries:    nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:    "Get Address Internal Error",
			queries: map[string]string{},
			mock: func(s *mocks.UseCase) {
				s.On("GetAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Get Address Error Custom",
			queries: map[string]string{
				"is_default":      "true",
				"is_shop_default": "true",
				"sort":            "asc",
				"sortBy":          "province",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/user/address/", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			if tc.authorized {
				c.Set("userID", "123456")
			}

			if tc.queries != nil && len(tc.queries) > 0 {
				u := url.Values{}
				for key, value := range tc.queries {
					u.Set(key, value)
				}
				c.Request.URL.RawQuery = u.Encode()
			}

			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetAddress(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_GetOrder(t *testing.T) {
	testCase := []struct {
		name       string
		queries    map[string]string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
		ParseUUID  bool
	}{
		{
			name:    "Success Get Order",
			queries: map[string]string{},
			mock: func(s *mocks.UseCase) {
				s.On("GetOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
			ParseUUID:  true,
		},
		{
			name:       "Unauthorized User",
			queries:    nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
			ParseUUID:  false,
		},
		{
			name:       "userID Error uuid Parse",
			queries:    nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
			ParseUUID:  false,
		},
		{
			name:    "Get Order Internal Error",
			queries: map[string]string{},
			mock: func(s *mocks.UseCase) {
				s.On("GetOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
			ParseUUID:  true,
		},
		{
			name: "Get Order Error Custom",
			queries: map[string]string{
				"sort": "asc",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
			ParseUUID:  true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/user/order", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			if tc.authorized {
				if tc.ParseUUID {
					c.Set("userID", "8302755e-25c5-4523-8498-7dc8b9e3a098")
				} else {
					c.Set("userID", "123456")
				}
			}

			if tc.queries != nil && len(tc.queries) > 0 {
				u := url.Values{}
				for key, value := range tc.queries {
					u.Set(key, value)
				}
				c.Request.URL.RawQuery = u.Encode()
			}

			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetOrder(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_GetOrderByOrderID(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Get Order By OrderID",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("GetOrderByOrderID", mock.Anything, mock.Anything).Return(&model.Order{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Invalid Order ID UUID parse",
			param:      "test",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: false,
		},
		{
			name:  "Get Order By OrderID Internal Error",
			param: uuid.Nil.String(),
			mock: func(s *mocks.UseCase) {
				s.On("GetOrderByOrderID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name:  "Get Order By OrderID Error Custom",
			param: uuid.Nil.String(),
			mock: func(s *mocks.UseCase) {
				s.On("GetOrderByOrderID", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/user/order/%s", tc.param), nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			if tc.authorized {
				c.Set("userID", "123456")
			}

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "order_id",
						Value: tc.param,
					},
				}
			}

			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetOrderByOrderID(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_ChangeOrderStatus(t *testing.T) {
	invalidRequestBody := struct {
		OrderID int `json:"order_id"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
		parseUUID  bool
	}{
		{
			name: "Success Change Order Status",
			body: body.ChangeOrderStatusRequest{
				OrderID:       "8302755e-25c5-4523-8498-7dc8b9e3a098",
				OrderStatusID: 1,
			},
			mock: func(s *mocks.UseCase) {
				s.On("ChangeOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
			parseUUID:  true,
		},
		{
			name:       "Invalid Request Body shouldBind error",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: false,
			parseUUID:  false,
		},
		{
			name:       "Invalid Request Body Validate",
			body:       body.ChangeOrderStatusRequest{},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: false,
			parseUUID:  false,
		},
		{
			name: "Unauthorized User",
			body: body.ChangeOrderStatusRequest{
				OrderID:       "8302755e-25c5-4523-8498-7dc8b9e3a098",
				OrderStatusID: 1,
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
			parseUUID:  false,
		},
		{
			name: "Invalid User ID UUID parse",
			body: body.ChangeOrderStatusRequest{
				OrderID:       "8302755e-25c5-4523-8498-7dc8b9e3a098",
				OrderStatusID: 1,
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
			parseUUID:  false,
		},
		{
			name: "Change Order Status Internal Error",
			body: body.ChangeOrderStatusRequest{
				OrderID:       "8302755e-25c5-4523-8498-7dc8b9e3a098",
				OrderStatusID: 1,
			},
			mock: func(s *mocks.UseCase) {
				s.On("ChangeOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
			parseUUID:  true,
		},
		{
			name: "Change Order Status Error Custom",
			body: body.ChangeOrderStatusRequest{
				OrderID:       "8302755e-25c5-4523-8498-7dc8b9e3a098",
				OrderStatusID: 1,
			},
			mock: func(s *mocks.UseCase) {
				s.On("ChangeOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
			parseUUID:  true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPatch, "/api/v1/user/order-status", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPatch(c, tc.body)
			if tc.authorized {
				if tc.parseUUID {
					c.Set("userID", "8302755e-25c5-4523-8498-7dc8b9e3a098")
				} else {
					c.Set("userID", "123456")
				}
			}

			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.ChangeOrderStatus(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_GetTransactionDetailByID(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Get Transaction Detail By ID",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("GetTransactionDetailByID", mock.Anything, mock.Anything, mock.Anything).Return(&body.TransactionDetailResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			param:      "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:  "Get Transaction Detail By ID Internal Error",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("GetTransactionDetailByID", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name:  "Get Transaction Detail By ID Error Custom",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("GetTransactionDetailByID", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/user/transaction/detail/%s", tc.param), nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")

			if tc.authorized {
				c.Set("userID", "123456")
			}

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "transaction_id",
						Value: tc.param,
					},
				}
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetTransactionDetailByID(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_ChangeTransactionPaymentMethod(t *testing.T) {
	invalidRequestBody := struct {
		TransactionID int `json:"transaction_id"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Change Transaction Payment Method",
			body: body.ChangeTransactionPaymentMethodReq{
				TransactionID: "123456",
				CardNumber:    "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateTransactionPaymentMethod", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Invalid Request Body ShouldBind Error",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Invalid Request Body Validate",
			body:       body.ChangeTransactionPaymentMethodReq{},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Change Transaction Payment Method Internal Error",
			body: body.ChangeTransactionPaymentMethodReq{
				TransactionID: "123456",
				CardNumber:    "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateTransactionPaymentMethod", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Change Transaction Payment Method Error Custom",
			body: body.ChangeTransactionPaymentMethodReq{
				TransactionID: "123456",
				CardNumber:    "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateTransactionPaymentMethod", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPut, "/api/v1/user/transaction", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPut(c, tc.body)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.ChangeTransactionPaymentMethod(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_EditUser(t *testing.T) {
	invalidRequestBody := struct {
		Username int `json:"username"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Edit User",
			body: body.EditUserRequest{
				Username:  "juww",
				PhoneNo:   "81770811472",
				FullName:  "test juww",
				Gender:    "M",
				BirthDate: "02-01-2006",
			},
			mock: func(s *mocks.UseCase) {
				s.On("EditUser", mock.Anything, mock.Anything, mock.Anything).Return(&model.User{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Invalid Request Body ShouldBind Error",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Invalid Request Body Validate",
			body:       body.EditUserRequest{},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Unauthorized User",
			body: body.EditUserRequest{
				Username:  "juww",
				PhoneNo:   "81770811472",
				FullName:  "test juww",
				Gender:    "M",
				BirthDate: "02-01-2006",
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "Edit User Internal Error",
			body: body.EditUserRequest{
				Username:  "juww",
				PhoneNo:   "81770811472",
				FullName:  "test juww",
				Gender:    "M",
				BirthDate: "02-01-2006",
			},
			mock: func(s *mocks.UseCase) {
				s.On("EditUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Edit User Error Custom",
			body: body.EditUserRequest{
				Username:  "juww",
				PhoneNo:   "81770811472",
				FullName:  "test juww",
				Gender:    "M",
				BirthDate: "02-01-2006",
			},
			mock: func(s *mocks.UseCase) {
				s.On("EditUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPut, "/api/v1/user/profile", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPut(c, tc.body)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.EditUser(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_EditEmail(t *testing.T) {
	invalidRequestBody := struct {
		Email int `json:"email"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Edit Email",
			body: body.EditEmailRequest{
				Email: "email@google.com",
			},
			mock: func(s *mocks.UseCase) {
				s.On("EditEmail", mock.Anything, mock.Anything, mock.Anything).Return(&model.User{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Invalid Request Body ShouldBind Error",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Invalid Request Body Validate",
			body:       body.EditUserRequest{},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Unauthorized User",
			body: body.EditEmailRequest{
				Email: "email@google.com",
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "Edit Email Internal Error",
			body: body.EditEmailRequest{
				Email: "email@google.com",
			},
			mock: func(s *mocks.UseCase) {
				s.On("EditEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Edit Email Error Custom",
			body: body.EditEmailRequest{
				Email: "email@google.com",
			},
			mock: func(s *mocks.UseCase) {
				s.On("EditEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/email", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.EditEmail(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_EditEmailUser(t *testing.T) {
	invalidRequestBody := struct {
		Email int `json:"email"`
	}{123}

	testCase := []struct {
		name       string
		queries    map[string]interface{}
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Edit Email User",
			queries: map[string]interface{}{
				"email": "email@google.com",
				"code":  "123456",
			},
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("EditEmailUser", mock.Anything, mock.Anything, mock.Anything).Return(&model.User{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Invalid Request Body ShouldBind Error",
			queries:    nil,
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: false,
		},
		{
			name:       "Invalid Request Body Validate",
			queries:    map[string]interface{}{},
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: false,
		},
		{
			name: "Unauthorized User",
			queries: map[string]interface{}{
				"email": "email@google.com",
				"code":  "123456",
			},
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "Edit Email User Internal Error",
			queries: map[string]interface{}{
				"email": "email@google.com",
				"code":  "123456",
			},
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("EditEmailUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Edit Email User Error Custom",
			queries: map[string]interface{}{
				"email": "email@google.com",
				"code":  "123456",
			},
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("EditEmailUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/user/email", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")

			if tc.authorized {
				c.Set("userID", "123456")
			}

			if tc.queries != nil && len(tc.queries) > 0 {
				u := url.Values{}
				for key, value := range tc.queries {
					u.Set(key, fmt.Sprintf("%v", value))
				}
				c.Request.URL.RawQuery = u.Encode()
			}

			if tc.body != nil {
				MockJsonPost(c, tc.body)
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.EditEmailUser(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_GetSealabsPay(t *testing.T) {
	testCase := []struct {
		name       string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Get Sealabs Pay",
			mock: func(s *mocks.UseCase) {
				s.On("GetSealabsPay", mock.Anything, mock.Anything).Return([]*model.SealabsPay{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "Get Sealabs Pay Internal Error",
			mock: func(s *mocks.UseCase) {
				s.On("GetSealabsPay", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Get Sealabs Pay Error Custom",
			mock: func(s *mocks.UseCase) {
				s.On("GetSealabsPay", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/user/sealab-pay", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetSealabsPay(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_AddSealabsPay(t *testing.T) {
	invalidRequestBody := struct {
		Name int `json:"name"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Add Sealabs Pay",
			body: body.AddSealabsPayRequest{
				CardNumber: "1234567890123456",
				Name:       "jww",
				IsDefault:  true,
				ActiveDate: "02-01-2006 15:04:05",
			},
			mock: func(s *mocks.UseCase) {
				s.On("AddSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Invalid Request Body ShouldBind Error",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Invalid Request Body Validate",
			body:       body.AddSealabsPayRequest{},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Add Sealabs Pay Internal Error",
			body: body.AddSealabsPayRequest{
				CardNumber: "1234567890123456",
				Name:       "jww",
				IsDefault:  true,
				ActiveDate: "02-01-2006 15:04:05",
			},
			mock: func(s *mocks.UseCase) {
				s.On("AddSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Add Sealabs Pay Error Custom",
			body: body.AddSealabsPayRequest{
				CardNumber: "1234567890123456",
				Name:       "jww",
				IsDefault:  true,
				ActiveDate: "02-01-2006 15:04:05",
			},
			mock: func(s *mocks.UseCase) {
				s.On("AddSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/sealab-pay", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.AddSealabsPay(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_userHandlers_PatchSealabsPay(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Patch Sealabs Pay",
			param: "1234567890123456",
			mock: func(s *mocks.UseCase) {
				s.On("PatchSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			param:      "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Invalid Request Body Validate",
			param:      "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name:  "Patch Sealabs Pay Internal Error",
			param: "1234567890123456",
			mock: func(s *mocks.UseCase) {
				s.On("PatchSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name:  "Patch Sealabs Pay Error Custom",
			param: "1234567890123456",
			mock: func(s *mocks.UseCase) {
				s.On("PatchSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, err := json.Marshal(tc.param)
			if err != nil {
				t.Error(err)
			}
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/user/sealab-pay/%s", tc.param), bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPatch(c, tc.param)

			if tc.authorized {
				c.Set("userID", "123456")
			}
			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "cardNumber",
						Value: tc.param,
					},
				}
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.PatchSealabsPay(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_DeleteSealabsPay(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Delete Sealabs Pay",
			param: "1234567890123456",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			param:      "1234567890123456",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Invalid Request Body Validate",
			param:      "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name:  "Delete Sealabs Pay Internal Error",
			param: "1234567890123456",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name:  "Delete Sealabs Pay Error Custom",
			param: "1234567890123456",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteSealabsPay", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/user/sealab-pay/%s", tc.param), nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")

			if tc.authorized {
				c.Set("userID", "123456")
			}
			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "cardNumber",
						Value: tc.param,
					},
				}
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.DeleteSealabsPay(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_GetUserProfile(t *testing.T) {
	testCase := []struct {
		name       string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Get User Profile",
			mock: func(s *mocks.UseCase) {
				s.On("GetUserProfile", mock.Anything, mock.Anything).Return(&body.ProfileResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "Get User Profile Internal Error",
			mock: func(s *mocks.UseCase) {
				s.On("GetUserProfile", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Get User Profile Error Custom",
			mock: func(s *mocks.UseCase) {
				s.On("GetUserProfile", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/user/profile", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetUserProfile(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

// func Test_userHandlers_UploadProfilePicture(t *testing.T) {
// 	// invalidRequestBody := struct {
// 	// 	Img int `form:"file"`
// 	// }{123}

// 	testCase := []struct {
// 		name       string
// 		body       interface{}
// 		mock       func(s *mocks.UseCase)
// 		expected   int
// 		authorized bool
// 	}{
// 		{
// 			name: "Success Upload Profile Picture",
// 			body: nil,
// 			mock: func(s *mocks.UseCase) {
// 				s.On("UploadProfilePicture", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 			},
// 			expected:   http.StatusOK,
// 			authorized: true,
// 		},
// 	}

// 	for _, tc := range testCase {
// 		t.Run(tc.name, func(t *testing.T) {
// 			rr := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(rr)

// 			r := httptest.NewRequest(http.MethodGet, "/api/v1/user/profile", nil)
// 			r.Header = make(http.Header)

// 			c.Request = r
// 			c.Request.Header.Set("Content-Type", "application/json")

// 			if tc.authorized {
// 				c.Set("userID", "123456")
// 			}

// 			c.Request = r
// 			s := mocks.NewUseCase(t)

// 			cfg := &config.Config{
// 				Logger: config.LoggerConfig{
// 					Development:       true,
// 					DisableCaller:     false,
// 					DisableStacktrace: false,
// 					Encoding:          "json",
// 					Level:             "info",
// 				},
// 			}

// 			appLogger := logger.NewAPILogger(cfg)
// 			appLogger.InitLogger()

// 			h := NewUserHandlers(cfg, s, appLogger)

// 			tc.mock(s)
// 			h.UploadProfilePicture(c)

// 			assert.Equal(t, rr.Code, tc.expected)
// 		})
// 	}
// }

func TestUserHandlers_VerifyPasswordChange(t *testing.T) {
	testCase := []struct {
		name       string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Verify Password Change",
			mock: func(s *mocks.UseCase) {
				s.On("VerifyPasswordChange", mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "Verify Password Change Internal Error",
			mock: func(s *mocks.UseCase) {
				s.On("VerifyPasswordChange", mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Verify Password Change Error Custom",
			mock: func(s *mocks.UseCase) {
				s.On("VerifyPasswordChange", mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/password", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, nil)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.VerifyPasswordChange(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_userHandlers_VerifyOTP(t *testing.T) {
	invalidRequestBody := struct {
		OTP int `json:"otp"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Verify Password Change",
			body: body.VerifyOTPRequest{
				OTP: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("VerifyOTP", mock.Anything, mock.Anything, mock.Anything).Return("test", nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Invalid Request Body ShouldBind Error",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Invalid Request Body Validate",
			body: body.VerifyOTPRequest{
				OTP: "",
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Verify Password Change Internal Error",
			body: body.VerifyOTPRequest{
				OTP: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("VerifyOTP", mock.Anything, mock.Anything, mock.Anything).Return("test", errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Verify Password Change Error Custom",
			body: body.VerifyOTPRequest{
				OTP: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("VerifyOTP", mock.Anything, mock.Anything, mock.Anything).Return("test", httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/verify", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.VerifyOTP(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_CompletedRejectedRefund(t *testing.T) {
	testCase := []struct {
		name     string
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name: "Success Completed Rejected Refund",
			mock: func(s *mocks.UseCase) {
				s.On("CompletedRejectedRefund", mock.Anything).Return(nil)
			},
			expected: http.StatusOK,
		},
		{
			name: "Completed Rejected Refund Internal Error",
			mock: func(s *mocks.UseCase) {
				s.On("CompletedRejectedRefund", mock.Anything).Return(errors.New("test"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name: "Completed Rejected Refund Error Custom",
			mock: func(s *mocks.UseCase) {
				s.On("CompletedRejectedRefund", mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/rejected-refund", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CompletedRejectedRefund(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_userHandlers_ChangePassword(t *testing.T) {
	// invalidRequestBody := struct {
	// 	NewPassword int `json:"password"`
	// }{123}

	testCase := []struct {
		name     string
		cookie   string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
	}{
		// {
		// 	name:   "Success Change Password",
		// 	cookie: "8302755e-25c5-4523-8498-7dc8b9e3a098",
		// 	body: body.ChangePasswordRequest{
		// 		NewPassword: "Tested9*",
		// 	},
		// 	mock: func(s *mocks.UseCase) {
		// 		s.On("ChangePassword", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		// 	},
		// 	expected: http.StatusOK,
		// },
		{
			name:     "Error Cookie",
			cookie:   "",
			body:     nil,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusForbidden,
		},
		{
			name:     "Error Cookie",
			cookie:   "test",
			body:     nil,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusForbidden,
		},
		// {
		// 	name:   "Change Password Internal Error",
		// 	cookie: "8302755e-25c5-4523-8498-7dc8b9e3a098",
		// 	body: body.ChangePasswordRequest{
		// 		NewPassword: "Tested9*",
		// 	},
		// 	mock: func(s *mocks.UseCase) {
		// 		s.On("ChangePassword", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
		// 	},
		// 	expected: http.StatusInternalServerError,
		// },
		// {
		// 	name:   "Change Password Error Custom",
		// 	cookie: "8302755e-25c5-4523-8498-7dc8b9e3a098",
		// 	body: body.ChangePasswordRequest{
		// 		NewPassword: "Tested9*",
		// 	},
		// 	mock: func(s *mocks.UseCase) {
		// 		s.On("ChangePassword", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
		// 	},
		// 	expected: http.StatusBadRequest,
		// },
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodPatch, "/api/v1/user/password", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPatch(c, tc.body)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}

			if tc.cookie != "" {
				// changePasswordToken, errToken := jwt.GenerateJWTChangePasswordToken(tc.cookie, cfg)
				// if errToken != nil {
				// 	log.Errorf("HandlerUser_Test, Error: %s", errToken)
				// 	return
				// }
				c.Request.Header.Set("Cookie", fmt.Sprintf("%s=%s", constant.ChangePasswordTokenCookie, tc.cookie))
				changePasswordToken, err := c.Cookie(constant.ChangePasswordTokenCookie)
				fmt.Println("test token", changePasswordToken)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.ChangePassword(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_userHandlers_WalletStepUp(t *testing.T) {
	invalidRequestBody := struct {
		Pin int `json:"pin"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Wallet StepUp",
			body: body.WalletStepUpRequest{
				Pin:    "123456",
				Amount: 123456,
			},
			mock: func(s *mocks.UseCase) {
				s.On("WalletStepUp", mock.Anything, mock.Anything, mock.Anything).Return("test", nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Invalid Request Body ShouldBind Error",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Invalid Request Body Validate",
			body:       body.WalletStepUpRequest{},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Wallet StepUp Internal Error",
			body: body.WalletStepUpRequest{
				Pin:    "123456",
				Amount: 123456,
			},
			mock: func(s *mocks.UseCase) {
				s.On("WalletStepUp", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Wallet StepUp Error Custom",
			body: body.WalletStepUpRequest{
				Pin:    "123456",
				Amount: 123456,
			},
			mock: func(s *mocks.UseCase) {
				s.On("WalletStepUp", mock.Anything, mock.Anything, mock.Anything).Return("", httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/wallet/step-up/pin", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}
			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.WalletStepUp(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestUserHandlers_ChangeWalletPinStepUp(t *testing.T) {
	invalidRequestBody := struct {
		Password int `json:"password"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Change Wallet Pin StepUp",
			body: body.ChangeWalletPinStepUpRequest{
				Password: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ChangeWalletPinStepUp", mock.Anything, mock.Anything, mock.Anything).Return("test", nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Invalid Request Body ShouldBind Error",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Invalid Request Body Validate",
			body:       body.WalletStepUpRequest{},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Change Wallet Pin StepUp Internal Error",
			body: body.ChangeWalletPinStepUpRequest{
				Password: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ChangeWalletPinStepUp", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Change Wallet Pin StepUp Error Custom",
			body: body.ChangeWalletPinStepUpRequest{
				Password: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ChangeWalletPinStepUp", mock.Anything, mock.Anything, mock.Anything).Return("", httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/wallet/step-up/password", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}
			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.ChangeWalletPinStepUp(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

// func Test_userHandlers_ChangeWalletPin(t *testing.T) {
// 	type fields struct {
// 		cfg    *config.Config
// 		userUC user.UseCase
// 		logger logger.Logger
// 	}
// 	type args struct {
// 		c *gin.Context
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := &userHandlers{
// 				cfg:    tt.fields.cfg,
// 				userUC: tt.fields.userUC,
// 				logger: tt.fields.logger,
// 			}
// 			h.ChangeWalletPin(tt.args.c)
// 		})
// 	}
// }

func Test_userHandlers_CreateSLPPayment(t *testing.T) {
	invalidRequestBody := struct {
		TransactionID int `json:"transaction_id"`
	}{123}

	testCase := []struct {
		name     string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name: "Success Change Wallet Pin StepUp",
			body: body.CreatePaymentRequest{
				TransactionID: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateSLPPayment", mock.Anything, mock.Anything).Return("test", nil)
			},
			expected: http.StatusOK,
		},
		{
			name:     "Invalid Request Body ShouldBind Error",
			body:     invalidRequestBody,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
		},
		{
			name:     "Invalid Request Body Validate",
			body:     body.CreatePaymentRequest{},
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusUnprocessableEntity,
		},
		{
			name: "Change Wallet Pin StepUp Internal Error",
			body: body.CreatePaymentRequest{
				TransactionID: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateSLPPayment", mock.Anything, mock.Anything).Return("", errors.New("test"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name: "Change Wallet Pin StepUp Error Custom",
			body: body.CreatePaymentRequest{
				TransactionID: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateSLPPayment", mock.Anything, mock.Anything).Return("", httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/wallet/step-up/password", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}
			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CreateSLPPayment(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

// func Test_userHandlers_CreateWalletPayment(t *testing.T) {
// 	type fields struct {
// 		cfg    *config.Config
// 		userUC user.UseCase
// 		logger logger.Logger
// 	}
// 	type args struct {
// 		c *gin.Context
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := &userHandlers{
// 				cfg:    tt.fields.cfg,
// 				userUC: tt.fields.userUC,
// 				logger: tt.fields.logger,
// 			}
// 			h.CreateWalletPayment(tt.args.c)
// 		})
// 	}
// }

func Test_userHandlers_SLPPaymentCallback(t *testing.T) {
	invalidRequestBody := struct {
		TxnID int `json:"txn_id"`
	}{123}

	testCase := []struct {
		name     string
		param    string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name:  "Success SLP Payment Callback",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			body: body.SLPCallbackRequest{
				Signature: "8343c086714a9950026bdc6d0c195fcee3141f7b0b1eddfb465bb8fda283076b",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected: http.StatusOK,
		},
		{
			name:     "Invalid Request Body ShouldBind Error",
			param:    "",
			body:     invalidRequestBody,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
		},
		{
			name:     "Invalid Request Body Validate",
			param:    "",
			body:     body.CreatePaymentRequest{},
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusUnprocessableEntity,
		},
		{
			name:  "Invalid Request Body Validate",
			param: "test",
			body: body.SLPCallbackRequest{
				Signature: "8343c086714a9950026bdc6d0c195fcee3141f7b0b1eddfb465bb8fda283076b",
			},
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
		},
		{
			name:  "SLP Payment Callback Internal Error",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			body: body.SLPCallbackRequest{
				Signature: "8343c086714a9950026bdc6d0c195fcee3141f7b0b1eddfb465bb8fda283076b",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name:  "SLP Payment Callback Error Custom",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			body: body.SLPCallbackRequest{
				Signature: "8343c086714a9950026bdc6d0c195fcee3141f7b0b1eddfb465bb8fda283076b",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/user/transaction/slp-payment/%s", tc.param), bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "id",
						Value: tc.param,
					},
				}
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}
			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.SLPPaymentCallback(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_userHandlers_WalletPaymentCallback(t *testing.T) {
	invalidRequestBody := struct {
		TxnID int `json:"txn_id"`
	}{123}

	testCase := []struct {
		name     string
		param    string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name:  "Success Wallet Payment Callback",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			body: body.SLPCallbackRequest{
				Signature: "8343c086714a9950026bdc6d0c195fcee3141f7b0b1eddfb465bb8fda283076b",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateWalletTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected: http.StatusOK,
		},
		{
			name:     "Invalid Request Body ShouldBind Error",
			param:    "",
			body:     invalidRequestBody,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
		},
		{
			name:     "Invalid Request Body Validate",
			param:    "",
			body:     body.CreatePaymentRequest{},
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusUnprocessableEntity,
		},
		{
			name:  "Invalid Request Body Validate",
			param: "test",
			body: body.SLPCallbackRequest{
				Signature: "8343c086714a9950026bdc6d0c195fcee3141f7b0b1eddfb465bb8fda283076b",
			},
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
		},
		{
			name:  "Wallet Payment Callback Internal Error",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			body: body.SLPCallbackRequest{
				Signature: "8343c086714a9950026bdc6d0c195fcee3141f7b0b1eddfb465bb8fda283076b",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateWalletTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name:  "Wallet Payment Callback Error Custom",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			body: body.SLPCallbackRequest{
				Signature: "8343c086714a9950026bdc6d0c195fcee3141f7b0b1eddfb465bb8fda283076b",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateWalletTransaction", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/user/transaction/wallet-payment/%s", tc.param), bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "id",
						Value: tc.param,
					},
				}
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}
			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.WalletPaymentCallback(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_userHandlers_GetTransactions(t *testing.T) {
	testCase := []struct {
		name       string
		queries    map[string]interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Get Transactions",
			mock: func(s *mocks.UseCase) {
				s.On("GetTransactionByUserID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "Get Transactions Internal Error",
			mock: func(s *mocks.UseCase) {
				s.On("GetTransactionByUserID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Get Transactions Error Custom",
			queries: map[string]interface{}{
				"sort":   "asc",
				"status": constant.OrderStatusWaitingToPay,
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetTransactionByUserID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodGet, "/api/v1/user/transaction", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")

			if tc.authorized {
				c.Set("userID", "123456")
			}

			if tc.queries != nil && len(tc.queries) > 0 {
				u := url.Values{}
				for key, value := range tc.queries {
					u.Set(key, fmt.Sprintf("%v", value))
				}
				c.Request.URL.RawQuery = u.Encode()
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}
			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetTransactions(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_userHandlers_GetTransaction(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Get Transaction",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("GetTransactionByID", mock.Anything, mock.Anything).Return(&body.GetTransactionByIDResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			param:      "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:  "Get Transaction Internal Error",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("GetTransactionByID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name:  "Get Transaction Error Custom",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("GetTransactionByID", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodGet, "/api/v1/user/transaction/:id", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")

			if tc.authorized {
				c.Set("userID", "123456")
			}

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "id",
						Value: tc.param,
					},
				}
			}
			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}
			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetTransaction(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_userHandlers_CreateRefundUser(t *testing.T) {
	invalidRequestBody := struct {
		OrderID int `json:"order_id"`
	}{123}
	imgTest := "image.jpg"

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Create Refund User",
			body: body.CreateRefundUserRequest{
				OrderID:        "8302755e-25c5-4523-8498-7dc8b9e3a098",
				Reason:         "reason test",
				Image:          &imgTest,
				IsSellerRefund: false,
				IsBuyerRefund:  true,
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateRefundUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Invalid Request Body ShouldBind Error",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Invalid Request Body Validate",
			body:       body.CreateRefundUserRequest{},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Create Refund User Internal Error",
			body: body.CreateRefundUserRequest{
				OrderID:        "8302755e-25c5-4523-8498-7dc8b9e3a098",
				Reason:         "reason test",
				Image:          &imgTest,
				IsSellerRefund: false,
				IsBuyerRefund:  true,
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateRefundUser", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Create Refund User Error Custom",
			body: body.CreateRefundUserRequest{
				OrderID:        "8302755e-25c5-4523-8498-7dc8b9e3a098",
				Reason:         "reason test",
				Image:          &imgTest,
				IsSellerRefund: false,
				IsBuyerRefund:  true,
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateRefundUser", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/refund", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}
			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CreateRefundUser(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_userHandlers_GetRefundOrder(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Get Refund Order",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("GetRefundOrder", mock.Anything, mock.Anything, mock.Anything).Return(&body.GetRefundThreadResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			param:      "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Invalid Request Body ShouldBind Error",
			param:      "test",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Get Refund Order Internal Error",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("GetRefundOrder", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name:  "Get Refund Order Error Custom",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("GetRefundOrder", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/user/refund/%s", tc.param), nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")

			if tc.authorized {
				c.Set("userID", "123456")
			}

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "refund_id",
						Value: tc.param,
					},
				}
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}
			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetRefundOrder(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_userHandlers_CreateRefundThreadUser(t *testing.T) {
	invalidRequestBody := struct {
		RefundID int `json:"refund_id"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Create Refund Thread User",
			body: body.CreateRefundThreadRequest{
				RefundID: "8302755e-25c5-4523-8498-7dc8b9e3a098",
				IsSeller: false,
				IsBuyer:  true,
				Text:     "testTest",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateRefundThreadUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Invalid Request Body ShouldBind Error",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Invalid Request Body Validate",
			body:       body.CreateRefundThreadRequest{},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Create Refund Thread User Internal Error",
			body: body.CreateRefundThreadRequest{
				RefundID: "8302755e-25c5-4523-8498-7dc8b9e3a098",
				IsSeller: false,
				IsBuyer:  true,
				Text:     "testTest",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateRefundThreadUser", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Create Refund Thread User Error Custom",
			body: body.CreateRefundThreadRequest{
				RefundID: "8302755e-25c5-4523-8498-7dc8b9e3a098",
				IsSeller: false,
				IsBuyer:  true,
				Text:     "testTest",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateRefundThreadUser", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/refund-thread", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}
			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CreateRefundThreadUser(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_userHandlers_CreateTransaction(t *testing.T) {
	invalidRequestBody := struct {
		WalletID int `json:"wallet_id"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Create Transaction",
			body: body.CreateTransactionRequest{
				WalletID:             "8302755e-25c5-4523-8498-7dc8b9e3a098",
				CardNumber:           "",
				VoucherMarketplaceID: "",
				CartItems:            []body.CartItem{},
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Return("test", nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Invalid Request Body ShouldBind Error",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: false,
		},
		{
			name:       "Invalid Request Body Validate",
			body:       body.CreateTransactionRequest{},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: false,
		},
		{
			name: "Unauthorized User",
			body: body.CreateTransactionRequest{
				WalletID:             "8302755e-25c5-4523-8498-7dc8b9e3a098",
				CardNumber:           "",
				VoucherMarketplaceID: "",
				CartItems:            []body.CartItem{},
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "Create Transaction Internal Error",
			body: body.CreateTransactionRequest{
				WalletID:             "8302755e-25c5-4523-8498-7dc8b9e3a098",
				CardNumber:           "",
				VoucherMarketplaceID: "",
				CartItems:            []body.CartItem{},
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Create Transaction Error Custom",
			body: body.CreateTransactionRequest{
				WalletID:             "8302755e-25c5-4523-8498-7dc8b9e3a098",
				CardNumber:           "",
				VoucherMarketplaceID: "",
				CartItems:            []body.CartItem{},
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Return("", httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/transaction", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}
			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CreateTransaction(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_userHandlers_ChangeWalletPinStepUpEmail(t *testing.T) {
	testCase := []struct {
		name       string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Change Wallet Pin StepUp Email",
			mock: func(s *mocks.UseCase) {
				s.On("ChangeWalletPinStepUpEmail", mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "Change Wallet Pin StepUp Email Internal Error",
			mock: func(s *mocks.UseCase) {
				s.On("ChangeWalletPinStepUpEmail", mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Change Wallet Pin StepUp Email Error Custom",
			mock: func(s *mocks.UseCase) {
				s.On("ChangeWalletPinStepUpEmail", mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/wallet/step-up/email", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, nil)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}
			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.ChangeWalletPinStepUpEmail(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_userHandlers_ChangeWalletPinStepUpVerify(t *testing.T) {
	invalidRequestBody := struct {
		OTP int `json:"otp" form:"otp"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Change Wallet Pin StepUp Verify",
			body: body.VerifyOTPRequest{
				OTP: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ChangeWalletPinStepUpVerify", mock.Anything, mock.Anything, mock.Anything).Return("test", nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Invalid Request Body ShouldBind Error",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Invalid Request Body Validate",
			body:       body.VerifyOTPRequest{},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Change Wallet Pin StepUp Verify Internal Error",
			body: body.VerifyOTPRequest{
				OTP: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ChangeWalletPinStepUpVerify", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Change Wallet Pin StepUp Verify Error Custom",
			body: body.VerifyOTPRequest{
				OTP: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ChangeWalletPinStepUpVerify", mock.Anything, mock.Anything, mock.Anything).Return("", httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			jsonValue, err := json.Marshal(tc.body)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodPost, "/api/v1/user/wallet/step-up/verify", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			c.Request = r
			s := mocks.NewUseCase(t)

			cfg := &config.Config{
				Logger: config.LoggerConfig{
					Development:       true,
					DisableCaller:     false,
					DisableStacktrace: false,
					Encoding:          "json",
					Level:             "info",
				},
			}
			appLogger := logger.NewAPILogger(cfg)
			appLogger.InitLogger()

			h := NewUserHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.ChangeWalletPinStepUpVerify(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}
