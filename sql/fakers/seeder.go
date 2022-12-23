package fakers

import (
	"database/sql"
	"murakali/pkg/postgre"
	"murakali/sql/fakers/table"
)

type Seeder struct {
	Seeder table.ISeeder
}

func RegisterSeeders() []Seeder {
	return []Seeder{
		{Seeder: table.NewRoleFaker([]string{"user", "seller", "admin"})},
		{Seeder: table.NewUserFaker(
			0,
			1,
			"M",
			[]string{"4c1d6464-3cc6-44d6-92d1-91aee337e025", "0c53ef3d-3682-4359-90e1-814eb6ab5191", "7950eca2-58d5-44f0-b873-22b23d8107da"},
			[]string{"fadhlan1337@gmail.com", "sammymanunggal@gmail.com", "user@gmail.com"})},
		{Seeder: table.NewUserFaker(0, 2, "M", []string{"f8d8d66a-e8eb-4633-bc2d-4ccd941fed47"}, []string{"seller@gmail.com"})},
		{Seeder: table.NewUserFaker(0, 3, "M", []string{"4df967a8-5b05-4d2a-bb72-da3921dce8fb"}, []string{"admin@gmail.com"})},
		{Seeder: table.NewUserFaker(10000, 1, "M", []string{}, []string{})},
		{Seeder: table.NewUserFaker(10000, 1, "F", []string{}, []string{})},
	}
}

func DBSeed(sqlDB *sql.DB) error {
	txDB := postgre.NewTxRepository(sqlDB)
	for _, seeder := range RegisterSeeders() {
		err := txDB.WithTransaction(func(transaction postgre.Transaction) error {
			if errSeeder := seeder.Seeder.GenerateData(transaction); errSeeder != nil {
				return errSeeder
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
