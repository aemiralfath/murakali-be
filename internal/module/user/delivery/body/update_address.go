package body

import (
	"github.com/google/uuid"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type UpdateAddressRequest struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ProvinceID    int    `json:"province_id"`
	CityID        int    `json:"city_id"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	SubDistrict   string `json:"sub_district"`
	AddressDetail string `json:"address_detail"`
	ZipCode       string `json:"zip_code"`
	IsDefault     bool   `json:"is_default"`
	IsShopDefault bool   `json:"is_shop_default"`
}

func (r *UpdateAddressRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"name":            "",
			"province_id":     "",
			"city_id":         "",
			"province":        "",
			"city":            "",
			"district":        "",
			"sub_district":    "",
			"address_detail":  "",
			"zip_code":        "",
			"is_default":      "",
			"is_shop_default": "",
		},
	}

	r.ID = strings.TrimSpace(r.ID)
	id, err := uuid.Parse(r.ID)
	if err != nil {
		unprocessableEntity = true
		entity.Fields["id"] = IDNotValidMessage
	}

	r.ID = id.String()
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		unprocessableEntity = true
		entity.Fields["name"] = FieldCannotBeEmptyMessage
	}

	if r.ProvinceID == 0 {
		unprocessableEntity = true
		entity.Fields["province_id"] = FieldCannotBeEmptyMessage
	}

	if r.CityID == 0 {
		unprocessableEntity = true
		entity.Fields["city_id"] = FieldCannotBeEmptyMessage
	}

	r.Province = strings.TrimSpace(r.Province)
	if r.Province == "" {
		unprocessableEntity = true
		entity.Fields["province"] = FieldCannotBeEmptyMessage
	}

	r.City = strings.TrimSpace(r.City)
	if r.City == "" {
		unprocessableEntity = true
		entity.Fields["city"] = FieldCannotBeEmptyMessage
	}

	r.District = strings.TrimSpace(r.District)
	if r.District == "" {
		unprocessableEntity = true
		entity.Fields["district"] = FieldCannotBeEmptyMessage
	}

	r.SubDistrict = strings.TrimSpace(r.SubDistrict)
	if r.SubDistrict == "" {
		unprocessableEntity = true
		entity.Fields["sub_district"] = FieldCannotBeEmptyMessage
	}

	r.AddressDetail = strings.TrimSpace(r.AddressDetail)
	if r.AddressDetail == "" {
		unprocessableEntity = true
		entity.Fields["address_detail"] = FieldCannotBeEmptyMessage
	}

	r.ZipCode = strings.TrimSpace(r.ZipCode)
	if r.ZipCode == "" {
		unprocessableEntity = true
		entity.Fields["zip_code"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
