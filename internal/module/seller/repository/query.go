package repository

const (
	GetTotalOrderQuery = `SELECT count(id) FROM "order" WHERE "shop_id" = $1 and "order_status_id"::text LIKE $2`

	GetOrdersQuery = `SELECT o.id,o.order_status_id,o.total_price,o.delivery_fee,o.resi_no,s.id,s.name,v.code,o.created_at
	from "order" o
	join "shop" s on s.id = o.shop_id
	left join "voucher" v on v.id = o.voucher_shop_id WHERE o.shop_id = $1 and "order_status_id"::text LIKE $2 ORDER BY o.created_at asc LIMIT $3 OFFSET $4
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

	GetOrderDetailQuery = `SELECT pd.id,pd.product_id,p.title,
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
	"order" set resi_no = $1 WHERE id = $2 
	AND shop_id = $3`

	GetTotalVoucherSellerQuery = `
	SELECT count(id) FROM "voucher" WHERE "shop_id" = $1
	`
	GetAllVoucherSellerQuery = `
	SELECT "v"."id", "v"."shop_id", "v"."code", "v"."quota", "v"."actived_date", "v"."expired_date",
		"v"."discount_percentage", "v"."discount_fix_price", "v"."min_product_price", "v"."max_discount_price",
		"v"."created_at", "v"."updated_at",  "v"."deleted_at"
	FROM "voucher" as "v"
	INNER JOIN "shop" as "s" ON "s"."id" = "v"."shop_id"
	WHERE "v"."shop_id" = $1
	`
	CreateVoucherSellerQuery = `INSERT INTO "voucher" 
    	(shop_id, code, quota, actived_date, expired_date, discount_percentage, discount_fix_price, min_product_price, max_discount_price)
    	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
)
