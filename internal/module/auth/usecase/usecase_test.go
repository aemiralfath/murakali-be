package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/auth/delivery/body"
	"murakali/internal/module/auth/mocks"
	"murakali/pkg/httperror"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
	"testing"
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
