package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/seller/delivery/body"
	"murakali/internal/module/seller/mocks"
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

func MockJsonPut(c *gin.Context, content interface{}) {
	c.Request.Method = "PUT"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
}

func Test_sellerHandlers_GetPerformance(t *testing.T) {
	testCase := []struct {
		name       string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
		userID     string
	}{
		{
			name: "Success Get Performance",
			mock: func(s *mocks.UseCase) {
				s.On("GetPerformance", mock.Anything, mock.Anything, mock.Anything).Return(&body.SellerPerformance{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:       "Unauthorized User",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
			userID:     "",
		},
		{
			name: "UUID Invalid",
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6acddd",
		},
		{
			name: "Error Get Performance",
			mock: func(s *mocks.UseCase) {
				s.On("GetPerformance", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4"},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/seller/performance", nil)
			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", tc.userID)
			}
			uuid.Parse(tc.userID)

			c.Set("userIDString", tc.userID)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetPerformance(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_GetAllSeller(t *testing.T) {
	testCase := []struct {
		name     string
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name: "Success Get All Seller",
			mock: func(s *mocks.UseCase) {
				s.On("GetAllSeller", mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected: http.StatusOK,
		},
		{
			name: "Error Get All Seller",
			mock: func(s *mocks.UseCase) {
				s.On("GetAllSeller", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/seller/", nil)
			r.Header = make(http.Header)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetAllSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_GetOrder(t *testing.T) {
	testCase := []struct {
		name       string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
		userID     string
	}{
		{
			name: "Success Get Order",
			mock: func(s *mocks.UseCase) {
				s.On("GetOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Success Get Order with filter",
			mock: func(s *mocks.UseCase) {
				s.On("GetOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:       "Unauthorized User",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
			userID:     "",
		},
		{
			name: "UUID Invalid",
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6acddd",
		},
		{
			name: "Error Get Order",
			mock: func(s *mocks.UseCase) {
				s.On("GetOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/seller/order", nil)
			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", tc.userID)
			}
			uuid.Parse(tc.userID)

			c.Set("userIDString", tc.userID)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetOrder(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_ChangeOrderStatus(t *testing.T) {
	invalidRequestBody := struct {
		OrderID int `json:"order_id"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
		userID     string
	}{
		{
			name: "Success Change Order Status",
			body: body.ChangeOrderStatusRequest{
				OrderID:       "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
				OrderStatusID: "1",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ChangeOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:       "body is nil",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:       "Invalid request",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Unauthorized User",
			body: body.ChangeOrderStatusRequest{
				OrderID:       "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
				OrderStatusID: "1",
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
			userID:     "",
		},
		{
			name: "UUID Invalid",
			body: body.ChangeOrderStatusRequest{
				OrderID:       "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
				OrderStatusID: "1",
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "1",
		},
		{
			name: "Error Change Order Status",
			body: body.ChangeOrderStatusRequest{
				OrderID:       "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
				OrderStatusID: "1",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ChangeOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Error Change Order Status HTTPError",
			body: body.ChangeOrderStatusRequest{
				OrderID:       "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
				OrderStatusID: "1",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ChangeOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
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

			r := httptest.NewRequest(http.MethodPatch, "/api/v1/seller/order-status", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)
			if tc.authorized {
				c.Set("userID", tc.userID)
			}
			uuid.Parse(tc.userID)

			c.Set("userIDString", tc.userID)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.ChangeOrderStatus(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_CancelOrderStatus(t *testing.T) {
	invalidRequestBody := struct {
		OrderID int `json:"order_id"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
		userID     string
	}{
		{
			name: "Success Change Order Status",
			body: body.CancelOrderStatus{
				OrderID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
				CancelNotes: "1",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CancelOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:       "body is nil",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:       "Invalid request",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Unauthorized User",
			body: body.CancelOrderStatus{
				OrderID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
				CancelNotes: "1",
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
			userID:     "",
		},
		{
			name: "Error Change Order Status",
			body: body.CancelOrderStatus{
				OrderID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
				CancelNotes: "1",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CancelOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Error Change Order Status HTTPError",
			body: body.CancelOrderStatus{
				OrderID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
				CancelNotes: "1",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CancelOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
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

			r := httptest.NewRequest(http.MethodPatch, "/api/v1/seller/order-cancel", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)
			if tc.authorized {
				c.Set("userID", tc.userID)
			}
			uuid.Parse(tc.userID)

			c.Set("userIDString", tc.userID)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CancelOrderStatus(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_GetOrderByOrderID(t *testing.T) {
	testCase := []struct {
		name     string
		param    string
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name:  "Success Get Order By Order ID",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetOrderByOrderID", mock.Anything, mock.Anything).Return(&model.Order{}, nil)
			},
			expected: http.StatusOK,
		},
		{
			name:  "Error Get Order By Order ID",
			param: "1",
			mock: func(s *mocks.UseCase) {
			},
			expected: http.StatusBadRequest,
		},
		{
			name:  "Error Get Order By Order ID HTTPError",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetOrderByOrderID", mock.Anything, mock.Anything).Return(&model.Order{}, httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
		{
			name:  "Error Get Order By Order ID Error",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetOrderByOrderID", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/seller/order%s", tc.param), nil)
			r.Header = make(http.Header)

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "order_id",
						Value: tc.param,
					},
				}
			}

			uuid.Parse(tc.param)

			c.Set("orderID", tc.param)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetOrderByOrderID(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_GetSellerBySellerID(t *testing.T) {
	testCase := []struct {
		name     string
		param    string
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name:  "Success Get Seller By Seller ID",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetSellerBySellerID", mock.Anything, mock.Anything).Return(&body.SellerResponse{}, nil)
			},
			expected: http.StatusOK,
		},
		{
			name:  "Error Get Seller By Seller ID",
			param: "1",
			mock: func(s *mocks.UseCase) {
			},
			expected: http.StatusBadRequest,
		},
		{
			name:  "Error Get Seller By Seller ID HTTPError",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetSellerBySellerID", mock.Anything, mock.Anything).Return(&body.SellerResponse{}, httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
		{
			name:  "Error Get Seller By Seller ID Error",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetSellerBySellerID", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/seller/%s", tc.param), nil)
			r.Header = make(http.Header)

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "seller_id",
						Value: tc.param,
					},
				}
			}

			uuid.Parse(tc.param)

			c.Set("sellerID", tc.param)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetSellerBySellerID(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_GetSellerByUserID(t *testing.T) {
	testCase := []struct {
		name     string
		param    string
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name:  "Success Get Seller By Seller ID",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetSellerByUserID", mock.Anything, mock.Anything).Return(&body.SellerResponse{}, nil)
			},
			expected: http.StatusOK,
		},
		{
			name:  "Error Get Seller By Seller ID",
			param: "1",
			mock: func(s *mocks.UseCase) {
			},
			expected: http.StatusBadRequest,
		},
		{
			name:  "Error Get Seller By Seller ID HTTPError",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetSellerByUserID", mock.Anything, mock.Anything).Return(&body.SellerResponse{}, httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
		{
			name:  "Error Get Seller By Seller ID Error",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetSellerByUserID", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/seller/user/%s", tc.param), nil)
			r.Header = make(http.Header)

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "user_id",
						Value: tc.param,
					},
				}
			}

			uuid.Parse(tc.param)

			c.Set("sellerID", tc.param)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetSellerByUserID(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_GetSellerDetailInformation(t *testing.T) {
	testCase := []struct {
		name       string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
		userID     string
	}{
		{
			name: "Success Get Seller Detail Information",
			mock: func(s *mocks.UseCase) {
				s.On("GetSellerByUserID", mock.Anything, mock.Anything).Return(&body.SellerResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Unauthorized Get Seller Detail Information",
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusUnauthorized,
			authorized: false,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Error Get Seller Detail Information",
			mock: func(s *mocks.UseCase) {
				s.On("GetSellerByUserID", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Error Get Seller Detail Information Error",
			mock: func(s *mocks.UseCase) {
				s.On("GetSellerByUserID", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/seller/information", nil)
			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", tc.userID)
			}
			uuid.Parse(tc.userID)

			c.Set("userIDString", tc.userID)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetSellerDetailInformation(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_UpdateSellerInformation(t *testing.T) {
	invalidRequestBody := struct {
		ShopName int `json:"shop_name"`
	}{123}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
		userID     string
	}{
		{
			name: "Success Change Seller Information",
			body: body.UpdateSellerInformationRequest{
				ShopName: "test",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateSellerInformationByUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:       "body is nil",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:       "Invalid request",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Unauthorized User",
			body: body.UpdateSellerInformationRequest{
				ShopName: "test",
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
			userID:     "",
		},
		{
			name: "Error Change Seller Information",
			body: body.UpdateSellerInformationRequest{
				ShopName: "test",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateSellerInformationByUserID", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Error Change Seller Information HTTPError",
			body: body.UpdateSellerInformationRequest{
				ShopName: "test",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateSellerInformationByUserID", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
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

			r := httptest.NewRequest(http.MethodPatch, "/api/v1/seller/order-cancel", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)
			if tc.authorized {
				c.Set("userID", tc.userID)
			}
			uuid.Parse(tc.userID)

			c.Set("userIDString", tc.userID)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.UpdateSellerInformation(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_GetCategoryBySellerID(t *testing.T) {
	var responseCategory []*body.CategoryResponse

	testCase := []struct {
		name     string
		param    string
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name:  "Success Get category by seller ID",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetCategoryBySellerID", mock.Anything, mock.Anything).Return(responseCategory, nil)
			},
			expected: http.StatusOK,
		},
		{
			name:  "Error Get category by seller ID",
			param: "1",
			mock: func(s *mocks.UseCase) {
			},
			expected: http.StatusBadRequest,
		},
		{
			name:  "Error Get category by seller ID HTTPError",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetCategoryBySellerID", mock.Anything, mock.Anything).Return(responseCategory, httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
		{
			name:  "Error Get category by seller ID Error",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetCategoryBySellerID", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/seller/%s/category", tc.param), nil)
			r.Header = make(http.Header)

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "seller_id",
						Value: tc.param,
					},
				}
			}

			uuid.Parse(tc.param)

			c.Set("sellerID", tc.param)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetCategoryBySellerID(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_GetCourierSeller(t *testing.T) {
	var response *body.CourierSellerResponse

	testCase := []struct {
		name       string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
		userID     string
	}{
		{
			name: "Success Get Courier by seller ID",
			mock: func(s *mocks.UseCase) {
				s.On("GetCourierSeller", mock.Anything, mock.Anything).Return(response, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Unauthorized Get Courier by seller ID",
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusUnauthorized,
			authorized: false,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Error Get Courier by seller ID HTTPError",
			mock: func(s *mocks.UseCase) {
				s.On("GetCourierSeller", mock.Anything, mock.Anything).Return(response, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Error Get Courier by seller ID Error",
			mock: func(s *mocks.UseCase) {
				s.On("GetCourierSeller", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/seller/courier", nil)
			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", tc.userID)
			}
			uuid.Parse(tc.userID)

			c.Set("userIDString", tc.userID)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetCourierSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_CreateCourierSeller(t *testing.T) {
	invalidRequestBody := struct {
		CourierID int `json:"courier_id"`
	}{123}
	RequestBody := body.CourierSellerRequest{
		CourierID: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
	}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
		userID     string
	}{
		{
			name: "Success Create Courier Seller",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("CreateCourierSeller", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:       "body is nil",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:       "Invalid request",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:       "Unauthorized User",
			body:       RequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
			userID:     "",
		},
		{
			name: "Error Create Courier Seller",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("CreateCourierSeller", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Error Create Courier Seller HTTPError",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("CreateCourierSeller", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
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

			r := httptest.NewRequest(http.MethodPatch, "/api/v1/seller/order-status", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)
			if tc.authorized {
				c.Set("userID", tc.userID)
			}
			uuid.Parse(tc.userID)

			c.Set("userIDString", tc.userID)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CreateCourierSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_DeleteCourierSellerByID(t *testing.T) {
	testCase := []struct {
		name     string
		param    string
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name:  "Success Get Order By Order ID",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteCourierSellerByID", mock.Anything, mock.Anything).Return(nil)
			},
			expected: http.StatusOK,
		},
		{
			name:  "Error Get Order By Order ID",
			param: "1",
			mock: func(s *mocks.UseCase) {
			},
			expected: http.StatusBadRequest,
		},
		{
			name:  "Error Get Order By Order ID HTTPError",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteCourierSellerByID", mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
		{
			name:  "Error Get Order By Order ID Error",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteCourierSellerByID", mock.Anything, mock.Anything).Return(errors.New("error"))
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/seller/courier/%s", tc.param), nil)
			r.Header = make(http.Header)

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "id",
						Value: tc.param,
					},
				}
			}

			uuid.Parse(tc.param)

			c.Set("sellerCourierID", tc.param)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.DeleteCourierSellerByID(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_UpdateResiNumberInOrderSeller(t *testing.T) {
	invalidRequestBody := struct {
		EstimateArriveAt string `json:"estimate_arrive_at"`
	}{"02-01-2006"}
	RequestBody := body.UpdateNoResiOrderSellerRequest{
		NoResi:               "123456789",
		EstimateArriveAt:     "02-01-2006 15:04:05",
		EstimateArriveAtTime: time.Now(),
	}
	param := "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4"

	testCase := []struct {
		name       string
		body       interface{}
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
		userID     string
	}{
		{
			name:  "Success Update Resi Number In Order Seller",
			body:  RequestBody,
			param: param,
			mock: func(s *mocks.UseCase) {
				s.On("UpdateResiNumberInOrderSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:  "Unauthorized Update Resi Number In Order Seller",
			body:  RequestBody,
			param: param,
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusUnauthorized,
			authorized: false,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:  "id Invalid",
			body:  RequestBody,
			param: "1",
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:       "Request Body Empty Update Resi Number In Order Seller",
			body:       nil,
			param:      param,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:       "Invalid Request Body Update Resi Number In Order Seller",
			body:       invalidRequestBody,
			param:      param,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:  "Internal Server Error Update Resi Number In Order Seller",
			body:  RequestBody,
			param: param,
			mock: func(s *mocks.UseCase) {
				s.On("UpdateResiNumberInOrderSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Internal Server Error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name:  "Internal Server Error Update Resi Number In Order Seller HTTPError",
			body:  RequestBody,
			param: param,
			mock: func(s *mocks.UseCase) {
				s.On("UpdateResiNumberInOrderSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
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

			r := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/seller/order-resi/%s", tc.param), bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)
			if tc.authorized {
				c.Set("userID", tc.userID)
			}
			uuid.Parse(tc.userID)

			c.Set("userIDString", tc.userID)

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "id",
						Value: tc.param,
					},
				}
			}

			uuid.Parse(tc.param)

			c.Set("orderID", tc.param)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.UpdateResiNumberInOrderSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_WithdrawalOrderBalance(t *testing.T) {
	param := "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4"

	testCase := []struct {
		name     string
		param    string
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name:  "Success Update Resi Number In Order Seller",
			param: param,
			mock: func(s *mocks.UseCase) {
				s.On("WithdrawalOrderBalance", mock.Anything, mock.Anything).Return(nil)
			},
			expected: http.StatusOK,
		},
		{
			name:     "Invalid id",
			param:    "1",
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
		},
		{
			name:  "Internal Server Error Update Resi Number In Order Seller",
			param: param,
			mock: func(s *mocks.UseCase) {
				s.On("WithdrawalOrderBalance", mock.Anything, mock.Anything).Return(errors.New("Internal Server Error"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name:  "Internal Server Error Update Resi Number In Order Seller HTTPError",
			param: param,
			mock: func(s *mocks.UseCase) {
				s.On("WithdrawalOrderBalance", mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/seller/withdrawal/%s", tc.param), nil)
			r.Header = make(http.Header)

			c.Request = r

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "id",
						Value: tc.param,
					},
				}
			}

			uuid.Parse(tc.param)

			c.Set("orderID", tc.param)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.WithdrawalOrderBalance(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_GetAllVoucherSeller(t *testing.T) {
	var response *pagination.Pagination

	testCase := []struct {
		name       string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
		userID     string
	}{
		{
			name: "Success Get Voucher by seller ID",
			mock: func(s *mocks.UseCase) {
				s.On("GetAllVoucherSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(response, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Unauthorized Get Voucher by seller ID",
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusUnauthorized,
			authorized: false,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Error Get Voucher by seller ID HTTPError",
			mock: func(s *mocks.UseCase) {
				s.On("GetAllVoucherSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(response, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
		{
			name: "Error Get Voucher by seller ID Error",
			mock: func(s *mocks.UseCase) {
				s.On("GetAllVoucherSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
			userID:     "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/seller/voucher", nil)
			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", tc.userID)
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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetAllVoucherSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_CreateVoucherSeller(t *testing.T) {
	invalidRequestBody := struct {
		Code string `json:"code"`
	}{""}
	RequestBody := body.CreateVoucherRequest{
		Code:               "test",
		Quota:              1,
		ActivedDate:        "02-01-2006 15:04:05",
		ExpiredDate:        "02-01-2006 15:04:05",
		DiscountPercentage: 1,
		DiscountFixPrice:   1,
		MinProductPrice:    1,
		MaxDiscountPrice:   1,
		ActiveDateTime:     time.Now(),
		ExpiredDateTime:    time.Now(),
	}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Create Voucher",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("CreateVoucherSeller", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized Create Voucher",
			body:       RequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Body Nil",
			body:       "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "Invalid Body",
			body:       invalidRequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Error Create Voucher HTTPError",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("CreateVoucherSeller", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Error Create Voucher Error",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("CreateVoucherSeller", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/seller/voucher", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)
			if tc.authorized {
				c.Set("userID", "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4")
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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CreateVoucherSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_DeleteVoucherSeller(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Delete Voucher",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteVoucherSeller", mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized Delete Voucher",
			param:      "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Param Nil",
			param:      "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Error Delete Voucher HTTPError",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteVoucherSeller", mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Error Delete Voucher Error",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("DeleteVoucherSeller", mock.Anything, mock.Anything).Return(errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/seller/voucher/%s", tc.param), nil)
			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4")
			}

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "id",
						Value: tc.param,
					},
				}
			}

			uuid.Parse(tc.param)

			c.Set("voucherShopID", tc.param)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.DeleteVoucherSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_UpdateVoucherSeller(t *testing.T) {
	// InvalidRequestBody := body.UpdateVoucherRequest{
	// 	VoucherID:          "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
	// 	Quota:              10,
	// 	ActivedDate:        "02-01-2006 15:04:05",
	// 	ExpiredDate:        "2021-01-02 00:00:00",
	// 	DiscountPercentage: 10,
	// 	DiscountFixPrice:   10000,
	// 	MinProductPrice:    10000,
	// 	MaxDiscountPrice:   10000,
	// }
	RequestBody := body.UpdateVoucherRequest{
		VoucherID:          "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		Quota:              10,
		ActivedDate:        "02-01-2006 15:04:05",
		ExpiredDate:        "03-01-2006 15:04:05",
		DiscountPercentage: 10,
		DiscountFixPrice:   10000,
		MinProductPrice:    10000,
		MaxDiscountPrice:   10000,
	}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Update Voucher",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("UpdateVoucherSeller", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "Failed Update Voucher",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("UpdateVoucherSeller", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Failed Update Voucher HTTP ERROR",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("UpdateVoucherSeller", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
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

			r := httptest.NewRequest(http.MethodPut, "/api/v1/seller/voucher", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)
			if tc.authorized {
				c.Set("userID", "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4")
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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.UpdateVoucherSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_UpdateOnDeliveryOrder(t *testing.T) {
	testCase := []struct {
		name     string
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name: "Success Update On Delivery Order",
			mock: func(s *mocks.UseCase) {
				s.On("UpdateOnDeliveryOrder", mock.Anything).Return(nil)
			},
			expected: http.StatusOK,
		},
		{
			name: "Failed Update On Delivery Order",
			mock: func(s *mocks.UseCase) {
				s.On("UpdateOnDeliveryOrder", mock.Anything).Return(errors.New("error"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name: "Failed Update On Delivery Order HTTP ERROR",
			mock: func(s *mocks.UseCase) {
				s.On("UpdateOnDeliveryOrder", mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPost, "/api/v1/seller/delivery", nil)
			r.Header = make(http.Header)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.UpdateOnDeliveryOrder(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_UpdateExpiredAtOrder(t *testing.T) {
	testCase := []struct {
		name     string
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name: "Success Update On Delivery Order",
			mock: func(s *mocks.UseCase) {
				s.On("UpdateExpiredAtOrder", mock.Anything).Return(nil)
			},
			expected: http.StatusOK,
		},
		{
			name: "Failed Update On Delivery Order",
			mock: func(s *mocks.UseCase) {
				s.On("UpdateExpiredAtOrder", mock.Anything).Return(errors.New("error"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name: "Failed Update On Delivery Order HTTP ERROR",
			mock: func(s *mocks.UseCase) {
				s.On("UpdateExpiredAtOrder", mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodPost, "/api/v1/seller/expired", nil)
			r.Header = make(http.Header)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.UpdateExpiredAtOrder(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_DetailVoucherSeller(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Get Detail Voucher",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetDetailVoucherSeller", mock.Anything, mock.Anything).Return(&model.Voucher{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized Get Detail Voucher",
			param:      "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Param Nil",
			param:      "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Error Get Detail Voucher HTTPError",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetDetailVoucherSeller", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Error Get Detail Voucher Error",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetDetailVoucherSeller", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/seller/voucher/%s", tc.param), nil)
			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4")
			}

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "id",
						Value: tc.param,
					},
				}
			}

			uuid.Parse(tc.param)

			c.Set("voucherShopID", tc.param)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.DetailVoucherSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_GetAllPromotionSeller(t *testing.T) {
	response := &pagination.Pagination{}

	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Get Promotion",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetAllPromotionSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(response, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized Get Promotion",
			param:      "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:  "Error Get Promotion HTTPError",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetAllPromotionSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Error Get Promotion Error",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetAllPromotionSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/seller/voucher/%s", tc.param), nil)
			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4")
			}

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "id",
						Value: tc.param,
					},
				}
			}

			uuid.Parse(tc.param)

			c.Set("voucherShopID", tc.param)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetAllPromotionSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_CreatePromotionSeller(t *testing.T) {
	InvalidRequestBody := body.CreatePromotionRequest{
		Name:             "tes t",
		ProductPromotion: []body.ProductPromotionData{},
		ActivedDate:      "02-01-2006 15:04:05",
		ExpiredDate:      "2021-01-02 00:00:00",
	}
	RequestBody := body.CreatePromotionRequest{
		Name: "test",
		ProductPromotion: []body.ProductPromotionData{
			{
				ProductID:          "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
				Quota:              10,
				MaxQuantity:        10,
				DiscountPercentage: 10,
				DiscountFixPrice:   10000,
				MinProductPrice:    10000,
				MaxDiscountPrice:   10000,
			},
		},
		ActivedDate: "02-01-2006 15:04:05",
		ExpiredDate: "02-01-2006 00:00:00",
	}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Create Promotion",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("CreatePromotionSeller", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized Create Promotion",
			body:       RequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Body Empty Create Promotion",
			body:       "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Invalid Body Create Promotion",
			body: InvalidRequestBody,
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Failed Create Promotion",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("CreatePromotionSeller", mock.Anything, mock.Anything, mock.Anything).Return(1, errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Failed Create Promotion HTTP ERROR",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("CreatePromotionSeller", mock.Anything, mock.Anything, mock.Anything).Return(1, httperror.New(http.StatusBadRequest, "test"))
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/seller/promotion", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)
			if tc.authorized {
				c.Set("userID", "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4")
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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CreatePromotionSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_UpdatePromotionSeller(t *testing.T) {
	InvalidRequestBody := body.UpdatePromotionRequest{
		PromotionID: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
	}
	RequestBody := body.UpdatePromotionRequest{
		PromotionID:        "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		ProductID:          "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		PromotionName:      "test",
		MaxQuantity:        1,
		ActivedDate:        "02-01-2006 15:04:05",
		ExpiredDate:        "02-01-2006 15:04:05",
		DiscountPercentage: 10,
		DiscountFixPrice:   10000,
		MinProductPrice:    10000,
		MaxDiscountPrice:   10000,
	}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Update Promotion",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("UpdatePromotionSeller", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized Update Promotion",
			body:       RequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Body Empty Update Promotion",
			body:       "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Invalid Body Update Promotion",
			body: InvalidRequestBody,
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Failed Update Promotion",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("UpdatePromotionSeller", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Failed Update Promotion HTTP ERROR",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("UpdatePromotionSeller", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
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

			r := httptest.NewRequest(http.MethodPut, "/api/v1/seller/promotion", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)
			if tc.authorized {
				c.Set("userID", "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4")
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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.UpdatePromotionSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_GetDetailPromotionSellerByID(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Get Detail Promotion By Detail Promotion ID",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetDetailPromotionSellerByID", mock.Anything, mock.Anything).Return(&body.PromotionDetailSeller{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized Get Detail Promotion By Detail Promotion ID",
			param:      "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:  "Error Get Detail Promotion By Detail Promotion ID",
			param: "1",
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Error Get Detail Promotion By Detail Promotion ID HTTPError",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetDetailPromotionSellerByID", mock.Anything, mock.Anything).Return(&body.PromotionDetailSeller{}, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Error Get Detail Promotion By Detail Promotion ID Error",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetDetailPromotionSellerByID", mock.Anything, mock.Anything).Return(&body.PromotionDetailSeller{}, errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/seller/promotion/%s", tc.param), nil)
			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4")
			}

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "id",
						Value: tc.param,
					},
				}
			}

			uuid.Parse(tc.param)

			c.Set("promotionShopID", tc.param)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetDetailPromotionSellerByID(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_GetProductWithoutPromotionSeller(t *testing.T) {
	testCase := []struct {
		name       string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Get Product Without Promotion",
			mock: func(s *mocks.UseCase) {
				s.On("GetProductWithoutPromotionSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized Get Product Without Promotion",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "Error Get Product Without Promotion HTTPError",
			mock: func(s *mocks.UseCase) {
				s.On("GetProductWithoutPromotionSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Error Get Product Without Promotion Error",
			mock: func(s *mocks.UseCase) {
				s.On("GetProductWithoutPromotionSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/seller/product/without-promotion", nil)
			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4")
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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetProductWithoutPromotionSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_CreateRefundThreadSeller(t *testing.T) {
	InvalidRequestBody := body.CreateRefundThreadRequest{
		RefundID: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
	}
	RequestBody := body.CreateRefundThreadRequest{
		RefundID: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
		IsSeller: true,
		IsBuyer:  true,
		Text:     "test",
	}

	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "Success Create Refund Thread",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("CreateRefundThreadSeller", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized Create Refund Thread",
			body:       RequestBody,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:       "Body Empty Create Refund Thread",
			body:       "",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "Invalid Body Create Refund Thread",
			body: InvalidRequestBody,
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "Failed Create Refund Thread",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("CreateRefundThreadSeller", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "Failed Create Refund Thread HTTP ERROR",
			body: RequestBody,
			mock: func(s *mocks.UseCase) {
				s.On("CreateRefundThreadSeller", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/seller/refund-thread", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPost(c, tc.body)
			if tc.authorized {
				c.Set("userID", "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4")
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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CreateRefundThreadSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func Test_sellerHandlers_GetRefundOrderSeller(t *testing.T) {
	testCase := []struct {
		name       string
		param      string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name:  "Success Get Refund Order",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetRefundOrderSeller", mock.Anything, mock.Anything, mock.Anything).Return(&body.GetRefundThreadResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "Unauthorized Get Refund Order",
			param:      "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name:  "Error Get Refund Order",
			param: "1",
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Error Get Refund Order HTTPError",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetRefundOrderSeller", mock.Anything, mock.Anything, mock.Anything).Return(&body.GetRefundThreadResponse{}, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:  "Error Get Refund Order Error",
			param: "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4",
			mock: func(s *mocks.UseCase) {
				s.On("GetRefundOrderSeller", mock.Anything, mock.Anything, mock.Anything).Return(&body.GetRefundThreadResponse{}, errors.New("error"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/seller/refund/%s", tc.param), nil)
			r.Header = make(http.Header)

			if tc.authorized {
				c.Set("userID", "4cf3a332-5d81-48a0-b935-cfa83a6b6ac4")
			}

			if tc.param != "" {
				c.Params = []gin.Param{
					{
						Key:   "refund_id",
						Value: tc.param,
					},
				}
			}

			uuid.Parse(tc.param)

			c.Set("refundID", tc.param)

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

			h := NewSellerHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetRefundOrderSeller(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}
