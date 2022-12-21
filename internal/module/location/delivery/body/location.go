package body

import "murakali/internal/model"

const (
	FieldCannotBeEmptyMessage = "Field cannot be empty."
)

type UnprocessableEntity struct {
	Fields map[string]string `json:"fields"`
}

type ProvinceResponse struct {
	Rows []model.Province `json:"rows"`
}

type RajaOngkirProvinceResponse struct {
	RajaOngkir struct {
		Query  []interface{} `json:"query"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
		Results []model.Province `json:"results"`
	} `json:"rajaongkir"`
}
