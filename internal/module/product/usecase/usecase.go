package usecase

import (
	"context"
	"database/sql"

	"math"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/product"
	"murakali/internal/module/product/delivery/body"
	"murakali/internal/util"
	"murakali/pkg/httperror"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"

	"github.com/google/uuid"
)

type productUC struct {
	cfg         *config.Config
	txRepo      *postgre.TxRepo
	productRepo product.Repository
}

func NewProductUseCase(cfg *config.Config, txRepo *postgre.TxRepo, productRepo product.Repository) product.UseCase {
	return &productUC{cfg: cfg, txRepo: txRepo, productRepo: productRepo}
}

func (u *productUC) UpdateProductMetadata(ctx context.Context) error {
	productFav, err := u.productRepo.GetFavoriteProduct(ctx)
	if err != nil {
		return err
	}

	productRating, err := u.productRepo.GetRatingProduct(ctx)
	if err != nil {
		return err
	}

	var errFav error
	for _, favorite := range productFav {
		if err := u.productRepo.UpdateProductFavorite(ctx, favorite.Product.ID.String(), *favorite.Count); err != nil {
			errFav = err
		}
	}

	var errRating error
	shopID := make(map[string]string, 0)
	for _, rating := range productRating {
		shopID[rating.Product.ShopID.String()] = rating.Product.ShopID.String()
		if err := u.productRepo.UpdateProductRating(ctx, rating.Product.ID.String(), *rating.Avg); err != nil {
			errRating = err
		}
	}

	var errShop error
	for _, id := range shopID {
		shopProductRating, errShop := u.productRepo.GetShopProductRating(ctx, id)
		if errShop == nil {
			errShop = u.productRepo.UpdateShopProductRating(ctx, shopProductRating)
		}
	}
	if errShop != nil {
		return errShop
	}

	if errFav != nil {
		return errFav
	}

	if errRating != nil {
		return errRating
	}

	return nil
}

func (u *productUC) GetCategories(ctx context.Context) ([]*body.CategoryResponse, error) {
	categories, err := u.productRepo.GetCategories(ctx)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	categoryResponse := make([]*body.CategoryResponse, 0)
	for _, category := range categories {
		childCategories, err := u.GetCategoriesByParentID(ctx, category.ID)
		if err != nil {
			return nil, err
		}

		res := &body.CategoryResponse{
			ID:            category.ID,
			ParentID:      category.ParentID,
			Name:          category.Name,
			PhotoURL:      category.PhotoURL,
			ChildCategory: childCategories,
		}
		categoryResponse = append(categoryResponse, res)
	}

	return categoryResponse, nil
}

func (u *productUC) GetBanners(ctx context.Context) ([]*model.Banner, error) {
	banners, err := u.productRepo.GetBanners(ctx)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	return banners, nil
}

func (u *productUC) GetCategoriesByName(ctx context.Context, name string) ([]*body.CategoryResponse, error) {
	categories, err := u.productRepo.GetCategoriesByName(ctx, name)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	categoryResponse := make([]*body.CategoryResponse, 0)
	for _, category := range categories {
		childCategories, err := u.GetCategoriesByParentID(ctx, category.ID)
		if err != nil {
			return nil, err
		}

		res := &body.CategoryResponse{
			ID:            category.ID,
			ParentID:      category.ParentID,
			Name:          category.Name,
			PhotoURL:      category.PhotoURL,
			ChildCategory: childCategories,
		}
		categoryResponse = append(categoryResponse, res)
	}

	return categoryResponse, nil
}

func (u *productUC) GetCategoriesByParentID(ctx context.Context, parentID uuid.UUID) ([]*body.CategoryResponse, error) {
	categories, err := u.productRepo.GetCategoriesByParentID(ctx, parentID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	categoryResponse := make([]*body.CategoryResponse, 0)
	for _, category := range categories {
		childCategories, err := u.GetCategoriesByParentID(ctx, category.ID)
		if err != nil {
			return nil, err
		}

		res := &body.CategoryResponse{
			ID:            category.ID,
			ParentID:      category.ParentID,
			Name:          category.Name,
			PhotoURL:      category.PhotoURL,
			ChildCategory: childCategories,
		}
		categoryResponse = append(categoryResponse, res)
	}

	return categoryResponse, nil
}

func (u *productUC) GetRecommendedProducts(ctx context.Context, pgn *pagination.Pagination) (*pagination.Pagination, error) {
	totalRows, err := u.productRepo.GetTotalProduct(ctx)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages

	products, promotions, vouchers, err := u.productRepo.GetRecommendedProducts(ctx, pgn)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	resultProduct := make([]*body.Products, 0)
	totalData := len(products)
	for i := 0; i < totalData; i++ {
		p := &body.Products{
			ID:                        products[i].ID,
			Title:                     products[i].Title,
			UnitSold:                  products[i].UnitSold,
			RatingAVG:                 products[i].RatingAVG,
			ThumbnailURL:              products[i].ThumbnailURL,
			MinPrice:                  products[i].MinPrice,
			MaxPrice:                  products[i].MaxPrice,
			PromoDiscountPercentage:   promotions[i].DiscountPercentage,
			PromoDiscountFixPrice:     promotions[i].DiscountFixPrice,
			PromoMinProductPrice:      promotions[i].MinProductPrice,
			PromoMaxDiscountPrice:     promotions[i].MaxDiscountPrice,
			VoucherDiscountPercentage: vouchers[i].DiscountPercentage,
			VoucherDiscountFixPrice:   vouchers[i].DiscountFixPrice,
			ShopName:                  products[i].ShopName,
			CategoryName:              products[i].CategoryName,
		}
		p = u.CalculateDiscountProduct(p)
		resultProduct = append(resultProduct, p)
	}
	pgn.Rows = resultProduct

	return pgn, nil
}

func (u *productUC) CalculateDiscountProduct(p *body.Products) *body.Products {
	if p.PromoMaxDiscountPrice == nil {
		return p
	}
	var maxDiscountPrice float64
	var minProductPrice float64
	var discountPercentage float64
	var discountFixPrice float64
	var resultDiscount float64
	if p.PromoMinProductPrice != nil {
		minProductPrice = *p.PromoMinProductPrice
	}

	maxDiscountPrice = *p.PromoMaxDiscountPrice
	if p.PromoDiscountPercentage != nil {
		discountPercentage = *p.PromoDiscountPercentage
		if p.MinPrice >= minProductPrice && *p.PromoDiscountPercentage > 0 {
			resultDiscount = math.Min(maxDiscountPrice,
				p.MinPrice*(discountPercentage/100.00))
		}
	}

	if p.PromoDiscountFixPrice != nil {
		discountFixPrice = *p.PromoDiscountFixPrice
		if p.MinPrice >= minProductPrice && discountFixPrice > 0 {
			resultDiscount = math.Max(resultDiscount, discountFixPrice)
			resultDiscount = math.Min(resultDiscount, maxDiscountPrice)
		}
	}

	if resultDiscount > 0 {
		p.ResultDiscount = &resultDiscount
		p.SubPrice = p.MinPrice - resultDiscount
	}

	return p
}

func (u *productUC) GetProductDetail(ctx context.Context, productID string) (*body.ProductDetailResponse, error) {
	productInfo, err := u.productRepo.GetProductInfo(ctx, productID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	promotionInfo, err := u.productRepo.GetPromotionInfo(ctx, productID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	details, err := u.productRepo.GetProductDetail(ctx, productID, promotionInfo)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	result := body.ProductDetailResponse{
		ProductInfo:   productInfo,
		PromotionInfo: promotionInfo,
		ProductDetail: details,
	}
	return &result, nil
}

func (u *productUC) GetProducts(ctx context.Context, pgn *pagination.Pagination, query *body.GetProductQueryRequest) (*pagination.Pagination, error) {
	totalRows, err := u.productRepo.GetAllTotalProduct(ctx, query)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages

	products, promotions, vouchers, err := u.productRepo.GetProducts(ctx, pgn, query)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	resultProduct := make([]*body.Products, 0)
	totalData := len(products)
	for i := 0; i < totalData; i++ {
		p := &body.Products{
			ID:                        products[i].ID,
			Title:                     products[i].Title,
			UnitSold:                  products[i].UnitSold,
			RatingAVG:                 products[i].RatingAVG,
			ThumbnailURL:              products[i].ThumbnailURL,
			MinPrice:                  products[i].MinPrice,
			MaxPrice:                  products[i].MaxPrice,
			ViewCount:                 products[i].ViewCount,
			PromoDiscountPercentage:   promotions[i].DiscountPercentage,
			PromoDiscountFixPrice:     promotions[i].DiscountFixPrice,
			PromoMinProductPrice:      promotions[i].MinProductPrice,
			PromoMaxDiscountPrice:     promotions[i].MaxDiscountPrice,
			VoucherDiscountPercentage: vouchers[i].DiscountPercentage,
			VoucherDiscountFixPrice:   vouchers[i].DiscountFixPrice,
			ShopName:                  products[i].ShopName,
			CategoryName:              products[i].CategoryName,
			ShopProvince:              products[i].ShopProvince,
			ListedStatus:              products[i].ListedStatus,
			CreatedAt:                 products[i].CreatedAt,
			UpdatedAt:                 products[i].UpdatedAt,
			SKU:                       products[i].SKU,
		}
		p = u.CalculateDiscountProduct(p)
		resultProduct = append(resultProduct, p)
	}
	pgn.Rows = resultProduct

	return pgn, nil
}

func (u *productUC) GetAllProductImage(ctx context.Context, productID string) ([]*body.GetImageResponse, error) {
	var images []*body.GetImageResponse
	productInfo, err := u.productRepo.GetProductInfo(ctx, productID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	details, err := u.productRepo.GetProductDetail(ctx, productID, nil)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	images = append(images, &body.GetImageResponse{
		URL: productInfo.ThumbnailURL,
	})

	for _, detail := range details {
		imageDetails, err := u.productRepo.GetAllImageByProductDetailID(ctx, detail.ProductDetailID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}

		for _, image := range imageDetails {
			images = append(images, &body.GetImageResponse{
				ProductDetailID: &detail.ProductDetailID,
				URL:             *image,
			})
		}
	}

	return images, nil
}

func (u *productUC) GetFavoriteProducts(
	ctx context.Context, pgn *pagination.Pagination, query *body.GetProductQueryRequest, userID string) (*pagination.Pagination, error) {
	totalRows, err := u.productRepo.GetAllFavoriteTotalProduct(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages

	products, promotions, vouchers, err := u.productRepo.GetFavoriteProducts(ctx, pgn, query, userID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	resultProduct := make([]*body.Products, 0)
	totalData := len(products)
	for i := 0; i < totalData; i++ {
		p := &body.Products{
			ID:                        products[i].ID,
			Title:                     products[i].Title,
			UnitSold:                  products[i].UnitSold,
			RatingAVG:                 products[i].RatingAVG,
			ThumbnailURL:              products[i].ThumbnailURL,
			MinPrice:                  products[i].MinPrice,
			MaxPrice:                  products[i].MaxPrice,
			ViewCount:                 products[i].ViewCount,
			PromoDiscountPercentage:   promotions[i].DiscountPercentage,
			PromoDiscountFixPrice:     promotions[i].DiscountFixPrice,
			PromoMinProductPrice:      promotions[i].MinProductPrice,
			PromoMaxDiscountPrice:     promotions[i].MaxDiscountPrice,
			VoucherDiscountPercentage: vouchers[i].DiscountPercentage,
			VoucherDiscountFixPrice:   vouchers[i].DiscountFixPrice,
			ShopName:                  products[i].ShopName,
			CategoryName:              products[i].CategoryName,
		}
		p = u.CalculateDiscountProduct(p)
		resultProduct = append(resultProduct, p)
	}
	pgn.Rows = resultProduct

	return pgn, nil
}

func (u *productUC) CheckProductIsFavorite(
	ctx context.Context, userID, productID string) bool {
	totalRows, _ := u.productRepo.CountUserFavoriteProduct(ctx, userID, productID)
	return totalRows > 0
}

func (u *productUC) CountSpecificFavoriteProduct(
	ctx context.Context, productID string) (int64, error) {
	totalRows, err := u.productRepo.CountSpecificFavoriteProduct(ctx, productID)
	if err != nil {
		return 0, err
	}

	return totalRows, err
}

func (u *productUC) CreateFavoriteProduct(ctx context.Context, productID, userID string) error {
	_, err := u.productRepo.GetProductInfo(ctx, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.ProductNotExistMessage)
		}
		return err
	}

	isExist, err := u.productRepo.FindFavoriteProduct(ctx, userID, productID)

	if err != nil {
		return err
	}

	if isExist {
		return httperror.New(http.StatusBadRequest, response.ProductAlreadyInFavMessage)
	}

	err = u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		err = u.productRepo.CreateFavoriteProduct(ctx, tx, userID, productID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *productUC) DeleteFavoriteProduct(ctx context.Context, productID, userID string) error {
	isExist, err := u.productRepo.FindFavoriteProduct(ctx, userID, productID)
	if err != nil {
		return err
	}

	if !isExist {
		return httperror.New(http.StatusBadRequest, response.ProductNotExistMessage)
	}
	err = u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		err = u.productRepo.DeleteFavoriteProduct(ctx, tx, userID, productID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *productUC) GetProductReviews(ctx context.Context, pgn *pagination.Pagination,
	productID string, query *body.GetReviewQueryRequest) (*pagination.Pagination, error) {
	totalRows, err := u.productRepo.GetTotalAllReviewProduct(ctx, productID, query)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages

	reviews, err := u.productRepo.GetProductReviews(ctx, pgn, productID, query)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	pgn.Rows = reviews

	return pgn, nil
}

func (u *productUC) GetTotalReviewRatingByProductID(ctx context.Context, productID string) (*body.AllRatingProduct, error) {
	ratings, err := u.productRepo.GetTotalReviewRatingByProductID(ctx, productID)

	valueRating := 0
	totalRating := 0

	for i := 0; i < len(ratings); i++ {
		valueRating += ratings[i].Rating * ratings[i].Count
		totalRating += ratings[i].Count
	}

	allTotalRating := &body.AllRatingProduct{
		AvgRating:     float64(valueRating) / float64(totalRating),
		TotalRating:   float64(totalRating),
		RatingProduct: ratings,
	}

	if err != nil {
		return nil, err
	}

	return allTotalRating, nil
}

func (u *productUC) CreateProduct(ctx context.Context, requestBody body.CreateProductRequest, userID string) error {
	shopID, errGet := u.productRepo.GetShopIDByUserID(ctx, userID)
	if errGet != nil {
		if errGet == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotExistMessage)
		}
		return errGet
	}

	err := u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		totalData := len(requestBody.ProductDetail)
		minPriceTemp, maxPriceTemp := requestBody.ProductDetail[0].Price, requestBody.ProductDetail[0].Price
		for i := 0; i < totalData; i++ {
			if requestBody.ProductDetail[i].Price < minPriceTemp {
				minPriceTemp = requestBody.ProductDetail[i].Price
			}
			if requestBody.ProductDetail[i].Price > maxPriceTemp {
				maxPriceTemp = requestBody.ProductDetail[i].Price
			}
		}

		var tempBodyProduct = body.CreateProductInfoForQuery{
			Title:        requestBody.ProductInfo.Title,
			Description:  requestBody.ProductInfo.Description,
			Thumbnail:    requestBody.ProductInfo.Thumbnail,
			CategoryID:   requestBody.ProductInfo.CategoryID,
			ListedStatus: requestBody.ProductInfo.ListedStatus,
			MinPrice:     minPriceTemp,
			MaxPrice:     maxPriceTemp,
			ShopID:       shopID,
			SKU:          util.SKUGenerator(requestBody.ProductInfo.Title),
		}

		productID, err := u.productRepo.CreateProduct(ctx, tx, tempBodyProduct)
		if err != nil {
			return err
		}

		for i := 0; i < totalData; i++ {
			productDetilID, err := u.productRepo.CreateProductDetail(ctx, tx, requestBody.ProductDetail[i], productID)
			if err != nil {
				return err
			}

			totalDataPhoto := len(requestBody.ProductDetail[i].Photo)
			if totalDataPhoto > 0 {
				for k := 0; k < totalDataPhoto; k++ {
					err = u.productRepo.CreatePhoto(ctx, tx, productDetilID, requestBody.ProductDetail[i].Photo[k])
					if err != nil {
						return err
					}
				}
			}

			totalDataVariant := len(requestBody.ProductDetail[i].VariantDetail)
			if totalDataVariant > 0 {
				for j := 0; j < totalDataVariant; j++ {
					variantDetailID, err := u.productRepo.CreateVariantDetail(ctx, tx, requestBody.ProductDetail[i].VariantDetail[j])
					if err != nil {
						return err
					}
					err = u.productRepo.CreateVariant(ctx, tx, productDetilID, variantDetailID)
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (u *productUC) UpdateListedStatus(ctx context.Context, productID string) error {
	listedStatus, err := u.productRepo.GetListedStatus(ctx, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, body.ProductNotFound)
		}
		return err
	}
	tempListedStatus := false
	if listedStatus {
		tempListedStatus = false
	} else {
		tempListedStatus = true
	}

	errTx := u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		if err := u.productRepo.UpdateListedStatus(ctx, tx, tempListedStatus, productID); err != nil {
			if err == sql.ErrNoRows {
				return httperror.New(http.StatusNotFound, body.UpdateProductFailed)
			}
			return err
		}
		return nil
	})

	if errTx != nil {
		return errTx
	}
	return nil
}

func (u *productUC) UpdateProductListedStatusBulk(ctx context.Context, productRequest body.UpdateProductListedStatusBulkRequest) error {
	errTx := u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		for i := 0; i < len(productRequest.ProductIDS); i++ {
			if err := u.productRepo.UpdateListedStatus(ctx, tx, productRequest.ListedStatus, productRequest.ProductIDS[i]); err != nil {
				if err == sql.ErrNoRows {
					return httperror.New(http.StatusNotFound, body.UpdateProductFailed)
				}
				return err
			}
		}
		return nil
	})

	if errTx != nil {
		return errTx
	}
	return nil
}

func (u *productUC) UpdateProduct(ctx context.Context, requestBody body.UpdateProductRequest, userID, productID string) error {
	errTx := u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		totalData := len(requestBody.ProductDetail)

		totalDataRemove := len(requestBody.ProductDetailRemove)
		if totalDataRemove > 0 {
			for i := 0; i < totalDataRemove; i++ {
				err := u.productRepo.DeleteProductDetail(ctx, tx, requestBody.ProductDetailRemove[i])
				if err != nil {
					return err
				}
			}
		}

		for i := 0; i < totalData; i++ {
			err := u.productRepo.UpdateProductDetail(ctx, tx, requestBody.ProductDetail[i], productID)
			if err != nil {
				return err
			}

			totalDataPhoto := len(requestBody.ProductDetail[i].Photo)
			if totalDataPhoto > 0 {
				err = u.productRepo.DeletePhoto(ctx, tx, requestBody.ProductDetail[i].ProductDetailID)
				if err != nil {
					return err
				}
				for k := 0; k < totalDataPhoto; k++ {
					err = u.productRepo.CreatePhoto(ctx, tx, requestBody.ProductDetail[i].ProductDetailID, requestBody.ProductDetail[i].Photo[k])
					if err != nil {
						return err
					}
				}
			}

			totalDataVariant := len(requestBody.ProductDetail[i].VariantDetailID)
			if totalDataVariant > 0 {
				for j := 0; j < totalDataVariant; j++ {
					errVariant := u.productRepo.UpdateVariant(ctx, tx,
						requestBody.ProductDetail[i].VariantDetailID[j].VariantID,
						requestBody.ProductDetail[i].VariantDetailID[j].VariantDetailID)
					if errVariant != nil {
						return errVariant
					}
				}
			}

			totalDataVariantRemove := len(requestBody.ProductDetail[i].VariantIDRemove)
			if totalDataVariantRemove > 0 {
				for j := 0; j < totalDataVariantRemove; j++ {
					err := u.productRepo.DeleteVariant(ctx, tx, requestBody.ProductDetail[i].VariantIDRemove[j])
					if err != nil {
						return err
					}
				}
			}
		}
		var rangePrice *body.RangePrice
		rangePrice, errMaxMin := u.productRepo.GetMaxMinPriceByID(ctx, productID)
		if errMaxMin != nil {
			return errMaxMin
		}
		var tempBodyProduct = body.UpdateProductInfoForQuery{
			Title:        requestBody.ProductInfo.Title,
			Description:  requestBody.ProductInfo.Description,
			Thumbnail:    requestBody.ProductInfo.Thumbnail,
			ListedStatus: requestBody.ProductInfo.ListedStatus,
			MinPrice:     rangePrice.MinPrice,
			MaxPrice:     rangePrice.MaxPrice,
		}
		err := u.productRepo.UpdateProduct(ctx, tx, tempBodyProduct, productID)
		if err != nil {
			return err
		}

		return nil
	})

	if errTx != nil {
		return errTx
	}
	return nil
}

func (u *productUC) CreateProductReview(ctx context.Context, reqBody body.ReviewProductRequest, userID string) error {
	gotReview, err := u.productRepo.GetProductReviews(ctx, &pagination.Pagination{}, reqBody.ProductID, &body.GetReviewQueryRequest{
		UserID: userID,
	})

	if len(gotReview) > 0 {
		return httperror.New(http.StatusBadRequest, body.ReviewAlreadyExist)
	}

	err = u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		err = u.productRepo.CreateProductReview(ctx, tx, userID, reqBody)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *productUC) DeleteProductReview(ctx context.Context, reviewID, userID string) error {
	gotReview, err := u.productRepo.FindReview(ctx, reviewID)
	if err != nil {
		return err
	}
	if gotReview == nil {
		return httperror.New(http.StatusBadRequest, body.ReviewNotExist)
	}
	err = u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		err = u.productRepo.DeleteReview(ctx, tx, reviewID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
