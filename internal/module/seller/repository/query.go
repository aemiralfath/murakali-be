package repository

const (
	GetTotalOrderQuery = `SELECT count(id) FROM "order" as "o" WHERE "o"."shop_id" = $1 and "o"."order_status_id"::text LIKE $2`

	GetTotalOrderWithVoucherIDQuery = `SELECT count(id) FROM "order" as "o" WHERE "o"."shop_id" = $1 and "o"."order_status_id"::text LIKE $2 
	AND "o"."voucher_shop_id" = $3
	`

	GetOrdersQuery = `SELECT o.id,o.order_status_id,o.total_price,o.delivery_fee,o.resi_no,s.id,s.name,v.code,o.created_at
	from "order" o
	join "shop" s on s.id = o.shop_id
	left join "voucher" v on v.id = o.voucher_shop_id 
	WHERE o.shop_id = $1 
	and "order_status_id"::text LIKE $2 
	ORDER BY o.created_at asc LIMIT $3 OFFSET $4
	`

	GetOrdersWithVoucherIDQuery = `SELECT o.id,o.order_status_id,o.total_price,o.delivery_fee,o.resi_no,s.id,s.name,v.code,o.created_at
	from "order" o
	join "shop" s on s.id = o.shop_id
	left join "voucher" v on v.id = o.voucher_shop_id 
	WHERE o.shop_id = $1 
	and "order_status_id"::text LIKE $2 
	AND "o"."voucher_shop_id" = $3
	ORDER BY o.created_at asc LIMIT $4 OFFSET $5
	`

	GetAddressByBuyerIDQuery = `SELECT
	"id", "user_id", "name", "province_id", "city_id", "province", "city", "district", "sub_district",  
	"address_detail", "zip_code", "is_default", "is_shop_default", "created_at", "updated_at"
	FROM "address" WHERE "user_id" = $1 AND "deleted_at" IS NULL AND is_default is true`

	GetAddressBySellerIDQuery = `SELECT
	"id", "user_id", "name", "province_id", "city_id", "province", "city", "district", "sub_district",  
	"address_detail", "zip_code", "is_default", "is_shop_default", "created_at", "updated_at"
	FROM "address" WHERE "user_id" = $1 AND "deleted_at" IS NULL AND is_shop_default is true`

	GetOrderByOrderID = `SELECT o.id,o.order_status_id,o.total_price,o.delivery_fee,o.resi_no,s.id,s.name,u2.phone_no,u2.username,v.code,o.created_at,t.invoice
	,c.name,c.code,c.service,c.description,u.username,u.phone_no
	from "order" o
	join "shop" s on s.id = o.shop_id
	join "courier" c on o.courier_id = c.id
	join "user" u on o.user_id = u.id
	join "user" u2 on s.user_Id = u2.id
	join transaction t on o.transaction_id = t.id
	left join "voucher" v on v.id = o.voucher_shop_id WHERE o.id = $1 ORDER BY o.created_at asc`

	GetBuyerIDByOrderIDQuery = `SELECT o.user_id from "order" o where o.id = $1`

	GetSellerIDByOrderIDQuery = `SELECT s.user_id from "order" o join shop s on o.shop_id = s.id where o.id = $1`

	GetOrderDetailQuery = `SELECT pd.id,pd.product_id,p.title, pd.weight,
		(select ph.url from "photo" ph 
			join product_detail pd on pd.id = ph.product_detail_id 
			join "order_item" oi on pd.id = oi.product_detail_id limit 1
		),oi.quantity,oi.item_price,oi.total_price
	from  "product_detail" pd 
	join "order_item" oi on pd.id = oi.product_detail_id 
	join "product" p on p.id = pd.product_id WHERE oi.order_id = $1 `

	GetOrderDetailProductVariant = `
		SELECT "vd"."name" as "name", "vd"."type" as "type" 
		FROM "variant_detail" as "vd"
		INNER JOIN "variant" as "v" ON "v"."variant_detail_id" = "vd"."id"
		INNER JOIN "product_detail" as "pd" ON "pd"."id" = "v"."product_detail_id"
		WHERE "pd"."id" = $1 AND "pd"."deleted_at" IS NULL
	`

	GetShopIDByUserQuery = `SELECT id from shop where user_id = $1 and deleted_at is null`

	GetShopIDByOrderQuery = `SELECT shop_id from "order" where id = $1 `

	ChangeOrderStatusQuery = `UPDATE "order" SET "order_status_id" = $1 WHERE "id" = $2`

	GetCourierSellerQuery = `
	SELECT "sp"."id" as "shop_courier_id",	"sp"."courier_id" as "courier_id", "sp"."deleted_at" as "deleted_at"
	FROM "shop_courier" as "sp"
	INNER JOIN "shop" as "s" ON "s"."id" = "sp"."shop_id"
	WHERE "s"."user_id" = $1;
	`

	GetOrderOnDeliveryQuery = `SELECT "id", "order_status_id", "arrived_at" FROM "order" WHERE "order_status_id" = $1 AND "arrived_at" <= current_timestamp`

	GetAllCourierQuery = `
	SELECT  "c"."id" as "courier_id","c"."name" as "name", "c"."code" as "code", "c"."service" as "service",
		"c"."description" as "description"
	FROM "courier" as "c"
	WHERE "c".deleted_at IS NULL;
	`

	GetShopIDByShopIDQuery = `SELECT s.id, s.user_id, s.name, s.total_product,
	 s.total_rating, s.rating_avg, s.created_at, u.photo_url 
	FROM "shop" s 
	JOIN "user" u ON u.id = s.user_id
	WHERE s.id = $1 AND s.deleted_at is null`

	GetShopDetailIDByUserIDQuery = `SELECT s.id, s.user_id, s.name, s.total_product,
	 s.total_rating, s.rating_avg, s.created_at, u.photo_url 
	FROM "shop" s 
	JOIN "user" u ON u.id = s.user_id
	WHERE s.user_id = $1 AND s.deleted_at is null`

	GetCourierByIDQuery                            = `SELECT id FROM "courier" WHERE id = $1 AND deleted_at IS NULL`
	GetShopIDByUserIDQuery                         = `SELECT id from "shop" WHERE user_id = $1 AND deleted_at IS NULL `
	GetCourierSellerNotNullByShopAndCourierIDQuery = `SELECT id from "shop_courier" WHERE shop_id = $1 AND courier_id = $2 `
	CreateCourierSellerQuery                       = `INSERT INTO "shop_courier" 
    	(shop_id, courier_id)
    	VALUES ($1, $2)`

	GetCourierSellerByIDQuery = `SELECT id FROM "shop_courier" WHERE id = $1 AND deleted_at IS NULL`
	DeleteCourierSellerQuery  = `UPDATE "shop_courier" set deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`

	UpdateCourierSellerQuery = `UPDATE "shop_courier" set deleted_at = NULL WHERE courier_id = $1 AND shop_id = $2 AND deleted_at IS NOT NULL`

	GetCategoryBySellerIDQuery = `SELECT c.id, c.name
	From shop s, product p, category c
	where s.id = p.shop_id
	and p.category_id = c.id
	and s.id = $1
	and c.deleted_at is null
	group by c.id`

	UpdateResiNumberInOrderSellerQuery = `UPDATE
	"order" set resi_no = $1, arrived_at = $2, order_status_id = $3 WHERE id = $4 
	AND shop_id = $5`

	GetTotalVoucherSellerQuery = `
	SELECT count(id) FROM "voucher" as "v" WHERE "v"."shop_id" = $1
	`
	GetAllVoucherSellerQuery = `
	SELECT "v"."id", "v"."shop_id", "v"."code", "v"."quota", "v"."actived_date", "v"."expired_date",
		"v"."discount_percentage", "v"."discount_fix_price", "v"."min_product_price", "v"."max_discount_price",
		"v"."created_at", "v"."updated_at",  "v"."deleted_at"
	FROM "voucher" as "v"
	INNER JOIN "shop" as "s" ON "s"."id" = "v"."shop_id"
	WHERE "v"."shop_id" = $1
	AND "v"."deleted_at" IS NULL
	ORDER BY "v"."created_at" DESC
	LIMIT $2 OFFSET $3
	
	`
	FilterVoucherOngoing = `
	 AND  ("v"."actived_date" <= now() AND "v"."expired_date" >= now())`

	FilterVoucherWillCome = `
	 AND (now() < "v"."actived_date" AND  now() < "v"."expired_date") `

	FilterVoucherHasEnded = `
	 AND (now() > "v"."actived_date" AND  now() > "v"."expired_date")  `

	CreateVoucherSellerQuery = `INSERT INTO "voucher" 
    	(shop_id, code, quota, actived_date, expired_date, discount_percentage, discount_fix_price, min_product_price, max_discount_price)
    	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	DeleteVoucherSellerQuery = `UPDATE "voucher" set deleted_at = now() WHERE "id" = $1 AND "shop_id" = $2 AND "deleted_at" IS NULL`

	GetAllVoucherSellerByIDandShopIDQuery = `
	SELECT "v"."id", "v"."shop_id", "v"."code", "v"."quota", "v"."actived_date", "v"."expired_date",
		"v"."discount_percentage", "v"."discount_fix_price", "v"."min_product_price", "v"."max_discount_price",
		"v"."created_at", "v"."updated_at",  "v"."deleted_at"
	FROM "voucher" as "v"
	INNER JOIN "shop" as "s" ON "s"."id" = "v"."shop_id"
	WHERE "v"."id"  = $1 AND "v"."shop_id" = $2 AND "v"."deleted_at" IS NULL
	`

	UpdateVoucherSellerQuery = `
		UPDATE "voucher" SET "quota" = $1, "actived_date" = $2, "expired_date" = $3, "discount_percentage" = $4,
			"discount_fix_price" = $5, "min_product_price" = $6, "max_discount_price" = $7, "updated_at" = now()
		WHERE "id" = $8
	`
	GetAllPromotionSellerQuery = `
	SELECT "promo"."id", "promo"."name", "p"."id", "p"."title", "p"."thumbnail_url", "promo"."discount_percentage",
		"promo"."discount_fix_price", "promo"."min_product_price", "promo"."max_discount_price", "promo"."quota", "promo"."max_quantity", 
		"promo"."actived_date", "promo"."expired_date", "promo"."created_at", "promo"."updated_at", "promo"."deleted_at"
	FROM "promotion" as "promo"
	INNER JOIN "product" as "p" ON "p"."id" = "promo"."product_id"
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	WHERE "s"."id" = $1
	`

	GetTotalPromotionSellerQuery = `
	SELECT count("promo"."id") FROM "promotion" as "promo" 
	INNER JOIN "product" as "p" ON "p"."id" = "promo"."product_id"
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	WHERE "s"."id" = $1
	`

	FilterWillComeQuery = ` AND ("promo"."actived_date" > now() AND "promo"."expired_date" > now())`
	FilterOngoingQuery  = ` AND ("promo"."actived_date" < now() AND "promo"."expired_date" > now())`
	FilterHasEndedQuery = ` AND ("promo"."actived_date" < now() AND "promo"."expired_date" < now())`

	GetProductPromotionQuery = `
	SELECT "p"."id", "promo"."id"
	FROM "product" as "p"
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	LEFT JOIN (
		SELECT "promotion"."id", "promotion"."actived_date", "promotion"."expired_date", "promotion"."product_id"
		FROM "promotion" WHERE now() BETWEEN "promotion"."actived_date" AND "promotion"."expired_date"
	) as "promo" ON "promo"."product_id" = "p"."id"
	WHERE "s"."id" = $1 AND "p"."id" = $2;
	`
	CreatePromotionSellerQuery = `
	INSERT INTO "promotion"
		(name, product_id, discount_percentage, discount_fix_price, min_product_price, max_discount_price,
		quota, max_quantity, actived_date, expired_date)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	GetPromotionSellerDetailByIDQuery = `
	SELECT "promo"."id", "promo"."name", "p"."id", "p"."title", "p"."thumbnail_url", "promo"."discount_percentage",
		"promo"."discount_fix_price", "promo"."min_product_price", "promo"."max_discount_price", "promo"."quota", "promo"."max_quantity", 
		"promo"."actived_date", "promo"."expired_date", "promo"."created_at", "promo"."updated_at", "promo"."deleted_at"
	FROM "promotion" as "promo"
	INNER JOIN "product" as "p" ON "p"."id" = "promo"."product_id"
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	WHERE "promo"."id" = $1 AND "s"."id" = $2 AND "p"."id" = $3
	`

	UpdatePromotionSellerQuery = `
	UPDATE "promotion" SET "name" = $1, "max_quantity" = $2, "discount_percentage" = $3,
		"discount_fix_price" = $4, "min_product_price" = $5, "max_discount_price" = $6,
		"actived_date" = $7, "expired_date" = $8, "updated_at" = now()
	WHERE "id" = $9
	`

	GetDetailPromotionSellerByIDQuery = `
	SELECT "promo"."id", "promo"."name", "p"."id", "p"."title", "p"."min_price", "p"."max_price" ,"p"."thumbnail_url", "promo"."discount_percentage",
		"promo"."discount_fix_price", "promo"."min_product_price", "promo"."max_discount_price", "promo"."quota", "promo"."max_quantity", 
		"promo"."actived_date", "promo"."expired_date", "promo"."created_at", "promo"."updated_at", "promo"."deleted_at"
	FROM "promotion" as "promo"
	INNER JOIN "product" as "p" ON "p"."id" = "promo"."product_id"
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	WHERE "promo"."id" = $1 AND "s"."id" = $2
	`

	UpdateOrderByID               = `UPDATE "order" SET "order_status_id" = $1 WHERE "id" = $2`
	UpdateTransactionByID         = `UPDATE "transaction" SET "paid_at" = $1, "canceled_at" = $2 WHERE "id" = $3`
	GetProductDetailByIDQuery     = `SELECT "id", "price", "stock", "weight", "size", "hazardous", "condition", "bulk_price" FROM "product_detail" WHERE "id" = $1 AND "deleted_at" IS NULL;`
	GetTransactionsExpiredQuery   = `SELECT "id", "voucher_marketplace_id", "wallet_id", "card_number", "invoice", "total_price", "paid_at", "canceled_at", "expired_at" FROM "transaction" WHERE "paid_at" IS NULL AND "canceled_at" IS NULL AND "expired_at" < current_timestamp`
	GetOrderItemsByOrderIDQuery   = `SELECT "id", "order_id", "product_detail_id", "quantity", "item_price", "total_price" FROM "order_item" WHERE "order_id" = $1`
	UpdateProductDetailStockQuery = `UPDATE "product_detail" SET "stock" = $1, "updated_at" = now() WHERE "id" = $2;`
	GetOrderByTransactionID       = `SELECT 
		"id", "transaction_id", "shop_id", "user_id", "courier_id", "voucher_shop_id", "order_status_id", "total_price", "delivery_fee", "resi_no", "created_at", "arrived_at" 
	FROM "order" WHERE "transaction_id" = $1`

	GetTotalProductWithoutPromotionQuery = `SELECT count("p"."id") FROM "product" as "p"
		INNER JOIN "category" as "c" ON "c"."id" = "p"."category_id"
		LEFT JOIN "promotion" as "promo" ON "promo"."product_id" = "p"."id"
		WHERE "p"."shop_id" = $1 AND ("promo"."id" is NULL OR
	("promo"."actived_date" < now() AND "promo"."expired_date" < now()));
	`
	GetProductWithoutPromotionQuery = `
	SELECT "p"."id", "p"."title", "p"."min_price", "c"."name", "p"."thumbnail_url", "p"."unit_sold", "p"."rating_avg" FROM "product" as "p"
		INNER JOIN "category" as "c" ON "c"."id" = "p"."category_id"
		LEFT JOIN "promotion" as "promo" ON "promo"."product_id" = "p"."id"
		WHERE "p"."shop_id" = $1 AND ("promo"."id" is NULL OR
		("promo"."actived_date" < now() AND "promo"."expired_date" < now()))
		LIMIT $2 OFFSET $3;
		`
)
