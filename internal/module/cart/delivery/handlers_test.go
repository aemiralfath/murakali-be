package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"murakali/config"
	"murakali/internal/module/cart/delivery/body"
	"murakali/internal/module/cart/mocks"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"murakali/pkg/pagination"
	"net/http"
	"net/http/httptest"
	"testing"

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

func MockJsonPUT(c *gin.Context, content interface{}) {
	c.Request.Method = "PUT"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
}

func TestCartHandlers_GetCartHoverHome(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success get cart hover",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCartHoverHome", mock.Anything, mock.Anything, mock.Anything).Return(&body.CartHomeResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "get cart hover unauthorized",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "get cart hover error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCartHoverHome", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "get cart hover error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCartHoverHome", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/cart/hover-home", nil)
			r.Header = make(http.Header)
			c.Request = r

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

			h := NewCartHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetCartHoverHome(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestCartHandlers_GetCartItems(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success get cart items",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCartItems", mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "get cart items unauthorized",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "get cart error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCartItems", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "get cart error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCartItems", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/cart/items", nil)
			r.Header = make(http.Header)
			c.Request = r

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

			h := NewCartHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetCartItems(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestCartHandlers_AddCartItems(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success update cart",
			body: body.AddCartItemRequest{
				ProductDetailID: "123456",
				Quantity:        12,
			},
			mock: func(s *mocks.UseCase) {
				s.On("AddCartItems", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "add cart unauthorized",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "invalid request",
			body: body.AddCartItemRequest{
				ProductDetailID: "123456",
				Quantity:        12,
			},
			mock: func(s *mocks.UseCase) {
				s.On("AddCartItems", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "invalid request entity",
			body: body.AddCartItemRequest{
				ProductDetailID: "",
				Quantity:        0,
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "add cart error internal",
			body: body.AddCartItemRequest{
				ProductDetailID: "123456",
				Quantity:        12,
			},
			mock: func(s *mocks.UseCase) {
				s.On("AddCartItems", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/cart/items", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			if tc.authorized {
				c.Set("userID", "123456")
			}

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

			h := NewCartHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.AddCartItems(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestCartHandlers_UpdateCartItems(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success update cart",
			body: body.CartItemRequest{
				ProductDetailID: "123456",
				Quantity:        12,
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateCartItems", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "update cart unauthorized",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "invalid request",
			body: body.CartItemRequest{
				ProductDetailID: "123456",
				Quantity:        12,
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateCartItems", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "invalid request entity",
			body: body.CartItemRequest{
				ProductDetailID: "",
				Quantity:        0,
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "update cart error internal",
			body: body.CartItemRequest{
				ProductDetailID: "123456",
				Quantity:        12,
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateCartItems", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
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

			r := httptest.NewRequest(http.MethodPut, "/api/v1/cart/items", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			if tc.authorized {
				c.Set("userID", "123456")
			}

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

			h := NewCartHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.UpdateCartItems(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestCartHandlers_DeleteCartItems(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success delete cart ",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("DeleteCartItems", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "get delete cart unauthorized",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "delete cart  error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("DeleteCartItems", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "delete cart  error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("DeleteCartItems", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodDelete, "/api/v1/cart/:id", nil)
			r.Header = make(http.Header)
			c.Request = r
			c.Params = []gin.Param{
				{
					Key:   "id",
					Value: "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
				},
			}

			s := mocks.NewUseCase(t)

			if tc.authorized {
				c.Set("userID", "123456")
			}

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

			h := NewCartHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.DeleteCartItems(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestCartHandlers_GetVoucherShop(t *testing.T) {
	testCase := []struct {
		name     string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
		parseID  bool
	}{
		{
			name: "success get voucher ",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetVoucherShop", mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected: http.StatusOK,
			parseID:  true,
		},
		{
			name:     "error parse id ",
			body:     nil,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
			parseID:  false,
		},
		{
			name: "get voucher  error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetVoucherShop", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected: http.StatusInternalServerError,
			parseID:  true,
		},
		{
			name: "get voucher  error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetVoucherShop", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
			parseID:  true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/cart/voucher/:shop_id", nil)
			r.Header = make(http.Header)
			c.Request = r

			if tc.parseID {
				c.Params = []gin.Param{
					{
						Key:   "shop_id",
						Value: "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
					},
				}
			} else {
				c.Params = []gin.Param{
					{
						Key:   "shop_id",
						Value: "123456",
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

			h := NewCartHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetVoucherShop(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}
func TestCartHandlers_GetVoucherMarketplace(t *testing.T) {
	testCase := []struct {
		name     string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name: "success get voucher ",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetVoucherMarketplace", mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected: http.StatusOK,
		},
		{
			name: "get voucher  error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetVoucherMarketplace", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name: "get voucher  error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetVoucherMarketplace", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/cart/voucher", nil)
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

			h := NewCartHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetVoucherMarketplace(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}
