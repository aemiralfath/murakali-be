package repository

const (
	GetSellerPerformanceMetadataQuery = `
	SELECT id AS shop_id, name AS shop_name, created_at AS shop_created_at, NOW() as report_updated_at
	FROM shop
	WHERE shop.id = $1 AND shop.deleted_at IS NULL
	`

	GetDailySalesQuery = `
	WITH calendar AS (
		SELECT 
			generate_series(date(NOW()) - INTERVAL '30 days', date(NOW()), '1 day')::date AS date
	)
	, daily_sales AS (
		SELECT
			shop_id,
			order_status_id,
			DATE(created_at) AS date,
			SUM(total_price) AS total_sales
		FROM "order"
		WHERE
			created_at >= (NOW() - INTERVAL '30 days') AND
			order_status_id = 7 AND
			shop_id = $1
		GROUP BY shop_id, date, order_status_id
	)
	
	SELECT 
		calendar.date,
		COALESCE(daily_sales.total_sales, 0) AS total_sales
	FROM calendar
	LEFT JOIN daily_sales ON calendar.date = daily_sales.date 
	ORDER BY calendar.date
	`

	GetDailyOrderQuery = `
	WITH calendar AS (
		SELECT
			generate_series(date(NOW()) - INTERVAL '30 days', date(NOW()), '1 day')::date AS date
	)
	, daily_order AS (
		SELECT
			DATE(created_at) AS date,
			SUM(CASE WHEN order_status_id = 7 THEN 1 ELSE 0 END) AS success_order,
			SUM(CASE WHEN order_status_id IN (8,9) THEN 1 ELSE 0 END) AS failed_order
		FROM "order"
		WHERE
			created_at >= (NOW() - INTERVAL '30 days') AND
			shop_id = $1
		GROUP BY date
	)
	
	SELECT 
		calendar.date, 
		COALESCE(daily_order.success_order, 0) as success_order,
		COALESCE(daily_order.failed_order, 0) as failed_order
		FROM calendar
		LEFT JOIN daily_order ON calendar.date = daily_order.date
	`

	GetMonthlyOrderQuery = `
	WITH month_order AS (
		SELECT 
			date_trunc('month', created_at) AS month,
			SUM(CASE WHEN order_status_id = 7 THEN 1 ELSE 0 END) AS success_order,
			SUM(CASE WHEN order_status_id IN (8,9) THEN 1 ELSE 0 END) AS failed_order
		FROM "order"
		WHERE 
			created_at >= (NOW() - INTERVAL '2 month') AND
			shop_id = $1
		GROUP BY month
	)
	
	SELECT 
		month,
		success_order,
		failed_order,
		NULL AS success_order_percent_change,
		NULL AS failed_order_percent_change
	FROM month_order
	ORDER BY month
	`

	GetTotalRatingQuery = `
	SELECT 
		SUM(CASE WHEN rating BETWEEN 1 AND 1.99 THEN 1 ELSE 0 END) as rating_1,
		SUM(CASE WHEN rating BETWEEN 2 AND 2.99 THEN 1 ELSE 0 END) as rating_2,
		SUM(CASE WHEN rating BETWEEN 3 AND 3.99 THEN 1 ELSE 0 END) as rating_3,
		SUM(CASE WHEN rating BETWEEN 4 AND 4.99 THEN 1 ELSE 0 END) as rating_4,
		SUM(CASE WHEN rating BETWEEN 5 AND 5.99 THEN 1 ELSE 0 END) as rating_5
	FROM review
	INNER JOIN product ON product.id = review.product_id
	WHERE shop_id = $1 AND review.deleted_at IS NULL
	GROUP BY shop_id
	`

	GetMostOrderedProductQuery = `
	SELECT 
    id, title, view_count, unit_sold, thumbnail_url
	FROM product
	WHERE shop_id = $1 AND deleted_at IS NULL
	ORDER BY unit_sold DESC
	LIMIT 5 
	`

	GetNumOrderByProvince = `
	WITH provinces AS (
		SELECT
		generate_series(1, 34) AS province_id)
			SELECT
				provinces.province_id,
				COALESCE(num_orders, 0) as num_orders
			FROM provinces
			LEFT JOIN (
				SELECT
					(json_extract_path_text(buyer_address::json, 'province_id')::integer) as province_id,
					COUNT(1) as num_orders
				FROM "order"
				WHERE shop_id = $1
				GROUP BY province_id
			) as orders ON provinces.province_id = orders.province_id
			ORDER BY num_orders DESC;
	`

	GetTotalSalesQuery = `
	SELECT 
		SUM(total_price) as total_sales,
		SUM(CASE WHEN is_withdraw = true THEN total_price ELSE 0 END) as withdrawn_sum,
		SUM(CASE WHEN is_withdraw = false THEN total_price ELSE 0 END) as withdrawable_sum
	FROM "order"
	WHERE shop_id = $1
	GROUP BY shop_id;
	`

	GetTotalAllSellerQuery = `SELECT count(s.id)
	FROM "shop" s 
	JOIN "user" u ON u.id = s.user_id
	WHERE s.name ILIKE $1 AND s.deleted_at is null`

	GetAllSellerQuery = `SELECT s.id, s.name,
	 s.total_product,
	 s.total_rating, 
	 s.rating_avg,
	  s.created_at,
	   u.photo_url 
	FROM "shop" s 
	JOIN "user" u ON u.id = s.user_id
	WHERE s.name ILIKE $1 AND s.deleted_at is null
	ORDER BY s.created_at desc
	LIMIT $2 OFFSET $3
	`

	GetUserByIDQuery   = `SELECT "id", "role_id", "email", "username", "phone_no", "fullname", "gender", "birth_date", "is_verify","photo_url" FROM "user" WHERE "id" = $1`
	GetShopByIDQuery   = `SELECT "id", "name", "user_id" FROM "shop" WHERE "id" = $1 AND "deleted_at" IS NULL;`
	GetTotalOrderQuery = `SELECT count(id) FROM "order" as "o" WHERE "o"."shop_id" = $1 and "o"."order_status_id"::text LIKE $2`

	GetTotalOrderWithVoucherIDQuery = `SELECT count(id) FROM "order" as "o" WHERE "o"."shop_id" = $1 and "o"."order_status_id"::text LIKE $2 
	AND "o"."voucher_shop_id" = $3
	`

	GetOrdersQuery = `SELECT o.id, o.is_withdraw, o.is_refund, o.order_status_id,o.total_price,o.delivery_fee,o.resi_no,s.id,s.name,v.code,o.created_at
	from "order" o
	join "shop" s on s.id = o.shop_id
	left join "voucher" v on v.id = o.voucher_shop_id 
	WHERE o.shop_id = $1 
	and "order_status_id"::text LIKE $2 
	`

	GetOrdersWithVoucherIDQuery = `SELECT o.id, o.is_withdraw, o.is_refund, o.order_status_id,o.total_price,o.delivery_fee,o.resi_no,s.id,s.name,v.code,o.created_at
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

	GetOrderByOrderID = `SELECT o.id, o.transaction_id, o.order_status_id, o.is_withdraw,o.is_refund,o.total_price,o.delivery_fee,o.resi_no,s.id,s.name,u2.phone_no,u2.username,v.code,o.created_at,t.invoice
	,c.name,c.code,c.service,c.description,u.username,u.phone_no
	from "order" o
	join "shop" s on s.id = o.shop_id
	join "courier" c on o.courier_id = c.id
	join "user" u on o.user_id = u.id
	join "user" u2 on s.user_Id = u2.id
	join transaction t on o.transaction_id = t.id
	left join "voucher" v on v.id = o.voucher_shop_id WHERE o.id = $1 ORDER BY o.created_at asc`

	GetBuyerIDByOrderIDQuery  = `SELECT o.user_id from "order" o where o.id = $1`
	CreateWalletHistoryQuery  = `INSERT INTO "wallet_history" (transaction_id, wallet_id, "from", "to", description, amount, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	GetSellerIDByOrderIDQuery = `SELECT s.user_id from "order" o join shop s on o.shop_id = s.id where o.id = $1`

	GetOrderDetailQuery = `SELECT pd.id,pd.product_id,p.title, pd.weight,
        p.thumbnail_url,oi.quantity,oi.item_price,oi.total_price
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
	CancelOrderStatusQuery = `UPDATE "order" SET "order_status_id" = $1, "cancel_notes" = $2, "is_refund" = $3 WHERE "id" = $4`

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

	UpdateShopInformationByUserIDQuery = `UPDATE "shop" SET "name" = $1 WHERE "user_id" = $2`

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

	CountCodeVoucher = `
	SELECT count(code) FROM "voucher" as "v" WHERE "v"."code" = $1  AND "v"."deleted_at" IS NULL
	`

	GetTotalVoucherSellerQuery = `
	SELECT count(id) FROM "voucher" as "v" WHERE "v"."shop_id" = $1 AND "v"."deleted_at" IS NULL
	`
	GetAllVoucherSellerQuery = `
	SELECT "v"."id", "v"."shop_id", "v"."code", "v"."quota", "v"."actived_date", "v"."expired_date",
		"v"."discount_percentage", "v"."discount_fix_price", "v"."min_product_price", "v"."max_discount_price",
		"v"."created_at", "v"."updated_at",  "v"."deleted_at"
	FROM "voucher" as "v"
	INNER JOIN "shop" as "s" ON "s"."id" = "v"."shop_id"
	WHERE "v"."shop_id" = $1
	AND "v"."deleted_at" IS NULL
	`
	OrderBySomething = ` 
	ORDER BY %s LIMIT %d OFFSET %d`

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

	CreateRefundSellerQuery = `INSERT INTO "refund" (order_id, is_seller_refund, reason, accepted_at) VALUES($1, $2, $3, $4)`

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
	OrderByCreatedAt = ` ORDER BY "promo"."created_at" DESC`

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

	UpdateOrderByID               = `UPDATE "order" SET "order_status_id" = $1, "is_withdraw" = $2 WHERE "id" = $3`
	UpdateTransactionByID         = `UPDATE "transaction" SET "paid_at" = $1, "canceled_at" = $2 WHERE "id" = $3`
	GetProductDetailByIDQuery     = `SELECT "id", "price", "stock", "weight", "size", "hazardous", "condition", "bulk_price" FROM "product_detail" WHERE "id" = $1 AND "deleted_at" IS NULL;`
	GetTransactionsExpiredQuery   = `SELECT "id", "voucher_marketplace_id", "wallet_id", "card_number", "invoice", "total_price", "paid_at", "canceled_at", "expired_at" FROM "transaction" WHERE "paid_at" IS NULL AND "canceled_at" IS NULL AND "expired_at" < current_timestamp`
	GetOrderItemsByOrderIDQuery   = `SELECT "id", "order_id", "product_detail_id", "quantity", "item_price", "total_price" FROM "order_item" WHERE "order_id" = $1`
	UpdateProductDetailStockQuery = `UPDATE "product_detail" SET "stock" = $1, "updated_at" = now() WHERE "id" = $2;`
	GetOrderByTransactionID       = `SELECT 
		"id", "transaction_id", "shop_id", "user_id", "courier_id", "voucher_shop_id", "order_status_id", "total_price", "delivery_fee", "resi_no", "created_at", "arrived_at" 
	FROM "order" WHERE "transaction_id" = $1`
	UpdateWalletBalanceQuery = `UPDATE "wallet" SET "balance" = $1, "updated_at" = $2 WHERE "id" = $3`

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

	GetWalletByUserIDQuery = `SELECT "id", "user_id", "balance", "pin", "attempt_count", "attempt_at", "unlocked_at", "active_date" FROM "wallet" WHERE "user_id" = $1 AND "deleted_at" IS NULL`

	GetOrderModelByIDQuery = `SELECT "id", "transaction_id", "shop_id", "user_id", "courier_id", "voucher_shop_id", "order_status_id", "total_price",
	"delivery_fee", "resi_no", "buyer_address", "shop_address", "cancel_notes", "is_withdraw", "is_refund", "created_at", "arrived_at"
	FROM "order" WHERE "id" = $1`

	GetRefundOrderByOrderIDQuery = `SELECT "id", "order_id", "is_seller_refund", "is_buyer_refund", "reason", "image", "accepted_at", "rejected_at", "refunded_at"
	FROM "refund" WHERE "order_id" = $1 ORDER BY "rejected_at" DESC LIMIT 1`

	GetRefundOrderByIDQuery = `SELECT "id", "order_id", "is_seller_refund", "is_buyer_refund", "reason", "image", "accepted_at", "rejected_at", "refunded_at"
	FROM "refund" WHERE "id" = $1 ORDER BY "rejected_at" DESC LIMIT 1`

	GetRefundThreadByRefundIDQuery = `SELECT "rt"."id", "rt"."refund_id", "rt"."user_id", "u"."username", "s"."name", "u"."photo_url", "rt"."is_seller", "rt"."is_buyer", "rt"."text", "rt"."created_at"
	FROM "refund_thread" as "rt"
    LEFT JOIN "user" as "u" ON "u"."id" = "rt"."user_id"
    LEFT JOIN "shop" as "s" ON "s"."user_id" = "u"."id"
	WHERE "refund_id" = $1 ORDER BY "created_at" ASC`

	CreateRefundThreadSellerQuery = `INSERT INTO "refund_thread" 
	(refund_id, user_id, is_seller, is_buyer, text)
	VALUES ($1, $2, $3, $4, $5)`

	UpdateRefundAcceptQuery = `UPDATE "refund" SET "accepted_at" = now() WHERE "id" = $1;`
	UpdateRefundRejectQuery = `UPDATE "refund" SET "rejected_at" = now() WHERE "id" = $1;`

	UpdateOrderRefundRejectedQuery = `UPDATE "order" SET "is_refund" = FALSE WHERE "id" = $1`
)
