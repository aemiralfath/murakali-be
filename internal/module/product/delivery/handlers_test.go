package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/product/delivery/body"
	"murakali/internal/module/product/mocks"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"murakali/pkg/pagination"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestProductHandlers_GetCategories(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success get categories ",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCategories", mock.Anything).Return([]*body.CategoryResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "get categories  error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCategories", mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "get categories  error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCategories", mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/product/category", nil)
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetCategories(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestProductHandlers_GetBanners(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success get banners ",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetBanners", mock.Anything).Return([]*model.Banner{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "get banners  error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetBanners", mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "get banners  error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetBanners", mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/product/banner", nil)
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetBanners(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestProductHandlers_GetCategoriesByNameLevelOne(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success get categories level one ",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCategoriesByName", mock.Anything, mock.Anything).Return([]*body.CategoryResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "get categories level one  error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCategoriesByName", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "get categories level one  error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCategoriesByName", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/product/category/:name_lvl_one", nil)
			r.Header = make(http.Header)
			c.Request = r

			s := mocks.NewUseCase(t)

			c.Params = []gin.Param{
				{
					Key:   "name_lvl_one",
					Value: "Hobby & Colection",
				},
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetCategoriesByNameLevelOne(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestProductHandlers_GetCategoriesByNameLevelTwo(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success get categories level two ",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCategoriesByName", mock.Anything, mock.Anything).Return([]*body.CategoryResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "get categories level two  error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCategoriesByName", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "get categories level two  error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCategoriesByName", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/product/category/:name_lvl_two", nil)
			r.Header = make(http.Header)
			c.Request = r

			s := mocks.NewUseCase(t)

			c.Params = []gin.Param{
				{
					Key:   "name_lvl_two",
					Value: "Hobby & Colection",
				},
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetCategoriesByNameLevelTwo(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestProductHandlers_GetCategoriesByNameLevelThree(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success get categories level three ",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCategoriesByName", mock.Anything, mock.Anything).Return([]*body.CategoryResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "get categories level three  error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCategoriesByName", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "get categories level three  error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCategoriesByName", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/product/category/:name_lvl_two", nil)
			r.Header = make(http.Header)
			c.Request = r

			s := mocks.NewUseCase(t)

			c.Params = []gin.Param{
				{
					Key:   "name_lvl_three",
					Value: "Hobby & Colection",
				},
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetCategoriesByNameLevelThree(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestProductHandlers_GetRecommendedProducts(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success get recommended products ",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetRecommendedProducts", mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "get recommended products  error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetRecommendedProducts", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "get recommended products  error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetRecommendedProducts", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/product/recommended", nil)
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetRecommendedProducts(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestProductHandlers_GetProducts(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		queries    map[string]string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success get  products query 1",
			body: nil,

			mock: func(s *mocks.UseCase) {
				s.On("GetProducts", mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "success get  products query 2",
			body: nil,
			queries: map[string]string{
				"limit":         "101",
				"listed_status": "1",
				"sort":          "asc",
				"sort_by":       "created_at",
				"min_price":     "0",
				"max_price":     "10000",
				"min_rating":    "8",
				"max_rating":    "9",
				"province_ids":  "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetProducts", mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "success get  products query 3",
			body: nil,
			queries: map[string]string{
				"limit":         "101",
				"listed_status": "2",
				"sort":          "asc",
				"sort_by":       "recommended",
				"min_price":     "0",
				"max_price":     "10000",
				"min_rating":    "8",
				"max_rating":    "9",
				"province_ids":  "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetProducts", mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "success get  products query 4",
			body: nil,
			queries: map[string]string{
				"limit":         "101",
				"listed_status": "3",
				"sort":          "asc",
				"sort_by":       "min_price",
				"min_price":     "0",
				"max_price":     "10000",
				"min_rating":    "1",
				"max_rating":    "4",
				"province_ids":  "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetProducts", mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "success get  products query 5",
			body: nil,
			queries: map[string]string{
				"limit":         "101",
				"listed_status": "3",
				"sort":          "asc",
				"sort_by":       "unit_sold",
				"min_price":     "0",
				"max_price":     "10000",
				"min_rating":    "1",
				"max_rating":    "4",
				"province_ids":  "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetProducts", mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "success get  products query 6",
			body: nil,
			queries: map[string]string{
				"limit":         "101",
				"listed_status": "3",
				"sort":          "asc",
				"sort_by":       "view_count",
				"min_price":     "0",
				"max_price":     "10000",
				"min_rating":    "1",
				"max_rating":    "4",
				"province_ids":  "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetProducts", mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "success get  products query 6",
			body: nil,
			queries: map[string]string{
				"limit":         "101",
				"listed_status": "3",
				"sort":          "asc",
				"sort_by":       "listed_status",
				"min_price":     "0",
				"max_price":     "10000",
				"min_rating":    "1",
				"max_rating":    "4",
				"province_ids":  "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetProducts", mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "get  products  error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetProducts", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "get  products  error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetProducts", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/product/", nil)
			r.Header = make(http.Header)
			c.Request = r

			s := mocks.NewUseCase(t)
			if tc.queries != nil && len(tc.queries) > 0 {
				u := url.Values{}
				for key, value := range tc.queries {
					u.Set(key, value)
				}
				c.Request.URL.RawQuery = u.Encode()
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetProducts(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestProductHandlers_GetFavoriteProducts(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success get favorites products ",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetFavoriteProducts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name:       "get favorites productsunauthorized",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "get favorites products  error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetFavoriteProducts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "get favorites products  error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetFavoriteProducts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/product/favorite", nil)
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetFavoriteProducts(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestProductHandlers_GetProductDetail(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success get product detail",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetProductDetail", mock.Anything, mock.Anything).Return(&body.ProductDetailResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "get product detail error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetProductDetail", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "get product detail error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetProductDetail", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/product/:product_id", nil)
			r.Header = make(http.Header)
			c.Request = r

			s := mocks.NewUseCase(t)

			c.Params = []gin.Param{
				{
					Key:   "product_id",
					Value: "123456",
				},
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetProductDetail(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestProductHandlers_GetAllProductImage(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success get product image",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetAllProductImage", mock.Anything, mock.Anything).Return([]*body.GetImageResponse{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "get product image error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetAllProductImage", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "get product image error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetAllProductImage", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/product/:product_id/picture", nil)
			r.Header = make(http.Header)
			c.Request = r

			s := mocks.NewUseCase(t)

			c.Params = []gin.Param{
				{
					Key:   "product_id",
					Value: "123456",
				},
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetAllProductImage(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}
