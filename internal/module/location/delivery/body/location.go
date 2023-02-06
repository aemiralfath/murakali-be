package body

import "murakali/internal/model"

const (
	InvalidIDMessage          = "Invalid id."
	FieldCannotBeEmptyMessage = "Field cannot be empty."
)

type UnprocessableEntity struct {
	Fields map[string]interface{} `json:"fields"`
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

type RajaOngkirCostResponse struct {
	Rajaongkir struct {
		Query struct {
			Origin      string `json:"origin,omitempty"`
			Destination string `json:"destination,omitempty"`
			Weight      int    `json:"weight,omitempty"`
			Courier     string `json:"courier,omitempty"`
		} `json:"query,omitempty"`
		Status struct {
			Code        int    `json:"code,omitempty"`
			Description string `json:"description,omitempty"`
		} `json:"status,omitempty"`
		OriginDetails struct {
			CityID     string `json:"city_id,omitempty"`
			ProvinceID string `json:"province_id,omitempty"`
			Province   string `json:"province,omitempty"`
			Type       string `json:"type,omitempty"`
			CityName   string `json:"city_name,omitempty"`
			PostalCode string `json:"postal_code,omitempty"`
		} `json:"origin_details,omitempty"`
		DestinationDetails struct {
			CityID     string `json:"city_id,omitempty"`
			ProvinceID string `json:"province_id,omitempty"`
			Province   string `json:"province,omitempty"`
			Type       string `json:"type,omitempty"`
			CityName   string `json:"city_name,omitempty"`
			PostalCode string `json:"postal_code,omitempty"`
		} `json:"destination_details,omitempty"`
		Results []struct {
			Code  string `json:"code,omitempty"`
			Name  string `json:"name,omitempty"`
			Costs []struct {
				Service     string `json:"service,omitempty"`
				Description string `json:"description,omitempty"`
				Cost        []struct {
					Value int    `json:"value,omitempty"`
					Etd   string `json:"etd,omitempty"`
					Note  string `json:"note,omitempty"`
				} `json:"cost,omitempty"`
			} `json:"costs,omitempty"`
		} `json:"results,omitempty"`
	} `json:"rajaongkir,omitempty"`
}
