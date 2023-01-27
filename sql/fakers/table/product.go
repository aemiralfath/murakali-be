package table

import (
	"math/rand"
	"murakali/internal/model"
	"murakali/pkg/postgre"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"murakali/internal/constant"
	"murakali/internal/util"
	"time"
)

const InsertProductQuery = `INSERT INTO "product" (id, category_id, shop_id, sku, title, description, view_count, favorite_count, unit_sold, listed_status, thumbnail_url, rating_avg, min_price, max_price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);`
const InsertProductDetailQuery = `INSERT INTO "product_detail" (id, product_id, price, stock, weight, size, hazardous, condition, bulk_price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
const InsertVariantDetailQuery = `INSERT INTO "variant_detail" (id, name, type) VALUES ($1, $2, $3)`
const InsertVariantQuery = `INSERT INTO "variant" (id, product_detail_id, variant_detail_id) VALUES ($1, $2, $3);`
const InsertProductDetailPhoto = `INSERT INTO "photo" (product_detail_id, url) VALUES ($1, $2);`
const InsertTransactionQuery = `INSERT INTO "transaction" (id, card_number, invoice, total_price, paid_at, expired_at) VALUES ($1, $2, $3, $4, $5, $6);`
const InsertOrderQuery = `INSERT INTO "order" (id, transaction_id, shop_id, user_id, courier_id, order_status_id, total_price, delivery_fee, resi_no, buyer_address, shop_address, is_withdraw, created_at, arrived_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);`
const InsertOrderItemQuery = `INSERT INTO "order_item" (order_id, product_detail_id, quantity, item_price, total_price, is_review) VALUES ($1, $2, $3, $4, $5, $6)`
const InsertOrderReviewQuery = `INSERT INTO "review" (user_id, product_id, comment, rating, created_at) VALUES ($1, $2, $3, $4, $5)`

type ProductFaker struct {
	Size       int
	ShopID     string
	CategoryID string
	UserID     string
	CourierID  string
	CardNumber string
	ID         []string
}

func NewProductFaker(size int, shopID, categoryID, userID, courierID, cardNumber string, id []string) ISeeder {
	return &ProductFaker{Size: size, ShopID: shopID, ID: id, CategoryID: categoryID, UserID: userID, CardNumber: cardNumber, CourierID: courierID}
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

	_, err = tx.Exec(InsertProductDetailPhoto, data.ID, "https://cf.shopee.co.id/file/76a0969b7d64065bc13493bf55df1849_tn")
	if err != nil {
		return err
	}

	return nil
}

func (f *ProductFaker) GenerateTransactions(tx postgre.Transaction, productDetail *model.ProductDetail) error {
	txID := uuid.New()
	deliveryFee := float64(8000)
	randomTime := time.Now().AddDate(0, 0, -1*rand.Intn(31))
	invoice, errInvoice := util.GenerateInvoice()
	if errInvoice != nil {
		return errInvoice
	}

	_, errTx := tx.Exec(InsertTransactionQuery, txID, f.CardNumber, invoice, productDetail.Price+deliveryFee, randomTime, randomTime)
	if errTx != nil {
		return errTx
	}

	orderID := uuid.New()
	buyerAddress := `{"id":"5a9a2b53-b79c-49d9-af51-595bbb998e15","user_id":"65ac2cb6-f25d-4bc8-aa4c-c59ea5da52b5","name":"Kost","province_id":6,"city_id":153,"province":"DKI Jakarta","city":"Jakarta Selatan","district":"Mampang Prapatan","sub_district":"Mampang Prapatan","address_detail":"JSR","zip_code":"12790","is_default":true,"is_shop_default":true,"created_at":"2023-01-25T08:02:33.621316Z","updated_at":{"Time":"2023-01-26T05:06:52.935979Z","Valid":true},"deleted_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}}`
	shopAddress := `{"id":"a78bc3c2-9ca3-473d-a450-3ed7f9689476","user_id":"9e3a2d9d-0479-4819-a40b-068453bcaf04","name":"Home","province_id":33,"city_id":327,"province":"Sumatera Selatan","city":"Palembang","district":"Ilir Timur II","sub_district":"2 Ilir","address_detail":"no 91","zip_code":"30118","is_default":true,"is_shop_default":true,"created_at":"2023-01-22T15:40:32.207751Z","updated_at":{"Time":"0001-01-01T00:00:00Z","Valid":false},"deleted_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}}`
	resiNo, err := util.GenerateRandomAlpaNumeric(16)
	if err != nil {
		return err
	}

	_, errOrder := tx.Exec(InsertOrderQuery, orderID, f.ShopID, f.UserID, f.CourierID, constant.OrderStatusCompleted, productDetail.Price, deliveryFee, resiNo, buyerAddress, shopAddress, true, randomTime, randomTime)
	if errOrder != nil {
		return errOrder
	}

	_, errItem := tx.Exec(InsertOrderItemQuery, orderID, productDetail.ID, 1, productDetail.Price, productDetail.Price, true)
	if errItem != nil {
		return errItem
	}

	_, errReview := tx.Exec(InsertOrderReviewQuery, f.UserID, productDetail.ProductID, faker.Paragraph(), rand.Intn(5-1)+1, randomTime)
	if errReview != nil {
		return errReview
	}

	return nil
}

func (f *ProductFaker) GenerateProductDetail(id, productID uuid.UUID, price float64) *model.ProductDetail {
	return &model.ProductDetail{
		ID:        id,
		ProductID: productID,
		Price:     price,
		Stock:     float64(rand.Intn(20000)),
		Weight:    float64(rand.Intn(10-1)+1) * 100,
		Size:      float64(rand.Intn(20000)),
		Hazardous: false,
		Condition: "new",
		BulkPrice: false,
	}
}

func (f *ProductFaker) GenerateProduct(id, categoryID, shopID uuid.UUID) *model.Product {
	name := faker.Username()
	price := (rand.Intn(1000-100) + 100) * 1000
	return &model.Product{
		ID:            id,
		CategoryID:    categoryID,
		ShopID:        shopID,
		SKU:           slug.Make(name),
		Title:         name,
		Description:   faker.Paragraph(),
		ViewCount:     int64(rand.Intn(20000)),
		FavoriteCount: 0,
		UnitSold:      1,
		ListedStatus:  true,
		ThumbnailURL:  "https://cf.shopee.co.id/file/76a0969b7d64065bc13493bf55df1849_tn",
		RatingAvg:     0,
		MinPrice:      float64(price),
		MaxPrice:      float64(price),
	}
}
