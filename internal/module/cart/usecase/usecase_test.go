package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/config"
	"net/http"

	"murakali/internal/model"
	"murakali/internal/module/cart/delivery/body"
	"murakali/internal/module/cart/mocks"
	"murakali/pkg/httperror"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCartUseCase_GetCartHoverHome(t *testing.T) {
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{

		{
			name: "success  get cart hover home",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				var temp float64 = 10
				var tempInt int = 10
				r.On("GetCartHoverHome", mock.Anything, mock.Anything, mock.Anything).Return([]*body.CartHome{
					{
						MaxDiscountPrice:   &temp,
						MinProductPrice:    &temp,
						Quota:              &tempInt,
						DiscountFixPrice:   &temp,
						ResultDiscount:     temp,
						Price:              temp,
						DiscountPercentage: &temp,
					},
				}, nil)
				r.On("GetTotalCart", mock.Anything, mock.Anything).Return(int64(1), nil)
			},
			expectedErr: nil,
		},

		{
			name: "get cart hover error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCartHoverHome", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "get total cart hover home error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCartHoverHome", mock.Anything, mock.Anything, mock.Anything).Return([]*body.CartHome{}, nil)
				r.On("GetTotalCart", mock.Anything, mock.Anything).Return(int64(0), fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewCartUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetCartHoverHome(context.Background(), "123456", 10)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestCartUseCase_GetCartItems(t *testing.T) {
	var temp float64 = 10
	var tempInt int = 10
	id, _ := uuid.Parse("989d94b7-58fc-4a76-ae01-1c1b47a0755c")
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success  get cart items",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {

				r.On("GetTotalCart", mock.Anything, mock.Anything).Return(int64(1), nil)
				r.On("GetCartItems", mock.Anything, mock.Anything, mock.Anything).Return([]*body.CartItemsResponse{{
					Shop: &body.ShopResponse{
						Name: "ASdsadsa",
						ID:   id,
					},
					Weight: 10,
					ProductDetails: []*body.ProductDetailResponse{
						{
							Title:        "test",
							ThumbnailURL: "asdassd",
							ProductPrice: temp,
							ProductStock: temp,
							Quantity:     10,
							Weight:       float64(tempInt),
							Promo: &body.PromoResponse{
								DiscountPercentage: &temp,
								DiscountFixPrice:   &temp,
								MinProductPrice:    &temp,
								MaxDiscountPrice:   &temp,
								ResultDiscount:     temp,
								SubPrice:           temp,
								Quota:              &tempInt,
							},
							Variant: map[string]string{
								"test": "test",
							},
						},
					},
				}}, []*body.ProductDetailResponse{{
					ID:           "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
					Title:        "Test",
					ThumbnailURL: "test",
					ProductPrice: temp,
					ProductStock: temp,
					Quantity:     temp,
					Weight:       temp,
					Variant: map[string]string{
						"test": "test",
					},
					Promo: &body.PromoResponse{
						DiscountPercentage: &temp,
						DiscountFixPrice:   &temp,
						MinProductPrice:    &temp,
						MaxDiscountPrice:   &temp,
						ResultDiscount:     temp,
						SubPrice:           temp,
						Quota:              &tempInt,
					},
				}}, []*body.PromoResponse{{
					DiscountPercentage: &temp,
					DiscountFixPrice:   &temp,
					MinProductPrice:    &temp,
					MaxDiscountPrice:   &temp,
					ResultDiscount:     temp,
					Quota:              &tempInt,
					SubPrice:           temp}}, nil)

			},
			expectedErr: nil,
		},
		{
			name: "get cart items error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalCart", mock.Anything, mock.Anything).Return(int64(0), fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "get cart items error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalCart", mock.Anything, mock.Anything).Return(int64(0), nil)
				r.On("GetCartItems", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewCartUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetCartItems(context.Background(), "123456", &pagination.Pagination{})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestCartUseCase_AddCartItems(t *testing.T) {
	passwordHash := "$2a$10$WKul/6gjYoYjOXuNVX4XGen1ZkWYb1PKFiI5vlZp5TFerZh6nTujG"
	testCase := []struct {
		name        string
		body        body.AddCartItemRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success add cart",
			body: body.AddCartItemRequest{ProductDetailID: "123456",
				Quantity: 1},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(&model.ProductDetail{Stock: 2}, nil)
				r.On("GetCartProductDetail", mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				r.On("CreateCart", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&model.CartItem{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "error get user by id",
			body: body.AddCartItemRequest{ProductDetailID: "123456",
				Quantity: 5},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.UserNotExistMessage),
		},
		{
			name: "error get product detail no rows",
			body: body.AddCartItemRequest{ProductDetailID: "123456",
				Quantity: 5},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.ProductDetailNotExistMessage),
		},
		{
			name: "error add cart no rows",
			body: body.AddCartItemRequest{ProductDetailID: "123456",
				Quantity: 5},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(&model.ProductDetail{Stock: 2}, nil)
				r.On("GetCartProductDetail", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error add cart quantity",
			body: body.AddCartItemRequest{ProductDetailID: "123456",
				Quantity: 5},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(&model.ProductDetail{Stock: 2}, nil)
				r.On("GetCartProductDetail", mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
			expectedErr: httperror.New(http.StatusBadRequest, response.QuantityReachedMaximum),
		},
		{
			name: "error add cart",
			body: body.AddCartItemRequest{ProductDetailID: "123456",
				Quantity: 1},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(&model.ProductDetail{Stock: 2}, nil)
				r.On("GetCartProductDetail", mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				r.On("CreateCart", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success add quantity cart",
			body: body.AddCartItemRequest{ProductDetailID: "123456",
				Quantity: 1},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(&model.ProductDetail{Stock: 2}, nil)
				r.On("GetCartProductDetail", mock.Anything, mock.Anything, mock.Anything).Return(&model.CartItem{}, nil)
				r.On("UpdateCartByID", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "error add quantity cart  quantity",
			body: body.AddCartItemRequest{ProductDetailID: "123456",
				Quantity: 5},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(&model.ProductDetail{Stock: 2}, nil)
				r.On("GetCartProductDetail", mock.Anything, mock.Anything, mock.Anything).Return(&model.CartItem{}, nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, response.QuantityReachedMaximum),
		},
		{
			name: "error add quantity cart update",
			body: body.AddCartItemRequest{ProductDetailID: "123456",
				Quantity: 1},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(&model.ProductDetail{Stock: 2}, nil)
				r.On("GetCartProductDetail", mock.Anything, mock.Anything, mock.Anything).Return(&model.CartItem{}, nil)
				r.On("UpdateCartByID", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewCartUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			err := u.AddCartItems(context.Background(), "123456", tc.body)
			assert.Equal(t, err, tc.expectedErr)
		})
	}
}

func TestCartUseCase_UpdateCartItems(t *testing.T) {
	passwordHash := "$2a$10$WKul/6gjYoYjOXuNVX4XGen1ZkWYb1PKFiI5vlZp5TFerZh6nTujG"
	testCase := []struct {
		name        string
		body        body.CartItemRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "error get user by id",
			body: body.CartItemRequest{ProductDetailID: "123456",
				Quantity: 5},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.UserNotExistMessage),
		},
		{
			name: "error get product detail no rows",
			body: body.CartItemRequest{ProductDetailID: "123456",
				Quantity: 5},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.ProductDetailNotExistMessage),
		},
		{
			name: "error add cart no rows",
			body: body.CartItemRequest{ProductDetailID: "123456",
				Quantity: 5},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(&model.ProductDetail{Stock: 2}, nil)
				r.On("GetCartProductDetail", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},

		{
			name: "success add quantity cart",
			body: body.CartItemRequest{ProductDetailID: "123456",
				Quantity: 1},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(&model.ProductDetail{Stock: 2}, nil)
				r.On("GetCartProductDetail", mock.Anything, mock.Anything, mock.Anything).Return(&model.CartItem{}, nil)
				r.On("UpdateCartByID", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "error add quantity cart  quantity",
			body: body.CartItemRequest{ProductDetailID: "123456",
				Quantity: 5},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(&model.ProductDetail{Stock: 2}, nil)
				r.On("GetCartProductDetail", mock.Anything, mock.Anything, mock.Anything).Return(&model.CartItem{}, nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, response.QuantityReachedMaximum),
		},
		{
			name: "error add quantity cart update",
			body: body.CartItemRequest{ProductDetailID: "123456",
				Quantity: 1},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Email: "a@test.com", Password: &passwordHash}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(&model.ProductDetail{Stock: 2}, nil)
				r.On("GetCartProductDetail", mock.Anything, mock.Anything, mock.Anything).Return(&model.CartItem{}, nil)
				r.On("UpdateCartByID", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewCartUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			err := u.UpdateCartItems(context.Background(), "123456", tc.body)
			assert.Equal(t, err, tc.expectedErr)
		})
	}
}

func TestCartUseCase_GetDeleteCart(t *testing.T) {
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success  delete cart",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(&model.ProductDetail{}, nil)
				r.On("GetCartProductDetail", mock.Anything, mock.Anything, mock.Anything).Return(&model.CartItem{}, nil)
				r.On("DeleteCartByID", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "error get user by id",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.UserNotExistMessage),
		},
		{
			name: "error get product detail by id",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.ProductDetailNotExistMessage),
		},
		{
			name: "error get cart product detail",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(&model.ProductDetail{}, nil)
				r.On("GetCartProductDetail", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error delete cart",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{}, nil)
				r.On("GetProductDetailByID", mock.Anything, mock.Anything).Return(&model.ProductDetail{}, nil)
				r.On("GetCartProductDetail", mock.Anything, mock.Anything, mock.Anything).Return(&model.CartItem{}, nil)
				r.On("DeleteCartByID", mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewCartUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			err := u.DeleteCartItems(context.Background(), "123456", "123456")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestCartUseCase_GetVoucherShop(t *testing.T) {
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{

		{
			name: "success get voucher shop",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalVoucherShop", mock.Anything, mock.Anything).Return(int64(1), nil)
				r.On("GetVoucherShop", mock.Anything, mock.Anything, mock.Anything).Return([]*model.Voucher{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "error count voucher shop",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalVoucherShop", mock.Anything, mock.Anything).Return(int64(1), fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error count voucher shop",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalVoucherShop", mock.Anything, mock.Anything).Return(int64(1), nil)
				r.On("GetVoucherShop", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewCartUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetVoucherShop(context.Background(), "123456", &pagination.Pagination{})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestCartUseCase_GetVoucherMarketplace(t *testing.T) {
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{

		{
			name: "success get voucher marketplace",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalVoucherMarketplace", mock.Anything, mock.Anything).Return(int64(1), nil)
				r.On("GetVoucherMarketplace", mock.Anything, mock.Anything, mock.Anything).Return([]*model.Voucher{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "error count voucher marketplace",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalVoucherMarketplace", mock.Anything, mock.Anything).Return(int64(1), fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error count voucher marketplace",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalVoucherMarketplace", mock.Anything, mock.Anything).Return(int64(1), nil)
				r.On("GetVoucherMarketplace", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewCartUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetVoucherMarketplace(context.Background(), &pagination.Pagination{})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
