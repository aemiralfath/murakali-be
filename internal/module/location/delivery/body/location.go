package body

import "murakali/internal/model"

type UnprocessableEntity struct {
	Fields map[string]string `json:"fields"`
}

type ProvinceResponse struct {
	Rows []model.Province `json:"rows"`
}

type CityResponse struct {
	Rows []model.City `json:"rows"`
}

type SubDistrictResponse struct {
	Rows []model.SubDistrict `json:"rows"`
}

type UrbanResponse struct {
	Rows []model.Urban `json:"rows"`
}

type StatusResponse struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type KodePosResponse struct {
	Code     int    `json:"code"`
	Status   bool   `json:"status"`
	Messages string `json:"messages"`
	Data     []struct {
		Province    string `json:"province"`
		City        string `json:"city"`
		Subdistrict string `json:"subdistrict"`
		Urban       string `json:"urban"`
		Postalcode  string `json:"postalcode"`
	} `json:"data"`
}

type RajaOngkirProvinceResponse struct {
	RajaOngkir struct {
		Query   []interface{}    `json:"query"`
		Status  StatusResponse   `json:"status"`
		Results []model.Province `json:"results"`
	} `json:"rajaongkir"`
}

type RajaOngkirCityResponse struct {
	Rajaongkir struct {
		Query struct {
			Province string `json:"province"`
		} `json:"query"`
		Status  StatusResponse `json:"status"`
		Results []struct {
			CityID     string `json:"city_id"`
			ProvinceID string `json:"province_id"`
			Province   string `json:"province"`
			Type       string `json:"type"`
			CityName   string `json:"city_name"`
			PostalCode string `json:"postal_code"`
		} `json:"results"`
	} `json:"rajaongkir"`
}
