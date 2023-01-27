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
		{Seeder: table.NewOrderStatusFaker([]string{"Waiting to Pay", "Waiting for Seller", "Processed", "On Delivery", "Delivered", "Received", "Completed", "Canceled", "Refunded"})},
		{Seeder: table.NewUserFaker(
			0,
			1,
			"M",
			[]string{"4c1d6464-3cc6-44d6-92d1-91aee337e025", "0c53ef3d-3682-4359-90e1-814eb6ab5191", "7950eca2-58d5-44f0-b873-22b23d8107da"},
			[]string{"fadhlan1337@gmail.com", "sammymanunggal@gmail.com", "user@gmail.com"},
			[]string{"1234567890123456", "2789760285732876", "2787884621261326"})},
		{Seeder: table.NewUserFaker(
			0,
			2,
			"M",
			[]string{"f8d8d66a-e8eb-4633-bc2d-4ccd941fed47", "87cb732a-4e09-461b-a6cc-be818353cae7", "942c718e-0ea9-4b8e-bbb8-ad5138ba9f6f", "9e3a2d9d-0479-4819-a40b-068453bcaf04", "fe52cf85-608a-4d83-b086-83da62dcccc5",
				"33998525-b174-4cd6-bb57-9778da7fe45b", "88012133-581a-4c0b-9048-7d61d969514f", "61daa36f-2bae-41ca-91fe-d4b6f92f6e76", "8bfb0a42-5a58-41d5-8c6e-bb39a187252d", "c11f2512-d882-4b3e-b0cb-1865d198c954",
				"4616cc5a-e76d-4c5f-bb64-45e7ddae5807", "12d1b664-4b1a-4968-bf33-ee07aa64bec3", "d9ba7b98-8c74-4af8-83bf-8fed2b48e9e6", "ace773c7-cb83-41bf-89a5-c8d8ab796ef6", "8f072f10-429f-482d-acfa-28bbc9ab8102",
				"e30a43f1-6229-4aed-a559-cf251c03fec9", "0d0f4d0f-b237-4082-b268-42887bb6ab79", "a75590bc-c2d0-4f9a-b14c-d316752a1684", "923af192-6d6c-4263-8551-cc1e848b177f", "2c3c960b-3599-40b4-b072-ad371a760f4e",
				"0279b09e-8ec7-49ae-8891-251c206315e9", "5512ba6a-5c9b-499a-93cb-cab64747df36", "3b2edaef-8d30-48af-9959-d4301ce24591", "20b24886-0429-4d74-a2ee-24167298f839", "8bb271be-648a-4356-be4b-050825590f6e",
				"ef659bc0-0c81-4f58-a28e-0e65f617f939", "bb43fe25-e268-469f-9e1f-6fcbbbc9a8c6", "f5d60e80-0083-4ba7-9c3f-d5790952466d", "a6fc1880-c93e-4e59-bf62-0418f0dc328e", "a6452c10-6d5e-4cf7-af0c-0e23417a2839",
				"f7315779-8c90-4aba-9764-e83a8ae3cbd3", "e5afa898-f5c2-4054-84c8-17961a63a0d6", "da68ab69-4066-42d4-8179-dfe45b62b9aa", "cf781b92-2324-4de5-bf0b-1c7ae3fba5f6", "2575471b-78ca-40f2-a3be-d4c2f6a4d66d",
				"a01b4451-c746-42b8-b06d-faf15b40e169", "dceebe0e-a38f-4645-b5f6-822cd6ddfaf6", "cdc8ad65-c382-47d7-8691-b4896d864a8e", "c07ce5df-7112-4f47-b276-38b0f9e93a9c", "1f5505c3-9197-4eea-a000-078f54350353",
				"66e94d1f-2e7f-400b-a66d-07db6eaccbf4", "719e2f28-624b-4f35-9eaf-cfea782cccaf", "01b3797f-3b33-4675-839c-7383da26d78b", "f61b4762-9bc5-4d27-b318-6fab625363f0", "bfbbcbec-e647-4ec4-9817-7644a87a556d"},
			[]string{"seller1@gmail.com", "seller2@gmail.com", "seller3@gmail.com", "seller4@gmail.com", "seller5@gmail.com",
				"seller6@gmail.com", "seller7@gmail.com", "seller8@gmail.com", "seller9@gmail.com", "seller10@gmail.com",
				"seller11@gmail.com", "seller12@gmail.com", "seller13@gmail.com", "seller14@gmail.com", "seller15@gmail.com",
				"seller16@gmail.com", "seller17@gmail.com", "seller18@gmail.com", "seller19@gmail.com", "seller20@gmail.com",
				"seller21@gmail.com", "seller22@gmail.com", "seller23@gmail.com", "seller24@gmail.com", "seller25@gmail.com",
				"seller26@gmail.com", "seller27@gmail.com", "seller28@gmail.com", "seller29@gmail.com", "seller30@gmail.com",
				"seller31@gmail.com", "seller32@gmail.com", "seller33@gmail.com", "seller34@gmail.com", "seller35@gmail.com",
				"seller36@gmail.com", "seller37@gmail.com", "seller38@gmail.com", "seller39@gmail.com", "seller40@gmail.com",
				"seller41@gmail.com", "seller42@gmail.com", "seller43@gmail.com", "seller44@gmail.com", "seller45@gmail.com"},
			[]string{"1234567890123453", "1234567890123454", "1234567890123455", "1234567890123457", "1234567890123458",
				"1234567890123459", "1234567890123460", "1234567890123461", "1234567890123462", "1234567890123463",
				"1234567890123464", "1234567890123465", "1234567890123466", "1234567890123467", "1234567890123468",
				"1234567890123469", "1234567890123470", "1234567890123471", "1234567890123472", "1234567890123473",
				"1234567890123474", "1234567890123475", "1234567890123476", "1234567890123477", "1234567890123478",
				"1234567890123479", "1234567890123480", "1234567890123481", "1234567890123482", "1234567890123483",
				"1234567890123484", "1234567890123485", "1234567890123486", "1234567890123487", "1234567890123488",
				"1234567890123489", "1234567890123490", "1234567890123491", "1234567890123492", "1234567890123493",
				"1234567890123494", "1234567890123495", "1234567890123496", "1234567890123497", "1234567890123498"})},
		{Seeder: table.NewUserFaker(0, 3, "M", []string{"4df967a8-5b05-4d2a-bb72-da3921dce8fb"}, []string{"admin@gmail.com"}, []string{"12345678901234616"})},
		{Seeder: table.NewCategoryFaker(
			0,
			[]string{"d92a0995-78cd-4eba-a855-dfc096ffec5b", "5d5bd121-adc2-4f62-9cad-d4172bec9a40", "5778e73c-f8b7-4c6b-a2f4-472079b164c5",
				"63f58102-9cb6-4249-b8d4-82f65f315c59", "f2a5281e-e9d1-4fd5-bff7-2afd995d5a59", "14a4a0d0-dc24-4ef3-ad18-5de3f19bb352",
				"66c5c7a8-729c-4d7a-b4b4-0de7a4b334ca", "80d9efde-1246-41f9-b768-743bf2949763", "2f575735-5232-4208-bb9c-bfcf091cae2d",
				"9b32fe3e-adfa-4bc7-82fa-6737080d44cd", "fa3bdd1d-b7d1-4cef-b737-be86d192162d", "49b298b7-aefd-452a-a08a-5181be8d3e1b"},
			[]string{"Elektronik", "Pakaian", "Sepatu",
				"Tas", "Aksesoris Fashion", "Hobi & Koleksi",
				"Kesehatan", "Makanan & Minuman", "Perawatan & Kecantikan",
				"Perlengkapan Dapur", "Otomotif", "Olahraga & Outdoor"},
			[]string{"https://cf.shopee.co.id/file/dcd61dcb7c1448a132f49f938b0cb553_tn", "https://cf.shopee.co.id/file/04dba508f1ad19629518defb94999ef9_tn", "https://cf.shopee.co.id/file/3c8ff51aab1692a80c5883972a679168_tn",
				"https://cf.shopee.co.id/file/47ed832eed0feb62fd28f08c9229440e_tn", "https://cf.shopee.co.id/file/1f18bdfe73df39c66e7326b0a3e08e87_tn", "https://cf.shopee.co.id/file/42394b78fac1169d67c6291973a3b132_tn",
				"https://cf.shopee.co.id/file/eb7d583e4b72085e71cd21a70ce47d7a_tn", "https://cf.shopee.co.id/file/7873b8c3824367239efb02d18eeab4f5_tn", "https://cf.shopee.co.id/file/2715b985ae706a4c39a486f83da93c4b_tn",
				"https://cf.shopee.co.id/file/c1494110e0383780cdea73ed890e0299_tn", "https://cf.shopee.co.id/file/27838b968afb76ca59dd8e8f57ece91f_tn", "https://cf.shopee.co.id/file/b2c24b49fd96704ed80b4f45080bfcac_tn"},
			[]string{"", "", "", "", "", "", "", "", "", "", "", ""})},
		{Seeder: table.NewCategoryFaker(
			0,
			[]string{"d99373d1-c55d-4769-a56e-f797db20235d", "159aa7d7-2fa0-4cc8-a708-3328d1d08eb5", "0774dbda-194f-439d-97e3-eec0e325fe5a",
				"1aaaed1f-9d23-47ef-8647-17b862becc27", "272085be-4887-498a-b7f6-85870fe93b40", "fb562584-ff7d-470b-a85a-0ee420a25850",
				"6ab15b0d-adf8-4c58-815a-194ba9c67e0d", "a7e9159f-af65-44a1-abe6-a86767a2f8f6", "9c5b448e-35c3-443a-bf13-65d9ac86fb19",
				"4a5d0bc1-5a43-4d0d-8b6a-4bf5e755846e", "f77cc578-bf01-4f79-8a6c-885e42d5ed37", "06a722ae-f0a1-4d94-80a0-6c2bc2b2597d",
				"de528431-8391-4c85-b89e-6fc6a77babba", "62648619-db60-42bc-bb98-5cf4bf2e06b4", "b09cab84-5327-415a-a076-3891a51aa211",
				"1a0b3314-2b2b-4fed-8a4a-55ce74a4ef19", "26bc0e0e-0b6c-487b-bc1e-823840cbaa52", "8ad11626-73ad-412c-8213-85ccbc5e180e",
				"fcf85764-8972-4e0a-9c27-1c3a3d4b7ae2", "f968b6cc-9e34-4891-b87b-7f78cc0f3aa5", "ed274237-32e2-4c3c-8606-7b44885e1ac1",
				"02cce80b-ecad-4cce-95f2-c5e3dd887390", "66a50383-baa3-4db9-871e-7ec151ac910f", "0a8476a5-b8bb-4a98-9a0a-64bc794e0c35",
				"458126d1-9230-48ba-992c-7f6287c35b26", "592186bd-6f25-4426-8ee1-fc78a7d98f56", "7f4a631f-f59d-4453-8c16-f023d050bbb9"},
			[]string{"Komputer", "Handphone", "Kamera",
				"Pakaian Pria", "Pakaian Wanita", "Outfit Hangat",
				"Sepatu Pria", "Sepatu Wanita", "Sepatu Anak",
				"Tas Pria", "Tas Wanita", "Kacamata",
				"Jam Tangan", "Gitar", "Buku",
				"Obat", "Masker", "Makanan",
				"Minuman", "Face Wash", "Sunscreen",
				"Panci", "Kompor", "Helm",
				"Spion", "Raket", "Bola"},
			[]string{"https://cf.shopee.co.id/file/290afeb96794a9870bda2a69562e2637_tn", "https://cf.shopee.co.id/file/5230277eefafad8611aaf703d3e99568_tn", "https://cf.shopee.co.id/file/9abe95c0c755968c5114f084ee11b8cb_tn",
				"https://cf.shopee.co.id/file/04dba508f1ad19629518defb94999ef9_tn", "https://cf.shopee.co.id/file/6d63cca7351ba54a2e21c6be1721fa3a_tn", "https://cf.shopee.co.id/file/d89df04fd3435962af59be0408ec4786",
				"https://cf.shopee.co.id/file/3c8ff51aab1692a80c5883972a679168_tn", "https://cf.shopee.co.id/file/0c2c105fb4317947269e4061c8694fb8_tn", "https://cf.shopee.co.id/file/0c2c105fb4317947269e4061c8694fb8_tn",
				"https://cf.shopee.co.id/file/47ed832eed0feb62fd28f08c9229440e_tn", "https://cf.shopee.co.id/file/522f86a1c4fa5a02c5996dc72ddd73b5_tn", "https://cf.shopee.co.id/file/1f18bdfe73df39c66e7326b0a3e08e87_tn",
				"https://cf.shopee.co.id/file/2bdf8cf99543342d4ebd8e1bdb576f80_tn", "https://cf.shopee.co.id/file/42394b78fac1169d67c6291973a3b132_tn", "https://cf.shopee.co.id/file/998c7682fd5e7a3563b2ad00aaa4e6f3_tn",
				"https://cf.shopee.co.id/file/eb7d583e4b72085e71cd21a70ce47d7a_tn", "https://cf.shopee.co.id/file/eb7d583e4b72085e71cd21a70ce47d7a_tn", "https://cf.shopee.co.id/file/7873b8c3824367239efb02d18eeab4f5_tn",
				"https://cf.shopee.co.id/file/7873b8c3824367239efb02d18eeab4f5_tn", "https://cf.shopee.co.id/file/2715b985ae706a4c39a486f83da93c4b_tn", "https://cf.shopee.co.id/file/2715b985ae706a4c39a486f83da93c4b_tn",
				"https://cf.shopee.co.id/file/c1494110e0383780cdea73ed890e0299_tn", "https://cf.shopee.co.id/file/c1494110e0383780cdea73ed890e0299_tn", "https://cf.shopee.co.id/file/27838b968afb76ca59dd8e8f57ece91f_tn",
				"https://cf.shopee.co.id/file/27838b968afb76ca59dd8e8f57ece91f_tn", "https://cf.shopee.co.id/file/b2c24b49fd96704ed80b4f45080bfcac_tn", "https://cf.shopee.co.id/file/b2c24b49fd96704ed80b4f45080bfcac_tn"},
			[]string{"d92a0995-78cd-4eba-a855-dfc096ffec5b", "d92a0995-78cd-4eba-a855-dfc096ffec5b", "d92a0995-78cd-4eba-a855-dfc096ffec5b",
				"5d5bd121-adc2-4f62-9cad-d4172bec9a40", "5d5bd121-adc2-4f62-9cad-d4172bec9a40", "5d5bd121-adc2-4f62-9cad-d4172bec9a40",
				"5778e73c-f8b7-4c6b-a2f4-472079b164c5", "5778e73c-f8b7-4c6b-a2f4-472079b164c5", "5778e73c-f8b7-4c6b-a2f4-472079b164c5",
				"63f58102-9cb6-4249-b8d4-82f65f315c59", "63f58102-9cb6-4249-b8d4-82f65f315c59", "f2a5281e-e9d1-4fd5-bff7-2afd995d5a59",
				"f2a5281e-e9d1-4fd5-bff7-2afd995d5a59", "14a4a0d0-dc24-4ef3-ad18-5de3f19bb352", "14a4a0d0-dc24-4ef3-ad18-5de3f19bb352",
				"66c5c7a8-729c-4d7a-b4b4-0de7a4b334ca", "66c5c7a8-729c-4d7a-b4b4-0de7a4b334ca", "80d9efde-1246-41f9-b768-743bf2949763",
				"80d9efde-1246-41f9-b768-743bf2949763", "2f575735-5232-4208-bb9c-bfcf091cae2d", "2f575735-5232-4208-bb9c-bfcf091cae2d",
				"9b32fe3e-adfa-4bc7-82fa-6737080d44cd", "9b32fe3e-adfa-4bc7-82fa-6737080d44cd", "fa3bdd1d-b7d1-4cef-b737-be86d192162d",
				"fa3bdd1d-b7d1-4cef-b737-be86d192162d", "49b298b7-aefd-452a-a08a-5181be8d3e1b", "49b298b7-aefd-452a-a08a-5181be8d3e1b"})},
		{Seeder: table.NewCategoryFaker(
			0,
			[]string{"a81c33b8-b429-4879-bbe1-1adc65987a57", "59c299d8-faae-44d7-b751-424fb3077072", "055e72b4-c0fb-4a19-b945-baa499daf3e6",
				"53ba58ec-6917-4ae6-b91f-6a3eadcaf0a6", "3d005413-587a-42ca-be1a-97106d684861", "29b42c97-7f1e-4933-8286-18592a9845b6"},
			[]string{"Laptop", "Mouse", "Keyboard",
				"Webcam", "Sweeter", "Jaket"},
			[]string{"https://cf.shopee.co.id/file/c139370836a9daa649da70876a326b58", "https://cf.shopee.co.id/file/6e70e5f93250a3e8344cda4fc79b0c3f", "https://cf.shopee.co.id/file/6e70e5f93250a3e8344cda4fc79b0c3f",
				"https://cf.shopee.co.id/file/45ee92cbf6243007a66f0f338058da80", "https://cf.shopee.co.id/file/19b8238c917f3dec99b689809ea43a79_tn", "https://cf.shopee.co.id/file/19b8238c917f3dec99b689809ea43a79_tn"},
			[]string{"d99373d1-c55d-4769-a56e-f797db20235d", "d99373d1-c55d-4769-a56e-f797db20235d", "d99373d1-c55d-4769-a56e-f797db20235d",
				"d99373d1-c55d-4769-a56e-f797db20235d", "fb562584-ff7d-470b-a85a-0ee420a25850", "fb562584-ff7d-470b-a85a-0ee420a25850"})},
		{Seeder: table.NewCourierFaker(
			[]string{"98c1921e-b80e-40f3-9cba-fe8806097517", "0d389020-f229-461e-9202-5788961fbb81", "4bf503dc-689e-4b66-8401-3f133f1d585a"},
			[]string{"JNE", "POS Indonesia", "TIKI"},
			[]string{"jne", "pos", "tiki"},
			[]string{"REG", "Pos Reguler", "REG"})},
		{Seeder: table.NewShopFaker(
			[]string{"e8854443-c2c7-488e-93d5-b9d93708b8a3", "07315003-5369-465f-9f05-09482d951645", "b61ef5a7-548c-4c81-a192-eadeb2af915f",
				"20d1015e-d03a-4326-bc23-427a861bbc4e", "ecd86fa9-c2a0-4adb-93e8-347b9fac3b56", "a050cfb3-957c-4b35-83cb-ff65095c6eb5",
				"1a21363c-bc64-4295-8ad2-cb5d6517c797", "735e9978-97cc-4427-9c24-2f8230429a7f", "cfd82da1-191e-40d8-a35e-725f9b1c8fb6",
				"41b415b2-56c1-4fdb-bc37-6522bd66840f", "14f695be-dcc7-4067-a410-82c50e7b2e27", "2f613376-09a7-4be3-8017-ac6e57a6e5ca",
				"f60b0a3c-22db-4842-9386-1ca5cc697d17", "ebe1fea2-bd5f-4c96-8fda-1b5299500307", "5612cf55-7a49-4ec8-8a30-179e2d8d4892",
				"a1dac882-e480-4960-b3fe-c5186843b867", "698145bc-c599-4ff4-952e-d1152720e382", "b6341f43-34d9-4ddb-967f-857fa56c3b6e",
				"a9c1386e-0a32-4733-adad-4e3ceb1dc4a8", "0e9bea36-aa1f-4898-ba80-0b1f63f1e68d", "89b8ec89-0f75-482f-995f-df48d70a9594",
				"be5b8d4c-e099-47d2-b412-298a2a778077", "98a0de8e-cc4a-4b6e-b58f-eb7a9afce71b", "3ff0ead2-855e-4124-884e-3bf2e5d16e3a",
				"f01315af-7b25-4ba3-8856-b757da8ebbdd", "f9bf8cea-d71c-4249-91dc-946b357946f3", "dd3904ce-12fe-4d23-95a3-833a6bfbe2c0",
				"97eec932-8e88-4cb2-8d2f-1ca9978c7ec9", "0d8b8893-357b-413a-bdb3-4348ea8291b1", "5e31142a-f1ba-4bc5-be42-f88b56c64e31",
				"2da16953-eef8-4ed4-a58d-9b5c3fcb7595", "a19337f9-0d62-490b-ab55-f18e9f05b194", "a145a28d-9b18-44f2-a708-75b374e32208",
				"36e4116b-f363-4bc1-829d-abd363e304da", "610c4e8a-141e-45d9-963a-f00972448942", "4b0feacd-6654-4f83-83b0-b8b6b5314713",
				"0cffde48-3d65-47b1-8c65-f1d9db33e08c", "d2fd9f2b-e4f1-4ac2-a0b8-57e189d73a33", "c1a4fcb7-00a5-4b55-afe6-69c80ba87e5c",
				"6a3718f6-434e-4cab-b014-829dc1ee9fb2", "f59fc2fa-5597-431f-bd8b-c61c900e4459", "48814df2-707e-4021-b386-7873be593fb1",
				"3c7f453e-3623-43a6-ab86-f460e20df33c", "50123b17-121e-46e6-a539-359fca00b1dc", "46346ddd-2ed6-4d50-98c6-d7de3353172d"},
			[]string{"f8d8d66a-e8eb-4633-bc2d-4ccd941fed47", "87cb732a-4e09-461b-a6cc-be818353cae7", "942c718e-0ea9-4b8e-bbb8-ad5138ba9f6f",
				"9e3a2d9d-0479-4819-a40b-068453bcaf04", "fe52cf85-608a-4d83-b086-83da62dcccc5", "33998525-b174-4cd6-bb57-9778da7fe45b",
				"88012133-581a-4c0b-9048-7d61d969514f", "61daa36f-2bae-41ca-91fe-d4b6f92f6e76", "8bfb0a42-5a58-41d5-8c6e-bb39a187252d",
				"c11f2512-d882-4b3e-b0cb-1865d198c954", "4616cc5a-e76d-4c5f-bb64-45e7ddae5807", "12d1b664-4b1a-4968-bf33-ee07aa64bec3",
				"d9ba7b98-8c74-4af8-83bf-8fed2b48e9e6", "ace773c7-cb83-41bf-89a5-c8d8ab796ef6", "8f072f10-429f-482d-acfa-28bbc9ab8102",
				"e30a43f1-6229-4aed-a559-cf251c03fec9", "0d0f4d0f-b237-4082-b268-42887bb6ab79", "a75590bc-c2d0-4f9a-b14c-d316752a1684",
				"923af192-6d6c-4263-8551-cc1e848b177f", "2c3c960b-3599-40b4-b072-ad371a760f4e", "0279b09e-8ec7-49ae-8891-251c206315e9",
				"5512ba6a-5c9b-499a-93cb-cab64747df36", "3b2edaef-8d30-48af-9959-d4301ce24591", "20b24886-0429-4d74-a2ee-24167298f839",
				"8bb271be-648a-4356-be4b-050825590f6e", "ef659bc0-0c81-4f58-a28e-0e65f617f939", "bb43fe25-e268-469f-9e1f-6fcbbbc9a8c6",
				"f5d60e80-0083-4ba7-9c3f-d5790952466d", "a6fc1880-c93e-4e59-bf62-0418f0dc328e", "a6452c10-6d5e-4cf7-af0c-0e23417a2839",
				"f7315779-8c90-4aba-9764-e83a8ae3cbd3", "e5afa898-f5c2-4054-84c8-17961a63a0d6", "da68ab69-4066-42d4-8179-dfe45b62b9aa",
				"cf781b92-2324-4de5-bf0b-1c7ae3fba5f6", "2575471b-78ca-40f2-a3be-d4c2f6a4d66d", "a01b4451-c746-42b8-b06d-faf15b40e169",
				"dceebe0e-a38f-4645-b5f6-822cd6ddfaf6", "cdc8ad65-c382-47d7-8691-b4896d864a8e", "c07ce5df-7112-4f47-b276-38b0f9e93a9c",
				"1f5505c3-9197-4eea-a000-078f54350353", "66e94d1f-2e7f-400b-a66d-07db6eaccbf4", "719e2f28-624b-4f35-9eaf-cfea782cccaf",
				"01b3797f-3b33-4675-839c-7383da26d78b", "f61b4762-9bc5-4d27-b318-6fab625363f0", "bfbbcbec-e647-4ec4-9817-7644a87a556d"},
			[]string{"d92a0995-78cd-4eba-a855-dfc096ffec5b", "5d5bd121-adc2-4f62-9cad-d4172bec9a40", "5778e73c-f8b7-4c6b-a2f4-472079b164c5",
				"63f58102-9cb6-4249-b8d4-82f65f315c59", "f2a5281e-e9d1-4fd5-bff7-2afd995d5a59", "14a4a0d0-dc24-4ef3-ad18-5de3f19bb352",
				"66c5c7a8-729c-4d7a-b4b4-0de7a4b334ca", "80d9efde-1246-41f9-b768-743bf2949763", "2f575735-5232-4208-bb9c-bfcf091cae2d",
				"9b32fe3e-adfa-4bc7-82fa-6737080d44cd", "fa3bdd1d-b7d1-4cef-b737-be86d192162d", "49b298b7-aefd-452a-a08a-5181be8d3e1b",
				"d99373d1-c55d-4769-a56e-f797db20235d", "159aa7d7-2fa0-4cc8-a708-3328d1d08eb5", "0774dbda-194f-439d-97e3-eec0e325fe5a",
				"1aaaed1f-9d23-47ef-8647-17b862becc27", "272085be-4887-498a-b7f6-85870fe93b40", "fb562584-ff7d-470b-a85a-0ee420a25850",
				"6ab15b0d-adf8-4c58-815a-194ba9c67e0d", "a7e9159f-af65-44a1-abe6-a86767a2f8f6", "9c5b448e-35c3-443a-bf13-65d9ac86fb19",
				"4a5d0bc1-5a43-4d0d-8b6a-4bf5e755846e", "f77cc578-bf01-4f79-8a6c-885e42d5ed37", "06a722ae-f0a1-4d94-80a0-6c2bc2b2597d",
				"de528431-8391-4c85-b89e-6fc6a77babba", "62648619-db60-42bc-bb98-5cf4bf2e06b4", "b09cab84-5327-415a-a076-3891a51aa211",
				"1a0b3314-2b2b-4fed-8a4a-55ce74a4ef19", "26bc0e0e-0b6c-487b-bc1e-823840cbaa52", "8ad11626-73ad-412c-8213-85ccbc5e180e",
				"fcf85764-8972-4e0a-9c27-1c3a3d4b7ae2", "f968b6cc-9e34-4891-b87b-7f78cc0f3aa5", "ed274237-32e2-4c3c-8606-7b44885e1ac1",
				"02cce80b-ecad-4cce-95f2-c5e3dd887390", "66a50383-baa3-4db9-871e-7ec151ac910f", "0a8476a5-b8bb-4a98-9a0a-64bc794e0c35",
				"458126d1-9230-48ba-992c-7f6287c35b26", "592186bd-6f25-4426-8ee1-fc78a7d98f56", "7f4a631f-f59d-4453-8c16-f023d050bbb9",
				"a81c33b8-b429-4879-bbe1-1adc65987a57", "59c299d8-faae-44d7-b751-424fb3077072", "055e72b4-c0fb-4a19-b945-baa499daf3e6",
				"53ba58ec-6917-4ae6-b91f-6a3eadcaf0a6", "3d005413-587a-42ca-be1a-97106d684861", "29b42c97-7f1e-4933-8286-18592a9845b6"},
			[]string{"Elektronik Shop", "Pakaian Shop", "Sepatu Shop",
				"Tas Shop", "Aksesoris Fashion Shop", "Hobi & Koleksi Shop",
				"Kesehatan Shop", "Makanan & Minuman Shop", "Perawatan & Kecantikan Shop",
				"Perlengkapan Dapur Shop", "Otomotif Shop", "Olahraga & Outdoor Shop",
				"Komputer Shop", "Handphone Shop", "Kamera Shop",
				"Pakaian Pria Shop", "Pakaian Wanita Shop", "Outfit Hangat Shop",
				"Sepatu Pria Shop", "Sepatu Wanita Shop", "Sepatu Anak Shop",
				"Tas Pria Shop", "Tas Wanita Shop", "Kacamata Shop",
				"Jam Tangan Shop", "Gitar Shop", "Buku Shop",
				"Obat Shop", "Masker Shop", "Makanan Shop",
				"Minuman Shop", "Face Wash Shop", "Sunscreen Shop",
				"Panci Shop", "Kompor Shop", "Helm Shop",
				"Spion Shop", "Raket Shop", "Bola Shop",
				"Laptop Shop", "Mouse Shop", "Keyboard Shop",
				"Webcam Shop", "Sweeter Shop", "Jaket Shop"},
			[]int{2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250,
				2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250,
				2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250, 2250},
			[]string{"98c1921e-b80e-40f3-9cba-fe8806097517", "0d389020-f229-461e-9202-5788961fbb81", "4bf503dc-689e-4b66-8401-3f133f1d585a"},
			"7950eca2-58d5-44f0-b873-22b23d8107da",
			"2787884621261326")},
		{Seeder: table.NewUserFaker(10000, 1, "M", []string{}, []string{}, []string{})},
		{Seeder: table.NewUserFaker(10000, 1, "F", []string{}, []string{}, []string{})},
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
