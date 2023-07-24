package usecase

import (
	"context"
	"database/sql"
	"errors"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/seller/delivery/body"
	"murakali/internal/module/seller/mocks"
	"murakali/pkg/httperror"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_sellerUC_GetPerformance(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		update      bool
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success get performance",
			userID: "123456",
			update: true,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("1234", nil)
				r.On("GetPerformance", mock.Anything, mock.Anything).Return(&body.SellerPerformance{}, nil)
				r.On("InsertPerformaceRedis", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetPerformance(context.Background(), tc.userID, tc.update)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}

}

func Test_sellerUC_GetAllSeller(t *testing.T) {
	value := int64(1)

	testCase := []struct {
		name        string
		shopName    string
		pgn         *pagination.Pagination
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:     "success get performance",
			shopName: "123456",
			pgn:      &pagination.Pagination{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalAllSeller", mock.Anything, mock.Anything).Return(value, nil)
				r.On("GetAllSeller", mock.Anything, mock.Anything, mock.Anything).Return([]*body.SellerResponse{}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetAllSeller(context.Background(), tc.shopName, tc.pgn)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_sellerUC_GetOrder(t *testing.T) {
	value := int64(1)

	testCase := []struct {
		name          string
		userID        string
		orderStatusID string
		voucherShopID string
		sortQuery     string
		pgn           *pagination.Pagination
		mock          func(t *testing.T, r *mocks.Repository)
		expectedErr   error
	}{
		{
			name:          "success get performance",
			userID:        "123456",
			orderStatusID: "123456",
			voucherShopID: "123456",
			sortQuery:     "123456",
			pgn:           &pagination.Pagination{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUser", mock.Anything, mock.Anything).Return("123", nil)
				r.On("GetTotalOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(value, nil)
				r.On("GetOrders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*model.Order{}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetOrder(context.Background(), tc.userID, tc.orderStatusID, tc.voucherShopID, tc.sortQuery, tc.pgn)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_sellerUC_ChangeOrderStatus(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		requestBody body.ChangeOrderStatusRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:        "success get performance",
			userID:      "123456",
			requestBody: body.ChangeOrderStatusRequest{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUser", mock.Anything, mock.Anything).Return("123", nil)
				r.On("GetShopIDByOrder", mock.Anything, mock.Anything).Return("123", nil)
				r.On("ChangeOrderStatus", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.ChangeOrderStatus(context.Background(), tc.userID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_sellerUC_GetSellerBySellerID(t *testing.T) {
	testCase := []struct {
		name        string
		sellerID    string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:     "success get performance",
			sellerID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSellerBySellerID", mock.Anything, mock.Anything).Return(&body.SellerResponse{}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetSellerBySellerID(context.Background(), tc.sellerID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_sellerUC_GetSellerByUserID(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success GetSellerByUserID",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSellerByUserID", mock.Anything, mock.Anything).Return(&body.SellerResponse{}, nil)
			},
			expectedErr: nil,
		},
		{
			name:   "failed GetSellerByUserID",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSellerByUserID", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusNotFound, body.SellerNotFoundMessage))
			},
			expectedErr: errors.New(body.SellerNotFoundMessage),
		},
		{
			name:   "failed no row GetSellerByUserID",
			userID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSellerByUserID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(body.SellerNotFoundMessage),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetSellerByUserID(context.Background(), tc.userID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_sellerUC_UpdateSellerInformationByUserID(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		shopName    string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:     "success UpdateSellerInformationByUserID",
			userID:   "123456",
			shopName: "test",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSellerByUserID", mock.Anything, mock.Anything).Return(&body.SellerResponse{}, nil)
				r.On("UpdateSellerInformationByUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			},
			expectedErr: nil,
		},
		{
			name:     "failed UpdateSellerInformationByUserID",
			userID:   "123456",
			shopName: "test",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSellerByUserID", mock.Anything, mock.Anything).Return(nil, httperror.New(http.StatusNotFound, body.SellerNotFoundMessage))
			},
			expectedErr: errors.New(body.SellerNotFoundMessage),
		},
		{
			name:     "failed no row UpdateSellerInformationByUserID",
			userID:   "123456",
			shopName: "test",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetSellerByUserID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(body.SellerNotFoundMessage),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.UpdateSellerInformationByUserID(context.Background(), tc.shopName, tc.userID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_sellerUC_DeleteCourierSellerByID(t *testing.T) {
	testCase := []struct {
		name          string
		shopCourierID string
		mock          func(t *testing.T, r *mocks.Repository)
		expectedErr   error
	}{
		{
			name:          "success UpdateSellerInformationByUserID",
			shopCourierID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCourierSellerByID", mock.Anything, mock.Anything).Return("", nil)
				r.On("DeleteCourierSellerByID", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:          "failed UpdateSellerInformationByUserID",
			shopCourierID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCourierSellerByID", mock.Anything, mock.Anything).Return("", httperror.New(http.StatusNotFound, body.SellerNotFoundMessage))
			},
			expectedErr: errors.New(body.SellerNotFoundMessage),
		},
		{
			name:          "failed no row UpdateSellerInformationByUserID",
			shopCourierID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCourierSellerByID", mock.Anything, mock.Anything).Return("", nil)
				r.On("DeleteCourierSellerByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(body.SellerNotFoundMessage),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.DeleteCourierSellerByID(context.Background(), tc.shopCourierID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_sellerUC_GetCategoryBySellerID(t *testing.T) {
	testCase := []struct {
		name        string
		shopID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:   "success UpdateSellerInformationByUserID",
			shopID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCategoryBySellerID", mock.Anything, mock.Anything).Return([]*body.CategoryResponse{}, nil)
			},
			expectedErr: nil,
		},
		{
			name:   "failed UpdateSellerInformationByUserID",
			shopID: "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCategoryBySellerID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: errors.New(body.CategoryNotFoundMessage),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetCategoryBySellerID(context.Background(), tc.shopID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_sellerUC_UpdateResiNumberInOrderSeller(t *testing.T) {
	testCase := []struct {
		name        string
		orderID     string
		userID      string
		requestBody body.UpdateNoResiOrderSellerRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:    "success UpdateResiNumberInOrderSeller",
			orderID: "123456",
			userID:  "123456",
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("", nil)
				r.On("UpdateResiNumberInOrderSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.UpdateResiNumberInOrderSeller(context.Background(), tc.userID, tc.orderID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_sellerUC_GetAllVoucherSeller(t *testing.T) {
	value := int64(1)

	testCase := []struct {
		name            string
		voucherStatusID string
		userID          string
		sortFilter      string
		pgn             *pagination.Pagination
		requestBody     body.UpdateNoResiOrderSellerRequest
		mock            func(t *testing.T, r *mocks.Repository)
		expectedErr     error
	}{
		{
			name:            "success UpdateResiNumberInOrderSeller",
			voucherStatusID: "123456",
			userID:          "123456",
			sortFilter:      "123456",
			pgn:             &pagination.Pagination{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("", nil)
				r.On("GetTotalVoucherSeller", mock.Anything, mock.Anything, mock.Anything).Return(value, nil)
				r.On("GetAllVoucherSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*model.Voucher{}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetAllVoucherSeller(context.Background(), tc.userID, tc.voucherStatusID, tc.sortFilter, tc.pgn)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

// func Test_sellerUC_CreateVoucherSeller(t *testing.T) {
// 	value := int64(1)

// 	testCase := []struct {
// 		name        string
// 		userID      string
// 		requestBody body.CreateVoucherRequest
// 		mock        func(t *testing.T, r *mocks.Repository)
// 		expectedErr error
// 	}{
// 		{
// 			name:   "success UpdateResiNumberInOrderSeller",
// 			userID: "123456",
// 			requestBody: body.CreateVoucherRequest{
// 				Code:               "123456",
// 				Quota:              1,
// 				ActivedDate:        "01-01-2021 00:00:00",
// 				ExpiredDate:        "01-01-2021 00:00:00",
// 				DiscountPercentage: 1,
// 				DiscountFixPrice:   1,
// 				MinProductPrice:    1,
// 				MaxDiscountPrice:   1,
// 			},
// 			mock: func(t *testing.T, r *mocks.Repository) {
// 				r.On("CountCodeVoucher", mock.Anything, mock.Anything).Return(value, nil)
// 				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("4cf3a332-5d81-48a0-b935-cfa83a6b6ac4", nil)
// 				r.On("CreateVoucherSeller", mock.Anything, mock.Anything).Return(nil)
// 			},
// 			expectedErr: nil,
// 		},
// 	}

// 	for _, tc := range testCase {
// 		t.Run(tc.name, func(t *testing.T) {
// 			sql, mock, _ := sqlmock.New()
// 			mock.ExpectBegin()
// 			r := mocks.NewRepository(t)
// 			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

// 			tc.mock(t, r)
// 			err := u.CreateVoucherSeller(context.Background(), tc.userID, tc.requestBody)
// 			if err != nil {
// 				assert.Equal(t, err.Error(), tc.expectedErr.Error())
// 			}
// 		})
// 	}
// }

func Test_sellerUC_UpdateVoucherSeller(t *testing.T) {
	testCase := []struct {
		name        string
		userID      string
		requestBody body.UpdateVoucherRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:        "success UpdateResiNumberInOrderSeller",
			userID:      "123456",
			requestBody: body.UpdateVoucherRequest{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("", nil)
				r.On("GetAllVoucherSellerByIDAndShopID", mock.Anything, mock.Anything).Return(&model.Voucher{}, nil)
				r.On("UpdateVoucherSeller", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.UpdateVoucherSeller(context.Background(), tc.userID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_sellerUC_GetDetailVoucherSeller(t *testing.T) {
	testCase := []struct {
		name            string
		voucherIDShopID *body.VoucherIDShopID
		mock            func(t *testing.T, r *mocks.Repository)
		expectedErr     error
	}{
		{
			name:            "success UpdateResiNumberInOrderSeller",
			voucherIDShopID: &body.VoucherIDShopID{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("", nil)
				r.On("GetAllVoucherSellerByIDAndShopID", mock.Anything, mock.Anything).Return(&model.Voucher{}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetDetailVoucherSeller(context.Background(), tc.voucherIDShopID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_sellerUC_DeleteVoucherSeller(t *testing.T) {
	testCase := []struct {
		name            string
		voucherIDShopID *body.VoucherIDShopID
		mock            func(t *testing.T, r *mocks.Repository)
		expectedErr     error
	}{
		{
			name:            "success DeleteVoucherSeller",
			voucherIDShopID: &body.VoucherIDShopID{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("", nil)
				r.On("GetAllVoucherSellerByIDAndShopID", mock.Anything, mock.Anything).Return(&model.Voucher{}, nil)
				r.On("DeleteVoucherSeller", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:            "failed DeleteVoucherSeller",
			voucherIDShopID: &body.VoucherIDShopID{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("", nil)
				r.On("GetAllVoucherSellerByIDAndShopID", mock.Anything, mock.Anything).Return(&model.Voucher{}, nil)
				r.On("DeleteVoucherSeller", mock.Anything, mock.Anything).Return(errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			name:            "failed no row DeleteVoucherSeller",
			voucherIDShopID: &body.VoucherIDShopID{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("", nil)
				r.On("GetAllVoucherSellerByIDAndShopID", mock.Anything, mock.Anything).Return(&model.Voucher{}, nil)
				r.On("DeleteVoucherSeller", mock.Anything, mock.Anything).Return(sql.ErrNoRows)
			},
			expectedErr: sql.ErrNoRows,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.DeleteVoucherSeller(context.Background(), tc.voucherIDShopID)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

// func Test_sellerUC_CancelOrderStatus(t *testing.T) {
// 	testCase := []struct {
// 		name        string
// 		userID      string
// 		requestBody body.CancelOrderStatus
// 		mock        func(t *testing.T, r *mocks.Repository)
// 		expectedErr error
// 	}{
// 		{
// 			name:   "success DeleteVoucherSeller",
// 			userID: "",
// 			mock: func(t *testing.T, r *mocks.Repository) {
// 				r.On("GetShopIDByUser", mock.Anything, mock.Anything).Return("", nil)
// 				r.On("GetOrderByOrderID", mock.Anything, mock.Anything).Return(&model.Order{}, nil)
// 				// r.On("CancelOrderStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 				// r.On("CreateRefundSeller", mock.Anything, mock.Anything, mock.Anything).Return(nil)

// 			},
// 			expectedErr: nil,
// 		},
// 	}

// 	for _, tc := range testCase {
// 		t.Run(tc.name, func(t *testing.T) {
// 			sql, mock, _ := sqlmock.New()
// 			mock.ExpectBegin()
// 			r := mocks.NewRepository(t)
// 			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

// 			tc.mock(t, r)
// 			err := u.CancelOrderStatus(context.Background(), tc.userID, tc.requestBody)
// 			if err != nil {
// 				assert.Equal(t, err.Error(), tc.expectedErr.Error())
// 			}
// 		})
// 	}
// }

func Test_sellerUC_GetAllPromotionSeller(t *testing.T) {
	value := int64(1)
	testCase := []struct {
		name            string
		userID          string
		promoStatusID   string
		pgn             *pagination.Pagination
		voucherIDShopID *body.VoucherIDShopID
		mock            func(t *testing.T, r *mocks.Repository)
		expectedErr     error
	}{
		{
			name:          "success DeleteVoucherSeller",
			userID:        "123456",
			promoStatusID: "1",
			pgn:           &pagination.Pagination{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("", nil)
				r.On("GetTotalPromotionSeller", mock.Anything, mock.Anything, mock.Anything).Return(value, nil)
				r.On("GetAllPromotionSeller", mock.Anything, mock.Anything, mock.Anything).Return([]*body.PromotionSellerResponse{}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetAllPromotionSeller(context.Background(), tc.userID, tc.promoStatusID, tc.pgn)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

// func Test_sellerUC_CreatePromotionSeller(t *testing.T) {
// 	testCase := []struct {
// 		name            string
// 		userID          string
// 		promoStatusID   string
// 		pgn             *pagination.Pagination
// 		voucherIDShopID *body.VoucherIDShopID
// 		requestBody     body.CreatePromotionRequest
// 		mock            func(t *testing.T, r *mocks.Repository)
// 		expectedErr     error
// 	}{
// 		{
// 			name:   "success DeleteVoucherSeller",
// 			userID: "123456",
// 			requestBody: body.CreatePromotionRequest{
// 				ProductPromotion: []body.ProductPromotionData{
// 					{
// 						ProductID:          "123456",
// 						Quota:              1,
// 						MaxQuantity:        1,
// 						DiscountPercentage: 1,
// 						DiscountFixPrice:   1,
// 						MinProductPrice:    1,
// 						MaxDiscountPrice:   1,
// 					},
// 				},
// 			},
// 			mock: func(t *testing.T, r *mocks.Repository) {
// 				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("123456", nil)
// 				r.On("GetProductPromotion", mock.Anything, mock.Anything).Return(&body.ProductPromotion{}, nil)
// 				r.On("CreatePromotionSeller", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 			},
// 			expectedErr: nil,
// 		},
// 	}

// 	for _, tc := range testCase {
// 		t.Run(tc.name, func(t *testing.T) {
// 			sql, mock, _ := sqlmock.New()
// 			mock.ExpectBegin()
// 			r := mocks.NewRepository(t)
// 			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

// 			tc.mock(t, r)
// 			_, err := u.CreatePromotionSeller(context.Background(), tc.userID, tc.requestBody)
// 			if err != nil {
// 				assert.Equal(t, err.Error(), tc.expectedErr.Error())
// 			}
// 		})
// 	}
// }

// func Test_sellerUC_UpdatePromotionSeller(t *testing.T) {
// 	testCase := []struct {
// 		name            string
// 		userID          string
// 		promoStatusID   string
// 		pgn             *pagination.Pagination
// 		voucherIDShopID *body.VoucherIDShopID
// 		requestBody     body.UpdatePromotionRequest
// 		mock            func(t *testing.T, r *mocks.Repository)
// 		expectedErr     error
// 	}{
// 		{
// 			name:        "success DeleteVoucherSeller",
// 			userID:      "123456",
// 			requestBody: body.UpdatePromotionRequest{},
// 			mock: func(t *testing.T, r *mocks.Repository) {
// 				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("", nil)
// 				r.On("GetPromotionSellerDetailByID", mock.Anything, mock.Anything).Return(&body.PromotionSellerResponse{}, nil)
// 				r.On("UpdatePromotionSeller", mock.Anything, mock.Anything).Return(nil)
// 			},
// 			expectedErr: nil,
// 		},
// 	}

// 	for _, tc := range testCase {
// 		t.Run(tc.name, func(t *testing.T) {
// 			sql, mock, _ := sqlmock.New()
// 			mock.ExpectBegin()
// 			r := mocks.NewRepository(t)
// 			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

// 			tc.mock(t, r)
// 			err := u.UpdatePromotionSeller(context.Background(), tc.userID, tc.requestBody)
// 			if err != nil {
// 				assert.Equal(t, err.Error(), tc.expectedErr.Error())
// 			}
// 		})
// 	}
// }

func Test_sellerUC_GetDetailPromotionSellerByID(t *testing.T) {
	float64Value := float64(1)
	testCase := []struct {
		name             string
		userID           string
		promoStatusID    string
		pgn              *pagination.Pagination
		voucherIDShopID  *body.VoucherIDShopID
		requestBody      body.UpdatePromotionRequest
		shopProductPromo *body.ShopProductPromo
		mock             func(t *testing.T, r *mocks.Repository)
		expectedErr      error
	}{
		{
			name:             "success DeleteVoucherSeller",
			shopProductPromo: &body.ShopProductPromo{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("008dc24d-1f30-4e13-823f-d62972f416df", nil)
				r.On("GetDetailPromotionSellerByID", mock.Anything, mock.Anything).Return(&body.PromotionDetailSeller{
					MinProductPrice:         &float64Value,
					MaxDiscountPrice:        &float64Value,
					DiscountPercentage:      &float64Value,
					MinPrice:                1000,
					MaxPrice:                10000,
					DiscountFixPrice:        &float64Value,
					ProductMinDiscountPrice: 1000,
					ProductSubMinPrice:      1000,
					ProductMaxDiscountPrice: 1000,
					ProductSubMaxPrice:      1000,
				}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetDetailPromotionSellerByID(context.Background(), tc.shopProductPromo)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_sellerUC_GetProductWithoutPromotionSeller(t *testing.T) {
	value := int64(1)
	testCase := []struct {
		name             string
		userID           string
		promoStatusID    string
		pgn              *pagination.Pagination
		voucherIDShopID  *body.VoucherIDShopID
		requestBody      body.UpdatePromotionRequest
		shopProductPromo *body.ShopProductPromo
		productName      string
		mock             func(t *testing.T, r *mocks.Repository)
		expectedErr      error
	}{
		{
			name:        "success DeleteVoucherSeller",
			userID:      "123456",
			productName: "test",
			pgn:         &pagination.Pagination{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("", nil)
				r.On("GetTotalProductWithoutPromotionSeller", mock.Anything, mock.Anything, mock.Anything).Return(value, nil)
				r.On("GetProductWithoutPromotionSeller", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*body.GetProductWithoutPromotion{}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetProductWithoutPromotionSeller(context.Background(), tc.userID, tc.productName, tc.pgn)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func Test_sellerUC_GetRefundOrderSeller(t *testing.T) {
	uuid, _ := uuid.Parse("008dc24d-1f30-4e13-823f-d62972f416df")
	bollTrue := true
	bollFalse := false

	testCase := []struct {
		name             string
		userID           string
		promoStatusID    string
		pgn              *pagination.Pagination
		voucherIDShopID  *body.VoucherIDShopID
		requestBody      body.UpdatePromotionRequest
		shopProductPromo *body.ShopProductPromo
		productName      string
		orderID          string
		mock             func(t *testing.T, r *mocks.Repository)
		expectedErr      error
	}{
		{
			name:    "success GetRefundOrderSeller",
			userID:  "008dc24d-1f30-4e13-823f-d62972f416df",
			orderID: "008dc24d-1f30-4e13-823f-d62972f416df",
			pgn:     &pagination.Pagination{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("008dc24d-1f30-4e13-823f-d62972f416df", nil)
				r.On("GetOrderModelByID", mock.Anything, mock.Anything).Return(&model.OrderModel{ShopID: uuid}, nil)
				r.On("GetRefundOrderByOrderID", mock.Anything, mock.Anything).Return(&model.Refund{IsBuyerRefund: &bollTrue, IsSellerRefund: &bollFalse}, nil)
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetRefundThreadByRefundID", mock.Anything, mock.Anything).Return([]*body.RThread{}, nil)
			},
			expectedErr: nil,
		},
		{
			name:    "success GetRefundOrderSeller",
			userID:  "008dc24d-1f30-4e13-823f-d62972f416df",
			orderID: "008dc24d-1f30-4e13-823f-d62972f416df",
			pgn:     &pagination.Pagination{},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("008dc24d-1f30-4e13-823f-d62972f416df", nil)
				r.On("GetOrderModelByID", mock.Anything, mock.Anything).Return(&model.OrderModel{ShopID: uuid}, nil)
				r.On("GetRefundOrderByOrderID", mock.Anything, mock.Anything).Return(&model.Refund{IsBuyerRefund: &bollFalse, IsSellerRefund: &bollTrue}, nil)
				r.On("GetShopByID", mock.Anything, mock.Anything).Return(&model.Shop{}, nil)
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetRefundThreadByRefundID", mock.Anything, mock.Anything).Return([]*body.RThread{}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			_, err := u.GetRefundOrderSeller(context.Background(), tc.userID, tc.orderID)
			if err != nil {
				assert.Equal(t, err, tc.expectedErr)
			}
		})
	}
}

func Test_sellerUC_CreateRefundThreadSeller(t *testing.T) {
	uuidString1, _ := uuid.Parse("008dc24d-1f30-4e13-823f-d62972f416df")
	// uuidString2, _ := uuid.Parse("008dc24d-1f30-4e13-823f-d62972f416de")
	testCase := []struct {
		name             string
		userID           string
		promoStatusID    string
		pgn              *pagination.Pagination
		voucherIDShopID  *body.VoucherIDShopID
		shopProductPromo *body.ShopProductPromo
		productName      string
		orderID          string
		requestBody      *body.CreateRefundThreadRequest
		mock             func(t *testing.T, r *mocks.Repository)
		expectedErr      error
	}{
		{
			name:   "success GetRefundOrderSeller",
			userID: "008dc24d-1f30-4e13-823f-d62972f416df",
			requestBody: &body.CreateRefundThreadRequest{
				RefundID: "008dc24d-1f30-4e13-823f-d62972f416df",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("008dc24d-1f30-4e13-823f-d62972f416df", nil)
				r.On("GetRefundOrderByID", mock.Anything, mock.Anything).Return(&model.Refund{OrderID: uuidString1}, nil)
				r.On("GetOrderModelByID", mock.Anything, mock.Anything).Return(&model.OrderModel{ShopID: uuidString1, OrderStatusID: 6}, nil)
				r.On("CreateRefundThreadSeller", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.CreateRefundThreadSeller(context.Background(), tc.userID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err, tc.expectedErr)
			}
		})
	}
}

func Test_sellerUC_UpdateRefundAccept(t *testing.T) {
	uuidString1, _ := uuid.Parse("008dc24d-1f30-4e13-823f-d62972f416df")
	// uuidString2, _ := uuid.Parse("008dc24d-1f30-4e13-823f-d62972f416de")
	testCase := []struct {
		name             string
		userID           string
		promoStatusID    string
		pgn              *pagination.Pagination
		voucherIDShopID  *body.VoucherIDShopID
		shopProductPromo *body.ShopProductPromo
		productName      string
		orderID          string
		requestBody      *body.UpdateRefundRequest
		mock             func(t *testing.T, r *mocks.Repository)
		expectedErr      error
	}{
		{
			name:   "success GetRefundOrderSeller",
			userID: "008dc24d-1f30-4e13-823f-d62972f416df",
			requestBody: &body.UpdateRefundRequest{
				RefundID: "008dc24d-1f30-4e13-823f-d62972f416df",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("008dc24d-1f30-4e13-823f-d62972f416df", nil)
				r.On("GetRefundOrderByID", mock.Anything, mock.Anything).Return(&model.Refund{OrderID: uuidString1}, nil)
				r.On("GetOrderModelByID", mock.Anything, mock.Anything).Return(&model.OrderModel{ShopID: uuidString1, OrderStatusID: 6}, nil)
				r.On("UpdateRefundAccept", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.UpdateRefundAccept(context.Background(), tc.userID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err, tc.expectedErr)
			}
		})
	}
}

func Test_sellerUC_UpdateRefundReject(t *testing.T) {
	uuidString1, _ := uuid.Parse("008dc24d-1f30-4e13-823f-d62972f416df")
	// uuidString2, _ := uuid.Parse("008dc24d-1f30-4e13-823f-d62972f416de")
	testCase := []struct {
		name             string
		userID           string
		promoStatusID    string
		pgn              *pagination.Pagination
		voucherIDShopID  *body.VoucherIDShopID
		shopProductPromo *body.ShopProductPromo
		productName      string
		orderID          string
		requestBody      *body.UpdateRefundRequest
		mock             func(t *testing.T, r *mocks.Repository)
		expectedErr      error
	}{
		{
			name:   "success GetRefundOrderSeller",
			userID: "008dc24d-1f30-4e13-823f-d62972f416df",
			requestBody: &body.UpdateRefundRequest{
				RefundID: "008dc24d-1f30-4e13-823f-d62972f416df",
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("008dc24d-1f30-4e13-823f-d62972f416df", nil)
				r.On("GetRefundOrderByID", mock.Anything, mock.Anything).Return(&model.Refund{OrderID: uuidString1}, nil)
				r.On("GetOrderModelByID", mock.Anything, mock.Anything).Return(&model.OrderModel{ShopID: uuidString1, OrderStatusID: 6}, nil)
				r.On("UpdateRefundReject", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateOrderRefundRejected", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewSellerUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.UpdateRefundReject(context.Background(), tc.userID, tc.requestBody)
			if err != nil {
				assert.Equal(t, err, tc.expectedErr)
			}
		})
	}
}
