package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/config"
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
