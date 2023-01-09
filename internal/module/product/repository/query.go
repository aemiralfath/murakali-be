package repository

const (
	GetCategoriesQuery           = `SELECT "id", "parent_id", "name", "photo_url" FROM "category" WHERE "parent_id" IS NULL AND "deleted_at" IS NULL`
	GetCategoriesByNameQuery     = `SELECT "id", "parent_id", "name", "photo_url" FROM "category" WHERE "name" = $1 AND "deleted_at" IS NULL`
	GetCategoriesByParentIdQuery = `SELECT "id", "parent_id", "name", "photo_url" FROM "category" WHERE "parent_id" = $1 AND "deleted_at" IS NULL`
	GetBannersQuery              = `SELECT "id", "title", "content", "image_url", "page_url", "is_active" FROM "banner" WHERE "is_active" = TRUE`
	GetTotalProductQuery         = `SELECT count(id) FROM "product" WHERE "deleted_at" IS NULL`
	GetRecommendedProductsQuery  = `
	SELECT "p"."title" as "title", "p"."unit_sold" as "unit_sold", "p"."rating_avg" as "rating_avg", "p"."thumbnail_url" as "thumbnail_url",
		"p"."min_price" as "min_price", "p"."max_price" as "max_price", "promo"."discount_percentage" as "promo_discount_percentage",  "promo"."discount_fix_price" as "promo_discount_fix_price",
		"promo"."min_product_price" as "promo_min_product_price",  "promo"."max_discount_price" as "promo_max_discount_price",
		"v"."discount_percentage" as "voucher_discount_percentage",  "v"."discount_fix_price" as "voucher_discount_fix_price", "s"."name" as "shop_name", "c"."name" as "category_name"
	FROM "product" as "p"
	LEFT JOIN (
		SELECT * FROM "promotion"
		WHERE now() BETWEEN "promotion"."actived_date" AND "promotion"."expired_date"
	) as "promo" ON "promo"."product_id" = "p"."id"
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	LEFT JOIN (
		SELECT * FROM "voucher"
		WHERE now() BETWEEN "voucher"."actived_date" AND "voucher"."expired_date"
	) as "v" ON "v"."shop_id" = "s"."id"
	INNER JOIN "category" as "c" ON "c"."id" = "p"."category_id"
	WHERE "p"."deleted_at" IS NULL
	ORDER BY "p"."unit_sold" DESC LIMIT $1 OFFSET $2;
	`
	GetProductInfoQuery = `select
	pr.id,pr.sku,pr.title,pr.description,pr.view_count,pr.favorite_count,pr.unit_sold,pr.listed_status,pr.thumbnail_url,pr.rating_avg,pr.min_price,pr.max_price
	,c.name,c.photo_url
	from 
	product pr 
	join product_detail b on pr.id = b.product_id 
	join category c on c.id = pr.category_id
	where pr.id = $1`

	GetProductDetailQuery = `select
	pd.id,pd.price,pd.stock,pd.weight,pd.size,pd.hazardous,pd.condition,pd.bulk_price,g.url
	from 
	product_detail pd
	join photo g on pd.id = g.product_detail_id
	where pd.product_id = $1`

	GetVariantDetailQuery = `select b.type,b.name from variant a join variant_detail b on a.variant_detail_id = b.id
	where a.product_detail_id = $1`

	GetPromotionDetailQuery = `select
	pro.name,pro.discount_percentage,pro.discount_fix_price,pro.min_product_price,pro.max_discount_price,pro.quota,pro.max_quantity,pro.actived_date,pro.expired_date
	from 
	promotion pro
	where pro.product_id = $1`

	GetProductsQuery = `
	SELECT "p"."id" as "product_id","p"."title" as "title", "p"."unit_sold" as "unit_sold", "p"."rating_avg" as "rating_avg", "p"."thumbnail_url" as "thumbnail_url",
		"p"."min_price" as "min_price", "p"."max_price" as "max_price", "p"."view_count" as "view_count", 
		"promo"."discount_percentage" as "promo_discount_percentage",  "promo"."discount_fix_price" as "promo_discount_fix_price",
		"promo"."min_product_price" as "promo_min_product_price",  "promo"."max_discount_price" as "promo_max_discount_price",
		"v"."discount_percentage" as "voucher_discount_percentage",  "v"."discount_fix_price" as "voucher_discount_fix_price", 
		"s"."name" as "shop_name", 
		"c"."name" as "category_name",
		"a"."province" as "province"
	FROM "product" as "p"
	LEFT JOIN (
		SELECT * FROM "promotion"
		WHERE now() BETWEEN "promotion"."actived_date" AND "promotion"."expired_date"
	) as "promo" ON "promo"."product_id" = "p"."id"
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	LEFT JOIN (
		SELECT * FROM "voucher"
		WHERE now() BETWEEN "voucher"."actived_date" AND "voucher"."expired_date"
	) as "v" ON "v"."shop_id" = "s"."id"
	INNER JOIN "category" as "c" ON "c"."id" = "p"."category_id"
	INNER JOIN "user" as "u" ON "u"."id" = "s"."user_id"
	INNER JOIN "address" as "a" ON "u"."id" = "a"."user_id"
	WHERE "p".title ILIKE $1 
	AND  "c".name ILIKE $2
	AND ("p".rating_avg BETWEEN $3 AND $4)
	AND ("p".min_price BETWEEN $5 AND $6)
	AND "p"."deleted_at" IS NULL 
	`

	GetProductsWithProvinceQuery = `
	SELECT "p"."id" as "product_id",
	"p"."title" as "title", "p"."unit_sold" as "unit_sold", "p"."rating_avg" as "rating_avg", "p"."thumbnail_url" as "thumbnail_url",
		"p"."min_price" as "min_price", "p"."max_price" as "max_price", "p"."view_count" as "view_count", 
		"promo"."discount_percentage" as "promo_discount_percentage",  "promo"."discount_fix_price" as "promo_discount_fix_price",
		"promo"."min_product_price" as "promo_min_product_price",  "promo"."max_discount_price" as "promo_max_discount_price",
		"v"."discount_percentage" as "voucher_discount_percentage",  "v"."discount_fix_price" as "voucher_discount_fix_price", 
		"s"."name" as "shop_name", 
		"c"."name" as "category_name",
		"a"."province" as "province"
	FROM "product" as "p"
	LEFT JOIN (
		SELECT * FROM "promotion"
		WHERE now() BETWEEN "promotion"."actived_date" AND "promotion"."expired_date"
	) as "promo" ON "promo"."product_id" = "p"."id"
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	LEFT JOIN (
		SELECT * FROM "voucher"
		WHERE now() BETWEEN "voucher"."actived_date" AND "voucher"."expired_date"
	) as "v" ON "v"."shop_id" = "s"."id"
	INNER JOIN "category" as "c" ON "c"."id" = "p"."category_id"
	INNER JOIN "user" as "u" ON "u"."id" = "s"."user_id"
	INNER JOIN "address" as "a" ON "u"."id" = "a"."user_id"
	WHERE "p".title ILIKE $1 
	AND  "c".name ILIKE $2
	AND ("p".rating_avg BETWEEN $3 AND $4)
	AND ("p".min_price BETWEEN $5 AND $6)
	AND "p"."deleted_at" IS NULL 
	AND ("a"."province_id"::text =any($7))
	`

	OrderBySomething = ` 
	ORDER BY %s LIMIT %d OFFSET %d`

	WhereShopIds = ` 
		AND "s"."id" = '%s'
	`

	GetAllTotalProductQuery = `
	SELECT count("p"."id") FROM "product" as "p" 
	INNER JOIN "category" as "c" ON "c"."id" = "p"."category_id" 
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	INNER JOIN "user" as "u" ON "u"."id" = "s"."user_id"
	INNER JOIN "address" as "a" ON "u"."id" = "a"."user_id"
	WHERE "p".title ILIKE $1 
	AND  "c".name ILIKE $2
	AND ("p".rating_avg BETWEEN $3 AND $4)
	AND ("p".min_price BETWEEN $5 AND $6)
	AND "p"."deleted_at" IS NULL`

	GetAllTotalProductWithProvinceQuery = `
	SELECT count("p"."id") FROM "product" as "p" 
	INNER JOIN "category" as "c" ON "c"."id" = "p"."category_id" 
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	INNER JOIN "user" as "u" ON "u"."id" = "s"."user_id"
	INNER JOIN "address" as "a" ON "u"."id" = "a"."user_id"
	WHERE "p".title ILIKE $1 
	AND  "c".name ILIKE $2
	AND ("p".rating_avg BETWEEN $3 AND $4)
	AND ("p".min_price BETWEEN $5 AND $6)
	AND "p"."deleted_at" IS NULL
	AND ("a"."province_id"::text =any($7))`

	GetFavoriteProductsQuery = `
	SELECT "p"."id" as "product_id","p"."title" as "title", "p"."unit_sold" as "unit_sold", "p"."rating_avg" as "rating_avg", "p"."thumbnail_url" as "thumbnail_url",
		"p"."min_price" as "min_price", "p"."max_price" as "max_price", "p"."view_count" as "view_count", "promo"."discount_percentage" as "promo_discount_percentage",  "promo"."discount_fix_price" as "promo_discount_fix_price",
		"promo"."min_product_price" as "promo_min_product_price",  "promo"."max_discount_price" as "promo_max_discount_price",
		"v"."discount_percentage" as "voucher_discount_percentage",  "v"."discount_fix_price" as "voucher_discount_fix_price", "s"."name" as "shop_name", "c"."name" as "category_name"
	FROM  "favorite" as "f"
	INNER JOIN "product" as "p" ON "p"."id" = "f"."product_id"
	LEFT JOIN (
		SELECT * FROM "promotion"
		WHERE now() BETWEEN "promotion"."actived_date" AND "promotion"."expired_date"
	) as "promo" ON "promo"."product_id" = "p"."id"
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	LEFT JOIN (
		SELECT * FROM "voucher"
		WHERE now() BETWEEN "voucher"."actived_date" AND "voucher"."expired_date"
	) as "v" ON "v"."shop_id" = "s"."id"
	INNER JOIN "category" as "c" ON "c"."id" = "p"."category_id"
	WHERE 
	"p".title ILIKE $1 
	AND  "c".name ILIKE $2
	AND ("p".rating_avg BETWEEN $3 AND $4)
	AND ("p".min_price BETWEEN $5 AND $6)
	AND "f"."user_id" = $7
	AND "p"."deleted_at" IS NULL
	ORDER BY %s LIMIT $8 OFFSET $9;
	`

	GetAllTotalFavoriteProductQuery = `
	SELECT count("p"."id") 
	FROM  "favorite" as "f"
	INNER JOIN "product" as "p" ON "p"."id" = "f"."product_id"
	INNER JOIN "category" as "c" ON "c"."id" = "p"."category_id" 
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	WHERE "p".title ILIKE $1 
	AND  "c".name ILIKE $2
	AND ("p".rating_avg BETWEEN $3 AND $4)
	AND ("p".min_price BETWEEN $5 AND $6)
	AND "f"."user_id" = $7
	 AND "p"."deleted_at" IS NULL `

	GetShopIDByUserIDQuery = `SELECT id from "shop" WHERE user_id = $1 AND deleted_at IS NULL `

	CreateProductQuery = `INSERT INTO "product" 
	(category_id, shop_id, sku, title,
	 description, view_count, favorite_count, 
	 unit_sold, listed_status, thumbnail_url,
	  rating_avg, min_price, max_price)
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING "id";`

	CreateProductDetailQuery = `INSERT INTO "product_detail" 
	(product_id, price, stock, weight, 
		size, hazardous, condition, bulk_price)
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING "id";`

	CreatePhotoQuery = `INSERT INTO "photo" 
	(product_detail_id, url)
	 VALUES ($1, $2) RETURNING "id";`

	CreateVideoQuery = `INSERT INTO "video" 
	(product_detail_id, url)
	 VALUES ($1, $2) RETURNING "id";`

	CreateVariantQuery = `INSERT INTO "variant" 
	(product_detail_id, variant_detail_id)
	 VALUES ($1, $2) RETURNING "id";`

	CreateProductCourierQuery = `INSERT 
	INTO "product_courier_whitelist" 
	(product_id, courier_id)
	 VALUES ($1, $2) RETURNING "id";`

	GetListedStatusQuery = `SELECT listed_status from "product" WHERE id = $1 AND deleted_at IS NULL `

	UpdateListedStatusQuery = `UPDATE "product" SET "listed_status" = $1 WHERE "id" = $2`
)
