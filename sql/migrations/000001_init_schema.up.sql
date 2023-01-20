CREATE
EXTENSION IF NOT EXISTS pg_trgm;
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "user"
(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "role_id" int,
    "username" varchar UNIQUE,
    "email" varchar UNIQUE NOT NULL,
    "phone_no" varchar UNIQUE,
    "fullname" varchar,
    "password" varchar NOT NULL DEFAULT '',
    "gender" varchar,
    "birth_date" timestamp,
    "photo_url" varchar,
    "is_sso" boolean DEFAULT false,
    "is_verify" boolean DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "shop"
(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "user_id" UUID,
    "name" varchar UNIQUE NOT NULL,
    "total_product" int,
    "total_rating" int,
    "rating_avg" float,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "shop_courier"
(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "shop_id" UUID,
    "courier_id" uuid,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "product_courier_whitelist"
(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "product_id" UUID,
    "courier_id" uuid,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "courier"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "name" varchar,
    "code" varchar,
    "service" varchar,
    "description" text,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "cart_item"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "user_id" UUID,
    "product_detail_id" UUID,
    "quantity" float,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "order_status"
(
    "id"
    serial
    PRIMARY
    KEY,
    "name"
    varchar,
    "created_at"
    timestamptz
    NOT
    NULL
    DEFAULT (
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "order_item"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "order_id" UUID,
    "product_detail_id" UUID,
    "quantity" int,
    "item_price" float,
    "total_price" float,
    "note" varchar NOT NULL DEFAULT '',
    "is_review" boolean NOT NULL DEFAULT FALSE
    );

CREATE TABLE IF NOT EXISTS "order"
(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "transaction_id" UUID,
    "shop_id" UUID,
    "user_id" UUID,
    "courier_id" UUID,
    "voucher_shop_id" uuid,
    "order_status_id" int,
    "total_price" float,
    "delivery_fee" float,
    "resi_no" varchar,
    "buyer_address" varchar NOT NULL DEFAULT '',
    "shop_address" varchar NOT NULL DEFAULT '',
    "cancel_notes" varchar NOT NULL DEFAULT '',
    "is_withdraw" boolean NOT NULL DEFAULT FALSE,
    "is_refund" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "arrived_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "transaction"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "voucher_marketplace_id" uuid,
    "wallet_id" UUID,
    "card_number" varchar,
    "invoice" varchar UNIQUE,
    "total_price" float,
    "paid_at" timestamptz,
    "canceled_at" timestamptz,
    "expired_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "promotion"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "name" varchar,
    "product_id" uuid,
    "discount_percentage" float,
    "discount_fix_price" int,
    "min_product_price" float,
    "max_discount_price" float,
    "quota" int,
    "max_quantity" int,
    "actived_date" timestamptz,
    "expired_date" timestamptz,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "voucher"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "shop_id" uuid,
    "code" varchar,
    "quota" int,
    "actived_date" timestamptz,
    "expired_date" timestamptz,
    "discount_percentage" int,
    "discount_fix_price" int,
    "min_product_price" float,
    "max_discount_price" float,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "favorite"
(
    "user_id"
    uuid,
    "product_id"
    uuid,
    "created_at"
    timestamptz
    NOT
    NULL
    DEFAULT (
    NOW
(
))
    );

CREATE TABLE IF NOT EXISTS "product"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "category_id" uuid,
    "shop_id" uuid,
    "sku" varchar,
    "title" varchar,
    "description" varchar,
    "view_count" bigint,
    "favorite_count" bigint,
    "unit_sold" bigint,
    "listed_status" boolean DEFAULT false,
    "thumbnail_url" varchar,
    "rating_avg" float,
    "min_price" float,
    "max_price" float,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "product_detail"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "product_id" uuid,
    "price" float,
    "stock" bigint,
    "weight" float,
    "size" float,
    "hazardous" boolean DEFAULT false,
    "condition" varchar,
    "bulk_price" boolean DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "variant"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "product_detail_id" uuid,
    "variant_detail_id" uuid,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "variant_detail"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "name" varchar,
    "type" varchar,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "category"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "parent_id" uuid,
    "name" varchar,
    "photo_url" varchar,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "review"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "user_id" uuid,
    "product_id" uuid,
    "comment" text,
    "rating" float,
    "image_url" varchar,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "photo"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "product_detail_id" uuid,
    "url" varchar
    );

CREATE TABLE IF NOT EXISTS "video"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "product_detail_id" uuid,
    "url" varchar
    );

CREATE TABLE IF NOT EXISTS "sealabs_pay"
(
    "card_number" varchar PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "user_id" UUID,
    "name" varchar,
    "is_default" boolean DEFAULT false,
    "active_date" timestamp,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "wallet"
(
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "user_id" UUID,
    "balance" float NOT NULL DEFAULT 0,
    "pin" varchar,
    "attempt_count" int,
    "attempt_at" timestamptz,
    "unlocked_at" timestamptz,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "active_date" timestamptz
    );

CREATE TABLE IF NOT EXISTS "wallet_history"
(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "transaction_id" UUID,
    "wallet_id" UUID,
    "from" varchar,
    "to" varchar,
    "amount" float,
    "description" text,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
))
    );

CREATE TABLE IF NOT EXISTS "address"
(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "user_id" UUID,
    "name" varchar,
    "province_id" int,
    "city_id" int,
    "province" varchar,
    "city" varchar,
    "district" varchar,
    "sub_district" varchar,
    "address_detail" varchar,
    "zip_code" varchar,
    "is_default" boolean DEFAULT false,
    "is_shop_default" boolean DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz,
    "deleted_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "role"
(
    "id"
    serial
    PRIMARY
    KEY,
    "name"
    varchar,
    "created_at"
    timestamptz
    NOT
    NULL
    DEFAULT (
    NOW
(
)),
    "updated_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "email_history"
(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "email" varchar UNIQUE NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT
(
    NOW
(
)),
    "updated_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "banner"
(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "title" varchar,
    "content" text,
    "image_url" varchar,
    "page_url" varchar,
    "is_active" boolean DEFAULT false
    );

CREATE TABLE IF NOT EXISTS "refund"
(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "order_id" UUID,
    "is_seller_refund" boolean DEFAULT false,
    "is_buyer_refund" boolean DEFAULT false,
    "reason" varchar NOT NULL DEFAULT '',
    "image" varchar NOT NULL DEFAULT '',
    "accepted_at" timestamptz,
    "rejected_at" timestamptz,
    "refunded_at" timestamptz
    );

CREATE TABLE IF NOT EXISTS "refund_thread"
(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4
(
),
    "refund_id" UUID,
    "user_id" UUID,
    "is_seller" boolean DEFAULT false,
    "is_buyer" boolean DEFAULT false,
    "text" varchar NOT NULL DEFAULT '',
    "created_at" timestamptz
    );

CREATE INDEX ON "refund_thread" ("refund_id");

CREATE INDEX ON "refund_thread" ("user_id");

CREATE INDEX ON "refund" ("order_id");

CREATE INDEX ON "user" ("email");

CREATE INDEX ON "user" ("phone_no");

CREATE INDEX ON "user" ("role_id");

CREATE INDEX ON "shop" ("user_id");

CREATE INDEX ON "shop" ("name");

CREATE INDEX ON "shop_courier" ("shop_id");

CREATE INDEX ON "shop_courier" ("courier_id");

CREATE INDEX ON "product_courier_whitelist" ("product_id");

CREATE INDEX ON "product_courier_whitelist" ("courier_id");

CREATE INDEX ON "courier" ("name");

CREATE INDEX ON "cart_item" ("user_id");

CREATE INDEX ON "cart_item" ("product_detail_id");

CREATE INDEX ON "order_status" ("name");

CREATE INDEX ON "order_item" ("order_id");

CREATE INDEX ON "order_item" ("product_detail_id");

CREATE INDEX ON "order" ("transaction_id");

CREATE INDEX ON "order" ("shop_id");

CREATE INDEX ON "order" ("user_id");

CREATE INDEX ON "order" ("courier_id");

CREATE INDEX ON "order" ("voucher_shop_id");

CREATE INDEX ON "order" ("order_status_id");

CREATE INDEX ON "transaction" ("voucher_marketplace_id");

CREATE INDEX ON "transaction" ("wallet_id");

CREATE INDEX ON "transaction" ("card_number");

CREATE INDEX ON "transaction" ("invoice");

CREATE INDEX ON "promotion" ("name");

CREATE INDEX ON "promotion" ("product_id");

CREATE INDEX ON "voucher" ("shop_id");

CREATE INDEX ON "voucher" ("code");

CREATE INDEX ON "favorite" ("user_id");

CREATE INDEX ON "favorite" ("product_id");

CREATE INDEX ON "product" ("category_id");

CREATE INDEX ON "product" ("shop_id");

CREATE INDEX ON "product" ("sku");

CREATE INDEX ON "product" ("title");

CREATE INDEX ON "product_detail" ("product_id");

CREATE INDEX ON "variant" ("product_detail_id");

CREATE INDEX ON "variant" ("variant_detail_id");

CREATE INDEX ON "variant_detail" ("name");

CREATE INDEX ON "category" ("name");

CREATE INDEX ON "category" ("parent_id");

CREATE INDEX ON "review" ("user_id");

CREATE INDEX ON "review" ("product_id");

CREATE INDEX ON "review" ("rating");

CREATE INDEX ON "photo" ("product_detail_id");

CREATE INDEX ON "video" ("product_detail_id");

CREATE INDEX ON "sealabs_pay" ("user_id");

CREATE INDEX ON "wallet" ("user_id");

CREATE INDEX ON "wallet_history" ("wallet_id");

CREATE INDEX ON "address" ("name");

CREATE INDEX ON "address" ("user_id");

CREATE INDEX ON "email_history" ("email");

CREATE INDEX ON "banner" ("title");

CREATE INDEX ON "promotion" ("actived_date", "expired_date");

CREATE INDEX ON "voucher" ("actived_date", "expired_date");

ALTER TABLE "user"
    ADD FOREIGN KEY ("role_id") REFERENCES "role" ("id");

ALTER TABLE "shop"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "shop_courier"
    ADD FOREIGN KEY ("shop_id") REFERENCES "shop" ("id");

ALTER TABLE "shop_courier"
    ADD FOREIGN KEY ("courier_id") REFERENCES "courier" ("id");

ALTER TABLE "product_courier_whitelist"
    ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");

ALTER TABLE "product_courier_whitelist"
    ADD FOREIGN KEY ("courier_id") REFERENCES "courier" ("id");

ALTER TABLE "cart_item"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "cart_item"
    ADD FOREIGN KEY ("product_detail_id") REFERENCES "product_detail" ("id");

ALTER TABLE "order_item"
    ADD FOREIGN KEY ("order_id") REFERENCES "order" ("id");

ALTER TABLE "refund"
    ADD FOREIGN KEY ("order_id") REFERENCES "order" ("id");

ALTER TABLE "refund_thread"
    ADD FOREIGN KEY ("refund_id") REFERENCES "refund" ("id");

ALTER TABLE "refund_thread"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "order_item"
    ADD FOREIGN KEY ("product_detail_id") REFERENCES "product_detail" ("id");

ALTER TABLE "order"
    ADD FOREIGN KEY ("transaction_id") REFERENCES "transaction" ("id");

ALTER TABLE "order"
    ADD FOREIGN KEY ("shop_id") REFERENCES "shop" ("id");

ALTER TABLE "order"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "order"
    ADD FOREIGN KEY ("courier_id") REFERENCES "courier" ("id");

ALTER TABLE "order"
    ADD FOREIGN KEY ("voucher_shop_id") REFERENCES "voucher" ("id");

ALTER TABLE "order"
    ADD FOREIGN KEY ("order_status_id") REFERENCES "order_status" ("id");

ALTER TABLE "transaction"
    ADD FOREIGN KEY ("voucher_marketplace_id") REFERENCES "voucher" ("id");

ALTER TABLE "transaction"
    ADD FOREIGN KEY ("wallet_id") REFERENCES "wallet" ("id");

ALTER TABLE "transaction"
    ADD FOREIGN KEY ("card_number") REFERENCES "sealabs_pay" ("card_number");

ALTER TABLE "promotion"
    ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");

ALTER TABLE "voucher"
    ADD FOREIGN KEY ("shop_id") REFERENCES "shop" ("id");

ALTER TABLE "favorite"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "favorite"
    ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");

ALTER TABLE "product"
    ADD FOREIGN KEY ("category_id") REFERENCES "category" ("id");

ALTER TABLE "product"
    ADD FOREIGN KEY ("shop_id") REFERENCES "shop" ("id");

ALTER TABLE "product_detail"
    ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");

ALTER TABLE "variant"
    ADD FOREIGN KEY ("product_detail_id") REFERENCES "product_detail" ("id");

ALTER TABLE "variant"
    ADD FOREIGN KEY ("variant_detail_id") REFERENCES "variant_detail" ("id");

ALTER TABLE "category"
    ADD FOREIGN KEY ("parent_id") REFERENCES "category" ("id");

ALTER TABLE "review"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "review"
    ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");

ALTER TABLE "photo"
    ADD FOREIGN KEY ("product_detail_id") REFERENCES "product_detail" ("id");

ALTER TABLE "video"
    ADD FOREIGN KEY ("product_detail_id") REFERENCES "product_detail" ("id");

ALTER TABLE "sealabs_pay"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "wallet"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "wallet_history"
    ADD FOREIGN KEY ("transaction_id") REFERENCES "transaction" ("id");

ALTER TABLE "wallet_history"
    ADD FOREIGN KEY ("wallet_id") REFERENCES "wallet" ("id");

ALTER TABLE "address"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");
