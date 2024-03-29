package repository

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/internal/model"
	"murakali/internal/module/product"
	"murakali/internal/module/product/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"

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

func (r *productRepo) GetFavoriteProduct(ctx context.Context) ([]*model.ProductFavorite, error) {
	productFav := make([]*model.ProductFavorite, 0)
	res, err := r.PSQL.QueryContext(ctx, GetFavoriteProductQuery)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		favorite := model.ProductFavorite{}
		favorite.Product = &model.Product{}
		if errScan := res.Scan(&favorite.Product.ID, &favorite.Product.Title, &favorite.Count); errScan != nil {
			return nil, errScan
		}
		productFav = append(productFav, &favorite)
	}

	if res.Err() != nil {
		return nil, err
	}

	return productFav, nil
}

func (r *productRepo) GetRatingProduct(ctx context.Context) ([]*model.ProductRating, error) {
	productRating := make([]*model.ProductRating, 0)
	res, err := r.PSQL.QueryContext(ctx, GetRatingProductQuery)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		rating := model.ProductRating{}
		rating.Product = &model.Product{}
		if errScan := res.Scan(&rating.Product.ID, &rating.Product.Title, &rating.Product.ShopID, &rating.Count, &rating.Avg); errScan != nil {
			return nil, errScan
		}
		productRating = append(productRating, &rating)
	}

	if res.Err() != nil {
		return nil, err
	}

	return productRating, nil
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
			return nil, errScan
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
			&productData.ID,
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
			&productInfo.ShopID,
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
		); errScan != nil {
			return nil, err
		}

		res3, err3 := r.PSQL.QueryContext(
			ctx, GetProductDetailPhotosQuery, detail.ProductDetailID)

		if err3 != nil {
			return nil, err3
		}

		var productURLs []string
		for res3.Next() {
			var url body.URL
			if errScan := res3.Scan(
				&url.URL,
			); errScan != nil {
				return nil, err
			}
			productURLs = append(productURLs, url.URL)
		}
		detail.ProductURL = productURLs

		if promo != nil && (*promo.PromotionQuota) > 0 {
			discountedPrice := 0.0
			if promo.PromotionDiscountPercentage != nil {
				discountedPrice = *detail.NormalPrice - (*detail.NormalPrice * (*promo.PromotionDiscountPercentage / float64(100)))
			}
			if promo.PromotionDiscountFixPrice != nil {
				discountedPrice = *detail.NormalPrice - *promo.PromotionDiscountFixPrice
			}
			if *detail.NormalPrice >= *promo.PromotionMinProductPrice {
				if promo.PromotionDiscountFixPrice != nil && *promo.PromotionDiscountFixPrice > *promo.PromotionMaxDiscountPrice {
					discountedPrice = *detail.NormalPrice - *promo.PromotionDiscountFixPrice
				}
				if discountedPrice > *promo.PromotionMaxDiscountPrice {
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

			mapVariant[variant.Name] = variant.Type
		}
		detail.Variant = mapVariant

		res4, err4 := r.PSQL.QueryContext(
			ctx, GetVariantInfoQuery, detail.ProductDetailID)

		if err4 != nil {
			return nil, err4
		}

		var variantInfos []body.VariantInfo
		for res4.Next() {
			var info body.VariantInfo
			if errScan := res4.Scan(
				&info.VariantID,
				&info.VariantDetailID,
				&info.Name,
			); errScan != nil {
				return nil, err
			}
			variantInfos = append(variantInfos, info)
		}
		detail.VariantInfos = variantInfos

		productDetail = append(productDetail, &detail)
	}
	if res.Err() != nil {
		return nil, err
	}

	return productDetail, nil
}

func (r *productRepo) GetAllImageByProductDetailID(ctx context.Context, productDetailID string) ([]*string, error) {
	res, err := r.PSQL.QueryContext(
		ctx, GetProductDetailPhotosQuery, productDetailID)

	if err != nil {
		return nil, err
	}

	var productURLs []*string

	for res.Next() {
		var url string
		if errScan := res.Scan(
			&url,
		); errScan != nil {
			return nil, err
		}
		productURLs = append(productURLs, &url)
	}

	return productURLs, nil
}

func (r *productRepo) GetTotalProduct(ctx context.Context) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(ctx, GetTotalProductQuery).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *productRepo) GetProducts(ctx context.Context, pgn *pagination.Pagination, query *body.GetProductQueryRequest) ([]*body.Products,
	[]*model.Promotion, []*model.Voucher, error) {
	products := make([]*body.Products, 0)
	promotions := make([]*model.Promotion, 0)
	vouchers := make([]*model.Voucher, 0)

	queryOrderBySomething := fmt.Sprintf(OrderBySomething, pgn.GetSort(), pgn.GetLimit(),
		pgn.GetOffset())
	var queryWhereProvinceIds, queryWhereShopIds string

	if query.Shop != "" {
		queryWhereShopIds = fmt.Sprintf(WhereShopIds, query.Shop)
	}

	var queryListedStatus string
	switch query.ListedStatus {
	case 0:
		queryListedStatus = ``
	case 1:
		queryListedStatus = WhereListedStatusTrue
	case 2:
		queryListedStatus = WhereListedStatusFalse
	}

	var res *sql.Rows
	var err error
	if len(query.Province) > 0 {
		res, err = r.PSQL.QueryContext(
			ctx, GetProductsWithProvinceQuery+queryWhereShopIds+queryWhereProvinceIds+queryListedStatus+queryOrderBySomething,
			query.Search,
			query.Category,
			query.MinRating,
			query.MaxRating,
			query.MinPrice,
			query.MaxPrice,
			query.Province,
		)
	} else {
		res, err = r.PSQL.QueryContext(
			ctx, GetProductsQuery+queryWhereShopIds+queryWhereProvinceIds+queryListedStatus+queryOrderBySomething,
			query.Search,
			query.Category,
			query.MinRating,
			query.MaxRating,
			query.MinPrice,
			query.MaxPrice,
		)
		if err != nil {
			return nil, nil, nil, err
		}
		defer res.Close()
	}
	if err != nil {
		return nil, nil, nil, err
	}
	defer res.Close()

	for res.Next() {
		var productData body.Products
		var promo model.Promotion
		var voucher model.Voucher

		if errScan := res.Scan(
			&productData.ID,
			&productData.Title,
			&productData.UnitSold,
			&productData.RatingAVG,
			&productData.ThumbnailURL,
			&productData.MinPrice,
			&productData.MaxPrice,
			&productData.ViewCount,
			&promo.DiscountPercentage,
			&promo.DiscountFixPrice,
			&promo.MinProductPrice,
			&promo.MaxDiscountPrice,
			&productData.ShopName,
			&productData.CategoryName,
			&productData.ShopProvince,
			&productData.ListedStatus,
			&productData.CreatedAt,
			&productData.UpdatedAt,
			&productData.SKU,
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

func (r *productRepo) GetAllTotalProduct(ctx context.Context, query *body.GetProductQueryRequest) (int64, error) {
	var total int64

	var queryWhereProvinceIds, queryWhereShopIds string

	if query.Shop != "" {
		queryWhereShopIds = fmt.Sprintf(WhereShopIds, query.Shop)
	}

	var queryListedStatus string
	switch query.ListedStatus {
	case 0:
		queryListedStatus = ``
	case 1:
		queryListedStatus = WhereListedStatusTrue
	case 2:
		queryListedStatus = WhereListedStatusFalse
	}

	if len(query.Province) > 0 {
		if err := r.PSQL.QueryRowContext(ctx,
			GetAllTotalProductWithProvinceQuery+queryWhereShopIds+queryWhereProvinceIds+queryListedStatus,
			query.Search,
			query.Category,
			query.MinRating,
			query.MaxRating,
			query.MinPrice,
			query.MaxPrice,
			query.Province,
		).Scan(&total); err != nil {
			return 0, err
		}
	} else {
		if err := r.PSQL.QueryRowContext(ctx,
			GetAllTotalProductQuery+queryWhereShopIds+queryWhereProvinceIds+queryListedStatus,
			query.Search,
			query.Category,
			query.MinRating,
			query.MaxRating,
			query.MinPrice,
			query.MaxPrice,
		).Scan(&total); err != nil {
			return 0, err
		}
	}

	return total, nil
}

func (r *productRepo) GetFavoriteProducts(
	ctx context.Context, pgn *pagination.Pagination, query *body.GetProductQueryRequest, userID string) ([]*body.Products,
	[]*model.Promotion, []*model.Voucher, error) {
	products := make([]*body.Products, 0)
	promotions := make([]*model.Promotion, 0)
	vouchers := make([]*model.Voucher, 0)

	q := fmt.Sprintf(GetFavoriteProductsQuery, pgn.GetSort())
	res, err := r.PSQL.QueryContext(
		ctx, q,
		query.Search,
		query.Category,
		query.MinRating,
		query.MaxRating,
		query.MinPrice,
		query.MaxPrice,
		userID,
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
			&productData.ID,
			&productData.Title,
			&productData.UnitSold,
			&productData.RatingAVG,
			&productData.ThumbnailURL,
			&productData.MinPrice,
			&productData.MaxPrice,
			&productData.ViewCount,
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

func (r *productRepo) GetAllFavoriteTotalProduct(ctx context.Context, query *body.GetProductQueryRequest, userID string) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(ctx,
		GetAllTotalFavoriteProductQuery,
		query.Search,
		query.Category,
		query.MinRating,
		query.MaxRating,
		query.MinPrice,
		query.MaxPrice,
		userID,
	).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}
func (r *productRepo) CountUserFavoriteProduct(ctx context.Context, userID, productID string) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(ctx, CountUserFavoriteProduct, userID, productID).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *productRepo) CountSpecificFavoriteProduct(ctx context.Context, productID string) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(ctx, CountSpecificFavoriteProduct, productID).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *productRepo) CreateFavoriteProduct(ctx context.Context, tx postgre.Transaction, userID, productID string) error {
	_, err := r.PSQL.ExecContext(ctx, CreateFavoriteProductQuery, userID, productID)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepo) DeleteFavoriteProduct(ctx context.Context, tx postgre.Transaction, userID, productID string) error {
	_, err := r.PSQL.ExecContext(ctx, DeleteFavoriteProductQuery, userID, productID)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepo) FindFavoriteProduct(ctx context.Context, userID, productID string) (bool, error) {
	var isExist bool

	if err := r.PSQL.QueryRowContext(ctx, CheckFavoriteProductIsExistQuery, userID, productID).Scan(&isExist); err != nil {
		return false, err
	}

	return isExist, nil
}

func (r *productRepo) GetProductReviews(ctx context.Context,
	pgn *pagination.Pagination, productID string, query *body.GetReviewQueryRequest) ([]*body.ReviewProduct, error) {
	reviews := make([]*body.ReviewProduct, 0)

	q := fmt.Sprintf(GetReviewProductQuery, query.GetValidate(), pgn.GetSort())
	res, err := r.PSQL.QueryContext(
		ctx, q,
		productID,
		pgn.GetLimit(),
		pgn.GetOffset())

	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var reviewData body.ReviewProduct

		if errScan := res.Scan(
			&reviewData.ID,
			&reviewData.UserID,
			&reviewData.ProductID,
			&reviewData.Comment,
			&reviewData.Rating,
			&reviewData.ImageURL,
			&reviewData.CreatedAt,
			&reviewData.PhotoURL,
			&reviewData.Username,
		); errScan != nil {
			return nil, err
		}

		reviews = append(reviews, &reviewData)
	}

	if res.Err() != nil {
		return nil, err
	}

	return reviews, err
}

func (r *productRepo) GetTotalAllReviewProduct(ctx context.Context, productID string, query *body.GetReviewQueryRequest) (int64, error) {
	var total int64
	q := fmt.Sprintf(GetAllTotalReviewProductQuery, query.GetValidate())
	res, err := r.PSQL.QueryContext(
		ctx, q,
		productID)

	if err != nil {
		return 0, err
	}
	defer res.Close()

	for res.Next() {
		if errScan := res.Scan(
			&total,
		); errScan != nil {
			return 0, err
		}
	}

	return total, nil
}

func (r *productRepo) GetTotalReviewRatingByProductID(ctx context.Context, productID string) ([]*body.RatingProduct, error) {
	reviewRating := make([]*body.RatingProduct, 5)
	reviewRating[0] = &body.RatingProduct{Rating: 5, Count: 0}
	reviewRating[1] = &body.RatingProduct{Rating: 4, Count: 0}
	reviewRating[2] = &body.RatingProduct{Rating: 3, Count: 0}
	reviewRating[3] = &body.RatingProduct{Rating: 2, Count: 0}
	reviewRating[4] = &body.RatingProduct{Rating: 1, Count: 0}

	res, err := r.PSQL.QueryContext(
		ctx, GetTotalReviewRatingByProductIDQuery,
		productID)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var reviewRatingData body.RatingProduct

		if errScan := res.Scan(
			&reviewRatingData.Rating,
			&reviewRatingData.Count,
		); errScan != nil {
			return nil, err
		}

		reviewRating[reviewRatingData.Rating-1] = &reviewRatingData
	}

	if res.Err() != nil {
		return nil, err
	}

	return reviewRating, nil
}

func (r *productRepo) FindReview(ctx context.Context, reviewID string) (*body.ReviewProduct, error) {
	var review body.ReviewProduct
	if err := r.PSQL.QueryRowContext(ctx, GetReviewProductByIDQuery, reviewID).Scan(
		&review.ID,
		&review.UserID,
		&review.ProductID,
		&review.Comment,
		&review.Rating,
		&review.ImageURL,
		&review.Username,
		&review.PhotoURL,
		&review.Username,
	); err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *productRepo) CreateProductReview(ctx context.Context, tx postgre.Transaction, userID string, reqBody body.ReviewProductRequest) error {
	_, err := tx.ExecContext(
		ctx,
		CreateReviewQuery,
		userID,
		reqBody.ProductID,
		reqBody.Comment,
		reqBody.Rating,
		reqBody.PhotoURL,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) DeleteReview(ctx context.Context, tx postgre.Transaction, reviewID string) error {
	_, err := tx.ExecContext(ctx, DeleteReviewByIDQuery, reviewID)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) GetShopIDByUserID(ctx context.Context, userID string) (string, error) {
	var shopID string
	if err := r.PSQL.QueryRowContext(ctx, GetShopIDByUserIDQuery, userID).Scan(&shopID); err != nil {
		return "", err
	}

	return shopID, nil
}

func (r *productRepo) CreateProduct(ctx context.Context, tx postgre.Transaction, requestBody body.CreateProductInfoForQuery) (string, error) {
	var productID *uuid.UUID
	err := tx.QueryRowContext(
		ctx,
		CreateProductQuery,
		requestBody.CategoryID,
		requestBody.ShopID,
		requestBody.SKU,
		requestBody.Title,
		requestBody.Description,
		0,
		0,
		0,
		requestBody.ListedStatus,
		requestBody.Thumbnail,
		0,
		requestBody.MinPrice,
		requestBody.MaxPrice).Scan(&productID)
	if err != nil {
		return "", err
	}

	return productID.String(), nil
}

func (r *productRepo) CreateProductDetail(ctx context.Context, tx postgre.Transaction,
	requestBody body.CreateProductDetailRequest, productID string) (string, error) {
	var productDetailID *uuid.UUID
	err := tx.QueryRowContext(
		ctx,
		CreateProductDetailQuery,
		productID,
		requestBody.Price,
		requestBody.Stock,
		requestBody.Weight,
		requestBody.Size,
		requestBody.Hazardous,
		requestBody.Codition,
		requestBody.BulkPrice,
	).Scan(&productDetailID)
	if err != nil {
		return "", err
	}
	return productDetailID.String(), nil
}

func (r *productRepo) CreatePhoto(ctx context.Context, tx postgre.Transaction, productDetailID, url string) error {
	_, err := tx.ExecContext(
		ctx,
		CreatePhotoQuery,
		productDetailID,
		url,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) CreateVariant(ctx context.Context, tx postgre.Transaction, productDetailID, variantDetailID string) error {
	_, err := tx.ExecContext(
		ctx,
		CreateVariantQuery,
		productDetailID,
		variantDetailID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) CreateVariantDetail(ctx context.Context, tx postgre.Transaction,
	requestBody body.VariantDetailRequest) (string, error) {
	var ID string
	err := tx.QueryRowContext(
		ctx,
		CreateVariantDetailQuery,
		requestBody.Name,
		requestBody.Type,
	).Scan(&ID)
	if err != nil {
		return "", err
	}

	return ID, nil
}

func (r *productRepo) GetListedStatus(ctx context.Context, productID string) (bool, error) {
	var listedStatus bool
	if err := r.PSQL.QueryRowContext(ctx, GetListedStatusQuery, productID).Scan(&listedStatus); err != nil {
		return false, err
	}

	return listedStatus, nil
}

func (r *productRepo) UpdateListedStatus(ctx context.Context, tx postgre.Transaction, listedStatus bool, productID string) error {
	temp, err := tx.ExecContext(ctx, UpdateListedStatusQuery, listedStatus, productID)
	if err != nil {
		return err
	}

	rowsAffected, _ := temp.RowsAffected()
	if rowsAffected == 0 {
		return httperror.New(http.StatusNotFound, response.ProductNotExistMessage)
	}
	return nil
}

func (r *productRepo) UpdateProduct(ctx context.Context, tx postgre.Transaction, requestBody body.UpdateProductInfoForQuery, productID string) error {
	_, err := tx.ExecContext(ctx, UpdateProductQuery,
		requestBody.Title,
		requestBody.Description,
		requestBody.Thumbnail,
		requestBody.MinPrice,
		requestBody.MaxPrice,
		requestBody.ListedStatus,
		productID)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) UpdateProductFavorite(ctx context.Context, productID string, favCount int) error {
	_, err := r.PSQL.ExecContext(ctx, UpdateProductFavoriteQuery, favCount, productID)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) UpdateProductRating(ctx context.Context, productID string, ratingAvg float64) error {
	_, err := r.PSQL.ExecContext(ctx, UpdateProductRatingQuery, ratingAvg, productID)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) UpdateProductDetail(ctx context.Context,
	tx postgre.Transaction, requestBody body.UpdateProductDetailRequest, productID string) error {
	_, err := tx.ExecContext(ctx,
		UpdateProductDetailQuery,
		requestBody.Price,
		requestBody.Stock,
		requestBody.Weight,
		requestBody.Size,
		requestBody.Hazardous,
		requestBody.Codition,
		requestBody.BulkPrice,
		requestBody.ProductDetailID,
		productID)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) DeletePhoto(ctx context.Context, tx postgre.Transaction, productDetailID string) error {
	_, err := r.PSQL.ExecContext(ctx, DeletePhotoByIDQuery, productDetailID)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) DeleteProductDetail(ctx context.Context, tx postgre.Transaction, productDetailID string) error {
	_, err := r.PSQL.ExecContext(ctx, DeleteProductDetailByIDQuery, productDetailID)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) DeleteVariant(ctx context.Context, tx postgre.Transaction, productID string) error {
	_, err := r.PSQL.ExecContext(ctx, DeleteVariantByIDQuery, productID)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) GetMaxMinPriceByID(ctx context.Context, productID string) (*body.RangePrice, error) {
	var rangePrice body.RangePrice
	if err := r.PSQL.QueryRowContext(ctx, GetMaxMinPriceQuery, productID).Scan(&rangePrice.MaxPrice, &rangePrice.MinPrice); err != nil {
		return nil, err
	}
	return &rangePrice, nil
}

func (r *productRepo) GetShopProductRating(ctx context.Context, shopID string) (*model.ShopProductRating, error) {
	var shopProductRating model.ShopProductRating
	if err := r.PSQL.QueryRowContext(ctx, GetShopProductRating, shopID).Scan(&shopProductRating.ShopID,
		&shopProductRating.Count, &shopProductRating.Avg); err != nil {
		return nil, err
	}
	return &shopProductRating, nil
}

func (r *productRepo) UpdateVariant(ctx context.Context, tx postgre.Transaction, variantID, variantDetailID string) error {
	_, err := tx.ExecContext(ctx,
		UpdateVariantQuery,
		variantDetailID,
		variantID)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) UpdateShopProductRating(ctx context.Context, shop *model.ShopProductRating) error {
	_, err := r.PSQL.ExecContext(ctx, UpdateShopProductRating, shop.Count, shop.Avg, shop.ShopID)
	if err != nil {
		return err
	}
	return nil
}
