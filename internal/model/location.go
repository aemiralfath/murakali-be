package model

type Province struct {
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
}

type City struct {
	CityID string `json:"city_id"`
	City   string `json:"city"`
}

type SubDistrict struct {
	SubDistrict string `json:"sub_district"`
}
