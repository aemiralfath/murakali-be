package repository

import (
	"context"
	"database/sql"
	"murakali/internal/model"
	"murakali/internal/module/product"
	"murakali/internal/module/product/delivery/body"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type productRepo struct {
	PSQL        *sql.DB
	RedisClient *redis.Client
}

func NewProductRepository(psql *sql.DB, client *redis.Client) product.Repository {
	return &productRepo{
		PSQL:        psql,
		RedisClient: client,
	}
}

func (r *productRepo) GetCategories(ctx context.Context) ([]*model.Category, error) {
	categories := make([]*model.Category, 0)
	res, err := r.PSQL.QueryContext(
		ctx, GetCategoriesQuery)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var category model.Category
		if errScan := res.Scan(
			&category.ID,
			&category.ParentID,
			&category.Name,
			&category.PhotoURL,
		); errScan != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	if res.Err() != nil {
		return nil, err
	}

	return categories, nil
}

func (r *productRepo) GetBanners(ctx context.Context) ([]*model.Banner, error) {
	banners := make([]*model.Banner, 0)

	res, err := r.PSQL.QueryContext(ctx, GetBannersQuery)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var banner model.Banner
		if errScan := res.Scan(
			&banner.ID,
			&banner.Title,
			&banner.Content,
			&banner.ImageURL,
			&banner.PageURL,
			&banner.IsActive); errScan != nil {
			return nil, err
		}

		banners = append(banners, &banner)
	}

	if res.Err() != nil {
		return nil, err
	}

	return banners, nil
}

func (r *productRepo) GetCategoriesByName(ctx context.Context, name string) ([]*model.Category, error) {
	categories := make([]*model.Category, 0)

	res, err := r.PSQL.QueryContext(
		ctx, GetCategoriesByNameQuery, name)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var category model.Category
		if errScan := res.Scan(
			&category.ID,
			&category.ParentID,
			&category.Name,
			&category.PhotoURL,
		); errScan != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	if res.Err() != nil {
		return nil, err
	}

	return categories, nil
}

func (r *productRepo) GetCategoriesByParentID(ctx context.Context, parentID uuid.UUID) ([]*model.Category, error) {
	categories := make([]*model.Category, 0)

	res, err := r.PSQL.QueryContext(
		ctx, GetCategoriesByParentIdQuery, parentID)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var category model.Category
		if errScan := res.Scan(
			&category.ID,
			&category.ParentID,
			&category.Name,
			&category.PhotoURL,
		); errScan != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	if res.Err() != nil {
		return nil, err
	}

	return categories, nil
}

func (r *productRepo) GetRecommendedProducts(ctx context.Context, limit int) ([]*model.Product, []*model.Promotion, []*model.Voucher, error) {
	products := make([]*model.Product, 0)
	promotions := make([]*model.Promotion, 0)
	vouchers := make([]*model.Voucher, 0)

	res, err := r.PSQL.QueryContext(
		ctx, GetRecommendedProductsQuery,
		limit)

	if err != nil {
		return nil, nil, nil, err
	}
	defer res.Close()

	for res.Next() {
		var productData model.Product
		var promo model.Promotion
		var voucher model.Voucher

		if errScan := res.Scan(
			&productData.Title,
			&productData.UnitSold,
			&productData.RatingAvg,
			&productData.ThumbnailURL,
			&productData.MinPrice,
			&productData.MaxPrice,
			&promo.DiscountPercentage,
			&promo.DiscountFixPrice,
			&promo.MinProductPrice,
			&promo.MaxDiscountPrice,
			&voucher.DiscountPercentage,
			&voucher.DiscountFixPrice,
		); errScan != nil {
			return nil, nil, nil, err
		}

		products = append(products, &productData)
		promotions = append(promotions, &promo)
		vouchers = append(vouchers, &voucher)
	}

	if res.Err() != nil {
		return nil, nil, nil, err
	}

	return products, promotions, vouchers, err
}

func (r *productRepo) GetProductInfo(ctx context.Context, productID string) (*body.ProductInfo, error) {
	var productInfo body.ProductInfo

	if err := r.PSQL.QueryRowContext(ctx, GetProductInfoQuery, productID).
		Scan(&productInfo.ProductID,
			&productInfo.SKU,
			&productInfo.Title,
			&productInfo.Description,
			&productInfo.ViewCount,
			&productInfo.FavoriteCount,
			&productInfo.UnitSold,
			&productInfo.ListedStatus,
			&productInfo.ThumbnailURL,
			&productInfo.RatingAVG,
			&productInfo.MinPrice,
			&productInfo.MaxPrice,
		); err != nil {

		return nil, err
	}

	return &productInfo, nil
}

func (r *productRepo) GetProductDetail(ctx context.Context, productID string) ([]*body.ProductDetail, error) {
	productDetail := make([]*body.ProductDetail, 0)

	res, err := r.PSQL.QueryContext(
		ctx, GetProductDetailQuery, productID)

	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {

		var detail body.ProductDetail

		if errScan := res.Scan(
			&detail.ProductDetailID,
			&detail.Price,
			&detail.Stock,
			&detail.Weight,
			&detail.Size,
			&detail.Hazardous,
			&detail.Condition,
			&detail.BulkPrice,
			&detail.ProductURL,
		); errScan != nil {
			return nil, err
		}

		variantDetail := make([]*body.VariantDetail, 0)

		res2, err := r.PSQL.QueryContext(
			ctx, GetVariantDetailQuery, detail.ProductDetailID)

		if err != nil {
			return nil, err
		}

		defer res2.Close()
		for res2.Next() {

			var variant body.VariantDetail

			if errScan := res2.Scan(
				&variant.Type,
				&variant.Name,
			); errScan != nil {
				return nil, err
			}

			variantDetail = append(variantDetail, &variant)
		}
		detail.Variant = variantDetail

		productDetail = append(productDetail, &detail)

	}
	if res.Err() != nil {
		return nil, err
	}

	return productDetail, nil

}
