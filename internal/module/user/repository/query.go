package repository

const (
	GetTotalAddressQuery        = `SELECT count(id) FROM "address" WHERE "user_id" = $1 AND "name" ILIKE $2 AND "deleted_at" IS NULL`
	GetTotalAddressDefaultQuery = `SELECT count(id) FROM "address" WHERE "user_id" = $1 AND "name" ILIKE $2 AND "deleted_at" IS NULL AND "is_default" = $3 AND "is_shop_default" = $4`
	GetDefaultAddressQuery      = `
		SELECT "id", "user_id", "name", "province_id", "city_id", "province", "city", "district", "sub_district",  
			"address_detail", "zip_code", "is_default", "is_shop_default", "created_at", "updated_at" 
		FROM "address" WHERE "user_id" = $1 AND "is_default" = $2 AND "deleted_at" IS NULL
	`
	GetDefaultShopAddressQuery = `
		SELECT "id", "user_id", "name", "province_id", "city_id", "province", "city", "district", "sub_district",  
		"address_detail", "zip_code", "is_default", "is_shop_default", "created_at", "updated_at" 
		FROM "address" WHERE "user_id" = $1 AND "is_shop_default" = $2 AND "deleted_at" IS NULL
	`
	UpdateDefaultAddressQuery     = `UPDATE "address" SET "is_default" = $1 WHERE "id" = $2`
	UpdateDefaultShopAddressQuery = `UPDATE "address" SET "is_shop_default" = $1 WHERE "id" = $2`
	DeleteAddressByIDQuery        = `DELETE FROM "address" WHERE "id" = $1`

	CreateAddressQuery = `INSERT INTO "address" 
    	(user_id, name, province_id, city_id, province, city, district, sub_district, address_detail, zip_code, is_default, is_shop_default)
    	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	GetAddressesDefaultQuery = `SELECT 
    	"id", "user_id", "name", "province_id", "city_id", "province", "city", "district", "sub_district",  
    	"address_detail", "zip_code", "is_default", "is_shop_default", "created_at", "updated_at" 
	FROM "address" WHERE "user_id" = $1 AND "name" ILIKE $2 AND "deleted_at" IS NULL AND "is_default" = $3 AND "is_shop_default" = $4 ORDER BY $5 LIMIT $6 OFFSET $7`

	GetAllAddressesQuery = `SELECT 
    	"id", "user_id", "name", "province_id", "city_id", "province", "city", "district", "sub_district",  
    	"address_detail", "zip_code", "is_default", "is_shop_default", "created_at", "updated_at" 
	FROM "address" WHERE "user_id" = $1 AND "name" ILIKE $2 AND "deleted_at" IS NULL ORDER BY $3 LIMIT $4 OFFSET $5`

	GetTotalOrderQuery = `SELECT count(id) FROM "order" WHERE "user_id" = $1 and "order_status_id"::text LIKE $2`

	GetOrdersQuery = `SELECT o.id,o.order_status_id,o.total_price,o.delivery_fee,o.resi_no,s.id,s.name,v.code,o.created_at
	from "order" o
	join "shop" s on s.id = o.shop_id
	left join "voucher" v on v.id = o.voucher_shop_id 
	WHERE o.user_id = $1 
	and "order_status_id"::text LIKE $2 
	ORDER BY o.created_at asc LIMIT $3 OFFSET $4
	`

	GetOrderDetailProductVariant = `
		SELECT "vd"."name" as "name", "vd"."type" as "type" 
		FROM "variant_detail" as "vd"
		INNER JOIN "variant" as "v" ON "v"."variant_detail_id" = "vd"."id"
		INNER JOIN "product_detail" as "pd" ON "pd"."id" = "v"."product_detail_id"
		WHERE "pd"."id" = $1 AND "pd"."deleted_at" IS NULL
	`

	GetTotalTransactionByUserIDQuery = `SELECT count(t.id) FROM "transaction" t, "order" o
	WHERE t.id = o.transaction_id
	AND o.user_id = $1
	GROUP BY o.user_id`

	GetTransactionByUserIDQuery = `
	SELECT t.id,t.voucher_marketplace_id,t.wallet_id,t.card_number,t.invoice,t.total_price,t.paid_at,t.canceled_at,t.expired_at
	from "transaction" t, "order" o
	WHERE t.id = o.transaction_id
	AND o.user_id = $1
	GROUP BY t.id
	ORDER BY t.expired_at DESC LIMIT $2 OFFSET $3
	`

	GetTransactionByUserIDNotPaidQuery = `
	SELECT t.id,t.voucher_marketplace_id,t.wallet_id,t.card_number,t.invoice,t.total_price,t.paid_at,t.canceled_at,t.expired_at
	from "transaction" t, "order" o
	WHERE t.id = o.transaction_id
	AND o.user_id = $1
	AND t.paid_at IS NULL AND t.canceled_at IS NULL
	GROUP BY t.id
	ORDER BY t.expired_at DESC LIMIT $2 OFFSET $3
	`

	GetOrdersByTransactionIDQuery = `SELECT o.id,o.order_status_id,o.total_price,o.delivery_fee,o.resi_no,s.id,s.name,v.code,o.created_at
	from "order" o
	join "shop" s on s.id = o.shop_id
	left join "voucher" v on v.id = o.voucher_shop_id WHERE o.transaction_id = $1
	`

	GetOrderDetailQuery = `SELECT pd.id,pd.product_id,p.title,ph.url,oi.quantity,oi.item_price,oi.total_price
	from  "product_detail" pd 
	join "photo" ph on pd.id = ph.product_detail_id join "order_item" oi on pd.id = oi.product_detail_id 
	join "product" p on p.id = pd.product_id WHERE oi.order_id = $1`

	GetAddressByIDQuery = `SELECT
		"id", "user_id", "name", "province_id", "city_id", "province", "city", "district", "sub_district",  
    	"address_detail", "zip_code", "is_default", "is_shop_default", "created_at", "updated_at"
	FROM "address" WHERE "id" = $1 AND "user_id" = $2 AND "deleted_at" IS NULL`

	UpdateAddressByIDQuery = `UPDATE "address" SET
		"name" = $1, "province_id" = $2, "city_id" = $3, "province" = $4, "city" = $5, "district" = $6,
		"sub_district" = $7, "address_detail" = $8, "zip_code" = $9, "is_default" = $10, "is_shop_default" = $11, "updated_at" = $12
	WHERE "id" = $13`

	GetSealabsPayByIdQuery         = `SELECT * from sealabs_pay where user_id = $1 and deleted_at is null`
	CreateSealabsPayQuery          = `INSERT INTO "sealabs_pay" (card_number, user_id, name, is_default,active_date) VALUES ($1, $2, $3, $4, $5)`
	CheckDefaultSealabsPayQuery    = `SELECT card_number from "sealabs_pay" where user_id = $1 and is_default is true and deleted_at is null`
	SetDefaultSealabsPayTransQuery = `UPDATE "sealabs_pay" set is_default = FALSE,updated_at = now() where card_number = $1`
	PatchSealabsPayQuery           = `UPDATE "sealabs_pay" set is_default = TRUE,updated_at = now() where card_number = $1`
	SetDefaultSealabsPayQuery      = `UPDATE "sealabs_pay" set is_default = FALSE where card_number <> $1 and user_id = $2`
	DeleteSealabsPayQuery          = `UPDATE "sealabs_pay" set "deleted_at" = now() where "card_number" = $1 and "is_default" = FALSE`
	GetUserByIDQuery               = `SELECT "id", "role_id", "email", "username", "phone_no", "fullname", "gender", "birth_date", "is_verify","photo_url" FROM "user" WHERE "id" = $1`
	GetUserPasswordByIDQuery       = `SELECT "id", "role_id", "email", "password", "username", "phone_no", "fullname", "gender", "birth_date", "is_verify","photo_url" FROM "user" WHERE "id" = $1`
	GetPasswordByIDQuery           = `SELECT "password" FROM "user" WHERE "id" = $1`
	CheckEmailHistoryQuery         = `SELECT "id", "email" FROM "email_history" WHERE "email" ILIKE $1`
	GetUserByUsernameQuery         = `SELECT "id", "email", "username", "is_verify" FROM "user" WHERE "username" ILIKE $1`
	GetUserByPhoneNoQuery          = `SELECT "id", "email", "phone_no", "is_verify" FROM "user" WHERE "phone_no" ILIKE $1`
	UpdateUserFieldQuery           = `UPDATE "user" SET "username" = $1, "fullname" = $2, "phone_no" = $3, "birth_date" = $4, "gender" = $5, "updated_at" = $6 WHERE "email" = $7`
	UpdateUserEmailQuery           = `UPDATE "user" SET "email" = $1, "updated_at" = $2 WHERE "id" = $3`
	CreateEmailHistoryQuery        = `INSERT INTO "email_history" (email) VALUES ($1)`
	CheckShopByIdQuery             = `SELECT count(id) from "shop" WHERE "user_id" = $1 and deleted_at IS NULL`
	CheckShopUniqueQuery           = `SELECT count(name) from "shop" WHERE "name" = $1 and deleted_at IS NULL`
	AddShopQuery                   = `INSERT INTO "shop" (user_id,name) VALUES ($1,$2) `
	UpdateRoleQuery                = `UPDATE "user" SET "role_id" = 2,updated_at = now() where id = $1`
	UpdateProfileImageQuery        = `UPDATE "user" SET "photo_url" = $1,updated_at = now() where id = $2`
	UpdatePasswordQuery            = `UPDATE "user" SET "password" = $1 WHERE "id" = $2`
	UpdateWalletPinQuery           = `UPDATE "wallet" SET "pin" = $1 WHERE "id" = $2`

	GetWalletUserQuery        = `SELECT "id", "user_id", "balance", "attempt_count", "attempt_at", "unlocked_at", "active_date" FROM "wallet" WHERE "id" = $1 AND "deleted_at" IS NULL;`
	GetWalletHistoryUserQuery = `SELECT "id", "from", "to", "amount", "description", "created_at" 
	FROM "wallet_history" 
	WHERE "wallet_id" = $1`
	GetWalletHistoryByIDQuery = `SELECT "id", "transaction_id", "wallet_id", "from", "to", "amount", "description", "created_at" 
	FROM "wallet_history" 
	WHERE "id" = $1`
	GetTotalWalletHistoryUserQuery = `SELECT count(id) FROM "wallet_history" WHERE "wallet_id" = $1;`
	GetSealabsPayUserQuery         = `SELECT "card_number", "user_id", "name", "is_default", "active_date" FROM "sealabs_pay" WHERE "user_id" = $1 AND "card_number" = $2 AND "deleted_at" IS NULL;`
	GetVoucherMarketplaceByIDQuery = `SELECT "id", "shop_id", "code", "quota", "actived_date", "expired_date", "discount_percentage", "discount_fix_price", "min_product_price", "max_discount_price" FROM "voucher"
		WHERE "id" = $1 AND "shop_id" is NULL AND "deleted_at" IS NULL AND now() BETWEEN "actived_date" AND "expired_date";`
	GetVoucherShopByIDQuery = `SELECT "id", "shop_id", "code", "quota", "actived_date", "expired_date", "discount_percentage", "discount_fix_price", "min_product_price", "max_discount_price" FROM "voucher"
		WHERE "id" = $1 AND "shop_id" = $2 AND "deleted_at" IS NULL AND now() BETWEEN "actived_date" AND "expired_date";`
	GetCourierShopByIDQuery = `SELECT "c"."id", "c"."name", "c"."code", "c"."service", "c"."description" FROM "courier" as "c"
		INNER JOIN "shop_courier" as sc ON "sc"."courier_id" = "c"."id"
		WHERE "c"."id" = $1 AND "sc"."shop_id" = $2 AND "c"."deleted_at" IS NULL;`
	GetProductDetailByIDQuery     = `SELECT "id", "product_id", "price", "stock", "weight", "size", "hazardous", "condition", "bulk_price" FROM "product_detail" WHERE "id" = $1 AND "deleted_at" IS NULL;`
	GetShopByIDQuery              = `SELECT "id", "name" FROM "shop" WHERE "id" = $1 AND "deleted_at" IS NULL;`
	CreateTransactionQuery        = `INSERT INTO "transaction" (voucher_marketplace_id, wallet_id, card_number, invoice, total_price, expired_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING "id";`
	CreateOrderQuery              = `INSERT INTO "order" (transaction_id, shop_id, user_id, courier_id, voucher_shop_id, order_status_id, total_price, delivery_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING "id";`
	CreateOrderItemQuery          = `INSERT INTO "order_item" (order_id, product_detail_id, quantity, item_price, total_price) VALUES ($1, $2, $3, $4, $5) RETURNING "id";`
	CreateWalletQuery             = `INSERT INTO "wallet" (user_id, balance, pin, attempt_count, active_date) VALUES ($1, $2, $3, $4, $5)`
	CreateWalletHistoryQuery      = `INSERT INTO "wallet_history" (transaction_id, wallet_id, "from", "to", description, amount, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	UpdateWalletBalanceQuery      = `UPDATE "wallet" SET "balance" = $1, "updated_at" = $2 WHERE "id" = $3`
	UpdateWalletQuery             = `UPDATE "wallet" SET "attempt_count" = $1, "attempt_at" = $2, "unlocked_at" = $3, "updated_at" = CURRENT_TIMESTAMP WHERE "id" = $4`
	GetWalletByUserIDQuery        = `SELECT "id", "user_id", "balance", "pin", "attempt_count", "attempt_at", "unlocked_at", "active_date" FROM "wallet" WHERE "user_id" = $1 AND "deleted_at" IS NULL`
	GetCartItemUserQuery          = `SELECT "id", "user_id", "product_detail_id", "quantity" FROM "cart_item" WHERE "user_id" = $1 AND "product_detail_id" = $2 AND "deleted_at" IS NULL;`
	UpdateProductDetailStockQuery = `UPDATE "product_detail" SET "stock" = $1, "updated_at" = now() WHERE "id" = $2;`
	DeleteCartItemByIDQuery       = `DELETE FROM "cart_item" WHERE "id" = $1`
	GetTransactionByIDQuery       = `SELECT "id", "voucher_marketplace_id", "wallet_id", "card_number", "invoice", "total_price", "paid_at", "canceled_at", "expired_at" FROM "transaction" WHERE "id" = $1;`
	UpdateTransactionByID         = `UPDATE "transaction" SET "paid_at" = $1, "canceled_at" = $2, "card_number" = $3 WHERE "id" = $4`
	UpdateOrderByID               = `UPDATE "order" SET "order_status_id" = $1 WHERE "id" = $2`
	GetOrderByTransactionID       = `SELECT 
		"id", "transaction_id", "shop_id", "user_id", "courier_id", "voucher_shop_id", "order_status_id", "total_price", "delivery_fee", "resi_no", "created_at", "arrived_at" 
	FROM "order" WHERE "transaction_id" = $1`
	GetOrderByTransactionIDUserID = `SELECT 
	"id", "transaction_id", "shop_id", "user_id", "courier_id", "voucher_shop_id", "order_status_id", "total_price", "delivery_fee", "resi_no", "created_at", "arrived_at" 
	FROM "order" WHERE "transaction_id" = $1
	AND "user_id" = $2`
	CheckUserSealabsPayQuery    = `SELECT count(1) from sealabs_pay where user_id = $1 and deleted_at is null`
	CheckDeletedSealabsPayQuery = `SELECT count(1) from sealabs_pay where card_number = $1 and deleted_at is not null`
	UpdateUserSealabsPayQuery   = `UPDATE "sealabs_pay" set user_id = $1, name = $2 ,updated_at = now(),deleted_at = null,is_default = true where card_number = $3`
	OrderBySomething            = ` 
	ORDER BY %s LIMIT %d OFFSET %d`

	GetOrderDetailQuery2 = `SELECT pd.id,pd.product_id,p.title, pd.weight,
	(select ph.url from "photo" ph 
			join product_detail pd on pd.id = ph.product_detail_id 
			join "order_item" oi on pd.id = oi.product_detail_id limit 1
		),oi.quantity,oi.item_price,oi.total_price
	from  "product_detail" pd 
	join "order_item" oi on pd.id = oi.product_detail_id 
	join "product" p on p.id = pd.product_id WHERE oi.order_id = $1 `

	GetOrderByOrderID = `SELECT o.id,o.order_status_id,o.total_price,o.delivery_fee,o.resi_no,s.id,s.name,u2.phone_no,u2.username,v.code,o.created_at,t.invoice
	,c.name,c.code,c.service,c.description,u.username,u.phone_no
	from "order" o
	join "shop" s on s.id = o.shop_id
	join "courier" c on o.courier_id = c.id
	join "user" u on o.user_id = u.id
	join "user" u2 on s.user_Id = u2.id
	join transaction t on o.transaction_id = t.id
	left join "voucher" v on v.id = o.voucher_shop_id WHERE o.id = $1 ORDER BY o.created_at asc`

	GetProductUnitSoldByOrderIDQuery = `
	SELECT 
  "product".id as product_id, 
  "product".unit_sold,
  "order_item".quantity
		FROM "order_item"
		JOIN "product_detail" ON "order_item"."product_detail_id" = "product_detail"."id"
		JOIN "product" ON "product_detail"."product_id" = "product"."id"
		WHERE "order_item"."order_id" = $1
`

	GetBuyerIDByOrderIDQuery = `SELECT o.user_id from "order" o where o.id = $1`

	GetSellerIDByOrderIDQuery = `SELECT s.user_id from "order" o join shop s on o.shop_id = s.id where o.id = $1`

	GetAddressByBuyerIDQuery = `SELECT
	"id", "user_id", "name", "province_id", "city_id", "province", "city", "district", "sub_district",  
	"address_detail", "zip_code", "is_default", "is_shop_default", "created_at", "updated_at"
	FROM "address" WHERE "user_id" = $1 AND "deleted_at" IS NULL AND is_default is true`

	GetAddressBySellerIDQuery = `SELECT
	"id", "user_id", "name", "province_id", "city_id", "province", "city", "district", "sub_district",  
	"address_detail", "zip_code", "is_default", "is_shop_default", "created_at", "updated_at"
	FROM "address" WHERE "user_id" = $1 AND "deleted_at" IS NULL AND is_shop_default is true`

	GetProductPromotionByProductIDQuery = `
	SELECT "promo"."id", "promo"."name", "promo"."product_id", "promo"."discount_percentage", "promo"."discount_fix_price",
		"promo"."min_product_price", "promo"."max_discount_price", "promo"."quota", "promo"."max_quantity", "promo"."actived_date",
		"promo"."expired_date", "promo"."created_at", "promo"."updated_at", "promo"."deleted_at"
	FROM "promotion" as "promo" 
	INNER JOIN "product" as "p" ON "p"."id" = "promo"."product_id"
	WHERE "p"."id" = $1 AND ("promo"."actived_date" < now() AND "promo"."expired_date" >= now())
	`

	UpdateVoucherQuotaQuery    = `UPDATE "voucher" SET "quota" = $1, "updated_at" = now() WHERE "id" = $2;`
	UpdatePromotionQuotaQuery  = `UPDATE "promotion" SET "quota" = $1, "updated_at" = now() WHERE "id" = $2;`
	UpdateProductUnitSoldQuery = `UPDATE "product" SET "unit_sold" = $1 "updated_at" = now() WHERE "id" = $2;`
)
