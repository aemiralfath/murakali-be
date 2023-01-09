package repository

const (
	GetCourierByID               = `SELECT "id", "name", "code", "service", "description", "created_at", "updated_at" FROM "courier" WHERE "id" = $1 AND "deleted_at" IS NULL`
	GetProductCourierWhitelistID = `SELECT "courier_id" FROM "product_courier_whitelist" WHERE "product_id" = $1 AND "deleted_at" IS NULL`
	GetShopCourierID             = `SELECT "courier_id" FROM "shop_courier" WHERE "shop_id" = $1 AND "deleted_at" IS NULL`
	GetShopByID                  = `SELECT "id", "user_id" FROM "shop" WHERE "id" = $1 AND "deleted_at" IS NULL`
	GetShopAddress               = `SELECT "id", "user_id", "city_id" FROM "address" WHERE "user_id" = $1 AND "is_shop_default" IS true AND "deleted_at" IS NULL`
)
