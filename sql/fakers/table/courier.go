package table

import (
	"github.com/google/uuid"
	"murakali/pkg/postgre"
)

type CourierFaker struct {
	ID      []string
	Name    []string
	Code    []string
	Service []string
}

func NewCourierFaker(id, name, code, service []string) ISeeder {
	return &CourierFaker{ID: id, Name: name, Code: code}
}

func (f *CourierFaker) GenerateData(tx postgre.Transaction) error {
	const InsertCourierQuery = `INSERT INTO "courier" (id, name, code, service, description) VALUES ($1, $2, $3, $4, $5)`
	for i, val := range f.ID {
		id, err := uuid.Parse(val)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(InsertCourierQuery, id, f.Name[i], f.Code[i], f.Service[i], "Layanan Reguler"); err != nil {
			return err
		}
	}
	return nil
}
