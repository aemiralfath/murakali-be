package product

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/product/delivery/body"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"

	"github.com/google/uuid"
)

type Repository interface {
	GetCategories(ctx context.Context) ([]*model.Category, error)
	GetBanners(ctx context.Context) ([]*model.Banner, error)
	GetCategoriesByName(ctx context.Context, name string) ([]*model.Category, error)
	GetCategoriesByParentID(ctx context.Context, parentID uuid.UUID) ([]*model.Category, error)
	GetRecommendedProducts(ctx context.Context, pgn *pagination.Pagination) ([]*body.Products, []*model.Promotion, []*model.Voucher, error)
	GetTotalProduct(ctx context.Context) (int64, error)
	GetProductInfo(ctx context.Context, productID string) (*body.ProductInfo, error)
	GetProductDetail(ctx context.Context, productID string, promo *body.PromotionInfo) ([]*body.ProductDetail, error)
	GetPromotionInfo(ctx context.Context, productID string) (*body.PromotionInfo, error)
	GetProducts(ctx context.Context, pgn *pagination.Pagination, query *body.GetProductQueryRequest) ([]*body.Products,
		[]*model.Promotion, []*model.Voucher, error)
	GetAllTotalProduct(ctx context.Context, query *body.GetProductQueryRequest) (int64, error)
	GetFavoriteProducts(ctx context.Context, pgn *pagination.Pagination, query *body.GetProductQueryRequest, userID string) ([]*body.Products,
		[]*model.Promotion, []*model.Voucher, error)
	GetAllFavoriteTotalProduct(ctx context.Context, query *body.GetProductQueryRequest, userID string) (int64, error)

	GetShopIDByUserID(ctx context.Context, userID string) (string, error)
	CreateProduct(ctx context.Context, tx postgre.Transaction, requestBody body.CreateProductInfoForQuery) (string, error)
	CreateProductDetail(ctx context.Context, tx postgre.Transaction, requestBody body.CreateProductDetailRequest, ProductID string) (string, error)
	CreatePhoto(ctx context.Context, tx postgre.Transaction, productDetailID string, URL string) error
	CreateVideo(ctx context.Context, tx postgre.Transaction, productDetailID string, URL string) error
	CreateVariant(ctx context.Context, tx postgre.Transaction, productDetailID string, variantDetailID string) error
	CreateProductCourier(ctx context.Context, tx postgre.Transaction, productDetailID string, courierID string) error

	GetListedStatus(ctx context.Context, productID string) (bool, error)
	UpdateListedStatus(ctx context.Context, listedStatus bool, productID string) error
}
