package repository

const (
	GetCategoriesQuery           = `SELECT "id", "parent_id", "name", "photo_url" FROM "category" WHERE "parent_id" IS NULL AND "deleted_at" IS NULL`
	GetCategoriesByNameQuery     = `SELECT "id", "parent_id", "name", "photo_url" FROM "category" WHERE "name" = $1 AND "deleted_at" IS NULL`
	GetCategoriesByParentIdQuery = `SELECT "id", "parent_id", "name", "photo_url" FROM "category" WHERE "parent_id" = $1 AND "deleted_at" IS NULL`
	GetBannersQuery              = `SELECT "id", "title", "content", "image_url", "page_url", "is_active" FROM "banner" WHERE "is_active" = TRUE`
	GetRecommendedProductsQuery  = `
	SELECT "p"."title" as "title", "p"."unit_sold" as "unit_sold", "p"."rating_avg" as "rating_avg", "p"."thumbnail_url" as "thumbnail_url",
		"p"."min_price" as "min_price", "p"."max_price" as "max_price", "promo"."discount_percentage" as "promo_discount_percentage",  "promo"."discount_fix_price" as "promo_discount_fix_price",
		"promo"."min_product_price" as "promo_min_product_price",  "promo"."max_discount_price" as "promo_max_discount_price",
		"v"."discount_percentage" as "voucher_discount_percentage",  "v"."discount_fix_price" as "voucher_discount_fix_price"
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
	WHERE "p"."deleted_at" IS NULL
	ORDER BY "p"."unit_sold" DESC
	limit $1;
	`
	GetProductInfoQuery = `select
	a.id,a.sku,a.title,a.description,a.view_count,a.favorite_count,a.unit_sold,a.listed_status,a.thumbnail_url,a.rating_avg,a.min_price,a.max_price
	,p.name,p.discount_percentage,p.discount_fix_price,p.min_product_price,p.max_discount_price,p.quota,p.max_quantity,p.actived_date,p.expired_date
	,e.parent_id,e.name,e.photo_url
	from 
	product a 
	join product_detail b on a.id = b.product_id 
	join photo g on b.id = g.product_detail_id
	join promotion p on a.id = p.product_id
	join category e on e.id = a.category_id
	where a.id = $1`

	GetProductDetailQuery = `select
	b.id,b.price,b.stock,b.weight,b.size,b.hazardous,b.condition,b.bulk_price,g.url
	from 
	product_detail b 
	join photo g on b.id = g.product_detail_id
	where b.product_id = $1`

	GetVariantDetailQuery = `select b.type,b.name from variant a join variant_detail b on a.variant_detail_id = b.id
	where a.product_detail_id = $1`
)
