package repository

const (
	GetTotalOrderQuery = `SELECT count(id) FROM "order" WHERE "shop_id" = $1 and "order_status_id"::text LIKE $2`

	GetOrdersQuery = `SELECT o.id,o.order_status_id,o.total_price,o.delivery_fee,o.resi_no,s.id,s.name,v.code,o.created_at
	from "order" o
	join "shop" s on s.id = o.shop_id
	left join "voucher" v on v.id = o.voucher_shop_id WHERE o.shop_id = $1 and "order_status_id"::text LIKE $2 ORDER BY o.created_at asc LIMIT $3 OFFSET $4
	`

	GetOrderByOrderID = `SELECT o.id,o.order_status_id,o.total_price,o.delivery_fee,o.resi_no,s.id,s.name,v.code,o.created_at
	from "order" o
	join "shop" s on s.id = o.shop_id
	join "voucher" v on v.id = o.voucher_shop_id WHERE o.id = $1 ORDER BY o.created_at asc`

	GetOrderDetailQuery = `SELECT pd.id,pd.product_id,p.title,ph.url,oi.quantity,oi.item_price,oi.total_price
	from  "product_detail" pd 
	join "photo" ph on pd.id = ph.product_detail_id join "order_item" oi on pd.id = oi.product_detail_id 
	join "product" p on p.id = pd.product_id WHERE oi.order_id = $1`

	GetShopIDByUserQuery = `SELECT id from shop where user_id = $1 and deleted_at is null`

	GetShopIDByOrderQuery = `SELECT shop_id from "order" where id = $1 `

	ChangeOrderStatusQuery = `UPDATE "order" SET "order_status_id" = $1 WHERE "id" = $2`



	GetCourierSellerQuery  = `
	SELECT "sp"."id" as "shop_courier_id", "c"."name" as "name", "c"."code" as "code", "c"."service" as "service",
		"c"."description" as "description"
	FROM "shop_courier" as "sp"
	INNER JOIN "courier" as "c" ON "c"."id" = "sp"."courier_id"
	INNER JOIN "shop" as "s" ON "s"."id" = "sp"."shop_id"
	WHERE "s"."user_id" = $1 AND "sp"."deleted_at" IS NULL
	ORDER BY "sp"."created_at" DESC ;
	`


	GetShopIDByShopIDQuery = `SELECT s.id, s.user_id, s.name, s.total_product,
	 s.total_rating, s.rating_avg, s.created_at, u.photo_url 
	FROM "shop" s 
	JOIN "user" u ON u.id = s.user_id
	WHERE s.id = $1 AND s.deleted_at is null`

	
	GetShopIDByUserIDQuery = `SELECT id from "shop" WHERE user_id = $1 AND deleted_at IS NULL `
	GetCourierSellerIDByUserIDQuery = `SELECT id from "shop_courier" WHERE shop_id = $1 AND courier_id = $2 AND deleted_at IS NULL `
	CreateCourierSellerQuery = `INSERT INTO "shop_courier" 
    	(shop_id, courier_id)
    	VALUES ($1, $2)`
	DeleteCourierSellerQuery  = `UPDATE "shop_courier" set deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`
)
