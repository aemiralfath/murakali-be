package model

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Address struct {
	ID            uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	UserID        uuid.UUID    `json:"user_id" db:"user_id" binding:"omitempty"`
	Name          string       `json:"name" db:"name" binding:"omitempty"`
	ProvinceID    int          `json:"province_id" db:"province_id" binding:"omitempty"`
	CityID        int          `json:"city_id" db:"city_id" binding:"omitempty"`
	Province      string       `json:"province" db:"province" binding:"omitempty"`
	City          string       `json:"city" db:"city" binding:"omitempty"`
	District      string       `json:"district" db:"district" binding:"omitempty"`
	SubDistrict   string       `json:"sub_district" db:"sub_district" binding:"omitempty"`
	AddressDetail string       `json:"address_detail" db:"address_detail" binding:"omitempty"`
	ZipCode       string       `json:"zip_code" db:"zip_code" binding:"omitempty"`
	IsDefault     bool         `json:"is_default" db:"is_default" binding:"omitempty"`
	IsShopDefault bool         `json:"is_shop_default" db:"is_shop_default" binding:"omitempty"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at" binding:"omitempty"`
	UpdatedAt     sql.NullTime `json:"updated_at" db:"updated_at" binding:"omitempty"`
	DeletedAt     sql.NullTime `json:"deleted_at" db:"deleted_at" binding:"omitempty"`
}
