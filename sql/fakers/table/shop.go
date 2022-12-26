package table

import (
	"github.com/google/uuid"
	"murakali/pkg/postgre"
)

const InsertShopQuery = `INSERT INTO "shop" (id, user_id, name, total_product, total_rating, rating_avg) VALUES ($1, $2, $3, $4, $5, $6);`

type ShopFaker struct {
	ID           []string
	UserID       []string
	CategoryID   []string
	Name         []string
	TotalProduct []int
}

func NewShopFaker(id, userID, categoryID, name []string, totalProduct []int) ISeeder {
	return &ShopFaker{ID: id, UserID: userID, CategoryID: categoryID, Name: name, TotalProduct: totalProduct}
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

		productFaker := NewProductFaker(f.TotalProduct[i], id.String(), categoryID.String(), []string{})
		if err := productFaker.GenerateData(tx); err != nil {
			return err
		}
	}

	return nil
}
