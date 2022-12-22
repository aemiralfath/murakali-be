INSERT INTO "role" (name)
VALUES ('user'),
       ('seller'),
       ('admin');

INSERT INTO "email_history" (email)
VALUES ('user@gmail.com'),
       ('seller@gmail.com'),
       ('admin@gmail.com');

INSERT INTO "user" (id, role_id, username, email, phone_no, fullname, password, is_verify)
VALUES ('7950eca2-58d5-44f0-b873-22b23d8107da', 1, 'user', 'user@gmail.com', '911', 'user',
        '$2a$10$cNhdZVN.pgsfK1xUQ00p7eK5Fh7iClrtJB9SY5un.H55Mi/dtQzCa', true),
       ('f8d8d66a-e8eb-4633-bc2d-4ccd941fed47', 2, 'seller', 'seller@gmail.com', '912', 'seller',
        '$2a$10$cNhdZVN.pgsfK1xUQ00p7eK5Fh7iClrtJB9SY5un.H55Mi/dtQzCa', true),
       ('4df967a8-5b05-4d2a-bb72-da3921dce8fb', 3, 'admin', 'admin@gmail.com', '913', 'admin',
        '$2a$10$cNhdZVN.pgsfK1xUQ00p7eK5Fh7iClrtJB9SY5un.H55Mi/dtQzCa', true);

INSERT INTO "address" (user_id, name, province_id, city_id, province, city, district, sub_district, address_detail, zip_code)
VALUES ('7950eca2-58d5-44f0-b873-22b23d8107da', 'test', 5, 39, 'DI Yogyakarta', 'Bantul', 'Pleret', 'Segoroyoso', 'no 91', '55791');

INSERT INTO "category" (id, name, photo_url)
VALUES ('d92a0995-78cd-4eba-a855-dfc096ffec5b', 'laptop', 'https://cf.shopee.co.id/file/c139370836a9daa649da70876a326b58'),
       ('5d5bd121-adc2-4f62-9cad-d4172bec9a40', 'electronik', 'https://cf.shopee.co.id/file/dcd61dcb7c1448a132f49f938b0cb553_tn'),
       ('5778e73c-f8b7-4c6b-a2f4-472079b164c5', 'pakaian pria', 'https://cf.shopee.co.id/file/04dba508f1ad19629518defb94999ef9_tn'),
       ('63f58102-9cb6-4249-b8d4-82f65f315c59', 'hobi & koleksi', 'https://cf.shopee.co.id/file/42394b78fac1169d67c6291973a3b132_tn');

INSERT INTO "category" (id, name, parent_id, photo_url)
VALUES ('d99373d1-c55d-4769-a56e-f797db20235d', 'outfit Hangat', '5778e73c-f8b7-4c6b-a2f4-472079b164c5', 'https://cf.shopee.co.id/file/d89df04fd3435962af59be0408ec4786'),
       ('159aa7d7-2fa0-4cc8-a708-3328d1d08eb5', 'sweeter', 'd99373d1-c55d-4769-a56e-f797db20235d', 'https://cf.shopee.co.id/file/19b8238c917f3dec99b689809ea43a79_tn'),
       ('0774dbda-194f-439d-97e3-eec0e325fe5a','mouse & keyboard', '5d5bd121-adc2-4f62-9cad-d4172bec9a40', 'https://cf.shopee.co.id/file/6e70e5f93250a3e8344cda4fc79b0c3f'),
       ('1aaaed1f-9d23-47ef-8647-17b862becc27','webcam',  '5d5bd121-adc2-4f62-9cad-d4172bec9a40', 'https://cf.shopee.co.id/file/45ee92cbf6243007a66f0f338058da80');

INSERT INTO "shop" (id, user_id, name, total_product, total_rating, rating_avg)
VALUES
('e8854443-c2c7-488e-93d5-b9d93708b8a3', 'f8d8d66a-e8eb-4633-bc2d-4ccd941fed47', 'jualan', 1, 1, 1)

INSERT INTO "product" (id, category_id, shop_id, sku, title, description, view_count, favorite_count, unit_sold, listed_status, thumbnail_url, rating_avg, min_price, max_price)
VALUES
('d6489799-9cc3-4480-9517-7b226a120f08', 'f7c9c484-18a8-4ad6-b214-86363225ead4', 'e8854443-c2c7-488e-93d5-b9d93708b8a3', 'sku product', 'nama produk', 'deskripsi produk', 10, 5, 1, 1, 'https://cf.shopee.co.id/file/76a0969b7d64065bc13493bf55df1849_tn', 1, 10000, 100000)

INSERT INTO "product_detail" (id, product_id, price, stock, weight, size, hazardous, condition, bulk_price)
VALUES
('0c53ef3d-3682-4359-90e1-814eb6ab5191', 'd6489799-9cc3-4480-9517-7b226a120f08', 10000, 4, 2, 2, false, 'bagus', false)

INSERT INTO "variant_detail" (id, name, type)
VALUES
('b11feaf3-8776-4a99-9230-7b90fa310ef5' ,'warna' ,'hitam'),
('049ffb82-3f3e-4dd4-bec2-216d43151f51' ,'warna' ,'putih');

INSERT INTO "variant" (id, product_detail_id, variant_detail_id)
VALUES
('0f8773e1-338e-4b58-ab4c-4523fecae9ca', '0c53ef3d-3682-4359-90e1-814eb6ab5191', 'b11feaf3-8776-4a99-9230-7b90fa310ef5');

INSERT INTO "promotion" (id, name, product_id, discount_percentage, discount_fix_price, min_product_price, max_discount_price, quota, max_quantity, actived_date, expired_date)
VALUES
('17d446f3-e35d-46c7-8d0c-252462ca6414', 'promo murah', 'd6489799-9cc3-4480-9517-7b226a120f08', 25, 30000, 50000, 30000, 10, 1, '2022-12-21 00:00:00-07', '2023-02-01 00:00:00-07');

INSERT INTO "cart_item" (id, user_id, product_detail_id, quantity)
VALUES
('d37f3e57-94d9-433c-a3c8-316f6b7194d8' ,'7950eca2-58d5-44f0-b873-22b23d8107da', '0c53ef3d-3682-4359-90e1-814eb6ab5191', 1);
