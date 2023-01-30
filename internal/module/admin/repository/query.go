package repository

const (
	CountCodeVoucher = `
	SELECT count(code) FROM "voucher" as "v" WHERE "v"."code" = $1  AND "v"."deleted_at" IS NULL
	`
	GetTotalVoucherQuery = `
	SELECT count(id) FROM "voucher" as "v" WHERE "v"."shop_id" IS NULL  AND "v"."deleted_at" IS NULL
	`
	GetTotalRefundsQuery = `SELECT count(id) FROM "refund" WHERE "accepted_at" IS NOT NULL AND "rejected_at" IS NULL AND "refunded_at" IS NULL`
	GetRefundsQuery      = `SELECT 
	"r"."id", "r"."order_id", "r"."is_seller_refund", "r"."is_buyer_refund", "r"."reason", "r"."image", 
	"r"."accepted_at", "r"."rejected_at", "r"."refunded_at", "o"."id", "o"."transaction_id", "o"."shop_id", "o"."user_id", 
	"o"."courier_id", "o"."voucher_shop_id", "o"."order_status_id", "o"."total_price", "o"."delivery_fee", "o"."resi_no", 
	"o"."created_at", "o"."arrived_at" 
	FROM "refund" as "r"
	INNER JOIN "order" as "o" on "r"."order_id" = "o"."id"
	WHERE "accepted_at" IS NOT NULL AND "rejected_at" IS NULL AND "refunded_at" IS NULL`

	GetAllVoucherQuery = `
	SELECT "v"."id", "v"."code", "v"."quota", "v"."actived_date", "v"."expired_date",
		"v"."discount_percentage", "v"."discount_fix_price", "v"."min_product_price", "v"."max_discount_price",
		"v"."created_at", "v"."updated_at",  "v"."deleted_at"
	FROM "voucher" as "v"
	WHERE "v"."shop_id"  IS NULL 
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

	CreateVoucherQuery = `INSERT INTO "voucher" 
    	( code, quota, actived_date, expired_date, discount_percentage, discount_fix_price, min_product_price, max_discount_price)
    	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	DeleteVoucherQuery = `UPDATE "voucher" set deleted_at = now() WHERE "id" = $1 AND "shop_id"  IS NULL  AND "deleted_at" IS NULL`

	GetVoucherByID = `
	SELECT "v"."id",  "v"."code", "v"."quota", "v"."actived_date", "v"."expired_date",
		"v"."discount_percentage", "v"."discount_fix_price", "v"."min_product_price", "v"."max_discount_price",
		"v"."created_at", "v"."updated_at",  "v"."deleted_at"
	FROM "voucher" as "v"
	WHERE "v"."id"  = $1 AND "v"."shop_id" IS NULL  AND "v"."deleted_at" IS NULL
	`

	GetRefundByIDQuery = `SELECT 
	"id", "order_id", "is_seller_refund", "is_buyer_refund", "reason", "image", 
	"accepted_at", "rejected_at", "refunded_at" FROM "refund" WHERE "id" = $1`

	UpdateVoucherQuery = `
		UPDATE "voucher" SET "quota" = $1, "actived_date" = $2, "expired_date" = $3, "discount_percentage" = $4,
			"discount_fix_price" = $5, "min_product_price" = $6, "max_discount_price" = $7, "updated_at" = now()
		WHERE "id" = $8
	`

	CreateWalletHistoryQuery      = `INSERT INTO "wallet_history" (transaction_id, wallet_id, "from", "to", description, amount, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	UpdateRefundQuery             = `UPDATE "refund" SET "refunded_at" = $1 WHERE "id" = $2`
	UpdateOrderByID               = `UPDATE "order" SET "order_status_id" = $1 WHERE "id" = $2`
	UpdateProductDetailStockQuery = `UPDATE "product_detail" SET "stock" = $1, "updated_at" = now() WHERE "id" = $2;`
	UpdateWalletBalanceQuery      = `UPDATE "wallet" SET "balance" = $1, "updated_at" = $2 WHERE "id" = $3`

	GetOrderByOrderIDQuery      = `SELECT o.id,o.order_status_id, o.user_id, o.transaction_id,o.total_price,o.delivery_fee,o.resi_no,o.created_at from "order" o WHERE o.id = $1`
	GetOrderItemsByOrderIDQuery = `SELECT "id", "order_id", "product_detail_id", "quantity", "item_price", "total_price" FROM "order_item" WHERE "order_id" = $1`
	GetProductDetailByIDQuery   = `SELECT "id", "price", "stock", "weight", "size", "hazardous", "condition", "bulk_price" FROM "product_detail" WHERE "id" = $1 AND "deleted_at" IS NULL;`
	GetWalletByUserIDQuery      = `SELECT "id", "user_id", "balance", "pin", "attempt_count", "attempt_at", "unlocked_at", "active_date" FROM "wallet" WHERE "user_id" = $1 AND "deleted_at" IS NULL`

	GetCategoriesQuery = `WITH RECURSIVE ctgry AS (
		SELECT id, parent_id, name, photo_url, created_at, updated_at, deleted_at, 1 as level
		FROM category
		WHERE parent_id IS NULL
		UNION ALL
		SELECT t.id, t.parent_id, t.name, t.photo_url, t.created_at, t.updated_at, t.deleted_at, ctgry.level + 1
		FROM category t
		JOIN ctgry ON t.parent_id = ctgry.id
	)
	SELECT id, parent_id, name, photo_url,level
	FROM ctgry
	WHERE level <= 3  and deleted_at is null`

	AddCategoryQuery = `INSERT INTO "category" 
	( parent_id, name, photo_url)
	VALUES ($1, $2, $3)`

	DeleteCategoryQuery       = `UPDATE "category" set deleted_at = now() WHERE "id" = $1 AND "deleted_at" IS NULL`
	EditCategoryQuery         = `UPDATE "category" set parent_id = $1, name = $2 , photo_url = $3,  updated_at = now() WHERE "id" = $4 AND "deleted_at" IS NULL`
	CountProductCategoryQuery = `SELECT count(1) from product where category_id = $1 and deleted_at is null`
	CountCategoryParentQuery  = `SELECT count(1) from category where parent_id = $1 and deleted_at is null`
	GetBannerQuery            = `SELECT id,title,content,image_url,page_url,is_active FROM "banner" order by id`
	AddBannerQuery            = `INSERT INTO "banner" 
	( title, content, image_url, page_url, is_active)
	VALUES ($1, $2, $3, $4, $5)`
	DeleteBannerQuery = `DELETE FROM "banner" WHERE id = $1`
	EditBannerQuery   = `UPDATE "banner" set is_active = $1 WHERE "id" = $2`
)
