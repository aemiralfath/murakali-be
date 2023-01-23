package repository

const (
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

	UpdateVoucherQuery = `
		UPDATE "voucher" SET "quota" = $1, "actived_date" = $2, "expired_date" = $3, "discount_percentage" = $4,
			"discount_fix_price" = $5, "min_product_price" = $6, "max_discount_price" = $7, "updated_at" = now()
		WHERE "id" = $8
	`
)
