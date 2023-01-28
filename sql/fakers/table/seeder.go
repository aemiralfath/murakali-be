package table

import "murakali/pkg/postgre"

type ISeeder interface {
	GenerateData(tx postgre.Transaction) error
}
