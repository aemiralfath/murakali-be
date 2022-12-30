package repository

import (
	"context"
	"database/sql"
	"murakali/internal/model"
	"murakali/internal/module/product"
	"murakali/internal/module/product/delivery/body"
	"murakali/pkg/pagination"

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

func (r *productRepo) GetRecommendedProducts(ctx context.Context, pgn *pagination.Pagination) ([]*body.Products,
	[]*model.Promotion, []*model.Voucher, error) {
	products := make([]*body.Products, 0)
	promotions := make([]*model.Promotion, 0)
	vouchers := make([]*model.Voucher, 0)

	res, err := r.PSQL.QueryContext(
		ctx, GetRecommendedProductsQuery,
		pgn.GetLimit(),
		pgn.GetOffset())

	if err != nil {
		return nil, nil, nil, err
	}
	defer res.Close()

	for res.Next() {
		var productData body.Products
		var promo model.Promotion
		var voucher model.Voucher

		if errScan := res.Scan(
			&productData.Title,
			&productData.UnitSold,
			&productData.RatingAVG,
			&productData.ThumbnailURL,
			&productData.MinPrice,
			&productData.MaxPrice,
			&promo.DiscountPercentage,
			&promo.DiscountFixPrice,
			&promo.MinProductPrice,
			&promo.MaxDiscountPrice,
			&voucher.DiscountPercentage,
			&voucher.DiscountFixPrice,
			&productData.ShopName,
			&productData.CategoryName,
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
			&productInfo.CategoryName,
			&productInfo.CategoryURL,
		); err != nil {
		return nil, err
	}

	return &productInfo, nil
}

func (r *productRepo) GetPromotionInfo(ctx context.Context, productID string) (*body.PromotionInfo, error) {
	var promotionInfo body.PromotionInfo

	if err := r.PSQL.QueryRowContext(ctx, GetPromotionDetailQuery, productID).
		Scan(&promotionInfo.PromotionName,
			&promotionInfo.PromotionDiscountPercentage,
			&promotionInfo.PromotionDiscountFixPrice,
			&promotionInfo.PromotionMinProductPrice,
			&promotionInfo.PromotionMaxDiscountPrice,
			&promotionInfo.PromotionQuota,
			&promotionInfo.PromotionMaxQuantity,
			&promotionInfo.PromotionActiveDate,
			&promotionInfo.PromotionExpiryDate,
		); err != nil {
		return nil, err
	}

	return &promotionInfo, nil
}

func (r *productRepo) GetProductDetail(ctx context.Context, productID string, promo *body.PromotionInfo) ([]*body.ProductDetail, error) {
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
			&detail.NormalPrice,
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

		if promo.PromotionDiscountPercentage != nil {
			if *detail.NormalPrice >= *promo.PromotionMinProductPrice {
				discountedPrice := *detail.NormalPrice - (*detail.NormalPrice * (*promo.PromotionDiscountPercentage / float64(100)))
				if discountedPrice > *promo.PromotionMaxDiscountPrice {
					discountedPrice = *detail.NormalPrice - *promo.PromotionMaxDiscountPrice
				}
				detail.DiscountPrice = &discountedPrice
			}
		}
		if promo.PromotionDiscountFixPrice != nil {
			if *detail.NormalPrice >= *promo.PromotionMinProductPrice {
				discountedPrice := *detail.NormalPrice - float64(*promo.PromotionDiscountFixPrice)
				if float64(*promo.PromotionDiscountFixPrice) > *promo.PromotionMaxDiscountPrice {
					discountedPrice = *detail.NormalPrice - *promo.PromotionMaxDiscountPrice
				}
				detail.DiscountPrice = &discountedPrice
			}
		}

		mapVariant := make(map[string]string, 0)

		res2, err2 := r.PSQL.QueryContext(
			ctx, GetVariantDetailQuery, detail.ProductDetailID)

		if err2 != nil {
			return nil, err2
		}

		for res2.Next() {
			var variant body.VariantDetail
			if errScan := res2.Scan(
				&variant.Type,
				&variant.Name,
			); errScan != nil {
				return nil, err
			}

			mapVariant[variant.Type] = variant.Name
		}
		detail.Variant = mapVariant

		productDetail = append(productDetail, &detail)
	}
	if res.Err() != nil {
		return nil, err
	}

	return productDetail, nil
}

func (r *productRepo) GetTotalProduct(ctx context.Context) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(ctx, GetTotalProductQuery).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}
