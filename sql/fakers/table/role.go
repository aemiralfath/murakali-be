package table

import (
	"murakali/pkg/postgre"
)

type RoleFaker struct {
	Name []string
}

func NewRoleFaker(name []string) ISeeder {
	return &RoleFaker{Name: name}
}

func (f *RoleFaker) GenerateData(tx postgre.Transaction) error {
	const InsertRoleQuery = `INSERT INTO "role" (name) VALUES ($1)`

	for _, val := range f.Name {
		_, err := tx.Exec(InsertRoleQuery, val)
		if err != nil {
			return err
		}
	}

	return nil
}
