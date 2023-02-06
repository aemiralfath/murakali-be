package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/admin/delivery/body"
	"murakali/internal/module/admin/mocks"
	"murakali/pkg/httperror"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdminUC_GetAllVoucher(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success get voucher",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalVoucher", mock.Anything, mock.Anything).Return(int64(1), nil)
				r.On("GetAllVoucher", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*model.Voucher{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "failed get voucher",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalVoucher", mock.Anything, mock.Anything).Return(int64(0), fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "failed get voucher",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalVoucher", mock.Anything, mock.Anything).Return(int64(1), nil)
				r.On("GetAllVoucher", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetAllVoucher(context.Background(), "123", "123", &pagination.Pagination{})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_GetRefunds(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success get voucher",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalRefunds", mock.Anything).Return(int64(1), nil)
				r.On("GetRefunds", mock.Anything, mock.Anything, mock.Anything).Return([]*model.RefundOrder{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "failed get voucher",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalRefunds", mock.Anything).Return(int64(0), fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "failed get voucher",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalRefunds", mock.Anything).Return(int64(1), nil)
				r.On("GetRefunds", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetRefunds(context.Background(), "123", &pagination.Pagination{})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_CreateVoucher(t *testing.T) {
	var temp float64 = 10
	datestring := "02-01-2006 15:04:05"
	date, _ := time.Parse("02-01-2006 15:04:05", datestring)
	testCase := []struct {
		name        string
		body        model.Voucher
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success create voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CountCodeVoucher", mock.Anything, mock.Anything).Return(int64(0), nil)
				r.On("CreateVoucher", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "failed create voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CountCodeVoucher", mock.Anything, mock.Anything).Return(int64(1), nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, body.CodeVoucherAlreadyExist),
		},
		{
			name: "failed create voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CountCodeVoucher", mock.Anything, mock.Anything).Return(int64(0), nil)
				r.On("CreateVoucher", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.CreateVoucher(context.Background(), body.CreateVoucherRequest{
				Code:               "123",
				Quota:              123,
				ActivedDate:        "02-01-2006 15:04:05",
				ExpiredDate:        "02-01-2006 15:04:05",
				DiscountPercentage: temp,
				DiscountFixPrice:   temp,
				MinProductPrice:    temp,
				MaxDiscountPrice:   temp,
			})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_UpdateVoucher(t *testing.T) {
	var temp float64 = 10
	datestring := "02-01-2006 15:04:05"
	date, _ := time.Parse("02-01-2006 15:04:05", datestring)
	testCase := []struct {
		name        string
		body        model.Voucher
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetVoucherByID", mock.Anything, mock.Anything).Return(&model.Voucher{}, nil)
				r.On("UpdateVoucher", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "failed update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetVoucherByID", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage))
			},
			expectedErr: httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage),
		},
		{
			name: "failed update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetVoucherByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetVoucherByID", mock.Anything, mock.Anything).Return(&model.Voucher{}, nil)
				r.On("UpdateVoucher", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.UpdateVoucher(context.Background(), body.UpdateVoucherRequest{
				VoucherID:          "123",
				Quota:              123,
				ActivedDate:        "02-01-2006 15:04:05",
				ExpiredDate:        "02-01-2006 15:04:05",
				DiscountPercentage: temp,
				DiscountFixPrice:   temp,
				MinProductPrice:    temp,
				MaxDiscountPrice:   temp,
			})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_GetDetailVoucher(t *testing.T) {
	var temp float64 = 10
	datestring := "02-01-2006 15:04:05"
	date, _ := time.Parse("02-01-2006 15:04:05", datestring)
	testCase := []struct {
		name        string
		body        model.Voucher
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetVoucherByID", mock.Anything, mock.Anything).Return(&model.Voucher{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "failed update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetVoucherByID", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage))
			},
			expectedErr: httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage),
		},
		{
			name: "failed update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetVoucherByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetDetailVoucher(context.Background(), "123")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_DeleteVoucher(t *testing.T) {
	var temp float64 = 10
	datestring := "02-01-2006 15:04:05"
	date, _ := time.Parse("02-01-2006 15:04:05", datestring)
	testCase := []struct {
		name        string
		body        model.Voucher
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetVoucherByID", mock.Anything, mock.Anything).Return(&model.Voucher{}, nil)
				r.On("DeleteVoucher", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "failed update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetVoucherByID", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage))
			},
			expectedErr: httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage),
		},
		{
			name: "failed update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetVoucherByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetVoucherByID", mock.Anything, mock.Anything).Return(&model.Voucher{}, nil)
				r.On("DeleteVoucher", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.DeleteVoucher(context.Background(), "123")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_GetCategories(t *testing.T) {
	var temp float64 = 10
	datestring := "02-01-2006 15:04:05"
	date, _ := time.Parse("02-01-2006 15:04:05", datestring)
	testCase := []struct {
		name        string
		body        model.Voucher
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCategories", mock.Anything).Return([]*body.CategoryResponse{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCategories", mock.Anything).Return([]*body.CategoryResponse{}, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetCategories(context.Background())
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_AddCategory(t *testing.T) {
	var temp float64 = 10
	datestring := "02-01-2006 15:04:05"
	date, _ := time.Parse("02-01-2006 15:04:05", datestring)
	testCase := []struct {
		name        string
		body        model.Voucher
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("AddCategory", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("AddCategory", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.AddCategory(context.Background(), body.CategoryRequest{
				ID:       "asd",
				ParentID: "asd",
				Name:     "asd",
				PhotoURL: "example.com",
				Level:    "1",
			})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_DeleteCategory(t *testing.T) {
	var temp float64 = 10
	datestring := "02-01-2006 15:04:05"
	date, _ := time.Parse("02-01-2006 15:04:05", datestring)
	testCase := []struct {
		name        string
		body        model.Voucher
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CountProductCategory", mock.Anything, mock.Anything).Return(0, nil)
				r.On("CountCategoryParent", mock.Anything, mock.Anything).Return(0, nil)
				r.On("DeleteCategory", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CountProductCategory", mock.Anything, mock.Anything).Return(0, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CountProductCategory", mock.Anything, mock.Anything).Return(1, nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, body.CategoryIsBeingUsed),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CountProductCategory", mock.Anything, mock.Anything).Return(0, nil)
				r.On("CountCategoryParent", mock.Anything, mock.Anything).Return(0, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CountProductCategory", mock.Anything, mock.Anything).Return(0, nil)
				r.On("CountCategoryParent", mock.Anything, mock.Anything).Return(1, nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, body.CategoryIsBeingUsed),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CountProductCategory", mock.Anything, mock.Anything).Return(0, nil)
				r.On("CountCategoryParent", mock.Anything, mock.Anything).Return(0, nil)
				r.On("DeleteCategory", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.DeleteCategory(context.Background(), "asd")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_GetBanner(t *testing.T) {
	var temp float64 = 10
	datestring := "02-01-2006 15:04:05"
	date, _ := time.Parse("02-01-2006 15:04:05", datestring)
	testCase := []struct {
		name        string
		body        model.Voucher
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetBanner", mock.Anything).Return([]*body.BannerResponse{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetBanner", mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetBanner(context.Background())
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_EditCategory(t *testing.T) {
	var temp float64 = 10
	datestring := "02-01-2006 15:04:05"
	date, _ := time.Parse("02-01-2006 15:04:05", datestring)
	testCase := []struct {
		name        string
		body        model.Voucher
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("EditCategory", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("EditCategory", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.EditCategory(context.Background(), body.CategoryRequest{
				ID:       "asd",
				ParentID: "asd",
				Name:     "asd",
				PhotoURL: "example.com",
				Level:    "1",
			})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_AddBanner(t *testing.T) {
	var temp float64 = 10
	datestring := "02-01-2006 15:04:05"
	date, _ := time.Parse("02-01-2006 15:04:05", datestring)
	testCase := []struct {
		name        string
		body        model.Voucher
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("AddBanner", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("AddBanner", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.AddBanner(context.Background(), body.BannerRequest{})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_DeleteBanner(t *testing.T) {
	var temp float64 = 10
	datestring := "02-01-2006 15:04:05"
	date, _ := time.Parse("02-01-2006 15:04:05", datestring)
	testCase := []struct {
		name        string
		body        model.Voucher
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("DeleteBanner", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("DeleteBanner", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.DeleteBanner(context.Background(), "123")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_EditBanner(t *testing.T) {
	var temp float64 = 10
	datestring := "02-01-2006 15:04:05"
	date, _ := time.Parse("02-01-2006 15:04:05", datestring)
	testCase := []struct {
		name        string
		body        model.Voucher
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("EditBanner", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("EditBanner", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.EditBanner(context.Background(), body.BannerIDRequest{})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_RefundOrder(t *testing.T) {
	ID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	var booll *bool
	boolltrue := true
	var temp float64 = 10
	var test *string
	var date2 sql.NullTime
	var date3 sql.NullTime

	date2.Valid = true
	date3.Valid = false

	datestring := "02-01-2025 15:04:05"
	date, _ := time.Parse("02-01-2025 15:04:05", datestring)
	date2.Time = date
	testCase := []struct {
		name        string
		body        model.Voucher
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{
					ID:             ID,
					OrderID:        ID,
					IsSellerRefund: &boolltrue,
					IsBuyerRefund:  booll,
					Reason:         "asd",
					Image:          test,
					AcceptedAt:     date3,
					RefundedAt:     date3,
					RejectedAt:     date3,
				}, nil)
				r.On("GetOrderByID", mock.Anything, mock.Anything).Return(&model.OrderModel{}, nil)
				r.On("UpdateRefund", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetOrderItemsByOrderID", mock.Anything, mock.Anything, mock.Anything).Return([]*model.OrderItem{
					{
						OrderID:         ID,
						ProductDetailID: ID,
						Quantity:        123,
						ItemPrice:       123,
						TotalPrice:      123,
						Note:            "123",
						IsReview:        true,
					},
				}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.ProductDetail{}, nil)
				r.On("UpdateProductDetailStock", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Wallet{}, nil)
				r.On("UpdateWalletBalance", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("InsertWalletHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("InsertWalletHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("123"))
			},
			expectedErr: fmt.Errorf("123"),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: httperror.New(http.StatusBadRequest, response.RefundNotFound),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{
					ID:             ID,
					OrderID:        ID,
					IsSellerRefund: booll,
					IsBuyerRefund:  booll,
					Reason:         "asd",
					Image:          test,
					AcceptedAt:     date2,
					RefundedAt:     date2,
					RejectedAt:     date2,
				}, nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, response.RefundRejected),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{
					ID:             ID,
					OrderID:        ID,
					IsSellerRefund: booll,
					IsBuyerRefund:  booll,
					Reason:         "asd",
					Image:          test,
					AcceptedAt:     date2,
					RefundedAt:     date2,
					RejectedAt:     date3,
				}, nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, response.RefundAlreadyFinished),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{}, nil)
				r.On("GetOrderByID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{}, nil)
				r.On("GetOrderByID", mock.Anything, mock.Anything).Return(&model.OrderModel{}, nil)
				r.On("UpdateRefund", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{}, nil)
				r.On("GetOrderByID", mock.Anything, mock.Anything).Return(&model.OrderModel{}, nil)
				r.On("UpdateRefund", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{}, nil)
				r.On("GetOrderByID", mock.Anything, mock.Anything).Return(&model.OrderModel{}, nil)
				r.On("UpdateRefund", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetOrderItemsByOrderID", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{}, nil)
				r.On("GetOrderByID", mock.Anything, mock.Anything).Return(&model.OrderModel{}, nil)
				r.On("UpdateRefund", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetOrderItemsByOrderID", mock.Anything, mock.Anything, mock.Anything).Return([]*model.OrderItem{
					{
						OrderID:         ID,
						ProductDetailID: ID,
						Quantity:        123,
						ItemPrice:       123,
						TotalPrice:      123,
						Note:            "123",
						IsReview:        true,
					},
				}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{}, nil)
				r.On("GetOrderByID", mock.Anything, mock.Anything).Return(&model.OrderModel{}, nil)
				r.On("UpdateRefund", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetOrderItemsByOrderID", mock.Anything, mock.Anything, mock.Anything).Return([]*model.OrderItem{
					{
						OrderID:         ID,
						ProductDetailID: ID,
						Quantity:        123,
						ItemPrice:       123,
						TotalPrice:      123,
						Note:            "123",
						IsReview:        true,
					},
				}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.ProductDetail{}, nil)
				r.On("UpdateProductDetailStock", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},

			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{}, nil)
				r.On("GetOrderByID", mock.Anything, mock.Anything).Return(&model.OrderModel{}, nil)
				r.On("UpdateRefund", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetOrderItemsByOrderID", mock.Anything, mock.Anything, mock.Anything).Return([]*model.OrderItem{
					{
						OrderID:         ID,
						ProductDetailID: ID,
						Quantity:        123,
						ItemPrice:       123,
						TotalPrice:      123,
						Note:            "123",
						IsReview:        true,
					},
				}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.ProductDetail{}, nil)
				r.On("UpdateProductDetailStock", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{}, nil)
				r.On("GetOrderByID", mock.Anything, mock.Anything).Return(&model.OrderModel{}, nil)
				r.On("UpdateRefund", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetOrderItemsByOrderID", mock.Anything, mock.Anything, mock.Anything).Return([]*model.OrderItem{
					{
						OrderID:         ID,
						ProductDetailID: ID,
						Quantity:        123,
						ItemPrice:       123,
						TotalPrice:      123,
						Note:            "123",
						IsReview:        true,
					},
				}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.ProductDetail{}, nil)
				r.On("UpdateProductDetailStock", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Wallet{}, nil)
				r.On("UpdateWalletBalance", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{}, nil)
				r.On("GetOrderByID", mock.Anything, mock.Anything).Return(&model.OrderModel{}, nil)
				r.On("UpdateRefund", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetOrderItemsByOrderID", mock.Anything, mock.Anything, mock.Anything).Return([]*model.OrderItem{
					{
						OrderID:         ID,
						ProductDetailID: ID,
						Quantity:        123,
						ItemPrice:       123,
						TotalPrice:      123,
						Note:            "123",
						IsReview:        true,
					},
				}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.ProductDetail{}, nil)
				r.On("UpdateProductDetailStock", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything, mock.Anything).Once().Return(&model.Wallet{}, nil)
				r.On("UpdateWalletBalance", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{}, nil)
				r.On("GetOrderByID", mock.Anything, mock.Anything).Return(&model.OrderModel{}, nil)
				r.On("UpdateRefund", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetOrderItemsByOrderID", mock.Anything, mock.Anything, mock.Anything).Return([]*model.OrderItem{
					{
						OrderID:         ID,
						ProductDetailID: ID,
						Quantity:        123,
						ItemPrice:       123,
						TotalPrice:      123,
						Note:            "123",
						IsReview:        true,
					},
				}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.ProductDetail{}, nil)
				r.On("UpdateProductDetailStock", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Wallet{}, nil)
				r.On("UpdateWalletBalance", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
				r.On("UpdateWalletBalance", mock.Anything, mock.Anything, mock.Anything).Once().Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{}, nil)
				r.On("GetOrderByID", mock.Anything, mock.Anything).Return(&model.OrderModel{}, nil)
				r.On("UpdateRefund", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetOrderItemsByOrderID", mock.Anything, mock.Anything, mock.Anything).Return([]*model.OrderItem{
					{
						OrderID:         ID,
						ProductDetailID: ID,
						Quantity:        123,
						ItemPrice:       123,
						TotalPrice:      123,
						Note:            "123",
						IsReview:        true,
					},
				}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.ProductDetail{}, nil)
				r.On("UpdateProductDetailStock", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Wallet{}, nil)
				r.On("UpdateWalletBalance", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("InsertWalletHistory", mock.Anything, mock.Anything, mock.Anything).Once().Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success update voucher",
			body: model.Voucher{
				Code:               "123",
				Quota:              123,
				ActivedDate:        date,
				ExpiredDate:        date,
				DiscountPercentage: &temp,
				DiscountFixPrice:   &temp,
				MinProductPrice:    &temp,
				MaxDiscountPrice:   &temp,
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetRefundByID", mock.Anything, mock.Anything).Return(&model.Refund{}, nil)
				r.On("GetOrderByID", mock.Anything, mock.Anything).Return(&model.OrderModel{}, nil)
				r.On("UpdateRefund", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetOrderItemsByOrderID", mock.Anything, mock.Anything, mock.Anything).Return([]*model.OrderItem{
					{
						OrderID:         ID,
						ProductDetailID: ID,
						Quantity:        123,
						ItemPrice:       123,
						TotalPrice:      123,
						Note:            "123",
						IsReview:        true,
					},
				}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.ProductDetail{}, nil)
				r.On("UpdateProductDetailStock", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("GetWalletByUserID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Wallet{}, nil)
				r.On("UpdateWalletBalance", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("InsertWalletHistory", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
				r.On("InsertWalletHistory", mock.Anything, mock.Anything, mock.Anything).Once().Return(fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewAdminUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.RefundOrder(context.Background(), "123")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}

}
