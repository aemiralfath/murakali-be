INSERT INTO "role" (name)
VALUES ('user'),
       ('seller'),
       ('admin');

INSERT INTO "email_history" (email)
VALUES ('fadhlan1337@gmail.com'),
       ('user@gmail.com'),
       ('seller@gmail.com'),
       ('admin@gmail.com');

INSERT INTO "user" (id, role_id, username, email, phone_no, fullname, password, is_verify)
VALUES ('4c1d6464-3cc6-44d6-92d1-91aee337e025', 1, 'fadhlan', 'fadhlan1337@gmail.com', '910', 'fadhlan',
        '$2a$10$cNhdZVN.pgsfK1xUQ00p7eK5Fh7iClrtJB9SY5un.H55Mi/dtQzCa', true),
       ('7950eca2-58d5-44f0-b873-22b23d8107da', 1, 'user', 'user@gmail.com', '911', 'user',
        '$2a$10$cNhdZVN.pgsfK1xUQ00p7eK5Fh7iClrtJB9SY5un.H55Mi/dtQzCa', true),
       ('f8d8d66a-e8eb-4633-bc2d-4ccd941fed47', 2, 'seller', 'seller@gmail.com', '912', 'seller',
        '$2a$10$cNhdZVN.pgsfK1xUQ00p7eK5Fh7iClrtJB9SY5un.H55Mi/dtQzCa', true),
       ('4df967a8-5b05-4d2a-bb72-da3921dce8fb', 3, 'admin', 'admin@gmail.com', '913', 'admin',
        '$2a$10$cNhdZVN.pgsfK1xUQ00p7eK5Fh7iClrtJB9SY5un.H55Mi/dtQzCa', true);

INSERT INTO "address" (user_id, name, province_id, city_id, province, city, district, sub_district, address_detail,
                       zip_code)
VALUES ('7950eca2-58d5-44f0-b873-22b23d8107da', 'test', 5, 39, 'DI Yogyakarta', 'Bantul', 'Pleret', 'Segoroyoso',
        'no 91', '55791');

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
('e8854443-c2c7-488e-93d5-b9d93708b8a3', 'f8d8d66a-e8eb-4633-bc2d-4ccd941fed47', 'jualan', 1, 1, 1);

INSERT INTO "product" (id, category_id, shop_id, sku, title, description, view_count, favorite_count, unit_sold, listed_status, thumbnail_url, rating_avg, min_price, max_price)
VALUES
('d6489799-9cc3-4480-9517-7b226a120f08', '0774dbda-194f-439d-97e3-eec0e325fe5a', 'e8854443-c2c7-488e-93d5-b9d93708b8a3', 'sku product', 'nama produk', 'deskripsi produk', 10, 5, 1, TRUE, 'https://cf.shopee.co.id/file/76a0969b7d64065bc13493bf55df1849_tn', 1, 10000, 100000);

INSERT INTO "product_detail" (id, product_id, price, stock, weight, size, hazardous, condition, bulk_price)
VALUES
('0c53ef3d-3682-4359-90e1-814eb6ab5191', 'd6489799-9cc3-4480-9517-7b226a120f08', 10000, 4, 2, 2, false, 'bagus', false);

INSERT INTO "variant_detail" (id, name, type)
VALUES
('b11feaf3-8776-4a99-9230-7b90fa310ef5' ,'warna' ,'hitam'),
('b11feaf3-8776-4a99-9230-7b90fa310ef1' ,'ukuran' ,'xl'),
('049ffb82-3f3e-4dd4-bec2-216d43151f51' ,'warna' ,'putih');

INSERT INTO "promotion" (id, name, product_id, discount_percentage, discount_fix_price, min_product_price, max_discount_price, quota, max_quantity, actived_date, expired_date)
VALUES
('17d446f3-e35d-46c7-8d0c-252462ca6414', 'promo murah', 'd6489799-9cc3-4480-9517-7b226a120f08', 25, 30000, 50000, 30000, 10, 1, '2022-12-21 00:00:00-07', '2023-02-01 00:00:00-07');

INSERT INTO "variant" (id, product_detail_id, variant_detail_id)
VALUES
('0f8773e1-338e-4b58-ab4c-4523fecae9cb', '0c53ef3d-3682-4359-90e1-814eb6ab5191', 'b11feaf3-8776-4a99-9230-7b90fa310ef1'),
('0f8773e1-338e-4b58-ab4c-4523fecae9ca', '0c53ef3d-3682-4359-90e1-814eb6ab5191', 'b11feaf3-8776-4a99-9230-7b90fa310ef5');

INSERT INTO "cart_item" (id, user_id, product_detail_id, quantity)
VALUES
('d37f3e57-94d9-433c-a3c8-316f6b7194d8' ,'7950eca2-58d5-44f0-b873-22b23d8107da', '0c53ef3d-3682-4359-90e1-814eb6ab5191', 1);

INSERT INTO "banner" (id, title, content, image_url, page_url, is_active)
VALUES
('9a4de0f5-0556-491b-ba1b-0873933262da' ,'festival hari ibu', 'festival di hari ibu', 'https://cf.shopee.co.id/file/776fa8ed99c660a7913666544a3c228d','https://shopee.co.id/m/festival-hari-ibu', true),
('37817b2c-c4b0-42ec-8422-055d44e47fbe' ,'tanggal tua', 'diskon di tanggal tua', 'https://cf.shopee.co.id/file/f73747bf997bd9d0f20f0b33727f018e','https://shopee.co.id/m/mall-super-category-day', true);

INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('f6f42748-dc03-4cae-a607-825a59af295a',NULL,'Sereal',NULL,'2022-12-22 21:20:08.264379+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('4e34d330-63ad-4413-a058-6777a99ba51b',NULL,'Baterai',NULL,'2022-12-22 21:20:08.276933+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('ccddcbb2-7a62-4e9d-9ce3-5b2e6e6435fb',NULL,'Tas & Jaring Helm',NULL,'2022-12-22 21:20:08.278861+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('0fdbf93b-f71b-4b41-8529-f011c1186f0f',NULL,'Produk Lainnya',NULL,'2022-12-22 21:20:08.28074+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('a25dad10-4005-48b5-b46f-442f421ef932',NULL,'Parfum Mobil',NULL,'2022-12-22 21:20:08.28285+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('b3777c29-de06-4b6f-abf9-9542af906a93',NULL,'Bunga Plastik',NULL,'2022-12-22 21:20:08.285801+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('92b839b7-4088-4509-818b-fa8dd740863a',NULL,'Wash and Wax',NULL,'2022-12-22 21:20:08.288731+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('48833905-98e8-4057-a076-efbc5585656a',NULL,'Sabun Cuci Piring',NULL,'2022-12-22 21:20:08.290746+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('be3c27fc-b20e-48c1-a5b9-f92227548d04',NULL,'Masker Medis',NULL,'2022-12-22 21:20:08.29343+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('4f1a49dd-4df8-4f49-adab-c76b3700e5b7',NULL,'Wiper & Wiper Cover Mobil',NULL,'2022-12-22 21:20:08.296112+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('31f8a451-4c64-422e-8774-5044ffa123bc',NULL,'Obat Herbal',NULL,'2022-12-22 21:20:08.298714+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('e0cef9ef-369e-461c-9bcf-af4b51c5e2f8',NULL,'Minyak',NULL,'2022-12-22 21:20:08.301366+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('2617472a-1ab3-4cb6-9562-34d2f076dc9b',NULL,'Mainan Montessori',NULL,'2022-12-22 21:20:08.303442+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('f3d60dfd-cb4c-4727-b509-e0019a34f57d',NULL,'Alat Diagnosa',NULL,'2022-12-22 21:20:08.305788+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('3e3a6fa8-d387-44d2-a8f1-4768bf658438',NULL,'Buket Bunga',NULL,'2022-12-22 21:20:08.310189+00',NULL,NULL);
INSERT INTO category(id,parent_id,name,photo_url,created_at,updated_at,deleted_at) VALUES ('7e2a3f49-87fc-4b42-aae0-e59e795beee0',NULL,'Tissue',NULL,'2022-12-22 21:20:08.312914+00',NULL,NULL);

INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('f2535899-134d-47e7-bbb5-43c979102ca6','f6f42748-dc03-4cae-a607-825a59af295a','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Cereal SIMBA Choco Chips Bag 1 kg',NULL,5397,821,10000,TRUE,' ',5,45648,55648,'2022-12-22 21:43:56.202214+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('f25293f2-7c97-4db9-ba23-38f0740ce146','e0cef9ef-369e-461c-9bcf-af4b51c5e2f8','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Triple Pack: Tropicana Slim Hokkaido Cheese Cookies (5 Sch)',NULL,9948,2387,10000,TRUE,' ',5,80400,90400,'2022-12-22 21:43:56.212156+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('ffb50fce-e84b-4e68-9316-69a301417b9e','4e34d330-63ad-4413-a058-6777a99ba51b','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'My Boring Battery AA 2800 mAh 4 Pack Baterai AA Rechargeable',NULL,1661,199,3000,TRUE,' ',4.9,125000,128000,'2022-12-22 21:43:56.213942+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('4bfc3bd8-7078-45f6-97f8-604a892c9e3c','ccddcbb2-7a62-4e9d-9ce3-5b2e6e6435fb','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Tali Helm jaring helm Tali Pengikat Helm elastis Helmet Rope Cargo Net',NULL,3713,234,7000,TRUE,' ',4.9,100000,107000,'2022-12-22 21:43:56.216901+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('b1ecdc29-7b82-4bab-b32e-025986dc3c97','0fdbf93b-f71b-4b41-8529-f011c1186f0f','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Karton Box / Kardus Untuk Packing',NULL,187,2832,3000,TRUE,' ',5,16000,19000,'2022-12-22 21:43:56.222758+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('4d4cba75-62af-46b6-bbff-7793f01ad332','a25dad10-4005-48b5-b46f-442f421ef932','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'TWIN PACK - Bagus Anti Bau Mobil 100 gr',NULL,4330,842,10000,TRUE,' ',4.9,31100,41100,'2022-12-22 21:43:56.224544+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('dd3a7739-05ed-401a-a685-aa62347a55f5','b3777c29-de06-4b6f-abf9-9542af906a93','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'BUKET BUNGA PEONY KOMBINASI BUNGA ARTIFICIAL - peach',NULL,151,587,1000,TRUE,' ',4.9,19500,20500,'2022-12-22 21:43:56.225595+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('c686c33b-b895-404d-811f-eb2c184ffbd9','92b839b7-4088-4509-818b-fa8dd740863a','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Back To Black Permanent Plastic & Trim Restorer Nanotech Protection',NULL,5609,9895,9000,TRUE,' ',4.9,85000,94000,'2022-12-22 21:43:56.226634+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('3e01b284-fdbb-443b-aa88-0e870d959e47','48833905-98e8-4057-a076-efbc5585656a','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Mama Lemon Sabun Cuci Piring Jeruk Nipis 680 ml x 2 pcs',NULL,10200,7448,10000,TRUE,' ',5,38000,48000,'2022-12-22 21:43:56.227543+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('e6fb2764-076f-4b3b-bc05-8aa125d537ed','be3c27fc-b20e-48c1-a5b9-f92227548d04','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Masker KN95 PM2.5 Earloop kn 95 filter 95% setara n95 4 ply isi 50 Pcs',NULL,46500,3874,100000,TRUE,' ',4.8,70000,170000,'2022-12-22 21:43:56.228411+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('92d69b54-e7ae-448c-b495-7cf7132f6ac8','4f1a49dd-4df8-4f49-adab-c76b3700e5b7','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Wiper Fluid YOSHIKI Cairan Air Wiper Pembersih Perawatan Kaca Mobil',NULL,2975,7589,6000,TRUE,' ',4.9,50000,56000,'2022-12-22 21:43:56.229258+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('9edbaf2d-a61e-41ff-b5a9-96219b5434c1','31f8a451-4c64-422e-8774-5044ffa123bc','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Tolak Angin Cair Dus 2x12''s Herbal - Masuk Angin Daya Tahan Tubuh',NULL,22500,5937,10000,TRUE,' ',5,118554,128554,'2022-12-22 21:43:56.230254+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('3b98fca4-200e-415a-a8e2-250a25de952a','e0cef9ef-369e-461c-9bcf-af4b51c5e2f8','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Extra Virgin Coconut Oil VCO 500ml - Minyak Kelapa Murni 100%',NULL,8190,1829,10000,TRUE,' ',5,20000,30000,'2022-12-22 21:43:56.231762+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('79977bbd-3acb-4ae9-8b12-efa8fae2e96a','48833905-98e8-4057-a076-efbc5585656a','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Mama Lemon Sabun Cuci Piring Lemon 680 ml x 2 pcs',NULL,11700,34,10000,TRUE,' ',5,38000,48000,'2022-12-22 21:43:56.235619+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('910dbb9b-53d5-4a23-b8ca-cf3f0caab169','e0cef9ef-369e-461c-9bcf-af4b51c5e2f8','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'ROSE BRAND Minyak Goreng 2 liter',NULL,5813,5774,10000,TRUE,' ',5,40000,50000,'2022-12-22 21:43:56.237079+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('137667bf-d00a-423c-b36d-55257bfd60d8','2617472a-1ab3-4cb6-9562-34d2f076dc9b','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Mainan Anak Tek Tek Mainan Lato Lato Jadul Edukasi Kato Kato - random3,8x3,8cm',NULL,36,38382,750,TRUE,' ',4.8,9000,9750,'2022-12-22 21:43:56.238227+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('fb06aa3c-f7b7-43bc-9e27-7f5b1ba69dba','f3d60dfd-cb4c-4727-b509-e0019a34f57d','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Sinocare BA-801 Alat Cek Tekanan Darah Otomatis',NULL,8664,82748,10000,TRUE,' ',4.9,399000,409000,'2022-12-22 21:43:56.239293+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('ebc4d9a3-8c9a-4731-b6d6-6359668fcd9d','3e3a6fa8-d387-44d2-a8f1-4768bf658438','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'buket bunga untuk wisuda,ulangtahun, DLL',NULL,77,85878,100,TRUE,' ',5,170000,170100,'2022-12-22 21:43:56.240314+00',NULL,NULL);
INSERT INTO product(id,category_id,shop_id,sku,title,description,view_count,favorite_count,unit_sold,listed_status,thumbnail_url,rating_avg,min_price,max_price,created_at,updated_at,deleted_at) VALUES ('855f81a1-c7aa-4428-9f26-74fd664de377','7e2a3f49-87fc-4b42-aae0-e59e795beee0','e8854443-c2c7-488e-93d5-b9d93708b8a3',NULL,'Nice Facial Tissue Tisue Tisu Wajah Muka Paket isi 5 x 180 sheets',NULL,3509,12736,10000,TRUE,' ',5,32000,42000,'2022-12-22 21:43:56.241554+00',NULL,NULL);

INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('ca8cf97a-0f13-4ac9-854d-a16b9e9a3d1a','Promo 1','f2535899-134d-47e7-bbb5-43c979102ca6',24,0,30000,50000,1,1,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.381502+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('9257e7ff-c84e-42ea-a6c7-64dfa31f9e72','Promo 2','f25293f2-7c97-4db9-ba23-38f0740ce146',25,0,150000,150000,2,1,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.387007+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('e1c7dd81-84ef-4c68-8509-67c8e49961e8','Promo 3','ffb50fce-e84b-4e68-9316-69a301417b9e',5,0,200000,125000,3,2,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.388819+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('1e64d57e-b53d-4f78-8907-9b72dbd40266','Promo 4','4bfc3bd8-7078-45f6-97f8-604a892c9e3c',0,3000,0,100000,4,2,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.390368+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('073dddc0-1314-4382-8fed-14a632444d24','Promo 5','b1ecdc29-7b82-4bab-b32e-025986dc3c97',3,0,0,16000,5,2,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.392647+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('71c0055f-b02e-4b71-8579-4e365aa94596','Promo 6','4d4cba75-62af-46b6-bbff-7793f01ad332',17,0,15000,30000,6,2,'2023-01-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.394003+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('3482b8f6-0854-4c86-bc93-a6bbd8b1c384','Promo 7','dd3a7739-05ed-401a-a685-aa62347a55f5',0,5000,15000,20000,7,3,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.395084+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('a1d85963-3b15-42a7-a6fb-62513a6b7a5d','Promo 8','c686c33b-b895-404d-811f-eb2c184ffbd9',6,0,15500,85000,8,3,'2023-01-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.396115+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('5e5ae8c5-cb7e-4dca-a9e9-9bc19176d5e8','Promo 9','3e01b284-fdbb-443b-aa88-0e870d959e47',42,0,0,38000,9,3,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.397102+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('7646e2a4-acfa-4804-a771-56589c858240','Promo 10','e6fb2764-076f-4b3b-bc05-8aa125d537ed',69,0,1000,70000,10,3,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.398182+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('a3f44539-0242-41ed-b654-e4030a26d777','Promo 11','92d69b54-e7ae-448c-b495-7cf7132f6ac8',0,5000,0,50000,11,4,'2023-01-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.399365+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('1e1713dc-56b2-43a0-902e-6d50816e6e31','Promo 12','9edbaf2d-a61e-41ff-b5a9-96219b5434c1',17,0,30000,10000,12,4,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.400441+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('86bcfb21-d729-4006-bc82-bb08c32357cc','Promo 13','3b98fca4-200e-415a-a8e2-250a25de952a',0,10000,15000,20000,13,4,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.401705+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('0fb7a375-84d6-47e2-8c60-f6f0cea48c7d','Promo 14','79977bbd-3acb-4ae9-8b12-efa8fae2e96a',42,0,0,38000,14,4,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.403468+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('276f1571-18d9-42fe-a7c7-e9fc7f62ab14','Promo 15','910dbb9b-53d5-4a23-b8ca-cf3f0caab169',21,0,150000,40000,15,5,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.405236+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('70b894eb-4bcd-40c3-a0c6-93e5e369e719','Promo 16','137667bf-d00a-423c-b36d-55257bfd60d8',0,1000,0,9000,16,5,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.407167+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('cfea1bbf-bfbc-4d2e-acbe-ae795c43d30f','Promo 17','fb06aa3c-f7b7-43bc-9e27-7f5b1ba69dba',52,0,150000,150000,17,5,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.409251+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('f61102ff-f2d8-46c0-8cf7-fefe5d31af1a','Promo 18','ebc4d9a3-8c9a-4731-b6d6-6359668fcd9d',0,25000,0,170000,18,5,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.410626+00',NULL,NULL);
INSERT INTO promotion(id,name,product_id,discount_percentage,discount_fix_price,min_product_price,max_discount_price,quota,max_quantity,actived_date,expired_date,created_at,updated_at,deleted_at) VALUES ('6df2615f-71b9-470e-8afc-b33fd3c70112','Promo 19','855f81a1-c7aa-4428-9f26-74fd664de377',6,0,20000,32000,19,6,'2023-01-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00','2022-12-23 02:26:21.411813+00',NULL,NULL);

INSERT INTO voucher(id,shop_id,code,quota,actived_date,expired_date,discount_percentage,discount_fix_price,min_product_price,max_discount_price,created_at,updated_at,deleted_at) VALUES ('59bfcd74-e278-4c70-a889-d9c8515bf71c','e8854443-c2c7-488e-93d5-b9d93708b8a3','ASD123',11,'2022-12-22 21:43:56.202214+00','2023-02-01 21:43:56.202214+00',0,5000,5000,100000,'2022-12-23 02:34:38.854025+00',NULL,NULL);

INSERT INTO "product_detail" (id, product_id, price, stock, weight, size, hazardous, condition, bulk_price)
VALUES
('0c53ef3d-3682-4359-90e1-814eb6ab5231', 'e6fb2764-076f-4b3b-bc05-8aa125d537ed', 70000, 111, 0.5, 0.1, false, 'perfect', false),
('0c53ef3d-3682-4359-90e1-814eb6ab5111', '855f81a1-c7aa-4428-9f26-74fd664de377', 100000, 100, 1, 1, false, 'jelek', false);

INSERT INTO "cart_item" (id, user_id, product_detail_id, quantity)
VALUES
('410a1545-1834-4dca-9624-8c1c7e1439de' ,'7950eca2-58d5-44f0-b873-22b23d8107da', '0c53ef3d-3682-4359-90e1-814eb6ab5231', 4),
('d37f3e57-94d9-433c-a3c8-316f6b719418' ,'7950eca2-58d5-44f0-b873-22b23d8107da', '0c53ef3d-3682-4359-90e1-814eb6ab5111', 2);

INSERT INTO "variant_detail" (id, name, type)
VALUES
('fe997ee9-6bb1-4ddc-bff9-4a2bc0d2ad24' ,'ukuran' ,'l'),
('4e18d196-f84a-4a95-a897-fe085d38347f' ,'warna' ,'hijau');

INSERT INTO "variant" (id, product_detail_id, variant_detail_id)
VALUES
('614fa874-d6a7-414a-9ebc-9b43cf765745', '0c53ef3d-3682-4359-90e1-814eb6ab5231', 'fe997ee9-6bb1-4ddc-bff9-4a2bc0d2ad24'),
('ff4b9968-39ba-48f3-ad77-088a118b7c4a', '0c53ef3d-3682-4359-90e1-814eb6ab5111', '4e18d196-f84a-4a95-a897-fe085d38347f');
