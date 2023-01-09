package repository

const (
	GetUserByIDQuery      = `SELECT "id", "role_id", "email", "username", "phone_no", "fullname", "gender", "birth_date", "is_verify","photo_url" FROM "user" WHERE "id" = $1`
	GetTotalCartQuery     = `SELECT count(id) FROM "cart_item" WHERE "user_id" = $1 AND "deleted_at" IS NULL`
	GetCartHoverHomeQuery = `
	SELECT "p"."title" as "title", "p"."thumbnail_url" as "thumbnail_url", "pd"."price" as "price", "promo"."discount_percentage" as "discount_percentage",
	"promo"."discount_fix_price" as "discount_fix_price", "promo"."min_product_price" as "min_product_price", "promo"."max_discount_price" as "max_discount_price",
	"cart"."quantity" as "quantity", array_agg("vd"."name") as "variant_name", array_agg("vd"."type") as "variant_type"
	FROM
	"cart_item" as "cart"
	INNER JOIN "product_detail" as "pd" ON "pd"."id" = "cart"."product_detail_id"
	INNER JOIN "product" as "p" ON "p"."id" = "pd"."product_id"
	LEFT JOIN (
		SELECT * FROM "promotion"
		WHERE now() BETWEEN "promotion"."actived_date" AND "promotion"."expired_date" 
	) as "promo" ON "promo"."product_id" = "p"."id"
	LEFT JOIN "variant" as "v" ON "v"."product_detail_id" = "pd"."id"
	LEFT JOIN "variant_detail" as "vd" ON "vd"."id" = "v"."variant_detail_id"
	WHERE 
	"cart"."user_id" = $1 AND "cart"."deleted_at" IS NULL
	GROUP BY "cart"."id", "pd"."id", "p"."id", "promo"."discount_percentage", "promo"."discount_fix_price", "promo"."min_product_price",
		"promo"."max_discount_price"
	ORDER BY "cart"."created_at" DESC
	LIMIT $2;
	`
	GetCartItemsQuery = `
	SELECT "ci"."id" as "id", "ci"."quantity" as "quantity", "pd"."id" as "product_detail_id", "p"."title" as "product_title", "s"."id" as "shop_id", "s"."name" as "shop_name", "p"."thumbnail_url" as "thumbnail_url", 
		"pd"."price" as "product_price", "pd"."stock" as "product_stock", "pd"."weight" as "product_weight",
		"promo"."discount_percentage" as "promo_discount_percentage", "promo"."discount_fix_price" as "promo_discount_fix_price",
		"promo"."min_product_price" as "promo_min_product_price", "promo"."max_discount_price" as "promo_max_discount_price", 
		array_agg("vd"."name") as "variant_name", array_agg("vd"."type") as "variant_type"
	FROM "cart_item" as "ci"
	INNER JOIN "product_detail" as "pd" ON "pd"."id" = "ci"."product_detail_id"
	LEFT JOIN "variant" as "v" ON "v"."product_detail_id" = "pd"."id"
	LEFT JOIN "variant_detail" as "vd" ON "vd"."id" = "v"."variant_detail_id"
	INNER JOIN "product" as "p" ON "p"."id" = "pd"."product_id"
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	LEFT JOIN (
		SELECT * FROM "promotion"
		WHERE now() BETWEEN "promotion"."actived_date" AND "promotion"."expired_date" 
	) as "promo" ON "promo"."product_id" = "p"."id"
	WHERE "ci"."user_id" = $1 AND "ci"."deleted_at" IS NULL
	GROUP BY "ci"."id", "pd"."id", "p"."id", "s"."id", "promo"."discount_percentage", "promo"."discount_fix_price", "promo"."min_product_price",
		"promo"."max_discount_price"
	ORDER BY "ci"."created_at" DESC LIMIT $2 OFFSET $3;
	`
	GetProductDetailByIDQuery = `
	SELECT "id", "product_id", "price", "stock", "weight", "size"
	FROM "product_detail"
	WHERE id = $1 AND "deleted_at" IS NULL;
	`
	GetCartProductDetailQuery = `
	SELECT "id", "user_id", "product_detail_id", "quantity"
	FROM "cart_item"
	WHERE "user_id" = $1 AND "product_detail_id" = $2 AND "deleted_at" IS NULL;
	`
	CreateCartQuery     = `INSERT INTO "cart_item" (user_id, product_detail_id, quantity) VALUES ($1, $2, $3) RETURNING "id";`
	UpdateCartByIDQuery = `UPDATE "cart_item" SET "quantity" = $1, "updated_at" = $2 WHERE "id" = $3;`
	DeleteCartByIDQuery = `DELETE FROM "cart_item" WHERE "id" = $1`
)
