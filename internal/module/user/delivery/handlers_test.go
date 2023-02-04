package delivery

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/user/delivery/body"
	"murakali/internal/module/user/mocks"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"murakali/pkg/pagination"
	"net/http"
	"net/http/httptest"
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
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Get Wallet History",
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
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},

		{
			name: "Get Wallet History Internal Error",
			mock: func(s *mocks.UseCase) {
				s.On("GetWalletHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Get Wallet History Custom  error",
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
		params     string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:   "Success Delete Address By ID",
			params: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			params:     "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "param id Parse Error",
			params:     "test",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:   "Delete Address By ID Internal Error",
			params: uuid.Nil.String(),
			mock: func(s *mocks.UseCase) {
				s.On("DeleteAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name:   "Delete Address By ID Custom error",
			params: uuid.Nil.String(),
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

			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/user/address/%s", tc.params), nil)

			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", "123456")
			}

			if tc.params != "" {
				c.Params = []gin.Param{
					{
						Key:   "id",
						Value: tc.params,
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
