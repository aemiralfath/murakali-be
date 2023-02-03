package repository

const (
	GetCategoriesQuery           = `SELECT "id", "parent_id", "name", "photo_url" FROM "category" WHERE "parent_id" IS NULL AND "deleted_at" IS NULL`
	GetCategoriesByNameQuery     = `SELECT "id", "parent_id", "name", "photo_url" FROM "category" WHERE "name" = $1 AND "deleted_at" IS NULL`
	GetCategoriesByParentIdQuery = `SELECT "id", "parent_id", "name", "photo_url" FROM "category" WHERE "parent_id" = $1 AND "deleted_at" IS NULL`
	GetBannersQuery              = `SELECT "id", "title", "content", "image_url", "page_url", "is_active" FROM "banner" WHERE "is_active" = TRUE`
	GetTotalProductQuery         = `SELECT count(id) FROM "product" 	WHERE listed_status = true `
	GetRecommendedProductsQuery  = `
	SELECT "p"."id" as "id", "p"."title" as "title", "p"."unit_sold" as "unit_sold", "p"."rating_avg" as "rating_avg", "p"."thumbnail_url" as "thumbnail_url",
		"p"."min_price" as "min_price", "p"."max_price" as "max_price", "promo"."discount_percentage" as "promo_discount_percentage",  "promo"."discount_fix_price" as "promo_discount_fix_price",
		"promo"."min_product_price" as "promo_min_product_price", "promo"."max_discount_price" as "promo_max_discount_price",
		"v"."discount_percentage" as "voucher_discount_percentage", "v"."discount_fix_price" as "voucher_discount_fix_price", "s"."name" as "shop_name", "c"."name" as "category_name"
	FROM "product" as "p"
	LEFT JOIN (
		SELECT * FROM "promotion"
		WHERE (now() BETWEEN "promotion"."actived_date" AND "promotion"."expired_date") AND "promotion"."quota" > 0
	) as "promo" ON "promo"."product_id" = "p"."id"
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	LEFT JOIN (
		SELECT * FROM "voucher"
		WHERE now() BETWEEN "voucher"."actived_date" AND "voucher"."expired_date" AND "voucher"."quota" > 0
	) as "v" ON "v"."shop_id" = "s"."id"
	INNER JOIN "category" as "c" ON "c"."id" = "p"."category_id"
	WHERE "p"."listed_status" = true
	ORDER BY "p"."unit_sold" DESC,
	"p"."rating_avg" DESC,
	"p"."view_count" DESC
	LIMIT $1 OFFSET $2;
	`
	GetProductInfoQuery = `select
	pr.id,pr.sku,pr.title,pr.description,pr.view_count,pr.favorite_count,pr.unit_sold,pr.listed_status,pr.thumbnail_url,pr.rating_avg,pr.min_price,pr.max_price,pr.shop_id
	,c.name,c.photo_url
	from 
	product pr 
	join product_detail b on pr.id = b.product_id 
	join category c on c.id = pr.category_id
	where pr.id = $1 and b.deleted_at is null`

	GetProductDetailQuery = `select
	pd.id,pd.price,pd.stock,pd.weight,pd.size,pd.hazardous,pd.condition,pd.bulk_price
	from 
	product_detail pd
	where pd.product_id = $1 and pd.deleted_at is null`

	GetProductDetailPhotosQuery = `select
	g.url
	from photo g 
	where product_detail_id = $1`

	GetVariantDetailQuery = `select b.type,b.name from variant a join variant_detail b on a.variant_detail_id = b.id
	where a.product_detail_id = $1`

	GetVariantInfoQuery = `select a.id,b.id,b.name from variant a join variant_detail b on a.variant_detail_id = b.id
	where a.product_detail_id = $1`

	GetPromotionDetailQuery = `
	SELECT "promo"."name", "promo"."discount_percentage", "promo"."discount_fix_price", "promo"."min_product_price", 
		"promo"."max_discount_price", "promo"."quota", "promo"."max_quantity", "promo"."actived_date", "promo"."expired_date"
	FROM "promotion" as "promo"
	WHERE "promo"."product_id" = $1 AND (now() BETWEEN "promo"."actived_date" AND "promo"."expired_date")`

	GetProductsQuery = `
	SELECT "p"."id" as "product_id","p"."title" as "title", "p"."unit_sold" as "unit_sold", "p"."rating_avg" as "rating_avg", "p"."thumbnail_url" as "thumbnail_url",
		"p"."min_price" as "min_price", "p"."max_price" as "max_price", "p"."view_count" as "view_count", 
		"promo"."discount_percentage" as "promo_discount_percentage",  "promo"."discount_fix_price" as "promo_discount_fix_price",
		"promo"."min_product_price" as "promo_min_product_price",  "promo"."max_discount_price" as "promo_max_discount_price",
		"v"."discount_percentage" as "voucher_discount_percentage",  "v"."discount_fix_price" as "voucher_discount_fix_price", 
		"s"."name" as "shop_name", 
		"c"."name" as "category_name",
		"a"."province" as "province",
		"p".listed_status,
		"p"."created_at",
		"p"."updated_at"
	FROM "product" as "p"
	LEFT JOIN (
		SELECT * FROM "promotion"
		WHERE (now() BETWEEN "promotion"."actived_date" AND "promotion"."expired_date") AND "promotion"."quota" > 0
	) as "promo" ON "promo"."product_id" = "p"."id"
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	LEFT JOIN (
		SELECT * FROM "voucher"
		WHERE (now() BETWEEN "voucher"."actived_date" AND "voucher"."expired_date") AND "voucher"."quota" > 0
	) as "v" ON "v"."shop_id" = "s"."id"
	INNER JOIN "category" as "c" ON "c"."id" = "p"."category_id"
	INNER JOIN "user" as "u" ON "u"."id" = "s"."user_id"
	INNER JOIN "address" as "a" ON "u"."id" = "a"."user_id"
	WHERE "p".title ILIKE $1 
	AND  "c".name ILIKE $2
	AND ("p".rating_avg BETWEEN $3 AND $4)
	AND ("p".min_price BETWEEN $5 AND $6)
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
		"a"."province" as "province",
		"p".listed_status,
		"p"."created_at",
		"p"."updated_at"
	FROM "product" as "p"
	LEFT JOIN (
		SELECT * FROM "promotion"
		WHERE (now() BETWEEN "promotion"."actived_date" AND "promotion"."expired_date") AND "promotion"."quota" > 0
	) as "promo" ON "promo"."product_id" = "p"."id"
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	LEFT JOIN (
		SELECT * FROM "voucher"
		WHERE (now() BETWEEN "voucher"."actived_date" AND "voucher"."expired_date") AND "voucher"."quota" > 0
	) as "v" ON "v"."shop_id" = "s"."id"
	INNER JOIN "category" as "c" ON "c"."id" = "p"."category_id"
	INNER JOIN "user" as "u" ON "u"."id" = "s"."user_id"
	INNER JOIN "address" as "a" ON "u"."id" = "a"."user_id"
	WHERE "p".title ILIKE $1 
	AND  "c".name ILIKE $2
	AND ("p".rating_avg BETWEEN $3 AND $4)
	AND ("p".min_price BETWEEN $5 AND $6)
	AND ("a"."province_id"::text =any($7))
	`

	OrderBySomething = ` 
	ORDER BY %s LIMIT %d OFFSET %d`

	WhereShopIds = ` 
		AND "s"."id" = '%s'
	`

	WhereListedStatusTrue = ` 
		AND "p"."listed_status" = true
	`
	WhereListedStatusFalse = ` 
		AND "p"."listed_status" = false
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
	`

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
		WHERE (now() BETWEEN "promotion"."actived_date" AND "promotion"."expired_date") AND "promotion"."quota" > 0
	) as "promo" ON "promo"."product_id" = "p"."id"
	INNER JOIN "shop" as "s" ON "s"."id" = "p"."shop_id"
	LEFT JOIN (
		SELECT * FROM "voucher"
		WHERE (now() BETWEEN "voucher"."actived_date" AND "voucher"."expired_date") AND "voucher"."quota" > 0
	) as "v" ON "v"."shop_id" = "s"."id"
	INNER JOIN "category" as "c" ON "c"."id" = "p"."category_id"
	WHERE 
	"p".title ILIKE $1 
	AND  "c".name ILIKE $2
	AND ("p".rating_avg BETWEEN $3 AND $4)
	AND ("p".min_price BETWEEN $5 AND $6)
	AND "f"."user_id" = $7
	AND "p"."deleted_at" IS NULL
	AND "p"."listed_status" = true
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
	 AND "p"."deleted_at" IS NULL
	 AND "p"."listed_status" = true `

	CountUserFavoriteProduct = `
	SELECT count(user_id) FROM "favorite" WHERE "user_id" = $1  AND "product_id"  = $2 
	`
	CountSpecificFavoriteProduct = `
	SELECT count(user_id) FROM "favorite" WHERE "product_id" = $1 
	`
	CreateFavoriteProductQuery = `
	INSERT INTO "favorite" ("user_id", "product_id")
	VALUES ($1, $2);`

	DeleteFavoriteProductQuery = `
	DELETE FROM "favorite"
	WHERE "user_id" = $1 AND "product_id" = $2;`

	CheckFavoriteProductIsExistQuery = `
	SELECT
		CASE WHEN EXISTS 
		(
			SELECT "user_id", "product_id"
			FROM "favorite"
			WHERE "user_id" = $1 AND "product_id" = $2
		)
		THEN 'TRUE'
		ELSE 'FALSE'
	END;
	`

	FindFavoriteProductQuery = `
	SELECT "user_id", "product_id"
	FROM "favorite"
	WHERE "user_id" = $1 AND "product_id" = $2;`

	GetAllTotalReviewProductQuery = `
	SELECT count(r.id)
	FROM review r
	WHERE r.product_id = $1
	%s
	and r.deleted_at IS NULL;`

	GetReviewProductQuery = `
	SELECT r.id, r.user_id, r.product_id, r.comment, r.rating, r.image_url, r.created_at, u.photo_url, u.username
	FROM review r
	INNER JOIN "user" u
	ON r.user_id = u.id
	WHERE r.product_id = $1
	%s
	AND r.deleted_at IS NULL
	ORDER BY %s LIMIT $2 OFFSET $3;`

	GetReviewProductByIDQuery = `
	SELECT r.id, r.user_id, r.product_id, r.comment, r.rating, r.image_url, r.created_at, u.photo_url, u.username
	FROM review r
	INNER JOIN "user" u
	ON r.user_id = u.id
	WHERE r.id = $1
	AND r.deleted_at IS NULL;`

	GetTotalReviewRatingByProductIDQuery = `
	SELECT r.rating, count(r.id) as count 
	FROM review r
	WHERE r.product_id = $1
	and r.deleted_at IS NULL
	group by r.rating;`

	CreateReviewQuery = `INSERT INTO "review" (user_id, product_id, comment, rating, image_url)
	VALUES ($1, $2, $3, $4, $5);`

	DeleteReviewByIDQuery = `UPDATE "review" set deleted_at = now() WHERE id = $1;`

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

	CreateVariantDetailQuery = `INSERT INTO "variant_detail" 
	(name, type)
	 VALUES ($1, $2) RETURNING "id";`

	CreateProductCourierQuery = `INSERT 
	INTO "product_courier_whitelist" 
	(product_id, courier_id)
	 VALUES ($1, $2) RETURNING "id";`

	GetRatingProductQuery = `SELECT 
    	"p"."id", "p"."title", count("r"."id"), avg("r"."rating")  
	FROM "product" as "p" INNER JOIN "review" as "r" on "p"."id" = "r"."product_id" 
	WHERE "r"."deleted_at" IS NULL AND "r"."created_at" >= (now() - interval '1 hour') GROUP BY "p"."id"`

	GetFavoriteProductQuery = `SELECT 
    "p"."id", "p"."title", count("p"."id") 
	FROM "product" as "p" INNER JOIN "favorite" as "f" on "p"."id" = "f"."product_id" 
	WHERE "f"."created_at" >= (now() - interval '1 hour')
	GROUP BY "p"."id"`

	GetListedStatusQuery = `SELECT listed_status from "product" WHERE id = $1 AND deleted_at IS NULL `

	UpdateListedStatusQuery = `UPDATE "product" SET "listed_status" = $1, "updated_at" = now() WHERE "id" = $2 `

	UpdateProductQuery = `UPDATE 
	"product" SET 
	"title" =$1,"description"=$2,
	"thumbnail_url"= $3,
	"min_price"=$4,
	"max_price"=$5,
	"listed_status"=$6, 
	"updated_at" = now()
	WHERE "id" = $7`

	UpdateProductFavoriteQuery = `UPDATE "product" SET "favorite_count" = $1 WHERE "id" = $2`
	UpdateProductRatingQuery   = `UPDATE "product" SET "rating_avg" = $1 WHERE "id" = $2`

	UpdateProductDetailQuery = `UPDATE 
	"product_detail" SET 
	"price" = $1,
	"stock" =$2,
	"weight"=$3,
	"size"= $4,
	"hazardous"=$5,
	"condition"=$6,
	"bulk_price"=$7,
	"updated_at" = now()
	WHERE "id" = $8 AND
	"product_id" = $9`

	DeleteProductDetailByIDQuery = `UPDATE "product_detail" set deleted_at = now() WHERE id = $1`

	DeleteVariantByIDQuery = `UPDATE "variant" set deleted_at = now() WHERE id = $1`

	DeletePhotoByIDQuery = `
	DELETE FROM "photo" WHERE "product_detail_id" = $1`

	GetMaxMinPriceQuery = `
	SELECT max(price), min(price) 
	FROM product_detail
	WHERE product_id = $1
	and deleted_at IS NULL;`

	UpdateVariantQuery = `UPDATE 
	"variant" SET  "variant_detail_id" = $1, "updated_at" = now()
	WHERE "id" = $2`
)
