INSERT INTO "banner" (id, title, content, image_url, page_url, is_active)
VALUES ('9a4de0f5-0556-491b-ba1b-0873933262da', 'festival hari ibu', 'festival di hari ibu',
        'https://cf.shopee.co.id/file/776fa8ed99c660a7913666544a3c228d', 'https://shopee.co.id/m/festival-hari-ibu',
        true),
       ('37817b2c-c4b0-42ec-8422-055d44e47fbe', 'tanggal tua', 'diskon di tanggal tua',
        'https://cf.shopee.co.id/file/f73747bf997bd9d0f20f0b33727f018e',
        'https://shopee.co.id/m/mall-super-category-day', true);

-- product more than 100k, courier, courier product, shop courier, transaction,

INSERT INTO "promotion" (id, name, product_id, discount_percentage, discount_fix_price, min_product_price,
                         max_discount_price, quota, max_quantity, actived_date, expired_date)
VALUES ('17d446f3-e35d-46c7-8d0c-252462ca6414', 'promo murah', 'd6489799-9cc3-4480-9517-7b226a120f08', 25, 30000, 50000,
        30000, 10, 1, '2022-12-21 00:00:00-07', '2023-02-01 00:00:00-07');

INSERT INTO "cart_item" (id, user_id, product_detail_id, quantity)
VALUES ('d37f3e57-94d9-433c-a3c8-316f6b7194d8', '7950eca2-58d5-44f0-b873-22b23d8107da',
        '0c53ef3d-3682-4359-90e1-814eb6ab5191', 1);

