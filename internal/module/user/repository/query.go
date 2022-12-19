package repository

const (
	GetUserByIDQuery              = `SELECT "id", "role_id", "email" FROM "user" WHERE "id" = $1 AND "deleted_at" IS NULL`
	GetTotalAddressQuery          = `SELECT count(id) FROM "address" WHERE "user_id" = $1 AND "name" ILIKE $2 AND "deleted_at" IS NULL`
	GetDefaultAddressQuery        = `SELECT "id", "user_id", "is_default" FROM "address" WHERE "user_id" = $1 AND "is_default" = $2 AND "deleted_at" IS NULL`
	GetDefaultShopAddressQuery    = `SELECT "id", "user_id", "is_shop_default" FROM "address" WHERE "user_id" = $1 AND "is_shop_default" = $2 AND "deleted_at" IS NULL`
	UpdateDefaultAddressQuery     = `UPDATE "address" SET "is_default" = $1 WHERE "id" = $2`
	UpdateDefaultShopAddressQuery = `UPDATE "address" SET "is_shop_default" = $1 WHERE "id" = $2`
	DeleteAddressByIDQuery        = `DELETE FROM "address" WHERE "id" = $1`

	CreateAddressQuery = `INSERT INTO "address" 
    	(user_id, name, province_id, city_id, province, city, district, sub_district, address_detail, zip_code, is_default, is_shop_default)
    	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	GetAddressesQuery = `SELECT 
    	"id", "user_id", "name", "province_id", "city_id", "province", "city", "district", "sub_district",  
    	"address_detail", "zip_code", "is_default", "is_shop_default", "created_at", "updated_at" 
	FROM "address" WHERE "user_id" = $1 AND "name" ILIKE $2 AND "deleted_at" IS NULL ORDER BY $3 LIMIT $4 OFFSET $5`

	GetAddressByIDQuery = `SELECT
		"id", "user_id", "name", "province_id", "city_id", "province", "city", "district", "sub_district",  
    	"address_detail", "zip_code", "is_default", "is_shop_default", "created_at", "updated_at"
	FROM "address" WHERE "id" = $1 AND "user_id" = $2 AND "deleted_at" IS NULL`

	UpdateAddressByIDQuery = `UPDATE "address" SET
		"name" = $1, "province_id" = $2, "city_id" = $3, "province" = $4, "city" = $5, "district" = $6,
		"sub_district" = $7, "address_detail" = $8, "zip_code" = $9, "is_default" = $10, "is_shop_default" = $11, "updated_at" = $12
	WHERE "id" = $13`
)
