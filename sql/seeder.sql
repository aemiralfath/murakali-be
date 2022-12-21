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

INSERT INTO "category" (name, photo_url)
VALUES ('laptop', 'https://cf.shopee.co.id/file/c139370836a9daa649da70876a326b58'),
       ('electronik', 'https://cf.shopee.co.id/file/dcd61dcb7c1448a132f49f938b0cb553_tn'),
       ('pakaian pria', 'https://cf.shopee.co.id/file/04dba508f1ad19629518defb94999ef9_tn'),
       ('hobi & koleksi', 'https://cf.shopee.co.id/file/42394b78fac1169d67c6291973a3b132_tn');
