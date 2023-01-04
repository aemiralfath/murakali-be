package repository

const (
	GetCourierByID               = `SELECT "id", "name", "code", "service", "description" FROM "courier" WHERE "id" = $1`
	GetProductCourierWhitelistID = `SELECT "courier_id" FROM "product_courier_whitelist" WHERE "product_id" = $1`
	GetShopCourierID             = `SELECT "courier_id" FROM "shop_courier" WHERE "shop_id" = $1`
)
