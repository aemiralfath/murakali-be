package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/auth/delivery/body"
	"murakali/internal/module/auth/mocks"
	"murakali/pkg/httperror"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthUseCase_Login(t *testing.T) {
	passwordHash := "$2a$10$WKul/6gjYoYjOXuNVX4XGen1ZkWYb1PKFiI5vlZp5TFerZh6nTujG"
	testCase := []struct {
		name        string
		body        body.LoginRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success login",
			body: body.LoginRequest{Email: "a@test.com", Password: "Tested8*"},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("InsertSessionRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
				r.On("InsertSessionRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "email not found",
			body: body.LoginRequest{Email: "a@test.com", Password: "Tested8*"},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByEmail", mock.Anything, mock.Anything).
					Return(nil, sql.ErrNoRows)
			},
			expectedErr: httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage),
		},
		{
			name: "email error",
			body: body.LoginRequest{Email: "a@test.com", Password: "Tested8*"},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByEmail", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "wrong password login",
			body: body.LoginRequest{Email: "a@test.com", Password: "Tested8"},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
			},
			expectedErr: httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage),
		},
		{
			name: "error access redis",
			body: body.LoginRequest{Email: "a@test.com", Password: "Tested8*"},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("InsertSessionRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error refresh redis",
			body: body.LoginRequest{Email: "a@test.com", Password: "Tested8*"},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("InsertSessionRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
				r.On("InsertSessionRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewAuthUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.Login(context.Background(), tc.body)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAuthUseCase_RegisterEmail(t *testing.T) {
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "check email history",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(&model.EmailHistory{}, fmt.Errorf("User already registered."))

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.EmailAlreadyExistMessage),
		},
		{
			name: "check email history",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error get user by email",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(nil, nil)
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAuthUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.RegisterEmail(context.Background(), body.RegisterEmailRequest{Email: "sammy@gmail.com"})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAuthUseCase_RegisterUser(t *testing.T) {
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success register user",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, nil)
				r.On("GetUserByPhoneNo", mock.Anything, mock.Anything).Return(nil, nil)
				r.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("CreateEmailHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, response.EmailAlreadyExistMessage),
		},
		{
			name: "get user email error ",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "get user email error no rows ",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.UserNotExistMessage),
		},
		{
			name: "get username error ",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "get username error same username",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(&model.User{}, nil)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.UserNameAlreadyExistMessage),
		},
		{
			name: "get phone number error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, nil)
				r.On("GetUserByPhoneNo", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "get phone number error same phone number",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, nil)
				r.On("GetUserByPhoneNo", mock.Anything, mock.Anything).Return(&model.User{}, nil)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.PhoneNoAlreadyExistMessage),
		},
		{
			name: "get update user error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, nil)
				r.On("GetUserByPhoneNo", mock.Anything, mock.Anything).Return(nil, nil)
				r.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAuthUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.RegisterUser(context.Background(), "sammy@gmail.com", body.RegisterUserRequest{})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAuthUseCase_ResetPasswordEmail(t *testing.T) {
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success check unique username",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(&model.EmailHistory{}, nil)
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{IsVerify: true}, nil)
				r.On("InsertNewOTPHashedKey", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, response.EmailAlreadyExistMessage),
		},
		{
			name: "error check email history",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error check email history",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(nil, nil)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.EmailNotExistMessage),
		},
		{
			name: "error get user by email",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(&model.EmailHistory{}, nil)
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error get user by email user nil",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(&model.EmailHistory{}, nil)
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, response.UserNotExistMessage),
		},
		{
			name: "error get user by email is not verify",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(&model.EmailHistory{}, nil)
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{IsVerify: false}, nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, response.UserNotVerifyMessage),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAuthUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.ResetPasswordEmail(context.Background(), body.ResetPasswordEmailRequest{})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAuthUseCase_ResetPasswordUser(t *testing.T) {

	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "error check email histoery",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error  get user by email",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(&model.EmailHistory{}, nil)
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error  get user by email ",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(&model.EmailHistory{}, nil)
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, response.UserNotExistMessage),
		},
		{
			name: "error  get user by email user nil",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(&model.EmailHistory{}, nil)
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, response.UserNotExistMessage),
		},
		{
			name: "error  get user by email user nil",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CheckEmailHistory", mock.Anything, mock.Anything).Return(&model.EmailHistory{}, nil)
				r.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{IsVerify: false}, nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, response.UserNotVerifyMessage),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAuthUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.ResetPasswordUser(context.Background(), "sammy@gmail.com", &body.ResetPasswordUserRequest{})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAuthUseCase_CheckUniqueUsername(t *testing.T) {
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success check unique username",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(&model.User{}, nil)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.EmailAlreadyExistMessage),
		},
		{
			name: "error get user by username",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error get user by username",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

			},
			expectedErr: nil,
		},
		{
			name: "error get user by username",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByUsername", mock.Anything, mock.Anything).Return(&model.User{}, nil)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.PhoneNoAlreadyExistMessage),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAuthUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.CheckUniqueUsername(context.Background(), "87738171235")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAuthUseCase_CheckUniquePhoneNo(t *testing.T) {
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success check unique phone no",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByPhoneNo", mock.Anything, mock.Anything).Return(&model.User{}, nil)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.EmailAlreadyExistMessage),
		},
		{
			name: "error get user by phone no",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByPhoneNo", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error get user by phone no",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByPhoneNo", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAuthUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.CheckUniquePhoneNo(context.Background(), "87738171235")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
