package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"murakali/config"
	"murakali/internal/module/location/delivery/body"
	"murakali/internal/module/location/mocks"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestLocationHandlers_GetShippingCost(t *testing.T) {
	invalidRequestBody := struct {
		ShopID int `json:"shop_id"`
	}{123}

	testCase := []struct {
		name     string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name: "success get shipping cost",
			body: body.GetShippingCostRequest{
				ShopID:      "168f5ed0-dda1-433e-b518-abe7674deb71",
				ProductIDS:  []string{"168f5ed0-dda1-433e-b518-abe7674deb71"},
				Weight:      1000,
				Destination: 12,
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetShippingCost", mock.Anything, mock.Anything).Return(&body.GetShippingCostResponse{}, nil)
			},
			expected: http.StatusOK,
		},
		{
			name: "internal server get shipping cost",
			body: body.GetShippingCostRequest{
				ShopID:      "168f5ed0-dda1-433e-b518-abe7674deb71",
				ProductIDS:  []string{"168f5ed0-dda1-433e-b518-abe7674deb71"},
				Weight:      1000,
				Destination: 12,
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetShippingCost", mock.Anything, mock.Anything).Return(&body.GetShippingCostResponse{}, errors.New("test"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name: "custom error get shipping cost",
			body: body.GetShippingCostRequest{
				ShopID:      "168f5ed0-dda1-433e-b518-abe7674deb71",
				ProductIDS:  []string{"168f5ed0-dda1-433e-b518-abe7674deb71"},
				Weight:      1000,
				Destination: 12,
			},
			mock: func(s *mocks.UseCase) {
				s.On("GetShippingCost", mock.Anything, mock.Anything).Return(&body.GetShippingCostResponse{}, httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
		{
			name:     "invalid request get shipping cost",
			body:     invalidRequestBody,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
		},
		{
			name: "invalid field get shipping cost",
			body: body.GetShippingCostRequest{
				ShopID:      "168f5ed0-dda1-433e-b518-abe7674deb71",
				ProductIDS:  []string{"168f5ed0-dda1-433e-b518-abe7674deb71"},
				Weight:      1000,
				Destination: 0,
			},
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusUnprocessableEntity,
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/location/cost", bytes.NewBuffer(jsonValue))
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

			h := NewLocationHandlers(cfg, s, appLogger)
			tc.mock(s)
			h.GetShippingCost(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestLocationHandlers_GetCity(t *testing.T) {
	testCase := []struct {
		name     string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
		isQuery  bool
	}{
		{
			name: "success get city",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCity", mock.Anything, mock.Anything).Return(&body.CityResponse{}, nil)
			},
			expected: http.StatusOK,
			isQuery:  true,
		},
		{
			name: "internal error get city",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCity", mock.Anything, mock.Anything).Return(&body.CityResponse{}, errors.New("test"))
			},
			expected: http.StatusInternalServerError,
			isQuery:  true,
		},
		{
			name: "custom error get city",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCity", mock.Anything, mock.Anything).Return(&body.CityResponse{}, httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
			isQuery:  true,
		},
		{
			name: "success get city",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetCity", mock.Anything, mock.Anything).Return(&body.CityResponse{}, nil)
			},
			expected: http.StatusOK,
			isQuery:  true,
		},
		{
			name:     "id error get city",
			body:     nil,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
			isQuery:  false,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/location/city", nil)
			r.Header = make(http.Header)

			if tc.isQuery {
				query := r.URL.Query()
				query.Add("province_id", "1")
				r.URL.RawQuery = query.Encode()
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

			h := NewLocationHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetCity(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestLocationHandlers_GetSubDistrict(t *testing.T) {
	testCase := []struct {
		name     string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
		isQuery  bool
	}{
		{
			name: "success get city",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetSubDistrict", mock.Anything, mock.Anything, mock.Anything).Return(&body.SubDistrictResponse{}, nil)
			},
			expected: http.StatusOK,
			isQuery:  true,
		},
		{
			name: "internal error get city",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetSubDistrict", mock.Anything, mock.Anything, mock.Anything).Return(&body.SubDistrictResponse{}, errors.New("test"))
			},
			expected: http.StatusInternalServerError,
			isQuery:  true,
		},
		{
			name: "custom error get city",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetSubDistrict", mock.Anything, mock.Anything, mock.Anything).Return(&body.SubDistrictResponse{}, httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
			isQuery:  true,
		},
		{
			name:     "query error get city",
			body:     nil,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
			isQuery:  false,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/location/city/subdistrict", nil)
			r.Header = make(http.Header)

			if tc.isQuery {
				query := r.URL.Query()
				query.Add("province", "jawa barat")
				query.Add("city", "bandung")
				r.URL.RawQuery = query.Encode()
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

			h := NewLocationHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetSubDistrict(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestLocationHandlers_GetUrban(t *testing.T) {
	testCase := []struct {
		name     string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
		isQuery  bool
	}{
		{
			name: "success get urban",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetUrban", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&body.UrbanResponse{}, nil)
			},
			expected: http.StatusOK,
			isQuery:  true,
		},
		{
			name: "internal error get urban",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetUrban", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&body.UrbanResponse{}, errors.New("test"))
			},
			expected: http.StatusInternalServerError,
			isQuery:  true,
		},
		{
			name: "custom error get urban",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetUrban", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&body.UrbanResponse{}, httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
			isQuery:  true,
		},
		{
			name:     "error query get urban",
			body:     nil,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
			isQuery:  false,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/location/city/subdistrict/urban", nil)
			r.Header = make(http.Header)

			if tc.isQuery {
				query := r.URL.Query()
				query.Add("province", "jawa barat")
				query.Add("city", "bandung")
				query.Add("subdistrict", "dago")
				r.URL.RawQuery = query.Encode()
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

			h := NewLocationHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetUrban(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestLocationHandlers_GetProvince(t *testing.T) {
	testCase := []struct {
		name     string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name: "success get province",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetProvince", mock.Anything).Return(&body.ProvinceResponse{}, nil)
			},
			expected: http.StatusOK,
		},
		{
			name: "internal error get province",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetProvince", mock.Anything).Return(&body.ProvinceResponse{}, errors.New("test"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name: "custom error get province",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GetProvince", mock.Anything).Return(&body.ProvinceResponse{}, httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/location/province", nil)
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

			h := NewLocationHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.GetProvince(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}
