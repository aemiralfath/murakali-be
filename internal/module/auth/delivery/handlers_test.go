package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"murakali/config"
	"murakali/internal/constant"
	"murakali/internal/model"
	"murakali/internal/module/auth/delivery/body"
	"murakali/internal/module/auth/mocks"
	"murakali/pkg/httperror"
	jwt2 "murakali/pkg/jwt"
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
			name: "login error internal",
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
			name: "login error custom",
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

func TestAuthHandlers_RefreshToken(t *testing.T) {
	testCase := []struct {
		name      string
		body      interface{}
		mock      func(s *mocks.UseCase)
		expected  int
		isCookie  bool
		cookieKey string
	}{
		{
			name: "success refresh",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("RefreshToken", mock.Anything, mock.Anything, mock.Anything).Return(&model.AccessToken{}, nil)
			},
			expected:  http.StatusOK,
			isCookie:  true,
			cookieKey: "",
		},
		{
			name:      "cookie error refresh",
			body:      nil,
			mock:      func(s *mocks.UseCase) {},
			expected:  http.StatusForbidden,
			isCookie:  false,
			cookieKey: "",
		},
		{
			name:      "jwt error refresh",
			body:      nil,
			mock:      func(s *mocks.UseCase) {},
			expected:  http.StatusForbidden,
			isCookie:  true,
			cookieKey: "error",
		},
		{
			name: "internal error refresh",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("RefreshToken", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected:  http.StatusInternalServerError,
			isCookie:  true,
			cookieKey: "",
		},
		{
			name: "custom error refresh",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("RefreshToken", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:  http.StatusBadRequest,
			isCookie:  true,
			cookieKey: "",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/auth/refresh", nil)

			r.Header = make(http.Header)
			if tc.isCookie {
				tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt2.RefreshClaims{}).SignedString([]byte(tc.cookieKey))
				r.AddCookie(&http.Cookie{Name: constant.RefreshTokenCookie, Value: tokenString})
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

			h := NewAuthHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.RefreshToken(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAuthHandlers_RegisterEmail(t *testing.T) {
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
			name: "success register",
			body: body.RegisterEmailRequest{
				Email: "emir@gmail.com",
			},
			mock: func(s *mocks.UseCase) {
				s.On("RegisterEmail", mock.Anything, mock.Anything).Return(&model.User{}, nil)
			},
			expected: http.StatusOK,
		},
		{
			name:     "invalid body register",
			body:     invalidRequestBody,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
		},
		{
			name: "invalid fields register",
			body: body.RegisterEmailRequest{
				Email: "",
			},
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusUnprocessableEntity,
		},
		{
			name: "internal error register",
			body: body.RegisterEmailRequest{
				Email: "emir@gmail.com",
			},
			mock: func(s *mocks.UseCase) {
				s.On("RegisterEmail", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name: "custom error register",
			body: body.RegisterEmailRequest{
				Email: "emir@gmail.com",
			},
			mock: func(s *mocks.UseCase) {
				s.On("RegisterEmail", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, "tst"))
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonValue))
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
			h.RegisterEmail(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAuthHandlers_RegisterUser(t *testing.T) {
	invalidRequestBody := struct {
		Username int `json:"username"`
	}{123}

	testCase := []struct {
		name      string
		body      interface{}
		mock      func(s *mocks.UseCase)
		expected  int
		isCookie  bool
		cookieKey string
	}{
		{
			name: "success register user",
			body: &body.RegisterUserRequest{
				Username: "emir",
				FullName: "emir",
				Password: "Tested9*",
				PhoneNo:  "83187115995",
			},
			mock: func(s *mocks.UseCase) {
				s.On("RegisterUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected:  http.StatusOK,
			isCookie:  true,
			cookieKey: "",
		},
		{
			name: "internal error register user",
			body: &body.RegisterUserRequest{
				Username: "emir",
				FullName: "emir",
				Password: "Tested9*",
				PhoneNo:  "83187115995",
			},
			mock: func(s *mocks.UseCase) {
				s.On("RegisterUser", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
			},
			expected:  http.StatusInternalServerError,
			isCookie:  true,
			cookieKey: "",
		},
		{
			name: "custom error register user",
			body: &body.RegisterUserRequest{
				Username: "emir",
				FullName: "emir",
				Password: "Tested9*",
				PhoneNo:  "83187115995",
			},
			mock: func(s *mocks.UseCase) {
				s.On("RegisterUser", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(http.StatusBadRequest, "test"))
			},
			expected:  http.StatusBadRequest,
			isCookie:  true,
			cookieKey: "",
		},
		{
			name: "body empty register user",
			body: &body.RegisterUserRequest{
				Username: "",
				FullName: "emir",
				Password: "Tested9*",
				PhoneNo:  "83187115995",
			},
			mock:      func(s *mocks.UseCase) {},
			expected:  http.StatusUnprocessableEntity,
			isCookie:  true,
			cookieKey: "",
		},
		{
			name:      "invalid body register user",
			body:      &invalidRequestBody,
			mock:      func(s *mocks.UseCase) {},
			expected:  http.StatusBadRequest,
			isCookie:  true,
			cookieKey: "",
		},
		{
			name: "cookie error user",
			body: &body.RegisterUserRequest{
				Username: "emir",
				FullName: "emir",
				Password: "Tested9*",
				PhoneNo:  "83187115995",
			},
			mock:      func(s *mocks.UseCase) {},
			expected:  http.StatusForbidden,
			isCookie:  false,
			cookieKey: "",
		},
		{
			name: "jwt error register user",
			body: &body.RegisterUserRequest{
				Username: "emir",
				FullName: "emir",
				Password: "Tested9*",
				PhoneNo:  "83187115995",
			},
			mock:      func(s *mocks.UseCase) {},
			expected:  http.StatusForbidden,
			isCookie:  true,
			cookieKey: "test",
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

			r := httptest.NewRequest(http.MethodPut, "/api/v1/auth/register", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)
			if tc.isCookie {
				tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt2.RegisterClaims{Email: "emir@gmail.com"}).SignedString([]byte(tc.cookieKey))
				r.AddCookie(&http.Cookie{Name: constant.RegisterTokenCookie, Value: tokenString})
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

			h := NewAuthHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.RegisterUser(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAuthHandlers_VerifyOTP(t *testing.T) {
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
			name: "success verify",
			body: body.VerifyOTPRequest{
				Email: "emir@gmail.com",
				OTP:   "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("VerifyOTP", mock.Anything, mock.Anything).Return("", nil)
			},
			expected: http.StatusOK,
		},
		{
			name: "internal error verify",
			body: body.VerifyOTPRequest{
				Email: "emir@gmail.com",
				OTP:   "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("VerifyOTP", mock.Anything, mock.Anything).Return("", errors.New("test"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name: "custom error verify",
			body: body.VerifyOTPRequest{
				Email: "emir@gmail.com",
				OTP:   "123456",
			},
			mock: func(s *mocks.UseCase) {
				s.On("VerifyOTP", mock.Anything, mock.Anything).Return("", httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "invalid field verify",
			body: body.VerifyOTPRequest{
				Email: "emir@gmail.com",
				OTP:   "",
			},
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusUnprocessableEntity,
		},
		{
			name:     "invalid body verify",
			body:     invalidRequestBody,
			mock:     func(s *mocks.UseCase) {},
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/verify", bytes.NewBuffer(jsonValue))
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
			h.VerifyOTP(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAuthHandlers_ResetPasswordEmail(t *testing.T) {
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
			name: "success reset password email",
			body: body.ResetPasswordEmailRequest{
				Email: "emir@gmail.com",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ResetPasswordEmail", mock.Anything, mock.Anything).Return(&model.User{}, nil)
			},
			expected: http.StatusOK,
		},
		{
			name: "internal error reset password email",
			body: body.ResetPasswordEmailRequest{
				Email: "emir@gmail.com",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ResetPasswordEmail", mock.Anything, mock.Anything).Return(&model.User{}, errors.New("test"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name: "custom error reset password email",
			body: body.ResetPasswordEmailRequest{
				Email: "emir@gmail.com",
			},
			mock: func(s *mocks.UseCase) {
				s.On("ResetPasswordEmail", mock.Anything, mock.Anything).Return(&model.User{}, httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "invalid field reset password email",
			body: body.ResetPasswordEmailRequest{
				Email: "",
			},
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusUnprocessableEntity,
		},
		{
			name:     "invalid body reset password email",
			body:     invalidRequestBody,
			mock:     func(s *mocks.UseCase) {},
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

			r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewBuffer(jsonValue))
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
			h.ResetPasswordEmail(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAuthHandlers_ResetPasswordUser(t *testing.T) {
	invalidRequestBody := struct {
		Password int `json:"password"`
	}{123}

	testCase := []struct {
		name      string
		body      interface{}
		mock      func(s *mocks.UseCase)
		expected  int
		isCookie  bool
		cookieKey string
	}{
		{
			name: "success reset pw user",
			body: &body.ResetPasswordUserRequest{Password: "Tested8*"},
			mock: func(s *mocks.UseCase) {
				s.On("ResetPasswordUser", mock.Anything, mock.Anything, mock.Anything).Return(&model.User{}, nil)
			},
			expected:  http.StatusOK,
			isCookie:  true,
			cookieKey: "",
		},
		{
			name: "internal error reset pw user",
			body: &body.ResetPasswordUserRequest{Password: "Tested8*"},
			mock: func(s *mocks.UseCase) {
				s.On("ResetPasswordUser", mock.Anything, mock.Anything, mock.Anything).Return(&model.User{}, errors.New("test"))
			},
			expected:  http.StatusInternalServerError,
			isCookie:  true,
			cookieKey: "",
		},
		{
			name: "custom error reset pw user",
			body: &body.ResetPasswordUserRequest{Password: "Tested8*"},
			mock: func(s *mocks.UseCase) {
				s.On("ResetPasswordUser", mock.Anything, mock.Anything, mock.Anything).Return(&model.User{}, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:  http.StatusBadRequest,
			isCookie:  true,
			cookieKey: "",
		},
		{
			name:      "invalid field reset pw user",
			body:      &body.ResetPasswordUserRequest{Password: ""},
			mock:      func(s *mocks.UseCase) {},
			expected:  http.StatusUnprocessableEntity,
			isCookie:  true,
			cookieKey: "",
		},
		{
			name:      "cookie invalid reset pw user",
			body:      &body.ResetPasswordUserRequest{Password: "Tested8*"},
			mock:      func(s *mocks.UseCase) {},
			expected:  http.StatusForbidden,
			isCookie:  false,
			cookieKey: "",
		},
		{
			name:      "cookie invalid key reset pw user",
			body:      &body.ResetPasswordUserRequest{Password: "Tested8*"},
			mock:      func(s *mocks.UseCase) {},
			expected:  http.StatusForbidden,
			isCookie:  true,
			cookieKey: "test",
		},
		{
			name:      "invalid body reset pw user",
			body:      invalidRequestBody,
			mock:      func(s *mocks.UseCase) {},
			expected:  http.StatusBadRequest,
			isCookie:  true,
			cookieKey: "",
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

			r := httptest.NewRequest(http.MethodPatch, "/api/v1/auth/reset-password", bytes.NewBuffer(jsonValue))
			r.Header = make(http.Header)
			if tc.isCookie {
				tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt2.ResetPasswordClaims{Email: "emir@gmail.com"}).SignedString([]byte(tc.cookieKey))
				r.AddCookie(&http.Cookie{Name: constant.ResetPasswordTokenCookie, Value: tokenString})
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

			h := NewAuthHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.ResetPasswordUser(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAuthHandlers_ResetPasswordVerifyOTP(t *testing.T) {
	testCase := []struct {
		name     string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
		isQuery  bool
	}{
		{
			name: "success reset pw verify",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("ResetPasswordVerifyOTP", mock.Anything, mock.Anything).Return("", nil)
			},
			expected: http.StatusOK,
			isQuery:  true,
		},
		{
			name: "internal error reset pw verify",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("ResetPasswordVerifyOTP", mock.Anything, mock.Anything).Return("", errors.New("test"))
			},
			expected: http.StatusInternalServerError,
			isQuery:  true,
		},
		{
			name: "custom error reset pw verify",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("ResetPasswordVerifyOTP", mock.Anything, mock.Anything).Return("", httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
			isQuery:  true,
		},
		{
			name:     "query err pw verify",
			body:     nil,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusUnprocessableEntity,
			isQuery:  false,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/auth/verify", nil)
			r.Header = make(http.Header)

			if tc.isQuery {
				query := r.URL.Query()
				query.Add("code", "123455")
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

			h := NewAuthHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.ResetPasswordVerifyOTP(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAuthHandlers_CheckUniqueUsername(t *testing.T) {
	testCase := []struct {
		name     string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
	}{
		{
			name: "success check username",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("CheckUniqueUsername", mock.Anything, mock.Anything).Return(true, nil)
			},
			expected: http.StatusOK,
		},
		{
			name: "internal error check username",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("CheckUniqueUsername", mock.Anything, mock.Anything).Return(true, errors.New("test"))
			},
			expected: http.StatusInternalServerError,
		},
		{
			name: "custom error check username",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("CheckUniqueUsername", mock.Anything, mock.Anything).Return(true, httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/auth/unique/username/:username", nil)
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
			h.CheckUniqueUsername(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAuthHandlers_CheckUniquePhoneNo(t *testing.T) {
	testCase := []struct {
		name     string
		body     interface{}
		mock     func(s *mocks.UseCase)
		expected int
		PhoneNo  string
	}{
		{
			name: "success check phone",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("CheckUniquePhoneNo", mock.Anything, mock.Anything).Return(true, nil)
			},
			expected: http.StatusOK,
			PhoneNo:  "83187115995",
		},
		{
			name: "internal error check phone",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("CheckUniquePhoneNo", mock.Anything, mock.Anything).Return(true, errors.New("test"))
			},
			expected: http.StatusInternalServerError,
			PhoneNo:  "83187115995",
		},
		{
			name: "custom error check phone",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("CheckUniquePhoneNo", mock.Anything, mock.Anything).Return(true, httperror.New(http.StatusBadRequest, "test"))
			},
			expected: http.StatusBadRequest,
			PhoneNo:  "83187115995",
		},
		{
			name:     "invalid string phone",
			body:     nil,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
			PhoneNo:  "831871159s",
		},
		{
			name:     "invalid string phone format",
			body:     nil,
			mock:     func(s *mocks.UseCase) {},
			expected: http.StatusBadRequest,
			PhoneNo:  "8318",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/auth/unique/phone-no/:phone_no", nil)
			r.Header = make(http.Header)

			c.Params = []gin.Param{
				{
					Key:   "phone_no",
					Value: tc.PhoneNo,
				},
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

			h := NewAuthHandlers(cfg, s, appLogger)

			tc.mock(s)
			h.CheckUniquePhoneNo(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}

func TestAuthHandlers_GoogleAuth(t *testing.T) {
	registerToken := "fjskdjfkldsjkfjdskljfkldsjkf"
	testCase := []struct {
		name       string
		body       interface{}
		mock       func(s *mocks.UseCase)
		expected   int
		queryCode  string
		queryState string
	}{
		{
			name: "success google auth register",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GoogleAuth", mock.Anything, mock.Anything, mock.Anything).Return(&model.GoogleAuthToken{RegisterToken: &registerToken}, nil)
			},
			expected:   http.StatusOK,
			queryCode:  "123456",
			queryState: "",
		},
		{
			name: "success google auth login",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GoogleAuth", mock.Anything, mock.Anything, mock.Anything).Return(&model.GoogleAuthToken{RegisterToken: nil, Token: &model.Token{&model.AccessToken{}, &model.RefreshToken{}}}, nil)
			},
			expected:   http.StatusOK,
			queryCode:  "123456",
			queryState: "",
		},
		{
			name: "internal error google auth register",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GoogleAuth", mock.Anything, mock.Anything, mock.Anything).Return(&model.GoogleAuthToken{RegisterToken: &registerToken}, errors.New("test"))
			},
			expected:   http.StatusInternalServerError,
			queryCode:  "123456",
			queryState: "",
		},
		{
			name: "custom error google auth register",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GoogleAuth", mock.Anything, mock.Anything, mock.Anything).Return(&model.GoogleAuthToken{RegisterToken: &registerToken}, httperror.New(http.StatusBadRequest, "test"))
			},
			expected:   http.StatusBadRequest,
			queryCode:  "123456",
			queryState: "",
		},
		{
			name: "failed google auth",
			body: nil,
			mock: func(s *mocks.UseCase) {
				s.On("GoogleAuth", mock.Anything, mock.Anything, mock.Anything).Return(&model.GoogleAuthToken{RegisterToken: &registerToken}, nil)
			},
			expected:   http.StatusOK,
			queryCode:  "123456",
			queryState: "",
		},
		{
			name: "failed google auth state",
			body: nil,
			mock: func(s *mocks.UseCase) {

				s.On("GoogleAuth", mock.Anything, mock.Anything, mock.Anything).Return(&model.GoogleAuthToken{RegisterToken: &registerToken}, nil)
			},
			expected:   http.StatusOK,
			queryCode:  "123456",
			queryState: "/home",
		},
		{
			name:       "failed google auth code",
			body:       nil,
			mock:       func(s *mocks.UseCase) {},
			expected:   http.StatusForbidden,
			queryCode:  "",
			queryState: "/home",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google-oauth", nil)
			r.Header = make(http.Header)

			query := r.URL.Query()
			query.Add("code", tc.queryCode)
			query.Add("state", tc.queryState)
			r.URL.RawQuery = query.Encode()

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
			h.GoogleAuth(c)

			assert.Equal(t, rr.Code, tc.expected)
		})
	}
}
