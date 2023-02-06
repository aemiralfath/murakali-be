package table

import (
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"murakali/pkg/postgre"
)

const InsertCategoryQuery = `INSERT INTO "category" (id, parent_id, name, photo_url) VALUES ($1, $2, $3, $4)`

type CategoryFaker struct {
	Size     int
	ID       []string
	ParentID []string
	Name     []string
	PhotoURL []string
}

func NewCategoryFaker(size int, id, name, photoURL []string, parentID []string) ISeeder {
	return &CategoryFaker{Size: size, ID: id, Name: name, PhotoURL: photoURL, ParentID: parentID}
}

func (f *CategoryFaker) GenerateData(tx postgre.Transaction) error {
	for i, val := range f.ID {
		id, err := uuid.Parse(val)
		if err != nil {
			return err
		}

		var parentID *string
		if f.ParentID[i] != "" {
			parentID = &f.ParentID[i]
		}

		if _, err := tx.Exec(InsertCategoryQuery, id, parentID, f.Name[i], f.PhotoURL[i]); err != nil {
			return err
		}
	}

	for i := 0; i < f.Size; i++ {
		if _, err := tx.Exec(InsertCategoryQuery, uuid.New(), nil, faker.Name(), "https://cf.shopee.co.id/file/c139370836a9daa649da70876a326b58"); err != nil {
			return err
		}
	}

	return nil
}
