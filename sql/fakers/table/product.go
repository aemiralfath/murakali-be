package table

import (
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"math/rand"
	"murakali/internal/model"
	"murakali/pkg/postgre"
)

const InsertProductQuery = `INSERT INTO "product" (id, category_id, shop_id, sku, title, description, view_count, favorite_count, unit_sold, listed_status, thumbnail_url, rating_avg, min_price, max_price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);`
const InsertProductDetailQuery = `INSERT INTO "product_detail" (id, product_id, price, stock, weight, size, hazardous, condition, bulk_price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
const InsertVariantDetailQuery = `INSERT INTO "variant_detail" (id, name, type) VALUES ($1, $2, $3)`
const InsertVariantQuery = `INSERT INTO "variant" (id, product_detail_id, variant_detail_id) VALUES ($1, $2, $3);`

type ProductFaker struct {
	Size       int
	ShopID     string
	CategoryID string
	MinPrice   int
	MaxPrice   int
	ID         []string
}

func NewProductFaker(size int, shopID, categoryID string, minPrice, maxPrice int, id []string) ISeeder {
	return &ProductFaker{Size: size, ShopID: shopID, ID: id, MinPrice: minPrice, MaxPrice: maxPrice, CategoryID: categoryID}
}

func (f *ProductFaker) GenerateData(tx postgre.Transaction) error {
	shopID, err := uuid.Parse(f.ShopID)
	if err != nil {
		return err
	}

	categoryID, err := uuid.Parse(f.CategoryID)
	if err != nil {
		return err
	}

	for _, val := range f.ID {
		id, err := uuid.Parse(val)
		if err != nil {
			return err
		}

		if err := f.GenerateDataProduct(tx, id, categoryID, shopID); err != nil {
			return err
		}
	}

	for i := 0; i < f.Size; i++ {
		if err := f.GenerateDataProduct(tx, uuid.New(), categoryID, shopID); err != nil {
			return err
		}
	}

	return nil
}

func (f *ProductFaker) GenerateDataProduct(tx postgre.Transaction, id, categoryID, shopID uuid.UUID) error {
	data := f.GenerateProduct(id, categoryID, shopID)
	_, err := tx.Exec(InsertProductQuery, data.ID, data.CategoryID, data.ShopID, data.SKU, data.Title, data.Description, data.ViewCount, data.FavoriteCount, data.UnitSold, data.ListedStatus, data.ThumbnailURL, data.RatingAvg, data.MinPrice, data.MaxPrice)
	if err != nil {
		return err
	}

	productDetailID := uuid.New()
	if err := f.GenerateDataProductDetail(tx, productDetailID, data.ID, data.MaxPrice); err != nil {
		return err
	}

	variantDetailID := uuid.New()
	if _, err := tx.Exec(InsertVariantDetailQuery, variantDetailID, "warna", "hitam"); err != nil {
		return err
	}

	if _, err := tx.Exec(InsertVariantQuery, uuid.New(), productDetailID, variantDetailID); err != nil {
		return err
	}

	return nil
}

func (f *ProductFaker) GenerateDataProductDetail(tx postgre.Transaction, id, productID uuid.UUID, price float64) error {
	data := f.GenerateProductDetail(id, productID, price)
	_, err := tx.Exec(InsertProductDetailQuery, data.ID, data.ProductID, data.Price, data.Stock, data.Weight, data.Size, data.Hazardous, data.Condition, data.BulkPrice)
	if err != nil {
		return err
	}

	return nil
}

func (f *ProductFaker) GenerateProductDetail(id, productID uuid.UUID, price float64) *model.ProductDetail {
	return &model.ProductDetail{
		ID:        id,
		ProductID: productID,
		Price:     price,
		Stock:     float64(rand.Intn(20000)),
		Weight:    float64(rand.Intn(10-1) + 1),
		Size:      float64(rand.Intn(20000)),
		Hazardous: false,
		Condition: "new",
		BulkPrice: false,
	}
}

func (f *ProductFaker) GenerateProduct(id, categoryID, shopID uuid.UUID) *model.Product {
	name := faker.Username()
	price := rand.Intn(f.MaxPrice-f.MinPrice) + f.MinPrice
	return &model.Product{
		ID:            id,
		CategoryID:    categoryID,
		ShopID:        shopID,
		SKU:           slug.Make(name),
		Title:         name,
		Description:   faker.Paragraph(),
		ViewCount:     int64(rand.Intn(20000)),
		FavoriteCount: 0,
		UnitSold:      0,
		ListedStatus:  true,
		ThumbnailURL:  "https://cf.shopee.co.id/file/76a0969b7d64065bc13493bf55df1849_tn",
		RatingAvg:     0,
		MinPrice:      float64(price),
		MaxPrice:      float64(price),
	}
}
