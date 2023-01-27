package table

import "murakali/pkg/postgre"

type OrderStatusFaker struct {
	Name []string
}

func NewOrderStatusFaker(name []string) ISeeder {
	return &OrderStatusFaker{Name: name}
}

func (f *OrderStatusFaker) GenerateData(tx postgre.Transaction) error {
	const InsertOrderStatusQuery = `INSERT INTO "order_status" (name) VALUES ($1)`

	for _, val := range f.Name {
		_, err := tx.Exec(InsertOrderStatusQuery, val)
		if err != nil {
			return err
		}
	}

	return nil
}
