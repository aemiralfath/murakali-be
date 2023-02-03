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
	"murakali/internal/model"
	"murakali/internal/module/auth/delivery/body"
	"murakali/internal/module/auth/mocks"
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

func TestAuthHandlers_Logout(t *testing.T) {
	testCase := []struct {
		name     string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name:     "success logout",
			body:     nil,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusOK,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/auth/logout", nil)
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

			h := NewAuthHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.Logout(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAuthHandlers_Login(t *testing.T) {
	invalidRequestBody := struct {
		Email int `json:"email"`
	}{123}

	testCase := []struct {
		name     string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name: "success login",
			body: body.LoginRequest{
				Email:    "emir@gmail.com",
				Password: "Tested8*",
			},
			mock: func(s *mocks.UseCase) {
				s.On("Login", mock.Anything, mock.Anything).Return(&model.Token{AccessToken: &model.AccessToken{}, RefreshToken: &model.RefreshToken{}}, nil)
			},
			expected: http.StatusOK,
		},
		{
			name:     "invalid request",
			body:     invalidRequestBody,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
		},
		{
			name: "invalid request entity",
			body: body.LoginRequest{
				Email:    "",
				Password: "",
			},
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusUnprocessableEntity,
		},
		{
			name: "register error internal",
			body: body.LoginRequest{
				Email:    "emir@gmail.com",
				Password: "Tested8*",
			},
			mock: func(s *mocks.UseCase) {
				s.On("Login", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name: "register error custom",
			body: body.LoginRequest{
				Email:    "emir@gmail.com",
				Password: "Tested8*",
			},
			mock: func(s *mocks.UseCase) {
				s.On("Login", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
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

			h := NewAuthHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.Login(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}
