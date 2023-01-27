package table

import (
	"github.com/google/uuid"
	"murakali/pkg/postgre"
)

const InsertShopQuery = `INSERT INTO "shop" (id, user_id, name, total_product, total_rating, rating_avg) VALUES ($1, $2, $3, $4, $5, $6);`
const InsertShopCourierQuery = `INSERT INTO "shop_courier" (shop_id, courier_id) VALUES ($1, $2)`

type ShopFaker struct {
	BuyerID      string
	BuyerCard    string
	ID           []string
	UserID       []string
	CategoryID   []string
	Name         []string
	TotalProduct []int
	CourierID    []string
}

func NewShopFaker(id, userID, categoryID, name []string, totalProduct []int, courierID []string, buyerID, buyerCard string) ISeeder {
	return &ShopFaker{ID: id, UserID: userID, CategoryID: categoryID, Name: name, TotalProduct: totalProduct, CourierID: courierID, BuyerID: buyerID, BuyerCard: buyerCard}
}

func (f *ShopFaker) GenerateData(tx postgre.Transaction) error {
	for i, val := range f.ID {
		id, err := uuid.Parse(val)
		if err != nil {
			return err
		}

		userID, err := uuid.Parse(f.UserID[i])
		if err != nil {
			return err
		}

		categoryID, err := uuid.Parse(f.CategoryID[i])
		if err != nil {
			return err
		}

		if _, err := tx.Exec(InsertShopQuery, id, userID, f.Name[i], f.TotalProduct[i], 0, 0); err != nil {
			return err
		}

		for _, cID := range f.CourierID {
			courierID, err := uuid.Parse(cID)
			if err != nil {
				return err
			}

			if _, err := tx.Exec(InsertShopCourierQuery, id, courierID); err != nil {
				return err
			}
		}

		productFaker := NewProductFaker(f.TotalProduct[i], id.String(), categoryID.String(), f.BuyerID, f.CourierID[0], f.BuyerCard, []string{})
		if err := productFaker.GenerateData(tx); err != nil {
			return err
		}
	}

	return nil
}
