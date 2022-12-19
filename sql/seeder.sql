INSERT INTO "role" (name)
VALUES ('user'),
       ('seller'),
       ('admin');

INSERT INTO "email_history" (email) VALUES ('test@gmail.com');

INSERT INTO "user" (id, role_id, username, email, phone_no, fullname, password, is_verify)
VALUES ('7950eca2-58d5-44f0-b873-22b23d8107da', 1, 'test', 'test@gmail.com', '911', 'test', '$2a$10$cNhdZVN.pgsfK1xUQ00p7eK5Fh7iClrtJB9SY5un.H55Mi/dtQzCa', true);

INSERT INTO "address" (user_id, name, province_id, city_id, province, city, district, sub_district, address_detail, zip_code)
VALUES ('7950eca2-58d5-44f0-b873-22b23d8107da', 'test', 5, 39, 'DI Yogyakarta', 'Bantul', 'Pleret', 'Segoroyoso', 'no 91', '55791');
