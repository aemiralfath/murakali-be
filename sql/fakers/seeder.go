package fakers

import (
	"database/sql"
	"murakali/pkg/postgre"
	"murakali/sql/fakers/table"
)

type Seeder struct {
	Seeder table.ISeeder
}

func RegisterSeeders() []Seeder {
	return []Seeder{
		{Seeder: table.NewRoleFaker([]string{"user", "seller", "admin"})},
		{Seeder: table.NewUserFaker(
			0,
			1,
			"M",
			[]string{"4c1d6464-3cc6-44d6-92d1-91aee337e025", "0c53ef3d-3682-4359-90e1-814eb6ab5191", "7950eca2-58d5-44f0-b873-22b23d8107da"},
			[]string{"fadhlan1337@gmail.com", "sammymanunggal@gmail.com", "user@gmail.com"})},
		{Seeder: table.NewUserFaker(
			0,
			2,
			"M",
			[]string{"f8d8d66a-e8eb-4633-bc2d-4ccd941fed47", "87cb732a-4e09-461b-a6cc-be818353cae7", "942c718e-0ea9-4b8e-bbb8-ad5138ba9f6f", "9e3a2d9d-0479-4819-a40b-068453bcaf04", "fe52cf85-608a-4d83-b086-83da62dcccc5", "33998525-b174-4cd6-bb57-9778da7fe45b", "88012133-581a-4c0b-9048-7d61d969514f", "61daa36f-2bae-41ca-91fe-d4b6f92f6e76",
				"8bfb0a42-5a58-41d5-8c6e-bb39a187252d", "c11f2512-d882-4b3e-b0cb-1865d198c954", "4616cc5a-e76d-4c5f-bb64-45e7ddae5807", "12d1b664-4b1a-4968-bf33-ee07aa64bec3"},
			[]string{"seller1@gmail.com", "seller2@gmail.com", "seller3@gmail.com", "seller4@gmail.com", "seller5@gmail.com", "seller6@gmail.com", "seller7@gmail.com", "seller8@gmail.com",
				"seller9@gmail.com", "seller10@gmail.com", "seller11@gmail.com", "seller12@gmail.com"})},
		{Seeder: table.NewUserFaker(0, 3, "M", []string{"4df967a8-5b05-4d2a-bb72-da3921dce8fb"}, []string{"admin@gmail.com"})},
		{Seeder: table.NewCategoryFaker(
			0,
			[]string{"d92a0995-78cd-4eba-a855-dfc096ffec5b", "5d5bd121-adc2-4f62-9cad-d4172bec9a40", "5778e73c-f8b7-4c6b-a2f4-472079b164c5", "63f58102-9cb6-4249-b8d4-82f65f315c59",
				"f2a5281e-e9d1-4fd5-bff7-2afd995d5a59", "14a4a0d0-dc24-4ef3-ad18-5de3f19bb352", "66c5c7a8-729c-4d7a-b4b4-0de7a4b334ca", "80d9efde-1246-41f9-b768-743bf2949763"},
			[]string{"laptop", "elektronik", "pakaian pria", "hobi & koleksi", "outfit Hangat", "sweeter", "mouse & keyboard", "webcam"},
			[]string{"https://cf.shopee.co.id/file/c139370836a9daa649da70876a326b58", "https://cf.shopee.co.id/file/dcd61dcb7c1448a132f49f938b0cb553_tn", "https://cf.shopee.co.id/file/04dba508f1ad19629518defb94999ef9_tn", "https://cf.shopee.co.id/file/42394b78fac1169d67c6291973a3b132_tn",
				"https://cf.shopee.co.id/file/d89df04fd3435962af59be0408ec4786", "https://cf.shopee.co.id/file/19b8238c917f3dec99b689809ea43a79_tn", "https://cf.shopee.co.id/file/6e70e5f93250a3e8344cda4fc79b0c3f", "https://cf.shopee.co.id/file/45ee92cbf6243007a66f0f338058da80"},
			[]string{"", "", "", "", "", "", "", ""})},
		{Seeder: table.NewCategoryFaker(
			0,
			[]string{"d99373d1-c55d-4769-a56e-f797db20235d", "159aa7d7-2fa0-4cc8-a708-3328d1d08eb5", "0774dbda-194f-439d-97e3-eec0e325fe5a", "1aaaed1f-9d23-47ef-8647-17b862becc27"},
			[]string{"outfit Hangat", "sweeter", "mouse & keyboard", "webcam"},
			[]string{"https://cf.shopee.co.id/file/d89df04fd3435962af59be0408ec4786", "https://cf.shopee.co.id/file/19b8238c917f3dec99b689809ea43a79_tn", "https://cf.shopee.co.id/file/6e70e5f93250a3e8344cda4fc79b0c3f", "https://cf.shopee.co.id/file/45ee92cbf6243007a66f0f338058da80"},
			[]string{"5778e73c-f8b7-4c6b-a2f4-472079b164c5", "d99373d1-c55d-4769-a56e-f797db20235d", "5d5bd121-adc2-4f62-9cad-d4172bec9a40", "5d5bd121-adc2-4f62-9cad-d4172bec9a40"})},
		{Seeder: table.NewCourierFaker(
			[]string{"98c1921e-b80e-40f3-9cba-fe8806097517", "0d389020-f229-461e-9202-5788961fbb81", "4bf503dc-689e-4b66-8401-3f133f1d585a"},
			[]string{"JNE", "POS Indonesia", "TIKI"},
			[]string{"jne", "pos", "tiki"},
			[]string{"REG", "Pos Reguler", "REG"})},
		{Seeder: table.NewShopFaker(
			[]string{"e8854443-c2c7-488e-93d5-b9d93708b8a3", "07315003-5369-465f-9f05-09482d951645", "b61ef5a7-548c-4c81-a192-eadeb2af915f", "20d1015e-d03a-4326-bc23-427a861bbc4e", "ecd86fa9-c2a0-4adb-93e8-347b9fac3b56", "a050cfb3-957c-4b35-83cb-ff65095c6eb5", "1a21363c-bc64-4295-8ad2-cb5d6517c797", "735e9978-97cc-4427-9c24-2f8230429a7f",
				"cfd82da1-191e-40d8-a35e-725f9b1c8fb6", "41b415b2-56c1-4fdb-bc37-6522bd66840f", "14f695be-dcc7-4067-a410-82c50e7b2e27", "2f613376-09a7-4be3-8017-ac6e57a6e5ca"},
			[]string{"f8d8d66a-e8eb-4633-bc2d-4ccd941fed47", "87cb732a-4e09-461b-a6cc-be818353cae7", "942c718e-0ea9-4b8e-bbb8-ad5138ba9f6f", "9e3a2d9d-0479-4819-a40b-068453bcaf04", "fe52cf85-608a-4d83-b086-83da62dcccc5", "33998525-b174-4cd6-bb57-9778da7fe45b", "88012133-581a-4c0b-9048-7d61d969514f", "61daa36f-2bae-41ca-91fe-d4b6f92f6e76",
				"8bfb0a42-5a58-41d5-8c6e-bb39a187252d", "c11f2512-d882-4b3e-b0cb-1865d198c954", "4616cc5a-e76d-4c5f-bb64-45e7ddae5807", "12d1b664-4b1a-4968-bf33-ee07aa64bec3"},
			[]string{"d92a0995-78cd-4eba-a855-dfc096ffec5b", "5d5bd121-adc2-4f62-9cad-d4172bec9a40", "5778e73c-f8b7-4c6b-a2f4-472079b164c5", "63f58102-9cb6-4249-b8d4-82f65f315c59", "d99373d1-c55d-4769-a56e-f797db20235d", "159aa7d7-2fa0-4cc8-a708-3328d1d08eb5", "0774dbda-194f-439d-97e3-eec0e325fe5a", "1aaaed1f-9d23-47ef-8647-17b862becc27",
				"f2a5281e-e9d1-4fd5-bff7-2afd995d5a59", "14a4a0d0-dc24-4ef3-ad18-5de3f19bb352", "66c5c7a8-729c-4d7a-b4b4-0de7a4b334ca", "80d9efde-1246-41f9-b768-743bf2949763"},
			[]string{"Laptop Shop", "Electronic Shop", "Man Fashion Shop", "Hobby & Collection Shop", "Warm Outfit Shop", "Sweeter Shop", "Mouse & Keyboard Shop", "Webcam Shop", "Warm Outfit Shop++", "Sweeter Shop++", "Mouse & Keyboard Shop++", "Webcam Shop++"},
			[]int{8500, 8500, 8500, 85100, 8500, 8500, 8500, 8500, 8500, 8500, 8500, 8500},
			[]string{"98c1921e-b80e-40f3-9cba-fe8806097517", "0d389020-f229-461e-9202-5788961fbb81", "4bf503dc-689e-4b66-8401-3f133f1d585a"})},
		{Seeder: table.NewUserFaker(10000, 1, "M", []string{}, []string{})},
		{Seeder: table.NewUserFaker(10000, 1, "F", []string{}, []string{})},
	}
}

func DBSeed(sqlDB *sql.DB) error {
	txDB := postgre.NewTxRepository(sqlDB)
	for _, seeder := range RegisterSeeders() {
		err := txDB.WithTransaction(func(transaction postgre.Transaction) error {
			if errSeeder := seeder.Seeder.GenerateData(transaction); errSeeder != nil {
				return errSeeder
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
