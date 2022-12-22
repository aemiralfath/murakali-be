package repository

const (
	GetTotalCartQuery     = `SELECT count(id) FROM "cart_item" WHERE "user_id" = $1 AND "deleted_at" IS NULL`
	GetCartHoverHomeQuery = `
	SELECT "p"."title" as "title", "p"."thumbnail_url" as "thumbnail_url", "pd"."price" as "price", "promo"."discount_percentage" as "discount_percentage",
	"promo"."discount_fix_price" as "discount_fix_price", "promo"."min_product_price" as "min_product_price", "promo"."max_discount_price" as "max_discount_price",
	"cart"."quantity" as "quantity", "vd"."name" as "variant_name", "vd"."type" as "variant_type"
	FROM
	"cart_item" as "cart" 
	INNER JOIN "user" as "u" ON "u"."id" = "cart"."user_id"
	INNER JOIN "product_detail" as "pd" ON "pd"."id" = "cart"."product_detail_id"
	INNER JOIN "product" as "p" ON "p"."id" = "pd"."product_id"
	INNER JOIN "promotion" as "promo" ON "promo"."product_id" = "p"."id"
	INNER JOIN "variant" as "v" ON "v"."product_detail_id" = "pd"."id"
	INNER JOIN "variant_detail" as "vd" ON "vd"."id" = "v"."variant_detail_id"
	WHERE 
	"cart"."user_id" = $1 AND "cart"."deleted_at" IS NULL
	ORDER BY "cart"."updated_at"
	LIMIT $2
	`
)
