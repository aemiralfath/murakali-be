package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func MockJsonPATCH(c *gin.Context, content interface{}) {
	c.Request.Method = "PATCH"
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
			name:       "get favorites products unauthorized",
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

func TestCartHandlers_CheckProductIsFavorite(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success check is favorite product",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CheckProductIsFavorite", mock.Anything, mock.Anything, mock.Anything).Return(true)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "check unauthorized",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},

		{
			name: "invalid request",
			body: body.GetProductRequest{
				ProductID: "   ",
			},
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusUnprocessableEntity,
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/product/favorite", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)
			if tc.authorized {
				c.Set("userID", "123456")
			}

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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CheckProductIsFavorite(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestCartHandlers_CountSpecificFavoriteProduct(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success count favorite product",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CountSpecificFavoriteProduct", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "count favorite error custom",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CountSpecificFavoriteProduct", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "invalid request",
			body: body.GetProductRequest{
				ProductID: "   ",
			},
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "count favorite product error internal",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CountSpecificFavoriteProduct", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), fmt.Errorf("test"))
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/product/favorite", bytes.NewBuffer(jsonValue))
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CountSpecificFavoriteProduct(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestCartHandlers_CreateFavoriteProduct(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success create favorite product",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateFavoriteProduct", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "create unauthorized",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "invalid request",
			body: body.GetProductRequest{
				ProductID: "   ",
			},
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},

		{
			name: "add create favorite product error internal",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateFavoriteProduct", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/product/favorite", bytes.NewBuffer(jsonValue))
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CreateFavoriteProduct(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestCartHandlers_DeleteFavoriteProduct(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success delete favorite product",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("DeleteFavoriteProduct", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "success delete unauthorized",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "invalid request",
			body: body.GetProductRequest{
				ProductID: "   ",
			},
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},

		{
			name: "add delete favorite product error internal",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("DeleteFavoriteProduct", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
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

			r := httptest.NewRequest(http.MethodDelete, "/api/v1/product/favorite", bytes.NewBuffer(jsonValue))
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.DeleteFavoriteProduct(c)

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

func TestCartHandlers_UpdateProductMetadata(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success update product meta data",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateProductMetadata", mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "update product error custom",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateProductMetadata", mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "update product meta data error internal",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateProductMetadata", mock.Anything).Return(errors.New("test"))
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

			r := httptest.NewRequest(http.MethodDelete, "/api/v1/product/favorite", bytes.NewBuffer(jsonValue))
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.UpdateProductMetadata(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestCartHandlers_CreateProduct(t *testing.T) {
	var temp float64 = 10
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success create  product",
			body: body.CreateProductRequest{
				ProductInfo: body.CreateProductInfo{
					Title:        "test",
					Description:  "description",
					Thumbnail:    "test",
					CategoryID:   "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
					ListedStatus: true,
				},
				ProductDetail: []body.CreateProductDetailRequest{{
					Price:     temp,
					Stock:     temp,
					Weight:    temp,
					Size:      temp,
					Hazardous: true,
					Codition:  "test",
					BulkPrice: true,
					Photo:     []string{"test", "test"},
					VariantDetail: []body.VariantDetailRequest{{
						Type: "color",
						Name: "big",
					}},
				}},
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "create product  error custom",
			body: body.CreateProductRequest{
				ProductInfo: body.CreateProductInfo{
					Title:        "test",
					Description:  "description",
					Thumbnail:    "test",
					CategoryID:   "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
					ListedStatus: true,
				},
				ProductDetail: []body.CreateProductDetailRequest{{
					Price:     temp,
					Stock:     temp,
					Weight:    temp,
					Size:      temp,
					Hazardous: true,
					Codition:  "test",
					BulkPrice: true,
					Photo:     []string{"test", "test"},
					VariantDetail: []body.VariantDetailRequest{{
						Type: "color",
						Name: "big",
					}},
				}},
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "create product unauthorized",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "invalid request",
			body: nil,
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},

		{
			name: "add create  product error internal",
			body: body.CreateProductRequest{
				ProductInfo: body.CreateProductInfo{
					Title:        "test",
					Description:  "description",
					Thumbnail:    "test",
					CategoryID:   "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
					ListedStatus: true,
				},
				ProductDetail: []body.CreateProductDetailRequest{{
					Price:     temp,
					Stock:     temp,
					Weight:    temp,
					Size:      temp,
					Hazardous: true,
					Codition:  "test",
					BulkPrice: true,
					Photo:     []string{"test", "test"},
					VariantDetail: []body.VariantDetailRequest{{
						Type: "color",
						Name: "big",
					}},
				}},
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/product", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			if tc.authorized {
				c.Set("userID", "989d94b7-58fc-4a76-ae01-1c1b47a0755c")
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CreateProduct(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestCartHandlers_UpdateListedStatus(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success update listed status product",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("UpdateListedStatus", mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "update listed status error custom",
			body: body.GetProductRequest{
				ProductID: "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateListedStatus", mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "add update listed status product error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("UpdateListedStatus", mock.Anything, mock.Anything).Return(errors.New("test"))
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

			r := httptest.NewRequest(http.MethodPut, "/api/v1/product/status/:id", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			c.Params = []gin.Param{
				{
					Key:   "id",
					Value: "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
				},
			}

			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPUT(c, tc.body)

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
			h.UpdateListedStatus(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestCartHandlers_UpdateListedStatusBulk(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success update listed status bulk product",
			body: body.UpdateProductListedStatusBulkRequest{
				ProductIDS:   []string{"989d94b7-58fc-4a76-ae01-1c1b47a0755c"},
				ListedStatus: true,
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateProductListedStatusBulk", mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "update listed status bulk error input",
			body: body.UpdateProductListedStatusBulkRequest{
				ProductIDS:   []string{"1234"},
				ListedStatus: true,
			},
			mock: func(s *mocks.UseCase) {

			},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},
		{
			name: "update listed status bulk error custom",
			body: body.UpdateProductListedStatusBulkRequest{
				ProductIDS:   []string{"989d94b7-58fc-4a76-ae01-1c1b47a0755c"},
				ListedStatus: true,
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateProductListedStatusBulk", mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "add update listed status bulk product error internal",
			body: body.UpdateProductListedStatusBulkRequest{
				ProductIDS:   []string{"989d94b7-58fc-4a76-ae01-1c1b47a0755c"},
				ListedStatus: true,
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateProductListedStatusBulk", mock.Anything, mock.Anything).Return(errors.New("test"))
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

			r := httptest.NewRequest(http.MethodPatch, "/api/v1/product/bulk-status", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r

			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPATCH(c, tc.body)

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
			h.UpdateListedStatusBulk(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestCartHandlers_UpdateProduct(t *testing.T) {
	var temp float64 = 10
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success update  product",
			body: body.UpdateProductRequest{
				ProductInfo: body.UpdateProductInfo{
					Title:        "test",
					Description:  "description",
					Thumbnail:    "test",
					ListedStatus: true,
				},
				ProductDetail: []body.UpdateProductDetailRequest{{
					ProductDetailID: "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
					Price:           temp,
					Stock:           temp,
					Weight:          temp,
					Size:            temp,
					Hazardous:       true,
					Codition:        "test",
					BulkPrice:       true,
					Photo:           []string{"test", "test"},
					VariantDetailID: []body.UpdateVariant{{
						VariantID:       "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
						VariantDetailID: "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
					}},
					VariantIDRemove: []string{"989d94b7-58fc-4a76-ae01-1c1b47a0755c"},
				}},
				ProductDetailRemove: []string{"989d94b7-58fc-4a76-ae01-1c1b47a0755c"},
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "update product  error custom",
			body: body.UpdateProductRequest{
				ProductInfo: body.UpdateProductInfo{
					Title:        "test",
					Description:  "description",
					Thumbnail:    "test",
					ListedStatus: true,
				},
				ProductDetail: []body.UpdateProductDetailRequest{{
					ProductDetailID: "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
					Price:           temp,
					Stock:           temp,
					Weight:          temp,
					Size:            temp,
					Hazardous:       true,
					Codition:        "test",
					BulkPrice:       true,
					Photo:           []string{"test", "test"},
					VariantDetailID: []body.UpdateVariant{{
						VariantID:       "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
						VariantDetailID: "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
					}},
					VariantIDRemove: []string{"989d94b7-58fc-4a76-ae01-1c1b47a0755c"},
				}},
				ProductDetailRemove: []string{"989d94b7-58fc-4a76-ae01-1c1b47a0755c"},
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},

		{
			name:       "update product unauthorized",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},

		{
			name: "invalid request",
			body: nil,
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},

		{
			name: "update  product error internal",
			body: body.UpdateProductRequest{
				ProductInfo: body.UpdateProductInfo{
					Title:        "test",
					Description:  "description",
					Thumbnail:    "test",
					ListedStatus: true,
				},
				ProductDetail: []body.UpdateProductDetailRequest{{
					ProductDetailID: "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
					Price:           temp,
					Stock:           temp,
					Weight:          temp,
					Size:            temp,
					Hazardous:       true,
					Codition:        "test",
					BulkPrice:       true,
					Photo:           []string{"test", "test"},
					VariantDetailID: []body.UpdateVariant{{
						VariantID:       "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
						VariantDetailID: "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
					}},
					VariantIDRemove: []string{"989d94b7-58fc-4a76-ae01-1c1b47a0755c"},
				}},
				ProductDetailRemove: []string{"989d94b7-58fc-4a76-ae01-1c1b47a0755c"},
			},
			mock: func(s *mocks.UseCase) {
				s.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
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

			r := httptest.NewRequest(http.MethodPut, "/api/v1/product/:product_id", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Request = r
			if tc.authorized {

				c.Set("userID", "989d94b7-58fc-4a76-ae01-1c1b47a0755c")

			}
			c.Params = []gin.Param{
				{
					Key:   "id",
					Value: "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
				},
			}

			c.Request.Header.Set("Content-Type", "application/json")
			MockJsonPUT(c, tc.body)

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
			h.UpdateProduct(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestProductHandlers_GetProductReviews(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		queries    map[string]string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success get product review query 1",
			body: nil,

			mock: func(s *mocks.UseCase) {
				s.On("GetProductReviews", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "success get product review query 2",
			body: nil,
			queries: map[string]string{
				"sort":         "asc",
				"show_comment": "false",
				"show_image":   "false",
				"rating":       "3",
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetProductReviews", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&pagination.Pagination{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "get product review error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetProductReviews", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "get product review error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetProductReviews", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/product/:product_id/review", nil)
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
			h.GetProductReviews(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestCartHandlers_CreateProductReview(t *testing.T) {

	var temp string = "test"
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success create review product",
			body: body.ReviewProductRequest{
				ProductID: "123456",
				Comment:   &temp,
				Rating:    1,
				PhotoURL:  &temp,
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateProductReview", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusCreated,
			authorized: true,
		},

		{
			name: "create review  error custom",
			body: body.ReviewProductRequest{
				ProductID: "123456",
				Comment:   &temp,
				Rating:    1,
				PhotoURL:  &temp,
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateProductReview", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name: "success create unauthorized",
			body: body.ReviewProductRequest{
				ProductID: "123456",
				Comment:   &temp,
				Rating:    1,
				PhotoURL:  &temp,
			},
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "invalid request",
			body: body.ReviewProductRequest{
				ProductID: "        ",
				Comment:   &temp,
				Rating:    1,
				PhotoURL:  &temp,
			},
			mock: func(s *mocks.UseCase) {
			},
			expected:   http.StatusUnprocessableEntity,
			authorized: true,
		},

		{
			name: "add create review product error internal",
			body: body.ReviewProductRequest{
				ProductID: "123456",
				Comment:   &temp,
				Rating:    1,
				PhotoURL:  &temp,
			},
			mock: func(s *mocks.UseCase) {
				s.On("CreateProductReview", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/product/favorite", bytes.NewBuffer(jsonValue))
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

			h := NewProductHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CreateProductReview(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestCartHandlers_DeleteProductReview(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success delete review product",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("DeleteProductReview", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "delete review  error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("DeleteProductReview", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
		{
			name:       "success delete unauthorized",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusUnauthorized,
			authorized: false,
		},
		{
			name: "add delete review product error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("DeleteProductReview", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
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

			r := httptest.NewRequest(http.MethodDelete, "/api/v1/product/review/:review_id", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)

			c.Params = []gin.Param{
				{
					Key:   "review_id",
					Value: "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
				},
			}
			c.Request = r
			if tc.authorized {
				c.Set("userID", "989d94b7-58fc-4a76-ae01-1c1b47a0755c")
			}

			c.Request.Header.Set("Content-Type", "application/json")

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
			h.DeleteProductReview(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestProductHandlers_GetTotalReviewRatingByProductID(t *testing.T) {
	testCase := []struct {
		name       string
		body       interface{}
		queries    map[string]string
		mock       func(s *mocks.UseCase)
		expected   int
		authorized bool
	}{
		{
			name: "success get product total review rating ",
			body: nil,

			mock: func(s *mocks.UseCase) {
				s.On("GetTotalReviewRatingByProductID", mock.Anything, mock.Anything).Return(&body.AllRatingProduct{}, nil)
			},
			expected:   http.StatusOK,
			authorized: true,
		},
		{
			name: "get product total review rating error internal",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetTotalReviewRatingByProductID", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			authorized: true,
		},
		{
			name: "get product total review rating error custom",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetTotalReviewRatingByProductID", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			authorized: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/product/:product_id/review/rating", nil)
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
			h.GetTotalReviewRatingByProductID(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}
