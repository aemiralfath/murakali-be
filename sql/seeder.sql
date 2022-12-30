-- TODO: add transaction seeder

INSERT INTO "banner" (id, title, content, image_url, page_url, is_active)
VALUES ('9a4de0f5-0556-491b-ba1b-0873933262da', 'festival hari ibu', 'festival di hari ibu',
        'https://cf.shopee.co.id/file/776fa8ed99c660a7913666544a3c228d', 'https://shopee.co.id/m/festival-hari-ibu',
        true),
       ('37817b2c-c4b0-42ec-8422-055d44e47fbe', 'tanggal tua', 'diskon di tanggal tua',
        'https://cf.shopee.co.id/file/f73747bf997bd9d0f20f0b33727f018e',
        'https://shopee.co.id/m/mall-super-category-day', true);

INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('e6fb2764-076f-4b3b-bc05-8aa125d537ed','159aa7d7-2fa0-4cc8-a708-3328d1d08eb5','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Masker KN95 PM2.5 Earloop kn 95 filter 95% setara n95 4 ply isi 50 Pcs',NULL,46500,3874,100000,TRUE,' ',4.8,70000,170000,'2022-12-22 21:43:56.228411+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('910dbb9b-53d5-4a23-b8ca-cf3f0caab169','63f58102-9cb6-4249-b8d4-82f65f315c59','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'ROSE BRAND Minyak Goreng 2 liter',NULL,5813,5774,10000,TRUE,' ',5,40000,50000,'2022-12-22 21:43:56.237079+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('855f81a1-c7aa-4428-9f26-74fd664de377','5d5bd121-adc2-4f62-9cad-d4172bec9a40','07315003-5369-465f-9f05-09482d951645',NULL,'Nice Facial Tissue Tisue Tisu Wajah Muka Paket isi 5 x 180 sheets',NULL,3509,12736,10000,TRUE,' ',5,32000,42000,'2022-12-22 21:43:56.241554+00',NULL,NULL);

INSERT INTO "variant_detail" (id, name, type)
VALUES ('2abaad66-0eae-46bf-8ec7-b6b1d8d6472e', 'ukuran', 's'),
       ('fe997ee9-6bb1-4ddc-bff9-4a2bc0d2ad24', 'ukuran', 'l'),
       ('4e18d196-f84a-4a95-a897-fe085d38347f', 'warna', 'hijau');

-- Update the id with generated id seeder (wait the seeder service in docker finish)
-- or you can insert the id manually first
INSERT INTO "product_detail" (id, product_id, price, stock, weight, size, hazardous, condition, bulk_price)
VALUES ('02d8b878-7124-471a-8f15-0d338ddcfa81', '910dbb9b-53d5-4a23-b8ca-cf3f0caab169', 50000, 5, 1, 1, false, 'wwww',
        false),
       ('0c53ef3d-3682-4359-90e1-814eb6ab5231', 'e6fb2764-076f-4b3b-bc05-8aa125d537ed', 70000, 111, 0.5, 0.1, false,
        'perfect', false),
       ('0c53ef3d-3682-4359-90e1-814eb6ab5111', '855f81a1-c7aa-4428-9f26-74fd664de377', 100000, 100, 1, 1, false,
        'jelek', false);
-- Update the id with generated id seeder

-- Update the id with generated id seeder (wait the seeder service in docker finish)
-- or you can insert the id manually first
INSERT INTO "variant" (id, product_detail_id, variant_detail_id)
VALUES ('f07df23d-7819-4c08-92b1-9702fda48f3e', '02d8b878-7124-471a-8f15-0d338ddcfa81',
        '2abaad66-0eae-46bf-8ec7-b6b1d8d6472e'),
       ('614fa874-d6a7-414a-9ebc-9b43cf765745', '0c53ef3d-3682-4359-90e1-814eb6ab5231',
        'fe997ee9-6bb1-4ddc-bff9-4a2bc0d2ad24'),
       ('ff4b9968-39ba-48f3-ad77-088a118b7c4a', '0c53ef3d-3682-4359-90e1-814eb6ab5111',
        '4e18d196-f84a-4a95-a897-fe085d38347f');
-- Update the id with generated id seeder

-- Update the id with generated id seeder (wait the seeder service in docker finish)
-- or you can insert the id manually first
INSERT INTO "promotion" (id, name, product_id, discount_percentage, discount_fix_price, min_product_price,
                         max_discount_price, quota, max_quantity, actived_date, expired_date)
VALUES ('17d446f3-e35d-46c7-8d0c-252462ca6414', 'promo murah', '855f81a1-c7aa-4428-9f26-74fd664de377', 25, 30000, 50000,
        30000, 10, 1, '2022-12-21 00:00:00-07', '2023-02-01 00:00:00-07');
-- Update the id with generated id seeder

-- Update the id with generated id seeder (wait the seeder service in docker finish)
-- or you can insert the id manually first
INSERT INTO voucher(id, shop_id, code, quota, actived_date, expired_date, discount_percentage, discount_fix_price,
                    min_product_price, max_discount_price, created_at, updated_at, deleted_at)
VALUES ('59bfcd74-e278-4c70-a889-d9c8515bf71c', 'e8854443-c2c7-488e-93d5-b9d93708b8a3', 'ASD123', 11,
        '2022-12-22 21:43:56.202214+00', '2023-02-01 21:43:56.202214+00', 0, 5000, 5000, 100000,
        '2022-12-23 02:34:38.854025+00', NULL, NULL),
        ('59bfcd74-e278-4c70-a889-d9c8515bf72c', '07315003-5369-465f-9f05-09482d951645', 'DSA123', 11,
        '2022-12-22 21:43:56.202214+00', '2023-02-01 21:43:56.202214+00', 0, 5000, 5000, 100000,
        '2022-12-23 02:34:38.854025+00', NULL, NULL);
-- Update the id with generated id seeder

-- Update the id with generated id seeder (wait the seeder service in docker finish)
-- or you can insert the id manually first
INSERT INTO "cart_item" (id, user_id, product_detail_id, quantity)
VALUES ('410a1545-1834-4dca-9624-8c1c7e1439de', '7950eca2-58d5-44f0-b873-22b23d8107da',
        '0c53ef3d-3682-4359-90e1-814eb6ab5231', 4),
       ('d37f3e57-94d9-433c-a3c8-316f6b7194d8', '7950eca2-58d5-44f0-b873-22b23d8107da',
        '02d8b878-7124-471a-8f15-0d338ddcfa81', 1),
       ('d37f3e57-94d9-433c-a3c8-316f6b719418', '7950eca2-58d5-44f0-b873-22b23d8107da',
        '0c53ef3d-3682-4359-90e1-814eb6ab5111', 2);
-- Update the id with generated id seeder

INSERT INTO "wallet" (id, user_id, balance, pin, attempt_count)
VALUES ('60a54e99-33a7-40d8-8ed0-979413a8c33d','7950eca2-58d5-44f0-b873-22b23d8107da', 1000000, '123456', 0);

INSERT INTO "order_status" (id, name)
VALUES (1, 'pending to pay');