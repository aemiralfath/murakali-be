package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/product"
	"murakali/internal/module/product/delivery/body"
	"murakali/pkg/postgre"

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
			fmt.Println("child category:", childCategories)
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

func (u *productUC) GetRecommendedProducts(ctx context.Context) (*body.RecommendedProductResponse, error) {

	limit := 18
	products, err := u.productRepo.GetRecommendedProducts(ctx, limit)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	for _, product := range products {
		product = u.CalculateDiscountProduct(product)
	}

	recommendedProductResponse := &body.RecommendedProductResponse{
		Limit:    limit,
		Products: products,
	}
	return recommendedProductResponse, nil
}

func (u *productUC) CalculateDiscountProduct(p *body.Products) *body.Products {

	if p.PromoMaxDiscountPrice == nil {
		return p
	}
	var maxDiscountPrice float64 = *p.PromoMaxDiscountPrice
	var minProductPrice float64
	var discountPercentage float64
	var discountFixPrice float64
	var resultDiscount float64
	if p.PromoMinProductPrice != nil {
		minProductPrice = *p.PromoMinProductPrice
	}

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
