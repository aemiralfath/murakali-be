package repository

const (
	GetUserByIDQuery = `SELECT "id", "role_id", "email" FROM "user" WHERE "id" = $1`
	GetTotalAddress  = `SELECT count(id) FROM "address" WHERE "user_id" = $1 AND "name" ILIKE $2`
	GetAddresses     = `SELECT 
    	"id", "user_id", "name", "province_id", "city_id", "province", "city", "district", "sub_district",  
    	"address_detail", "zip_code", "is_default", "is_shop_default", "created_at", "updated_at" 
	FROM "address" WHERE "user_id" = $1 AND "name" ILIKE $2 ORDER BY $3 LIMIT $4 OFFSET $5`
)
