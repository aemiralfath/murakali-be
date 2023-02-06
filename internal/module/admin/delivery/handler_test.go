package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/admin/delivery/body"
	"murakali/internal/module/admin/mocks"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"murakali/pkg/pagination"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

func TestAdminHandlers_GetAllVoucher(t *testing.T) {
	testCase := []struct {
		name       string
		queries    map[string]string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Get Wallet",
			queries: map[string]string{
				"sort": "asc",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetAllVoucher", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Success Get Wallet",
			queries: map[string]string{
				"sort": "asc",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetAllVoucher", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Success Get Wallet",
			mock: func(s *mocks.UseCase) {
				s.On("GetAllVoucher", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Success Get Wallet",
			queries: map[string]string{
				"voucher_status": "1",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetAllVoucher", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Success Get Wallet",
			queries: map[string]string{
				"voucher_status": "2",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetAllVoucher", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Success Get Wallet",
			queries: map[string]string{
				"voucher_status": "3",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetAllVoucher", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Success Get Wallet",
			queries: map[string]string{
				"voucher_status": "4",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetAllVoucher", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Failed Get Wallet",
			mock: func(s *mocks.UseCase) {
				s.On("GetAllVoucher", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/admin/voucher", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")

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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetAllVoucher(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAdminHandlers_GetRefunds(t *testing.T) {
	testCase := []struct {
		name       string
		queries    map[string]string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Get Refund",
			queries: map[string]string{
				"sort": "asc",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetRefunds", mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Success Get Wallet",
			queries: map[string]string{
				"sort": "asc",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetRefunds", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Success Get Refund",
			queries: map[string]string{
				"sort": "desc",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetRefunds", mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Failed Get Refund",
			mock: func(s *mocks.UseCase) {
				s.On("GetRefunds", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/admin/refund", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")

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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetRefunds(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}

}

func TestAdminHandlers_CreateVoucher(t *testing.T) {
	invalidRequestBody := struct {
		Code               int    `json:"code"`
		Quota              string `json:"quota"`
		ActivedDate        int    `json:"actived_date"`
		ExpiredDate        int    `json:"expired_date"`
		DiscountPercentage string `json:"discount_percentage"`
		DiscountFixPrice   string `json:"discount_fix_price"`
		MinProductPrice    string `json:"min_product_price"`
		MaxDiscountPrice   string `json:"max_discount_price"`
		ActiveDateTime     string
		ExpiredDateTime    string
	}{123, "123", 123, 123, "123", "123", "123", "123", "123", "123"}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Create Voucher",
			body: body.CreateVoucherRequest{
				Code:               "123",
				Quota:              123,
				ActivedDate:        "02-01-2006 15:04:05",
				ExpiredDate:        "02-01-2006 15:04:05",
				DiscountPercentage: 123,
				DiscountFixPrice:   123,
				MinProductPrice:    123,
				MaxDiscountPrice:   123,
				ActiveDateTime:     time.Now(),
				ExpiredDateTime:    time.Now(),
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateVoucher", mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Success Create Voucher",
			body: body.CreateVoucherRequest{
				Code:               "123",
				Quota:              123,
				ActivedDate:        "02-01-2006 15:04:05",
				ExpiredDate:        "02-01-2006 15:04:05",
				DiscountPercentage: 123,
				DiscountFixPrice:   123,
				MinProductPrice:    123,
				MaxDiscountPrice:   123,
				ActiveDateTime:     time.Now(),
				ExpiredDateTime:    time.Now(),
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateVoucher", mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Failed Create Voucher",
			body: body.CreateVoucherRequest{
				Code:               "123",
				Quota:              123,
				ActivedDate:        "123",
				ExpiredDate:        "123",
				DiscountPercentage: 123,
				DiscountFixPrice:   123,
				MinProductPrice:    123,
				MaxDiscountPrice:   123,
				ActiveDateTime:     time.Now(),
				ExpiredDateTime:    time.Now(),
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name:       "Failed Create Voucher",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Failed Create Voucher",
			body: body.CreateVoucherRequest{
				Code:               "123",
				Quota:              123,
				ActivedDate:        "02-01-2006 15:04:05",
				ExpiredDate:        "02-01-2006 15:04:05",
				DiscountPercentage: 123,
				DiscountFixPrice:   123,
				MinProductPrice:    123,
				MaxDiscountPrice:   123,
				ActiveDateTime:     time.Now(),
				ExpiredDateTime:    time.Now(),
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateVoucher", mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
			},
			expected:   http.StatusInternalServerError,
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/admin/voucher", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CreateVoucher(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAdminHandlers_DeleteVoucher(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Delete Voucher",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteVoucher", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:  "Success Get Wallet",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteVoucher", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			param:      "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: false,
		},
		{
			name:  "Failed Delete Voucher",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteVoucher", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
			},
			expected:   http.StatusInternalServerError,
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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.DeleteVoucher(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAdminHandlers_UpdateVoucher(t *testing.T) {
	invalidRequestBody := struct {
		VoucherID          int    `json:"voucher_id"`
		Quota              string `json:"quota"`
		ActivedDate        int    `json:"actived_date"`
		ExpiredDate        int    `json:"expired_date"`
		DiscountPercentage string `json:"discount_percentage"`
		DiscountFixPrice   string `json:"discount_fix_price"`
		MinProductPrice    string `json:"min_product_price"`
		MaxDiscountPrice   string `json:"max_discount_price"`
		ActiveDateTime     int
		ExpiredDateTime    int
	}{123, "123", 123, 123, "123", "123", "123", "123", 123, 123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Update Voucher",
			body: body.UpdateVoucherRequest{
				VoucherID:          "8302755e-25c5-4523-8498-7dc8b9e3a098",
				Quota:              123,
				ActivedDate:        "02-01-2006 15:04:05",
				ExpiredDate:        "02-01-2006 15:04:05",
				DiscountPercentage: 123,
				DiscountFixPrice:   123,
				MinProductPrice:    123,
				MaxDiscountPrice:   123,
				ActiveDateTime:     time.Now(),
				ExpiredDateTime:    time.Now(),
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateVoucher", mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Success Get Wallet",
			body: body.UpdateVoucherRequest{
				VoucherID:          "8302755e-25c5-4523-8498-7dc8b9e3a098",
				Quota:              123,
				ActivedDate:        "02-01-2006 15:04:05",
				ExpiredDate:        "02-01-2006 15:04:05",
				DiscountPercentage: 123,
				DiscountFixPrice:   123,
				MinProductPrice:    123,
				MaxDiscountPrice:   123,
				ActiveDateTime:     time.Now(),
				ExpiredDateTime:    time.Now(),
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateVoucher", mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Failed Update Voucher",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: false,
		},
		{
			name: "Success Update Voucher",
			body: body.UpdateVoucherRequest{
				VoucherID:          "8302755e-25c5-4523-8498-7dc8b9e3a098",
				Quota:              123,
				ActivedDate:        "02-01-2006 15:04:05",
				ExpiredDate:        "02-01-2006 15:04:05",
				DiscountPercentage: 123,
				DiscountFixPrice:   123,
				MinProductPrice:    123,
				MaxDiscountPrice:   123,
				ActiveDateTime:     time.Now(),
				ExpiredDateTime:    time.Now(),
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateVoucher", mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
			},
			expected:   http.StatusInternalServerError,
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

			r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/admin/voucher"), bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPut(c, tc.body)

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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.UpdateVoucher(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAdminHandlers_GetDetailVoucher(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:       "Failed Get Voucher",
			param:      "8302755e-25c5-4523-8498-7dc8b9e3a09",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Success Get Voucher",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("GetDetailVoucher", mock.Anything, mock.Anything).Return(&model.Voucher{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:  "Success Get Voucher",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("GetDetailVoucher", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Failed Get Voucher",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("GetDetailVoucher", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/admin/voucher/%s", tc.param), nil)
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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetDetailVoucher(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAdminHandlers_GetCategories(t *testing.T) {
	testCase := []struct {
		name       string
		queries    map[string]string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Get Categories",
			mock: func(s *mocks.UseCase) {
				s.On("GetCategories", mock.Anything).Return([]*body.CategoryResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Failed Get Categories",
			mock: func(s *mocks.UseCase) {
				s.On("GetCategories", mock.Anything).Return(nil, fmt.Errorf("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/admin/categories", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")

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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetCategories(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}

}

func TestAdminHandlers_AddCategory(t *testing.T) {
	invalidRequestBody := struct {
		ID            int `json:"id"`
		ParentID      int `json:"parent_id" `
		ParentIDValue int
		Name          int `json:"name" `
		PhotoURL      int `json:"photo_url"`
		Level         int `json:"level" `
	}{123, 123, 123, 123, 123, 123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Create Voucher",
			body: body.CategoryRequest{
				ID:       "123",
				ParentID: "123",
				Name:     "123",
				PhotoURL: "http://example.com",
				Level:    "1",
			},
			mock: func(s *mocks.UseCase) {
				s.On("AddCategory", mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Success Get Wallet",
			body: body.CategoryRequest{
				ID:       "123",
				ParentID: "123",
				Name:     "123",
				PhotoURL: "http://example.com",
				Level:    "1",
			},
			mock: func(s *mocks.UseCase) {
				s.On("AddCategory", mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Failed Create Voucher",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Failed Create Voucher",
			body: body.CategoryRequest{
				ID:       "123",
				ParentID: "123",
				Name:     "123",
				PhotoURL: "http://example.com",
				Level:    "1",
			},
			mock: func(s *mocks.UseCase) {
				s.On("AddCategory", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expected:   http.StatusInternalServerError,
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/admin/voucher", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.AddCategory(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAdminHandlers_DeleteCategory(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Delete Category",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteCategory", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:  "Success Get Wallet",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteCategory", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			param:      "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: false,
		},
		{
			name:  "Failed Delete Category",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteCategory", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/admin/category/%s", tc.param), nil)

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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.DeleteCategory(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAdminHandlers_DeleteBanner(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Delete Banner",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteBanner", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:  "Success Get Wallet",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteBanner", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Unauthorized User",
			param:      "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: false,
		},
		{
			name:  "Failed Delete Banner",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteBanner", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/admin/Banner/%s", tc.param), nil)

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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.DeleteBanner(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAdminHandlers_EditCategory(t *testing.T) {
	invalidRequestBody := struct {
		ID            int `json:"id"`
		ParentID      int `json:"parent_id" `
		ParentIDValue int
		Name          int `json:"name" `
		PhotoURL      int `json:"photo_url"`
		Level         int `json:"level" `
	}{123, 123, 123, 123, 123, 123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Create Voucher",
			body: body.CategoryRequest{
				ID:       "123",
				ParentID: "123",
				Name:     "123",
				PhotoURL: "http://example.com",
				Level:    "1",
			},
			mock: func(s *mocks.UseCase) {
				s.On("EditCategory", mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Success Get Wallet",
			body: body.CategoryRequest{
				ID:       "123",
				ParentID: "123",
				Name:     "123",
				PhotoURL: "http://example.com",
				Level:    "1",
			},
			mock: func(s *mocks.UseCase) {
				s.On("EditCategory", mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Failed Create Voucher",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Failed Create Voucher",
			body: body.CategoryRequest{
				ID:       "123",
				ParentID: "123",
				Name:     "123",
				PhotoURL: "http://example.com",
				Level:    "1",
			},
			mock: func(s *mocks.UseCase) {
				s.On("EditCategory", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expected:   http.StatusInternalServerError,
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

			r := httptest.NewRequest(http.MethodPut, "/api/v1/admin/voucher", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.EditCategory(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAdminHandlers_GetBanner(t *testing.T) {
	testCase := []struct {
		name       string
		queries    map[string]string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Get Banner",
			mock: func(s *mocks.UseCase) {
				s.On("GetBanner", mock.Anything).Return([]*body.BannerResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Success Get Wallet",
			mock: func(s *mocks.UseCase) {
				s.On("GetBanner", mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Failed Get Banner",
			mock: func(s *mocks.UseCase) {
				s.On("GetBanner", mock.Anything).Return(nil, fmt.Errorf("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			r := httptest.NewRequest(http.MethodGet, "/api/v1/admin/Banner", nil)
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")

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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetBanner(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}

}

func TestAdminHandlers_AddBanner(t *testing.T) {
	invalidRequestBody := struct {
		Title    int    `json:"title" `
		Content  int    `json:"content"`
		ImageURL int    `json:"image_url"`
		PageURL  int    `json:"page_url"`
		IsActive string `json:"is_active"`
	}{123, 123, 123, 123, "123"}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Create Voucher",
			body: body.BannerRequest{
				Title:    "123",
				Content:  "123",
				ImageURL: "http://example.com",
				PageURL:  "http://example.com",
				IsActive: true,
			},
			mock: func(s *mocks.UseCase) {
				s.On("AddBanner", mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Success Get Wallet",
			body: body.BannerRequest{
				Title:    "123",
				Content:  "123",
				ImageURL: "http://example.com",
				PageURL:  "http://example.com",
				IsActive: true,
			},
			mock: func(s *mocks.UseCase) {
				s.On("AddBanner", mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Failed Create Voucher",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Failed Create Voucher",
			body: body.BannerRequest{
				Title:    "123",
				Content:  "123",
				ImageURL: "http://example.com",
				PageURL:  "http://example.com",
				IsActive: true,
			},
			mock: func(s *mocks.UseCase) {
				s.On("AddBanner", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expected:   http.StatusInternalServerError,
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/admin/voucher", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.AddBanner(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAdminHandlers_EditBanner(t *testing.T) {
	invalidRequestBody := struct {
		ID       int    `json:"id"`
		IsActive string `json:"is_active"`
	}{123, "123"}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Create Banner",
			body: body.BannerIDRequest{
				ID:       "123",
				IsActive: true,
			},
			mock: func(s *mocks.UseCase) {
				s.On("EditBanner", mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Success Get Wallet",
			body: body.BannerIDRequest{
				ID:       "123",
				IsActive: true,
			},
			mock: func(s *mocks.UseCase) {
				s.On("EditBanner", mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Failed Create Banner",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Failed Create Voucher",
			body: body.BannerIDRequest{
				ID:       "123",
				IsActive: true,
			},
			mock: func(s *mocks.UseCase) {
				s.On("EditBanner", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expected:   http.StatusInternalServerError,
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

			r := httptest.NewRequest(http.MethodPut, "/api/v1/admin/voucher", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)

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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.EditBanner(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAdminHandlers_RefundOrder(t *testing.T) {

	testCase := []struct {
		name       string
		param      string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Refund Order",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("RefundOrder", mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Failed Refund Order",
			param:      "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Failed Refund Order",
			param: "8302755e-25c5-4523-8498-7dc8b9e3a098",
			mock: func(s *mocks.UseCase) {
				s.On("RefundOrder", mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/admin/refund/%s", tc.param), nil)

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

			h := NewAdminHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.RefundOrder(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}
