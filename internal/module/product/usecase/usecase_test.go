package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/product/delivery/body"
	"murakali/internal/module/product/mocks"
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

func TestProductUseCase_GetCategories(t *testing.T) {
	dateString := "2021-11-23"
	date, _ := time.Parse("2006-01-02", dateString)
	id, _ := uuid.Parse("989d94b7-58fc-4a76-ae01-1c1b47a0755c")
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success  get categories",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {

				r.On("GetCategories", mock.Anything).Return(
					[]*model.Category{{ID: id, ParentID: id, Name: "test", PhotoURL: "test",
						CreatedAt: date,
					}}, nil)
				r.On("GetCategoriesByParentID", mock.Anything, mock.Anything).Once().Return(
					[]*model.Category{
						{
							ID:        id,
							ParentID:  id,
							Name:      "test",
							PhotoURL:  "test",
							CreatedAt: date,
						},
					}, nil)

				r.On("GetCategoriesByParentID", mock.Anything, mock.Anything).Once().Return(
					[]*model.Category{}, nil)

			},
			expectedErr: nil,
		},
		{
			name: "error get categories no rows",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCategories", mock.Anything).Return(
					nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error get categories child no rows",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCategories", mock.Anything).Return(
					[]*model.Category{{ID: id, ParentID: id, Name: "test", PhotoURL: "test",
						CreatedAt: date,
					}}, nil)
				r.On("GetCategoriesByParentID", mock.Anything, mock.Anything).Once().Return(
					nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error get categories child no rows",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCategories", mock.Anything).Return(
					[]*model.Category{{ID: id, ParentID: id, Name: "test", PhotoURL: "test",
						CreatedAt: date,
					}}, nil)
				r.On("GetCategoriesByParentID", mock.Anything, mock.Anything).Once().Return(
					[]*model.Category{
						{
							ID:        id,
							ParentID:  id,
							Name:      "test",
							PhotoURL:  "test",
							CreatedAt: date,
						},
					}, nil)

				r.On("GetCategoriesByParentID", mock.Anything, mock.Anything).Once().Return(
					nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetCategories(context.Background())
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestProductUseCase_GetBanner(t *testing.T) {

	id, _ := uuid.Parse("989d94b7-58fc-4a76-ae01-1c1b47a0755c")
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{

		{
			name: "success  get banner home",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetBanners", mock.Anything).Return([]*model.Banner{{ID: id,
					Title:   "test",
					Content: "test", ImageURL: "test", PageURL: "test", IsActive: true}}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "get banner error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetBanners", mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetBanners(context.Background())
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestProductUseCase_GetCategoriesByName(t *testing.T) {
	dateString := "2021-11-23"
	date, _ := time.Parse("2006-01-02", dateString)
	id, _ := uuid.Parse("989d94b7-58fc-4a76-ae01-1c1b47a0755c")
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{

		{
			name: "success  get categories by name",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCategoriesByName", mock.Anything, mock.Anything).
					Return([]*model.Category{
						{ID: id, ParentID: id, Name: "test", PhotoURL: "test",
							CreatedAt: date,
						}}, nil)
				r.On("GetCategoriesByParentID", mock.Anything, mock.Anything).Once().Return(
					[]*model.Category{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "get categories by name error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCategoriesByName", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "get categories child by name error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetCategoriesByName", mock.Anything, mock.Anything).
					Return([]*model.Category{
						{ID: id, ParentID: id, Name: "test", PhotoURL: "test",
							CreatedAt: date,
						}}, nil)
				r.On("GetCategoriesByParentID", mock.Anything, mock.Anything).Once().Return(
					nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetCategoriesByName(context.Background(), "test")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestProductUseCase_GetRecommendedProducts(t *testing.T) {
	var temp float64 = 10
	dateString := "2021-11-23"
	date, _ := time.Parse("2006-01-02", dateString)
	id, _ := uuid.Parse("989d94b7-58fc-4a76-ae01-1c1b47a0755c")
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{

		{
			name: "success  get reccomended product",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalProduct", mock.Anything).
					Return(int64(1), nil)
				r.On("GetRecommendedProducts", mock.Anything, mock.Anything).
					Return([]*body.Products{
						{
							ID:                        id,
							Title:                     "test",
							UnitSold:                  10,
							RatingAVG:                 temp,
							ThumbnailURL:              "test",
							MinPrice:                  temp,
							MaxPrice:                  temp,
							ViewCount:                 10,
							SubPrice:                  temp,
							PromoDiscountPercentage:   &temp,
							PromoDiscountFixPrice:     &temp,
							PromoMaxDiscountPrice:     &temp,
							ResultDiscount:            &temp,
							VoucherDiscountPercentage: &temp,
							VoucherDiscountFixPrice:   &temp,
						}},
						[]*model.Promotion{{
							ID:                 id,
							Name:               "test",
							ProductID:          id,
							DiscountPercentage: &temp,
							DiscountFixPrice:   &temp,
							MinProductPrice:    &temp,
							MaxDiscountPrice:   &temp,
						}}, []*model.Voucher{{ID: id, ShopID: id, Code: "test",
							Quota: 10, ActivedDate: date, ExpiredDate: date}}, nil)

			},
			expectedErr: nil,
		},
		{
			name: "get  get reccomended product error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalProduct", mock.Anything).
					Return(int64(0), fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success  get reccomended product",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalProduct", mock.Anything).
					Return(int64(1), nil)
				r.On("GetRecommendedProducts", mock.Anything, mock.Anything).
					Return(nil, nil, nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetRecommendedProducts(context.Background(), &pagination.Pagination{})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestProductUseCase_GetProductDetail(t *testing.T) {
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{

		{
			name: "success  get product detail",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProductInfo", mock.Anything, mock.Anything).
					Return(&body.ProductInfo{}, nil)

				r.On("GetPromotionInfo", mock.Anything, mock.Anything).
					Return(&body.PromotionInfo{}, nil)
				r.On("GetProductDetail", mock.Anything, mock.Anything, mock.Anything).
					Return([]*body.ProductDetail{{
						ProductDetailID: "989d94b7-58fc-4a76-ae01-1c1b47a0755c"}}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "get product detail error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProductInfo", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "get product detail error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProductInfo", mock.Anything, mock.Anything).
					Return(&body.ProductInfo{}, nil)
				r.On("GetPromotionInfo", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "get product detail error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProductInfo", mock.Anything, mock.Anything).
					Return(&body.ProductInfo{}, nil)

				r.On("GetPromotionInfo", mock.Anything, mock.Anything).
					Return(&body.PromotionInfo{}, nil)
				r.On("GetProductDetail", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetProductDetail(context.Background(), "989d94b7-58fc-4a76-ae01-1c1b47a0755c")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestProductUseCase_GetProducts(t *testing.T) {
	var temp float64 = 10
	dateString := "2021-11-23"
	date, _ := time.Parse("2006-01-02", dateString)
	id, _ := uuid.Parse("989d94b7-58fc-4a76-ae01-1c1b47a0755c")
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{

		{
			name: "success  get  product",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetAllTotalProduct", mock.Anything, mock.Anything).
					Return(int64(1), nil)
				r.On("GetProducts", mock.Anything, mock.Anything, mock.Anything).
					Return([]*body.Products{
						{
							ID:                        id,
							Title:                     "test",
							UnitSold:                  10,
							RatingAVG:                 temp,
							ThumbnailURL:              "test",
							MinPrice:                  temp,
							MaxPrice:                  temp,
							ViewCount:                 10,
							SubPrice:                  temp,
							PromoDiscountPercentage:   &temp,
							PromoDiscountFixPrice:     &temp,
							PromoMaxDiscountPrice:     &temp,
							ResultDiscount:            &temp,
							VoucherDiscountPercentage: &temp,
							VoucherDiscountFixPrice:   &temp,
						}},
						[]*model.Promotion{{
							ID:                 id,
							Name:               "test",
							ProductID:          id,
							DiscountPercentage: &temp,
							DiscountFixPrice:   &temp,
							MinProductPrice:    &temp,
							MaxDiscountPrice:   &temp,
						}}, []*model.Voucher{{ID: id, ShopID: id, Code: "test",
							Quota: 10, ActivedDate: date, ExpiredDate: date}}, nil)

			},
			expectedErr: nil,
		},
		{
			name: "get  get  product error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetAllTotalProduct", mock.Anything, mock.Anything).
					Return(int64(0), fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success  get  product",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetAllTotalProduct", mock.Anything, mock.Anything).
					Return(int64(1), nil)
				r.On("GetProducts", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, nil, nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetProducts(context.Background(), &pagination.Pagination{}, &body.GetProductQueryRequest{})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestProductUseCase_GetAllProductImage(t *testing.T) {
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{

		{
			name: "success  get  product image",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProductInfo", mock.Anything, mock.Anything).
					Return(&body.ProductInfo{}, nil)
				r.On("GetProductDetail", mock.Anything, mock.Anything, mock.Anything).
					Return([]*body.ProductDetail{
						{
							ProductDetailID: "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
						}},
						nil)
				r.On("GetAllImageByProductDetailID", mock.Anything, mock.Anything, mock.Anything).
					Return([]*string{},
						nil)

			},
			expectedErr: nil,
		},
		{
			name: "get  product image info error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProductInfo", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "get  product image detail error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProductInfo", mock.Anything, mock.Anything).
					Return(&body.ProductInfo{}, nil)
				r.On("GetProductDetail", mock.Anything, mock.Anything, mock.Anything).
					Return(
						nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "get  product image  error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProductInfo", mock.Anything, mock.Anything).
					Return(&body.ProductInfo{}, nil)
				r.On("GetProductDetail", mock.Anything, mock.Anything, mock.Anything).
					Return([]*body.ProductDetail{
						{
							ProductDetailID: "989d94b7-58fc-4a76-ae01-1c1b47a0755c",
						}},
						nil)
				r.On("GetAllImageByProductDetailID", mock.Anything, mock.Anything, mock.Anything).
					Return(
						nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetAllProductImage(context.Background(), "989d94b7-58fc-4a76-ae01-1c1b47a0755c")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestProductUseCase_GetFavoriteProducts(t *testing.T) {
	var temp float64 = 10
	dateString := "2021-11-23"
	date, _ := time.Parse("2006-01-02", dateString)
	id, _ := uuid.Parse("989d94b7-58fc-4a76-ae01-1c1b47a0755c")
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{

		{
			name: "success  get  product",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetAllFavoriteTotalProduct", mock.Anything, mock.Anything, mock.Anything).
					Return(int64(1), nil)
				r.On("GetFavoriteProducts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return([]*body.Products{
						{
							ID:                        id,
							Title:                     "test",
							UnitSold:                  10,
							RatingAVG:                 temp,
							ThumbnailURL:              "test",
							MinPrice:                  temp,
							MaxPrice:                  temp,
							ViewCount:                 10,
							SubPrice:                  temp,
							PromoDiscountPercentage:   &temp,
							PromoDiscountFixPrice:     &temp,
							PromoMaxDiscountPrice:     &temp,
							ResultDiscount:            &temp,
							VoucherDiscountPercentage: &temp,
							VoucherDiscountFixPrice:   &temp,
						}},
						[]*model.Promotion{{
							ID:                 id,
							Name:               "test",
							ProductID:          id,
							DiscountPercentage: &temp,
							DiscountFixPrice:   &temp,
							MinProductPrice:    &temp,
							MaxDiscountPrice:   &temp,
						}}, []*model.Voucher{{ID: id, ShopID: id, Code: "test",
							Quota: 10, ActivedDate: date, ExpiredDate: date}}, nil)

			},
			expectedErr: nil,
		},
		{
			name: "get  get  product error",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetAllFavoriteTotalProduct", mock.Anything, mock.Anything, mock.Anything).
					Return(int64(0), fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success  get  product",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetAllFavoriteTotalProduct", mock.Anything, mock.Anything, mock.Anything).
					Return(int64(1), nil)
				r.On("GetFavoriteProducts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, nil, nil, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetFavoriteProducts(context.Background(), &pagination.Pagination{}, &body.GetProductQueryRequest{}, "989d94b7-58fc-4a76-ae01-1c1b47a0755c")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestProductUseCase_CountSpecificFavoriteProduct(t *testing.T) {

	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{

		{
			name: "success count specific favorite product",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CountSpecificFavoriteProduct", mock.Anything, mock.Anything).
					Return(int64(1), nil)

			},
			expectedErr: nil,
		},
		{
			name: "error count specific favorite product",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("CountSpecificFavoriteProduct", mock.Anything, mock.Anything).
					Return(int64(0), fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.CountSpecificFavoriteProduct(context.Background(), "989d94b7-58fc-4a76-ae01-1c1b47a0755c")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestCartUseCase_CreateFavoriteProduct(t *testing.T) {
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success create favorite",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProductInfo", mock.Anything, mock.Anything).Return(&body.ProductInfo{}, nil)
				r.On("FindFavoriteProduct", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
				r.On("CreateFavoriteProduct", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			},
			expectedErr: nil,
		},
		{
			name: "error get product info",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProductInfo", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.ProductNotExistMessage),
		},
		{
			name: "error find favorite product",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProductInfo", mock.Anything, mock.Anything).Return(&body.ProductInfo{}, nil)
				r.On("FindFavoriteProduct", mock.Anything, mock.Anything, mock.Anything).Return(true, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "success create favorite",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetProductInfo", mock.Anything, mock.Anything).Return(&body.ProductInfo{}, nil)
				r.On("FindFavoriteProduct", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
				r.On("CreateFavoriteProduct", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.CreateFavoriteProduct(context.Background(), "123456", "123456")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestCartUseCase_DeleteFavoriteProduct(t *testing.T) {
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success delete favorite",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("FindFavoriteProduct", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
				r.On("DeleteFavoriteProduct", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			},
			expectedErr: nil,
		},
		{
			name: "error find favorite product",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("FindFavoriteProduct", mock.Anything, mock.Anything, mock.Anything).Return(false, fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error find favorite product 2",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("FindFavoriteProduct", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.ProductNotExistMessage),
		},
		{
			name: "error find favorite product 2",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("FindFavoriteProduct", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
				r.On("DeleteFavoriteProduct", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.DeleteFavoriteProduct(context.Background(), "123456", "123456")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestProductUseCase_GetProductReviews(t *testing.T) {

	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{

		{
			name: "success  get product reviews",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalAllReviewProduct", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)
				r.On("GetProductReviews", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*body.ReviewProduct{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "error get product reviews",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalAllReviewProduct", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "error get product reviews",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalAllReviewProduct", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)
				r.On("GetProductReviews", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetProductReviews(context.Background(), &pagination.Pagination{}, "123456", &body.GetReviewQueryRequest{})
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestProductUseCase_GetTotalReviewRatingByProductID(t *testing.T) {
	testCase := []struct {
		name        string
		body        interface{}
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{

		{
			name: "success  get product reviews",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalReviewRatingByProductID", mock.Anything, mock.Anything).Return([]*body.RatingProduct{{}, {}}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "success  get product reviews",
			body: nil,
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetTotalReviewRatingByProductID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		// {
		// 	name: "error get product reviews",
		// 	body: nil,
		// 	mock: func(t *testing.T, r *mocks.Repository) {
		// 		r.On("GetTotalAllReviewProduct", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), fmt.Errorf("test"))
		// 	},
		// 	expectedErr: fmt.Errorf("test"),
		// },
		// {
		// 	name: "error get product reviews",
		// 	body: nil,
		// 	mock: func(t *testing.T, r *mocks.Repository) {
		// 		r.On("GetTotalAllReviewProduct", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)
		// 		r.On("GetProductReviews", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
		// 	},
		// 	expectedErr: fmt.Errorf("test"),
		// },
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{}, r)

			tc.mock(t, r)
			_, err := u.GetTotalReviewRatingByProductID(context.Background(), "123456")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestCartUseCase_CreateProduct(t *testing.T) {
	var temp float64 = 10
	testCase := []struct {
		name        string
		body        body.CreateProductRequest
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name: "success create favorite",
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
					Photo:     []string{"test"},
					VariantDetail: []body.VariantDetailRequest{{
						Type: "color",
						Name: "big",
					}},
				}},
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("123456", nil)
				r.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return("123456", nil)
				r.On("CreateProductDetail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("123456", nil)
				r.On("CreatePhoto", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("CreateVariantDetail", mock.Anything, mock.Anything, mock.Anything).Return("123456", nil)
				r.On("CreateVariant", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			},
			expectedErr: nil,
		},
		{
			name: "success create favorite",
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
					Photo:     []string{"test"},
					VariantDetail: []body.VariantDetailRequest{{
						Type: "color",
						Name: "big",
					}},
				}},
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("", sql.ErrNoRows)

			},
			expectedErr: httperror.New(http.StatusBadRequest, response.UserNotExistMessage),
		},
		{
			name: "success create favorite",
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
					Photo:     []string{"test"},
					VariantDetail: []body.VariantDetailRequest{{
						Type: "color",
						Name: "big",
					}},
				}},
			},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetShopIDByUserID", mock.Anything, mock.Anything).Return("123456", nil)
				r.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return("", fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.CreateProduct(context.Background(), tc.body, "123456")
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
			name: "Delete Product successfully",
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
				r.On("FindReview", mock.Anything, mock.Anything).Return(&body.ReviewProduct{}, nil)
				r.On("DeleteReview", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Delete Product successfully",
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
				r.On("FindReview", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name: "Delete Product successfully",
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
				r.On("FindReview", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, body.ReviewNotExist),
		},
		{
			name: "Delete Product successfully",
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
				r.On("FindReview", mock.Anything, mock.Anything).Return(&body.ReviewProduct{}, nil)
				r.On("DeleteReview", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.DeleteProductReview(context.Background(), "123", "123")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_CreateProductReview(t *testing.T) {
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
			name: "Delete Product successfully",
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
				r.On("GetProductReviews", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*body.ReviewProduct{}, nil)
				r.On("CreateProductReview", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Delete Product successfully",
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
				r.On("GetProductReviews", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*body.ReviewProduct{{}, {}}, nil)
			},
			expectedErr: httperror.New(http.StatusBadRequest, body.ReviewAlreadyExist),
		},
		{
			name: "Delete Product successfully",
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
				r.On("GetProductReviews", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*body.ReviewProduct{}, nil)
				r.On("CreateProductReview", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.CreateProductReview(context.Background(), body.ReviewProductRequest{}, "123")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_UpdateProduct(t *testing.T) {
	testCase := []struct {
		name        string
		reqBody     body.UpdateProductRequest
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:    "Delete Product successfully",
			reqBody: body.UpdateProductRequest{ProductDetailRemove: []string{"123", "123"}},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("DeleteProductDetail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name:    "Delete Product successfully",
			reqBody: body.UpdateProductRequest{ProductDetailRemove: []string{"123", "123"}, ProductDetail: []body.UpdateProductDetailRequest{{}, {}}},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("DeleteProductDetail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				r.On("UpdateProductDetail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("test"))

			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.UpdateProduct(context.Background(), tc.reqBody, "123", "123")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_UpdateProductListedStatusBulk(t *testing.T) {

	testCase := []struct {
		name        string
		reqBody     body.UpdateProductListedStatusBulkRequest
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:    "Delete Product successfully",
			reqBody: body.UpdateProductListedStatusBulkRequest{ProductIDS: []string{"123", "123"}, ListedStatus: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("UpdateListedStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "Delete Product successfully",
			reqBody: body.UpdateProductListedStatusBulkRequest{ProductIDS: []string{"123", "123"}, ListedStatus: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("UpdateListedStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows)
			},
			expectedErr: httperror.New(http.StatusNotFound, body.UpdateProductFailed),
		},
		{
			name:    "Delete Product successfully",
			reqBody: body.UpdateProductListedStatusBulkRequest{ProductIDS: []string{"123", "123"}, ListedStatus: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("UpdateListedStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.UpdateProductListedStatusBulk(context.Background(), tc.reqBody)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestAdminUC_UpdateListedStatus(t *testing.T) {
	// var temp float64 = 10

	testCase := []struct {
		name        string
		reqBody     body.UpdateProductListedStatusBulkRequest
		userID      string
		mock        func(t *testing.T, r *mocks.Repository)
		expectedErr error
	}{
		{
			name:    "Delete Product successfully",
			reqBody: body.UpdateProductListedStatusBulkRequest{ProductIDS: []string{"123", "123"}, ListedStatus: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetListedStatus", mock.Anything, mock.Anything).Return(true, nil)
				r.On("UpdateListedStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "Delete Product successfully",
			reqBody: body.UpdateProductListedStatusBulkRequest{ProductIDS: []string{"123", "123"}, ListedStatus: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetListedStatus", mock.Anything, mock.Anything).Return(true, sql.ErrNoRows)
			},
			expectedErr: httperror.New(http.StatusBadRequest, body.ProductNotFound),
		},
		{
			name:    "Delete Product successfully",
			reqBody: body.UpdateProductListedStatusBulkRequest{ProductIDS: []string{"123", "123"}, ListedStatus: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetListedStatus", mock.Anything, mock.Anything).Return(true, fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
		{
			name:    "Delete Product successfully",
			reqBody: body.UpdateProductListedStatusBulkRequest{ProductIDS: []string{"123", "123"}, ListedStatus: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetListedStatus", mock.Anything, mock.Anything).Return(false, nil)
				r.On("UpdateListedStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "Delete Product successfully",
			reqBody: body.UpdateProductListedStatusBulkRequest{ProductIDS: []string{"123", "123"}, ListedStatus: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetListedStatus", mock.Anything, mock.Anything).Return(false, nil)
				r.On("UpdateListedStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows)
			},
			expectedErr: httperror.New(http.StatusNotFound, body.UpdateProductFailed),
		},
		{
			name:    "Delete Product successfully",
			reqBody: body.UpdateProductListedStatusBulkRequest{ProductIDS: []string{"123", "123"}, ListedStatus: true},
			mock: func(t *testing.T, r *mocks.Repository) {
				r.On("GetListedStatus", mock.Anything, mock.Anything).Return(false, nil)
				r.On("UpdateListedStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
			},
			expectedErr: fmt.Errorf("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sql, mock, _ := sqlmock.New()
			mock.ExpectBegin()
			r := mocks.NewRepository(t)
			u := NewProductUseCase(&config.Config{}, &postgre.TxRepo{PSQL: sql}, r)

			tc.mock(t, r)
			err := u.UpdateListedStatus(context.Background(), "123")
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
